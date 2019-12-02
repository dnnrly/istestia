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

	"github.com/dnnrly/istestia/markdown"
	"github.com/dnnrly/istestia/run"
	"github.com/spf13/cobra"
)

var testFile string
var isMarkdown bool

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run a single test file",
	RunE:  testCmdRun,
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVarP(
		&testFile,
		"file", "f",
		"",
		"the file that you would like to test against",
	)
	testCmd.Flags().BoolVarP(
		&isMarkdown,
		"markdown", "m",
		false,
		"set the input as markdown",
	)
}

func testCmdRun(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	tmpFiles := []string{}

	switch {
	case testFile == "" && len(args) == 0:
		return errors.New("must specify file with tests or string with tests")
	case isMarkdown:
		md := ""
		if testFile != "" {
			data, err := ioutil.ReadFile(testFile)
			if err != nil {
				return err
			}
			md = string(data)
		} else {
			md = args[0]
		}
		blocks, err := markdown.Extract(md)
		if err != nil {
			return err
		}
		for _, b := range blocks {
			if b.Type != "go" {
				return fmt.Errorf("language %s not supported", b.Type)
			}
			tmpFile := newTmpName(dir)
			err := writeFile(b.Contents, tmpFile)
			if err != nil {
				return err
			}
			defer os.Remove(tmpFile)
			tmpFiles = append(tmpFiles, tmpFile)
		}
	case testFile != "" && !isMarkdown:
		tmpFile := newTmpName(dir)
		err := copyFile(testFile, tmpFile)
		if err != nil {
			return err
		}
		defer os.Remove(tmpFile)
		tmpFiles = append(tmpFiles, tmpFile)
	case testFile == "" && !isMarkdown:
		tmpFile := newTmpName(dir)
		err := ioutil.WriteFile(tmpFile, []byte(args[0]), 0644)
		if err != nil {
			return err
		}
		defer os.Remove(tmpFile)
		tmpFiles = append(tmpFiles, tmpFile)
	}

	failed := false

	for _, tmpFile := range tmpFiles {
		code, err := parseFile(tmpFile)
		if err != nil {
			return err
		}

		tests := &visitor{}
		ast.Walk(tests, code)

		for _, t := range *tests {
			errTest := run.DoTest(t)
			if errTest != nil {
				failed = true
			}
		}
	}

	// We want to exit with an error but not with usage
	if failed {
		for _, t := range tmpFiles {
			os.Remove(t)
		}
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

func newTmpName(dir string) string {
	return fmt.Sprintf(
		"%s%cistestia_%s_test.go",
		dir,
		os.PathSeparator,
		randString(5),
	)
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

func writeFile(contents, dst string) error {
	err := ioutil.WriteFile(dst, []byte(contents), 0644)
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
