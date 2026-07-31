import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { verifyManifest } from "./verify.mjs";

const sha256 = (value) => createHash("sha256").update(value).digest("hex");

test("verifyManifest accepts matching files and reports a changed path", () => {
  const root = mkdtempSync(join(tmpdir(), "tinycolor-verify-"));
  try {
    writeFileSync(join(root, "first.txt"), "first");
    writeFileSync(join(root, "second.txt"), "second");
    const manifest = join(root, "manifest.sha256");
    writeFileSync(manifest, `${sha256("first")}  first.txt\n${sha256("second")}  second.txt\n`);

    assert.deepEqual(verifyManifest(manifest, root), []);

    writeFileSync(join(root, "second.txt"), "changed");
    assert.deepEqual(verifyManifest(manifest, root), ["mismatch: second.txt"]);
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});

test("verifyManifest reports a malformed manifest line", () => {
  const root = mkdtempSync(join(tmpdir(), "tinycolor-verify-"));
  try {
    const manifest = join(root, "manifest.sha256");
    writeFileSync(manifest, "not a manifest entry\n");

    assert.deepEqual(verifyManifest(manifest, root), ["invalid: not a manifest entry"]);
  } finally {
    rmSync(root, { recursive: true, force: true });
  }
});
