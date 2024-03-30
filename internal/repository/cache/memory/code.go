package memory

import (
	"context"
	"fmt"
	"sync"
)

type Memory struct {
	sMap *sync.Map
}

func NewMemory() *Memory {
	return &Memory{
		sMap: &sync.Map{},
	}
}

func (m Memory) Set(ctx context.Context, biz, phone, code string) error {
	m.sMap.Store(m.getKey(biz, phone), code)

	return nil
}

func (m Memory) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {

	code, ok := m.sMap.Load(m.getKey(biz, phone))
	if !ok {
		// 没有数据
		return false, nil
	}

	if inputCode != code {
		return false, nil
	}

	return true, nil

}

func (m Memory) getKey(biz, phone string) string {
	return fmt.Sprintf("phone:code:%s:%s", biz, phone)
}
