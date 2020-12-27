package dao

import "time"

type UserEntity struct {
	ID        int64
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}

func (UserEntity) TableName() string {
	return "user"
}

type UserDAO interface {
	SelectByEmail(email string) (user *UserEntity, err error)
	Save(user *UserEntity) (err error)
}

type UserDAOImpl struct {
}

func (u *UserDAOImpl) SelectByEmail(email string) (user *UserEntity, err error) {
	user = &UserEntity{}
	err = db.Where("email = ?", email).First(user).Error
	return
}

func (u *UserDAOImpl) Save(user *UserEntity) (err error) {
	return db.Create(user).Error
}
