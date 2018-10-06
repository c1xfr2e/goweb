package uuid

import (
	"encoding/hex"

	"github.com/pborman/uuid"
)

// New new uuid(v4).
//
// Return string is the hex format of uuid byte array.
func New() string {
	return hex.EncodeToString(uuid.NewRandom())
}
