package sybase_storage

import (
	"context"
	"strings"
	"time"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
)

// SybaseStorage 基于 Sybase 的存储引擎实现
type SybaseStorage struct {
	*sql_based_storage.SqlBasedStorage
}

func NewSybaseStorage(ctx context.Context, options *sql_based_storage.SqlBasedStorageOptions) (*SybaseStorage, error) {
	// 如果调用方没有显式指定方言（还是默认的 Sql92Provider），换成 Sybase 自己的方言
	// 提供者：Sql92Provider.NowTimestampSql 是 MySQL 专有语法，在 Sybase 上会直接报语法错误
	if _, ok := options.SqlProvider.(*sql_based_storage.Sql92Provider); ok || options.SqlProvider == nil {
		options.SetSqlProvider(NewSybaseSqlProvider())
	}
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

// CreateWithVersion Sybase 唯一约束冲突时返回 ErrLockAlreadyExists
// 注意：未接入真实 Sybase 驱动验证，错误消息按 Sybase ASE 官方文档做的最佳努力匹配。
func (x *SybaseStorage) CreateWithVersion(ctx context.Context, lockId string, version storage.Version, lockInformation *storage.LockInformation) error {
	err := x.SqlBasedStorage.CreateWithVersion(ctx, lockId, version, lockInformation)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "Attempt to insert duplicate key row") || strings.Contains(msg, "duplicate key") {
			return storage_lock.ErrLockAlreadyExists
		}
	}
	return err
}

const StorageName = "sybase-storage"

func (x *SybaseStorage) GetName() string {
	return StorageName
}

func (x *SybaseStorage) Capabilities() []storage.StorageCapability {
	return []storage.StorageCapability{
		storage.CapabilityCAS,
		storage.CapabilityReliableTime,
		storage.CapabilityWatch,
	}
}

func (x *SybaseStorage) Watch(ctx context.Context, lockId string) (<-chan *storage.LockInformation, error) {
	return storage.NewPollingLockWatcher(x, 500*time.Millisecond).Watch(ctx, lockId)
}
