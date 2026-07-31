# Compatibility Matrix

Target: this checkout's `mod.js` and `test.js` from `bgrins/TinyColor`.

## Current evidence

| Area | Status | Evidence | Notes |
|---|---|---|---|
| Oracle baseline | Ready | `mod.js`, `test.js`, Node 24.18.0 | JavaScript source is untouched. |
| Node adapter | Passing | `node compat/js-runner.test.mjs` | Imports local `mod.js`; success/error responses are exclusive. |
| Go adapter | Passing | `go test ./...` from `go/` | JSONL schema shared with Node adapter. |
| HEX/RGB/name parsing | Planned | Wave 1 | Include invalid and permissive syntax. |
| HSL/HSV parsing | Planned | Wave 1 | Include ratios, percentages, and wrapping. |
| Conversion/formatting | Planned | Wave 2 | Exact strings and rounding. |
| Analysis/readability | Planned | Wave 2–3 | Preserve WCAG defaults. |
| Manipulation/palettes | Planned | Wave 3 | Preserve defaults and ordering. |
| CLI/CI/benchmarks | Planned | Wave 4 | No result claimed yet. |
| Original Deno suite | Blocked locally | `deno` unavailable | Run when Deno is installed or in CI. |

## Mismatch record format

```text
Case: <stable case id>
Operation: <adapter operation>
Input: <JSON>
JavaScript: <JSON or error>
Go: <JSON or error>
Owner: <A/B/C/D>
Status: open | fixed | accepted-difference
Reason: <required for accepted difference>
```

## Phase 1 smoke result

On 2026-08-01, the fixed Phase 1 corpus passed with **9/9 cases** and **0
mismatches**:

```powershell
$env:GOCACHE = 'R:\Code\TinyColor\.cache\go-build'
node compat/js-runner.test.mjs
Set-Location go; go test ./...; go vet ./...; Set-Location ..
node compat/run.mjs compat/cases/smoke.jsonl
```

This is only the fixed Phase 1 corpus; it is not a full TinyColor parity claim.
