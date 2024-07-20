package group

import "time"

type Groups struct {
	Id         int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 自增ID
	GroupName  string    `gorm:"column:group_name" json:"group_name"`            // 所在群聊名称
	Subscribe  string    `json:"subscribe" gorm:"column:subscribe"`
	Deleted    int       `gorm:"column:deleted;default:0;NOT NULL" json:"deleted"`                         // 是否删除
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 修改时间
}

type Subscribe struct {
	News bool `json:"news"`
	Tips bool `json:"tips"`
}

func (group *Groups) TableName() string {
	return "group"
}
