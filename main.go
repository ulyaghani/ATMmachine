package main

import (
	"bufio"
	"fmt"
	"mesin-atm/controllers"
	"mesin-atm/db"
	"mesin-atm/models"
	"os"
	"strconv"
	"strings"
)

func main() {
	db.ConnectDB()
	for {
		fmt.Println("==== Selamat Datang di Mesin ATM ====")
		fmt.Println("[1] Login")
		fmt.Println("[2] Register")
		fmt.Println("Pilih Menu: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var name, pin string
			fmt.Println("Masukkan Nama anda: ")
			fmt.Scanln(&name)
			fmt.Println("Masukkan PIN anda: ")
			fmt.Scanln(&pin)

			currentUser, err := controllers.Login(name, pin)
			if err != nil {
				fmt.Println("Login gagal:", err)
				continue
			}
			mainMenu(currentUser)
		case 2:
			controllers.RegisterAccount()
		case 3:
			fmt.Println("Pilihan Tidak Valid, Coba Lagi")
		}
	}
}

func mainMenu(currentUser *models.Account) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("===== HALOO!", currentUser.Name, "=====")
		fmt.Println("[1] Cek Saldo")
		fmt.Println("[2] Setor Uang")
		fmt.Println("[3] Tarik Tunai")
		fmt.Println("[4] Transfer")
		fmt.Println("[5] Histori Transaksi")
		fmt.Println("[6] Logout")
		fmt.Print("Pilih Menu: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Input tidak valid, coba lagi.")
			continue
		}

		fmt.Println()

		switch choice {
		case 1:
			controllers.CheckBalance(currentUser)
		case 2:
			var amount float64
			fmt.Print("Masukkan jumlah setor: ")
			fmt.Scanln(&amount)
			controllers.Deposit(currentUser, amount)
		case 3:
			var amount float64
			fmt.Print("Masukkan jumlah tarik: ")
			fmt.Scanln(&amount)
			controllers.Withdraw(currentUser, amount)
		case 4:
			var toAccountID int
			var amount float64
			fmt.Print("Masukkan ID akun penerima: ")
			fmt.Scanln(&toAccountID)
			fmt.Print("Masukkan jumlah transfer: ")
			fmt.Scanln(&amount)
			controllers.Transfer(currentUser, toAccountID, amount)
		case 5:
			controllers.TransactionHistory(currentUser)
		case 6:
			fmt.Println("Logout...")
			return
		default:
			fmt.Println("Pilihan tidak valid")

		}

		fmt.Println()
	}
}
