package sybase_storage

import (
	"context"
	"testing"

	"github.com/storage-lock/go-storage"
	"github.com/stretchr/testify/assert"
)

func TestSybaseStorage_GetName(t *testing.T) {
	s := &SybaseStorage{}
	assert.Equal(t, StorageName, s.GetName())
}

func TestSybaseStorage_Capabilities(t *testing.T) {
	s := &SybaseStorage{}
	caps := s.Capabilities()
	assert.Contains(t, caps, storage.CapabilityCAS)
	assert.Contains(t, caps, storage.CapabilityReliableTime)
}

func TestSybaseStorage_Close(t *testing.T) {
	s := &SybaseStorage{}
	err := s.Close(context.Background())
	assert.NoError(t, err)
}
