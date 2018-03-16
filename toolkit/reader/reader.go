package reader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

type MultipleReader interface {
	Reader() io.ReadCloser
}

type myMyltipleReader struct {
	data []byte
}

func (rr *myMyltipleReader) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(rr.data))
}

func NewMultipleReader(reader io.Reader) (MultipleReader, error) {
	var (
		data []byte
		err  error
	)
	if reader != nil {
		data, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("multiple reader: couldn't create a new one: %s", err)
		}
	} else {
		data = []byte{}
	}
	return &myMyltipleReader{
		data: data,
	}, nil
}
