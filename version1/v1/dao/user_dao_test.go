package dao

import (
	"testing"
)

func TestUserDAOImpl_Save(t *testing.T) {
	userDAO := &UserDAOImpl{}

	if err := InitMysql("192.168.190.139", "3306", "root", "mysqlly", "user"); err != nil {
		t.Error(err)
		t.FailNow()
	}

	user := &UserEntity{
		Username: "Candy",
		Password: "cd",
		Email:    "Ricoly@qq.com",
	}

	if err := userDAO.Save(user); err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("new User ID is %d", user.ID)
}

func TestUserDAOImpl_SelectByEmail(t *testing.T) {

	userDAO := &UserDAOImpl{}

	if err := InitMysql("192.168.190.139", "3306", "root", "mysqlly", "user"); err != nil {
		t.Error(err)
		t.FailNow()
	}

	user, err := userDAO.SelectByEmail("Ricoly@qq.com")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("result uesrname is %s", user.Username)

}
