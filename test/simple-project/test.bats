BIN="$(pwd)/istestia"

@test "Passes for a passing test" {
   cd $BATS_TEST_DIRNAME
   run ${BIN} test --file passing_test.fixture
   [ "$status" -eq 0 ]
}

@test "Fails for a failing test" {
   cd $BATS_TEST_DIRNAME
   run ${BIN} test --file failing_test.fixture
   [ "$status" -eq 1 ]
}

@test "Passes for a passing test as string" {
   cd $BATS_TEST_DIRNAME
   run ${BIN} test "$(cat passing_test.fixture)"
   [ "$status" -eq 0 ]
}

@test "Fails for a failing test as string" {
   cd $BATS_TEST_DIRNAME
   run ${BIN} test "$(cat failing_test.fixture)"
   [ "$status" -eq 1 ]
}

@test "Passes for a passing test from pipe" {
   cd $BATS_TEST_DIRNAME
   run $(cat passing_test.fixture | ${BIN} test -)
   [ "$status" -eq 0 ]
}

@test "Fails for a failing test from pipe" {
   cd $BATS_TEST_DIRNAME
   run $(cat failing_test.fixture | ${BIN} test -)
   [ "$status" -eq 1 ]
}
