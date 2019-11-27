BIN="$(pwd)/istestia"

setup() {
   rm -f istestia_*_.go
}

teardown() {
   rm -f istestia_*_.go
}

@test "Passes for a passing test" {
   cd $BATS_TEST_DIRNAME
   run ${BIN} test --markdown --file passing_test.fixture
   [ "$status" -eq 0 ]
   [ "$(find . | grep -c istestia_)" -eq 0 ]
}

