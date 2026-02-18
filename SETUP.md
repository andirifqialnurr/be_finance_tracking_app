# Setup Instructions for Finance Tracking API

## Quick Start Guide

### 1. Install PostgreSQL

Download dan install PostgreSQL dari: https://www.postgresql.org/download/

Atau menggunakan Docker:
```bash
docker run --name finance-db -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15
```

### 2. Create Database

```bash
# Login ke PostgreSQL
psql -U postgres

# Buat database
CREATE DATABASE finance_tracking;

# Keluar
\q
```

### 3. Setup Environment Variables

Copy file `.env.example` ke `.env` dan sesuaikan dengan konfigurasi PostgreSQL Anda:

```bash
cp .env.example .env
```

Edit `.env` file, minimal update password:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=YOUR_POSTGRES_PASSWORD  # <-- Update ini
DB_NAME=finance_tracking
SERVER_PORT=8080
```

### 4. Install Go Dependencies

```bash
go mod download
```

### 5. Run Database Seeder (Optional)

Untuk populate initial categories sesuai kebutuhan Anda:

```bash
go run scripts/seed.go
```

Ini akan create 8 kategori default:
- Makan (Rp 1.240.000)
- Bensin (Rp 175.000)
- Listrik (Rp 200.000)
- GitHub Copilot (Rp 180.000)
- Kuota Internet (Rp 100.000)
- Netflix (Rp 120.000)
- Spotify (Rp 60.000)
- Fore Coffee (Rp 24.000)

**Total: Rp 2.099.000 / bulan**

### 6. Run Application

```bash
go run main.go
```

Server akan berjalan di: http://localhost:8080

### 7. Test API

```bash
# Health check
curl http://localhost:8080/health

# Get all categories
curl http://localhost:8080/api/v1/categories

# Create income
curl -X POST http://localhost:8080/api/v1/incomes \
  -H "Content-Type: application/json" \
  -d '{
    "source": "Gaji Februari",
    "amount": 8000000,
    "date": "2026-02-01T00:00:00Z"
  }'
```

## Troubleshooting

### Error: "database does not exist"
Pastikan database `finance_tracking` sudah dibuat.

### Error: "connection refused"
Pastikan PostgreSQL service sudah running.

**Windows:**
```bash
# Check service
Get-Service postgresql*

# Start service jika belum running
Start-Service postgresql-x64-15
```

**Mac/Linux:**
```bash
# Check status
sudo systemctl status postgresql

# Start service
sudo systemctl start postgresql
```

### Error: "password authentication failed"
Update DB_PASSWORD di `.env` file sesuai password PostgreSQL Anda.

## Development

### Hot Reload dengan Air

Install Air:
```bash
go install github.com/cosmtrek/air@latest
```

Buat file `.air.toml`:
```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ."
  bin = "tmp/main"
  include_ext = ["go"]
  exclude_dir = ["tmp", "vendor"]
```

Run with hot reload:
```bash
air
```

## Next Steps

✅ Backend API sudah ready
✅ Database schema sudah setup
✅ Auto-allocation logic sudah implement

Selanjutnya:
1. Test semua endpoints
2. Buat frontend (React/Vue/Flutter)
3. Deploy ke production

Happy coding! 🚀
