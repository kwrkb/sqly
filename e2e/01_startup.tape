# E2E: Startup + Sidebar + Table list
# Verifies: launch in INSERT mode, switch to NORMAL, open sidebar, see tables

Output e2e/recordings/01_startup.mp4
Output e2e/recordings/01_startup.gif

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

# 1. Launch asql — should start in INSERT mode
Type "./asql /tmp/asql-e2e.db"
Enter
Wait+Screen /INSERT/
Sleep 1s

# 2. Switch to NORMAL mode
Escape
Wait+Screen /NORMAL/
Sleep 500ms

# 3. Open table sidebar
Type "t"
Wait+Screen /Tables/
Sleep 500ms

# 4. Verify table names are listed
Wait+Screen /users/
Wait+Screen /posts/
Wait+Screen /tags/
Sleep 1s

# 5. Close sidebar and quit
Escape
Sleep 500ms
Type "q"
Sleep 500ms
