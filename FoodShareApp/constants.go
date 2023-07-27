package foodshare

import (
	"database/sql/driver"

)

type DonatableObjType string

const (
	FoodStuff            DonatableObjType = "FoodStuff"
	Cloths               DonatableObjType = "Cloths"
	MedicalSupplies      DonatableObjType = "MedicalSupplies"
	SchoolSupplies       DonatableObjType = "SchoolSupplies"
	PersonalCareSupplies DonatableObjType = "PersonalCareSupplies"
	BooksAndToys         DonatableObjType = "BooksAndToys"
)

func (e *DonatableObjType) Scan(value interface{}) error {
	*e = DonatableObjType(value.([]byte))
	return nil
}

func (e DonatableObjType) Value() (driver.Value, error) {
	return string(e), nil
}

type DonationStatusType string

const (
	Pending  DonationStatusType = "Pending"
	PickedUp DonationStatusType = "PickedUp"
)

func (e *DonationStatusType) Scan(value interface{}) error {
	*e = DonationStatusType(value.([]byte))
	return nil
}

func (e DonationStatusType) Value() (driver.Value, error) {
	return string(e), nil
}


type RequestFromType string

const (
	WareHouse RequestFromType = "WareHouse"
	Community RequestFromType = "Community"
)

func (e *RequestFromType) Scan(value interface{}) error {
	*e = RequestFromType(value.([]byte))
	return nil
}

func (e RequestFromType) Value() (driver.Value, error) {
	return string(e), nil
}

type RequestStatusType string


const (
	PartiallyFulfilled RequestStatusType = "PartiallyFulfilled"
	FullyFulfilled     RequestStatusType = "FullyFulfilled"
)


func (e *RequestStatusType) Scan(value interface{}) error {
	*e = RequestStatusType(value.([]byte))
	return nil
}

func (e RequestStatusType) Value() (driver.Value, error) {
	return string(e), nil
}
