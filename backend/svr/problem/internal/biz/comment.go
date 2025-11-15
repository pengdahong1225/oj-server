package biz

import (
	"oj-server/svr/problem/internal/model"
)

// 仓库接口由data层去实现
type CommentRepo interface {
	// check
	// obj是否存在
	AssertObj(id int64) bool

	// query
	// 查询根评论列表
	QueryRootComment(objId int64, page, pageSize int) (int64, []model.Comment, error)
	// 查询子评论列表
	QueryChildComment(objId, rootId, rootCommentId int64, cursor int32) (int64, []model.Comment, error)

	// save
	// 保存根评论
	SaveRootComment(pbComment *model.Comment)
	// 保存子评论
	SaveChildComment(pbComment *model.Comment)

	// like
	CommentLike(objId, commentId int64)
}

type CommentUseCase struct {
	repo CommentRepo
}

func NewCommentUseCase(repo CommentRepo) *CommentUseCase {
	return &CommentUseCase{
		repo: repo,
	}
}
func (uc *CommentUseCase) AssertObj(id int64) bool {
	return uc.repo.AssertObj(id)
}
func (uc *CommentUseCase) SaveRootComment(comment *model.Comment) {
	uc.repo.SaveRootComment(comment)
}
func (uc *CommentUseCase) SaveChildComment(comment *model.Comment) {
	uc.repo.SaveChildComment(comment)
}
func (uc *CommentUseCase) QueryRootComment(objId int64, page, pageSize int) (int64, []model.Comment, error) {
	return uc.repo.QueryRootComment(objId, page, pageSize)
}
func (uc *CommentUseCase) QueryChildComment(objId, rootId, rootCommentId int64, cursor int32) (int64, []model.Comment, error) {
	return uc.repo.QueryChildComment(objId, rootId, rootCommentId, cursor)
}
func (uc *CommentUseCase) CommentLike(objId, commentId int64) {
	uc.repo.CommentLike(objId, commentId)
}
