package id

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/ncfex/dcart-auth/internal/application/ports/id"
)

type deterministicIDGenerator struct {
	namespace string
}

func NewDeterministicIDGenerator(namespace string) id.IDGenerator {
	return &deterministicIDGenerator{
		namespace: namespace,
	}
}

func (g *deterministicIDGenerator) GenerateFromData(data []byte) string {
	hasher := sha256.New()

	if g.namespace != "" {
		hasher.Write([]byte(g.namespace))
	}

	hasher.Write(data)

	return hex.EncodeToString(hasher.Sum(nil))[:32]
}
