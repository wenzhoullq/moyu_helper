package dao

import (
	"testing"
	"weixin_LLM/init/config"
	"weixin_LLM/init/db"
	"weixin_LLM/init/redis"
	"weixin_LLM/lib"
)

func initTest(confAddress string) error {
	err := config.ConfigInit(confAddress)
	if err != nil {
		return err
	}
	err = redis.InitRedis()
	if err != nil {
		return err
	}
	err = db.InitDB()
	if err != nil {
		return err
	}
	return nil
}

func TestRank(t *testing.T) {
	err := initTest("../../config/configTest.toml")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	um := NewWxDao()
	f, err := um.IncrKey("test" + lib.GetCurYearAndMonth())
	if err != nil {
		t.Error(err)
		panic(err)
	}
	t.Log(f)
	f, err = um.IncrKey("test" + lib.GetCurYearAndMonth())
	if err != nil {
		t.Error(err)
		panic(err)
	}
	_, err = um.Expire("test"+lib.GetCurYearAndMonth(), 10)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func TestSign(t *testing.T) {
	err := initTest("../../config/configTest.toml")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	um := NewWxDao()
	f, err := um.AddBit("test"+lib.GetCurYearAndMonth(), lib.GetCurDay())
	if err != nil {
		t.Error(err)
		panic(err)
	}
	t.Log(f)

	//f, err := um.AddSet("test"+lib.GetCurYearAndMonth(), lib.GetCurDay())
	//if err != nil {
	//	t.Error(err)
	//	panic(err)
	//}
}
