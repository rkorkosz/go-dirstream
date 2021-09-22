[![Go](https://github.com/rkorkosz/go-dirstream/actions/workflows/go.yml/badge.svg)](https://github.com/rkorkosz/go-dirstream/actions/workflows/go.yml)

# DirStream

## What DirStream is?

DirStream is an `io.ReadCloser` implementation that takes a directory as an entry point and iterates over every file in that directory. Each file is read and streamed as one continuous byte stream.

## Usage

If you have logs in multiple files under `logs` directory you can print all logs to stdout by doing this:

```go
package main

import dirstream "github.com/rkorkosz/go-dirstream"

func main() {
    ds, err := dirstream.NewDirReader("logs")
    if err != nil {
        panic(err)
    }
    defer ds.Close()
    // print logs from all files in stdout
    io.Copy(os.Stdout, ds)
}
```

You can also use an implementation of `fs.FS` interface to create `DirStream` by using `NewFSReader` function.
