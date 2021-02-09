package dao

import (
	"testing"
)

func TestUserDaoImpl_Add(t *testing.T) {
	if err := InitMysql("192.168.142.128", "3306", "root", "mysqlly", "user"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	userDao := NewUserDaoImpl()
	user := &UserEntity{
		Username: "Phatumai",
		Password: "ptm",
		Email:    "patumai@qq.com",
	}
	if err := userDao.Add(user); err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("new user id is %d", user.Id)
}
