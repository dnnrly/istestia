package run

import (
	"fmt"
	"os/exec"
	"strings"
)

func DoTest(name string) error {
	test := name
	if strings.HasPrefix(test, "Test") {
		test = strings.Replace(test, "Test", "", 1)
	}

	cmd := exec.Command("go", "test", "-run", test)
	output, err := cmd.CombinedOutput()
	fmt.Printf("Test output:\n%s\n", string(output))

	return err
}
