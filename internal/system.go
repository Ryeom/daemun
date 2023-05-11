package internal

import (
	"runtime"
)

var operateModes = []string{"develop", "stage", "product"}

func IsOperationMode() bool {
	if runtime.GOOS == "darwin" {
		return false
	}
	if !Contains(operateModes, GetMode()) {
		return false
	}
	return true
}
