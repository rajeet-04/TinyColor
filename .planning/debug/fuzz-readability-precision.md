---
status: investigating
trigger: "Investigate and fix next deterministic Phase 5 fuzz parity mismatch: readability for {r:127,g:64,b:15,a:0} against '#80007f' returns 1.188086751976723 in JavaScript and 1.1880867519767233 in Go."
created: 2026-08-01T00:00:00+05:30
updated: 2026-08-01T00:04:00+05:30
---

## Current Focus

hypothesis: confirmed: V8 Math.pow and Go math.Pow use different algorithms whose final rounding differs across the finite 8-bit luminance channel domain
test: inspect the primary V8 implementation and compare all 245 nonlinear channel values to select the smallest exact shared-method implementation
expecting: algorithm evidence plus the bounded mismatch distribution will show whether a small compatible calculation exists or an exact transfer table is required
next_action: inspect V8 ieee754 pow source and evaluate exact implementation options

## Symptoms

expected: JavaScript and Go exact JSON response match for readability; JavaScript returns 1.188086751976723.
actual: Go returns 1.1880867519767233; broad seed-1 one-second fuzz run has 388 divergences and fuzz-32 is first.
errors: exact floating-point JSON mismatch for operation readability
reproduction: exact request is in .superpowers/sdd/05-02-PLAN/task-1-report.md; run through both adapters and typed Readability with input {r:127,g:64,b:15,a:0}, args.other '#80007f'.
started: uncovered after the fuzz saturation root-cause fix.

## Eliminated

## Evidence

- timestamp: 2026-08-01T00:01:00+05:30
  checked: repository status and protected fuzz files
  found: fuzz/README.md is modified and fuzz/harness.mjs plus fuzz/harness.test.mjs are untracked pre-existing work
  implication: preserve those files without editing, staging, reverting, or narrowing them
- timestamp: 2026-08-01T00:01:00+05:30
  checked: debug knowledge base
  found: only prior saturation rounding case exists and its keywords/root cause do not match readability precision
  implication: no known-pattern candidate applies
- timestamp: 2026-08-01T00:01:00+05:30
  checked: complete JavaScript getLuminance/readability and Go Color.Luminance/Readability implementations plus callers
  found: both round normalized RGB before luminance and use the same WCAG constants; Go caches the two luminances once while JavaScript recomputes them for Max and Min; IsReadable and MostReadable both call shared Readability
  implication: parsing and exact arithmetic intermediates must be measured before changing the shared method
- timestamp: 2026-08-01T00:02:00+05:30
  checked: exact requests through both adapters
  found: normalized RGB inspections are identical; first luminance is identical at 0.08213307169664993, but #80007f luminance is 0.06121500300950977 in JavaScript and 0.061215003009509765 in Go
  implication: parsing, alpha, and final ratio ordering are ruled out; the first divergence is inside shared Color.Luminance for integer RGB 128,0,127
- timestamp: 2026-08-01T00:03:00+05:30
  checked: new exact typed regression and exact operations JSONL row before production changes
  found: typed test fails with 1.1880867519767233 and corpus has exactly one mismatch at readability-object-precision
  implication: RED is confirmed and independently covers both public Go API and adapter parity
- timestamp: 2026-08-01T00:04:00+05:30
  checked: channel-level gamma expansion raw float bits for #80007f
  found: channel 128 has identical normalized/base input but V8 Math.pow yields bits 3fcba1511e3e632d while Go math.Pow yields 3fcba1511e3e632c; channel 127 matches at 3fcb2a60a1263b0a
  implication: the root cause is Go math.Pow not reproducing the JavaScript oracle's power rounding, not weighted luminance summation
- timestamp: 2026-08-01T00:04:00+05:30
  checked: all nonlinear integer channels 11 through 255
  found: 93 of 245 channel powers differ, with signed deltas in both directions and two cases differing by two ULPs
  implication: unconditional Nextafter or decimal/output rounding would be incorrect; the fix must reproduce the bounded transfer function or V8 algorithm

## Resolution

root_cause:
fix:
verification:
files_changed: []
