package biz

import "oj-server/module/proto/pb"

// 仓库接口由data层去实现
type CommentRepo interface {
	// check
	// obj是否存在
	AssertObj(id int64) bool
	// 根评论是否存在，校验root_id和root_comment_id
	AssertRoot(rootCommentId int64, rootId int64) bool
	// 回复评论是否存在，校验reply_id和reply_comment_id
	AssertReply(replyCommentId int64, replyId int64) bool

	// query
	// 查询根评论列表
	QueryRootComment(page, pageSize int) (int32, []pb.Comment, error)

	// save
	// 保存根评论
	SaveRootComment(pbComment *pb.Comment)
	// 保存子评论
	SaveChildComment(pbComment *pb.Comment)
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
func (uc *CommentUseCase) AssertRoot(rootCommentId int64, rootId int64) bool {
	return uc.repo.AssertRoot(rootCommentId, rootId)
}
func (uc *CommentUseCase) AssertReply(replyCommentId int64, replyId int64) bool {
	return uc.repo.AssertReply(replyCommentId, replyId)
}
func (uc *CommentUseCase) SaveRootComment(pbComment *pb.Comment) {
	uc.repo.SaveRootComment(pbComment)
}
func (uc *CommentUseCase) SaveChildComment(pbComment *pb.Comment) {
	uc.repo.SaveChildComment(pbComment)
}
func (uc *CommentUseCase) QueryRootComment(page, pageSize int) (int32, []pb.Comment, error) {
	return uc.repo.QueryRootComment(page, pageSize)
}
