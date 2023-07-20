package foodshare

import (
	"time"

	accounts "github.com/DrAnonymousNet/foodshare/auth"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Donation struct {
	gorm.Model
	UID             uuid.UUID
	Title           string
	DonorID         uint8
	User            accounts.User    `gorm:"foreignKey:DonorID"`
	DonatedObjType  DonatableObjType `gorm:"type:ENUM('FoodStuff', 'Cloths', 'MedicalSupplies', 'SchoolSupplies', 'PersonalCareSupplies', 'BooksAndToys')"`
	DonationDate    time.Time
	PickUpAddress   DonationStatusType `gorm:"type:ENUM('Pending', 'PickedUp')"`
	ItemDescription string
}

type DonationRequest struct {
	gorm.Model
	UID                uuid.UUID
	RequestorID        uint8
	User               accounts.User `gorm:"foreignKey:RequestorID"`
	RequestDescription string
	Quantity           uint8
	RequestDate        time.Time
	DeliveryAddress    string
	RequestStatus      RequestStatusType `gorm:"type:ENUM('PartiallyFulfilled, 'FullyFulfilled')"`
	RequestFrom        RequestFromType   `gorm:"type:ENUM('WareHouse', 'Community')"`
}
