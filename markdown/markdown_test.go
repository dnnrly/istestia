package markdown

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustLoad(f string) string {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		panic(fmt.Sprintf("Unable read file %s: %v\n", f, err))
	}

	return string(data)
}

func TestExtract_single(t *testing.T) {
	data, err := Extract(mustLoad("test/simple.md"))

	assert.NoError(t, err)
	require.Equal(t, 1, len(data))
	assert.Equal(t, "go", data[0].Type)
	assert.Equal(t, "package pkg\n", data[0].Contents)
}
