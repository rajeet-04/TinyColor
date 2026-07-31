# Project State

**Updated:** 2026-08-01

- Planning initialized from the local `bgrins/TinyColor` checkout.
- JavaScript source and tests are the immutable oracle.
- Node 24.18.0 and Go 1.26.1 are available locally.
- `deno` is not installed; no original-source test pass is claimed.
- Phase 1 completed in commit `78423e3`: Go/Node JSONL runners, protocol tests,
  and a fixed corpus passed 9/9 with zero mismatches. This is not whole-library
  parity.
- Phase 2 completed locally: the normalized parser covers HEX, RGB(A), HSL(A),
  HSV(A), names, objects, invalid input, and FromRatio. Differential evidence:
  26/26 HEX/RGB/name cases, 23/23 complete parser cases, and 9/9 smoke cases.
- The Go module lives in `src/`; `tests/original/manifest.sha256` pins the
  unmodified root JavaScript oracle.
- Phase 3 completed locally: conversion/representation and analysis behavior
  passed a 35/35 exact corpus; all Phase 1–3 corpora are zero-mismatch. The
  compatibility protocol now preserves successful JSON `false` results.
- Phase 4 completed locally: modifiers, Mix, WCAG readability, and palette
  operations passed the full gate. Exact differential evidence is 9/9 smoke,
  26/26 HEX/RGB, 23/23 parser, 35/35 conversion, and 58/58 operations with
  zero mismatches. `deno` remains unavailable, so no original-source suite
  pass is claimed.
