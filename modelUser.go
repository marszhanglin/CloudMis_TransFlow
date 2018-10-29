// echoRouteSaveflow
package main

import (
	"time"
)

func dbsaveUser(user *User) error {
	//tx := db.Begin()
	user.UserId = getTimeUUID()
	if err := db.Table("t_user").Create(&user).Error; err != nil {
		glogError(err.Error())
		//tx.Rollback()
		return err
	}
	//tx.Commit()
	return nil
}

func dbQueryUserByName(userName string) (isExit bool, userrtn User) {

	queryUsers := []User{}
	db.Table("t_user").Select(" * ").Where(" user_name = ? ", userName).Find(&queryUsers)
	if len(queryUsers) > 1 {
		glogError("存在多个同名用户")
		return false, queryUsers[0]
	} else if len(queryUsers) == 1 {
		return true, queryUsers[0]
	} else {
		user := &User{}
		return false, *user
	}
}

func dbQueryUserByShopId(shopId int64, userType int64) ([]QueryUser, error) {

	queryUsers := []QueryUser{}
	err := db.Table("t_user").Select("user_id,user_name,user_type,shop_id ").Where(" shop_id = ? and user_type = ? ", shopId, userType).Find(&queryUsers).Error
	return queryUsers, err
}

func dbDeleteUserByName(userName string) error {
	tx := db.Begin()
	if err := tx.Table("t_user").Where(" user_name = ? ", userName).Delete(User{}).Error; err != nil {
		glogError(err.Error())
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

type User struct {
	UserId     int64     `gorm:"column:user_id;primary_key" validate:"-"`
	UserName   string    `gorm:"column:user_name" validate:"required,max=64"`
	Password   string    `gorm:"column:password" validate:"required,max=64"`
	UserStatus int64     `gorm:"column:user_status" validate:"-"`
	UserType   int64     `gorm:"column:user_type" validate:"required"`
	OperatorId int64     `gorm:"column:operator_id" validate:"max=64"`
	ShopId     int64     `gorm:"column:shop_id" validate:"-"`
	CreTime    time.Time `gorm:"column:cre_time" validate:"-"`
	UpdTime    time.Time `gorm:"column:upd_time" validate:"-"`
	Remark     string    `gorm:"column:remark" validate:"max=256"`
}

type QueryUser struct {
	UserId     int64  `gorm:"column:user_id;primary_key" validate:"-"`
	UserName   string `gorm:"column:user_name" validate:"required,max=64"`
	UserStatus int64  `gorm:"column:user_status" validate:"-"`
	UserType   int64  `gorm:"column:user_type" validate:"required"`
	OperatorId int64  `gorm:"column:operator_id" validate:"max=64"`
	ShopId     int64  `gorm:"column:shop_id" validate:"-"`
}
