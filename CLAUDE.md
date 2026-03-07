# CLAUDE.md

This file provides guidance to AI assistants when working with code in this repository.

## Project Overview

**asql = Data Observation CLI**
asql は Go 製の軽量 TUI SQL クライアント。データを「速く・見やすく・並べて触れる」ことで違和感や仮説に気づくための「データ観察ツール（顕微鏡）」であり、巨大な分析基盤ではない。
- **Module**: `github.com/kwrkb/asql`
- **Framework**: Bubble Tea (Charmbracelet)
- **Support**: SQLite, MySQL, PostgreSQL
- **AI**: OpenAI 互換 API による Text-to-SQL 補助機能

## Commands

```bash
# ビルド
go build

# 実行
./asql <sqlite-file-path>
./asql "mysql://user:pass@host:3306/dbname"
./asql "postgres://user:pass@host:5432/dbname"

# テスト
go test ./...

# E2E テスト（VHS）— 実行前に e2e/README.md を必ず読むこと
bash e2e/run.sh

# 静的解析
go vet ./...
```

## Architecture

- **internal/db/** — データベース抽象層。
  - `adapter.go` — `DBAdapter` インターフェース（`Type`, `Query`, `Tables`, `Schema`, `Close`）。
  - `dbutil/` — 全アダプタ共通のユーティリティ（`returnsRows` 判定、値の文字列化など）。
  - `sqlite/`, `mysql/`, `postgres/` — 各 DB のアダプタ実装。
- **internal/ui/** — TUI 層。`model.go` (Bubble Tea) を中心に、NORMAL/INSERT/SIDEBAR/AI/EXPORT のモード管理。`export.go` にエクスポート関連の UI ロジックを分離。
- **internal/export/** — CSV/JSON/Markdown フォーマット変換ロジック。
- **internal/ai/** — LLM クライアント。スキーマ情報をプロンプトに注入。
- **internal/config/** — `~/.config/asql/config.yaml` の管理。

## Design Principles (from VISION.md)

1. **軽さを最優先**: 起動・反応速度を損なわない。
2. **思考を止めない**: キーボード中心。SQL を書き直さずに探索できる UX。
3. **ノイズを排除**: UI は目立たず、情報は必要なものだけ。
4. **観察を加速**: 比較や結合、基礎統計へのアクセスを容易にする。
5. **Bring Data Strategy**: 異種 DB を直接統合せず、ローカルに持ち寄って比較・結合する。

## Workflow Files

### VISION.md
- プロジェクトの理念、原則、非目標（Non-Goals）を定義する最上位文書。
- 機能追加の要否は `VISION.md` の決定ルール（Decision Rule）に従う。

### PLAN.md
- **未完了のタスクのみ**を管理するアクティブなロードマップ。
- フェーズごとに P0 (Must) / P1 (Should) 等の優先順位を付ける。

### HISTORY.md
- 完了したタスクの永続的な記録。
- **PLAN.md で完了したタスクは、随時こちらに移動して記録すること。**

### LESSONS.md
- 開発中に得た知見や回避した問題、ユーザーからの指摘事項を記録する。
- 同じミスを繰り返さないための「プロジェクトの知恵袋」として活用する。

### status/
- 品質監査スキルの出力先。`review.md`, `quality-report.md` などのレポートを格納。
