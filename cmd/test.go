package cmd

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dnnrly/istestia/run"
	"github.com/spf13/cobra"
)

var testFile string

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run a single test file",
	RunE:  testCmdRun,
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVarP(&testFile, "file", "f", "", "the file that you would like to test against")
}

func testCmdRun(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	tmpFile := fmt.Sprintf("%s%cistestia_%s_test.go", dir, os.PathSeparator, randString(5))
	if testFile == "" {
		if len(args) == 0 {
			return errors.New("must specify file with tests")
		}

		err = ioutil.WriteFile(tmpFile, []byte(args[0]), 0644)
	} else {
		err = copyFile(testFile, tmpFile)
	}
	defer os.Remove(tmpFile)

	if err != nil {
		return err
	}

	code, err := parseFile(tmpFile)
	if err != nil {
		return err
	}

	tests := &visitor{}
	ast.Walk(tests, code)

	failed := false
	for _, t := range *tests {
		errTest := run.DoTest(t)
		if errTest != nil {
			failed = true
		}
	}

	// We want to exit with an error but not with usage
	if failed {
		os.Remove(tmpFile)
		fmt.Fprintf(os.Stderr, "Tests failed\n")
		os.Exit(1)
	}

	return nil
}

type visitor []string

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.FuncDecl:
		n := d.Name.Name
		if strings.HasPrefix(n, "Test") || strings.HasPrefix(n, "Example") {
			*v = append(*v, n)
		}
	}

	return v
}

func parseFile(src string) (*ast.File, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, src, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
