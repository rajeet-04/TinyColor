.PHONY: build test verify fmt-check hashes fuzz bench clean

GO_CACHE ?= $(CURDIR)/.cache/go-build
GO = GOCACHE=$(GO_CACHE) go

build:
	$(GO) -C src build -o ../bin/tinycolor ./cmd/tinycolor-compat

test:
	node tests/port/adapter.test.mjs
	$(GO) -C src test ./...
	node compat/run.mjs compat/cases/smoke.jsonl
	node compat/run.mjs compat/cases/parser-hex-rgb.jsonl
	node compat/run.mjs compat/cases/parser.jsonl
	node compat/run.mjs compat/cases/conversion.jsonl
	node compat/run.mjs compat/cases/operations.jsonl

verify: fmt-check hashes test
	$(GO) -C src vet ./...

fmt-check:
	node -e "const { readdirSync } = require('node:fs'); const { join } = require('node:path'); const { execFileSync } = require('node:child_process'); const files = []; const walk = (dir) => readdirSync(dir, { withFileTypes: true }).forEach((entry) => entry.isDirectory() ? walk(join(dir, entry.name)) : entry.name.endsWith('.go') && files.push(join(dir, entry.name))); walk('src'); const output = execFileSync('gofmt', ['-l', ...files], { encoding: 'utf8' }).trim(); if (output) { console.error(output); process.exit(1); }"

hashes:
	node tests/original/verify.mjs

fuzz:
	node fuzz/harness.mjs --duration 60 --seed 20260801

bench:
	node bench/run.mjs

clean:
	node -e "const fs = require('node:fs'); for (const file of ['bin/tinycolor', 'bin/tinycolor.exe']) fs.rmSync(file, { force: true });"
