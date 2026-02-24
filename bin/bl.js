#!/usr/bin/env node
const { execFileSync } = require("child_process");
const path = require("path");

const binary = path.join(__dirname, "bl-binary");

try {
  execFileSync(binary, process.argv.slice(2), { stdio: "inherit" });
} catch (e) {
  if (e.status !== undefined) {
    process.exit(e.status);
  }
  console.error("bl-cli: バイナリが見つかりません。再インストールしてください: npm install -g @kimmaru10/bl-cli");
  process.exit(1);
}
