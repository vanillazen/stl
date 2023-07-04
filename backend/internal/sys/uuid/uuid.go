package uuid

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/vanillazen/stl/backend/internal/sys/errors"
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

var (
	NotValidUUIDErr = errors.NewError("not a valid UUID")
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

func (uid UUID) String() string {
	if !Validate(uid.Val) {
		return Nil.Val
	}

	return uid.Val
}

func formatUUID(uuid []byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func Parse(uuidStr string) (uid UUID, err error) {
	if !Validate(uuidStr) {
		return Nil, NotValidUUIDErr
	}

	return UUID{Val: uuidStr}, nil
}

func Validate(uuid string) bool {
	trimmedUUID := strings.ReplaceAll(uuid, "-", "")

	if len(trimmedUUID) != 32 {
		return false
	}

	_, err := fmt.Sscanf(trimmedUUID, "%x")
	return err == nil
}
