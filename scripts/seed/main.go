package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const baseURL = "http://localhost:8081/api/v1"

// M is a shorthand for JSON object maps
type M map[string]interface{}

type APIResp struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error"`
}

var httpClient = &http.Client{Timeout: 15 * time.Second}

// -- HTTP helpers --------------------------------------------------------------

func apiCall(method, path string, body M) (json.RawMessage, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var ar APIResp
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	if !ar.Success {
		return nil, fmt.Errorf("API [%s %s]: %s", method, path, ar.Error)
	}
	return ar.Data, nil
}

func post(path string, body M) (json.RawMessage, error) { return apiCall("POST", path, body) }

// idOf extracts "id" field from a raw JSON object
func idOf(data json.RawMessage) string {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return ""
	}
	if id, ok := obj["id"]; ok {
		return fmt.Sprintf("%v", id)
	}
	return ""
}

// ts formats a UTC timestamp as RFC3339
func ts(year, month, day, hour, min int) string {
	return time.Date(year, time.Month(month), day, hour, min, 0, 0, time.UTC).Format(time.RFC3339)
}

func rp(amount float64) string { return fmt.Sprintf("Rp %.0f", amount) }

func sep(label string) {
	log.Println("\n" + strings.Repeat("─", 60))
	log.Printf("  %s", label)
	log.Println(strings.Repeat("─", 60))
}

// ─── Entry point ───────────────────────────────────────────────────────────────

func main() {
	// Health check
	resp, err := httpClient.Get("http://localhost:8081/health")
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("❌  Server tidak berjalan di localhost:8081. Jalankan server terlebih dahulu!")
	}
	resp.Body.Close()
	log.Println("✅  Server sehat. Memulai seeding...")

	// ─── STEP 1: Categories ────────────────────────────────────────────────────
	sep("STEP 1 — Expense Categories")

	dailyMakan := 40_000.0
	catDefs := []M{
		{"name": "Makan", "type": "DAILY_CONTINUOUS", "daily_amount": dailyMakan, "allocation_priority": 1, "metadata": M{"icon": "food", "note": "Rp40.000/hari"}},
		{"name": "Bensin", "type": "USAGE_BASED", "monthly_budget": 175_000.0, "allocation_priority": 2, "metadata": M{"icon": "fuel"}},
		{"name": "Listrik", "type": "USAGE_BASED", "monthly_budget": 200_000.0, "allocation_priority": 3, "metadata": M{"icon": "electric", "type": "PLN token"}},
		{"name": "GitHub Copilot", "type": "SUBSCRIPTION", "monthly_budget": 180_000.0, "allocation_priority": 4, "metadata": M{"renewal": "monthly", "vendor": "GitHub"}},
		{"name": "Kuota Internet", "type": "SUBSCRIPTION", "monthly_budget": 100_000.0, "allocation_priority": 5, "metadata": M{"provider": "Telkomsel", "quota": "50GB"}},
		{"name": "Netflix", "type": "SUBSCRIPTION", "monthly_budget": 120_000.0, "allocation_priority": 6, "metadata": M{"plan": "Premium"}},
		{"name": "Spotify", "type": "SUBSCRIPTION", "monthly_budget": 60_000.0, "allocation_priority": 7, "metadata": M{"plan": "Individual"}},
		{"name": "Fore Coffee", "type": "SUBSCRIPTION", "monthly_budget": 24_000.0, "allocation_priority": 8, "metadata": M{"membership": "Fore+"}},
	}

	catIDs := make(map[string]string) // name → uuid
	for _, def := range catDefs {
		data, err := post("/categories", def)
		if err != nil {
			log.Fatalf("  ❌ Category [%v]: %v", def["name"], err)
		}
		catIDs[def["name"].(string)] = idOf(data)
		log.Printf("  ✅ %-20s → %s", def["name"], catIDs[def["name"].(string)])
	}

	// ─── STEP 2: Budget Alerts (one per category) ───────────────────────────────
	sep("STEP 2 — Budget Alerts")

	alertThresholds := map[string]int{
		"Makan": 75, "Bensin": 75,
		"GitHub Copilot": 90, "Netflix": 90, "Spotify": 90, "Fore Coffee": 90,
	}
	for name, catID := range catIDs {
		threshold := 80
		if t, ok := alertThresholds[name]; ok {
			threshold = t
		}
		data, err := post("/alerts", M{
			"category_id": catID, "threshold_percentage": threshold, "is_enabled": true,
		})
		if err != nil {
			log.Printf("  ⚠️  Alert [%s]: %v", name, err)
		} else {
			log.Printf("  ✅ Alert %-20s threshold=%d%% → %s", name, threshold, idOf(data))
		}
	}

	// ─── STEP 3: Accounts ─────────────────────────────────────────────────────
	sep("STEP 3 — Accounts")

	salaryType := "SALARY"
	goalAmt := 25_000_000.0
	goalLabel := "Dana darurat 3 bulan gaji"

	bcaData, err := post("/accounts", M{
		"name": "BCA - Gajian", "type": "CARD", "income_type": salaryType,
		"color": "#1565C0", "description": "Rekening utama penerimaan gaji & auto-alokasi budget",
		"initial_balance": 0,
	})
	if err != nil {
		log.Fatalf("  ❌ Account BCA: %v", err)
	}
	bcaID := idOf(bcaData)
	log.Printf("  ✅ %-25s → %s", "BCA - Gajian", bcaID)

	cashData, err := post("/accounts", M{
		"name": "Dompet Cash", "type": "CASH",
		"color": "#2E7D32", "description": "Uang tunai untuk kebutuhan harian (makan & bensin)",
		"initial_balance": 0,
	})
	if err != nil {
		log.Fatalf("  ❌ Account Cash: %v", err)
	}
	cashID := idOf(cashData)
	log.Printf("  ✅ %-25s → %s", "Dompet Cash", cashID)

	savingsData, err := post("/accounts", M{
		"name": "Tabungan Darurat", "type": "SAVINGS",
		"color": "#E65100", "description": "Dana darurat untuk keperluan mendesak",
		"goal_amount": goalAmt, "goal_label": goalLabel, "initial_balance": 0,
	})
	if err != nil {
		log.Fatalf("  ❌ Account Savings: %v", err)
	}
	savingsID := idOf(savingsData)
	log.Printf("  ✅ %-25s → %s", "Tabungan Darurat", savingsID)

	// ─── Helpers ────────────────────────────────────────────────────────────────
	// doTopUp  → POST /accounts/:id/topup  (updates balance, triggers auto-alloc for SALARY)
	// doSpent  → POST /accounts/:id/spent  (deducts balance + updates category budget)
	// doTransfer → POST /transfers         (moves balance between accounts)

	type topUpEntry struct {
		source, desc string
		amount       float64
		day, hour    int
	}
	type spentEntry struct {
		cat, desc string
		accountID string
		amount    float64
		day, hr   int
	}
	type transferEntry struct {
		from, to string
		amount   float64
		day      int
		note     string
	}

	doTopUp := func(y, mo int, accountID string, e topUpEntry) {
		_, err := post(fmt.Sprintf("/accounts/%s/topup", accountID), M{
			"source": e.source, "amount": e.amount,
			"date": ts(y, mo, e.day, e.hour, 0), "description": e.desc,
		})
		if err != nil {
			log.Printf("    ❌ TopUp [%s]: %v", e.source, err)
		} else {
			log.Printf("    ✅ TopUp  %-35s %s", e.source, rp(e.amount))
		}
	}

	doSpent := func(y, mo int, e spentEntry) {
		catID, ok := catIDs[e.cat]
		if !ok {
			log.Printf("    ⚠️  Kategori tidak ditemukan: %s", e.cat)
			return
		}
		_, err := post(fmt.Sprintf("/accounts/%s/spent", e.accountID), M{
			"category_id": catID, "amount": e.amount,
			"date": ts(y, mo, e.day, e.hr, 0), "description": e.desc,
		})
		if err != nil {
			log.Printf("    ❌ Spent [%s] %s: %v", e.cat, e.desc, err)
		} else {
			log.Printf("    ✅ Spent  [%-15s] %-35s %s", e.cat, e.desc, rp(e.amount))
		}
	}

	doTransfer := func(y, mo int, e transferEntry) {
		_, err := post("/transfers", M{
			"from_account_id": e.from, "to_account_id": e.to,
			"amount": e.amount, "note": e.note,
			"date": ts(y, mo, e.day, 10, 0),
		})
		if err != nil {
			log.Printf("    ❌ Transfer [%s]: %v", e.note, err)
		} else {
			log.Printf("    ✅ Transfer %-35s %s", e.note, rp(e.amount))
		}
	}

	// Langganan bulanan tetap (always from BCA)
	subscriptions := func(y, mo int) {
		for _, e := range []spentEntry{
			{"GitHub Copilot", "GitHub Copilot Monthly Subscription", bcaID, 180_000, 1, 0},
			{"Fore Coffee", "Fore+ Membership Monthly", bcaID, 24_000, 1, 0},
			{"Kuota Internet", "Paket Internet 50GB Telkomsel", bcaID, 100_000, 2, 10},
			{"Netflix", "Netflix Premium Plan", bcaID, 120_000, 5, 0},
			{"Spotify", "Spotify Individual Plan", bcaID, 60_000, 7, 0},
		} {
			doSpent(y, mo, e)
		}
	}

	// ─── STEP 4: Oktober 2025 ─────────────────────────────────────────────────
	sep("STEP 4 — Oktober 2025  |  BCA only")
	{
		y, mo := 2025, 10
		doTopUp(y, mo, bcaID, topUpEntry{"Gaji Oktober 2025", "Gaji bulanan transfer bank", 8_000_000, 1, 9})
		subscriptions(y, mo)
		for _, e := range []spentEntry{
			{"Makan", "Nasi Warteg + Minum", bcaID, 42_000, 2, 12},
			{"Makan", "Ayam Geprek + Es Teh", bcaID, 55_000, 5, 19},
			{"Makan", "Nasi Padang siang", bcaID, 38_000, 9, 13},
			{"Makan", "Makan malam berdua", bcaID, 70_000, 13, 19},
			{"Makan", "Soto ayam + Nasi", bcaID, 45_000, 18, 12},
			{"Makan", "Mie Ayam + Es Jeruk", bcaID, 35_000, 22, 13},
			{"Makan", "Nasi Goreng Spesial", bcaID, 40_000, 26, 20},
			{"Bensin", "Isi bensin Pertalite", bcaID, 50_000, 4, 8},
			{"Bensin", "Isi bensin Pertamax", bcaID, 50_000, 12, 17},
			{"Bensin", "Isi bensin full tank", bcaID, 50_000, 24, 7},
			{"Listrik", "Tagihan listrik PLN Oktober", bcaID, 178_000, 14, 15},
		} {
			doSpent(y, mo, e)
		}
	}

	// ─── STEP 5: November 2025 ────────────────────────────────────────────────
	sep("STEP 5 — November 2025  |  BCA + freelance")
	{
		y, mo := 2025, 11
		doTopUp(y, mo, bcaID, topUpEntry{"Gaji November 2025", "Gaji bulanan transfer bank", 8_000_000, 1, 9})
		doTopUp(y, mo, bcaID, topUpEntry{"Freelance Design", "Pembayaran proyek desain UI mobile app", 1_500_000, 15, 14})
		subscriptions(y, mo)
		for _, e := range []spentEntry{
			{"Makan", "Makan siang Warteg", bcaID, 50_000, 3, 12},
			{"Makan", "Nasi Padang + Lauk", bcaID, 40_000, 7, 13},
			{"Makan", "Makan malam Geprek", bcaID, 65_000, 11, 20},
			{"Makan", "Soto + Nasi + Minum", bcaID, 48_000, 16, 12},
			{"Makan", "Bakso Komplit", bcaID, 35_000, 20, 13},
			{"Makan", "Makan malam perayaan proyek", bcaID, 120_000, 28, 19},
			{"Bensin", "Isi bensin Pertalite pagi", bcaID, 50_000, 5, 7},
			{"Bensin", "Isi bensin Pertamax sore", bcaID, 50_000, 15, 18},
			{"Bensin", "Isi bensin full tank", bcaID, 50_000, 25, 8},
			{"Listrik", "Tagihan listrik PLN November", bcaID, 192_000, 12, 15},
		} {
			doSpent(y, mo, e)
		}
	}

	// ─── STEP 6: Desember 2025 ────────────────────────────────────────────────
	sep("STEP 6 — Desember 2025  |  Transfer aktif")
	{
		y, mo := 2025, 12
		doTopUp(y, mo, bcaID, topUpEntry{"Gaji Desember 2025", "Gaji bulanan transfer bank", 8_000_000, 1, 9})
		doTopUp(y, mo, bcaID, topUpEntry{"Bonus Akhir Tahun 2025", "Bonus kinerja tahunan dari perusahaan", 3_000_000, 20, 10})
		doTransfer(y, mo, transferEntry{bcaID, cashID, 2_000_000, 1, "Uang cash bulanan Desember"})
		doTransfer(y, mo, transferEntry{bcaID, savingsID, 2_000_000, 2, "Tabungan darurat Desember"})
		subscriptions(y, mo)
		for _, e := range []spentEntry{
			{"Makan", "Makan siang Warung Padang", cashID, 55_000, 2, 12},
			{"Makan", "Makan malam restoran keluarga", cashID, 85_000, 8, 19},
			{"Makan", "Nasi Warteg + Lauk", cashID, 45_000, 14, 12},
			{"Makan", "Dinner perayaan bonus", cashID, 75_000, 20, 20},
			{"Makan", "Makan siang Natal", cashID, 60_000, 25, 13},
			{"Makan", "Makan malam tahun baru", cashID, 150_000, 31, 20},
			{"Bensin", "Isi bensin Pertalite", cashID, 50_000, 3, 8},
			{"Bensin", "Isi bensin Pertamax Turbo", cashID, 50_000, 12, 7},
			{"Bensin", "Isi bensin untuk mudik", cashID, 100_000, 22, 17},
			{"Listrik", "Tagihan listrik PLN Desember", bcaID, 205_000, 12, 15},
		} {
			doSpent(y, mo, e)
		}
	}

	// ─── STEP 7: Januari 2026 ─────────────────────────────────────────────────
	sep("STEP 7 — Januari 2026  |  Awal tahun baru")
	{
		y, mo := 2026, 1
		doTopUp(y, mo, bcaID, topUpEntry{"Gaji Januari 2026", "Gaji bulanan transfer bank", 8_000_000, 1, 9})
		doTransfer(y, mo, transferEntry{bcaID, cashID, 2_000_000, 1, "Uang cash bulanan Januari"})
		doTransfer(y, mo, transferEntry{bcaID, savingsID, 1_500_000, 2, "Tabungan darurat Januari"})
		subscriptions(y, mo)
		for _, e := range []spentEntry{
			{"Makan", "Makan siang Warteg", cashID, 40_000, 3, 12},
			{"Makan", "Ayam Geprek + Minum", cashID, 52_000, 7, 13},
			{"Makan", "Nasi Padang siang", cashID, 38_000, 12, 12},
			{"Makan", "Makan malam berdua", cashID, 68_000, 17, 19},
			{"Makan", "Gado-gado + Es Teh", cashID, 30_000, 20, 13},
			{"Makan", "Mie Goreng Spesial", cashID, 45_000, 25, 19},
			{"Bensin", "Isi bensin Pertalite pagi", cashID, 50_000, 4, 7},
			{"Bensin", "Isi bensin Pertamax sore", cashID, 50_000, 16, 18},
			{"Bensin", "Isi bensin full tank", cashID, 50_000, 28, 8},
			{"Listrik", "Tagihan listrik PLN Januari", bcaID, 182_000, 12, 15},
		} {
			doSpent(y, mo, e)
		}
	}

	// ─── STEP 8: Februari 2026 (current) ─────────────────────────────────────
	sep("STEP 8 — Februari 2026  |  Bulan ini")
	{
		y, mo := 2026, 2
		doTopUp(y, mo, bcaID, topUpEntry{"Gaji Februari 2026", "Gaji bulanan transfer bank", 8_000_000, 1, 9})
		doTopUp(y, mo, bcaID, topUpEntry{"Freelance Website", "Pembayaran proyek website development client", 3_500_000, 10, 14})
		doTransfer(y, mo, transferEntry{bcaID, cashID, 2_000_000, 1, "Uang cash bulanan Februari"})
		doTransfer(y, mo, transferEntry{bcaID, savingsID, 2_000_000, 2, "Tabungan darurat Februari"})
		subscriptions(y, mo)
		for _, e := range []spentEntry{
			{"Makan", "Makan siang Warteg", cashID, 45_000, 2, 12},
			{"Makan", "Nasi Padang + Es Teh", cashID, 38_000, 5, 13},
			{"Makan", "Ayam Geprek + Minuman", cashID, 52_000, 8, 19},
			{"Makan", "Valentine dinner", cashID, 200_000, 14, 19},
			{"Makan", "Soto ayam + Nasi", cashID, 48_000, 15, 12},
			{"Makan", "Makan siang kantin", cashID, 35_000, 18, 12},
			{"Makan", "Makan berdua di Restoran", cashID, 75_000, 22, 20},
			{"Bensin", "Isi bensin Pertamax Shell", cashID, 50_000, 3, 8},
			{"Bensin", "Isi bensin Pertalite", cashID, 50_000, 10, 7},
			{"Bensin", "Isi bensin full tank", cashID, 50_000, 20, 18},
			{"Listrik", "Tagihan listrik token PLN Februari", bcaID, 185_000, 12, 15},
		} {
			doSpent(y, mo, e)
		}

		// Budget Reallocation: sisa Listrik Feb dialokasikan ke Makan
		reallocData, err := post("/budgets/reallocate", M{
			"from_category_id": catIDs["Listrik"],
			"to_category_id":   catIDs["Makan"],
			"amount":           15_000.0,
			"reason":           "Token PLN lebih hemat bulan ini, realokasi sisa ke makan",
			"month":            2, "year": 2026,
		})
		if err != nil {
			log.Printf("  ⚠️  Realokasi: %v", err)
		} else {
			log.Printf("  ✅ Realokasi Listrik→Makan Rp15.000 → %s", idOf(reallocData))
		}
	}

	// ─── STEP 9: Register Device ──────────────────────────────────────────────
	sep("STEP 9 — Notification: Register Device")
	_, err = post("/notifications/register-device", M{"onesignal_player_id": "demo-player-id-001"})
	if err != nil {
		log.Printf("  ⚠️  Register device: %v", err)
	} else {
		log.Println("  ✅ Device demo-player-id-001 registered")
	}

	// ─── STEP 10: Notification Settings ─────────────────────────────────────
	sep("STEP 10 — Notification Settings")
	// Valid types: ALLOCATION_REMINDER | BUDGET_ALERT | SAVINGS_GOAL | SCHEDULED_FUND
	notifDefs := []M{
		{
			"type": "ALLOCATION_REMINDER", "title": "Waktunya Set Alokasi Gaji",
			"body":         "Gaji sudah masuk! Lakukan topup ke akun BCA dan auto-alokasi budget bulan ini.",
			"day_of_month": 1, "time_of_day": "09:30", "is_enabled": true,
			"onesignal_player_id": "demo-player-id-001",
		},
		{
			"type": "BUDGET_ALERT", "title": "Cek Budget Pengeluaran",
			"body":         "Pantau sisa budget kategori pengeluaranmu di pertengahan bulan ini.",
			"day_of_month": 15, "time_of_day": "20:00", "is_enabled": true,
			"onesignal_player_id": "demo-player-id-001",
		},
		{
			"type": "SAVINGS_GOAL", "title": "Progress Tabungan Darurat",
			"body":         "Yuk cek seberapa dekat kamu ke target Rp25.000.000 tabungan darurat!",
			"day_of_month": 25, "time_of_day": "19:00", "is_enabled": true,
			"onesignal_player_id": "demo-player-id-001",
		},
		{
			"type": "SCHEDULED_FUND", "title": "Scheduled Fund Akan Berjalan",
			"body":         "Transfer otomatis bulanan ke Dompet Cash & Tabungan akan dieksekusi besok.",
			"day_of_month": 0, "time_of_day": "20:00", "is_enabled": true,
			"onesignal_player_id": "demo-player-id-001",
		},
	}
	for _, def := range notifDefs {
		data, err := post("/notifications/settings", def)
		if err != nil {
			log.Printf("  ❌ Notification [%v]: %v", def["title"], err)
		} else {
			log.Printf("  ✅ %-45s → %s", def["title"], idOf(data))
		}
	}

	// ─── STEP 11: Scheduled Funds ─────────────────────────────────────────────
	sep("STEP 11 — Scheduled Funds")
	sfDefs := []M{
		{
			"account_id": cashID, "from_account_id": bcaID,
			"schedule_type": "TRANSFER", "amount": 2_000_000.0, "day_of_month": 1,
			"description": "Auto transfer uang cash bulanan dari BCA ke Dompet",
		},
		{
			"account_id": savingsID, "from_account_id": bcaID,
			"schedule_type": "TRANSFER", "amount": 1_500_000.0, "day_of_month": 2,
			"description": "Auto tabungan darurat bulanan — target Rp25jt",
		},
	}
	for _, def := range sfDefs {
		data, err := post("/scheduled-funds", def)
		if err != nil {
			log.Printf("  ❌ Scheduled fund: %v", err)
		} else {
			log.Printf("  ✅ %-55s → %s", def["description"], idOf(data))
		}
	}

	// ─── Summary ──────────────────────────────────────────────────────────────
	log.Println()
	log.Println(strings.Repeat("═", 60))
	log.Println("  SEEDING SELESAI!")
	log.Println(strings.Repeat("═", 60))
	log.Printf("  Periode          : Oktober 2025 – Februari 2026 (5 bulan)")
	log.Printf("  Kategori         : %d", len(catDefs))
	log.Printf("  Budget Alerts    : %d (satu per kategori)", len(catIDs))
	log.Printf("  Akun             : BCA - Gajian | Dompet Cash | Tabungan Darurat")
	log.Printf("  Total Pemasukan  : Rp 52.000.000 (gaji + freelance + bonus)")
	log.Printf("  Transfer         : BCA→Cash Rp8jt | BCA→Tabungan Rp7.5jt (Des–Feb)")
	log.Printf("  Budget Realokasi : 1 (Listrik→Makan Feb 2026 Rp15.000)")
	log.Printf("  Notifikasi       : %d settings", len(notifDefs))
	log.Printf("  Scheduled Funds  : %d", len(sfDefs))
	log.Println()
	log.Printf("  Estimasi saldo akhir Februari 2026:")
	log.Printf("     BCA - Gajian       : ≈ Rp 32.155.000")
	log.Printf("     Dompet Cash        : ≈ Rp  4.224.000")
	log.Printf("     Tabungan Darurat   : ≈ Rp  5.500.000")
	log.Println()
	log.Println("  Swagger: http://localhost:8081/swagger/index.html")
}
