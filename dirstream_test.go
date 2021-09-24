package dirstream

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestDirReaderTwoFiles(t *testing.T) {
	dir, err := os.MkdirTemp("", t.Name())
	requireNoErr(t, err)
	defer os.RemoveAll(dir)

	content1 := []byte("file 1 content\n")
	createFileWithContent(t, dir, "file1", content1)

	content2 := []byte("file 2 content\nmore file 2 content\n")
	createFileWithContent(t, dir, "file2", content2)

	dr, err := NewDirReader(dir)
	requireNoErr(t, err)
	defer dr.Close()

	expected := []byte("file 1 content\nfile 2 content\nmore file 2 content\n")
	var buf bytes.Buffer
	_, err = io.Copy(&buf, dr)
	requireNoErr(t, err)
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Errorf("want %s, got %s", string(expected), buf.String())
	}
	err = dr.Err()
	requireNoErr(t, err)
}

func createFileWithContent(t *testing.T, dir, name string, content []byte) {
	t.Helper()
	file, err := os.CreateTemp(dir, name)
	requireNoErr(t, err)
	defer file.Close()
	_, err = file.Write(content)
	requireNoErr(t, err)
}

func requireNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}
