package biz

import (
	"oj-server/module/db/model"
)

// 仓库接口由data层去实现
type UserRepo interface {
	// 创建用户
	CreateNewUser(user *model.UserInfo) (int64, error)
	// 获取用户信息
	GetUserInfoByMobile(mobile int64) (*model.UserInfo, error)
	// 重置用户密码
	ResetUserPassword(mobile int64, password string) error
}

type UserUseCase struct {
	repo UserRepo
}

func NewUserUseCase(repo UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (uc *UserUseCase) CreateNewUser(user *model.UserInfo) (int64, error) {
	return uc.repo.CreateNewUser(user)
}
func (uc *UserUseCase) GetUserInfoByMobile(mobile int64) (*model.UserInfo, error) {
	return uc.repo.GetUserInfoByMobile(mobile)
}
func (uc *UserUseCase) ResetUserPassword(mobile int64, password string) error {
	return uc.repo.ResetUserPassword(mobile, password)
}
