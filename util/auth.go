/*
 * @Date: 2022-08-19 09:58:54
 */
package util

const AdminUserKey = "caiseshi_admin_user"

type BaseAuth struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
