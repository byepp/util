package jsonutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/byepp/util/fileutil"
)

func LoadFile(filePath string, e interface{}) error {
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, e)
}

func SaveFile(filePath string, e interface{}) error {
	if e == nil {
		e = struct{}{}
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(e)
	if err != nil {
		return err
	}
	return fileutil.FileWriteAll(filePath, buf.Bytes())
}

// 保存为带缩进格式的Json
func SaveIndentFile(filePath string, e interface{}) error {
	if e == nil {
		e = struct{}{}
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "\t")
	err := enc.Encode(e)
	if err != nil {
		return err
	}
	return fileutil.FileWriteAll(filePath, buf.Bytes())
}
