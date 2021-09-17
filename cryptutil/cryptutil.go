package cryptutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"math/big"
)

const (
	// RSA default E
	defaultExponent = 65537
)

var (
	ErrWrongInputLength = errors.New("helpers/crypt: wrong input length")

	ErrWrongIVSize = errors.New("helpers/crypt: wrong IV length. IV length must equal block size")

	ErrWrongInputParameter = errors.New("helpers/crypt: wrong intput parameter")
)

type (
	AES256CFBKey    [32]byte
	TripleDESECBKey [24]byte
)

func CalcMD5(input string) string {
	output := md5.Sum([]byte(input))
	return hex.EncodeToString(output[:])
}

func AES256CFBEncrypt(input []byte, key AES256CFBKey) ([]byte, error) {
	var iv = key[:aes.BlockSize]
	encrypted := make([]byte, len(input))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted, input)
	return encrypted, nil
}

func AES256CFBDecrypt(encrypted []byte, key AES256CFBKey) ([]byte, error) {
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	var iv = key[:aes.BlockSize]
	decrypted := make([]byte, len(encrypted))
	var block cipher.Block
	block, err = aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(decrypted, encrypted)
	return decrypted, nil
}

// 使用PKCS7填充
func PKCS7Padding(in []byte, blockSize int) []byte {
	padding := blockSize - len(in)%blockSize
	buffer := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(in, buffer...)
}

// 去除PKCS7填充
func PKCS7Trimming(in []byte) ([]byte, error) {
	if len(in) == 0 {
		return nil, ErrWrongInputLength
	}

	padding := in[len(in)-1]
	if int(padding) > len(in) {
		return nil, ErrWrongInputLength
	}
	return in[:len(in)-int(padding)], nil
}

// AESCBCEncrypt AESCBC加密, 使用PKCS7填充
func AESCBCEncrypt(in, key, iv []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != c.BlockSize() {
		return nil, ErrWrongIVSize
	}

	in = PKCS7Padding(in, c.BlockSize())
	out := make([]byte, len(in))

	encrypter := cipher.NewCBCEncrypter(c, iv)
	encrypter.CryptBlocks(out, in)

	return out, nil
}

// AESCBCDecrypt AESCBC解密, 并去除PKCS7填充
func AESCBCDecrypt(in, key, iv []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != c.BlockSize() {
		return nil, ErrWrongIVSize
	}

	if len(in)%c.BlockSize() != 0 {
		return nil, ErrWrongInputLength
	}

	out := make([]byte, len(in))

	decrypter := cipher.NewCBCDecrypter(c, iv)
	decrypter.CryptBlocks(out, in)

	out, err = PKCS7Trimming(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// 3DESECB加密, 使用PKCS7填充
func EasyTripleDESECBEncrypt(in []byte, passphrase string) ([]byte, error) {
	keyBytes := md5.Sum([]byte(passphrase))
	var key [32]byte
	hex.Encode(key[:], keyBytes[:])
	c, err := des.NewTripleDESCipher(key[:24])
	if err != nil {
		return nil, err
	}

	in = PKCS7Padding(in, c.BlockSize())
	var out []byte
	dst := make([]byte, c.BlockSize())
	for i := 0; i < len(in); i += c.BlockSize() {
		src := in[i : i+c.BlockSize()]
		c.Encrypt(dst, src)

		out = append(out, dst...)
	}

	return out, nil
}

// 3DESECB解密, 并去除PKCS7填充
func EasyTripleDESECBDecrypt(in []byte, passphrase string) ([]byte, error) {
	keyBytes := md5.Sum([]byte(passphrase))
	var key [32]byte
	hex.Encode(key[:], keyBytes[:])
	c, err := des.NewTripleDESCipher(key[:24])
	if err != nil {
		return nil, err
	}

	if len(in)%c.BlockSize() != 0 {
		return nil, ErrWrongInputLength
	}

	var out []byte
	dst := make([]byte, c.BlockSize())
	for i := 0; i < len(in); i += c.BlockSize() {
		src := in[i : i+c.BlockSize()]
		c.Decrypt(dst, src)

		out = append(out, dst...)
	}

	out, err = PKCS7Trimming(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// 3DESECB加密, 使用PKCS7填充
func TripleDESECBEncrypt(in []byte, key TripleDESECBKey) ([]byte, error) {
	c, err := des.NewTripleDESCipher(key[:])
	if err != nil {
		return nil, err
	}

	in = PKCS7Padding(in, c.BlockSize())
	var out []byte
	dst := make([]byte, c.BlockSize())
	for i := 0; i < len(in); i += c.BlockSize() {
		src := in[i : i+c.BlockSize()]
		c.Encrypt(dst, src)

		out = append(out, dst...)
	}

	return out, nil
}

// 3DESECB解密, 并去除PKCS7填充
func TripleDESECBDecrypt(in []byte, key TripleDESECBKey) ([]byte, error) {
	c, err := des.NewTripleDESCipher(key[:])
	if err != nil {
		return nil, err
	}

	if len(in)%c.BlockSize() != 0 {
		return nil, ErrWrongInputLength
	}

	var out []byte
	dst := make([]byte, c.BlockSize())
	for i := 0; i < len(in); i += c.BlockSize() {
		src := in[i : i+c.BlockSize()]
		c.Decrypt(dst, src)

		out = append(out, dst...)
	}

	out, err = PKCS7Trimming(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// TripleDESCBCEncrypt 3DESCBC加密, 使用PKCS7填充
func TripleDESCBCEncrypt(in, key, iv []byte) ([]byte, error) {
	c, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != c.BlockSize() {
		return nil, ErrWrongIVSize
	}

	in = PKCS7Padding(in, c.BlockSize())
	out := make([]byte, len(in))

	encrypter := cipher.NewCBCEncrypter(c, iv)
	encrypter.CryptBlocks(out, in)

	return out, nil
}

// TripleDESCBCDecrypt 3DESCBC解密, 并去除PKCS7填充
func TripleDESCBCDecrypt(in, key, iv []byte) ([]byte, error) {
	c, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != c.BlockSize() {
		return nil, ErrWrongIVSize
	}

	if len(in)%c.BlockSize() != 0 {
		return nil, ErrWrongInputLength
	}

	out := make([]byte, len(in))

	decrypter := cipher.NewCBCDecrypter(c, iv)
	decrypter.CryptBlocks(out, in)

	out, err = PKCS7Trimming(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// DESCBCEncrypt DESCBC加密, 使用PKCS7填充
func DESCBCEncrypt(in, key, iv []byte) ([]byte, error) {
	c, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != c.BlockSize() {
		return nil, ErrWrongIVSize
	}

	in = PKCS7Padding(in, c.BlockSize())
	out := make([]byte, len(in))

	encrypter := cipher.NewCBCEncrypter(c, iv)
	encrypter.CryptBlocks(out, in)

	return out, nil
}

// DESCBCDecrypt DESCBC解密, 并去除PKCS7填充
func DESCBCDecrypt(in, key, iv []byte) ([]byte, error) {
	c, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != c.BlockSize() {
		return nil, ErrWrongIVSize
	}

	if len(in)%c.BlockSize() != 0 {
		return nil, ErrWrongInputLength
	}

	out := make([]byte, len(in))

	decrypter := cipher.NewCBCDecrypter(c, iv)
	decrypter.CryptBlocks(out, in)

	out, err = PKCS7Trimming(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// PublicKeyFromBytes 通过字节集生成公钥, e为0时默认使用65537
func PublicKeyFromBytes(n []byte, e int) *rsa.PublicKey {
	if e == 0 {
		e = defaultExponent
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(n),
		E: e,
	}
}

// PublicKeyFromString 通过字符串生成公钥, e为0时默认使用65537
func PublicKeyFromString(n string, e int, base int) (*rsa.PublicKey, error) {
	i, ok := new(big.Int).SetString(n, base)
	if !ok {
		return nil, ErrWrongInputParameter
	}

	if e == 0 {
		e = defaultExponent
	}

	return &rsa.PublicKey{
		N: i,
		E: e,
	}, nil
}

// RSAEncryptPKCS1v15 RSA加密, 并用PKCS1, V1.5填充
func RSAEncryptPKCS1v15(in []byte, key *rsa.PublicKey) (data []byte, err error) {
	k := ((key.N.BitLen() + 7) / 8) - 11
	for {
		if len(in) > k {
			arr := in[:k]
			arr, err := rsa.EncryptPKCS1v15(rand.Reader, key, arr)
			if err != nil {
				return nil, err
			}

			data = append(data, arr...)
			in = in[k:]
		} else {
			arr, err := rsa.EncryptPKCS1v15(rand.Reader, key, in)
			if err != nil {
				return nil, err
			}

			data = append(data, arr...)

			break
		}
	}

	return data, nil
}

// RSAEncrypt RSA加密
func RSAEncrypt(in []byte, key []byte) ([]byte, error) {
	pub := PublicKeyFromBytes(key, 65537)
	return RSAEncryptPKCS1v15(in, pub)
}

// RSAEncryptNoPadding RSA使用非填充方式加密
func RSAEncryptNoPadding(input []byte, key *rsa.PublicKey) []byte {
	c := new(big.Int)
	c.Exp(new(big.Int).SetBytes(input), big.NewInt(int64(key.E)), key.N)
	return c.Bytes()
}

func GenerateAES256CFBKey() (ret AES256CFBKey) {
	_, _ = rand.Read(ret[:])
	return ret
}

func Generate3DESECBKey() (ret TripleDESECBKey) {
	_, _ = rand.Read(ret[:])
	return ret
}
