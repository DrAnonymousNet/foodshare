package notifications

import (
	"time"

	foodshare "github.com/DrAnonymousNet/foodshare/FoodShareApp"
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	RequestID           uint8
	DonationRequest     foodshare.DonationRequest `gorm:"foreignkey:RequestID"`
	NotificationMessage string
	TimeStamp           time.Time
}
