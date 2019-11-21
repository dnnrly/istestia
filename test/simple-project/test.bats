BIN="$(pwd)/istestia"

@test "Passes for a passing test" {
   cd $BATS_TEST_DIRNAME
   run ${BIN} test --file passing_test.fixture
   echo "$output" >&3
   [ "$status" -eq 0 ]
}

@test "Fails for a failing test" {
   cd $BATS_TEST_DIRNAME
   run ${BIN} test --file failing_test.fixture
   echo "$output" >&3
   [ "$status" -eq 1 ]
}
