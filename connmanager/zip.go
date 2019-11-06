package connmanager

import (
	"compress/gzip"
	"bytes"
)

func GzipEncode(in []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	writer := gzip.NewWriter(buf)
	defer writer.Close()

	_, err := writer.Write(in)
	if err != nil {
		return nil, err
	}
	writer.Flush()
	return buf.Bytes(), nil
}