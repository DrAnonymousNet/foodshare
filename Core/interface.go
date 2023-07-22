package core

import "net/http"

type Model interface {
	isModel() bool
}

type Serializer interface {
	Save() error
	Update(Instance Model) error
	Validate() error
	Bind(r *http.Request) error
	Render(w http.ResponseWriter, r *http.Request) error
}
