package models

import (
	"encoding/json"
	"time"
)

type SwaggerDefault struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            interface{}
}

type SwaggerListHotels struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	limit           string
	total           string
	page            string
	data            []Hotels
}

type SwaggerDetailHotels struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            Hotels
}

type SwaggerCreateHotels struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            Hotels
}

//tmp
type UserAppAccount struct {
	Id        int
	RoomId    int
	BookingId int
	LoginName string
	Status    int8
}

//booking
type SwaggerListBooking struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            []BookingData
	limit           string
	total           string
	page            string
}

type SwaggerDetailBooking struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            BookingData
}

type SwaggerCreateBooking struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            BookingData
}

//customer
type SwaggerListGuest struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            []Guests
}

type SwaggerGuest struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            Guests
}

//room
type Roomdata struct {
	Id                int `orm:"column(room_id)"`
	HotelId           int
	RoomTypeId        int
	BedTypeId         int
	RoomNumber        string
	FloorNumber       string
	Price             string
	RoomSize          string
	Status            int8
	MaxCapacity       int
	NumberOfBed       int
	NumberOfExtraBed  int
	SmokingAllow      int
	PetAllow          int
	Files             string
	RoomMultiLanguage map[string]roomMultiLanguage
	DeletedAt         int8
	CreatedUser       int
	UpdatedUser       int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type roomMultiLanguage struct {
	RoomName            string
	Description         string
	BathroomDescription string
}

type SwaggerCreateRoom struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            Roomdata
}

type SwaggerDetailRoom struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            Roomdata
}

type SwaggerListRoom struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            []Roomdata
}

//admin account
type SwaggerListAccount struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            []AdminAccounts
}

type SwaggerCreateAccount struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            AdminAccounts
}

type SwaggerDetailAccount struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            AdminAccounts
}

//common
type SwaggerGetConfig struct {
	success         string
	code            string
	message         string
	detailErrorCode []string
	data            json.RawMessage //{"config_key1": "value1", "config_key2": "value2"}
}
