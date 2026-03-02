# Multi AI Quality Report - 2026-02-28

## Summary

- Total findings: 18
- Critical: 0 | High: 4 | Medium: 6 | Low: 8
- Consensus findings: 0 (Gemini: timeout, Codex: output truncated — Claude subagent only)
- AIs used: Claude (Gemini/Codex は結果なし)

---

## Consensus Findings

なし（有効な結果を返したAIが1つのみ）

---

## All Findings by Category

### Security

- [High] `/internal/db/sqlite/adapter.go:37-60` - Raw SQL を無検証で ExecContext/QueryContext に渡している。ローカル TUI 用途では問題ないが、将来のネットワーク利用・マルチユーザー対応時にSQLインジェクションリスクとなる。DDL/DML（DROP, DELETE）実行前に確認プロンプトが存在しない。
  - **Recommendation**: 破壊的ステートメント（DROP/DELETE/TRUNCATE）はUI層で確認ダイアログを挟む。threat boundary をドキュメント化する。
- [Low] `/internal/db/sqlite/adapter.go:20` - Ping() 前にパスが有効な SQLite ファイルかを検証していない。パストラバーサル文字列も受け入れてしまう。
  - **Recommendation**: `os.Stat` で通常ファイルか確認。DSN に `?_foreign_keys=on` を追加する。

### Code Quality

- [Medium] `/internal/ui/model.go:280-311` - `renderStatusBar` が `View()` コール毎に `lipgloss.NewStyle()` インスタンスを複数生成。フレームレンダリング毎に高頻度アロケーション発生。
  - **Recommendation**: 静的ベーススタイルを `NewModel` 時に1回構築してキャッシュし、動的部分のみ毎回適用する。
- [Medium] `/internal/ui/model.go:237-247` - `syncViewport` が複数のコードパスから冗長に呼ばれている。将来の編集で二重同期バグを招くリスク。
  - **Recommendation**: `Update` 末尾の単一コールサイトに集約し、各ハンドラメソッドから削除する。
- [Low] `/internal/ui/model.go:17` - `mode` 型の switch に `default` ブランチがなく、無効モード値がサイレントに fallthrough する。
  - **Recommendation**: `default: panic("unknown mode")` またはエラーログを追加し、開発時に検出できるようにする。
- [Low] `/internal/ui/model.go:55-115` - `NewModel` が60行超で textarea 初期化（スタイル設定25行以上）と table 初期化を両方担当。
  - **Recommendation**: `newTextarea()` / `newTable()` ヘルパーを抽出し、`NewModel` をアセンブリ専用にする。

### Bug

- [High] `/internal/db/sqlite/adapter.go:52-55` - `res.RowsAffected()` のエラーをサイレントに破棄し、`"statement executed"` を返す。呼び出し元がエラーを検知できない。
  - **Recommendation**: エラーを呼び出し元に返すか、最低限ログに記録する。
- [Medium] `/internal/ui/model.go:313-321` - クエリタイムアウトが `5*time.Second` にハードコード。長時間クエリがサイレントにキャンセルされ、ユーザーには raw エラー文字列のみ表示。
  - **Recommendation**: タイムアウトを定数化し、"Query timed out after 5s" など分かりやすいメッセージを表示する。
- [Medium] `/internal/ui/model.go:121-153` - `Update` は値レシーバ、`resize`/`syncViewport`/`applyResult`/`setStatus` はポインタレシーバで混在。設計意図が不明瞭で、将来の貢献者を混乱させる。
  - **Recommendation**: Bubble Tea イディオム（値レシーバ、モデルを返す）に統一する。
- [Low] `/internal/ui/model.go:265-266` - `result.Rows` が空の場合のセンチネル行 `{"(no rows)"}` がセル1つのみで、複数カラムのテーブルではレンダリング異常またはパニックの可能性。
  - **Recommendation**: `make(table.Row, len(columns))` でカラム数分パディングする。

### Testing

- [High] `/internal/db/sqlite/adapter.go` - `Query`/`queryRows`/`Open` のテストが0件。DB実際の操作パス（SELECT/DML ルーティング、`rows.Err()`、コンテキストキャンセル、`RowsAffected` パス）が全て未テスト。
  - **Recommendation**: インメモリ SQLite（`:memory:`）を使った統合テストを追加する。
- [High] `/internal/ui/model.go` - `model_test.go` が `columnWidth` のみをテスト。`Update`/`updateNormal`/`updateInsert`/`applyResult`/`resize`/`renderStatusBar`/`executeQueryCmd` が全て未テスト。
  - **Recommendation**: モック `DBAdapter` を使い、モード切替・クエリ発行・結果適用のユニットテストを追加する。
- [Medium] `/internal/db/sqlite/adapter_test.go:64-90` - `TestStringifyValue` に `bool false`、非UTF-8 `[]byte`、カスタム構造体のテストケースが欠如。
  - **Recommendation**: 上記ケースをテーブルに追加する。
- [Low] `/internal/db/sqlite/adapter_test.go` - エクスポートされた `Open` と `Adapter.Query` のテストが存在しない。
  - **Recommendation**: `package sqlite_test` から公開 API のテストを追加する。

### Config

- [Medium] `/internal/ui/model.go:315` - `5 * time.Second` が magic number として `executeQueryCmd` にハードコード。名前付き定数なし。
  - **Recommendation**: `const queryTimeout = 5 * time.Second` を定義し参照する。
- [Low] `/internal/ui/model.go:24-35` - UI カラー値が `var` で定義されており、実行時に変更可能。`const` が意図に合っている。
  - **Recommendation**: `lipgloss.Color` は `type Color string` なので `const` ブロックに変更する。
- [Low] `/home/yugosasaki/code/asql/main.go:15-18` - `--version` / `--help` フラグなし。引数解析がアドホック実装で、オプション追加時に互換性が壊れる。
  - **Recommendation**: 早期に `flag` パッケージまたは軽量 CLI ライブラリ（cobra 等）を採用する。

---

## Per-AI Raw Results

<details>
<summary>Claude Subagent Output</summary>

SUMMARY: Critical=0 High=4 Medium=6 Low=8

High:
- Security: adapter.go:37-60 - Raw SQL 無検証実行
- Bug: adapter.go:52-55 - RowsAffected エラー破棄
- Testing: adapter.go - DB操作テスト0件
- Testing: model.go - UI操作テスト0件 (columnWidthのみ)

Medium:
- Quality: model.go:280-311 - View()毎のスタイルアロケーション
- Quality: model.go:237-247 - syncViewport 冗長呼び出し
- Bug: model.go:313-321 - 5秒タイムアウトハードコード
- Bug: model.go:121-153 - 値/ポインタレシーバ混在
- Testing: adapter_test.go:64-90 - stringifyValue テスト不足
- Config: model.go:315 - magic number タイムアウト

Low (8件): Security:adapter.go:20, Quality:model.go:17, Quality:model.go:55-115, Bug:model.go:265-266, Testing:adapter_test.go, Config:model.go:24-35, Config:main.go:15-18

</details>

<details>
<summary>Gemini Output</summary>

タイムアウト（180秒）により結果なし。起動・認証は成功したが、レビュー出力を返す前に終了。

</details>

<details>
<summary>Codex Output</summary>

ファイル読み込みは完了したが出力が途中で打ち切られ、findings を返さなかった。

</details>
