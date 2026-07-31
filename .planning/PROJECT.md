# TinyColor Go Port

## Purpose

Create a fresh Go port of the local `bgrins/TinyColor` checkout with behavioral
parity demonstrated through unchanged-source differential testing.

## Success criteria

- The upstream JavaScript source is untouched.
- Go handles every behavior covered by the defined compatibility corpus.
- A reproducible command reports compatibility results and mismatch details.
- A documented Go API and CLI are usable from a clean checkout.
- CI, benchmarks, attribution, and known differences are present.

## Constraints

- Target source: local `mod.js` / `test.js`; MIT attribution retained.
- Target language: Go 1.26+; standard library first.
- Node is available as local oracle runtime; Deno is currently unavailable.

