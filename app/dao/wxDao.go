package dao

import (
	"context"
	r "github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"time"
	"weixin_LLM/dto/user"
	"weixin_LLM/init/db"
	"weixin_LLM/init/redis"
)

type WxDao struct {
	*gorm.DB
	redisClient *r.Client
	context.Context
}

func NewWxDao(ops ...func(c *WxDao)) *WxDao {
	wxDao := &WxDao{
		DB:          db.DB,
		redisClient: redis.RedisClient,
		Context:     redis.Ctx,
	}
	for _, op := range ops {
		op(wxDao)
	}
	return wxDao
}

func (wd *WxDao) UpdateUserID(user *user.User) error {
	if err := wd.Table(user.TableName()).Where("group_name = ? and user_name = ?", user.GroupName, user.UserName).Update("user_id", user.UserId).Error; err != nil {
		return err
	}
	return nil
}

func (wd *WxDao) GetUsersByGroupName(groupName string) ([]*user.User, error) {
	u := &user.User{}
	users := make([]*user.User, 0)
	if err := wd.Table(u.TableName()).Where("group_name = ?", groupName).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (wd *WxDao) GetUserByUserNameAndGroupNameAndUserId(displayName string, groupName string, userId string) (*user.User, error) {
	user := &user.User{}
	if err := wd.Table(user.TableName()).Where("user_name = ? and group_name = ? and user_id = ?", displayName, groupName, userId).Find(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (wd *WxDao) UpdateUserReward(user *user.User) error {
	if err := wd.Table(user.TableName()).Where("user_name = ? and group_name = ? and user_id = ?", user.UserName, user.GroupName, user.UserId).Update("reward", user.Reward).Error; err != nil {
		return err
	}
	return nil
}

func (wd *WxDao) UpdateUserExtra(user *user.User) error {
	if err := wd.Table(user.TableName()).Where("user_name = ? and group_name = ? and user_id = ?", user.UserName, user.GroupName, user.UserId).Update("extra", user.Extra).Error; err != nil {
		return err
	}
	return nil
}

func (wd *WxDao) UpdateUserName(user *user.User) error {
	if err := wd.Table(user.TableName()).Where("user_name = ? and group_name = ? and user_id = ?", user.UserName, user.GroupName, user.UserId).Update("user_name", user.UserName).Error; err != nil {
		return err
	}
	return nil
}

func (wd *WxDao) AddUser(user *user.User) error {
	if err := wd.Table(user.TableName()).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (wd *WxDao) Expire(key string, exp int64) (bool, error) {
	return wd.redisClient.Expire(wd.Context, key, time.Duration(exp*1000000000)).Result()
}

func (wd *WxDao) AddBit(key string, value int64) (int64, error) {
	return wd.redisClient.SetBit(wd.Context, key, value, 1).Result()
}

func (wd *WxDao) GetBit(key string, value int64) (int64, error) {
	return wd.redisClient.GetBit(wd.Context, key, value).Result()
}

func (wd *WxDao) IncrKey(key string) (int64, error) {
	incrResult, err := wd.redisClient.Incr(wd.Context, key).Result()
	if err != nil {
		return 0, err
	}
	return incrResult, nil
}

func (wd *WxDao) GetString(key string) (string, error) {
	str, err := wd.redisClient.Get(wd.Context, key).Result()
	if err != nil {
		return "", err
	}
	return str, nil
}
func (wd *WxDao) SetString(key string, value interface{}, exp int64) error {
	err := wd.redisClient.Set(wd.Context, key, value, time.Duration(exp*1000000000)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (wd *WxDao) DelString(key string) (int64, error) {
	return wd.redisClient.Del(wd.Context, key).Result()
}
