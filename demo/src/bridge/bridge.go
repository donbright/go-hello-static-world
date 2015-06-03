package bridge

import "fmt"

// #cgo CFLAGS: -I/tmp/go-static-linking/include
// #cgo LDFLAGS: /tmp/go-static-linking/build/libgb.a
// #include <junk.h>
import "C"

func Run() {
  fmt.Printf("Invoking c library...\n")
  C.x(10)
  fmt.Printf("Done\n")
}
