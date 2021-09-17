package gziputil

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"io/ioutil"
)

func GzipEncode(in []byte) []byte {
	var buffer bytes.Buffer
	w := gzip.NewWriter(&buffer)
	defer w.Close()
	_, err := w.Write(in)
	if err != nil {
		return nil
	}
	err = w.Flush()
	if err != nil {
		return nil
	}

	return buffer.Bytes()
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	ret, err := ioutil.ReadAll(reader)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	return ret, nil
}

func ZlibEncode(in []byte) []byte {
	var buffer bytes.Buffer
	w, _ := zlib.NewWriterLevel(&buffer, zlib.BestCompression)
	w.Write(in)
	w.Flush()
	defer w.Close()

	return buffer.Bytes()
}

func ZlibDecode(in []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	ret := new(bytes.Buffer)
	_, err = io.Copy(ret, r)
	//return ioutil.ReadAll(r)

	return ret.Bytes(), err
}
