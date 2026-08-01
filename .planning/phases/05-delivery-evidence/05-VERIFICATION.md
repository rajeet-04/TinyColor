---
phase: 5
status: locally-verified
verified: 2026-08-01
requirements: [DEL-01, QLT-03]
---

# Phase 5 verification

## Goal

A judge can build and use one runnable Go artifact, reproduce exact parity
evidence, inspect honest performance/safety numbers, and distinguish observed
local evidence from external actions.

## Roadmap success criteria

| Criterion | Status | Evidence |
|---|---|---|
| CI runs format, vet, Go tests, and differential checks | Configured; exact final run pending | `.github/workflows/port.yml`, `Makefile`; local constituent commands passed. |
| CLI supports parse, convert, lighten, palette, contrast, and JSON | Local pass | `make build` equivalent plus all five JSON invocations exited zero; implementation under `src/cmd/tinycolor-compat`. |
| Benchmarks, docs, attribution, ownership, and known differences are complete | Local pass | `bench/results.json`, `bench/methodology.md`, `README.md`, `LICENSE`, `DECISIONS.md`, `COMPATIBILITY.md`, `docs/TEAM-OWNERSHIP.md`. |

## DEL-01 evidence

- One-command build is defined by `make build`; direct Windows build produced
  `bin/tinycolor.exe` and all five public commands ran successfully.
- `tests/original/manifest.sha256` verified all three immutable oracle files.
- Five fixed corpora passed 164/164 with zero mismatches.
- `fuzz/log.txt` validates a 60.012-second, seed-20260801 run with 1,091,630
  cases and zero divergences; a fresh one-second smoke passed 11,103 cases.
- `bench/results.json` contains 20 cold starts, 1,000 persistent requests, and
  both implementations' startup p99, latency p99, throughput, and peak RSS.
- `docs/DEMO.md` is ready; recording/upload and its URL remain human actions.

## QLT-03 evidence

- 15 Node verification, fuzz, and benchmark tests passed.
- All Go packages passed tests and vet; package coverage is recorded in
  `COMPATIBILITY.md` without claiming JavaScript coverage parity.
- `rg -n '\bunsafe\b' src -g '*.go'` returned zero source occurrences.
- Deno source suite passed 45, failed 0, ignored 1.
- `git diff --check` passed, source diff was empty, and benchmark workload hash
  matched the committed normalized SHA-256.

## External checkpoint

- Public repository: pending unauthenticated verification.
- Exact final-commit GitHub Actions run: pending push and observation.
- Demo video: not supplied.

Phase 5 remains open until public clone access and the exact successful CI run
are recorded. The missing video URL remains explicitly unclaimed.
