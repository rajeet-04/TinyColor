# tinycolor Go package

From the module directory, add the package to a Go program and import
`github.com/rajeet-04/tinycolor-go/tinycolor`.

```go
color, err := tinycolor.FromCompat("rgba(255, 0, 0, .5)", false)
if err != nil { panic(err) }
fmt.Println(color.ToHex8String()) // #FF000080
fmt.Println(color.ToHSLString())  // hsla(0, 100%, 50%, 0.5)

lighter := color.Clone()
lighter.Lighten(10)
fmt.Println(lighter.ToRGBString())
```

For compatibility testing from the repository root, run each fixed corpus with
`node compat/run.mjs compat/cases/operations.jsonl`, or run the full local gate
with `make verify`. This repository does not use npm scripts; the Node command
is the equivalent differential-test harness.
