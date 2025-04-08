package util

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/huandu/xstrings"
	"github.com/shopspring/decimal"
)

func AbiStringDecoder(data string) string {
	data = TrimHex(data)
	dataSlice := DataAnalysis(data)
	if len(dataSlice) == 1 {
		return string(bytes.TrimRight(HexToBytes(data), "\u0000"))
	}
	if len(dataSlice) < 3 {
		return ""
	}

	offset := StringToInt(HexToNumStr(dataSlice[0])) % 32
	length := StringToInt(HexToNumStr(dataSlice[1])) * 2
	data = data[128+offset*64 : 128+offset*64+length]
	return string(bytes.TrimRight(HexToBytes(data), "\u0000"))
}

func DataAnalysis(log string) []string {
	log = strings.TrimPrefix(log, "0x")
	logLength := len(log)
	var logSlice []string
	for i := 0; i < logLength/64; i++ {
		logSlice = append(logSlice, log[i*64:(i+1)*64])
	}
	return logSlice
}

func Padding(str string) string {
	str = strings.TrimPrefix(str, "0x")
	return xstrings.RightJustify(str, 64, "0")
}
func PaddingLeft(str string) string {
	str = strings.TrimPrefix(str, "0x")
	return xstrings.LeftJustify(str, 64, "0")
}

func BigIntToHex(b string) string {
	decimalTokenId, err := decimal.NewFromString(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", decimalTokenId.Coefficient())
}
