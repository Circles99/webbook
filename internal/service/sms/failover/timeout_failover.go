package failover

import (
	"context"
	"sync/atomic"
	"webbook/internal/service/sms"
)

type TimeoutFailoverSmsService struct {
	// 你的服务商
	svcs []sms.Service

	idx int32

	// 连续超时的个数
	cnt int32

	// 阈值
	// 连续超过这个数字就要切换
	threshold int32
}

func NewTimeoutFailoverSmsService(svcs []sms.Service) sms.Service {
	return &TimeoutFailoverSmsService{
		svcs: svcs,
	}
}

func (t TimeoutFailoverSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)

	if cnt > t.threshold {
		// 这里要切换
		// 这里取余就是为了防止越界，下标溢出
		newIdx := idx + 1%int32(len(t.svcs))
		// 原子操作修改idx
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 成功修改idx
			atomic.StoreInt32(&t.cnt, 0)
		}

		// else 出现并发了

		idx = atomic.LoadInt32(&t.idx)
	}

	svc := t.svcs[idx]

	err := svc.Send(ctx, tplId, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return nil
	case nil:
		// 连续状态被打断
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		// 不知道什么错误，可以考虑换下一个
		// / 超时 可能是偶发的，尽量再试试
		// 非超市，直接下一个
		return err
	}
	return nil
}
