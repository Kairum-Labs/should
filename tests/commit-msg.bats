#!/usr/bin/env bash
# tests/commit-msg.bats — BATS tests for scripts/commit-msg
#
# Run locally:
#   bats tests/commit-msg.bats
#
# BATS is installed automatically in the CI workflow; to install locally:
#   brew install bats-core        # macOS
#   sudo apt install bats         # Ubuntu/Debian

SCRIPT="$(cd "$(dirname "$BATS_TEST_FILENAME")/.." && pwd)/scripts/commit-msg"

# ── helpers ────────────────────────────────────────────────────────────────

# write_msg <content>  — create a temp file with the given content and store
# its path in $MSG_FILE (cleaned up after each test via teardown).
setup() {
  MSG_FILE="$(mktemp)"
}

teardown() {
  rm -f "$MSG_FILE"
}

write_msg() {
  printf '%s' "$1" > "$MSG_FILE"
}

# ── tests ──────────────────────────────────────────────────────────────────

@test "passes when subject is exactly 72 characters" {
  write_msg "$(printf '%0.s-' {1..72})"   # 72 dashes
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 0 ]
}

@test "passes when subject is shorter than 72 characters" {
  write_msg "feat: add short commit message"
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 0 ]
}

@test "fails when subject is 73 characters (one over limit)" {
  write_msg "$(printf '%0.s-' {1..73})"   # 73 dashes
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 1 ]
}

@test "fails when subject is much longer than 72 characters" {
  write_msg "feat: this is a very long commit subject line that clearly exceeds the 72 character limit by quite a lot"
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 1 ]
}

@test "error output reports the actual character count" {
  write_msg "$(printf '%0.s-' {1..80})"   # 80 dashes
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 1 ]
  [[ "$output" == *"80"* ]]
}

@test "error output includes the limit (72)" {
  write_msg "$(printf '%0.s-' {1..80})"
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 1 ]
  [[ "$output" == *"72"* ]]
}

@test "ignores Git comment lines (lines starting with #)" {
  write_msg "$(printf '%s\n' \
    '# This is a git comment' \
    '# On branch main' \
    'feat: real subject line')"
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 0 ]
}

@test "ignores leading blank lines before the subject" {
  write_msg "$(printf '%s\n' '' '' 'feat: subject after blank lines')"
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 0 ]
}

@test "fails on an empty commit message" {
  write_msg ""
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 1 ]
}

@test "fails on a message that contains only blank lines" {
  write_msg "$(printf '\n\n\n')"
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 1 ]
}

@test "fails on a message that contains only Git comments" {
  write_msg "$(printf '%s\n' '# comment one' '# comment two')"
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 1 ]
}

@test "fails when commit message file argument is missing" {
  run bash "$SCRIPT"
  [ "$status" -eq 1 ]
}

@test "fails when commit message file does not exist" {
  run bash "$SCRIPT" "/tmp/does-not-exist-$RANDOM"
  [ "$status" -eq 1 ]
}

@test "subject at exactly 72 chars is not reported as over limit" {
  write_msg "$(printf '%0.s-' {1..72})"
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 0 ]
  [[ "$output" != *"Error"* ]]
}

@test "only the subject line is checked, not the body" {
  # Subject is 40 chars; body line is 120 chars — should still pass
  LONG_BODY=$(printf '%0.sa' {1..120})
  write_msg "$(printf '%s\n\n%s' 'feat: normal subject line' "$LONG_BODY")"
  run bash "$SCRIPT" "$MSG_FILE"
  [ "$status" -eq 0 ]
}
