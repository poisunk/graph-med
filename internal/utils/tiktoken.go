package utils

import (
	"github.com/pkoukk/tiktoken-go"
	tiktoken_loader "github.com/pkoukk/tiktoken-go-loader"
)

var tke *tiktoken.Tiktoken

func init() {
	encoding := "cl100k_base"
	tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())

	var err error
	tke, err = tiktoken.GetEncoding(encoding)
	if err != nil {
		panic(err)
	}
}

func NumTokens(text string) int {
	tokens := tke.Encode(text, nil, nil)
	return len(tokens)
}
