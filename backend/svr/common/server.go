package common

type IServer interface {
	Init() error
	Run()
	Stop()
}
