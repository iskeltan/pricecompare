package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// CompareURL is a model
type CompareURL struct {
	gorm.Model
	URL      string `gorm:"type:varchar(255)"`
	User     User   `gorm:"foreignkey:UserID; not null"`
	UserID   uint
	Schedule time.Time
	// Schedule TODO: schedule was here
}

func (CompareUrl CompareURL) CreateNew() CompareURL {
	conn, _ := Database()
	CompareUrl.Schedule := 2 * time.Hour
	return CompareUrl
}
