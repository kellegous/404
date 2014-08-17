package main

import (
  "crypto/rand"
  "crypto/sha256"
  "encoding/base64"
  "fmt"
  "io"
)

func main() {
  buf := make([]byte, sha256.BlockSize)

  if _, err := io.ReadFull(rand.Reader, buf); err != nil {
    panic(err)
  }

  fmt.Println(base64.StdEncoding.EncodeToString(buf))
}
