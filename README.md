# bl-cli

> Backlog をターミナルから操作する CLI ツール

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**bl** は [Nulab Backlog](https://backlog.com/) の課題管理をコマンドラインから行うための CLI ツールです。  
GitHub CLI (`gh`) にインスパイアされた直感的なコマンド体系で、ブラウザを開かずに課題の作成・更新・コメントができます。

## デモ

```bash
# 課題を作成
$ bl issue create --summary "ログイン機能の実装" --type "タスク" --priority "中"
PROJ-123

# ブランチ名から課題を自動推測して操作
$ git checkout -b feature/PROJ-123-add-login
$ bl issue edit --status "処理中"
✔ PROJ-123 のステータスを「処理中」に変更しました

# コメントを追加
$ bl issue comment --body "実装完了しました。レビューお願いします"
✔ PROJ-123 にコメントを追加しました

# ブラウザで課題を確認
$ bl issue view --web
```

## 特徴

- **gh ライクな操作感** — `bl issue list`、`bl issue create` など、GitHub CLI に慣れていればすぐ使える
- **ブランチ名から課題キーを自動推測** — `feature/PROJ-123-xxx` ブランチにいれば、課題キーの入力を省略できる
- **インタラクティブ UI** — プロジェクトや担当者をリストから選択。課題キーを覚えていなくても操作可能
- **`--web` でブラウザ連携** — ターミナルからワンコマンドで Backlog の Web UI を開ける
- **シングルバイナリ** — Go 製。依存なしでインストール・配布が簡単

## インストール

### Homebrew

```bash
brew tap <org>/bl-cli
brew install bl-cli
```

### GitHub Releases

[Releases](https://github.com/<org>/bl-cli/releases) からお使いの OS に合ったバイナリをダウンロードしてください。

### Go

```bash
go install github.com/<org>/bl-cli@latest
```

## セットアップ

```bash
# Backlog の API キーとスペース URL を設定
bl auth login
```

API キーは Backlog の「個人設定 > API」から発行できます。

```bash
# デフォルトプロジェクトを設定（インタラクティブに選択）
bl project set
```

## 使い方

### 課題の一覧

```bash
# 自分に割り当てられた未完了の課題
bl issue list

# ステータスで絞り込み
bl issue list --status "処理中"

# マイルストーンで絞り込み
bl issue list --milestone "v1.0"
```

### 課題の詳細

```bash
# 課題キーを指定して表示
bl issue view PROJ-123

# ブランチ名から自動推測
bl issue view

# ブラウザで開く
bl issue view --web
```

### 課題の作成

```bash
# インタラクティブに作成（種別・優先度・担当者をリストから選択）
bl issue create

# オプション指定で作成
bl issue create --summary "バグ修正" --type "バグ" --priority "高" --assignee "yamada"
```

### 課題の更新

```bash
# ステータス変更
bl issue edit PROJ-123 --status "処理中"

# 期日変更
bl issue edit PROJ-123 --due-date "2026-03-31"

# 担当者変更
bl issue edit PROJ-123 --assignee "yamada"

# ブランチ名から推測して更新
bl issue edit --status "完了"
```

### コメント

```bash
# インラインでコメント追加
bl issue comment PROJ-123 --body "対応しました"

# エディタを起動してコメント入力
bl issue comment PROJ-123

# コメント一覧
bl issue comment list PROJ-123
```

### ブランチ名からの課題キー自動推測

git ブランチ名に課題キーが含まれている場合、自動的に抽出します。

```
feature/PROJ-123-add-login    → PROJ-123
fix/PROJ-456-bugfix           → PROJ-456
hotfix/PROJ-789               → PROJ-789
```

ブランチにいる状態で課題キーを省略すると、自動推測が働きます。

```bash
$ git checkout feature/PROJ-123-add-login
$ bl issue view              # PROJ-123 の詳細を表示
$ bl issue edit --status "処理中"  # PROJ-123 を更新
$ bl issue comment --body "完了"   # PROJ-123 にコメント
```

## 設定

設定ファイルは `~/.config/bl/config.yaml` に保存されます。

```yaml
space_url: https://myteam.backlog.com
api_key: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
default_project: MYPROJ
```

## コマンド一覧

| コマンド | 説明 |
|---------|------|
| `bl auth login` | 認証情報を設定 |
| `bl auth logout` | 認証情報を削除 |
| `bl auth status` | 認証状態を確認 |
| `bl project list` | プロジェクト一覧 |
| `bl project set` | デフォルトプロジェクトを設定 |
| `bl project current` | 現在のデフォルトプロジェクトを表示 |
| `bl issue list` | 課題一覧 |
| `bl issue view` | 課題の詳細を表示 |
| `bl issue create` | 課題を作成 |
| `bl issue edit` | 課題を更新 |
| `bl issue comment` | コメントを追加 |
| `bl issue comment list` | コメント一覧 |

## 開発

```bash
# クローン
git clone https://github.com/<org>/bl-cli.git
cd bl-cli

# ビルド
go build -o bl .

# 実行
./bl --help
```

## ライセンス

[MIT](LICENSE)
