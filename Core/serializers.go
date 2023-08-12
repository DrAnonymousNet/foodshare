package core

// import (
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"reflect"

// 	"github.com/go-playground/validator"
// )



// type Serializer struct {
// 	Instance *Model
// 	Instances []Model
// 	Meta SerializerMeta
// }

// type SerializerMeta struct {
// 	ModelStruct Model
// }

// type SerializerResponse struct{
// 	Serializer
// }


// // This function retrieves the fields that needs to be renders in the Json Response.
// // It creates the DonationResponse objects from the data in the DonationSerializer.
// func (serializer Serializer) ResponsePayload(r *http.Request) SerializerResponse {

// 	// //user := r.Context().Value("User").(*auth.User)
// 	// if serializer.Instance != nil {
// 	// 	serializerInstanceValue := reflect.ValueOf(serializer.Instance).Elem()
// 	// 	serializerType := reflect.TypeOf(serializer)
// 	// 	if serializerInstanceValue.Kind() != reflect.Struct {
// 	// 		for i := 0; i < serializerType.NumField(); i++ {
// 	// 			field := serializerType.Field(i)
// 	// 			fieldName := field.Name
	
// 	// 			if fieldName == "Model" {
// 	// 				continue
// 	// 			}
// 	// 			if fieldName == "User" {
// 	// 				continue
// 	// 			}
// 	// 			value := reflect.ValueOf(serializer).Elem().FieldByName(fieldName)
// 	// 			fmt.Printf("Field Name: %s, Value: %v, Type: %s\n", fieldName, value.Interface(), value.Type())
	
// 	// 			if value.IsValid() {
// 	// 				fmt.Println("Reached here")
// 	// 				serializerField := serializerInstanceValue.FieldByName(fieldName)
// 	// 				if serializerField.CanSet() {
// 	// 					serializerField.Set(value)
// 	// 				}
// 	// 			}
// 	// 		}
// 	// 	}		

// 	// }

// 	 responsePayload := SerializerResponse{
// 	 	Serializer: serializer,
// 	 }
// 	// 	UID:                serializer.UID,
// 	// 	DonorUID:           serializer.Instance.DonorUID,
// 	// }
// 	return responsePayload
// }

// // This Binds the data that are not sent as part of request body but are
// // required to create the model associated with the serializer. It also calls
// // the validate method.
// func (serializer Serializer) Bind(r *http.Request) error {
// 	// The Bind method calls the validate method to validate the stuct fields as describe above

// 	// err := serializer.Validate(r)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	//Donor := r.Context().Value("User").(*auth.User)
// 	//serializer.DonorUID = Donor.UID

// 	return nil
// }

// func (serializer Serializer) Validate(r *http.Request) error {
// 	errorMessage := ""
// 	if serializer.Instance == nil {
// 	validate := validator.New()
// 	if err := validate.Struct(serializer); err != nil {
// 		validationErrors := err.(validator.ValidationErrors)
// 		// Return validation errors to the client
// 		for _, err := range validationErrors {
// 			errorMessage += fmt.Sprintf("Field %s: Validation Error (%s) \n", err.Field(), err.Tag())
// 		}
// 	}
// 	}

// 	// Validate Fields in the struct
// 	serializerValue := reflect.ValueOf(serializer).Elem() // Get the actual struct value
// 	serializerType := serializerValue.Type()       // Get the type of the struct

// 	// Loop through the fields of the struct
// 	for i := 0; i < serializerType.NumField(); i++ {
// 		field := serializerType.Field(i)
// 		fieldName := field.Name

// 		if fieldName == "Instance" || fieldName == "Instances" {
// 			continue
// 		}
// 		if fieldName == "User" {
// 			continue
// 		}
// 		log.Println(fieldName, field)
// 		value := serializerValue.Field(i)

// 		// Check if the field is valid and has a validate{field name} method
// 		method := reflect.ValueOf(serializer).MethodByName("Validate" + fieldName)
// 		if method.IsValid() && method.Type().NumIn() == 1 {
// 			log.Panicln(value, "FRMVRVVT")
// 			err := method.Call([]reflect.Value{value})
// 			log.Println(err[1])
// 			if err[1].Interface() != nil { // Check if the error is not nil
// 				errorMessage += fmt.Sprintf("Field %s: Validation Error (%s) \n", fieldName, err[1].Interface())
// 			}
// 		}
// 	}

// 	if errorMessage != "" {
// 		return errors.New(errorMessage)
// 	}
// 	return nil
// }


// // Update the donation object
// // During the update operation, an Instance of the donation will be passed in the serializer
// func (serializer Serializer) Update() error {
// 	if serializer.Instance == nil {
// 		return errors.New("an instance must be passed during update")
// 	}

// 	serializerValue := reflect.ValueOf(serializer.Instance).Elem()
// 	serializerType := reflect.TypeOf(serializer)
// 	if serializerValue.Kind() != reflect.Struct {
// 		for i := 0; i < serializerType.NumField(); i++ {
// 			field := serializerType.Field(i)
// 			fieldName := field.Name

// 			if fieldName == "Model" {
// 				continue
// 			}
// 			if fieldName == "User" {
// 				continue
// 			}
// 			value := reflect.ValueOf(serializer).Elem().FieldByName(fieldName)
// 			fmt.Printf("Field Name: %s, Value: %v, Type: %s\n", fieldName, value.Interface(), value.Type())

// 			if value.IsValid() {
// 				fmt.Println("Reached here")
// 				modelField := serializerValue.FieldByName(fieldName)
// 				if modelField.CanSet() {
// 					modelField.Set(value)
// 				}
// 			}
// 		}
// 	}
// 	updatedField := serializer.GetUpdatedFields()
// 	err := DB.Model(&serializer.Instance).Updates(updatedField).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (serializer Serializer) Delete() error {
// 	// Delete the donation object
// 	err := DB.Model(&serializer.Instance).Delete(&serializer.Instance).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Retrieves the fields sent in the request body of the update request
// func (serializer Serializer) GetUpdatedFields() map[string]interface{} {
// 	data := map[string]interface{}{

// 	}
// 	serializerType := reflect.TypeOf(serializer)
// 	serializerValue := reflect.ValueOf(serializer)
	
// 	// Get the none readonly fields that can be sent in the request body
// 	for i := 0; i < serializerType.NumField(); i++ {
// 		field := serializerType.Field(i)
// 		tag := field.Tag.Get("json")
// 		if tag != "-" {
// 			data[tag] = serializerValue.Field(i).Interface()
// 		}
// 	}

// 	// Remove the fields that are not sent in the request body
// 	for key, value := range data {
// 		log.Println(key, value)
// 		val := reflect.ValueOf(value)
// 		if !val.IsValid() || reflect.DeepEqual(value, reflect.Zero(val.Type()).Interface()) {
// 			delete(data, key)
// 		}
// 	}
// 	return data

// }

// // Creates the instance of the model associated with the serializer
// func (serializer *Serializer) Save(r *http.Request, serializerD SerializerInterface) error {
// 	data := map[string]interface{}{}

// 	serializerValue := reflect.ValueOf(serializerD).Elem()
// 	serializerType := reflect.TypeOf(serializerD).Elem()
// 	log.Println(serializerType.NumField(), "serializerType.NumField()")
// 	for i := 0; i < serializerType.NumField(); i++ {
// 		field := serializerType.Field(i)
// 		tag := field.Tag.Get("json")
// 		if tag != "-" {
// 			data[tag] = serializerValue.Field(i).Interface()
// 		}
// 	}

// 	// Create an instance of the concrete type that implements the Model interface
// 	modelType := reflect.TypeOf(serializerD.Meta.ModelStruct)
// 	log.Println(modelType, "mODEL TYPE", serializerD)
// 	instance := reflect.New(modelType.Elem()).Interface()
// 	reflectValue := reflect.ValueOf(instance).Elem()

// 	for key, value := range data {
// 		fieldName := ""
// 		for i := 0; i < reflectValue.NumField(); i++ {
// 			tag := reflectValue.Type().Field(i).Tag.Get("json")
// 			if tag == key {
// 				fieldName = reflectValue.Type().Field(i).Name
// 				break
// 			}
// 		}

// 		if fieldName == "" {
// 			continue
// 		}

// 		// Access the unexported field and set its value
// 		unexportedField := reflectValue.FieldByName(fieldName)
// 		if unexportedField.IsValid() && unexportedField.CanSet() {
// 			unexportedField.Set(reflect.ValueOf(value))
// 		}
// 	}

// 	err := DB.Create(instance).Error
// 	if err != nil {
// 		return err
// 	}
// 	serializerD.Instance = instance.(*Model)

// 	return nil
// }



// func (response SerializerResponse) Render(w http.ResponseWriter, r *http.Request) error {

// 	return nil
// }

// func ListResponsePayload(r *http.Request, models []Model) []SerializerResponse {
// 	var responsePayload []SerializerResponse
// 	for _, model := range models {
// 		serializerResponse := SerializerResponse{}
// 		serializerValue := reflect.ValueOf(serializerResponse).Elem()
// 		serializerType := reflect.TypeOf(serializerResponse)

// 		// Copy the model data into the serializer
// 		for j := 0; j < serializerType.NumField(); j++ {
// 			field := serializerType.Field(j)
// 			tag := field.Tag.Get("json")
// 			if tag != "-" {
// 				serializerValue.Field(j).Set(reflect.ValueOf(model).FieldByName(field.Name))
// 			}
// 		}

// 		responsePayload = append(responsePayload, serializerResponse)
// 	}
// 	return responsePayload

// }