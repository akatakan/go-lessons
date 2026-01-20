package models

type User struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	KullaniciAdi string `json:"kullanici_adi" gorm:"unique"`
	Email        string `json:"email" gorm:"unique"`
	PasswordHash string `json:"-"`
}
