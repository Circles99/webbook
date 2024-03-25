package ioc

import (
	"webbook/internal/service/sms"
	"webbook/internal/service/sms/tencent"
)

func InitSms() sms.Service {
	return tencent.NewService(nil, "", "")
}
