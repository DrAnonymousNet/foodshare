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


type DonationHTTPRequestBody struct{
	UID 		   	uuid.UUID `json:"-"` 
	Title           string `json:"title" validate:"required"`
	DonorID         uint8	`json:"donor_id" validate:"required"`
	DonatedObjType  DonatableObjType `json:"donation_type"`
	DonationDate    time.Time	`json:"donation_date"`
	PickUpAddress   DonationStatusType `json:"pickup_address" validate:"required"`
	ItemDescription string `json:"item_description"`
}

func (d *DonationHTTPRequestBody)Bind(r *http.Request) error {
	err := d.Validate(r)
	if err != nil {
		return err
	}
	return nil
}

func (d *DonationHTTPRequestBody)Validate(r *http.Request) error {
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
	donationRequestType := reflect.TypeOf(d)

	if donationRequestType.Kind() != reflect.Struct {
		for i := 0; i < donationRequestType.NumField(); i++ {
			field := donationRequestType.Field(i)
			fieldName := field.Name
			value := reflect.ValueOf(donationRequestType).FieldByName(fieldName)
			method := reflect.ValueOf(donationRequestType).MethodByName("Validate" + fieldName)
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

func (d *DonationHTTPRequestBody)ValidateDonorID(ID int) (int, error){
	var user auth.User
	core.DB.Model(&auth.User{}).Where("ID = ?", ID)
	if user.ID == 0 {
		return 0, errors.New("donor does not exist")
	}
	return ID, nil
}

func (d *DonationHTTPRequestBody)ValidateDonation(DonationDate time.Time)(time.Time, error){
	if d.DonationDate.Before(time.Now()){
		return time.Time{}, errors.New("donation date cannot be in the past")
	}
	return DonationDate, nil
}
