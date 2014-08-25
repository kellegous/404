package main

import (
  "crypto/rand"
  "crypto/sha256"
  "encoding/json"
  "fmt"
  "four04/config"
  "io"
)

func main() {
  sigKey := make([]byte, sha256.BlockSize)
  if _, err := io.ReadFull(rand.Reader, sigKey); err != nil {
    panic(err)
  }

  encKey := make([]byte, 32)
  if _, err := io.ReadFull(rand.Reader, encKey); err != nil {
    panic(err)
  }

  cfg := config.Config{
    HmacKey: sigKey,
    AesKey:  encKey,
    DbPath:  "dat",
  }

  b, err := json.MarshalIndent(&cfg, "", "  ")
  if err != nil {
    panic(err)
  }

  fmt.Printf("%s\n", b)
}
