# E2E: Export overlay display
# Verifies: execute query, open export overlay, see format options

Output e2e/recordings/04_export.mp4

Set Shell "bash"
Set Width 1200
Set Height 600
Set FontSize 16
Set TypingSpeed 50ms
Set WaitTimeout 10s

# --- Setup ---
Hide
Type "python3 docs/setup-demo-db.py /tmp/asql-e2e.db /tmp/asql-e2e-staging.db >/dev/null 2>&1"
Enter
Sleep 1s
Show

# 1. Launch and execute a query first (export needs results)
Type "./asql /tmp/asql-e2e.db"
Enter
Wait+Screen /INSERT/
Ctrl+l
Sleep 500ms
Type "SELECT * FROM users"
Sleep 500ms
Ctrl+j
Wait+Screen /row/
Sleep 1s

# 2. Switch to NORMAL and open export
Escape
Wait+Screen /NORMAL/
Sleep 500ms
Type "e"
Wait+Screen /Export Results/
Sleep 500ms

# 3. Verify export format options are visible
Wait+Screen /CSV/
Wait+Screen /JSON/
Wait+Screen /Markdown/
Sleep 1s

# 4. Close export and quit
Escape
Sleep 500ms
Type "q"
Sleep 500ms
