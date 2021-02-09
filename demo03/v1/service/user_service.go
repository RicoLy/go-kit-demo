package service

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"go-kit-demo/demo03/v1/dao"
	"go-kit-demo/demo03/v1/redis"
	"go.uber.org/zap"
	"log"
	"time"
)

type LoginVO struct {
	Email string
	Password string
}

type UserInfoDTO struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
}

type RegisterUserVO struct {
	Username string
	Password string
	Email string
}

var (
	ErrUserExisted = errors.New("user is existed")
	ErrPassword    = errors.New("email and password are not match")
	ErrRegistering = errors.New("email is registering")
)

type UserService interface {
	Login(ctx context.Context, vo *LoginVO) (userInfoDTO *UserInfoDTO, err error)
	Register(ctx context.Context, vo *RegisterUserVO) (user *UserInfoDTO, err error)
}

type UserServiceImpl struct {
	userDAO dao.UserDao
	logger *zap.Logger
}

func NewUserServiceImpl(userDao dao.UserDao, log *zap.Logger) UserService {
	var service UserService
	service = &UserServiceImpl{
		userDAO: userDao,
		logger:  log,
	}
	service = NewLogMiddlewareServer(log)(service)

	return service
}

func (u UserServiceImpl) Login(ctx context.Context, vo *LoginVO) (userInfoDTO *UserInfoDTO, err error) {
	user, err := u.userDAO.First("email = ?", vo.Email)
	if  err == nil {
		if user.Password == vo.Password {
			return &UserInfoDTO{
				Id:       user.Id,
				Username: user.Username,
				Email:    user.Email,
			}, nil
		} else {
			return nil, ErrPassword
		}
	} else {
		log.Printf("err: %s", err)
	}
	return nil, err
}

func (u UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVO) (userInfoDTO *UserInfoDTO, err error) {
	lock := redis.GetRedisLock(vo.Email, time.Duration(5) * time.Second)
	err = lock.Lock()
	if err != nil {
		log.Printf("err : %s", err)
		return nil, ErrRegistering
	}
	defer lock.Unlock()
	user, err := u.userDAO.First("email = ?", vo.Email)
	if (err == nil && &user == nil) || err == gorm.ErrRecordNotFound {
		newUser := &dao.UserEntity{
			Username:  vo.Username,
			Password:  vo.Password,
			Email:     vo.Email,
		}
		if err = u.userDAO.Add(newUser); err == nil {
			return &UserInfoDTO{
				Id:       newUser.Id,
				Username: newUser.Username,
				Email:    newUser.Email,
			}, nil
		}
	}

	if err != nil {
		err = ErrUserExisted
	}
	return nil, err
}

























