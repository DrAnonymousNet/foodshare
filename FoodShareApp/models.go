package foodshare

import (
	"time"

	auth "github.com/DrAnonymousNet/foodshare/Auth"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Donation struct {
	gorm.Model
	UID             uuid.UUID
	Title           string
	DonorID         uint8
	User            *auth.User        `gorm:"foreignKey:DonorID;references:ID"`
	DonatedObjType  DonatableObjType `sql:"type:donatable_obj_type"`
	DonationDate    time.Time
	PickUpAddress   DonationStatusType `sql:"type:donation_status_type"`
	ItemDescription string
}

type DonationRequest struct {
	gorm.Model
	UID                uuid.UUID `gorm:"default:generate_uuid_v4"`
	RequestorID        uint8
	User               auth.User `gorm:"foreignKey:RequestorID"`
	RequestDescription string
	Quantity           uint8
	RequestDate        time.Time
	DeliveryAddress    string
	RequestStatus      RequestStatusType `sql:"type:request_status_type"`
	RequestFrom        RequestFromType   `sql:"type:request_from_type"`
}
