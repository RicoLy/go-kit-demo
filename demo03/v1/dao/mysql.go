package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var db *gorm.DB

func InitMysql(host, port, user, password, dbName string) (err error) {
	if db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbName)); err != nil {
		log.Println(err)
		return
	}
	return 
}

// 页码结构体
type Pagination struct {
	Page      int64 `json:"page" example:"0"`      // 当前页
	PageSize  int64 `json:"pageSize" example:"20"` // 每页条数
	Total     int64 `json:"total"`                 // 总条数
	TotalPage int64 `json:"totalPage"`             // 总页数
}