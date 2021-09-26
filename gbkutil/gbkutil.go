package gbkutil

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

// 查询GBK字符混合长度
func RuneLen(s string) (ret int) {
	ret = 0
	for _, c := range s {
		l := utf8.RuneLen(c)
		if l == 3 {// GBK
			ret += 2
		} else {
			ret += l
		}
	}
	return
}

// 左对齐不足补充空格
func PadLeftSpace(s string, width int) string {
	return PadLeft(s, width, " ")
}

// 左对齐不足补充指定字符串
func PadLeft(s string, width int, pad string) string {
	padWidth := (width-RuneLen(s)) / RuneLen(pad)
	if padWidth <= 0 {
		return s
	}
	return strings.Repeat(pad, padWidth) + s
}

// 右对齐不足补充空格
func PadRightSpace(s string, width int) string {
	return PadRight(s, width, " ")
}

// 右对齐不足补充指定字符串
func PadRight(s string, width int, pad string) string {
	padWidth := (width-RuneLen(s)) / RuneLen(pad)
	if padWidth <= 0 {
		return s
	}
	return s + strings.Repeat(pad, padWidth)
}

// 剧中对齐不足补充空格
func PadCenterSpace(s string, width int) string {
	return PadCenter(s, width, " ")
}

// 居中对齐不足补充指定字符串
func PadCenter(s string, width int, pad string) (ret string) {
	unitWidth := RuneLen(pad)
	sLen := RuneLen(s)
	leftSpaceWidth := (width - sLen) / unitWidth / 2
	rightSpaceWidth := (width-leftSpaceWidth*unitWidth-sLen)/unitWidth
	if leftSpaceWidth > 0 {
		ret += strings.Repeat(pad, leftSpaceWidth)
	}
	ret += s
	if rightSpaceWidth > 0 {
		ret += strings.Repeat(pad, rightSpaceWidth)
	}
	return
}

// GBK字符串转UTF8
func ToUTF8(gbkStr string) string {
	return BytesToUTF8([]byte(gbkStr))
}

// GBK字节数组转UTF8
func BytesToUTF8(gbkBytes []byte) string {
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(gbkBytes), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return ""
	}
	return string(data)
}

// UTF8字符串转GBK
func FromUTF8(utf8Str string) string {
	return FromUTF8Bytes([]byte(utf8Str))
}

// UTF8字节数组转GBK
func FromUTF8Bytes(utf8Bytes []byte) string {
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(utf8Bytes), simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		return ""
	}
	return string(data)
}
