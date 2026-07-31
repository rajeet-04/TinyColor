# GSD Debug Knowledge Base

Resolved debug sessions. Used by `gsd-debugger` to surface known-pattern hypotheses at the start of new investigations.

---

## fuzz-saturate-rounding — Saturate negative rounding differed at an RGB half boundary
- **Date:** 2026-08-01
- **Error patterns:** saturate, rounding, #400140, #202020, #212121, RGB 32, RGB 33, fuzz divergence
- **Root cause:** setHSL called a direct hslToRGB conversion, bypassing TinyColor's object-input percentage conversion and Bound01 truncation before RGB output rounding.
- **Fix:** Route setHSL through the existing hslColor/parser conversion and remove the duplicate direct converter.
- **Files changed:** src/tinycolor/color.go, src/tinycolor/operations_test.go, compat/cases/operations.jsonl
---
