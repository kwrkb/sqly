# E2E: Query execution + result display
# Verifies: type query, execute, see result rows and data

Output e2e/recordings/02_query_exec.mp4

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

# 1. Launch asql
Type "./asql /tmp/asql-e2e.db"
Enter
Wait+Screen /INSERT/

# 2. Clear editor and type a query
Ctrl+l
Sleep 500ms
Type "SELECT * FROM users"
Sleep 500ms

# 3. Execute query
Ctrl+j
Wait+Screen /row/
Sleep 1s

# 4. Verify data is displayed
Wait+Screen /alice@example\.com/
Wait+Screen /bob@example\.com/
Wait+Screen /charlie@example\.com/
Sleep 1s

# 5. Quit
Escape
Sleep 500ms
Type "q"
Sleep 500ms
