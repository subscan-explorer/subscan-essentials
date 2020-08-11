package validator

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestValidate(t *testing.T) {
	model := new(struct {
		BlockNum  int    `json:"block_num" validate:"required"`
		BlockHash string `json:"block_hash" validate:"required"`
	})

	t1 := []byte{123, 34, 98, 108, 111, 99, 107, 95, 104, 97, 115, 104, 34, 58, 34, 102, 102, 102, 34, 44, 34, 98, 108, 111, 99, 107, 95, 110, 117, 109, 34, 58, 49, 125}
	err := Validate(t1, model)
	assert.NoError(t, err)

	t2 := map[string]interface{}{"block_num": 1, "block_hash": ""}
	err = Validate(t2, model)
	assert.Error(t, err)

	t3 := new(bytes.Buffer)
	t3.WriteString(`{"block_hash":"ttt","block_num":2}`)
	err = Validate(ioutil.NopCloser(t3), model)
	assert.NoError(t, err)

	t4 := map[string]interface{}{"block_num": "1", "block_hash": 12}
	err = Validate(t4, model)
	assert.Error(t, err)

}
