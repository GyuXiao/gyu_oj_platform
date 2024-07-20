// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entity

import (
	"time"
)

const TableNameUser = "user"

// User 用户表
type User struct {
	ID         int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:唯一 id" json:"id,string"`                                                                        // 唯一 id
	Username   string    `gorm:"column:username;type:varchar(256);not null;comment:用户昵称" json:"username"`                                                                                   // 用户昵称
	Password   string    `gorm:"column:password;type:varchar(512);not null;comment:用户密码" json:"password"`                                                                                   // 用户密码
	AvatarURL  string    `gorm:"column:avatarUrl;type:varchar(1024);not null;default:https://gyu-pic-bucket.oss-cn-shenzhen.aliyuncs.com/gyustudio_icon.jpg;comment:用户头像" json:"avatarUrl"` // 用户头像
	Email      string    `gorm:"column:email;type:varchar(256);comment:用户邮箱" json:"email"`                                                                                                  // 用户邮箱
	Phone      string    `gorm:"column:phone;type:varchar(256);comment:手机号" json:"phone"`                                                                                                   // 手机号
	UserRole   int64     `gorm:"column:userRole;type:tinyint;not null;comment:用户角色 0 - 普通用户 1 - 管理员" json:"userRole"`                                                                       // 用户角色 0 - 普通用户 1 - 管理员
	IsDelete   int64     `gorm:"column:isDelete;type:tinyint;not null;comment:是否删除 0 - 未删除 1- 删除" json:"isDelete"`                                                                          // 是否删除 0 - 未删除 1- 删除
	CreateTime time.Time `gorm:"column:createTime;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createTime"`                                                         // 创建时间
	UpdateTime time.Time `gorm:"column:updateTime;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updateTime"`                                                         // 更新时间
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
