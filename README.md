# Istia

This tool allows you to run a test that has been pulled from a PR description.

Imagine this, you are a user of an open source project and you find a problem. You manage
to reproduce the problem in to a neat test case that exposes the bug. Wouldn't it be nice
if, as a maintainer, you would be able to pull this new test case directly from the PR.

Examples:

```bash
$ export GO111MODULE=on
$ cd ~/project-dir
$ istestia test --file 'passes_test.go'
$ istestia test --file 'fails_test.go'
tests failed
$ cat passes_test.go | istestia test
$ cat fails_test.go | istestia test
tests failed
$ istestia test 'package project

import "testing"

func TestPass(t *testing.T) {
}
'
$ istestia test 'package project

import "testing"

func TestFail(t *testing.T) {
    t.Fail()
}
'
tests failed
$ cd ~/another-project
$ cat ~/project-dir/passes_test.go | istestia test
build failed
```
