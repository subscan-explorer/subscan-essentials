package ipfs

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

type gatewayTransport func(*http.Request) (*http.Response, error)

func (f gatewayTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// stubGateway keeps gateway responses deterministic without changing production endpoints.
func stubGateway(t *testing.T, expectedURL string, status int, body string) {
	t.Helper()
	original := http.DefaultTransport
	http.DefaultTransport = gatewayTransport(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, expectedURL, req.URL.String())
		assert.Equal(t, http.MethodGet, req.Method)
		_, bounded := req.Context().Deadline()
		assert.True(t, bounded, "gateway requests must retain the bounded context")
		return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	t.Cleanup(func() { http.DefaultTransport = original })
}

func Test_OpenFile(t *testing.T) {
	assert.Error(t, verifyCid("123234567"))
	// test open file success
	ctx := context.Background()
	stubGateway(t, "https://bafkreidyeivj7adnnac6ljvzj2e3rd5xdw3revw4da7mx2ckrstapoupoq.ipfs.nftstorage.link", http.StatusOK, "fixture metadata")
	data, err := OpenFile(ctx, "bafkreidyeivj7adnnac6ljvzj2e3rd5xdw3revw4da7mx2ckrstapoupoq")
	assert.NoError(t, err)
	assert.Equal(t, "fixture metadata", string(data))
	_, err = OpenFile(ctx, "fff")
	assert.Error(t, err)
}

func Test_CheckUriImageExt(t *testing.T) {
	ext, err := CheckUriImageExt("ipfs://bafybeiftkn4roy2j3b2qipq3i4oavyvnjubvnavxmnkt3mmj6fnw5bvngq/26c96f4851a04e99845409bdd81a1868.png")
	assert.Equal(t, "png", ext)
	assert.NoError(t, err)
}

func Test_OpenArFile(t *testing.T) {
	stubGateway(t, "https://www.arweave.net/sDPTwwsvz_FtwuHUWkJ2lXzVSiQDldK_tgTYVIgxA3M", http.StatusOK, "arweave metadata")
	_, err := OpenArFile(context.Background(), "sDPTwwsvz_FtwuHUWkJ2lXzVSiQDldK_tgTYVIgxA3M")
	assert.NoError(t, err)
}

func Test_OpenFileGatewayError(t *testing.T) {
	stubGateway(t, "https://bafkreidyeivj7adnnac6ljvzj2e3rd5xdw3revw4da7mx2ckrstapoupoq.ipfs.nftstorage.link", http.StatusTooManyRequests, "rate limited")
	_, err := OpenFile(context.Background(), "bafkreidyeivj7adnnac6ljvzj2e3rd5xdw3revw4da7mx2ckrstapoupoq")
	assert.ErrorContains(t, err, "http error: 429")
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
