package kasalink

import "testing"

func TestKasaTurnPlugOn(t *testing.T) {
	var JSONString = KasaTurnPlugOn("abc123", 1, 3, 4)
	t.Logf("%s\n", JSONString)
}

func TestKasaTurnPlugOff(t *testing.T) {
	var JSONString = KasaTurnPlugOff("abc123", 1, 3, 4)
	t.Logf("%s\n", JSONString)
}
