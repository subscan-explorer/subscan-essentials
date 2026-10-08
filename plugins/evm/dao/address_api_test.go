package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	gorilla "github.com/gorilla/websocket"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/metadata"
	ws "github.com/itering/substrate-api-rpc/websocket"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const sampleNative = "5400e2f7f5669b26998d8e4d3c1a2c8a2d0a9af827ca54a1cc3509105035c32e"
const sampleH160 = "0x2ee6ace1cdc464a61d3f3910a6853ab5d093a107"

func mappingMetadata(t *testing.T) {
	t.Helper()
	// Minimal decoded metadata for the storage used by the resolver.
	m := metadata.Latest(&metadata.RuntimeRaw{Spec: 1, Raw: "0x6d6574610900"})
	old := *m
	t.Cleanup(func() { *m = old })
	if err := json.Unmarshal([]byte(`{"metadata":{"modules":[{"name":"Revive","prefix":"Revive","storage":[{"name":"OriginalAccount","type":{"origin":"MapType","map_type":{"key":"H160","hasher":"Identity","value":"AccountId"}}}]}]}}`), m); err != nil {
		t.Fatal(err)
	}
}

func TestBalanceAPIWithIsolatedDatabase(t *testing.T) {
	mappingMetadata(t)
	oldSG, oldEVM := sg, util.IsEvmChain
	t.Cleanup(func() { sg = oldSG; util.IsEvmChain = oldEVM; ws.Close() })
	util.IsEvmChain = false
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	// Exercise GORM's real SELECT through its existing MySQL dialect on SQLite.
	// The schema and rows are test-only; no configured service is contacted.
	_, err = db.Exec("CREATE TABLE balance_accounts (id INTEGER PRIMARY KEY, address TEXT, nonce INTEGER, balance TEXT, locked TEXT, reserved TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	fallback := reviveAccount("0x1111111111111111111111111111111111111111")
	for _, row := range []struct{ addr, amount string }{{sampleNative, "53700000000"}, {fallback, "1000"}} {
		if _, err = db.Exec("INSERT INTO balance_accounts(address,nonce,balance,locked,reserved) VALUES (?,0,?,'0','0')", row.addr, row.amount); err != nil {
			t.Fatal(err)
		}
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sg = &Storage{db: gdb}
	var rpcError atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := (&gorilla.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			var request struct {
				ID     int      `json:"id"`
				Method string   `json:"method"`
				Params []string `json:"params"`
			}
			if c.ReadJSON(&request) != nil {
				return
			}
			response := map[string]interface{}{"jsonrpc": "2.0", "id": request.ID, "result": nil}
			const prefix = "0x735f040a5d490f1107ad9c56f5ca00d2c56ab6c1f203b345fe5879f819627723"
			if request.Method != "state_getStorage" || len(request.Params) != 1 || !strings.HasPrefix(request.Params[0], prefix) {
				response["error"] = map[string]interface{}{"code": -32602, "message": "bad storage key"}
			} else if rpcError.Load() {
				response["error"] = map[string]interface{}{"code": -32000, "message": "unavailable"}
			} else if request.Params[0] == prefix+util.TrimHex(sampleH160) {
				response["result"] = "0x" + sampleNative
			}
			if c.WriteJSON(response) != nil {
				return
			}
		}
	}))
	defer server.Close()
	ws.Close()
	ws.SetEndpoint("ws" + strings.TrimPrefix(server.URL, "http"))
	srv := &ApiSrv{}
	result, err := srv.API_GetAccounts(context.Background(), []string{sampleH160, "0x1111111111111111111111111111111111111111"})
	if err != nil {
		t.Fatal(err)
	}
	if result[sampleH160].Address != sampleNative || result[sampleH160].Balance.String() != "53700000000" {
		t.Fatalf("mapped balance: %+v", result)
	}
	if result["0x1111111111111111111111111111111111111111"].Balance.String() != "1000" {
		t.Fatalf("fallback balance: %+v", result)
	}
	ws.Close()
	rpcError.Store(true)
	if _, err = srv.API_GetAccounts(context.Background(), []string{sampleH160}); err == nil {
		t.Fatal("RPC error became a balance response")
	}
}

func TestLiveReviveMapping(t *testing.T) {
	path := os.Getenv("SUBSCAN_MAPPING_METADATA")
	if path == "" {
		t.Skip("opt-in: set SUBSCAN_MAPPING_METADATA to a live metadata JSON response")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var replies []struct {
		ID     int             `json:"id"`
		Result json.RawMessage `json:"result"`
	}
	if err = json.Unmarshal(raw, &replies); err != nil {
		t.Fatal(err)
	}
	var chain, encoded string
	for _, r := range replies {
		if r.ID == 1 {
			json.Unmarshal(r.Result, &chain)
		}
		if r.ID == 2 {
			json.Unmarshal(r.Result, &encoded)
		}
	}
	if chain != "Polkadot Asset Hub" && chain != "Asset Hub Polkadot" {
		t.Fatalf("unexpected chain %q", chain)
	}
	current := metadata.Latest(nil)
	if current == nil {
		metadata.Latest(&metadata.RuntimeRaw{Spec: 1, Raw: encoded})
	} else {
		old := *current
		t.Cleanup(func() { *current = old })
		decoder := scalecodec.MetadataDecoder{}
		decoder.Init(util.HexToBytes(encoded))
		if err := decoder.Process(); err != nil {
			t.Fatal(err)
		}
		*current = metadata.Instant(decoder.Metadata)
	}
	ws.Close()
	ws.SetEndpoint("wss://polkadot-asset-hub-rpc.polkadot.io")
	defer ws.Close()
	got, err := readReviveOriginalAccount(context.Background(), sampleH160)
	if err != nil || got != sampleNative {
		t.Fatalf("live mapping: %q %v", got, err)
	}
	t.Logf("%s: %s -> %s", chain, sampleH160, got)
}
