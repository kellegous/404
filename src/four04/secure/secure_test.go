package secure

import (
  "bytes"
  "crypto/aes"
  "crypto/sha256"
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

func testEncryptAndDecrypt(t *testing.T) {
  sigKey := make([]byte, sha256.BlockSize)
  encKey := make([]byte, aes.BlockSize)
  fillUp(sigKey)
  fillUp(encKey)

  msgs := [][]byte{
    []byte{1},
    []byte{2},
    []byte("test message"),
    []byte{0, 0, 0, 0},
  }

  for _, msg := range msgs {
    buf, err := Encrypt(msg, encKey, sigKey)
    if err != nil {
      t.Error(err)
      continue
    }

    org, _, err := Decrypt(buf, encKey, sigKey)
    if err != nil {
      t.Error(err)
      continue
    }

    if !bytes.Equal(msg, org) {
      t.Error("message not restored through decrypt: %v vs %v", msg, org)
    }
  }
}
