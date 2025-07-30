package biz

type UserRepo interface {
}

type UserUseCase struct {
	repo UserRepo
}
