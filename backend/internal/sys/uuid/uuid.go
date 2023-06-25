package uuid

import (
	"crypto/rand"
	"fmt"
	"strings"
)

type (
	UUID struct {
		Val string
	}
)

var (
	Nil = UUID{
		Val: "00000000-0000-0000-0000-000000000000",
	}
)

// NewUUID should return an additional error value.
// TODO: Implement additional error return value
func NewUUID(uuidStr string) UUID {
	ok := Validate(uuidStr)
	if !ok {
		return Nil
	}

	return UUID{Val: uuidStr}
}

func New() (uid UUID, err error) {
	uuid := make([]byte, 16)
	_, err = rand.Read(uuid)
	if err != nil {
		return uid, err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant RFC 4122

	return UUID{
		Val: formatUUID(uuid),
	}, nil
}

func formatUUID(uuid []byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func Validate(uuid string) bool {
	trimmedUUID := strings.ReplaceAll(uuid, "-", "")

	if len(trimmedUUID) != 32 {
		return false
	}

	_, err := fmt.Sscanf(trimmedUUID, "%x")
	return err == nil
}
