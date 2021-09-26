package yamlutil

import (
	"bytes"
	"io/ioutil"

	"github.com/byepp/util/fileutil"

	"gopkg.in/yaml.v3"
)

func LoadFile(filePath string, e interface{}) error {
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, e)
}

func SaveFile(filePath string, e interface{}) error {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	err := enc.Encode(e)
	if err != nil {
		return err
	}
	return fileutil.FileWriteAll(filePath, buf.Bytes())
}
