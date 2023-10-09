# LazyMultiReader

*LazyMultiReader* is an alternative package for *io.MultiReader*, where the composing reader streams are only opened when they are read from. It can be used when the underlying readers are prone to timeouts, especially for socket/network connections.

## Installation

```
go get -u github.com/shivanthzen/lazymultireader
```

## Usage

```
package main

import (
	"io"
	"log"
	"os"
	"strings"
)



type stringReadOpener struct {
	strings.Reader
}

func (b *stringReadOpener) Open() error {
	// nop

	return nil
}

func NewStringReadOpener(s string) *stringReadOpener {
	return &stringReadOpener{
		*strings.NewReader(s),
	}
}

func main() {
	r1 := NewStringReadOpener("first reader ")
	r2 := NewStringReadOpener("second reader ")
	r3 := NewStringReadOpener("third reader\n")
	r := lazymultireader.NewLazyMultiReader(r1, r2, r3)

	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}

}

```
