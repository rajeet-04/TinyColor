# Five-minute delivery demo

1. Show `git status --short` and `node tests/original/verify.mjs` to establish
   the immutable JavaScript oracle.
2. Build the CLI: `go -C src build -o ../bin/tinycolor.exe ./cmd/tinycolor-compat`.
3. Demonstrate the five commands: `parse`, `convert`, `lighten`, `palette`, and
   `contrast`, each with `--json`.
4. Run `node compat/run.mjs compat/cases/operations.jsonl` and show 71/71
   passing cases.
5. Show `fuzz/log.txt`, then reproduce a short session with
   `node fuzz/harness.mjs --duration 1 --seed 20260801`.
6. Show `bench/results.json` and `bench/methodology.md`, explaining that the
   figures are same-host observations.
7. Show `COMPATIBILITY.md` and `DECISIONS.md`, including known external gaps.

Manual, unverified submission actions: record/upload the live video, verify
public clone access, and link a successful CI run for the final commit.
