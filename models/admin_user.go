/*
 * @Date: 2022-08-19 10:08:10
 */
package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gin/drivers"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

const PasswordPrefix = "caiseshi_"
const secretKey = "eyJhbGciOiJSUzI1NiIsInR"

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

func GetAdminUserByToken(token string) (*AdminUser, error) {
	_one := &AdminUser{}
	if err := drivers.Mysql().Model(&AdminUser{}).Where("remember_token = ?", token).First(&_one).Error; err != nil {
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

func (a *AdminUser) UpdateRememberMe(rememberMe bool) {
	if rememberMe {
		newToken, err := GenerateToken(a.Name)
		if err != nil {
			return
		}
		a.RememberToken = newToken
	} else {
		a.RememberToken = ""
	}
	drivers.Mysql().Save(&a)
}

type Claims struct {
	Username string
	jwt.StandardClaims
}

// GenerateToken 生成Token值
func GenerateToken(userName string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(86400 * 14 * time.Second)
	claims := Claims{
		Username: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	return token, err
}

// token: "eyJhbGciO...解析token"
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	fmt.Println("-=====")

	return nil, err
}
