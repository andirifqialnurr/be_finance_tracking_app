# Finance Tracking App - Arsitektur & Flow

## 📋 Deskripsi Umum
Aplikasi finance tracking pribadi untuk membantu tracking pemasukan dan pengeluaran dengan sistem alokasi otomatis dan kategori yang dapat disesuaikan.

---

## 🎯 Fitur Utama

### 1. Income Management (Pemasukan)
- Catat pemasukan dari berbagai sumber (gaji, project, dll)
- Auto-allocate pemasukan ke berbagai kategori pengeluaran
- History pemasukan dengan detail sumber

### 2. Expense Management (Pengeluaran)
- **Tipe Pengeluaran:**
  - **SUBSCRIPTION**: Pembayaran bulanan tetap (Netflix, Spotify, dll)
  - **DAILY_CONTINUOUS**: Pengeluaran harian berulang (Makan)
  - **USAGE_BASED**: Berbasis penggunaan dengan budget tetap (Bensin, Listrik)
  - **ONE_TIME**: Pengeluaran sekali bayar

- **Budget Allocation**: Setiap kategori punya budget yang terisi otomatis saat ada pemasukan
- **Budget Tracking**: Track sisa budget per kategori
- **Expense Recording**: Catat setiap pengeluaran yang mengurangi budget kategori

### 3. Budget Auto-Allocation
- Ketika ada pemasukan, otomatis dialokasikan ke semua kategori aktif
- Alokasi berdasarkan persentase atau jumlah tetap per kategori
- Real-time update sisa budget setelah expense dicatat

---

## 🏗️ Arsitektur Sistem

```
┌─────────────────────────────────────────────────────────┐
│                     Mobile/Web App                      │
│                    (React/Flutter)                      │
└─────────────────────────────────────────────────────────┘
                            │
                            │ HTTP/REST API
                            ▼
┌─────────────────────────────────────────────────────────┐
│                   Backend API (Go)                      │
│  ┌──────────────────────────────────────────────────┐  │
│  │              API Handlers Layer                  │  │
│  │  - Income Handler                                │  │
│  │  - Expense Handler                               │  │
│  │  - Category Handler                              │  │
│  │  - Budget Handler                                │  │
│  └──────────────────────────────────────────────────┘  │
│                           │                             │
│  ┌──────────────────────────────────────────────────┐  │
│  │           Business Logic Layer                   │  │
│  │  - Budget Allocation Service                     │  │
│  │  - Category Management Service                   │  │
│  │  - Transaction Service                           │  │
│  └──────────────────────────────────────────────────┘  │
│                           │                             │
│  ┌──────────────────────────────────────────────────┐  │
│  │              Data Access Layer                   │  │
│  │  - Repository Pattern                            │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                PostgreSQL Database                      │
│  - incomes                                              │
│  - expense_categories                                   │
│  - category_budgets                                     │
│  - expenses                                             │
│  - budget_allocations                                   │
└─────────────────────────────────────────────────────────┘
```

---

## 📊 Database Schema

### Table: `incomes`
Menyimpan semua pemasukan
```sql
id              UUID PRIMARY KEY
source          VARCHAR(100)      -- "Gaji", "Project X", "Freelance", dll
amount          DECIMAL(15,2)     -- Jumlah pemasukan
date            TIMESTAMP         -- Tanggal pemasukan
description     TEXT              -- Deskripsi opsional
created_at      TIMESTAMP
updated_at      TIMESTAMP
```

### Table: `expense_categories`
Master kategori pengeluaran (user-defined)
```sql
id              UUID PRIMARY KEY
name            VARCHAR(100)      -- "Makan", "Netflix", "Bensin", dll
type            VARCHAR(50)       -- "SUBSCRIPTION", "DAILY_CONTINUOUS", "USAGE_BASED", "ONE_TIME"
monthly_budget  DECIMAL(15,2)     -- Budget bulanan untuk kategori ini
allocation_priority INT           -- Urutan prioritas alokasi (1 = tertinggi)
is_active       BOOLEAN           -- Kategori aktif atau tidak
metadata        JSONB             -- Extra config: {daily_amount: 40000, refill_count: 5, etc}
created_at      TIMESTAMP
updated_at      TIMESTAMP
```

### Table: `category_budgets`
Budget real-time per kategori (reset setiap bulan)
```sql
id              UUID PRIMARY KEY
category_id     UUID REFERENCES expense_categories(id)
month           INT               -- 1-12
year            INT               -- 2026, 2027, etc
allocated_amount DECIMAL(15,2)    -- Total budget yang sudah dialokasikan
spent_amount    DECIMAL(15,2)     -- Total yang sudah dipakai
remaining_amount DECIMAL(15,2)    -- Sisa budget (allocated - spent)
created_at      TIMESTAMP
updated_at      TIMESTAMP

UNIQUE(category_id, month, year)
```

### Table: `expenses`
Detail setiap pengeluaran
```sql
id              UUID PRIMARY KEY
category_id     UUID REFERENCES expense_categories(id)
amount          DECIMAL(15,2)     -- Jumlah pengeluaran
date            TIMESTAMP         -- Tanggal pengeluaran
description     TEXT              -- Deskripsi (misal: "Isi bensin Shell", "Makan siang Warteg")
created_at      TIMESTAMP
updated_at      TIMESTAMP
```

### Table: `budget_allocations`
History alokasi budget dari income ke kategori
```sql
id              UUID PRIMARY KEY
income_id       UUID REFERENCES incomes(id)
category_id     UUID REFERENCES expense_categories(id)
allocated_amount DECIMAL(15,2)    -- Jumlah yang dialokasikan
month           INT
year            INT
created_at      TIMESTAMP
```

---

## 🔄 Flow Diagram

### Flow 1: Pemasukan & Auto-Allocation
```
┌─────────────────┐
│  User Input     │
│  Income Data    │
│  - Source       │
│  - Amount       │
│  - Date         │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────┐
│  Simpan ke table `incomes`      │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│  Get semua active expense_categories    │
│  (is_active = true)                     │
└────────┬────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│  Hitung alokasi per kategori            │
│  Berdasarkan monthly_budget ratio       │
│                                          │
│  Formula:                                │
│  allocation = (category.monthly_budget / │
│                total_monthly_budgets) *  │
│                income.amount             │
└────────┬────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│  Loop setiap kategori:                  │
│  1. Update/Create category_budgets      │
│     - allocated_amount += allocation    │
│     - remaining_amount += allocation    │
│  2. Insert ke budget_allocations        │
└────────┬────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│  Return success response        │
│  + Allocation breakdown         │
└─────────────────────────────────┘
```

### Flow 2: Catat Pengeluaran
```
┌─────────────────┐
│  User Input     │
│  Expense Data   │
│  - Category     │
│  - Amount       │
│  - Date         │
│  - Description  │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────┐
│  Validasi:                      │
│  - Cek category exists & active │
│  - Cek budget mencukupi         │
│  (remaining_amount >= amount)   │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│  Simpan ke table `expenses`     │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│  Update category_budgets:       │
│  - spent_amount += amount       │
│  - remaining_amount -= amount   │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│  Return success response        │
│  + Updated budget info          │
└─────────────────────────────────┘
```

### Flow 3: Tambah/Edit Kategori
```
┌─────────────────┐
│  User Input     │
│  Category Data  │
│  - Name         │
│  - Type         │
│  - Budget       │
│  - Priority     │
│  - Metadata     │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────┐
│  Validasi input                 │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│  Simpan ke expense_categories   │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│  Return created category        │
└─────────────────────────────────┘
```

---

## 🔌 API Endpoints

### Income APIs
```
POST   /api/v1/incomes
  Body: {
    "source": "Gaji Januari",
    "amount": 8000000,
    "date": "2026-01-25T00:00:00Z",
    "description": "Gaji bulanan"
  }
  Response: {
    "income": {...},
    "allocations": [
      {"category": "Makan", "allocated": 1240000},
      {"category": "Netflix", "allocated": 120000},
      ...
    ]
  }

GET    /api/v1/incomes?month=1&year=2026
  Response: [{...}, {...}]

GET    /api/v1/incomes/:id
  Response: {...}
```

### Expense Category APIs
```
POST   /api/v1/categories
  Body: {
    "name": "Netflix",
    "type": "SUBSCRIPTION",
    "monthly_budget": 120000,
    "allocation_priority": 1,
    "metadata": {}
  }

GET    /api/v1/categories
  Response: [{...}, {...}]

GET    /api/v1/categories/:id
  Response: {...}

PUT    /api/v1/categories/:id
  Body: {
    "name": "Netflix Premium",
    "monthly_budget": 150000
  }

DELETE /api/v1/categories/:id
  (Soft delete: set is_active = false)
```

### Expense APIs
```
POST   /api/v1/expenses
  Body: {
    "category_id": "uuid",
    "amount": 35000,
    "date": "2026-02-17T12:00:00Z",
    "description": "Isi bensin Shell"
  }

GET    /api/v1/expenses?category_id=uuid&month=2&year=2026
  Response: [{...}, {...}]

GET    /api/v1/expenses/:id
  Response: {...}
```

### Budget APIs
```
GET    /api/v1/budgets?month=2&year=2026
  Response: [
    {
      "category": {...},
      "allocated_amount": 1240000,
      "spent_amount": 480000,
      "remaining_amount": 760000,
      "percentage_used": 38.7
    },
    ...
  ]

GET    /api/v1/budgets/summary?month=2&year=2026
  Response: {
    "total_income": 8000000,
    "total_allocated": 2200000,
    "total_spent": 1850000,
    "total_remaining": 350000,
    "unallocated": 5800000
  }
```

---

## 💡 Contoh Use Case

### Scenario 1: Dapat Gaji
1. User input pemasukan: Rp 8.000.000 dari "Gaji Februari"
2. System otomatis alokasikan:
   - Makan: Rp 1.240.000
   - Bensin: Rp 175.000
   - Listrik: Rp 200.000
   - Copilot: Rp 180.000
   - Internet: Rp 100.000
   - Netflix: Rp 120.000
   - Spotify: Rp 60.000
   - Fore: Rp 24.000
   - **Total alokasi: Rp 2.099.000**
   - **Sisa uang bebas: Rp 5.901.000**

3. User bisa lihat budget per kategori di dashboard

### Scenario 2: Isi Bensin Pertama Kali
1. User buka app, pilih "Catat Pengeluaran"
2. Pilih kategori "Bensin"
3. Input amount: Rp 35.000
4. Deskripsi: "Isi bensin Shell Sudirman"
5. Submit
6. System:
   - Budget bensin: Rp 175.000 → Rp 140.000 (sisa 4 kali isi)
   - Catat expense baru

### Scenario 3: Makan Sehari
1. User catat pengeluaran makan: Rp 40.000
2. System:
   - Budget makan: Rp 1.240.000 → Rp 1.200.000 (sisa 30 hari)

### Scenario 4: Tambah Kategori Baru
1. User mau track "Kopi" dengan budget Rp 300.000/bulan
2. Buat kategori baru:
   - Name: "Kopi"
   - Type: DAILY_CONTINUOUS
   - Monthly Budget: Rp 300.000
   - Priority: 5
3. Kategori tersimpan, dan pemasukan berikutnya akan include alokasi ke "Kopi"

---

## 🎨 Dashboard/UI Features

### 1. Dashboard Home
- **Total Income bulan ini**
- **Total Allocated bulan ini**
- **Total Spent bulan ini**
- **Remaining Budget**
- **Chart**: Income vs Expense trend (monthly)
- **Quick Action**: Tambah Income, Tambah Expense

### 2. Budget Overview
- List semua kategori dengan:
  - Progress bar (spent/allocated)
  - Remaining amount
  - Percentage used
  - Color coding:
    - Green: < 50% used
    - Yellow: 50-80% used
    - Red: > 80% used

### 3. Expense History
- Filter by:
  - Category
  - Date range
  - Month/Year
- Sorting: Date, Amount
- Detail per expense

### 4. Income History
- List semua pemasukan
- Lihat breakdown alokasi per income

### 5. Category Management
- CRUD kategori
- Set priority
- Set monthly budget
- Enable/disable kategori

---

## 🔐 Future Enhancements

1. **Multi-Currency Support**
2. **Recurring Transactions**: Auto-create expense untuk subscription
3. **Budget Alerts**: Notifikasi ketika budget hampir habis
4. **Export Reports**: PDF/Excel
5. **Savings Goals**: Set target tabungan
6. **Analytics**: Spending patterns, trends
7. **Multi-User**: Family budget tracking
8. **Bank Integration**: Auto-sync dengan rekening bank

---

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin / Echo / Chi
- **Database**: PostgreSQL 15+
- **ORM**: GORM
- **API**: RESTful
- **Auth**: JWT (untuk future multi-user)

### Frontend (Future)
- **Web**: React + TypeScript + Tailwind CSS
- **Mobile**: Flutter / React Native

### DevOps
- **Containerization**: Docker
- **Database Migration**: golang-migrate
- **API Documentation**: Swagger/OpenAPI

---

## 📝 Development Phases

### Phase 1: Core Backend (1-2 minggu)
- ✅ Setup project structure
- ✅ Database schema & migrations
- ✅ Models & repositories
- ✅ Core APIs: Income, Category, Expense, Budget
- ✅ Auto-allocation logic
- ✅ Unit tests

### Phase 2: Additional Features (1 minggu)
- [ ] Advanced filtering & search
- [ ] Budget analytics endpoints
- [ ] Data validation & error handling
- [ ] API documentation

### Phase 3: Frontend Development (2-3 minggu)
- [ ] Web/Mobile UI
- [ ] Dashboard
- [ ] Forms & CRUD operations
- [ ] Charts & visualizations

### Phase 4: Polish & Deploy (1 minggu)
- [ ] Testing (E2E)
- [ ] Performance optimization
- [ ] Deploy to production
- [ ] User documentation

---

## ✅ Kesimpulan

Arsitektur ini dirancang untuk:
1. **Flexible**: User bisa custom kategori sendiri
2. **Automated**: Auto-allocation dari income
3. **Real-time tracking**: Budget update langsung
4. **Scalable**: Bisa dikembangkan ke multi-user atau fitur advanced
5. **Simple**: Focus pada core functionality dulu

Apakah arsitektur dan flow ini sudah sesuai dengan kebutuhan Anda? Ada yang perlu diubah atau ditambahkan?
