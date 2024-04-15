package transfer

import (
	"Agen/data"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func Send(connection *gorm.DB, sender *data.Users, receiver data.Users, nominal int) error {

	receiver.Saldo += nominal
	query2 := connection.Table("users").Where("hp = ?", receiver.HP).Update("saldo", receiver.Saldo)
	if err := query2.Error; err != nil || query2.RowsAffected == 0 {
		return fmt.Errorf("failed to change the receiver's saldo")
	}

	sender.Saldo -= nominal
	query1 := connection.Table("users").Where("hp = ?", sender.HP).Update("saldo", sender.Saldo)
	if err := query1.Error; err != nil || query1.RowsAffected == 0 {
		return fmt.Errorf("failed to change the sender's saldo")
	}

	err := connection.Create(&data.Transfer{
		HP_Pengirim: sender.HP,
		HP_Penerima: receiver.HP,
		Nominal:     nominal,
		CreatedAt:   time.Time{}}).Error
	if err != nil {
		return err
	}

	return nil
}

func HistoryTransfer(connection *gorm.DB, user data.Users) []data.Transfer {
	var result []data.Transfer
	connection.Where("hp_pengirim = ?", user.HP).Find(&result)
	return result
}
