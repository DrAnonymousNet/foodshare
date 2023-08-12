package core

import "net/http"

type Model interface {
	isModel() bool
}

type SerializerResponseInterface interface{
	Render(w http.ResponseWriter, r *http.Request) error

}

type SerializerInterface interface {
	Save() error
	Update() error
	Bind(r *http.Request) error
	Delete() error
	Validate() error
	ResponsePayload(r *http.Request) SerializerResponseInterface
	GetUpdatedFields() map[string]interface{}

}


