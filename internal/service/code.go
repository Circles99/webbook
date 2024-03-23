package service

import (
	"context"
	"fmt"
	"math/rand"
	"webbook/internal/repository"
	"webbook/internal/service/sms"
)

const codeTplId = "1877555"

var ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
var ErrSetCodeTooMany = repository.ErrSetCodeTooMany

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{repo: repo, smsSvc: smsSvc}
}

func (u *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := u.generateCode()
	err := u.repo.Store(ctx, biz, phone, code)
	if err != nil {
		// 有问题
		return err
	}
	// 发送出去
	return u.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
}

func (u *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return u.repo.Verify(ctx, biz, phone, inputCode)
}

func (u *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	// 不够6位的，加上前导0
	return fmt.Sprintf("%06d", num)
}
