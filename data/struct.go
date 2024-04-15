package data

import "time"

type Users struct {
	HP        string `gorm:"primarykey"`
	Password  string
	Nama      string
	Alamat    string
	Saldo     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Topup struct {
	TopupID   uint `grom:"primarykey"`
	HP        string
	Nominal   int
	CreatedAt time.Time
}

type Transfer struct {
	TransferID  uint `gorm:"primarykey"`
	HP_Pengirim string
	HP_Penerima string
	Nominal     int
	CreatedAt   time.Time
}
