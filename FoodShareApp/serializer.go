package foodshare

import (
	"errors"
	"fmt"
	"log"

	//"fmt"
	"net/http"
	"reflect"
	"time"

	auth "github.com/DrAnonymousNet/foodshare/Auth"
	core "github.com/DrAnonymousNet/foodshare/Core"
	//"gorm.io/gorm"

	//"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type DonationSerializer struct {
	UID             uuid.UUID  `json:"-"`
	Title           string     `json:"title" validate:"required"`
	DonorID         uint       `json:"-"`
	DonatedObjType  string     `json:"donation_obj_type"`
	DonationDate    time.Time  `json:"donation_date" validate:"required"`
	PickUpAddress   string     `json:"pick_up_address" validate:"required"`
	ItemDescription string     `json:"item_description"`
	Instance        *Donation  `json:"-"`
	Instances       []Donation `json:"-"`
	// Also note that we only include a pointer to the original Donation struct in our Request struct.
	// This indirection avoids having to allocate a new copy of Donation.
}

type DonationResponse struct {
	DonationSerializer
	UID      uuid.UUID `json:"uid"`
	DonorUID uuid.UUID `json:"donor_uid"`
}

func (d *DonationSerializer) ResponsePayload(r *http.Request) DonationResponse {
	user := r.Context().Value("User").(*auth.User)
	responsePayload := DonationResponse{
		DonationSerializer: *d,
		UID:                d.UID,
		DonorUID:           user.UID,
	}
	return responsePayload
}
func (d *DonationSerializer) Bind(r *http.Request) error {
	// The Bind method calls the validate method to validate the stuct fields as describe above

	//err := d.Validate(r)
	//if err != nil {
	//	return err
	//}
	Donor := r.Context().Value("User").(*auth.User)
	d.DonorID = Donor.ID

	return nil
}

func (d *DonationSerializer) Validate(r *http.Request) error {
	// Validate the tags in the request body
	// The validation of the serializer fields includes the validation from the validate struct tags
	// and the validate{field name} that is implemented on the struct.

	//errorMessage := ""

	//validate := validator.New()
	//if err := validate.Struct(d); err != nil {
	//	validationErrors := err.(validator.ValidationErrors)
	// Return validation errors to the client
	//	for _, err := range validationErrors {
	//		errorMessage += fmt.Sprintf("Field %s: Validation Error (%s) \n", err.Field(), err.Tag())
	//	}
	//}

	// Validate Fields in of the struct
	//donationSerializer := reflect.TypeOf(*d)
	// Check if the struct is a struct
	//if donationSerializer.Kind() == reflect.Struct {
	// Loop through the fields of the struct
	//	for i := 0; i < donationSerializer.NumField() - 1; i++ {

	//		field := donationSerializer.Field(i)
	//		fieldName := field.Name
	//		value := reflect.ValueOf(donationSerializer).FieldByName(fieldName)

	// Check if the field is valid and has a validate{field name} method
	//		method := reflect.ValueOf(donationSerializer).MethodByName("Validate" + fieldName)
	//		if method.IsValid() && method.Type().NumIn() == 1 {
	//			err := method.Call([]reflect.Value{value})
	//			if err != nil {
	//				errorMessage += fmt.Sprintf("Field %s: Validation Error (%s) \n", fieldName, err)
	//			}
	//		}
	//	}

	//}
	//if errorMessage != "" {
	//	return errors.New(errorMessage)
	//}
	return nil

}

func (d *DonationSerializer) ValidateDonorID(ID int) (int, error) {
	// Validate the DonorID field
	var user auth.User
	core.DB.Model(&auth.User{}).Where("ID = ?", ID)
	if user.ID == 0 {
		return 0, errors.New("donor does not exist")
	}
	return ID, nil
}

func (d *DonationSerializer) ValidateDonationDate(DonationDate time.Time) (time.Time, error) {
	if d.DonationDate.Before(time.Now()) {
		return time.Time{}, errors.New("donation date cannot be in the past")
	}
	return DonationDate, nil
}
func (d *DonationSerializer) Update() error {
	// Update the donation object
	// During the update operation, an Instance of the donation will be passed in the serializer

	donationValue := reflect.ValueOf(d.Instance).Elem()
	serializerType := reflect.TypeOf(d).Elem()
	if donationValue.Kind() != reflect.Struct {
		for i := 0; i < serializerType.NumField(); i++ {
			field := serializerType.Field(i)
			fieldName := field.Name

			if fieldName == "Model" {
				continue
			}
			if fieldName == "User" {
				continue
			}
			value := reflect.ValueOf(d).Elem().FieldByName(fieldName)
			fmt.Printf("Field Name: %s, Value: %v, Type: %s\n", fieldName, value.Interface(), value.Type())

			if value.IsValid() {
				fmt.Println("Reached here")
				donationField := donationValue.FieldByName(fieldName)
				if donationField.CanSet() {
					donationField.Set(value)
				}
			}
		}
	}
	dataMap := d.GetDataMap()
	err := core.DB.Model(&d.Instance).Updates(dataMap).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DonationSerializer) Delete() error {
	// Delete the donation object
	err := core.DB.Model(&d.Instance).Delete(&d.Instance).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DonationSerializer) GetDataMap() map[string]interface{} {
	data := map[string]interface{}{
		"uid":              d.UID,
		"title":            d.Title,
		"donor_id":         d.DonorID,
		"donated_obj_type": d.DonatedObjType,
		"donation_date":    d.DonationDate,
		"pick_up_address":  d.PickUpAddress,
		"item_description": d.ItemDescription,
	}
	for key, value := range data {
		val := reflect.ValueOf(value)
		if !val.IsValid() || reflect.DeepEqual(value, reflect.Zero(val.Type()).Interface()) {
			delete(data, key)
		}
	}
	return data

}

func (d *DonationSerializer) Save(r *http.Request) error {
	//var user auth.User
	//core.DB.Model(&auth.User{}).Where("ID = ?", d.DonorID)

	donationRequest := Donation{
		UID:             uuid.New(),
		Title:           d.Title,
		DonorID:         d.DonorID,
		DonatedObjType:  d.DonatedObjType,
		DonationDate:    d.DonationDate,
		PickUpAddress:   d.PickUpAddress,
		ItemDescription: d.ItemDescription,
		//User:            &user,
	}
	err := core.DB.Model(&Donation{}).Create(&donationRequest).Error
	if err != nil {
		return err
	}
	d.UID = donationRequest.UID

	return nil
}
func (d *DonationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if d.Instance != nil {
		donationObj := reflect.TypeOf(d.Instance).Elem()
		if donationObj.Kind() != reflect.Struct {
			log.Println("Invalid type:", donationObj)
			return errors.New("invalid type")
		}

		donationValue := reflect.ValueOf(d.Instance).Elem()
		donationSerializer := reflect.ValueOf(d).Elem()

		for i := 0; i < donationObj.NumField(); i++ {
			field := donationObj.Field(i)
			fieldName := field.Name

			// Skip the "Model" field to avoid the panic
			if fieldName == "Model" {
				continue
			}
			if fieldName == "User" {
				//donationSerializer.FieldByName("DonorID").Set(donationValue.FieldByName("User").FieldByName("ID"))
				continue
			}

			// Get the value of the field in the donation object
			value := donationValue.FieldByName(fieldName)
			fmt.Printf("Field Name: %s, Value: %v, Type: %s\n", fieldName, value.Interface(), value.Type())

			if value.IsValid() && value.CanSet() {
				// Set the value of the field in the donation response
				donationSerializer.FieldByName(fieldName).Set(value)
			}
		}

		// Special handling for gorm.Model
		//if modelField := donationValue.FieldByName("Model"); modelField.IsValid() && modelField.Type() == reflect.TypeOf(gorm.Model{}) {
		//	donationSerializer.FieldByName("ID").Set(modelField.FieldByName("ID"))
		//	donationSerializer.FieldByName("CreatedAt").Set(modelField.FieldByName("CreatedAt"))
		//	donationSerializer.FieldByName("UpdatedAt").Set(modelField.FieldByName("UpdatedAt"))
		//	donationSerializer.FieldByName("DeletedAt").Set(modelField.FieldByName("DeletedAt"))
		//}
	}

	return nil
}
