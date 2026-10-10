package data

import (
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserInfo struct {
	ID        int64     `gorm:"<-:false;primary_key;autoIncrement;column:id"`
	CreateAt  time.Time `gorm:"<-:false;column:create_at"`
	DeletedAt time.Time `gorm:"<-:false;column:delete_at"`
	Mobile    int64     `gorm:"column:mobile;unique"`
	PassWord  string    `gorm:"column:password"`
	NickName  string    `gorm:"column:nickname"`
	Email     string    `gorm:"column:email"`
	Gender    int32     `gorm:"column:gender"`
	Role      int32     `gorm:"column:role"`
	AvatarUrl string    `gorm:"column:avatar_url"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

func (up *Repo) CreateNewUser(user *UserInfo) (int64, error) {
	var u UserInfo
	result := up.db_.Where("mobile=?", user.Mobile).Find(&u)
	if result.Error != nil {
		logrus.Errorf("查询错误, err: %s", result.Error.Error())
		return -1, status.Errorf(codes.Internal, "query user faild")
	}
	if result.RowsAffected > 0 {
		return -1, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	result = up.db_.Create(&user)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return -1, status.Errorf(codes.Internal, "create user faild")
	}
	return user.ID, nil
}
func (up *Repo) GetUserInfoByUid(mobile int64) (*UserInfo, error) {
	var user UserInfo
	result := up.db_.Where("id=?", mobile).Find(&user)
	if result.Error != nil {
		logrus.Errorf("查询错误, err: %s", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query user faild")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &user, nil
}
func (up *Repo) GetUserInfoByMobile(mobile int64) (*UserInfo, error) {
	var user UserInfo
	result := up.db_.Where("mobile=?", mobile).Find(&user)
	if result.Error != nil {
		logrus.Errorf("查询错误, err: %s", result.Error.Error())
		return nil, status.Errorf(codes.Internal, "query user faild")
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &user, nil
}
func (up *Repo) ResetUserPassword(mobile int64, password string) error {
	/*
		update user_info set password = '123456'
		where mobile = ?;
	*/
	result := up.db_.Model(&UserInfo{}).Where("mobile=?", mobile).Update("password", password)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return status.Errorf(codes.Internal, "update user faild")
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "user not found")
	}
	return nil
}

// 分页查询用户列表
// 查询{id，mobile，nickname，create_at, role}
// @page 页码
// @page_size 单页数量
// @keyword 关键字(昵称)
func (up *Repo) QueryUserList(page, pageSize int, keyword string) (int64, []UserInfo, error) {
	var count int64
	query := up.db_.Model(&UserInfo{})
	if keyword != "" {
		query = query.Where("nickname LIKE ?", "%"+keyword+"%")
	}
	result := query.Count(&count)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return -1, nil, status.Errorf(codes.Internal, "query failed")
	}
	offSet := (page - 1) * pageSize
	var users []UserInfo
	result = query.Select("id, mobile, nickname, email, create_at, role").Offset(offSet).Limit(pageSize).Find(&users)
	if result.Error != nil {
		logrus.Errorf("query user list failed: %s", result.Error.Error())
		return -1, nil, status.Errorf(codes.Internal, "query failed")
	}
	return count, users, nil
}
