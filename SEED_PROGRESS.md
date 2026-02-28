# Seed Progress

Status terakhir seeding data demo untuk pengujian semua endpoint Swagger.

---

## Checklist Steps

| Step | Status | Detail |
|------|--------|--------|
| 1. Expense Categories | ✅ Done | 8 kategori: Makan, Bensin, Listrik, GitHub Copilot, Kuota Internet, Netflix, Spotify, Fore Coffee |
| 2. Budget Alerts | ✅ Done | 1 alert per kategori (threshold 75–90%) |
| 3. Accounts | ✅ Done | BCA - Gajian (CARD/SALARY), Dompet Cash (CASH), Tabungan Darurat (SAVINGS + goal Rp25jt) |
| 4. Oktober 2025 | ✅ Done | TopUp BCA Rp8jt, subscriptions, makan/bensin/listrik dari BCA |
| 5. November 2025 | ✅ Done | TopUp BCA Rp8jt + freelance Rp1.5jt, subscriptions, expenses BCA |
| 6. Desember 2025 | ✅ Done | TopUp Rp8jt + bonus Rp3jt, transfer BCA→Cash Rp2jt + BCA→Sav Rp2jt, expenses Cash & BCA |
| 7. Januari 2026 | ✅ Done | TopUp Rp8jt, transfer BCA→Cash Rp2jt + BCA→Sav Rp1.5jt, expenses Cash & BCA |
| 8. Februari 2026 | ✅ Done | TopUp Rp8jt + freelance Rp3.5jt, transfers, expenses, **realokasi Listrik→Makan Rp15k** |
| 9. Register Device | ✅ Done | `POST /notifications/register-device` — OneSignal player id `demo-player-id-001` |
| 10. Notification Settings | ✅ Done | 4 settings: ALLOCATION_REMINDER, BUDGET_ALERT, SAVINGS_GOAL, SCHEDULED_FUND |
| 11. Scheduled Funds | ✅ Done | 2 schedule: auto-transfer ke Cash (tgl 1) dan Tabungan (tgl 2) |

---

## Endpoint Coverage Matrix

### Accounts
| Endpoint | Dicakup Oleh | Data |
|----------|-------------|------|
| `GET /accounts` | Step 3 | 3 akun |
| `POST /accounts` | Step 3 | BCA, Cash, Savings |
| `GET /accounts/:id` | Step 3 | detail tiap akun |
| `PUT /accounts/:id` | — | belum diseed (manual via Swagger) |
| `DELETE /accounts/:id` | — | tidak relevan untuk demo |
| `POST /accounts/:id/topup` | Step 4–8 | 7 topup (gaji + freelance + bonus) |
| `POST /accounts/:id/spent` | Step 4–8 | 50+ pengeluaran tersebar 5 bulan |
| `GET /accounts/summary` | Step 3–8 | tersedia setelah ada akun + saldo |
| `POST /accounts/:id/transfer` | — | manual (transfer sudah via `POST /transfers`) |
| `GET /accounts/:id/transfers` | Step 6–8 | BCA memiliki 6 transfer keluar |

### Categories
| Endpoint | Dicakup Oleh | Data |
|----------|-------------|------|
| `GET /categories` | Step 1 | 8 kategori |
| `POST /categories` | Step 1 | dibuat semua |
| `GET /categories/:id` | Step 1 | tiap kategori |
| `PUT /categories/:id` | — | manual |
| `DELETE /categories/:id` | — | tidak relevan |

### Budgets
| Endpoint | Dicakup Oleh | Data |
|----------|-------------|------|
| `GET /budgets` | Step 4–8 | budget dibuat otomatis saat topup SALARY |
| `GET /budgets/:id` | Step 4–8 | ada per bulan per kategori |
| `PUT /budgets/:id` | — | manual |
| `POST /budgets/reallocate` | Step 8 | 1 realokasi: Listrik→Makan Feb 2026 |
| `GET /budgets/reallocations` | Step 8 | harus ada 1 entri |

### Alerts
| Endpoint | Dicakup Oleh | Data |
|----------|-------------|------|
| `GET /alerts` | Step 2 | 8 alert |
| `POST /alerts` | Step 2 | 1 per kategori |
| `GET /alerts/:id` | Step 2 | detail |
| `PUT /alerts/:id` | — | manual |
| `DELETE /alerts/:id` | — | tidak relevan |

### Transfers
| Endpoint | Dicakup Oleh | Data |
|----------|-------------|------|
| `GET /transfers` | Step 6–8 | 6 transfer tersebar Des–Feb |
| `POST /transfers` | Step 6–8 | BCA→Cash dan BCA→Tabungan |
| `GET /transfers/:id` | Step 6–8 | detail tiap transfer |
| `DELETE /transfers/:id` | — | tidak relevan |

### Notifications
| Endpoint | Dicakup Oleh | Data |
|----------|-------------|------|
| `GET /notifications` | Step 10 | 4 setting notifikasi |
| `POST /notifications/settings` | Step 10 | ALLOCATION_REMINDER, BUDGET_ALERT, SAVINGS_GOAL, SCHEDULED_FUND |
| `PATCH /notifications/settings/:id` | — | manual |
| `DELETE /notifications/settings/:id` | — | tidak relevan |
| `POST /notifications/register-device` | Step 9 | player id `demo-player-id-001` |
| `POST /notifications/test` | — | manual via Swagger |

### Scheduled Funds
| Endpoint | Dicakup Oleh | Data |
|----------|-------------|------|
| `GET /scheduled-funds` | Step 11 | 2 scheduled fund |
| `POST /scheduled-funds` | Step 11 | auto-transfer Cash & Tabungan |
| `GET /scheduled-funds/:id` | Step 11 | detail |
| `PUT /scheduled-funds/:id` | — | manual |
| `DELETE /scheduled-funds/:id` | — | tidak relevan |

### Analytics & Reports
| Endpoint | Dicakup Oleh | Data |
|----------|-------------|------|
| `GET /analytics/overview` | Step 4–8 | 5 bulan data (Okt 2025–Feb 2026) |
| `GET /analytics/spending-trends` | Step 4–8 | trend pengeluaran per kategori |
| `GET /analytics/income-vs-expense` | Step 4–8 | perbandingan bulanan |
| `GET /analytics/category-breakdown` | Step 4–8 | breakdown per kategori |
| `GET /reports/monthly` | Step 4–8 | laporan bulanan per bulan |
| `GET /reports/export` | Step 4–8 | export PDF/Excel tersedia |
| `GET /statistics/overview` | Step 4–8 | ringkasan statistik |
| `GET /statistics/monthly` | Step 4–8 | statistik bulanan |

### Transactions (Read-Only)
| Endpoint | Dicakup Oleh | Data |
|----------|-------------|------|
| `GET /transactions` | Step 4–8 | semua transaksi gabungan |
| `GET /transactions/:id` | Step 4–8 | detail |

---

## Estimasi Saldo Akhir (Februari 2026)

| Akun | Estimasi Saldo |
|------|---------------|
| BCA - Gajian | ≈ Rp 32.155.000 |
| Dompet Cash | ≈ Rp 4.224.000 |
| Tabungan Darurat | ≈ Rp 5.500.000 |
| **Total** | **≈ Rp 41.879.000** |

---

## Bug yang Diperbaiki dalam Seed Rewrite

| # | Bug Lama | Fix |
|---|----------|-----|
| 1 | Menggunakan `POST /incomes` langsung | Diganti `POST /accounts/:id/topup` → saldo ter-update + auto-alokasi SALARY |
| 2 | Menggunakan `POST /expenses` langsung | Diganti `POST /accounts/:id/spent` → saldo berkurang + budget kategori ter-update |
| 3 | Tipe notifikasi `EXPENSE_REMINDER` & `SAVINGS_REMINDER` invalid | Diganti dengan `ALLOCATION_REMINDER`, `BUDGET_ALERT`, `SAVINGS_GOAL`, `SCHEDULED_FUND` |
| 4 | Tidak ada budget alert | Ditambah 8 alert (1 per kategori) di Step 2 |
| 5 | Tidak ada budget reallocation | Ditambah 1 realokasi Listrik→Makan Rp15k di Step 8 |
| 6 | Tidak ada `POST /notifications/register-device` | Ditambah di Step 9 |

---

## Cara Menjalankan Seed

```bash
# Pastikan server berjalan terlebih dahulu
go run main.go

# Jalankan seed (terminal baru)
go run scripts/seed/main.go
```

> **Catatan:** Jalankan cleanup script dulu bila database sudah berisi data lama:
> ```bash
> go run scripts/cleanup/main.go
> ```
