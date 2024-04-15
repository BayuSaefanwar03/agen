package main

import (
	"Agen/data"
	"Agen/database"
	"Agen/printFormat"
	"Agen/topup"
	"Agen/transfer"
	"Agen/users"
	"fmt"
	"time"
)

var db = database.InitMysql()

func main() {
	database.Migrate(db)
	var input int

	for input != 99 {
		fmt.Println("\n-----Agent-----")
		fmt.Println(" [1] Login")
		fmt.Println(" [2] Register")
		fmt.Println(" [3] Forget Password")
		fmt.Println("[99] Exit")
		fmt.Print("==> ")
		printFormat.ScanInt(&input)
		fmt.Printf("")
		switch input {
		case 1:
			if user, success := login(); success {
				Agent(user)

				// fmt.Println("Selamat Datang", user)
			}
		case 2:
			register()
		case 3:
		case 99:
		default:
			fmt.Println("[Error] Wrong Input")
		}
	}
	fmt.Println("Thank You & GoodBye")
}

func Agent(user *data.Users) {
	fmt.Println("\nWelcome", user.Nama)
	var input int
	for input != 99 {
		fmt.Printf("\nSaldo Rp.%s\n", printFormat.ToRP(user.Saldo))
		fmt.Println("===HOME===")
		fmt.Println("[1] View Profile")
		fmt.Println("[2] Update Profile")
		fmt.Println("[3] Delete Account")
		fmt.Println("[4] Top Up")
		fmt.Println("[5] Transfer")
		fmt.Println("[6] History Top Up")
		fmt.Println("[7] History Transfer")
		fmt.Println("[8] View Profile Other User")
		fmt.Println("[99] Logout")
		fmt.Printf("=> ")
		printFormat.ScanInt(&input)
		fmt.Println("")
		switch input {
		case 1:
			ViewProfile(*user, true)
		case 2:
			updateProfile(user)
		case 3:
			if do := deleteAccount(*user); do {
				input = 99
			}
		case 4:
			Topup(user)
		case 5:
			Transfer(user)
		case 6:
			Topup_History(*user)
		case 7:
			Transfer_History(*user)
		case 8:
			View_Other_Profile()
		case 99:
		default:
			fmt.Println("[Error] Wrong Input")
		}
	}
	fmt.Println("GoodBye", user.Nama)
}

func Transfer(user *data.Users) {
	var to string
	var nominal int
	fmt.Println("|--[Transfer]")
	fmt.Print("|--> To number phone account: ")
	printFormat.Scan(&to)
	fmt.Print("|--> Rp.")
	printFormat.ScanInt(&nominal)
	fmt.Println("")

	if to == user.HP {
		fmt.Println("[Error] Can't send to yourself")
		return
	} else if nominal <= 0 {
		fmt.Println("[Error] Nominal cannot be less than one")
		return
	} else if nominal > user.Saldo {
		fmt.Println("[Error] Nominal > Your Saldo")
		return
	}
	receiver, _ := users.View(db, to)
	if receiver.Nama == "" {
		fmt.Println("[Error] The recipient's cellphone number was not found")
		return
	}

	err := transfer.Send(db, user, receiver, nominal)
	if err != nil {
		fmt.Println("[Error]", err)
	} else {
		fmt.Println("Successfully saldo transfer")
	}
}

func Topup(user *data.Users) {
	var nominal int
	fmt.Println("|--[Topup saldo]")
	fmt.Println("|")
	fmt.Println("| press [enter] to cancel")
	fmt.Println("|--Nominal :")
	fmt.Print("|--> Rp.")
	printFormat.ScanInt(&nominal)
	fmt.Println("")
	if nominal < 0 {
		fmt.Println("[Error] Nominal cannot be less than zero")
		return
	} else if nominal == 0 {
		fmt.Println("Cencel")
		return
	}

	success, err := topup.TopUp(db, user, nominal)
	if err != nil {
		fmt.Println("[Error]", err)
	} else if !success {
		fmt.Println("[Error] Failed to add saldo")
	} else {
		fmt.Println("Successfully added")
	}
}

func Topup_History(user data.Users) {
	data := topup.HistoryTopup(db, user)

	if len(data) == 0 {
		fmt.Println("No topup history")
	}

	for i := 0; i < len(data); i++ {
		if i == 0 {
			fmt.Printf("|%-2s|%-10s|%-19s|\n", "ID", "Nominal", "Date")
		}
		fmt.Printf("|%-2d|%-10d|%-19s|\n", data[i].TopupID, data[i].Nominal, data[i].CreatedAt.Format("01/02/2006 03:04:05"))
	}
}

func updateProfile(user *data.Users) {
	var input int
	user_old := *user
	fmt.Println("[--Edit Profile]")
	fmt.Println("[1] Nama")
	fmt.Println("[2] No Hp")
	fmt.Println("[3] Alamat")
	fmt.Println("[4] Password")
	fmt.Println("Press enter to Exit")
	fmt.Println("=> ")
	printFormat.ScanInt(&input)
	fmt.Println("")

	switch input {
	case 1:
		fmt.Print("Nama -> ")
		printFormat.Scan(&user.Nama)
	case 2:
		fmt.Print("No Hp -> ")
		printFormat.Scan(&user.HP)
	case 3:
		fmt.Print("Alamat -> ")
		printFormat.Scan(&user.Alamat)
	case 4:
		fmt.Print("Password -> ")
		printFormat.Scan(&user.Password)
	case 0:
		return
	}
	if user.Nama == "" || user.HP == "" || user.Alamat == "" {
		fmt.Println("[Error] Tidak Boleh Kosong")
	} else {
		user.UpdatedAt = time.Time{}
		success, err := users.Update(db, user_old, *user)
		if err != nil {
			fmt.Println("[Error]", err)
		} else if !success {
			fmt.Println("[Error] Failed to change Profile")
		} else {
			fmt.Println("Successfully changed profile")
		}
	}

}

func ViewProfile(user data.Users, access bool) {
	if access {
		fmt.Println("--[View Profile]")
	} else {
		fmt.Printf("--[See %s's Profile]\n", user.Nama)
	}
	fmt.Println("=> Nama   :", user.Nama)
	fmt.Println("=> No Hp  :", user.HP)
	fmt.Println("=> Alamat :", user.Alamat)

	if access {
		fmt.Println("=> Create At    :", user.CreatedAt.Format("01/03/2006 03:04:05"))
		fmt.Println("=> Last Update  :", user.UpdatedAt.Format("01/03/2006 03:04:05"))
	}
}

func login() (*data.Users, bool) {
	var hp, password string
	fmt.Println("***Login***")
	fmt.Print("=> No Hp : ")
	printFormat.Scan(&hp)
	fmt.Print("=> Password : ")
	printFormat.Scan(&password)

	user, err := users.Login(db, hp, password)
	if err != nil {
		fmt.Println("[Error]", err)
		return &data.Users{}, false
	} else if user.Nama == "" {
		fmt.Println("[Error] Invalid No Hp or Password")
		return &data.Users{}, false
	} else {
		return &user, true
	}

}

func register() {
	var user data.Users
	fmt.Println("***Register***")
	fmt.Print("=> Name : ")
	printFormat.Scan(&user.Nama)
	fmt.Print("=> No Hp : ")
	printFormat.Scan(&user.HP)
	fmt.Print("=> Password : ")
	printFormat.Scan(&user.Password)
	fmt.Print("=> Address : ")
	printFormat.Scan(&user.Alamat)
	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}

	success, err := users.Register(db, user)

	if err != nil {
		fmt.Println("[Error]", err)
	} else if !success {
		fmt.Println("[Error] Nomor Hp telah di Daftarkan")
	} else {
		fmt.Println("Success !!!, Registrasi Berhasil")
	}
}

func deleteAccount(user data.Users) bool {
	var input int
	fmt.Println("-- [Delete Account]")
	fmt.Println("Are you sure? ")
	fmt.Println("[1] Yes ")
	fmt.Println("[2] No ")
	fmt.Print("=> ")
	printFormat.ScanInt(&input)
	fmt.Println("")
	switch input {
	case 1:
		success, err := users.DeleteAccount(db, user)
		if err != nil {
			fmt.Println("[Error]", err)
		} else if !success {
			fmt.Println("[Error] Failed to Delete Account")
		} else {
			fmt.Println(user.Nama, "Deleted Success!!!")
			return true
		}

	case 2:
		fmt.Println("Thank you staying here")
	case 0:
		return false
	}
	return false
}

func Transfer_History(user data.Users) {
	data := transfer.HistoryTransfer(db, user)

	if len(data) == 0 {
		fmt.Println("No topup history")
	}

	for i := 0; i < len(data); i++ {
		if i == 0 {
			fmt.Printf("|%-3s|%-14s|%-14s|%-10s|%-19s|\n", "ID", "Sender", "Receiver", "Nominal", "Date")
		}
		fmt.Printf("|%-3d|%-14s|%-14s|%-10d|%-19s|\n", data[i].TransferID, data[i].HP_Penerima, data[i].HP_Penerima, data[i].Nominal, data[i].CreatedAt.Format("01/02/2006 03:04:05"))
	}
}

func View_Other_Profile() {
	var hp string
	fmt.Println("|--[View Profile Other User]")
	fmt.Print("|--> No Hp : ")
	printFormat.Scan(&hp)
	fmt.Println("")

	receiver, _ := users.View(db, hp)
	if receiver.Nama == "" {
		fmt.Println("[Error] The recipient's cellphone number was not found")
		return
	}

	View_Profile(receiver, false)
}
func View_Profile(user data.Users, access bool) {
	if access {
		fmt.Println("|--[View Profile]")
	} else {
		fmt.Printf("|--[See %s's Profile]\n", user.Nama)
	}
	fmt.Println("|--> Name         :", user.Nama)
	fmt.Println("|--> Phone Number :", user.HP)
	fmt.Println("|--> Address      :", user.Alamat)
	if access {
		fmt.Println("|--> Create At    :", user.CreatedAt.Format("01/02/2006 03:04:05"))
		fmt.Println("|--> Last Update  :", user.UpdatedAt.Format("01/02/2006 03:04:05"))
	}
}
