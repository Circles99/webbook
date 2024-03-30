package memory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemory_Set(t *testing.T) {

	tests := []struct {
		biz     string
		phone   string
		code    string
		name    string
		wantErr bool
	}{
		{
			biz:     "login",
			phone:   "12345",
			code:    "666",
			name:    "添加成功",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemory()
			if err := m.Set(context.Background(), tt.biz, tt.phone, tt.code); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_Verify(t *testing.T) {
	m := NewMemory()
	tests := []struct {
		biz     string
		phone   string
		code    string
		name    string
		before  func(t *testing.T)
		wantOk  bool
		wantErr error
	}{
		{
			biz:     "login",
			phone:   "12345",
			code:    "666",
			name:    "验证成功",
			wantErr: nil,
			wantOk:  true,
			before: func(t *testing.T) {
				err := m.Set(context.Background(), "login", "12345", "666")
				require.NoError(t, err)
			},
		},
		{
			biz:     "login",
			phone:   "12345",
			code:    "666",
			name:    "验证失败",
			wantErr: nil,
			wantOk:  false,
			before: func(t *testing.T) {
				err := m.Set(context.Background(), "login", "12345", "777")
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.before(t)
			ok, err := m.Verify(context.Background(), tt.biz, tt.phone, tt.code)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantOk, ok)

		})
	}
}
