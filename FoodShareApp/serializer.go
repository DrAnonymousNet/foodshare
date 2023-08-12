package foodshare

import (
	"errors"
	"fmt"
	"log"

	"net/http"
	"reflect"
	"time"

	auth "github.com/DrAnonymousNet/foodshare/Auth"
	core "github.com/DrAnonymousNet/foodshare/Core"
	"github.com/go-playground/validator"


	"github.com/google/uuid"
)

type DonationRequestSerializer struct {
	UID                uuid.UUID 
	RequestorID        uint8
	User               auth.User 
	RequestDescription string
	Quantity           uint8
	RequestDate        time.Time
	DeliveryAddress    string
	RequestStatus      string //RequestStatusType `sql:"type:request_status_type"`
	RequestFrom        string //RequestFromType   `sql:"type:request_from_type"`
}
type DonationSerializer struct {
	UID             uuid.UUID  `json:"-"`
	Title           string     `json:"title" validate:"required"`
	DonorUID        uuid.UUID  `json:"-" serializer:"readonly"`
	DonatedObjType  string     `json:"donation_obj_type"`
	DonationDate    time.Time  `json:"donation_date" validate:"required"`
	PickUpAddress   string     `json:"pick_up_address" validate:"required"`
	ItemDescription string     `json:"item_description"`
	Instance        *Donation  `json:"-"`
	Instances       []Donation `json:"-"`
	// Also note that we only include a pointer to the original Donation struct in our Request struct.
	// This indirection avoids having to allocate a new copy of Donation.
}



// This renders the data to be rendered and adds the read only fields.
// It implements the Render method and compose of the Serializer class
type DonationResponse struct {
	DonationSerializer
	UID      uuid.UUID `json:"uid"`
	DonorUID uuid.UUID `json:"donor_uid"`
}

// This function retrieves the fields that needs to be renders in the Json Response.
// It creates the DonationResponse objects from the data in the DonationSerializer.
func (d *DonationSerializer) ResponsePayload(r *http.Request) DonationResponse {

	//user := r.Context().Value("User").(*auth.User)
	if d.Instance != nil {
		d.UID = d.Instance.UID
		d.Title = d.Instance.Title
		d.DonatedObjType = d.Instance.DonatedObjType
		d.DonationDate = d.Instance.DonationDate
		d.PickUpAddress = d.Instance.PickUpAddress
		d.ItemDescription = d.Instance.ItemDescription
		d.DonorUID = d.Instance.DonorUID
	}

	responsePayload := DonationResponse{
		DonationSerializer: *d,
		UID:                d.UID,
		DonorUID:           d.Instance.DonorUID,
	}
	return responsePayload
}

// This Binds the data that are not sent as part of request body but are
// required to create the model associated with the serializer. It also calls
// the validate method.
func (d *DonationSerializer) Bind(r *http.Request) error {
	// The Bind method calls the validate method to validate the stuct fields as describe above

	err := d.Validate(r)
	if err != nil {
		return err
	}
	Donor := r.Context().Value("User").(*auth.User)
	d.DonorUID = Donor.UID

	return nil
}

func (d *DonationSerializer) Validate(r *http.Request) error {

	errorMessage := ""
	// This is not an update operation
	if d.Instance == nil{
	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		// Return validation errors to the client
		for _, err := range validationErrors {
			errorMessage += fmt.Sprintf("Field %s: Validation Error (%s) \n", err.Field(), err.Tag())
		}
	}
	}

	// Validate Fields in the struct
	donationValue := reflect.ValueOf(d).Elem() // Get the actual struct value
	donationType := donationValue.Type()       // Get the type of the struct
	updatedFields := d.GetUpdatedFields()
	// Loop through the fields of the struct
	for i := 0; i < donationType.NumField(); i++ {
		field := donationType.Field(i)
		fieldName := field.Name

		if fieldName == "Instance" || fieldName == "Instances" {
			continue
		}
		if fieldName == "User" {
			continue
		}
		log.Println(fieldName, field)
		value := donationValue.Field(i)

		// Check if the field is valid and has a validate{field name} method
		method := reflect.ValueOf(d).MethodByName("Validate" + fieldName)
		_, ok := updatedFields[fieldName]
		if method.IsValid() && method.Type().NumIn() == 1 && ok {
			err := method.Call([]reflect.Value{value})
			log.Println(err[1])
			if err[1].Interface() != nil { // Check if the error is not nil
				errorMessage += fmt.Sprintf("Field %s: Validation Error (%s) \n", fieldName, err[1].Interface())
			}
		}
	}

	if errorMessage != "" {
		return errors.New(errorMessage)
	}
	return nil
}

func (d *DonationSerializer) ValidateDonationDate(DonationDate time.Time) (time.Time, error) {
	if d.DonationDate.Before(time.Now()) {
		return time.Time{}, errors.New("donation date cannot be in the past")
	}
	return DonationDate, nil
}

// Update the donation object
// During the update operation, an Instance of the donation will be passed in the serializer
func (d *DonationSerializer) Update() error {
	if d.Instance == nil {
		return errors.New("an instance must be passed during update")
	}

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
	updatedField := d.GetUpdatedFields()
	err := core.DB.Model(&d.Instance).Updates(updatedField).Error
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

// // Retrieves the fields sent in the request body of the update request
 func (d *DonationSerializer) GetUpdatedFields() map[string]interface{} {
	data := map[string]interface{}{
		"uid":              d.UID,
		"title":            d.Title,
		"donor_uid":        d.DonorUID,
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

// // Creates the instance of the model associated with the serializer
 func (d *DonationSerializer) Save(r *http.Request) error {
	//Donor := r.Context().Value("User").(*auth.User)
	//d.DonorUID = Donor.UID

	donationRequest := Donation{
		UID:             uuid.New(),
		Title:           d.Title,
		DonatedObjType:  d.DonatedObjType,
		DonationDate:    d.DonationDate,
		PickUpAddress:   d.PickUpAddress,
		ItemDescription: d.ItemDescription,
		DonorUID:        d.DonorUID,
	}
	err := core.DB.Model(&Donation{}).Create(&donationRequest).Error
	if err != nil {
		return err
	}
	d.UID = donationRequest.UID
	d.Instance = &donationRequest

	return nil
}
func (response DonationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ListResponsePayload(r *http.Request, donations []Donation) []DonationResponse {
	var responsePayload []DonationResponse
	for _, donation := range donations {
		donationSerializer := DonationSerializer{
			UID:             donation.UID,
			Title:           donation.Title,
			DonorUID:        donation.DonorUID,
			DonatedObjType:  donation.DonatedObjType,
			DonationDate:    donation.DonationDate,
			PickUpAddress:   donation.PickUpAddress,
			ItemDescription: donation.ItemDescription,
			Instance:        &donation,
		}
		responsePayload = append(responsePayload, donationSerializer.ResponsePayload(r))
	}
	return responsePayload
}
