BIN="$(pwd)/istestia"

setup() {
   rm -f istestia_*_.go
}

teardown() {
   rm -f istestia_*_.go
}

@test "Passes for a passing test from markdown" {
   cd $BATS_TEST_DIRNAME
   run ${BIN} test --markdown --file passing_test.fixture
   [ "$status" -eq 0 ]
   [ "$(find . | grep -c istestia_)" = "0" ]
}

@test "Passes test when passing markdown as arg" {
   cd $BATS_TEST_DIRNAME
   run $(cat passing_test.fixture | ${BIN} test --markdown)
   [ "$status" -eq 0 ]
   [ "$(find . | grep -c istestia_)" = "0" ]
}

@test "Fails for a failing test from markdown" {
   cd $BATS_TEST_DIRNAME
   run ${BIN} test --markdown --file failing_test.fixture
   [ "$status" -ne 0 ]
   [ "$(find . | grep -c istestia_)" = "0" ]
}

@test "Passes for a multiple passing tests from markdown" {
   cd $BATS_TEST_DIRNAME
   run ${BIN} test --markdown --file multiple_test.fixture
   [ "$status" -eq 0 ]
   [ "$(for l in ${lines[@]} ; do echo $l ; done | grep -c PASS)" = "2" ]
   [ "$(find . | grep -c istestia_)" = "0" ]
}

@test "Fails for wrong language in markdown" {
   cd $BATS_TEST_DIRNAME
   run ${BIN} test --markdown --file wrong_language_test.fixture
   [ "$status" -ne 0 ]
   [ "$(find . | grep -c istestia_)" = "0" ]
   [ "$(echo $output | grep -c 'not supported')" = "1" ]
}

