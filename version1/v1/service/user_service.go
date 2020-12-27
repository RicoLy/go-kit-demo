package service

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"go-kit-demo/version1/v1/dao"
	"go-kit-demo/version1/v1/redis"
	"go.uber.org/zap"
	"log"
	"time"
)

type LoginVO struct {
	Email string
	Password string
}

type UserInfoDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUserVO struct {
	Username string
	Password string
	Email    string
}

var (
	ErrUserExisted = errors.New("user is existed")
	ErrPassword    = errors.New("email and password are not match")
	ErrRegistering = errors.New("email is registering")
)

type UserService interface {
	Login(ctx context.Context, vo *LoginVO) (user *UserInfoDTO, err error)
	Register(ctx context.Context, vo *RegisterUserVO) (user *UserInfoDTO, err error)
}

type UserServiceImpl struct {
	userDAO dao.UserDAO
	logger *zap.Logger
}

func NewUserServiceImpl(userDAO dao.UserDAO, log *zap.Logger) UserService {
	var server UserService
	server = &UserServiceImpl{
		userDAO: userDAO,
		logger:  log,
	}
	server = NewLogMiddlewareServer(log)(server)
	return server
}

func (userService UserServiceImpl) Login(ctx context.Context, vo *LoginVO) (*UserInfoDTO, error) {
	user, err := userService.userDAO.SelectByEmail(vo.Email)
	if  err == nil {
		if user.Password == vo.Password {
			return &UserInfoDTO{
				ID:       user.ID,
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

func (userService UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO,  error) {
	lock := redis.GetRedisLock(vo.Email, time.Duration(5) * time.Second)
	err := lock.Lock()
	if err != nil {
		log.Printf("err : %s", err)
		return nil, ErrRegistering
	}
	defer lock.Unlock()
	existUser, err := userService.userDAO.SelectByEmail(vo.Email)

	if (err == nil && existUser == nil) || err == gorm.ErrRecordNotFound {
		newUser := &dao.UserEntity{
			Username:  vo.Username,
			Password:  vo.Password,
			Email:     vo.Email,
		}

		if err = userService.userDAO.Save(newUser); err == nil {
			return &UserInfoDTO{
				ID:       newUser.ID,
				Username: newUser.Username,
				Email:    newUser.Email,
			}, nil
		}
	}
	if err == nil {
		err = ErrUserExisted
	}
	return nil, err
}

























