# E2E: Mode transition cycle
# Verifies: INSERT -> NORMAL -> SIDEBAR -> NORMAL -> INSERT

Output e2e/recordings/03_mode_transitions.mp4

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

# 1. Launch — INSERT mode
Type "./asql /tmp/asql-e2e.db"
Enter
Wait+Screen /INSERT/
Sleep 1s

# 2. INSERT -> NORMAL
Escape
Wait+Screen /NORMAL/
Sleep 500ms

# 3. NORMAL -> SIDEBAR
Type "t"
Wait+Screen /SIDEBAR/
Sleep 500ms

# 4. SIDEBAR -> NORMAL (Escape closes sidebar)
Escape
Wait+Screen /NORMAL/
Sleep 500ms

# 5. NORMAL -> INSERT
Type "i"
Wait+Screen /INSERT/
Sleep 500ms

# 6. Quit
Escape
Sleep 500ms
Type "q"
Sleep 500ms
