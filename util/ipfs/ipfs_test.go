package ipfs

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_OpenFile(t *testing.T) {
	assert.Error(t, verifyCid("123234567"))
	// test open file success
	ctx := context.Background()
	data, err := OpenFile(ctx, "bafkreidyeivj7adnnac6ljvzj2e3rd5xdw3revw4da7mx2ckrstapoupoq")
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	_, err = OpenFile(ctx, "fff")
	assert.Error(t, err)
}

func Test_CheckUriImageExt(t *testing.T) {
	ext, err := CheckUriImageExt("ipfs://bafybeiftkn4roy2j3b2qipq3i4oavyvnjubvnavxmnkt3mmj6fnw5bvngq/26c96f4851a04e99845409bdd81a1868.png")
	assert.Equal(t, "png", ext)
	assert.NoError(t, err)
}

func Test_OpenArFile(t *testing.T) {
	_, err := OpenArFile(context.Background(), "sDPTwwsvz_FtwuHUWkJ2lXzVSiQDldK_tgTYVIgxA3M")
	assert.NoError(t, err)
}

func Test_TrimMetadataUri(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ipfs://ipfs/QmTzQ1Nj5x1", "QmTzQ1Nj5x1"},
		{"ipfs://QmTzQ1Nj5x1", "QmTzQ1Nj5x1"},
		{"ar://QmTzQ1Nj5x1", "QmTzQ1Nj5x1"},
		{"ipfs/QmTzQ1Nj5x1", "QmTzQ1Nj5x1"},
		{"https://ipfs.io/QmTzQ1Nj5x1", "QmTzQ1Nj5x1"},
		{"https://arweave.net/QmTzQ1Nj5x1", "QmTzQ1Nj5x1"},
		{"https://example.com/path/to/resource", "path/to/resource"},
	}

	for _, test := range tests {
		result := TrimMetadataUri(test.input)
		assert.Equal(t, test.expected, result)
	}
}
