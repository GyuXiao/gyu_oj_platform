// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entity

import (
	"time"
)

const TableNameQuestionSubmit = "question_submit"

// QuestionSubmit 题目提交表
type QuestionSubmit struct {
	ID         int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:id" json:"id,string"`                   // id
	Language   string    `gorm:"column:language;type:varchar(128);not null;comment:编程语言" json:"language"`                           // 编程语言
	SubmitCode string    `gorm:"column:submitCode;type:text;not null;comment:用户提交的代码" json:"submitCode"`                            // 用户提交的代码
	JudgeInfo  string    `gorm:"column:judgeInfo;type:text;comment:判题信息（json 对象）" json:"judgeInfo"`                                 // 判题信息（json 对象）
	Status     int64     `gorm:"column:status;type:int;not null;comment:判题状态（0 - 待判题、1 - 判题中、2 - 成功、3 - 失败）" json:"status"`         // 判题状态（0 - 待判题、1 - 判题中、2 - 成功、3 - 失败）
	QuestionID int64     `gorm:"column:questionId;type:bigint;not null;comment:题目 id" json:"questionId"`                            // 题目 id
	UserID     int64     `gorm:"column:userId;type:bigint;not null;comment:创建用户 id" json:"userId"`                                  // 创建用户 id
	CreateTime time.Time `gorm:"column:createTime;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createTime"` // 创建时间
	UpdateTime time.Time `gorm:"column:updateTime;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updateTime"` // 更新时间
	IsDelete   int64     `gorm:"column:isDelete;type:tinyint;not null;comment:是否删除" json:"isDelete"`                                // 是否删除
}

// TableName QuestionSubmit's table name
func (*QuestionSubmit) TableName() string {
	return TableNameQuestionSubmit
}
