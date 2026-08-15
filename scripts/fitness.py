#!/usr/bin/env python3
"""Fitness function executor for the gofr_layout repo.

Maps repo guardrails to external signals (build / test / vet / fmt / diff checks).
Each check is graded: "block" fails the whole run (exit 1), "warn" only reports.
This executor is the machine part of the Harness layer (see docs/fitness/).

Usage:
    python3 scripts/fitness.py           # human-readable report
    python3 scripts/fitness.py --json    # one JSON line per check, for scripts/loop.sh
"""
import argparse
import json
import os
import re
import subprocess
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
PROTO_DIR = os.path.join(ROOT, "proto")
SPECS_DIR = os.path.join(ROOT, "docs", "specs")
SPEC_STATUS_OK = {"draft", "active", "done", "archived"}


def run(cmd, cwd=ROOT, env=None):
    e = dict(os.environ)
    if env:
        e.update(env)
    p = subprocess.run(cmd, cwd=cwd, env=e, capture_output=True, text=True)
    return p.returncode, (p.stdout or "") + (p.stderr or "")


def merge_base(cwd):
    """Merge base against the default branch, or None on a repo with no commits."""
    for ref in ("origin/main", "origin/master", "main", "HEAD~1"):
        p = subprocess.run(["git", "merge-base", "HEAD", ref], cwd=cwd,
                           capture_output=True, text=True)
        if p.returncode == 0:
            return p.stdout.strip()
    return None


def changed_files(cwd=ROOT):
    """Files changed on HEAD relative to the merge base."""
    base = merge_base(cwd)
    if base is None:
        return []
    p = subprocess.run(["git", "diff", "--name-only", base + "...HEAD"], cwd=cwd,
                       capture_output=True, text=True)
    return p.stdout.splitlines()


def check_build():
    # go.work drives resolution (gofr/gokit/proto siblings); go.mod replace is the
    # GOWORK=off fallback covered by check_vet below.
    code, out = run(["go", "build", "./..."])
    return ("ok" if code == 0 else "fail", "build", out if code else "go build ./...")


def check_vet():
    code, out = run(["go", "vet", "./..."], env={"GOWORK": "off"})
    return ("ok" if code == 0 else "fail", "vet", out if code else "GOWORK=off go vet ./...")


def check_test_unit():
    # cmd/test is the integration binary (boots against a local DB); excluded here.
    code, out = run(["go", "test", "./internal/...", "./cmd/api",
                     "./cmd/consumer", "./cmd/script"])
    return ("ok" if code == 0 else "fail", "test-unit",
            out if code else "go test ./internal/... ./cmd/api ./cmd/consumer ./cmd/script")


def check_test_integration():
    code, out = run(["go", "test", "./cmd/test"])
    if code == 0:
        return ("ok", "test-integration", "go test ./cmd/test")
    return ("warn", "test-integration",
            "go test ./cmd/test (needs local DB)\n" + out)


def check_fmt():
    files = [f for f in changed_files() if f.endswith(".go")]
    if not files:
        return ("ok", "fmt", "gofmt (no changed .go files)")
    code, out = run(["gofmt", "-l"] + files)
    # gofmt -l exits 0 even when files need formatting; the list is the signal.
    if code == 0 and not out.strip():
        return ("ok", "fmt", "gofmt -l on changed files")
    return ("warn", "fmt", out or "gofmt exited %d" % code)


def check_wire_sync():
    """Contract: provider sets / wire.go changed -> wire_gen.go must be regenerated."""
    files = changed_files()
    sources = [f for f in files if f.endswith("wire.go") or f.endswith("wireProviderSet.go")]
    if not sources:
        return ("ok", "wire-sync", "no provider-set changes")
    gens = [f for f in files if f.endswith("wire_gen.go")]
    if not gens:
        return ("fail", "wire-sync",
                "provider sets changed without regenerated wire_gen.go: %s" % ", ".join(sources))
    return ("ok", "wire-sync", "wire_gen.go regenerated")


def check_proto_sync():
    """Contract (in-repo ./proto): a changed .proto must have fresh generated output.
    Generated *_pb.go are gitignored (regenerated via `cd proto && make all`), so git
    diff cannot see them; compare file mtimes on disk instead."""
    if not os.path.isdir(PROTO_DIR):
        return ("skip", "proto-sync", "proto/ not present")
    protos = [f for f in changed_files()
              if f.startswith("proto/") and f.endswith(".proto")]
    if not protos:
        return ("ok", "proto-sync", "no proto changes in range")
    stale = []
    for rel in protos:
        abs_proto = os.path.join(ROOT, rel)
        if not os.path.exists(abs_proto):
            continue
        proto_mtime = os.path.getmtime(abs_proto)
        d = os.path.dirname(abs_proto)
        gens = [f for f in os.listdir(d) if f.endswith(".pb.go")]
        if not any(os.path.getmtime(os.path.join(d, g)) >= proto_mtime for g in gens):
            stale.append(rel)
    if stale:
        return ("warn", "proto-sync",
                "proto changed without regenerated output (cd proto && make all): %s"
                % ", ".join(stale))
    return ("ok", "proto-sync", "generated output fresh for changed protos")


def _change_names():
    d = os.path.join(SPECS_DIR, "changes")
    if not os.path.isdir(d):
        return []
    return [x for x in sorted(os.listdir(d))
            if x != "archive" and os.path.isdir(os.path.join(d, x))]


def _spec_status(change):
    p = os.path.join(SPECS_DIR, "changes", change, ".spec.yaml")
    try:
        txt = open(p, encoding="utf-8").read()
    except OSError:
        return None
    m = re.search(r"^status:\s*(\S+)", txt, re.M)
    return m.group(1) if m else None


def _reqs(text):
    """Requirement summaries: (header, scenario_count, has_then)."""
    reqs = []
    cur = None
    for ln in text.splitlines():
        if ln.startswith("### Requirement:"):
            cur = {"line": ln, "scenarios": 0, "then": False}
            reqs.append(cur)
        elif ln.startswith("#### Scenario:"):
            if cur:
                cur["scenarios"] += 1
        elif ln.startswith("- **THEN"):
            if cur:
                cur["then"] = True
    return reqs


def check_spec():
    """Structural gate for spec artifacts. Drafts are lenient; others strict."""
    issues = []
    for change in _change_names():
        st = _spec_status(change)
        if st is None:
            issues.append("changes/%s: missing .spec.yaml" % change)
            continue
        if st not in SPEC_STATUS_OK:
            issues.append("changes/%s: illegal status %r" % (change, st))
        if st == "draft":
            continue
        p = os.path.join(SPECS_DIR, "changes", change, "spec.md")
        if not os.path.exists(p):
            issues.append("changes/%s: missing spec.md" % change)
            continue
        text = open(p, encoding="utf-8").read()
        if not re.search(r"^## (ADDED|MODIFIED|REMOVED) Requirements", text, re.M):
            issues.append("changes/%s/spec.md: no ADDED/MODIFIED/REMOVED Requirements" % change)
        for req in _reqs(text):
            if not re.search(r"### Requirement: R\d+", req["line"]):
                issues.append("changes/%s/spec.md: %r lacks R-id" % (change, req["line"].strip()))
            if req["scenarios"] == 0:
                issues.append("changes/%s/spec.md: %r has no Scenario" % (change, req["line"].strip()))
            elif not req["then"]:
                issues.append("changes/%s/spec.md: %r Scenario lacks THEN" % (change, req["line"].strip()))
    lib_dir = os.path.join(SPECS_DIR, "lib")
    if os.path.isdir(lib_dir):
        for fn in sorted(x for x in os.listdir(lib_dir) if x.endswith(".md")):
            text = open(os.path.join(lib_dir, fn), encoding="utf-8").read()
            if "## Requirements" not in text:
                issues.append("lib/%s: missing ## Requirements" % fn)
            for req in _reqs(text):
                if req["scenarios"] == 0:
                    issues.append("lib/%s: %r has no Scenario" % (fn, req["line"].strip()))
                elif not req["then"]:
                    issues.append("lib/%s: %r Scenario lacks THEN" % (fn, req["line"].strip()))
    if issues:
        return ("fail", "spec-validate", "\n".join(issues))
    return ("ok", "spec-validate", "docs/specs structure valid")


def check_spec_tasks():
    """Coverage: every requirement R# in a change spec has a task; task ids unique."""
    issues = []
    for change in _change_names():
        st = _spec_status(change)
        if st == "draft":
            continue
        spec_p = os.path.join(SPECS_DIR, "changes", change, "spec.md")
        tasks_p = os.path.join(SPECS_DIR, "changes", change, "tasks.md")
        if not os.path.exists(spec_p):
            continue
        spec_text = open(spec_p, encoding="utf-8").read()
        req_ids = set(re.findall(r"### Requirement: (R\d+)", spec_text))
        if not req_ids:
            continue
        if not os.path.exists(tasks_p):
            issues.append("changes/%s: missing tasks.md" % change)
            continue
        tasks_text = open(tasks_p, encoding="utf-8").read()
        missing = sorted(req_ids - set(re.findall(r"\bR\d+\b", tasks_text)))
        if missing:
            issues.append("changes/%s/tasks.md: requirements without tasks: %s"
                          % (change, ", ".join(missing)))
        task_ids = re.findall(r"^- \[ \] (T\d+)", tasks_text, re.M)
        dup = sorted({t for t in task_ids if task_ids.count(t) > 1})
        if dup:
            issues.append("changes/%s/tasks.md: duplicate task ids %s"
                          % (change, ", ".join(dup)))
    if issues:
        return ("fail", "spec-tasks", "\n".join(issues))
    return ("ok", "spec-tasks", "requirements -> tasks coverage ok")


CHECKS = [
    check_build,
    check_vet,
    check_test_unit,
    check_test_integration,
    check_fmt,
    check_wire_sync,
    check_proto_sync,
    check_spec,
    check_spec_tasks,
]

ICON = {"ok": "[ok]  ", "warn": "[warn]", "fail": "[FAIL]", "skip": "[skip]"}


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--json", action="store_true")
    args = ap.parse_args()

    results = [fn() for fn in CHECKS]
    for status, label, output in results:
        if args.json:
            print(json.dumps({"name": label, "status": status, "output": output},
                             ensure_ascii=False))
        else:
            first = output.splitlines()[0] if output else ""
            print("%s %-17s %s" % (ICON[status], label, first))

    failed = [r for r in results if r[0] == "fail"]
    if failed:
        print("FITNESS FAILED (block): %s" % ", ".join(r[1] for r in failed))
        sys.exit(1)
    print("fitness passed: %d checks" % len(results))


if __name__ == "__main__":
    main()
