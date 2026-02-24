const https = require("https");
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");
const os = require("os");

const REPO = "KimMaru10/bl-cli";
const BIN_DIR = path.join(__dirname, "bin");
const BIN_PATH = path.join(BIN_DIR, "bl-binary");

function getPlatform() {
  const platform = os.platform();
  if (platform === "darwin") return "darwin";
  if (platform === "linux") return "linux";
  throw new Error(`Unsupported platform: ${platform}`);
}

function getArch() {
  const arch = os.arch();
  if (arch === "x64") return "amd64";
  if (arch === "arm64") return "arm64";
  throw new Error(`Unsupported architecture: ${arch}`);
}

function getDownloadUrl(version, platform, arch) {
  const ext = platform === "darwin" ? "zip" : "tar.gz";
  const tag = version.startsWith("v") ? version : `v${version}`;
  const ver = tag.replace(/^v/, "");
  return `https://github.com/${REPO}/releases/download/${tag}/bl-cli_${ver}_${platform}_${arch}.${ext}`;
}

function httpsGet(url) {
  return new Promise((resolve, reject) => {
    https.get(url, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        httpsGet(res.headers.location).then(resolve).catch(reject);
        return;
      }
      if (res.statusCode !== 200) {
        reject(new Error(`HTTP ${res.statusCode} for ${url}`));
        return;
      }
      const chunks = [];
      res.on("data", (chunk) => chunks.push(chunk));
      res.on("end", () => resolve(Buffer.concat(chunks)));
      res.on("error", reject);
    }).on("error", reject);
  });
}

function getVersion() {
  const pkg = require("./package.json");
  return pkg.version;
}

async function install() {
  const platform = getPlatform();
  const arch = getArch();
  const version = getVersion();
  const url = getDownloadUrl(version, platform, arch);
  const ext = platform === "darwin" ? "zip" : "tar.gz";

  console.log(`Downloading bl ${version} for ${platform}/${arch}...`);

  const data = await httpsGet(url);

  if (!fs.existsSync(BIN_DIR)) {
    fs.mkdirSync(BIN_DIR, { recursive: true });
  }

  const tmpFile = path.join(os.tmpdir(), `bl-cli.${ext}`);
  fs.writeFileSync(tmpFile, data);

  try {
    if (ext === "zip") {
      execSync(`unzip -o -j "${tmpFile}" bl -d "${BIN_DIR}"`, { stdio: "pipe" });
      fs.renameSync(path.join(BIN_DIR, "bl"), BIN_PATH);
    } else {
      execSync(`tar -xzf "${tmpFile}" -C "${BIN_DIR}" bl`, { stdio: "pipe" });
      fs.renameSync(path.join(BIN_DIR, "bl"), BIN_PATH);
    }

    fs.chmodSync(BIN_PATH, 0o755);
    console.log("bl installed successfully.");
  } finally {
    fs.unlinkSync(tmpFile);
  }
}

install().catch((err) => {
  console.error("Failed to install bl:", err.message);
  process.exit(1);
});
