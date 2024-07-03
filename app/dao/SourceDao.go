package dao

import (
	"github.com/jinzhu/gorm"
	"weixin_LLM/dto/source"
	"weixin_LLM/init/db"
	"weixin_LLM/lib/constant"
)

type SourceDao struct {
	*gorm.DB
}

func NewSourceDao(ops ...func(c *SourceDao)) *SourceDao {
	sourceDao := &SourceDao{
		DB: db.DB,
	}
	for _, op := range ops {
		op(sourceDao)
	}
	return sourceDao
}

func (sd *SourceDao) GetNotExpSources() ([]*source.Source, error) {
	sourceList := make([]*source.Source, 0)
	source := &source.Source{}
	if err := sd.Table(source.TableName()).Where("status = ?", constant.SourceNorMal).Find(&sourceList).Error; err != nil {
		return nil, err
	}
	return sourceList, nil
}
func (sd *SourceDao) CreateSource(source *source.Source) error {
	if err := sd.Table(source.TableName()).Create(source).Error; err != nil {
		return err
	}
	return nil
}

func (sd *SourceDao) UpdateSource(source *source.Source) error {
	if err := sd.Table(source.TableName()).Where("id = ? ", source.Id).Update(source).Error; err != nil {
		return err
	}
	return nil
}
