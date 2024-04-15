package users

import (
	"Agen/data"

	"gorm.io/gorm"
)

// func (u *Users) GantiPassword(connection *gorm.DB, newPassword string) (bool, error) {
// 	query := connection.Table("users").Where("hp = ?", u.HP).Update("password", newPassword)
// 	if err := query.Error; err != nil {
// 		return false, err
// 	}

// 	return query.RowsAffected > 0, nil
// }

func Register(connection *gorm.DB, newUser data.Users) (bool, error) {
	// query := connection.Create(&newUser)
	var user data.Users
	query := connection.Where("hp = ?", newUser.HP).Find(&user)
	if err := query.Error; err != nil {
		return false, err
	} else if user.HP != "" {
		return false, nil
	} else {
		err2 := connection.Create(&newUser).Error
		if err2 != nil {
			return false, err2
		}
		return true, nil
	}
}

func Login(connection *gorm.DB, hp string, password string) (data.Users, error) {
	var result data.Users
	err := connection.Find(&result, data.Users{HP: hp, Password: password}).Error
	if err != nil {
		return data.Users{}, err
	}

	return result, nil
}

func Update(connection *gorm.DB, user_old data.Users, user_new data.Users) (bool, error) {
	query := connection.Model(user_old).Updates(&user_new)
	if err := query.Error; err != nil {
		return false, err
	} else if !(query.RowsAffected > 0) {
		return false, nil
	} else {
		return true, nil
	}
}

func DeleteAccount(db *gorm.DB, user data.Users) (bool, error) {
	query := db.Delete(user)
	if err := query.Error; err != nil {
		return false, err
	} else if !(query.RowsAffected > 0) {
		return false, nil
	} else {
		return true, nil
	}
}

func View(db *gorm.DB, hp string) (data.Users, error) {
	var result data.Users
	err := db.Select("hp", "nama", "alamat", "saldo").Find(&result, data.Users{HP: hp}).Error
	if err != nil {
		return data.Users{}, err
	}

	return result, nil
}
