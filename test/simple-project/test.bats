BIN="$(pwd)/istestia"

@test "Passes for a passing test" {
   run ${BIN} --file passing_test.fixture
   echo "$output" >&3
   [ "$status" -eq 0 ]
}

@test "Fails for a failing test" {
   skip
   run ${BIN} --file failing_test.fixture
   echo "$output" >&3
   [ "$status" -eq 1 ]
}
