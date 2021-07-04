package model

import "gorm.io/gorm"

// Company is a listing for a company
type Company struct {
	gorm.Model
	CompanyName    string `gorm:"unique;"`
	CompanyWebsite string
	Listings       []Listing
}

// Listing is a single job listing
type Listing struct {
	gorm.Model
	CompanyID     uint
	Approved      bool
	ListingURL    string
	Reported      bool
	Archived      bool
	ScreenshotURL string
}
