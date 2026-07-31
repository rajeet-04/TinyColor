.PHONY: build test

build:
	go -C src build -o ../bin/tinycolor-compat ./cmd/tinycolor-compat

test:
	node tests/port/adapter.test.mjs
	go -C src test ./...
	go -C src vet ./...
	node compat/run.mjs compat/cases/parser-hex-rgb.jsonl
	node compat/run.mjs compat/cases/parser.jsonl
	node compat/run.mjs compat/cases/smoke.jsonl
