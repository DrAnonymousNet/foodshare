package foodshare

import (
	"errors"
	"net/http"

	"strings"

	core "github.com/DrAnonymousNet/foodshare/Core"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

type DonationHandler struct{
	filters map[string][]core.FilterType //Fields and the filter properties like containes, exact, etc.
	model Donation
	queryset *gorm.DB
}

func SetViewSet() *DonationHandler{
	return &DonationHandler{
		filters : map[string][]core.FilterType{
			"donor_id":        {core.Equal},
			"donation_type":   {core.Equal},
			"donation_date":   {core.Equal, core.LessEqual, core.GreaterEqual, core.GreaterThan, core.LessThan},
			"pickup_address":  {core.Equal},
			"item_description": {core.Contains, core.IContains},
			"title":           {core.Contains, core.IContains},
		},
		model: Donation{},
		queryset: core.DB.Model(&Donation{}),
	}
}


func (d *DonationHandler)ParseFilter(key string) error{
	splitResult := strings.Split(key, "__")
	fieldName, filter := splitResult[0], splitResult[1]
	//Does the Model support this filter?
	_, ok := d.filters[fieldName]
	if !ok{
		return errors.New("invalid filter")
	}
	//Does the filter exist for this field?
	filterTypeExist := false
	filters := d.filters[fieldName]
	for _, filterType := range filters{
		if filter == string(filterType){
			filterTypeExist = true
			break
		}
	}
	if !filterTypeExist{
		return errors.New("filter type {filter} does not exist for this field")
	}
	if fieldName != "" && filter != ""{
		d.queryset = d.queryset.Where(fieldName + " = ?", filter)
	}else if fieldName != "" && filter == ""{
		d.queryset = d.queryset.Where(fieldName + " = ?", "equal")
	}
	return nil

}

// @Accept json
// @Produce json
// @Param requestBody body DonationSerializer true "Donation"
// @Success 200 {object} DonationSerializer
// @Summary Create a donation
// @Description create a donation
// @Router /donations [post]
func (d *DonationHandler) CreateDonation(w http.ResponseWriter, r *http.Request) {
	requestBody := &DonationSerializer{}
	if err := render.Bind(r, requestBody); err != nil {
		render.Render(w, r, core.ErrInvalidRequest(err))
		return
	}
	err := requestBody.Save(r)
	if err != nil {
		render.Render(w, r, core.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, requestBody)

}

// @Summary Get a donation
// @Description get a donation
// @Produce json
// @Param uid path string true "Donation UID"
// @Success 200 {object} DonationSerializer
// @Router /donations/{uid} [get]
func (d *DonationHandler) GetDonation(w http.ResponseWriter, r *http.Request) {
	donation_object, err := d.getObject(r)
	if err != nil {
		render.Render(w, r, core.ErrNotFound)
		return
	}
	serializer := DonationSerializer{Instance: &donation_object}
	render.Status(r, http.StatusOK)
	render.Render(w, r, &serializer)
}

// @Accept json
// @Summary Update a donation
// @Description update a donation
// @Produce json
// @Param uid path string true "Donation UID"
// @Success 200
// @Router /donations/{uid} [patch]
func (d *DonationHandler) UpdateDonation(w http.ResponseWriter, r *http.Request) {
	donation_object, err := d.getObject(r)
	if err != nil {
		render.Render(w, r, core.ErrNotFound)
		return
	}
	requestBody := &DonationSerializer{Instance: &donation_object}
	if err := render.Bind(r, requestBody); err != nil {
		render.Render(w, r, core.ErrInvalidRequest(err))
		return
	}
	err = requestBody.Update()
	if err != nil {
		render.Render(w, r, core.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, requestBody)
}

// @Summary Delete a donation
// @Description delete a donation
// @Produce json
// @Param uid path string true "Donation UID"
// @Success 204
// @Router /donations/{uid} [delete]
func (d *DonationHandler) DeleteDonation(w http.ResponseWriter, r *http.Request) {
	donation_object, err := d.getObject(r)
	if err != nil {
		render.Render(w, r, core.ErrNotFound)
		return
	}
	err = core.DB.Model(&Donation{}).Delete(&donation_object).Error
	if err != nil {
		render.Render(w, r, core.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusNoContent)
}
// @Summary List donations
// @Description list donations
// @Produce json
// @Success 200 {object} DonationSerializer
// @Router /donations [get]
func (d *DonationHandler) ListDonations(w http.ResponseWriter, r *http.Request) {
	var donations []Donation
	qs, err := d.getQuerySet(r)
	if err != nil {
		render.Render(w, r, core.ErrInvalidRequest(err))
		return
	}
	qs.Find(&donations)
	serializer := DonationSerializer{Instances: donations}
	render.Status(r, http.StatusOK)
	render.Render(w, r, &serializer)
}

func (d *DonationHandler) getObject(r *http.Request) (Donation, error) {
	uid := chi.URLParam(r, "uid")
	var donation Donation
	err := core.DB.Model(&Donation{}).Where("uid = ?", uid).First(&donation).Error
	return donation, err
}

func (d *DonationHandler)getQuerySet(r *http.Request) (*gorm.DB, error) {
	filters := r.URL.Query()
	for key := range filters{
		err := d.ParseFilter(key)
		if err != nil{
			return nil , err
		}
	}
	return d.queryset, nil
}