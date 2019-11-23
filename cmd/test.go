package cmd

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	testCmd.Flags().StringVarP(&testFile, "file", "f", "", "the file that you would like to test against")
}

func testCmdRun(cmd *cobra.Command, args []string) error {
	if testFile == "" {
		return errors.New("must specify file with tests")
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	tmpFile := fmt.Sprintf("%s%cistestia_%s_test.go", dir, os.PathSeparator, randString(5))

	err = copyFile(testFile, tmpFile)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

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
