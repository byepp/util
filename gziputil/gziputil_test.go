package gziputil

import (
	"bytes"
	"encoding/base64"
	"math/rand"
	"testing"
)

func TestGzip(t *testing.T) {
	//out := GzipEncode([]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55})
	//input, err := GzipDecode(out)
	//t.Log(input, err)

	for i := 0; i < 100; i++ {
		var encode bytes.Buffer
		l := rand.Intn(1000000)
		for j := 0; j < l; j++ {
			encode.WriteByte(byte(rand.Intn(0xFF)))
		}
		encoded := GzipEncode(encode.Bytes())
		decoded, err := GzipDecode(encoded)
		if err != nil {
			t.Error("GzipDecode Failed", err)
			continue
		}
		if bytes.Compare(decoded, encode.Bytes()) != 0 {
			t.Error("Compare Failed", decoded, encode.Bytes())
			continue
		}
		t.Logf("Compare Success (%v) => (%v)\n", decoded, encode.Bytes())
	}
}

func TestZlib(t *testing.T) {
	out1 := ZlibEncode([]byte("empty"))
	t.Log(out1)
	in1, err := ZlibDecode(out1)
	t.Log(in1, err)

	out := ZlibEncode([]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99})
	t.Log(out)
	input, err := ZlibDecode(out)
	t.Log(input, err)
}

func TestEasyNetQ(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString("ImhlbGxvIg==")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	t.Logf("%s", b)

	b, err = base64.StdEncoding.DecodeString("eyJ3b3JsZCI6IndvcmxkMSJ9")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	t.Logf("%s", b)
}
