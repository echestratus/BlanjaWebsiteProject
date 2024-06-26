package models

import (
	"gofiber-marketplace/src/configs"

	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	UserID        uint   `json:"user_id" validate:"required"`
	User          User   `gorm:"foreignKey:UserID" validate:"-"`
	Name          string `json:"name" validate:"required,max=50"`
	MainAddress   string `json:"main_address" validate:"required"`
	DetailAddress string `json:"detail_address" validate:"required"`
	Phone         string `json:"phone" validate:"required,numeric,max=15"`
	PostalCode    string `json:"postal_code" validate:"required,numeric,max=8"`
	Primary       string `json:"primary" validate:"required,oneof=on off"`
	City          string `json:"city" validate:"required"`
}

func SelectAllAddresses() []*Address {
	var addresses []*Address
	configs.DB.Preload("User").Find(&addresses)
	return addresses
}

func SelectAddressesbyUserID(user_id int) []*Address {
	var addresses []*Address
	configs.DB.Preload("User").Where("user_id = ?", user_id).Find(&addresses)
	return addresses
}

func SelectAddressbyId(id int) *Address {
	var address Address
	configs.DB.Preload("User").First(&address, "id = ?", id)
	return &address
}

func CreateAddress(address *Address) error {
	result := configs.DB.Create(&address)
	return result.Error
}

func SetOtherAddressesPrimaryOff(userID uint, currentAddressID uint) error {
	query := configs.DB.Model(&Address{}).Where("user_id = ?", userID)

	if currentAddressID != 0 {
		query = query.Where("id != ?", currentAddressID)
	}

	result := query.Update("primary", "off")
	return result.Error
}

func SetPrimaryOnForFirstAddress(userID uint) error {
	var address Address
	result := configs.DB.Where("user_id = ?", userID).First(&address)
	if result.Error != nil {
		return result.Error
	}

	address.Primary = "on"
	return configs.DB.Save(&address).Error
}

func UpdateAddress(id int, updatedAddress *Address) error {
	result := configs.DB.Model(&Address{}).Where("id = ?", id).Updates(updatedAddress)
	return result.Error
}

func DeleteAddress(id int) error {
	result := configs.DB.Delete(&Address{}, "id = ?", id)
	return result.Error
}
