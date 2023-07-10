package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/vanillazen/stl/backend/internal/sys/errors"
)

type UUID struct {
	Val string
}

var (
	Nil = UUID{
		Val: "00000000-0000-0000-0000-000000000000",
	}
)

var (
	NotValidUUIDErr = errors.New("not a valid UUID")
)

// NewUUID generates a new UUID.
func NewUUID() UUID {
	uuid := make([]byte, 16)
	_, _ = rand.Read(uuid)

	// Set the version (4) and variant (RFC 4122) bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return UUID{Val: formatUUID(uuid)}
}

func Must() UUID {
	return NewUUID()
}

func Parse(s string) (UUID, error) {
	if !Validate(s) {
		return Nil, NotValidUUIDErr
	}
	return UUID{Val: s}, nil
}

func MustParse(s string) UUID {
	uuid, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return uuid
}

func (uid UUID) Equal(u2 UUID) bool {
	return uid.Val == u2.Val
}

func (uid UUID) Compare(u2 UUID) int {
	return strings.Compare(uid.Val, u2.Val)
}

func (uid UUID) Nil() bool {
	return uid.Val == "00000000-0000-0000-0000-000000000000"
}

func (uid UUID) String() string {
	return uid.Val
}

func (uid UUID) MarshalText() ([]byte, error) {
	return []byte(uid.Val), nil
}

func ParseBytes(input []byte) (UUID, error) {
	return Parse(string(input))
}

func Validate(uuid string) bool {
	trimmedUUID := strings.ReplaceAll(uuid, "-", "")

	if len(trimmedUUID) != 32 {
		return false
	}

	_, err := hex.DecodeString(trimmedUUID)
	return err == nil
}

func formatUUID(uuid []byte) string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func (uid *UUID) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		uid.Val = string(v)
	case string:
		uid.Val = v
	case nil:
		uid.Val = ""
	default:
		return fmt.Errorf("unsupported scan, storing driver.Value type %T into type UUID", value)
	}

	return nil
}
