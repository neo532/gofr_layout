#!/usr/bin/env bash
# One verification round of the edit loop: run fitness, record state, stop on failure.
# The agent edits between rounds; this script owns "verify + record" so that:
#   - state lives outside the model (written to .claude/state/loop.json);
#   - a failing round stops the loop instead of drifting forward.
# Exit 0 = converged (all block checks green), 1 = repair needed.
set -u

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STATE_DIR="${ROOT}/.claude/state"
STATE_FILE="${STATE_DIR}/loop.json"
cd "${ROOT}"

# $1: "0" converged / "1" failed
record() {
  python3 - "${STATE_FILE}" "$1" <<'PY'
import glob, json, os, re, sys, datetime
path, ok = sys.argv[1], sys.argv[2] == "0"
try:
    with open(path) as f:
        state = json.load(f)
except Exception:
    state = {}
state["last_round"] = {
    "at": datetime.datetime.now().isoformat(timespec="seconds"),
    "ok": ok,
}
active = []
for y in glob.glob("docs/specs/changes/*/.spec.yaml"):
    try:
        txt = open(y, encoding="utf-8").read()
        m = re.search(r"^status:\s*(\S+)", txt, re.M)
        if m and m.group(1) == "active":
            active.append(os.path.basename(os.path.dirname(y)))
    except Exception:
        pass
state["active_specs"] = active
with open(path, "w") as f:
    json.dump(state, f, indent=2, ensure_ascii=False)
PY
}

mkdir -p "${STATE_DIR}"
echo "==> loop: running fitness round"
if ! python3 scripts/fitness.py --json > "${STATE_DIR}/fitness.out.json"; then
  record 1
  echo "==> loop: fitness FAILED — stop and repair. Inspect ${STATE_DIR}/fitness.out.json"
  exit 1
fi
record 0
echo "==> loop: converged — all block checks green at $(date +%H:%M:%S)"
