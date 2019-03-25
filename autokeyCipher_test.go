package kasalink

import "testing"

func TestKasaCrypt(t *testing.T) {
	t.Logf("%v", encrypt(getSysInfo))
}

func TestKasaDeCrypt(t *testing.T) {
	t.Logf("%s", decrypt(encrypt(getSysInfo)))
}
