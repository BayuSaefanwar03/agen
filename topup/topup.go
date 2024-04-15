package topup

import (
	"Agen/data"
	"time"

	"gorm.io/gorm"
)

func TopUp(db *gorm.DB, user *data.Users, nominal int) (bool, error) {

	err := db.Create(&data.Topup{HP: user.HP, Nominal: nominal, CreatedAt: time.Time{}}).Error
	if err != nil {
		return false, err
	} else {
		user.Saldo += nominal
		query := db.Table("users").Where("hp=?", user.HP).Update("saldo", user.Saldo)
		if err := query.Error; err != nil {
			return false, err
		}

		return query.RowsAffected > 0, nil
	}
}

func HistoryTopup(connection *gorm.DB, user data.Users) []data.Topup {
	var result []data.Topup
	connection.Where("hp = ?", user.HP).Find(&result)
	return result
}
