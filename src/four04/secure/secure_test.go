package secure

import (
  "bytes"
  "testing"
)

func fillUp(b []byte) {
  for i, n := 0, len(b); i < n; i++ {
    b[i] = byte(i)
  }
}

func TestSignAndVerify(t *testing.T) {
  var key [32]byte
  fillUp(key[:])

  msgs := [][]byte{
    []byte{1},
    []byte{2},
    []byte("test message"),
    []byte{0, 0, 0, 0},
  }

  for _, msg := range msgs {
    buf, err := Sign(msg, key[:])
    if err != nil {
      t.Error(err)
    }

    org, _, err := Verify(buf, key[:])
    if err != nil {
      t.Error(err)
    }

    if !bytes.Equal(msg, org) {
      t.Errorf("message not restored through verify: %v vs %v", msg, org)
    }
  }
}
