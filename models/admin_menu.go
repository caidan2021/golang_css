/*
 * @Date: 2022-09-02 19:05:03
 */
package models

import "gin/drivers"

type AdminMenu struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" binding:"required"`
	Icon     string `json:"icon" binding:"required"`
	ParentId int    `json:"parentId" binding:"required"`
	PagePath string `json:"pagePath" binding:"required"`
	ApiPath  string `json:"apiPath" binding:"required"`

	UnixModelTimeWithDel
}

func (AdminMenu) TableName() string {
	return "css_admin_menu"
}

func (AdminMenu) GetMenuList() []*AdminMenu {
	_ls := []*AdminMenu{}
	if err := drivers.Mysql().Model(&AdminMenu{}).Find(&_ls).Error; err != nil {
		return nil
	}
	return _ls
}
