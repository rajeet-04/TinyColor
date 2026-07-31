# TinyColor Go port

This repository ports the behavior of the pinned `bgrins/TinyColor` checkout
to Go. Build the compatibility CLI with one command:

```powershell
make build
```

Run the current parity checks with `make test`. The immutable JavaScript oracle
stays at the root and its kickoff hashes are recorded in
`tests/original/manifest.sha256`.
