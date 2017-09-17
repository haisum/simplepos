package stringutils

import (
	"fmt"
)

func ByteToString(byteArray []byte) string {
	return fmt.Sprintf("%s", byteArray)
}
