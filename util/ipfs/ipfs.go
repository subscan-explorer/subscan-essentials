package ipfs

import (
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/itering/subscan/util"
	"net/url"
	"strings"
	"time"
)

// https://ipfs.github.io/public-gateway-checker/
func OpenFile(ctx context.Context, id string) ([]byte, error) {
	subCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	var gateways = "ipfs.nftstorage.link"
	defer cancel()
	if err := verifyCid(id); err != nil {
		return nil, fmt.Errorf("cid %s verify failed %s", id, err)
	}
	return util.HttpGet(subCtx, fmt.Sprintf("https://%s.%s", id, gateways))

}
func verifyCid(id string) error {
	ids := strings.Split(id, "/")
	if len(ids) == 0 {
		return fmt.Errorf("invalid cid")
	}
	_, err := cid.Decode(ids[0])
	return err
}

func CheckUriImageExt(uri string) (string, error) {
	if strings.HasPrefix(uri, "ar://") {
		return "png", nil
	}
	uri = TrimMetadataUri(uri)
	if verifyCid(uri) != nil {
		return "", fmt.Errorf("cid %s verify failed", uri)
	}
	exts := []string{"svg", "png", "jpg", "bmp", "gif", "webp"}
	if split := strings.Split(uri, "."); len(split) >= 2 {
		if util.StringInSliceFold(split[len(split)-1], exts) {
			return split[len(split)-1], nil
		}
	}
	// default return png
	return "png", nil
}

func TrimMetadataUri(data string) string {
	uri := strings.ReplaceAll(data, "ipfs://ipfs/", "")
	uri = strings.ReplaceAll(uri, "ipfs://", "")
	uri = strings.ReplaceAll(uri, "ar://", "")
	uri = strings.ReplaceAll(uri, "ipfs/", "")
	uri = strings.ReplaceAll(uri, "https://ipfs.io/", "")
	uri = strings.ReplaceAll(uri, "https://arweave.net/", "")
	urls, err := url.Parse(uri)
	if err != nil {
		return uri
	}
	return strings.TrimPrefix(urls.Path, "/")
}

func OpenArFile(ctx context.Context, id string) ([]byte, error) {
	subCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	endpoint := fmt.Sprintf("https://%s/%s", "www.arweave.net", id)
	return util.HttpGet(subCtx, endpoint)
}
