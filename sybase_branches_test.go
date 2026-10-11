package sybase_storage

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
)

type stubProvider struct{ *sql_based_storage.Sql92Provider }

func (stubProvider) CreateTableSql(context.Context, string) (string, []any) {
	return "SELECT 1", nil
}

func TestSybaseSqlProvider(t *testing.T) {
	sqlText, args := NewSybaseSqlProvider().NowTimestampSql(context.Background(), "locks")
	assert.Contains(t, sqlText, "DATEDIFF")
	assert.Nil(t, args)
}

func TestNewSybaseStorage_NilOptions(t *testing.T) {
	assert.Panics(t, func() { _, _ = NewSybaseStorage(context.Background(), nil) })
}

func TestNewSybaseStorage_InvalidOptions(t *testing.T) {
	_, err := NewSybaseStorage(context.Background(), &sql_based_storage.SqlBasedStorageOptions{})
	assert.Error(t, err)
}

func TestNewSybaseStorage_InitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec("CREATE TABLE").WillReturnError(errors.New("init failed"))
	_, err = NewSybaseStorage(context.Background(), sql_based_storage.NewSqlBasedStorageOptions().
		SetConnectionManager(storage.NewFixedSqlDBConnectionManager(db)))
	assert.EqualError(t, err, "init failed")
}

func TestNewSybaseStorage_CustomProvider(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(0, 0))
	s, err := NewSybaseStorage(context.Background(), sql_based_storage.NewSqlBasedStorageOptions().
		SetConnectionManager(storage.NewFixedSqlDBConnectionManager(db)).
		SetSqlProvider(stubProvider{Sql92Provider: sql_based_storage.NewSql92Provider()}))
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestSybaseStorage_CreateWithVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
	s, err := NewSybaseStorage(context.Background(), sql_based_storage.NewSqlBasedStorageOptions().
		SetConnectionManager(storage.NewFixedSqlDBConnectionManager(db)))
	assert.NoError(t, err)
	info := &storage.LockInformation{OwnerId: "owner"}

	mock.ExpectExec("INSERT INTO").WillReturnError(errors.New("Attempt to insert duplicate key row"))
	assert.ErrorIs(t, s.CreateWithVersion(context.Background(), "lock", 1, info), storage_lock.ErrLockAlreadyExists)
	mock.ExpectExec("INSERT INTO").WillReturnError(errors.New("duplicate key"))
	assert.ErrorIs(t, s.CreateWithVersion(context.Background(), "lock", 1, info), storage_lock.ErrLockAlreadyExists)
	mock.ExpectExec("INSERT INTO").WillReturnError(errors.New("other"))
	assert.EqualError(t, s.CreateWithVersion(context.Background(), "lock", 1, info), "other")
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	assert.NoError(t, s.CreateWithVersion(context.Background(), "lock", 1, info))
}

func TestSybaseStorage_Watch_EmptyLockID(t *testing.T) {
	_, err := (&SybaseStorage{}).Watch(context.Background(), "")
	assert.Error(t, err)
}
