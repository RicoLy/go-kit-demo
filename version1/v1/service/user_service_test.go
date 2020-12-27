package service

import (
	"context"
	"go-kit-demo/version1/v1/dao"
	"go-kit-demo/version1/v1/redis"
	"testing"
)

func TestUserServiceImpl_Login(t *testing.T) {


	err := dao.InitMysql("192.168.190.139", "3306", "root", "mysqlly", "user")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("192.168.190.139","6379", "")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}


	userService := &UserServiceImpl{
		userDAO: &dao.UserDAOImpl{},
	}

	user, err := userService.Login(context.Background(), "aoho@mail.com", "aoho")

	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)

}

func TestUserServiceImpl_Register(t *testing.T) {


	err := dao.InitMysql("192.168.190.139", "3306", "root", "mysqlly", "user")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("192.168.190.139","6379", "" )
	if err != nil{
		t.Error(err)
		t.FailNow()
	}


	userService := &UserServiceImpl{
		userDAO: &dao.UserDAOImpl{},
	}

	user, err := userService.Register(context.Background(),
		&RegisterUserVO{
			Username:"aoho",
			Password:"aoho",
			Email:"aoho@mail.com",
		})

	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)

}