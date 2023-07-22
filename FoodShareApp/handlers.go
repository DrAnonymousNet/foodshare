package foodshare

import (
	"net/http"

	core "github.com/DrAnonymousNet/foodshare/Core"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type DonationHandler struct{}

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

func (d *DonationHandler) ListDonations(w http.ResponseWriter, r *http.Request) {
	var donations []Donation
	err := core.DB.Model(&Donation{}).Find(&donations).Error
	if err != nil {
		render.Render(w, r, core.ErrInvalidRequest(err))
		return
	}
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
