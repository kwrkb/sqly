[English](README.md)

# asql

**データ観察**のための軽量 TUI SQL クライアント — 生データを素早く見て、並べ替えて、探索し、違和感や仮説に気づくためのツール。[Bubble Tea](https://github.com/charmbracelet/bubbletea) で構築。SQLite、MySQL、PostgreSQL をサポート。

## デモ

![asql デモ](docs/demo.gif)

## インストール

[GitHub Releases](https://github.com/kwrkb/asql/releases) からビルド済みバイナリをダウンロードできます。

Go でインストール:

```bash
go install github.com/kwrkb/asql@latest
```

またはソースからビルド:

```bash
git clone https://github.com/kwrkb/asql
cd asql
go build -o asql .
```

## 使い方

```bash
# SQLite
asql <sqlite-ファイルパス>

# MySQL
asql "mysql://user:password@host:3306/dbname"

# PostgreSQL
asql "postgres://user:password@host:5432/dbname"
```

## 主な機能

- **型情報付きヘッダ** — カラム名と型を並べて表示（`name text`、`age integer`）
- **NULL / 空文字の区別** — NULL は `NULL`、空文字は `""` で表示し混同を防止
- **インプレースソート** — `s` キーでソート切替（None → Asc → Desc）、NULL は常に末尾
- **クエリ履歴** — INSERT モードで `Ctrl+P` / `Ctrl+N` で過去のクエリを呼び出し
- **ページング表示** — ステータスバーに現在位置とカラム情報を表示（`col:name 1/100`）
- **テーブルサイドバー** — テーブル一覧をブラウズし、ワンキーで SELECT を挿入
- **エクスポート** — CSV / JSON / Markdown でコピー、またはファイル保存
- **AI アシスタント** — OpenAI 互換 API で自然言語から SQL を生成

## キーバインド

| キー | モード | 動作 |
|------|--------|------|
| `i` | NORMAL | INSERT モードに入る |
| `Esc` | INSERT | NORMAL モードに戻る |
| `Ctrl+Enter` / `Ctrl+J` | INSERT | クエリを実行 |
| `Ctrl+P` / `Ctrl+N` | INSERT | クエリ履歴の前 / 次 |
| `j` / `k` | NORMAL | 結果行を移動 |
| `h` / `l` | NORMAL | カラムを選択（ソート用） |
| `s` | NORMAL | 選択カラムのソートを切替 |
| `PgUp` / `PgDn` | NORMAL | ページ移動 |
| `t` | NORMAL | テーブルサイドバーを開く |
| `j` / `k` | SIDEBAR | テーブルを移動 |
| `Enter` | SIDEBAR | テーブルの SELECT クエリを挿入 |
| `Esc` / `t` | SIDEBAR | サイドバーを閉じる |
| `e` | NORMAL | エクスポートメニューを開く |
| `j` / `k` | EXPORT | 選択肢を移動 |
| `Enter` | EXPORT | エクスポート実行 |
| `Esc` | EXPORT | キャンセル |
| `Ctrl+K` | NORMAL | AI アシスタントを開く |
| `Enter` | AI | SQL を生成 |
| `Esc` | AI | キャンセル |
| `Ctrl+C` | *全モード* | 実行中のクエリ/AI をキャンセル、または終了 |
| `q` | NORMAL | 終了 |

## エクスポート

クエリ実行後、NORMAL モードで `e` を押すとエクスポートメニューが開きます。対応フォーマット:

- **Copy as CSV** — クリップボードにコピー
- **Copy as JSON** — クリップボードにコピー（オブジェクト配列）
- **Copy as Markdown** — クリップボードにコピー（GFM テーブル）
- **Save to File (CSV)** — カレントディレクトリに `result_YYYYMMDD_HHMMSS.csv` を保存

## AI アシスタント（Text-to-SQL）

OpenAI 互換 API を利用して、自然言語から SQL を生成できます。`~/.config/asql/config.yaml` に設定ファイルを作成してください:

```yaml
ai:
  ai_endpoint: http://localhost:11434/v1   # Ollama
  ai_model: llama3
  ai_api_key: ""                           # 省略可（Ollama は不要）
```

NORMAL モードで `Ctrl+K` を押すと AI プロンプトが開きます。データベースのスキーマ情報が自動的にコンテキストに含まれるため、正確なテーブル名・カラム名で SQL が生成されます。

設定ファイルがない場合、AI 機能はサイレントに無効化されます。

## 開発

```bash
go test ./...
go build
go vet ./...
```

## ライセンス

MIT — [LICENSE](LICENSE) を参照
