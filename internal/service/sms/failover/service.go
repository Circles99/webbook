package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"webbook/internal/service/sms"
)

// 轮询Sms
type FailoverSmsService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSmsService(svcs []sms.Service) sms.Service {
	return &FailoverSmsService{svcs: svcs}
}

func (f FailoverSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)

		// 发送成功
		// 有个问题 永远打在头一个机器上
		if err == nil {
			return nil
		}

		//// 超时了
		//if err == context.DeadlineExceeded {
		//
		//}

		// 要做好监控
		log.Println(err)
	}

	return errors.New("全部服务商失败 ")
}

func (f FailoverSmsService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {

	// 我取下一个节点来作为起始节点
	// 利用取余打散，避免每次都从第一个节点发送
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))

	for i := idx; i < idx+length; i++ {
		svc := f.svcs[int(i%length)]
		err := svc.Send(ctx, tplId, args, numbers...)

		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			// 超时了
			return err
		default:
			log.Println(err)
		}
	}

	return errors.New("全部服务商失败 ")
}
