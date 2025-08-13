package biz

import (
	"oj-server/module/db/model"
)

// 仓库接口由data层去实现
type ProblemRepo interface {
	// 创建题目
	CreateProblem(problem *model.Problem) (int64, error)
	// 查询题目列表，偏移量分页
	QueryProblemList(page, pageSize int, keyword, tag string) (int64, []model.Problem, error)
	// 查询题目数据
	QueryProblemData(id int64) (*model.Problem, error)
	// 更新题目
	UpdateProblem(problem *model.Problem) error
	// 更新题目状态
	UpdateProblemStatus(id int64, st int32) error
	// 删除题目
	DeleteProblem(id int64) error
	// 查询标签列表
	QueryTagList() ([]string, error)
}

type ProblemUseCase struct {
	repo ProblemRepo
}

func NewProblemUseCase(repo ProblemRepo) *ProblemUseCase {
	return &ProblemUseCase{
		repo: repo,
	}
}

func (pc *ProblemUseCase) CreateProblem(problem *model.Problem) (int64, error) {
	return pc.repo.CreateProblem(problem)
}
func (pc *ProblemUseCase) QueryProblemList(page, pageSize int, keyword, tag string) (int64, []model.Problem, error) {
	return pc.repo.QueryProblemList(page, pageSize, keyword, tag)
}
func (pc *ProblemUseCase) QueryProblemData(id int64) (*model.Problem, error) {
	return pc.repo.QueryProblemData(id)
}
func (pc *ProblemUseCase) UpdateProblem(problem *model.Problem) error {
	return pc.repo.UpdateProblem(problem)
}
func (pc *ProblemUseCase) UpdateProblemStatus(id int64, st int32) error {
	return pc.repo.UpdateProblemStatus(id, st)
}
func (pc *ProblemUseCase) DeleteProblem(id int64) error {
	return pc.repo.DeleteProblem(id)
}
func (pc *ProblemUseCase) QueryTagList() ([]string, error) {
	return pc.repo.QueryTagList()
}
