/*
 * @Date: 2022-08-19 10:08:10
 */
package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"gin/drivers"

	"gorm.io/gorm"
)

const PasswordPrefix = "caiseshi_"

type AdminUser struct {
	ID            int64  `json:"id"`
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email"`
	Password      string `json:"password" binding:"required"`
	RememberToken string `json:"rememberToken"`
	UnixModelTime
}

func (AdminUser) TableName() string {
	return "css_admin_user"
}

func GetAdminUserByNameAndPw(name, pw string) (*AdminUser, error) {
	_one := &AdminUser{}
	if err := drivers.Mysql().Model(&AdminUser{}).Where("name = ?", name).Where("password = ?", _one.EncryptionPw(pw)).First(&_one).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return _one, nil
}

func (AdminUser) NameExists(name string, id int64) bool {
	query := drivers.Mysql().Model(&AdminUser{}).Where("name = ?", name)
	if id != 0 {
		query = query.Where("id != ?", id)
	}
	_l := []*AdminUser{}
	query.Find(&_l)
	return len(_l) > 0
}

func (AdminUser) EncryptionPw(password string) string {
	password = PasswordPrefix + password
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
