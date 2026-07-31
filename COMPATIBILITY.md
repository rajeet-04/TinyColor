# Compatibility Matrix

Target: this checkout's `mod.js` and `test.js` from `bgrins/TinyColor`.

## Current evidence

| Area | Status | Evidence | Notes |
|---|---|---|---|
| Oracle baseline | Ready | `mod.js`, `test.js`, Node 24.18.0 | JavaScript source is untouched. |
| Node adapter | Passing | `node tests/port/adapter.test.mjs` | Imports local `mod.js`; success/error responses are exclusive. |
| Go adapter | Passing | `go test ./...` from `src/` | JSONL schema shared with Node adapter. |
| HEX/RGB/name parsing | Passing | `compat/cases/parser-hex-rgb.jsonl` | Source-derived HEX/RGB/name/object corpus. |
| HSL/HSV parsing | Passing | `compat/cases/parser.jsonl` | Includes ratios, percentages, wrapping, and object precedence. |
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
node tests/port/adapter.test.mjs
Set-Location src; go test ./...; go vet ./...; Set-Location ..
node compat/run.mjs compat/cases/smoke.jsonl
```

This is only the fixed Phase 1 corpus; it is not a full TinyColor parity claim.

## Phase 2 parser slice result

On 2026-08-01, the HEX/RGB/name parser corpus passed with **26/26 cases** and
**0 mismatches**:

```powershell
$env:GOCACHE = 'R:\Code\TinyColor\.cache\go-build'
Set-Location src; go test ./...; go vet ./...; Set-Location ..
node compat/run.mjs compat/cases/parser-hex-rgb.jsonl
node compat/run.mjs compat/cases/smoke.jsonl
```

This records only the completed parser slice; the full HSL/HSV corpus remains
Phase 2 Plan 02 work.

## Phase 2 full parser result

`node compat/run.mjs compat/cases/parser.jsonl` passed with **23/23 cases** and
**0 mismatches** on 2026-08-01. The unchanged Phase 1 smoke corpus also passed
**9/9** with **0 mismatches**.
