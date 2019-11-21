package cmd

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

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

	err = os.Remove(tmpFile)
	if err != nil {
		return err
	}

	return nil
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
