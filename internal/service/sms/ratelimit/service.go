package ratelimit

import (
	"context"
	"fmt"
	"webbook/internal/service/sms"
	"webbook/pkg/ratelimit"
)

type RatelimitSmsService struct {
	// 被装饰的
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSmsService(service sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSmsService{
		svc:     service,
		limiter: limiter,
	}
}

func (s RatelimitSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 使用装饰器模式对已有代码进行功能上的添加
	// 防止侵入式代码到Service的send种
	limit, err := s.limiter.Limit(ctx, "send_key")
	if err != nil {
		// 系统错误，可以限流，也可以不限流，根据下游考虑
		return err
	}

	if limit {
		return fmt.Errorf("触发了限流")
	}

	err = s.svc.Send(ctx, tplId, args, numbers...)
	return err
}
