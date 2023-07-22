package foodshare

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	auth "github.com/DrAnonymousNet/foodshare/Auth"
	core "github.com/DrAnonymousNet/foodshare/Core"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type DonationSerializer struct {
	UID             uuid.UUID          `json:"-"`
	Title           string             `json:"title" validate:"required"`
	DonorID         uint8              `json:"donor_id" validate:"required"`
	DonatedObjType  DonatableObjType   `json:"donation_type"`
	DonationDate    time.Time          `json:"donation_date"`
	PickUpAddress   DonationStatusType `json:"pickup_address" validate:"required"`
	ItemDescription string             `json:"item_description"`
	Instance        *Donation          `json:"-"`
	Instances       []Donation         `json:"-"`
	// Also note that we only include a pointer to the original Donation struct in our Request struct.
	// This indirection avoids having to allocate a new copy of Donation.
}

func (d *DonationSerializer) Bind(r *http.Request) error {
	err := d.Validate(r)
	if err != nil {
		return err
	}
	return nil
}

func (d *DonationSerializer) Validate(r *http.Request) error {
	// Validate the tags in the request body
	errorMessage := ""

	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		// Return validation errors to the client
		for _, err := range validationErrors {
			errorMessage += fmt.Sprintf("Field %s: Validation Error (%s) \n", err.Field(), err.Tag())
		}
	}

	// Validate Fields in of the struct
	donationSerializer := reflect.TypeOf(d)

	if donationSerializer.Kind() != reflect.Struct {
		for i := 0; i < donationSerializer.NumField(); i++ {
			field := donationSerializer.Field(i)
			fieldName := field.Name
			value := reflect.ValueOf(donationSerializer).FieldByName(fieldName)
			method := reflect.ValueOf(donationSerializer).MethodByName("Validate" + fieldName)
			if method.IsValid() && method.Type().NumIn() == 1 {
				err := method.Call([]reflect.Value{value})
				if err != nil {
					errorMessage += fmt.Sprintf("Field %s: Validation Error (%s) \n", fieldName, err)
				}
			}
		}

	}
	if errorMessage != "" {
		return errors.New(errorMessage)
	}
	return nil

}

func (d *DonationSerializer) ValidateDonorID(ID int) (int, error) {
	var user auth.User
	core.DB.Model(&auth.User{}).Where("ID = ?", ID)
	if user.ID == 0 {
		return 0, errors.New("donor does not exist")
	}
	return ID, nil
}

func (d *DonationSerializer) ValidateDonation(DonationDate time.Time) (time.Time, error) {
	if d.DonationDate.Before(time.Now()) {
		return time.Time{}, errors.New("donation date cannot be in the past")
	}
	return DonationDate, nil
}

func (d *DonationSerializer) Update() error {
	donation := reflect.TypeOf(d.Instance)
	payload := reflect.TypeOf(d)
	if donation.Kind() != reflect.Struct {
		for i := 0; i < payload.NumField(); i++ {
			field := donation.Field(i)
			fieldName := field.Name
			value := reflect.ValueOf(payload).FieldByName(fieldName)
			if value.IsValid() && value.Type().NumIn() == 1 {
				reflect.ValueOf(donation).FieldByName(fieldName).Set(value)
			}
		}
	}
	err := core.DB.Model(&Donation{}).Save(&donation).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DonationSerializer) Save(r *http.Request) error {
	var user auth.User
	core.DB.Model(&auth.User{}).Where("ID = ?", d.DonorID)

	donationRequest := Donation{
		UID:             uuid.New(),
		Title:           d.Title,
		DonorID:         d.DonorID,
		DonatedObjType:  d.DonatedObjType,
		DonationDate:    d.DonationDate,
		PickUpAddress:   d.PickUpAddress,
		ItemDescription: d.ItemDescription,
		User:            user,
	}
	err := core.DB.Model(&Donation{}).Create(&donationRequest).Error
	if err != nil {
		return err
	}
	d.UID = donationRequest.UID

	return nil
}

func (d *DonationSerializer) Render(w http.ResponseWriter, r *http.Request) error {
	if d.Instance != nil { //The serializer is only meant to be user for retrieval of donation data
		donationObj := reflect.TypeOf(d.Instance)
		donationSerializer := reflect.TypeOf(d)
		if donationObj.Kind() != reflect.Struct {
			for i := 0; i < donationSerializer.NumField(); i++ {
				field := donationObj.Field(i)
				fieldName := field.Name

				//Get the value of the field in the donation object
				value := reflect.ValueOf(donationObj).FieldByName(fieldName)
				if value.IsValid() && value.Type().NumIn() == 1 {
					//Set the value of the field in the donation
					reflect.ValueOf(donationSerializer).FieldByName(fieldName).Set(value)
				}
			}

		}
	}

	return nil
}
