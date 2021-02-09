package dao

import (
	"github.com/jinzhu/gorm"
	"math"
	"time"
)

type UserEntity struct {
	Id        int64
	Username  string
	Password  string
	Email     string
	Money     int64
	CreatedAt time.Time
}

// 表名
func (u *UserEntity) TableName() string {
	return "user"
}

type UserDao interface {
	Add(model *UserEntity) (err error)
	Save(model *UserEntity) (err error)
	Delete(query interface{}, args ...interface{}) (err error)
	First(query interface{}, args ...interface{}) (model UserEntity, err error)
	Find(query interface{}, page *Pagination, args ...interface{}) (models []UserEntity, err error)
	Count(where interface{}, args ...interface{}) (count int64, err error)
}

type UserDaoImpl struct {
	DB *gorm.DB
}

func NewUserDaoImpl() *UserDaoImpl {
	return &UserDaoImpl{
		DB: db,
	}
}

//--------------------------增删改相关业务------------------------

// 添加记录
func (u *UserDaoImpl) Add(model *UserEntity) (err error) {
	return u.DB.Create(model).Error
}

// 更新保存记录
func (u *UserDaoImpl) Save(model *UserEntity) (err error) {
	return u.DB.Save(model).Error
}

// 软删除：结构体需要继承Base model 有delete_at字段
func (u *UserDaoImpl) Delete(query interface{}, args ...interface{}) (err error) {
	return u.DB.Where(query, args...).Delete(&UserEntity{}).Error
}

//--------------------------查询相关业务------------------------

// 根据条件获取单挑记录
func (u *UserDaoImpl) First(query interface{}, args ...interface{}) (model UserEntity, err error) {
	err = u.DB.Where(query, args...).First(&model).Error
	return
}

// 获取列表 数据量大时Count数据需另外请求接口
func (u *UserDaoImpl) Find(query interface{}, page *Pagination, args ...interface{}) (models []UserEntity, err error) {
	switch page {
	case nil:
		err = u.DB.Find(&models).Error
	default:
		err = u.DB.Model(UserEntity{}).Where(query, args...).
			Count(&page.Total).Offset((page.Page - 1) * page.PageSize).
			Limit(page.PageSize).Find(&models).Error

		page.TotalPage = int64(math.Ceil(float64(page.Total / page.PageSize)))
	}
	return
}

// 获取总记录条数
func (u *UserDaoImpl) Count(where interface{}, args ...interface{}) (count int64, err error) {
	err = u.DB.Model(&UserEntity{}).Where(where, args...).Count(&count).Error
	return
}
