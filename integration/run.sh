#!/bin/sh
# Integration checks for yup-cat, run inside a Debian (GNU coreutils) container.
#
# parity_stdin INPUT ARGS...  — yup-cat reading stdin must match GNU `cat`.
# parity_file  INPUT ARGS...  — yup-cat reading file operands must match GNU.
set -eu

fails=0
sample='alpha
beta

gamma'

parity_stdin() {
	in=$1
	shift
	ours=$(printf '%s\n' "$in" | yup-cat "$@" 2>/dev/null || true)
	gnu=$(printf '%s\n' "$in" | cat "$@" 2>/dev/null || true)
	if [ "$ours" = "$gnu" ]; then
		printf 'ok    parity  cat %s < stdin\n' "$*"
	else
		printf 'FAIL  parity  cat %s < stdin\n        gnu:  %s\n        ours: %s\n' "$*" "$gnu" "$ours"
		fails=$((fails + 1))
	fi
}

parity_file() {
	in=$1
	shift
	printf '%s\n' "$in" > /tmp/a.txt
	printf 'one\ntwo\n' > /tmp/b.txt
	ours=$(yup-cat "$@" 2>/dev/null || true)
	gnu=$(cat "$@" 2>/dev/null || true)
	if [ "$ours" = "$gnu" ]; then
		printf 'ok    parity  cat %s\n' "$*"
	else
		printf 'FAIL  parity  cat %s\n        gnu:  %s\n        ours: %s\n' "$*" "$gnu" "$ours"
		fails=$((fails + 1))
	fi
}

# stdin: identity, number-all (-n), number-nonblank (-b).
parity_stdin "$sample"
parity_stdin "$sample" -n
parity_stdin "$sample" -b
# -b overrides -n (GNU precedence).
parity_stdin "$sample" -n -b

# file operands: single, multiple, and concatenation order.
printf '%s\n' "$sample" > /tmp/a.txt
printf 'one\ntwo\n' > /tmp/b.txt
parity_file "$sample" /tmp/a.txt
parity_file "$sample" /tmp/a.txt /tmp/b.txt
parity_file "$sample" -n /tmp/a.txt /tmp/b.txt

if [ "$fails" -ne 0 ]; then
	printf '\n%s check(s) failed\n' "$fails"
	exit 1
fi
printf '\nall checks passed\n'
