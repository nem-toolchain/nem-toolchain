package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidChainId(t *testing.T) {
	assert.True(t, IsValidChainId(TestnetId))
	assert.True(t, IsValidChainId(MainnetId))
	assert.True(t, IsValidChainId(MijinId))

	assert.False(t, IsValidChainId(0x10))
	assert.False(t, IsValidChainId(0x00))
}
