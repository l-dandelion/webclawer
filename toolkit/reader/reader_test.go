package reader

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestReaderNew(t *testing.T) {
	expectedData := "8dS7DkH7D"
	rr, err := NewMultipleReader(strings.NewReader(expectedData))
	if err != nil {
		t.Fatalf("An error occurs when new multiple reader: %s", err)
	}
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, rr.Reader())
	if err != nil {
		t.Fatalf("An error occurs when copy multiple reader: %s", err)
	}
	content1 := buffer.String()
	if content1 != expectedData {
		t.Fatalf("Inconsistent data: expected: %s, actual: %s",
			expectedData, content1)
	}
	content2 := buffer.String()
	if content2 != expectedData {
		t.Fatalf("Inconsistent data: expected: %s, actual: %s",
			expectedData, content2)
	}
}
