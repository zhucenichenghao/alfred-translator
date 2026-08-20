package provider

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"strconv"
)

type SaltSource func() (string, error)

func randomSalt() (string, error) {
	var data [8]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return strconv.FormatUint(binary.BigEndian.Uint64(data[:]), 10), nil
}

func md5Hex(text string) string {
	sum := md5.Sum([]byte(text))
	return hex.EncodeToString(sum[:])
}

func signature(key, query, salt, secret string) string {
	return md5Hex(key + query + salt + secret)
}
