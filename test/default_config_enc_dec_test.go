package test

import (
	"fmt"
	"github.com/Ryeom/daemun/internal"
	"testing"
)

var gatewayAesKey = ""

func TestEncryption(t *testing.T) {
	list := map[string]string{"": ""}
	for k, _ := range list {
		encValue, _ := internal.EncryptAES(k, []byte(gatewayAesKey))
		fmt.Println(k, " = ", encValue)
	}

}
func TestDecryption(t *testing.T) {
	list := map[string]string{"": ""}
	for k, v := range list {
		decValue, _ := internal.DecryptAES(v, []byte(gatewayAesKey))
		fmt.Println(k, " = ", decValue)
	}
}

func TestContains(t *testing.T) {
	fmt.Println(internal.Contains([]any{"asdf", "zxcv"}, "asdf"))
	fmt.Println(internal.Contains([]any{1, 2}, 2))
}
