package stringutils

import "testing"

func TestByteToString(t *testing.T) {
	b := []byte(`hello world`)
	if ByteToString(b) != "hello world" {
		t.Errorf("Got %s, Expected %s", ByteToString(b), "hello world")
	}
}
