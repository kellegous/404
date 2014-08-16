package main

import (
  "fmt"
  "github.com/DanielMorsing/rocksdb"
)

func main() {
  o := rocksdb.NewOptions()
  o.SetCreateIfMissing(true)

  db, err := rocksdb.Open("fake.db", o)
  if err != nil {
    panic(err)
  }

  fmt.Println(db)
}
