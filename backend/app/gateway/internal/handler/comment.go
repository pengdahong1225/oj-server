package handler

import (
	"oj-server/app/gateway/internal/define"
	"oj-server/proto/pb"

	"github.com/gin-gonic/gin"
)

func HandleGetRootCommentList(ctx *gin.Context)  {}
func HandleGetChildCommentList(ctx *gin.Context) {}
func HandleCreateComment(ctx *gin.Context) {
	form, ok := validate(ctx, define.AddCommentForm{})
	if !ok {
		return
	}
	resp := &define.Response{
		ErrCode: pb.Error_EN_Success,
	}
}
func HandleDeleteComment(ctx *gin.Context) {}
func HandleLikeComment(ctx *gin.Context)   {}
