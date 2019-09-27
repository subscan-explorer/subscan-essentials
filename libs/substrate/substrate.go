package substrate

import (
	"fmt"
	"github.com/gorilla/websocket"
	"math/big"
	"math/rand"
	"strings"
	"subscan-end/libs/substrate/protos/codec_protos"
	"subscan-end/libs/substrate/storage"
	"subscan-end/utiles"
)

type Websocket struct {
	Provider string `json:"provider"`
}

func (w *Websocket) GetFreeBalance(module, accountId string) (balanceValue *big.Int, err error) {
	module = strings.ToUpper(string(module[0])) + string(module[1:])
	if balanceValue, err := GetStorage(nil, module, "FreeBalance", utiles.TrimHex(accountId)); err == nil {
		return balanceValue.ToBigInt(), nil
	}
	return nil, err

}

func (w *Websocket) GetChainInfo(c *websocket.Conn) (err error) {
	v := &JsonRpcResult{}
	if err = c.WriteMessage(websocket.TextMessage, SystemChain(3001)); err != nil {
		return err
	}
	_ = c.ReadJSON(v)
	chain, _ := v.ToString()
	if err = c.WriteMessage(websocket.TextMessage, SystemName(3002)); err != nil {
		return err
	}
	_ = c.ReadJSON(v)
	name, _ := v.ToString()
	if err = c.WriteMessage(websocket.TextMessage, SystemVersion(3003)); err != nil {
		return err
	}
	_ = c.ReadJSON(v)
	version, _ := v.ToString()
	fmt.Println(chain, name, version)
	return
}

// Get storage with hash
func GetStorageAt(c *websocket.Conn, hash, section, method string, arg ...string) (r storage.StateStorage, err error) {
	if c == nil {
		if c, err = initWebsocket(); err != nil {
			return
		}
	}
	storageKey, scaleType := encodeStorageKey(section, method, arg...)
	v := &JsonRpcResult{}
	if err = c.WriteMessage(websocket.TextMessage, StateGetStorageAt(rand.Intn(10000), utiles.AddHex(storageKey), hash)); err != nil {
		return
	}
	_ = c.ReadJSON(v)
	if dataHex, err := v.ToString(); err == nil {
		decodeMsg, err := codec_protos.DecodeStorage(dataHex, scaleType)
		return storage.StateStorage(decodeMsg), err
	}
	return r, err
}

// Get storage without hash
func GetStorage(c *websocket.Conn, section, method string, arg ...string) (r storage.StateStorage, err error) {
	if c == nil {
		if c, err = initWebsocket(); err != nil {
			return
		}
	}
	storageKey, scaleType := encodeStorageKey(section, method, arg...)
	v := &JsonRpcResult{}
	if err = c.WriteMessage(websocket.TextMessage, StateGetStorage(rand.Intn(10000), utiles.AddHex(storageKey))); err != nil {
		return
	}
	_ = c.ReadJSON(v)
	if dataHex, err := v.ToString(); err == nil {
		decodeMsg, err := codec_protos.DecodeStorage(dataHex, scaleType)
		return storage.StateStorage(decodeMsg), err
	} else {
		return r, err
	}
}
