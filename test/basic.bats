BIN=./istestia

@test "Can run the application" {
    run ${BIN}
    [ $status -eq 0 ]
}

