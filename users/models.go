package users

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name string
	TgId int
}
