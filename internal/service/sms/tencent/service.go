package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func NewService(client *sms.Client, appId, signName string) *Service {
	return &Service{
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
		client:   client,
	}
}

func (s Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](tplId)
	req.PhoneNumberSet = slice.Map[string, *string](numbers, func(idx int, src string) *string {
		return &src
	})

	req.TemplateParamSet = slice.Map[string, *string](args, func(idx int, src string) *string {
		return &src
	})

	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}

	// 多个手机号有多个状态，有些可能成功，有些可能没成功
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送信息失败 %s, %s", *status.Code, *status.Message)
		}
	}
	return nil
}
