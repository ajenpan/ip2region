package ip2region

import (
	"testing"
)

func BenchmarkMemorySearch(B *testing.B) {
	region, err := New("../../conf/ip2region.db ")
	if err != nil {
		B.Error(err)
	}
	for i := 0; i < B.N; i++ {
		region.MemorySearch("127.0.0.1")
	}
}

func TestIp2long(t *testing.T) {
	ip, err := StrIP2Int("127.0.0.1")
	if err != nil {
		t.Error(err)
	}
	if ip != 2130706433 {
		t.Error("result error")
	}
	t.Log(ip)
}
func TestMemorySearch(t *testing.T) {
	region, err := New("../../conf/ip2region.db ")
	if err != nil {
		t.Error(err)
	}
	_, err = region.MemorySearch("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
}
