# E2E: Error message display
# Verifies: invalid SQL shows error in status bar

Output e2e/recordings/05_error.mp4

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

# 1. Launch
Type "./asql /tmp/asql-e2e.db"
Enter
Wait+Screen /INSERT/

# 2. Type invalid SQL and execute
Ctrl+l
Sleep 500ms
Type "SELECT * FROM nonexistent_table"
Sleep 500ms
Ctrl+j
Wait+Screen /no such table/
Sleep 1s

# 3. Quit
Escape
Sleep 500ms
Type "q"
Sleep 500ms
