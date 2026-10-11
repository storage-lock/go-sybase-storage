package sybase_storage

import (
	"context"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
)

// SybaseSqlProvider Sybase（ASE）数据库的 SQL 方言提供者
// Sybase 不支持 MySQL 的 UNIX_TIMESTAMP(NOW())。
type SybaseSqlProvider struct {
	*sql_based_storage.Sql92Provider
}

func NewSybaseSqlProvider() *SybaseSqlProvider {
	return &SybaseSqlProvider{
		Sql92Provider: sql_based_storage.NewSql92Provider(),
	}
}

// NowTimestampSql Sybase 获取当前时间的 Unix 时间戳
// 修复：原来直接继承 Sql92Provider 的 "SELECT UNIX_TIMESTAMP(NOW())"，Sybase 不支持这两个
// MySQL 专有函数，执行即报语法错误，导致 GetTime 完全不可用。用 Sybase 的 DATEDIFF(ss, ...)
// 函数换算 1970-01-01 到 GETDATE() 的秒数差作为 Unix 时间戳。
// 注意：未接入真实 Sybase 驱动验证，是按 Sybase ASE SQL 官方文档做的最佳努力实现。
func (x *SybaseSqlProvider) NowTimestampSql(ctx context.Context, tableFullName string) (string, []any) {
	return `SELECT DATEDIFF(ss, '19700101', GETDATE())`, nil
}
