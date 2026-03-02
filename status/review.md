# Code Review - 2026-02-28

## Summary
- Total findings: 6
- Critical: 0 | High: 1 | Medium: 3 | Low: 2
- Issues created: 4

## Findings

### Bug
- [High] `internal/db/sqlite/adapter.go:114` - `RETURNING` 句付き INSERT/UPDATE/DELETE で結果セットが返されない → Issue #1
- [Medium] `internal/db/sqlite/adapter.go:148` - `[]byte` の無条件 `string` 変換（バイナリ BLOB で TUI 表示が壊れる可能性） → Issue #3

### Quality
- [Medium] `internal/db/sqlite/adapter.go:107` - `returnsRows` が先頭キーワードのみに依存し保守性が低い → Issue #2
- [Low] `internal/ui/model.go:124` - `Update` 関数が肥大化（イベント処理の混在）

### Testing
- [Medium] `internal/db/sqlite/adapter.go:41` - `Query` メソッドのエッジケーステスト不足 → Issue #4

### Config
- [Low] `internal/ui/model.go:21` - `queryTimeout` (5秒) がコード内にハードコード（CLI フラグ非対応）

## Created Issues
- #1: [High] Bug: RETURNING句付きDML文で結果セットが返されない ✅ Resolved in PR #5
- #2: [Medium] Quality: returnsRows の保守性が低い ✅ Resolved in PR #5
- #3: [Medium] Bug: BLOBデータの unsafe string 変換 ✅ Resolved in PR #5
- #4: [Medium] Testing: Query メソッドのエッジケーステスト不足 ✅ Resolved in PR #5

## Resolution (2026-02-28)
PR #5 `fix/gemini-audit-issues` で Issue #1-#4 を一括解決:
- `containsReturning()` スキャナ追加（文字列/識別子/コメントスキップ + 単語境界チェック）
- `returnsRows()` を先頭キーワード + RETURNING句の2段階判定に変更
- `stringifyValue()` の `[]byte` を `utf8.Valid()` で判定し非UTF-8は hex 表示
- `TestContainsReturning` 新規、既存テストにエッジケース追加

## Low Severity (記録のみ・未対応)
- `internal/ui/model.go:124` - `Update` 関数の分割検討（現状は正常動作）
- `internal/ui/model.go:21` - `queryTimeout` をCLIフラグ化（初期リリース段階のため後回し）
