package address

import (
	"github.com/itering/subscan/util"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

var testCases = []struct {
	pk      string
	prefix  int
	address string
}{
	{"3a370c6e4af506123c30e091a1cbfbc3728e1ec5fc47d87457fbb0b504903260", 0, "12KL8YptX9SuUCZGrsNrSRzp3zHNqbwLqmfN8vubtj1z1Bqv"},
	{"8c5d672931b073b1a7bca99e6aaa7cfedb6c8e8b5279e3c0d6eb5eb5768fbd39", 2, "FkMy71y6X8eYTqn1QKJTf39eaLzH1ud2ZjUaqV3fYNdFPqG"},
	{"5eea80719566804542730a997f2b4f94766428e46f2f919248782c6d6377901c", 5, "Y5kTQ89XcPoodYGeTUZbRM3usR3JRDL1hy9eYsrdiRgpmGh"},
	{"0890a6c7e0c98bc7a7466c5c07eeaec85784627c1fb4360b5071c8da267c383e", 42, "5CFwDB4oKXFv1EU3ziJXG61gdwqbWViups5e4pEzRi7zAVCp"},
	{"124de43e638b9cb913e7f8b619f72172824c5fdd43cd3e9c31127d365a75223c", 7, "hvcoyfLo51QTJW3dSfP5gm7fiRa8BKX6kvLaBM673pzeepv"},
	{"124de43e638b9cb913e7f8b619f72172824c5fdd43c", 16, ""},
}

func Test_SS58Address(t *testing.T) {
	for _, test := range testCases {
		util.AddressType = strconv.Itoa(test.prefix)
		assert.Equal(t, SS58Address(test.pk), test.address)
	}
}
