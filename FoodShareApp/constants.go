package foodshare

type DonatableObjType string
type DonationStatusType string
type RequestStatusType string
type RequestFromType string

const (
	FoodStuff            DonatableObjType = "FoodStuff"
	Cloths               DonatableObjType = "Cloths"
	MedicalSupplies      DonatableObjType = "MedicalSupplies"
	SchoolSupplies       DonatableObjType = "SchoolSupplies"
	PersonalCareSupplies DonatableObjType = "PersonalCareSupplies"
	BooksAndToys         DonatableObjType = "BooksAndToys"
)

const (
	Pending  DonationStatusType = "Pending"
	PickedUp DonationStatusType = "PickedUp"
)

const (
	PartiallyFulfilled RequestStatusType = "PartiallyFulfilled"
	FullyFulfilled     RequestStatusType = "FullyFulfilled"
)

const (
	WareHouse RequestFromType = "WareHouse"
	Community RequestFromType = "Community"
)
