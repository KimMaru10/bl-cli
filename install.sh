#!/bin/bash
#
# bl-cli インストールスクリプト
#
# 使い方:
#   curl -fsSL https://raw.githubusercontent.com/KimMaru10/bl-cli/main/install.sh | bash
#
set -e

REPO="KimMaru10/bl-cli"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="bl"

# カラー出力
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

info()  { printf "${GREEN}==>${NC} %s\n" "$1"; }
warn()  { printf "${YELLOW}警告:${NC} %s\n" "$1"; }
error() { printf "${RED}エラー:${NC} %s\n" "$1" >&2; }

# OS とアーキテクチャを判定
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    error "サポートされていないアーキテクチャです: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  darwin) EXT="zip" ;;
  linux)  EXT="tar.gz" ;;
  *)
    error "サポートされていない OS です: $OS"
    exit 1
    ;;
esac

# 最新バージョンを取得
info "最新バージョンを取得しています..."
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | head -1 \
  | sed -E 's/.*"v?([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
  error "バージョンの取得に失敗しました"
  exit 1
fi

info "bl-cli v${VERSION} をインストールします"

# ダウンロード
FILENAME="bl-cli_${VERSION}_${OS}_${ARCH}.${EXT}"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}"
TMP_DIR=$(mktemp -d)
trap "rm -rf ${TMP_DIR}" EXIT

info "ダウンロード中: ${FILENAME}"
if ! curl -fsSL -o "${TMP_DIR}/${FILENAME}" "${URL}"; then
  error "ダウンロードに失敗しました: ${URL}"
  exit 1
fi

# 解凍
info "展開中..."
cd "${TMP_DIR}"
if [ "$EXT" = "zip" ]; then
  unzip -q "${FILENAME}"
else
  tar -xzf "${FILENAME}"
fi

if [ ! -f "${BINARY_NAME}" ]; then
  error "バイナリが見つかりません: ${BINARY_NAME}"
  exit 1
fi

# インストール
info "インストール先: ${INSTALL_DIR}/${BINARY_NAME}"
if [ -w "${INSTALL_DIR}" ]; then
  mv "${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
  chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
else
  warn "管理者権限が必要です。Mac のログインパスワードを入力してください"
  sudo mv "${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
  sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
fi

# macOS Gatekeeper 回避
if [ "$OS" = "darwin" ]; then
  if [ -w "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    xattr -d com.apple.quarantine "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null || true
  else
    sudo xattr -d com.apple.quarantine "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null || true
  fi
fi

# 確認
if ! command -v bl >/dev/null 2>&1; then
  error "インストールに失敗しました。${INSTALL_DIR} が PATH に含まれているか確認してください"
  exit 1
fi

echo ""
info "✔ インストール完了！ (bl v${VERSION})"
echo ""
echo "次のステップ:"
echo "  1. bl auth login     # Backlog に認証"
echo "  2. bl project set    # デフォルトプロジェクトを設定"
echo "  3. bl mcp setup      # Claude Desktop に連携を登録"
echo ""
