#!/usr/bin/env python3
"""Archive a spec change: merge its spec.md delta into the capability library.

Usage:
    python3 scripts/spec-archive.py <change-name>
    # e.g. python3 scripts/spec-archive.py 2026-08-12-user-crud

Mirrors OpenSpec archive:
- ADDED/MODIFIED/REMOVED Requirements from changes/<name>/spec.md are applied to
  docs/specs/lib/<capability>.md (capability read from .spec.yaml).
- MODIFIED/REMOVED referencing a requirement id absent from the lib is a hard error.
- On success the change is marked archived and moved to changes/archive/<name>.
"""
import os
import re
import shutil
import sys

SPECS = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
                     "docs", "specs")
CHANGES = os.path.join(SPECS, "changes")
LIB = os.path.join(SPECS, "lib")


def read(path):
    with open(path, encoding="utf-8") as f:
        return f.read()


def write(path, text):
    with open(path, "w", encoding="utf-8") as f:
        f.write(text)


def req_id(line):
    m = re.search(r"R\d+", line)
    return m.group(0) if m else None


def parse_delta(text):
    """Return {'ADDED': [[lines...]], 'MODIFIED': [...], 'REMOVED': [...]}."""
    ops = {"ADDED": [], "MODIFIED": [], "REMOVED": []}
    current = None
    block = None
    for ln in text.splitlines():
        m = re.match(r"^## (ADDED|MODIFIED|REMOVED) Requirements", ln)
        if m:
            if block is not None and current is not None:
                ops[current].append(block)
                block = None
            current = m.group(1)
            continue
        if current is None:
            continue
        if ln.startswith("### Requirement:"):
            if block is not None:
                ops[current].append(block)
            block = [ln]
        elif block is not None:
            block.append(ln)
    if block is not None and current is not None:
        ops[current].append(block)
    return ops


def find_block(lines, rid):
    """(start, end) of the requirement block with id rid, else None."""
    for i, ln in enumerate(lines):
        if ln.startswith("### Requirement:") and req_id(ln) == rid:
            j = i + 1
            while j < len(lines) and not lines[j].startswith("### Requirement:") \
                    and not lines[j].startswith("## "):
                j += 1
            return i, j
    return None


def insert_added(lines, block):
    try:
        h = next(i for i, ln in enumerate(lines) if ln == "## Requirements")
    except StopIteration:
        if lines and lines[-1] != "":
            lines.append("")
        lines.append("## Requirements")
        h = len(lines) - 1
    end = len(lines)
    for j in range(h + 1, len(lines)):
        if lines[j].startswith("## ") and lines[j] != "## Requirements":
            end = j
            break
    lines[end:end] = [""] + block
    return lines


def apply(lib_text, ops):
    lines = lib_text.splitlines()
    for op, blocks in ops.items():
        for block in blocks:
            rid = req_id(block[0])
            if rid is None:
                raise SystemExit("delta requirement without R-id: %r" % block[0])
            span = find_block(lines, rid)
            if op == "ADDED":
                if span is not None:
                    raise SystemExit("ADDED %s already exists in lib" % rid)
                lines = insert_added(lines, block)
            elif op == "MODIFIED":
                if span is None:
                    raise SystemExit("MODIFIED %s not found in lib" % rid)
                lines[span[0]:span[1]] = block
            elif op == "REMOVED":
                if span is None:
                    raise SystemExit("REMOVED %s not found in lib" % rid)
                del lines[span[0]:span[1]]
    return "\n".join(lines).rstrip() + "\n"


def main():
    if len(sys.argv) != 2:
        print(__doc__)
        sys.exit(1)
    name = sys.argv[1]
    change = os.path.join(CHANGES, name)
    if not os.path.isdir(change):
        sys.exit("change not found: %s" % change)
    cfg_path = os.path.join(change, ".spec.yaml")
    spec_path = os.path.join(change, "spec.md")
    if not (os.path.exists(cfg_path) and os.path.exists(spec_path)):
        sys.exit("%s: missing .spec.yaml or spec.md" % name)
    cfg = read(cfg_path)
    m = re.search(r"^capability:\s*(\S+)", cfg, re.M)
    if m is None:
        sys.exit("%s: .spec.yaml missing capability" % name)
    st = re.search(r"^status:\s*(\S+)", cfg, re.M)
    if st is None or st.group(1) == "archived":
        sys.exit("%s: already archived or no status" % name)
    lib_path = os.path.join(LIB, m.group(1) + ".md")
    if not os.path.exists(lib_path):
        os.makedirs(LIB, exist_ok=True)
        write(lib_path, "# %s 能力库\n\n## Requirements\n" % m.group(1))
    ops = parse_delta(read(spec_path))
    write(lib_path, apply(read(lib_path), ops))
    write(cfg_path, re.sub(r"^status:\s*\S+", "status: archived", cfg, flags=re.M))
    archive_dir = os.path.join(CHANGES, "archive")
    os.makedirs(archive_dir, exist_ok=True)
    shutil.move(change, os.path.join(archive_dir, name))
    summary = "; ".join(
        "%s %s" % (op, [req_id(b[0]) for b in blocks])
        for op, blocks in ops.items() if blocks)
    print("archived %s -> %s: %s" % (name, lib_path, summary))


if __name__ == "__main__":
    main()
