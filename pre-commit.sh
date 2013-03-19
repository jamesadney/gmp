#!/bin/bash

## Install the hook ##
# `ln -s ../../pre-commit.sh .git/hooks/pre-commit`

run_tests() {
	# use `set -e` to exit the function as soon as a command fails.
	# the parens keep the `set -e` from affecting the calling script which
	# would make it exit before calling `stash pop`.
	( set -e
		./check-gofmt.sh # from `misc/git` of Go 1.1
		go test
	)
}

git stash --quiet --keep-index # only test staged changes
run_tests; RESULT=$?
git stash pop --quiet # reapply unstaged changes

exit $RESULT
