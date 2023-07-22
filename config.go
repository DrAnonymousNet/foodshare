package main

import (
	auth "github.com/DrAnonymousNet/foodshare/Auth"
	foodshare "github.com/DrAnonymousNet/foodshare/FoodShareApp"
	notifications "github.com/DrAnonymousNet/foodshare/Notifications"
)

//Model Registration

func getModels() []interface{} {
	Models := []interface{}{
		auth.User{},
		notifications.Notification{},
		foodshare.DonationRequest{},
		foodshare.Donation{},
	}
	return Models
}
