package controllers

import (
	"fmt"
	"mesin-atm/db"
	"mesin-atm/models"
	"strings"
	"time"
)

func RegisterAccount() {
	var name, pin string

	fmt.Print("Masukkan nama: ")
	fmt.Scanln(&name)

	fmt.Print("Masukkan PIN (4 digit): ")
	fmt.Scanln(&pin)

	if len(pin) != 4 {
		fmt.Println("PIN harus 4 digit!")
		return
	}

	account := models.Account{
		Name:      name,
		Pin:       pin,
		Balance:   0,
		CreatedAt: time.Now(),
	}

	result := db.Conn.Create(&account)

	if result.Error != nil {
		fmt.Println("Gagal registrasi:", result.Error)
		return
	}

	fmt.Println("Registrasi berhasil! Nomor akun Anda adalah:", account.ID)
}

var currentUser *models.Account

// Fungsi login untuk memverifikasi akun

func Login(name string, pin string) (*models.Account, error) {
	var account models.Account

	// Normalisasi input
	name = strings.ToLower(strings.TrimSpace(name))
	pin = strings.TrimSpace(pin)

	// Query dengan LOWER agar tidak case sensitive
	result := db.Conn.Where("LOWER(name) = ? AND pin = ?", name, pin).First(&account)
	if result.Error != nil {
		return nil, fmt.Errorf("login gagal: %v", result.Error)
	}

	currentUser = &account
	return &account, nil
}

// controllers/account_controller.go
func CheckBalance(account *models.Account) {
	if currentUser == nil {
		fmt.Println("Anda belum login")
		return
	}

	// Tampilkan saldo pengguna yang sedang login
	fmt.Printf("Saldo Anda: %.2f\n", currentUser.Balance)
}

// controllers/account_controller.go
func Deposit(account *models.Account, amount float64) {
	if currentUser == nil {
		fmt.Println("Anda belum login")
		return
	}

	// Tambahkan saldo
	currentUser.Balance += amount
	result := db.Conn.Save(&currentUser) // Update saldo di database
	if result.Error != nil {
		fmt.Println("Gagal setor tunai:", result.Error)
		return
	}

	// Catat transaksi deposit
	transaction := models.Transaction{
		AccountID: currentUser.ID,
		Type:      "deposit",
		Amount:    amount,
		CreatedAt: time.Now(),
	}
	db.Conn.Create(&transaction)

	fmt.Println("Setoran tunai berhasil! Saldo Anda sekarang:", currentUser.Balance)
}

// controllers/account_controller.go
func Withdraw(account *models.Account, amount float64) {
	if currentUser == nil {
		fmt.Println("Anda belum login")
		return
	}

	// Pastikan saldo mencukupi
	if currentUser.Balance < amount {
		fmt.Println("Saldo Anda tidak mencukupi")
		return
	}

	// Kurangi saldo
	currentUser.Balance -= amount
	result := db.Conn.Save(&currentUser) // Update saldo di database
	if result.Error != nil {
		fmt.Println("Gagal tarik tunai:", result.Error)
		return
	}

	// Catat transaksi withdraw
	transaction := models.Transaction{
		AccountID: currentUser.ID,
		Type:      "withdraw",
		Amount:    amount,
		CreatedAt: time.Now(),
	}
	db.Conn.Create(&transaction)

	fmt.Println("Tarik tunai berhasil! Saldo Anda sekarang:", currentUser.Balance)
}

// controllers/account_controller.go
func Transfer(account *models.Account, toAccountID int, amount float64) {
	if currentUser == nil {
		fmt.Println("Anda belum login")
		return
	}

	// Pastikan saldo mencukupi
	if currentUser.Balance < amount {
		fmt.Println("Saldo Anda tidak mencukupi untuk transfer")
		return
	}

	// Cari akun penerima
	var recipient models.Account
	result := db.Conn.First(&recipient, toAccountID)
	if result.Error != nil {
		fmt.Println("Akun penerima tidak ditemukan:", result.Error)
		return
	}

	// Kurangi saldo pengirim
	currentUser.Balance -= amount
	db.Conn.Save(&currentUser) // Update saldo pengirim

	// Tambah saldo penerima
	recipient.Balance += amount
	db.Conn.Save(&recipient) // Update saldo penerima

	// Catat transaksi transfer_out untuk pengirim
	transactionOut := models.Transaction{
		AccountID: currentUser.ID,
		Type:      "transfer_out",
		Amount:    amount,
		TargetID:  &toAccountID,
		CreatedAt: time.Now(),
	}
	db.Conn.Create(&transactionOut)

	// Catat transaksi transfer_in untuk penerima
	transactionIn := models.Transaction{
		AccountID: recipient.ID,
		Type:      "transfer_in",
		Amount:    amount,
		TargetID:  &currentUser.ID,
		CreatedAt: time.Now(),
	}
	db.Conn.Create(&transactionIn)

	fmt.Println("Transfer berhasil! Saldo Anda sekarang:", currentUser.Balance)
}

// Fungsi untuk menampilkan riwayat transaksi
func TransactionHistory(account *models.Account) {
	if currentUser == nil {
		fmt.Println("Anda belum login")
		return
	}

	// Ambil semua transaksi yang terkait dengan akun yang sedang login
	var transactions []models.Transaction
	result := db.Conn.Where("account_id = ?", currentUser.ID).Order("created_at desc").Find(&transactions)
	if result.Error != nil {
		fmt.Println("Gagal mengambil riwayat transaksi:", result.Error)
		return
	}

	// Tampilkan riwayat transaksi
	if len(transactions) == 0 {
		fmt.Println("Tidak ada riwayat transaksi.")
		return
	}

	fmt.Println("\n===== RIWAYAT TRANSAKSI =====")
	for _, transaction := range transactions {
		var transactionType string
		if transaction.Type == "deposit" {
			transactionType = "Setor Uang"
		} else if transaction.Type == "withdraw" {
			transactionType = "Tarik Tunai"
		} else if transaction.Type == "transfer_out" {
			transactionType = "Transfer Keluar"
		} else if transaction.Type == "transfer_in" {
			transactionType = "Transfer Masuk"
		}

		fmt.Printf("\nTanggal: %s\n", transaction.CreatedAt.Format("02-01-2006 15:04:05"))
		fmt.Printf("Tipe Transaksi: %s\n", transactionType)
		fmt.Printf("Jumlah: %.2f\n", transaction.Amount)
		if transaction.Type == "transfer_out" || transaction.Type == "transfer_in" {
			fmt.Printf("Target Akun: %d\n", *transaction.TargetID)
		}
	}
}
