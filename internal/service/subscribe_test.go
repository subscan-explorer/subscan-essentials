package service

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/itering/substrate-api-rpc/rpc"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
)

type Buffer struct {
	bytes.Buffer
	io.ReaderFrom // conflicts with and hides bytes.Buffer's ReaderFrom.
	io.WriterTo   // conflicts with and hides bytes.Buffer's WriterTo.
}

func (*Buffer) Close() error {
	return nil
}

type TestConn struct {
	*websocket.Conn
	Connected bool
	writer    *Buffer // the current writer returned to the application
}

func (t *TestConn) MarkUnusable() {}

func (t *TestConn) Dial(urlStr string, reqHeader http.Header) {
	conn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	t.Conn = conn
	if err != nil {
		panic(err)
	}
	t.Connected = true
}

func (t *TestConn) IsConnected() bool {
	return t.Connected
}

func (t *TestConn) Close() {
	t.Connected = false
	if t.Conn != nil {
		t.Conn.Close()
	}
}

func (t *TestConn) ReadMessage() (messageType int, message []byte, err error) {
	return
}

func (t *TestConn) WriteMessage(messageType int, data []byte) error {
	if strings.EqualFold(string(data), `{"id":3,"method":"chain_subscribeFinalizedHeads","params":[],"jsonrpc":"2.0"}`) {
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}
	wb := new(Buffer)
	t.writer = wb
	if _, err := t.writer.Write(data); err != nil {
		return err
	}
	return t.writer.Close()
}

func (t *TestConn) ReadJSON(v interface{}) error {
	rb := new(Buffer)
	_, _ = io.Copy(rb, t.writer)
	var q rpc.Param
	_ = json.NewDecoder(rb).Decode(&q)
	switch q.Method {
	case "chain_getBlockHash":
		_ = json.Unmarshal([]byte(`{"result":"0xb6469b0823e733fde54437857f1b13b1a42ba21bb7e574941decedcf1571aa69"}`), v)
	case "chain_getBlock":
		_ = json.Unmarshal([]byte(`{"jsonrpc":"2.0","result":{"block":{"extrinsics":["0x280402000b70221dbd7301","0x1c0407003a004c00"],"header":{"digest":{"logs":["0x0642414245b5010106000000d160dc0f000000002e49423cc4c599cd47bed729976270c2b0e881927cb042779b7a2cad8ca4780282400a0fd2f0828e036ce768e07489b04788c62894ee8a1ec52f46e6aeaaf00ffb552f3b4d0ef5f33fe5ce9f94e4ab5eeb8636fc2a361b270933680e6e722305","0x00904d4d52525e2ce848d09d49bf88364ead19d42c278a079aaa1ba12d9e201f3f118f0f3097","0x05424142450101eaab19690f14a5c087f1602e594ae0622e8555d831544f8027109d5d1e8da359a269629819b89e2473e85a8d85f7210738a85143d01159a3e535653327e5fc8c"]},"extrinsicsRoot":"0x349025d9ecccca59046687048c5c0c658c65c065e5407e37948bdd94de834fda","number":"0x130011","parentHash":"0xc7d6e3b1f63085a89f7286791e471e2477741f8cb162b0d5e81a05af6332482d","stateRoot":"0x21cc68be9e14c43b000777a1ea591cebd401a646701d1002787cffe185ca31c9"}},"justification":null},"id":2}`), v)
	case "state_getStorageAt":
		if q.Id == wsEvent {
			_ = json.Unmarshal([]byte(`{"result":"0x080000000000000080e36a0900000000020000000100000000000000000000000000020000"}`), v)
		} else {
			// validator list
			_ = json.Unmarshal([]byte(`{"result":"0x7cb4f7f03bebc56ebe96bc52ea5ed3159d45a0ce3a8d7f082983c33ef13327474780a5d9612f5504f3e04a31ca19f1d6108ca77252bd05940031eb446953409c1ae0785340e721a52a5e406c3381a03ead07afb38ac447edc53267a4a0424b8176ca3bace1d97d5e997970fa351059c8580ce87991e71dc640ad5234b66c8cad250f9986176437e7e354fb505884f4b3286aea730ea7840f25d722558bb93c4ffbc4429847f3598f40008d0cbab53476a2f19165696aa41002778524b3ecf82938822d3fc99cc284a8b243328783442f65d799e546c295da6aee4a7dd9e6897519a6f7230e743089cdf1257e973fd65691ddcb4aa62ab5b0d579c0708d6fc28e6af6e367a71d6b817abfefd1338dbe5fd78f3f1816f237e4931d70c54b76351607762fd1b5bd1db180f27e9ee1c7a4196d85b86c6f9f55048e7733d4816bc2ff576482213046a034a1fbed5786a5bfac55782c7edfa9937ccb3d89ca28ab77c64710d904a7b054b178813f165ff68820b762cc8740d05baf4d15d1373c63720a49e8692be8e3c3af9e9eba0e6d6c73e2ad825ea1e5515f81726734798b6093e8034cc2fffa77764557754a249e577c86700e45442cc10032dade5489c02add7006be60e61c9a14ee4682d798c883dfec2df4ff34ad1c2e1bc8efcaa5a96f021355b61c73bc97b9ee9be396c585ba7d4df131e5eeecf5e1cd8732214c60de9d99090f998cb2bded511a9bc7c9ace3b0962b28ea29b1f92d1a6e6e51078be5c87ae9522dec17f772828eaa4719764cf057117614a80912ec5465087ee53c1ad3941dc6a8cab00d9f7a8b828ef0b8c982bb7374c261ee6453950c296b14d53ec15e6cfce6923ace2bb25f44fc064d51f7e5a054c0cc458f163c9ad3156abe3bcf29420f99888eef6cc5eb2de7c60f4cd71d9eabe094d61e0a0d2e106af6db88ddf663e30eeb7bbb002ead6d129b37e872b85d09dd6c09af7edac6c38f5dad8eb1c45e0f998513444a7c70ee486cc2d8d5546c388c982af8b6f701580007cd4dbf2c650f998d4d11d273909f1b5f659201245329d5c883914a22ed5b109b8e088f602f0f99845c37f38d70d0d28c78016a8ac0623096c32b19605a342dcdd3483ccada0f99850f381fdfb5b965617619e59ed082a4515c22967c9c1c94621969362abaced6b08a2f340fc4fda789e5f15d51ccbe5109c87924b6a51fb156ddc74a886e6e6844ba5c73db6c4c6b67ea59c2787dd6bd2f9b8139a69c33e14a722d1e801d68a2a0899f9bcd220193750b3ab17743bb62734832805b49625dd2dd9b21c30c8a9673ef43c3f41d8b772fd89a5d9abfa21afa4d9ecdeff1d182beb1b8e76b140f9980c64699fb44f34ee07e37e619e90fd172770c6f9b819ca5cf21494ef604"}`), v)
		}
	case "chain_getRuntimeVersion":
		_ = json.Unmarshal([]byte(`{"result":{"specVersion":5}}`), v)
	case "payment_queryInfo":
		_ = json.Unmarshal([]byte(`{"result":{"partialFee":1000000}}`), v)
	}
	return nil
}

func TestService_Subscribe(t *testing.T) {
	tc := TestConn{}
	interrupt := make(chan os.Signal, 1)
	testSrv.Subscribe(&tc, interrupt)
}
