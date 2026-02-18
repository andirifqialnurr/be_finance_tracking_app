# Finance Tracking API

## Daftar Isi
1. [Pengenalan](#pengenalan)
2. [Struktur Aplikasi](#struktur-aplikasi)
3. [Teknologi yang Digunakan](#teknologi-yang-digunakan)
4. [Cara Instalasi](#cara-instalasi)
5. [Menjalankan Aplikasi](#menjalankan-aplikasi)
6. [API Endpoints](#api-endpoints)
7. [Contoh Penggunaan](#contoh-penggunaan)
8. [Penjelasan Kode](#penjelasan-kode)

---

## Pengenalan

Finance Tracking App adalah aplikasi backend untuk melacak transaksi keuangan pribadi. Aplikasi ini dibangun menggunakan **Golang** dengan arsitektur yang clean dan mudah dipahami.

### Fitur Utama:
- ✅ **Create** - Menambah transaksi baru (pemasukan/pengeluaran)
- ✅ **Read** - Melihat semua transaksi atau transaksi spesifik
- ✅ **Update** - Mengubah data transaksi yang sudah ada
- ✅ **Delete** - Menghapus transaksi
- ✅ **Statistics** - Melihat ringkasan keuangan (total pemasukan, pengeluaran, saldo)

---

## Struktur Aplikasi

```
be_finance_tracking_app/
│
├── main.go                          # Entry point aplikasi
├── go.mod                           # Go module dependencies
├── .gitignore                       # Git ignore file
│
├── database/
│   └── database.go                  # Koneksi dan setup database
│
├── models/
│   └── transaction.go               # Model data Transaction
│
└── handlers/
    └── transaction_handler.go       # Handler untuk CRUD operations
```

---

## Teknologi yang Digunakan

| Teknologi | Versi | Fungsi |
|-----------|-------|--------|
| **Go** | 1.21 | Bahasa pemrograman utama |
| **Gorilla Mux** | 1.8.1 | HTTP router dan URL matcher |
| **SQLite3** | 1.14.19 | Database (embedded database) |

### Mengapa SQLite?
- 📦 Tidak perlu instalasi database server terpisah
- 🚀 Cepat untuk development dan testing
- 💾 Database disimpan dalam satu file `.db`
- ✨ Cocok untuk aplikasi kecil hingga menengah

---

## Cara Instalasi

### 1. Prasyarat
Pastikan Anda telah menginstal:
- **Go** (versi 1.21 atau lebih tinggi)
- **Git** (optional, untuk version control)

### 2. Install Dependencies

Buka terminal/command prompt di folder project, kemudian jalankan:

```bash
go mod download
```

Perintah ini akan mengunduh semua dependencies yang diperlukan:
- `github.com/gorilla/mux` - untuk routing HTTP
- `github.com/mattn/go-sqlite3` - untuk database SQLite

### 3. Build Aplikasi (Optional)

Untuk membuat executable file:

```bash
go build -o finance-app.exe
```

---

## Menjalankan Aplikasi

### Cara 1: Menggunakan `go run`

```bash
go run main.go
```

### Cara 2: Menggunakan executable (jika sudah di-build)

```bash
./finance-app.exe
```

Jika berhasil, Anda akan melihat output:

```
Database connected successfully
Tables created successfully
Server running on http://localhost:8080
```

Server akan berjalan di **http://localhost:8080** 🚀

---

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

### 1. Get All Transactions
**Mendapatkan semua transaksi**

```
GET /api/v1/transactions
```

**Response Example:**
```json
{
  "success": true,
  "message": "Transactions fetched successfully",
  "data": [
    {
      "id": 1,
      "title": "Gaji Bulanan",
      "amount": 5000000,
      "type": "income",
      "category": "Salary",
      "description": "Gaji bulan Februari",
      "date": "2026-02-01",
      "created_at": "2026-02-01T10:00:00Z",
      "updated_at": "2026-02-01T10:00:00Z"
    }
  ]
}
```

---

### 2. Get Transaction by ID
**Mendapatkan satu transaksi berdasarkan ID**

```
GET /api/v1/transactions/{id}
```

**Parameter:**
- `id` (integer) - ID transaksi

**Response Example:**
```json
{
  "success": true,
  "message": "Transaction fetched successfully",
  "data": {
    "id": 1,
    "title": "Gaji Bulanan",
    "amount": 5000000,
    "type": "income",
    "category": "Salary",
    "description": "Gaji bulan Februari",
    "date": "2026-02-01",
    "created_at": "2026-02-01T10:00:00Z",
    "updated_at": "2026-02-01T10:00:00Z"
  }
}
```

---

### 3. Create Transaction
**Membuat transaksi baru**

```
POST /api/v1/transactions
```

**Request Body:**
```json
{
  "title": "Belanja Bulanan",
  "amount": 500000,
  "type": "expense",
  "category": "Groceries",
  "description": "Belanja kebutuhan rumah tangga",
  "date": "2026-02-10"
}
```

**Field Explanation:**
- `title` (string, required) - Judul transaksi
- `amount` (number, required) - Jumlah uang
- `type` (string, required) - Tipe transaksi: `"income"` atau `"expense"`
- `category` (string, required) - Kategori transaksi
- `description` (string, optional) - Deskripsi tambahan
- `date` (string, required) - Tanggal transaksi (format: YYYY-MM-DD)

**Response Example:**
```json
{
  "success": true,
  "message": "Transaction created successfully",
  "data": {
    "id": 2,
    "title": "Belanja Bulanan",
    "amount": 500000,
    "type": "expense",
    "category": "Groceries",
    "description": "Belanja kebutuhan rumah tangga",
    "date": "2026-02-10",
    "created_at": "2026-02-10T15:30:00Z",
    "updated_at": "2026-02-10T15:30:00Z"
  }
}
```

---

### 4. Update Transaction
**Mengupdate transaksi yang sudah ada**

```
PUT /api/v1/transactions/{id}
```

**Parameter:**
- `id` (integer) - ID transaksi yang akan diupdate

**Request Body:**
```json
{
  "title": "Belanja Bulanan (Revised)",
  "amount": 600000,
  "type": "expense",
  "category": "Groceries",
  "description": "Belanja kebutuhan rumah tangga + tambahan",
  "date": "2026-02-10"
}
```

**Response Example:**
```json
{
  "success": true,
  "message": "Transaction updated successfully",
  "data": {
    "id": 2,
    "title": "Belanja Bulanan (Revised)",
    "amount": 600000,
    "type": "expense",
    "category": "Groceries",
    "description": "Belanja kebutuhan rumah tangga + tambahan",
    "date": "2026-02-10",
    "created_at": "2026-02-10T15:30:00Z",
    "updated_at": "2026-02-11T09:00:00Z"
  }
}
```

---

### 5. Delete Transaction
**Menghapus transaksi**

```
DELETE /api/v1/transactions/{id}
```

**Parameter:**
- `id` (integer) - ID transaksi yang akan dihapus

**Response Example:**
```json
{
  "success": true,
  "message": "Transaction deleted successfully"
}
```

---

### 6. Get Statistics
**Mendapatkan statistik keuangan**

```
GET /api/v1/statistics
```

**Response Example:**
```json
{
  "success": true,
  "message": "Statistics fetched successfully",
  "data": {
    "total_income": 5000000,
    "total_expense": 600000,
    "balance": 4400000,
    "transactions": 2
  }
}
```

**Field Explanation:**
- `total_income` - Total semua pemasukan
- `total_expense` - Total semua pengeluaran
- `balance` - Saldo (pemasukan - pengeluaran)
- `transactions` - Jumlah total transaksi

---

## Contoh Penggunaan

### Menggunakan cURL (Command Line)

#### 1. Create Transaction (Pemasukan)
```bash
curl -X POST http://localhost:8080/api/v1/transactions ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Gaji Bulanan\",\"amount\":5000000,\"type\":\"income\",\"category\":\"Salary\",\"description\":\"Gaji Februari 2026\",\"date\":\"2026-02-01\"}"
```

#### 2. Create Transaction (Pengeluaran)
```bash
curl -X POST http://localhost:8080/api/v1/transactions ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Bayar Listrik\",\"amount\":300000,\"type\":\"expense\",\"category\":\"Utilities\",\"description\":\"Tagihan listrik\",\"date\":\"2026-02-05\"}"
```

#### 3. Get All Transactions
```bash
curl http://localhost:8080/api/v1/transactions
```

#### 4. Get Transaction by ID
```bash
curl http://localhost:8080/api/v1/transactions/1
```

#### 5. Update Transaction
```bash
curl -X PUT http://localhost:8080/api/v1/transactions/1 ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Gaji Bulanan (Updated)\",\"amount\":5500000,\"type\":\"income\",\"category\":\"Salary\",\"description\":\"Gaji + Bonus\",\"date\":\"2026-02-01\"}"
```

#### 6. Delete Transaction
```bash
curl -X DELETE http://localhost:8080/api/v1/transactions/1
```

#### 7. Get Statistics
```bash
curl http://localhost:8080/api/v1/statistics
```

---

### Menggunakan Postman atau Thunder Client

1. **Install Postman** atau **Thunder Client** (VS Code Extension)
2. Buat request baru dengan method dan URL sesuai endpoint di atas
3. Untuk POST/PUT, pilih **Body** → **raw** → **JSON**, kemudian masukkan data

---

## Penjelasan Kode

### 1. main.go
File ini adalah **entry point** aplikasi. Fungsi utamanya:
- Menginisialisasi koneksi database
- Membuat router menggunakan Gorilla Mux
- Mendefinisikan semua API routes
- Menjalankan HTTP server di port 8080

**Penjelasan Detail:**
```go
database.InitDB()           // Koneksi ke database dan buat tabel
defer database.CloseDB()    // Tutup koneksi saat aplikasi berhenti

router := mux.NewRouter()   // Buat router baru
api := router.PathPrefix("/api/v1").Subrouter()  // Prefix untuk semua endpoint

// Definisi routes
api.HandleFunc("/transactions", handlers.GetAllTransactions).Methods("GET")
// ... dst
```

---

### 2. database/database.go
File ini menangani:
- Koneksi ke database SQLite
- Pembuatan tabel `transactions`
- Menutup koneksi database

**Schema Tabel Transactions:**
```sql
CREATE TABLE transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,     -- ID unik (auto increment)
    title TEXT NOT NULL,                       -- Judul transaksi
    amount REAL NOT NULL,                      -- Jumlah uang
    type TEXT NOT NULL,                        -- 'income' atau 'expense'
    category TEXT NOT NULL,                    -- Kategori
    description TEXT,                          -- Deskripsi (optional)
    date TEXT NOT NULL,                        -- Tanggal transaksi
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Waktu dibuat
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP   -- Waktu diupdate
);
```

**Constraint:**
- `type` hanya bisa `"income"` atau `"expense"`
- `CHECK(type IN ('income', 'expense'))`

---

### 3. models/transaction.go
File ini mendefinisikan **struktur data** yang digunakan:

#### Transaction Model
```go
type Transaction struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Amount      float64   `json:"amount"`
    Type        string    `json:"type"`        // "income" atau "expense"
    Category    string    `json:"category"`
    Description string    `json:"description"`
    Date        string    `json:"date"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**Tag `json`** digunakan untuk:
- Mengubah struct Go ke JSON (marshaling)
- Mengubah JSON ke struct Go (unmarshaling)

#### Statistics Model
```go
type Statistics struct {
    TotalIncome  float64 `json:"total_income"`
    TotalExpense float64 `json:"total_expense"`
    Balance      float64 `json:"balance"`
    Transactions int     `json:"transactions"`
}
```

#### Response Model
```go
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`  // omitempty: tidak tampil jika nil
}
```

---

### 4. handlers/transaction_handler.go
File ini berisi **business logic** untuk CRUD operations:

#### GetAllTransactions
- Query semua data dari database
- Sort berdasarkan tanggal terbaru
- Return array of transactions

```go
query := `SELECT ... FROM transactions ORDER BY date DESC, created_at DESC`
rows, err := database.DB.Query(query)
// Loop through rows dan scan ke struct
```

#### GetTransactionByID
- Cari transaksi berdasarkan ID
- Return single transaction
- Return 404 jika tidak ditemukan

```go
err = database.DB.QueryRow(query, id).Scan(...)
if err == sql.ErrNoRows {
    respondWithError(w, http.StatusNotFound, "Transaction not found")
}
```

#### CreateTransaction
- Validasi input (required fields)
- Validasi type harus "income" atau "expense"
- Insert ke database
- Return data yang baru dibuat

```go
if t.Type != "income" && t.Type != "expense" {
    respondWithError(w, http.StatusBadRequest, "Type must be 'income' or 'expense'")
}
```

#### UpdateTransaction
- Validasi ID dan input
- Update data di database
- Return 404 jika transaksi tidak ditemukan
- Update `updated_at` otomatis

```go
query := `UPDATE transactions SET ... WHERE id = ?`
result, err := database.DB.Exec(query, ...)
rowsAffected, _ := result.RowsAffected()
```

#### DeleteTransaction
- Hapus transaksi berdasarkan ID
- Return 404 jika tidak ditemukan

```go
result, err := database.DB.Exec(query, id)
rowsAffected, _ := result.RowsAffected()
```

#### GetStatistics
- Hitung total income dengan `SUM(amount) WHERE type = 'income'`
- Hitung total expense dengan `SUM(amount) WHERE type = 'expense'`
- Hitung balance: income - expense
- Hitung jumlah transaksi dengan `COUNT(*)`

```go
stats.Balance = stats.TotalIncome - stats.TotalExpense
```

---

## Error Handling

Aplikasi memiliki error handling yang baik:

### Status Codes yang Digunakan:
- **200 OK** - Request berhasil (GET, PUT, DELETE)
- **201 Created** - Data berhasil dibuat (POST)
- **400 Bad Request** - Input tidak valid
- **404 Not Found** - Data tidak ditemukan
- **500 Internal Server Error** - Error di server

### Format Error Response:
```json
{
  "success": false,
  "message": "Error message here"
}
```

### Format Success Response:
```json
{
  "success": true,
  "message": "Success message here",
  "data": { ... }
}
```

---

## Tips Pengembangan Lebih Lanjut

### 1. Tambah Fitur Filter
Tambahkan filter berdasarkan:
- Tanggal (range)
- Kategori
- Tipe (income/expense)

### 2. Pagination
Untuk data yang banyak, tambahkan pagination:
```go
api.HandleFunc("/transactions?page=1&limit=10", handlers.GetAllTransactions)
```

### 3. Authentication
Tambahkan JWT authentication untuk keamanan:
```go
import "github.com/dgrijalva/jwt-go"
```

### 4. Export Data
Tambahkan endpoint untuk export ke CSV/Excel:
```go
api.HandleFunc("/export/csv", handlers.ExportCSV).Methods("GET")
```

### 5. Dashboard
Tambahkan endpoint untuk data dashboard:
- Grafik income vs expense per bulan
- Top categories
- Recent transactions

### 6. Kategori Custom
Buat tabel `categories` terpisah untuk manage kategori:
```sql
CREATE TABLE categories (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL
);
```

### 7. Multi-User
Tambahkan tabel `users` dan relasikan dengan transactions:
```sql
ALTER TABLE transactions ADD COLUMN user_id INTEGER;
```

---

## Troubleshooting

### Problem: "go: cannot find main module"
**Solution:** Pastikan Anda di folder yang benar dan file `go.mod` ada.

### Problem: "cannot find package github.com/gorilla/mux"
**Solution:** Run `go mod download` atau `go get github.com/gorilla/mux`

### Problem: "database locked"
**Solution:** Pastikan hanya satu instance aplikasi yang berjalan.

### Problem: Port 8080 sudah digunakan
**Solution:** Ubah port di `main.go`:
```go
port := ":8081"  // Ganti dengan port lain
```

---

## Testing API

### Test dengan cURL (Windows PowerShell):
Ganti `^` dengan `` ` `` (backtick) di PowerShell:

```powershell
curl -X POST http://localhost:8080/api/v1/transactions `
  -H "Content-Type: application/json" `
  -d '{"title":"Test","amount":1000,"type":"income","category":"Test","date":"2026-02-11"}'
```

### Test dengan Browser:
Untuk GET request, bisa langsung buka di browser:
- http://localhost:8080/api/v1/transactions
- http://localhost:8080/api/v1/statistics

---

## Kesimpulan

Anda telah berhasil membuat **Finance Tracking App** dengan fitur CRUD lengkap menggunakan Golang! 🎉

### Yang Telah Dipelajari:
✅ Struktur project Go yang clean  
✅ Setup database dengan SQLite  
✅ Membuat REST API dengan Gorilla Mux  
✅ CRUD operations (Create, Read, Update, Delete)  
✅ Error handling dan response formatting  
✅ Database queries dengan Go  

### Next Steps:
1. Coba semua endpoint dengan cURL atau Postman
2. Tambahkan fitur-fitur baru sesuai kebutuhan
3. Deploy aplikasi ke server (heroku, railway, dll)
4. Buat frontend untuk aplikasi ini (React, Vue, atau Angular)

---

## Kontak & Support

Jika ada pertanyaan atau butuh bantuan, jangan ragu untuk bertanya! 😊

**Happy Coding!** 🚀

---

*Dokumentasi ini dibuat pada: 11 Februari 2026*
