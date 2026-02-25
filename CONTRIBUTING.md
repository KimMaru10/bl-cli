# Contributing to bl-cli

bl-cli へのコントリビューションを歓迎します！

## 開発環境のセットアップ

### 必要なもの

- Go 1.22 以上
- Git

### ビルド手順

```bash
git clone https://github.com/KimMaru10/bl-cli.git
cd bl-cli
go build -o bl .
./bl --help
```

### テスト用の Backlog 設定

```bash
./bl auth login
# → スペース URL と API キーを入力
./bl project set
# → デフォルトプロジェクトを選択
```

## コントリビューションの流れ

1. **Issue を確認する**: 既存の Issue を確認し、取り組みたいものがあればコメントしてください
2. **Fork する**: リポジトリを Fork してください
3. **ブランチを作成する**: `feature/xxx` や `fix/xxx` のような名前でブランチを作成してください
4. **実装する**: コードを変更してください
5. **ビルドを確認する**: `go build -o bl .` が通ることを確認してください
6. **PR を作成する**: main ブランチに向けて PR を作成してください

## コーディング規約

### ディレクトリ構成

```
cmd/           CLI コマンド（Cobra）
internal/api/  Backlog API クライアント
internal/tui/  TUI コンポーネント（Bubble Tea）
internal/git/  Git 連携
```

### コミットメッセージ

```
[#Issue番号]タイプ：内容
```

タイプ:
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメント
- `chore`: メンテナンス
- `refactor`: リファクタリング

例: `[#21]feat：npm 配布対応`

### コード

- エラーメッセージは日本語
- コメントは英語
- `go fmt` でフォーマットしてください
- 新しいコマンドは `NewXxxCmd()` ファクトリパターンに従ってください

## Issue の作成

### バグ報告

- 再現手順を記載してください
- `bl --version` の出力を含めてください
- OS とアーキテクチャを記載してください

### 機能要望

- ユースケースを記載してください
- 既存のコマンドとの関連があれば記載してください

## 質問

Issue で「question」ラベルをつけて質問してください。
