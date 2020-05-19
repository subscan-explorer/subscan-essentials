package util

import "github.com/go-kratos/kratos/pkg/ecode"

var (
	ParamsError           = ecode.New(10001)
	InvalidAccountAddress = ecode.New(10002)
	RecordNotFound        = ecode.New(10004)
)

func init() {
	ecode.Register(map[int]string{
		0:     "Success",
		10001: "Params Error",
		10002: "Invalid Account Address",
		10004: "Record Not Found",
	})

}
