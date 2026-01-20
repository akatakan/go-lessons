package models

type User struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	KullaniciAdi string `json:"kullanici_adi" gorm:"unique" validate:"required,min=3,max=20"`
	Email        string `json:"email" gorm:"unique" validate:"required,email"`
	PasswordHash string `json:"-"`
}
