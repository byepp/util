package cryptutil

import (
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypt(t *testing.T) {
	c := 10
	b := make([]byte, c)
	_, err := rand.Read(b)
	if !assert.Nil(t, err) {
		return
	}
	// The slice should now contain random bytes instead of only zeroes.
	assert.NotEqualValues(t, b, make([]byte, c))
}

func TestGenAES256Key(t *testing.T) {
	in := "我是中文测试"
	key := GenerateAES256CFBKey()
	out, err := AES256CFBEncrypt([]byte(in), key)
	if !assert.Nil(t, err, "加密失败") {
		return
	}
	dec, err := AES256CFBDecrypt(out, key)
	if !assert.Nil(t, err, "解密失败") {
		return
	}
	assert.Equal(t, string(dec), in, "解密后对比失败")
}

var (
	myAesKey = AES256CFBKey{
		0x93, 0x46, 0x0E, 0xD6, 0x90, 0xC4, 0x9C, 0x2A,
		0xCB, 0x74, 0x29, 0xB1, 0x6A, 0x23, 0xD0, 0xF3,
		0xA6, 0xAE, 0x98, 0xFA, 0xA5, 0x9E, 0x8D, 0xDD,
		0x50, 0x61, 0x73, 0x4E, 0x2C, 0xB5, 0xD3, 0x80,
	}
)

func Encrypt(text string) (string, error) {
	encrypted, err := AES256CFBEncrypt([]byte(text), myAesKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(encrypted), nil
}

func Decrypt(encrypted string) (string, error) {
	src, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	decrypted, err := AES256CFBDecrypt(src, myAesKey)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func TestAES256(t *testing.T) {
	in := "52e1459f-9c30-4222-86f5-9b407dd7c199"
	out, err := Encrypt(in)
	if !assert.Nil(t, err, "加密失败") {
		return
	}
	dec, err := Decrypt(out)
	if !assert.Nil(t, err, "解密失败") {
		return
	}
	assert.Equal(t, dec, in, "解密后对比失败")
}

func TestGen3DESECBKeyKey(t *testing.T) {
	in := "我是中文测试"
	key := Generate3DESECBKey()
	out, err := TripleDESECBEncrypt([]byte(in), key)
	if !assert.Nil(t, err, "加密失败") {
		return
	}
	dec, err := TripleDESECBDecrypt(out, key)
	if !assert.Nil(t, err, "解密失败") {
		return
	}
	assert.Equal(t, string(dec), in, "解密后对比失败")
}

func TestDESECBEncrypt(t *testing.T) {
	in := "我是中文"
	key := TripleDESECBKey{0xDC, 0xD3, 0x44, 0xCF, 0x9A, 0xA0, 0x45, 0x63, 0xD2, 0x27, 0xC9, 0x7A, 0x68, 0xCC, 0xB7, 0xA3, 0x97, 0x9C, 0x53, 0x65, 0x67, 0x27, 0xEA, 0xBB}
	out, err := TripleDESECBEncrypt([]byte(in), key)
	if !assert.Nil(t, err, "加密失败") {
		return
	}
	dec, err := TripleDESECBDecrypt(out, key)
	if !assert.Nil(t, err, "解密失败") {
		return
	}
	assert.Equal(t, string(dec), in, "解密后对比失败")
}

const lettersRunes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+-=`[]"

func RandStringRunes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = lettersRunes[rand.Intn(len(lettersRunes))]
	}
	return string(b)
}
