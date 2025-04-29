package token

type Token struct {
	TokenId string `json:"token_id"`

	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

var defaultToken *Token

func GetDefaultToken() *Token {
	return defaultToken
}

func SetDefault(t *Token) {
	defaultToken = t
}
