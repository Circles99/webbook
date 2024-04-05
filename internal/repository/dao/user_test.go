package dao

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDao_Insert(t *testing.T) {

	tests := []struct {
		name    string
		mock    func() *sql.DB
		ctx     context.Context
		user    User
		wantErr error
	}{
		{
			name: "插入成功",
			mock: func() *sql.DB {
				mockDb, mock, err := sqlmock.New()

				res := sqlmock.NewResult(3, 1)

				// 这个*是正则表达式
				// 这个语法的意思就是 只要是insert 到users的语句
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnResult(res)
				require.NoError(t, err)
				return mockDb
			},
			ctx:     context.Background(),
			user:    User{},
			wantErr: nil,
		},

		{
			name: "邮箱冲突",
			mock: func() *sql.DB {
				mockDb, mock, err := sqlmock.New()

				// 这个*是正则表达式
				// 这个语法的意思就是 只要是insert 到users的语句
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(&mysql.MySQLError{Number: 1062})
				require.NoError(t, err)
				return mockDb
			},
			ctx:     context.Background(),
			user:    User{},
			wantErr: ErrUserDuplicate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open(gormmysql.New(gormmysql.Config{
				Conn: tt.mock(),
				// 初始化的时候select Version
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// mock Db 不需要pin
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			require.NoError(t, err)
			d := NewUserDao(db)
			err = d.Insert(tt.ctx, tt.user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
