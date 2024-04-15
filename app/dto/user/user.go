package user

import "time"

type User struct {
	Id         int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`                           // 自增ID
	Status     int       `gorm:"column:status;default:0;NOT NULL" json:"status"`                           // 状态
	Deleted    int       `gorm:"column:deleted;default:0;NOT NULL" json:"deleted"`                         // 是否删除
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 修改时间
	UserName   string    `gorm:"column:user_name" json:"user_name"`                                        // 用户昵称
	GroupName  string    `gorm:"column:group_name" json:"group_name"`                                      // 所在群聊名称
	UserId     string    `gorm:"column:user_id" json:"user_id"`                                            // 临时ID
	Reward     int       `gorm:"column:reward;default:0" json:"reward"`                                    // 金币
}

func (user *User) TableName() string {
	return "user"
}
