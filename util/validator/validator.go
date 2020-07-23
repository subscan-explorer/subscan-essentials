package validator

import (
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"io"
	"io/ioutil"
)

var validate = validator.New()

func Validate(data interface{}, model interface{}) (err error) {
	var b []byte
	switch v := data.(type) {
	case []byte:
		b = v
	case io.ReadCloser:
		b, err = ioutil.ReadAll(v)
	default:
		b, _ = json.Marshal(data)
	}
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, model); err != nil {
		return err
	}
	return validate.Struct(model)
}
