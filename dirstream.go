package dirstream

import (
	"io"
	"io/fs"
	"os"
)

// DirReader creates a reader that combines all
// the files in a directory with just one reader
type DirReader struct {
	fs          fs.FS
	currentFile io.ReadCloser
	files       chan string
	err         error
}

// NewDirReader create a new instance of dir reader
// using a path to file system directory
func NewDirReader(dir string) (*DirReader, error) {
	return NewFSReader(os.DirFS(dir))
}

// NewFSReader create a new instance of dir reader
// using a file system implementation
func NewFSReader(fileSystem fs.FS) (*DirReader, error) {
	dr := DirReader{
		fs:    fileSystem,
		files: make(chan string),
	}
	go func() {
		defer close(dr.files)
		err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.Type().IsRegular() {
				dr.files <- path
			}
			return nil
		})
		if err != nil {
			dr.err = err
		}
	}()
	err := dr.openNextFile()
	return &dr, err
}

func (dr *DirReader) Read(p []byte) (n int, err error) {
	n, err = dr.currentFile.Read(p)
	if err == io.EOF {
		dr.currentFile.Close()
		err = dr.openNextFile()
	}
	return n, err
}

// Close closes currently opened file
func (dr *DirReader) Close() error {
	return dr.currentFile.Close()
}

// Err returns an error that may happen while walking the file system
func (dr *DirReader) Err() error {
	return dr.err
}

func (dr *DirReader) openNextFile() error {
	path, ok := <-dr.files
	if !ok {
		return io.EOF
	}
	var err error
	dr.currentFile, err = dr.fs.Open(path)
	return err
}
