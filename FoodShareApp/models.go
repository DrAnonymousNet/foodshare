package foodshare

import (
	"time"

	auth "github.com/DrAnonymousNet/foodshare/Auth"
	core "github.com/DrAnonymousNet/foodshare/Core"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Donation struct {
	gorm.Model
	core.ModelStruct
	UID             uuid.UUID
	Title           string
	DonorUID        uuid.UUID
	User            *auth.User `gorm:"foreignKey:DonorUID;references:UID"`
	DonatedObjType  string     //DonatableObjType `sql:"type:enum_donatable_obj_type"`
	DonationDate    time.Time
	PickUpAddress   string //DonationStatusType `sql:"type:enum_donation_status_type"`
	ItemDescription string
}

type DonationRequest struct {
	gorm.Model
	core.ModelStruct
	UID                uuid.UUID `gorm:"default:generate_uuid_v4"`
	RequestorID        uint8
	User               auth.User `gorm:"foreignKey:RequestorID"`
	RequestDescription string
	Quantity           uint8
	RequestDate        time.Time
	DeliveryAddress    string
	RequestStatus      string //RequestStatusType `sql:"type:request_status_type"`
	RequestFrom        string //RequestFromType   `sql:"type:request_from_type"`
}
