package service

import (
	"context"
	"go-kit-demo/demo03/v1/dao"
	"go-kit-demo/demo03/v1/redis"
	"go-kit-demo/demo03/v1/utils"
	"testing"
)

func TestUserServiceImpl_Login(t *testing.T) {
	if err := dao.InitMysql("192.168.142.128", "3306", "root", "mysqlly", "user"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := redis.InitRedis("192.168.142.128", "6379", ""); err != nil {
		t.Error(err)
		t.FailNow()
	}
	utils.NewLoggerServer()
	logger := utils.GetLogger()
	userDao := dao.NewUserDaoImpl()
	userService := NewUserServiceImpl(userDao, logger)
	user, err := userService.Login(context.Background(), &LoginVO{
		Email:    "Ricoly@qq.com",
		Password: "ly",
	})

	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.Id)
}

func TestUserServiceImpl_Register(t *testing.T) {
	if err := dao.InitMysql("192.168.142.128", "3306", "root", "mysqlly", "user"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := redis.InitRedis("192.168.142.128", "6379", ""); err != nil {
		t.Error(err)
		t.FailNow()
	}
	utils.NewLoggerServer()
	logger := utils.GetLogger()
	userDao := dao.NewUserDaoImpl()
	userService := NewUserServiceImpl(userDao, logger)
	user, err := userService.Register(context.Background(),&RegisterUserVO{
		Username: "bbaava",
		Password: "123",
		Email:    "abcfdac@qq.com",
	})
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.Id)
}