package failover

import (
	"context"
	"webbook/internal/service/sms"
)

type TimeoutFailoverSmsService struct {
	// 你的服务商
	svcs []sms.Service
	// 连续超时的个数
	cnt int32
}

func NewTimeoutFailoverSmsService() sms.Service {

}

func (t TimeoutFailoverSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {

}
