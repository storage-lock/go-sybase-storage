package sybase_storage

import (
	"context"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
)

// SybaseStorage 基于 Sybase 的存储引擎实现
type SybaseStorage struct {
	*sql_based_storage.SqlBasedStorage
}

func NewSybaseStorage(ctx context.Context, options *sql_based_storage.SqlBasedStorageOptions) (*SybaseStorage, error) {
	baseStorage, err := sql_based_storage.NewSqlBasedStorage(options)
	if err != nil {
		return nil, err
	}
	s := &SybaseStorage{SqlBasedStorage: baseStorage}
	if err := s.Init(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

const StorageName = "sybase-storage"

func (x *SybaseStorage) GetName() string {
	return StorageName
}

func (x *SybaseStorage) Capabilities() []storage.StorageCapability {
	return []storage.StorageCapability{
		storage.CapabilityCAS,
		storage.CapabilityReliableTime,
	}
}
