package source

import "time"

type Source struct {
	Id         int32     `gorm:"column:id;primary_key;AUTO_INCREMENT"json:"id"`                            // 自增ID
	Status     int       `gorm:"column:status;default:0;NOT NULL" json:"status"`                           // 状态
	Deleted    int       `gorm:"column:deleted;default:0;NOT NULL" json:"deleted"`                         // 是否删除
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 修改时间
	SourceName string    `gorm:"column:source_name" json:"source_name"`                                    // 资源名称
	SourceDesc string    `gorm:"column:source_desc" json:"source_desc"`                                    // 资源描述
	SourceLink string    `gorm:"column:source_link" json:"source_link"`                                    // 资源链接
	SourceType int       `gorm:"column:source_type;default:0" json:"source_type"`                          // 资源类型
	SourceExp  string    `gorm:"column:source_exp" json:"source_exp"`                                      // 资源过期时间
}

func (source *Source) TableName() string {
	return "source"
}
