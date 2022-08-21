package models

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/plugin/soft_delete"
	"time"
)

type Tabler interface {
	TableName() string
}

type UnixModelTimeWithDel struct {
	CreatedAt int                   `json:"createdAt"`
	UpdatedAt *DefaultTimeFormat    `json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
}

type UnixModelTime struct {
	CreatedAt int                `json:"createdAt"`
	UpdatedAt *DefaultTimeFormat `json:"updatedAt"`
}

type BaseModelDefTime struct {
	BaseModelDefTimeWithoutDel
	DeletedAt *DefaultTimeFormat `json:"-"`
}
type BaseModelDefTimeWithoutDel struct {
	CreatedAt DefaultTimeFormat  `gorm:"not null" json:"createdAt"`
	UpdatedAt *DefaultTimeFormat `json:"updatedAt"`
}

type DefaultTimeFormat struct {
	time.Time
}

func (t DefaultTimeFormat) MarshalJSON() ([]byte, error) {
	tune := t.Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
}

func (t *DefaultTimeFormat) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DefaultTimeFormat{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t DefaultTimeFormat) Value() (driver.Value, error) {
	var _time time.Time
	if t.Time.UnixNano() == _time.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}
