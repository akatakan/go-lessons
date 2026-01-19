package models

type User struct {
	ID           int    `json:"id"`
	KullaniciAdi string `json:"kullanici_adi"`
	Email        string `json:"email"`
}
