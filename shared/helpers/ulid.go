package helpers

import (
	"time"

	"math/rand"

	"github.com/oklog/ulid/v2"
)

func GenerateULID() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
