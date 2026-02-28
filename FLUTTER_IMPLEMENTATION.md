# Flutter Frontend Implementation Guide
## Finance Tracking App

> Base URL: `http://localhost:8081/api/v1`  
> Backend: Go + Gin + PostgreSQL  
> Frontend: Flutter (Dart)

---

## Daftar Isi

1. [Project Setup & Dependencies](#1-project-setup--dependencies)
2. [Struktur Folder](#2-struktur-folder)
3. [Core Layer (Network + Config)](#3-core-layer-network--config)
4. [Models (Dart Data Classes)](#4-models-dart-data-classes)
5. [Repositories (API Calls)](#5-repositories-api-calls)
6. [State Management (Riverpod)](#6-state-management-riverpod)
7. [Implementasi Layar per Fitur](#7-implementasi-layar-per-fitur)
8. [Progress Checklist](#8-progress-checklist)

---

## 1. Project Setup & Dependencies

### `pubspec.yaml`

```yaml
name: finance_tracking_app
description: Personal finance tracking mobile app
version: 1.0.0+1

environment:
  sdk: ">=3.0.0 <4.0.0"

dependencies:
  flutter:
    sdk: flutter

  # HTTP & Networking
  dio: ^5.4.0                    # HTTP client utama
  pretty_dio_logger: ^1.3.1      # log request/response di dev

  # State Management
  flutter_riverpod: ^2.5.1       # state management
  riverpod_annotation: ^2.3.3    # code generation untuk riverpod

  # Serialization
  freezed_annotation: ^2.4.1     # immutable models
  json_annotation: ^4.8.1        # json serialization

  # Local Storage
  flutter_secure_storage: ^9.0.0 # simpan token/base URL
  shared_preferences: ^2.2.2     # simpan preferensi user

  # UI & Charts
  fl_chart: ^0.68.0              # chart untuk analytics
  intl: ^0.19.0                  # format tanggal & currency (IDR)
  cached_network_image: ^3.3.1   # image caching
  shimmer: ^3.0.0                # loading skeleton

  # Navigation
  go_router: ^13.2.0             # declarative routing

  # Utilities
  equatable: ^2.0.5              # value equality
  dartz: ^0.10.1                 # functional: Either<Failure, T>
  uuid: ^4.3.3                   # generate UUID client-side jika perlu

dev_dependencies:
  flutter_test:
    sdk: flutter
  build_runner: ^2.4.8
  freezed: ^2.5.2
  json_serializable: ^6.7.1
  riverpod_generator: ^2.4.0
  custom_lint: ^0.6.4
  riverpod_lint: ^2.3.7
```

---

## 2. Struktur Folder

```
lib/
├── main.dart
├── app.dart                          # MaterialApp + GoRouter setup
│
├── core/
│   ├── config/
│   │   └── app_config.dart           # baseURL, timeout
│   ├── network/
│   │   ├── dio_client.dart           # Dio singleton + interceptors
│   │   ├── api_exception.dart        # custom exception class
│   │   └── network_info.dart         # cek koneksi internet
│   ├── error/
│   │   ├── failure.dart              # sealed class Failure
│   │   └── exception.dart
│   └── utils/
│       ├── currency_formatter.dart   # format IDR
│       ├── date_formatter.dart       # format tanggal Indonesia
│       └── validators.dart
│
├── models/                           # Dart data classes (Freezed)
│   ├── account.dart
│   ├── income.dart
│   ├── expense.dart
│   ├── category.dart
│   ├── budget.dart
│   ├── transaction.dart
│   ├── transfer.dart
│   ├── analytics.dart
│   ├── report.dart
│   ├── alert.dart
│   ├── notification_setting.dart
│   ├── scheduled_fund.dart
│   └── statistics.dart
│
├── repositories/                     # Akses API langsung
│   ├── account_repository.dart
│   ├── income_repository.dart
│   ├── expense_repository.dart
│   ├── category_repository.dart
│   ├── budget_repository.dart
│   ├── transaction_repository.dart
│   ├── transfer_repository.dart
│   ├── analytics_repository.dart
│   ├── report_repository.dart
│   ├── alert_repository.dart
│   ├── notification_repository.dart
│   ├── scheduled_fund_repository.dart
│   └── statistics_repository.dart
│
├── providers/                        # Riverpod providers
│   ├── account_provider.dart
│   ├── income_provider.dart
│   ├── expense_provider.dart
│   ├── category_provider.dart
│   ├── budget_provider.dart
│   ├── transfer_provider.dart
│   ├── analytics_provider.dart
│   ├── report_provider.dart
│   └── statistics_provider.dart
│
├── screens/
│   ├── home/
│   │   ├── home_screen.dart          # Dashboard utama
│   │   └── widgets/
│   │       ├── balance_card.dart
│   │       ├── quick_stats_row.dart
│   │       └── recent_transactions.dart
│   │
│   ├── accounts/
│   │   ├── accounts_screen.dart      # Daftar akun
│   │   ├── account_detail_screen.dart
│   │   ├── account_form_screen.dart
│   │   └── widgets/
│   │       └── account_card.dart
│   │
│   ├── transactions/
│   │   ├── transactions_screen.dart  # Feed semua transaksi
│   │   ├── add_income_screen.dart
│   │   ├── add_expense_screen.dart
│   │   └── widgets/
│   │       └── transaction_tile.dart
│   │
│   ├── transfer/
│   │   ├── transfer_screen.dart      # Daftar transfer
│   │   └── add_transfer_screen.dart
│   │
│   ├── budget/
│   │   ├── budget_screen.dart        # Overview budget
│   │   ├── budget_detail_screen.dart
│   │   ├── reallocation_screen.dart
│   │   └── widgets/
│   │       └── budget_progress_bar.dart
│   │
│   ├── analytics/
│   │   ├── analytics_screen.dart     # Charts & insights
│   │   └── widgets/
│   │       ├── spending_pattern_chart.dart
│   │       ├── category_pie_chart.dart
│   │       └── budget_performance_chart.dart
│   │
│   ├── reports/
│   │   ├── reports_screen.dart
│   │   └── report_detail_screen.dart
│   │
│   ├── categories/
│   │   ├── categories_screen.dart
│   │   └── category_form_screen.dart
│   │
│   ├── scheduled_funds/
│   │   ├── scheduled_funds_screen.dart
│   │   └── scheduled_fund_form_screen.dart
│   │
│   ├── notifications/
│   │   └── notification_settings_screen.dart
│   │
│   └── settings/
│       └── settings_screen.dart
│
└── shared/
    ├── widgets/
    │   ├── app_button.dart
    │   ├── amount_input.dart          # input dengan format IDR
    │   ├── account_selector.dart      # dropdown pilih akun
    │   ├── category_selector.dart
    │   ├── date_picker_field.dart
    │   ├── warning_banner.dart        # SOFT warning (budget)
    │   └── loading_overlay.dart
    └── themes/
        └── app_theme.dart
```

---

## 3. Core Layer (Network + Config)

### `lib/core/config/app_config.dart`

```dart
class AppConfig {
  static const String baseUrl = 'http://10.0.2.2:8081/api/v1'; // Android emulator
  // static const String baseUrl = 'http://localhost:8081/api/v1'; // Web/iOS simulator
  static const Duration connectTimeout = Duration(seconds: 10);
  static const Duration receiveTimeout = Duration(seconds: 30);
}
```

> **Catatan emulator Android**: `10.0.2.2` adalah alias ke `localhost` host mesin.  
> iOS Simulator gunakan `localhost` langsung.

---

### `lib/core/network/dio_client.dart`

```dart
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../config/app_config.dart';

final dioProvider = Provider<Dio>((ref) {
  final dio = Dio(BaseOptions(
    baseUrl: AppConfig.baseUrl,
    connectTimeout: AppConfig.connectTimeout,
    receiveTimeout: AppConfig.receiveTimeout,
    headers: {'Content-Type': 'application/json'},
  ));

  // Interceptor: log di debug mode
  dio.interceptors.add(PrettyDioLogger(
    requestHeader: false,
    requestBody: true,
    responseBody: true,
    error: true,
  ));

  // Interceptor: handle error global
  dio.interceptors.add(InterceptorsWrapper(
    onError: (DioException e, ErrorInterceptorHandler handler) {
      // map HTTP status ke ApiException
      handler.next(e);
    },
  ));

  return dio;
});
```

---

### `lib/core/error/failure.dart`

```dart
sealed class Failure {
  final String message;
  const Failure(this.message);
}

class ServerFailure extends Failure {
  const ServerFailure(super.message);
}

class NetworkFailure extends Failure {
  const NetworkFailure() : super('Tidak ada koneksi internet');
}

class ValidationFailure extends Failure {
  final Map<String, String> errors;
  const ValidationFailure(this.errors) : super('Validasi gagal');
}

class InsufficientBalanceFailure extends Failure {
  const InsufficientBalanceFailure() : super('Saldo akun tidak mencukupi');
}
```

---

## 4. Models (Dart Data Classes)

Semua model menggunakan **Freezed** untuk immutable + copyWith + toJson/fromJson.

### `lib/models/account.dart`

```dart
import 'package:freezed_annotation/freezed_annotation.dart';
part 'account.freezed.dart';
part 'account.g.dart';

enum AccountType { card, cash, savings }
enum IncomeType { salary, none }

@freezed
class Account with _$Account {
  const factory Account({
    required String id,
    required String name,
    required AccountType type,
    required double balance,
    String? incomeType,     // "SALARY" | null
    double? goalAmount,     // untuk SAVINGS
    String? color,
    @Default(false) bool isArchived,
    required DateTime createdAt,
  }) = _Account;

  factory Account.fromJson(Map<String, dynamic> json) =>
      _$AccountFromJson(json);
}

@freezed
class AccountSummary with _$AccountSummary {
  const factory AccountSummary({
    required Account account,
    required double totalIncome,
    required double totalExpense,
    required double totalTransferIn,
    required double totalTransferOut,
    required double netFlow,
  }) = _AccountSummary;

  factory AccountSummary.fromJson(Map<String, dynamic> json) =>
      _$AccountSummaryFromJson(json);
}
```

---

### `lib/models/income.dart`

```dart
@freezed
class Income with _$Income {
  const factory Income({
    required String id,
    required double amount,
    required String source,
    String? description,
    required String accountId,
    String? categoryId,
    required DateTime date,
    required DateTime createdAt,
  }) = _Income;

  factory Income.fromJson(Map<String, dynamic> json) =>
      _$IncomeFromJson(json);
}

@freezed
class IncomeCreateResult with _$IncomeCreateResult {
  const factory IncomeCreateResult({
    required Income income,
    @Default([]) List<BudgetAllocation> allocations,
    String? alreadyAllocatedWarning,   // SOFT warning
  }) = _IncomeCreateResult;

  factory IncomeCreateResult.fromJson(Map<String, dynamic> json) =>
      _$IncomeCreateResultFromJson(json);
}

@freezed
class BudgetAllocation with _$BudgetAllocation {
  const factory BudgetAllocation({
    required String categoryId,
    required String categoryName,
    required double allocatedAmount,
    required double percentage,
  }) = _BudgetAllocation;

  factory BudgetAllocation.fromJson(Map<String, dynamic> json) =>
      _$BudgetAllocationFromJson(json);
}
```

---

### `lib/models/expense.dart`

```dart
@freezed
class Expense with _$Expense {
  const factory Expense({
    required String id,
    required double amount,
    required String description,
    required String categoryId,
    required String accountId,
    required DateTime date,
    required DateTime createdAt,
  }) = _Expense;

  factory Expense.fromJson(Map<String, dynamic> json) =>
      _$ExpenseFromJson(json);
}

@freezed
class ExpenseCreateResult with _$ExpenseCreateResult {
  const factory ExpenseCreateResult({
    required Expense expense,
    String? budgetWarning,    // SOFT: "exceeds remaining budget" | "no budget allocated"
  }) = _ExpenseCreateResult;

  factory ExpenseCreateResult.fromJson(Map<String, dynamic> json) =>
      _$ExpenseCreateResultFromJson(json);
}
```

---

### `lib/models/category.dart`

```dart
enum CategoryType { dailyContinuous, usageBased, subscription }

@freezed
class ExpenseCategory with _$ExpenseCategory {
  const factory ExpenseCategory({
    required String id,
    required String name,
    required CategoryType type,
    double? dailyAmount,       // untuk DAILY_CONTINUOUS
    double? monthlyBudget,     // auto-computed = dailyAmount * 30
    String? description,
    required DateTime createdAt,
  }) = _ExpenseCategory;

  factory ExpenseCategory.fromJson(Map<String, dynamic> json) =>
      _$ExpenseCategoryFromJson(json);
}
```

---

### `lib/models/transfer.dart`

```dart
@freezed
class Transfer with _$Transfer {
  const factory Transfer({
    required String id,
    required String fromAccountId,
    required String toAccountId,
    required double amount,
    String? note,
    required DateTime transferDate,
    required DateTime createdAt,
  }) = _Transfer;

  factory Transfer.fromJson(Map<String, dynamic> json) =>
      _$TransferFromJson(json);
}
```

---

### `lib/models/statistics.dart`

```dart
@freezed
class MonthlyStats with _$MonthlyStats {
  const factory MonthlyStats({
    required int year,
    required int month,
    required double totalIncome,
    required double totalExpense,
    required double netSavings,
    required double savingsRate,
    required List<CategorySpending> topCategories,
  }) = _MonthlyStats;

  factory MonthlyStats.fromJson(Map<String, dynamic> json) =>
      _$MonthlyStatsFromJson(json);
}

@freezed
class OverviewStats with _$OverviewStats {
  const factory OverviewStats({
    required double totalBalance,
    required double monthIncome,
    required double monthExpense,
    required List<Account> accounts,
  }) = _OverviewStats;

  factory OverviewStats.fromJson(Map<String, dynamic> json) =>
      _$OverviewStatsFromJson(json);
}
```

---

## 5. Repositories (API Calls)

### `lib/repositories/account_repository.dart`

```dart
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/account.dart';
import '../core/network/dio_client.dart';

abstract class IAccountRepository {
  Future<List<Account>> getAccounts();
  Future<Account> createAccount(Map<String, dynamic> data);
  Future<Account> updateAccount(String id, Map<String, dynamic> data);
  Future<AccountSummary> getAccountSummary(String id);
  Future<Account> topUp(String id, double amount, String? note);
  Future<Account> archiveAccount(String id);
}

class AccountRepository implements IAccountRepository {
  final Dio _dio;
  AccountRepository(this._dio);

  @override
  Future<List<Account>> getAccounts() async {
    final res = await _dio.get('/accounts');
    final List data = res.data['data'] ?? res.data;
    return data.map((e) => Account.fromJson(e)).toList();
  }

  @override
  Future<Account> createAccount(Map<String, dynamic> data) async {
    final res = await _dio.post('/accounts', data: data);
    return Account.fromJson(res.data['data']);
  }

  @override
  Future<Account> updateAccount(String id, Map<String, dynamic> data) async {
    final res = await _dio.patch('/accounts/$id', data: data);
    return Account.fromJson(res.data['data']);
  }

  @override
  Future<AccountSummary> getAccountSummary(String id) async {
    final res = await _dio.get('/accounts/$id/summary');
    return AccountSummary.fromJson(res.data['data']);
  }

  @override
  Future<Account> topUp(String id, double amount, String? note) async {
    final res = await _dio.post('/accounts/$id/topup', data: {
      'amount': amount,
      if (note != null) 'note': note,
    });
    return Account.fromJson(res.data['data']);
  }

  @override
  Future<Account> archiveAccount(String id) async {
    final res = await _dio.post('/accounts/$id/archive');
    return Account.fromJson(res.data['data']);
  }
}

final accountRepositoryProvider = Provider<IAccountRepository>((ref) {
  return AccountRepository(ref.watch(dioProvider));
});
```

---

### `lib/repositories/expense_repository.dart`

```dart
class ExpenseRepository {
  final Dio _dio;
  ExpenseRepository(this._dio);

  // POST /expenses
  // Response: { data: { expense: {...}, budget_warning: "..." } }
  Future<ExpenseCreateResult> createExpense({
    required double amount,
    required String description,
    required String categoryId,
    required String accountId,
    required DateTime date,
  }) async {
    final res = await _dio.post('/expenses', data: {
      'amount': amount,
      'description': description,
      'category_id': categoryId,
      'account_id': accountId,             // WAJIB - Phase 4
      'date': date.toIso8601String(),
    });
    return ExpenseCreateResult.fromJson(res.data['data']);
  }

  // GET /expenses?month=2&year=2026
  Future<List<Expense>> getExpenses({int? month, int? year}) async {
    final res = await _dio.get('/expenses', queryParameters: {
      if (month != null) 'month': month,
      if (year != null) 'year': year,
    });
    final List data = res.data['data'] ?? [];
    return data.map((e) => Expense.fromJson(e)).toList();
  }
}
```

---

### `lib/repositories/income_repository.dart`

```dart
// POST /incomes
// Response: { data: { income: {...}, allocations: [...], already_allocated_warning: "..." } }
Future<IncomeCreateResult> createIncome({
  required double amount,
  required String source,
  required String accountId,          // WAJIB - Phase 4
  String? description,
  required DateTime date,
}) async {
  final res = await _dio.post('/incomes', data: {
    'amount': amount,
    'source': source,
    'account_id': accountId,
    if (description != null) 'description': description,
    'date': date.toIso8601String(),
  });
  return IncomeCreateResult.fromJson(res.data['data']);
}
```

---

### `lib/repositories/transfer_repository.dart`

```dart
// POST /transfers
Future<Transfer> createTransfer({
  required String fromAccountId,
  required String toAccountId,
  required double amount,
  String? note,
  required DateTime transferDate,
}) async {
  final res = await _dio.post('/transfers', data: {
    'from_account_id': fromAccountId,
    'to_account_id': toAccountId,
    'amount': amount,
    if (note != null) 'note': note,
    'transfer_date': transferDate.toIso8601String(),
  });
  return Transfer.fromJson(res.data['data']);
}

// GET /transfers?month=2&year=2026
Future<List<Transfer>> getTransfers({int? month, int? year}) async {
  final res = await _dio.get('/transfers', queryParameters: {
    if (month != null) 'month': month,
    if (year != null) 'year': year,
  });
  final List data = res.data['data'] ?? [];
  return data.map((e) => Transfer.fromJson(e)).toList();
}

// DELETE /transfers/:id  -> cancel/reverse
Future<void> cancelTransfer(String id) async {
  await _dio.delete('/transfers/$id');
}
```

---

### `lib/repositories/analytics_repository.dart`

```dart
// GET /analytics/spending-pattern?month=2&year=2026
Future<SpendingPattern> getSpendingPattern(int month, int year) async { ... }

// GET /analytics/category-comparison?month=2&year=2026
Future<List<CategoryComparison>> getCategoryComparison(int month, int year) async { ... }

// GET /analytics/top-spending?month=2&year=2026&limit=5
Future<List<TopSpending>> getTopSpending(int month, int year, {int limit = 5}) async { ... }

// GET /analytics/budget-performance?month=2&year=2026
Future<List<BudgetPerformance>> getBudgetPerformance(int month, int year) async { ... }
```

---

### `lib/repositories/report_repository.dart`

```dart
// GET /reports/monthly?month=2&year=2026
Future<MonthlyReport> getMonthlyReport(int month, int year) async { ... }

// GET /reports/monthly/export?month=2&year=2026&format=pdf
// Returns raw bytes (PDF or Excel)
Future<Uint8List> exportMonthlyReport(int month, int year, String format) async {
  final res = await _dio.get('/reports/monthly/export',
    queryParameters: {'month': month, 'year': year, 'format': format},
    options: Options(responseType: ResponseType.bytes),
  );
  return Uint8List.fromList(res.data);
}

// GET /reports/yearly?year=2026
Future<YearlyReport> getYearlyReport(int year) async { ... }
```

---

## 6. State Management (Riverpod)

### `lib/providers/account_provider.dart`

```dart
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../repositories/account_repository.dart';
import '../models/account.dart';
part 'account_provider.g.dart';

// Daftar semua akun (non-archived)
@riverpod
Future<List<Account>> accounts(AccountsRef ref) async {
  return ref.watch(accountRepositoryProvider).getAccounts();
}

// Summary satu akun
@riverpod
Future<AccountSummary> accountSummary(AccountSummaryRef ref, String id) async {
  return ref.watch(accountRepositoryProvider).getAccountSummary(id);
}

// ==================== Notifier untuk mutasi ====================

@riverpod
class AccountNotifier extends _$AccountNotifier {
  @override
  AsyncValue<List<Account>> build() {
    return const AsyncValue.loading();
  }

  Future<void> loadAccounts() async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(
      () => ref.read(accountRepositoryProvider).getAccounts(),
    );
  }

  Future<Account> createAccount(Map<String, dynamic> data) async {
    final repo = ref.read(accountRepositoryProvider);
    final account = await repo.createAccount(data);
    ref.invalidate(accountsProvider);  // refresh list
    return account;
  }

  Future<Account> topUp(String id, double amount) async {
    final repo = ref.read(accountRepositoryProvider);
    final updated = await repo.topUp(id, amount, null);
    ref.invalidate(accountsProvider);
    ref.invalidate(accountSummaryProvider(id));
    return updated;
  }
}
```

---

### `lib/providers/expense_provider.dart`

```dart
@riverpod
class ExpenseNotifier extends _$ExpenseNotifier {
  @override
  AsyncValue<ExpenseCreateResult?> build() => const AsyncValue.data(null);

  // Mengembalikan result termasuk budget_warning untuk ditampilkan di UI
  Future<ExpenseCreateResult> createExpense({
    required double amount,
    required String description,
    required String categoryId,
    required String accountId,
    required DateTime date,
  }) async {
    state = const AsyncValue.loading();
    try {
      final repo = ref.read(expenseRepositoryProvider);
      final result = await repo.createExpense(
        amount: amount,
        description: description,
        categoryId: categoryId,
        accountId: accountId,
        date: date,
      );
      state = AsyncValue.data(result);

      // Invalidate data terkait
      ref.invalidate(accountsProvider);
      ref.invalidate(monthlyStatsProvider);

      return result;
    } on DioException catch (e) {
      // HARD error dari server (400: insufficient balance)
      final msg = e.response?.data['error'] ?? 'Gagal mencatat pengeluaran';
      state = AsyncValue.error(ServerFailure(msg), StackTrace.current);
      rethrow;
    }
  }
}
```

---

### `lib/providers/statistics_provider.dart`

```dart
@riverpod
Future<OverviewStats> overviewStats(OverviewStatsRef ref) async {
  return ref.watch(statisticsRepositoryProvider).getOverview();
}

@riverpod
Future<MonthlyStats> monthlyStats(MonthlyStatsRef ref, {int? month, int? year}) async {
  final now = DateTime.now();
  return ref.watch(statisticsRepositoryProvider).getMonthlyStats(
    month: month ?? now.month,
    year: year ?? now.year,
  );
}
```

---

## 7. Implementasi Layar per Fitur

---

### 7.1 Home Screen (Dashboard)

**API yang dipakai:**
- `GET /statistics/overview` → total balance, income bulan ini, expense bulan ini
- `GET /statistics/monthly?month=X&year=Y` → savings rate, top categories
- `GET /transactions?limit=5` → 5 transaksi terakhir

**Widget utama:**
```
HomeScreen
├── BalanceCard           → totalBalance dari overview
├── QuickStatsRow         → income | expense | savings rate
├── AccountsCarousel      → slide tiap akun + balance
├── RecentTransactions    → list 5 terakhir
└── FAB (+ tambah cepat)  → bottom sheet pilih Income/Expense/Transfer
```

**Contoh `BalanceCard`:**
```dart
class BalanceCard extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final overview = ref.watch(overviewStatsProvider);
    return overview.when(
      loading: () => const Shimmer(...),
      error: (e, _) => Text('Gagal memuat'),
      data: (stats) => Card(
        child: Column(children: [
          Text('Total Balance', style: ...),
          Text(
            CurrencyFormatter.format(stats.totalBalance),
            style: Theme.of(context).textTheme.headlineLarge,
          ),
        ]),
      ),
    );
  }
}
```

---

### 7.2 Accounts Screen

**API yang dipakai:**
- `GET /accounts` → list semua akun
- `POST /accounts` → buat akun baru
- `PATCH /accounts/:id` → edit akun
- `GET /accounts/:id/summary` → detail flow per akun
- `POST /accounts/:id/topup` → tambah saldo manual
- `POST /accounts/:id/archive` → arsipkan

**Field saat buat akun:**
```dart
// CreateAccountRequest
{
  "name": "BCA Gajian",
  "type": "CARD",               // CARD | CASH | SAVINGS
  "balance": 5000000,
  "income_type": "SALARY",      // hanya untuk CARD
  "goal_amount": 25000000,      // hanya untuk SAVINGS
  "color": "#1565C0"
}
```

**AccountDetailScreen** menampilkan:
- Total income/expense masuk ke akun ini
- Transfer in/out
- History transaksi di akun ini (`GET /transfers/account?account_id=xxx`)

---

### 7.3 Tambah Pemasukan (Income)

**API:** `POST /incomes`

**Form fields:**
- Jumlah (amount) — wajib
- Sumber (source) — wajib, e.g. "Gaji Februari"
- Akun tujuan (account_id) — dropdown dari `GET /accounts` — **WAJIB**
- Tanggal
- Deskripsi (opsional)

**Flow setelah submit:**
```
1. Panggil POST /incomes
2. Cek response:
   - result.allocations.isNotEmpty → tampilkan bottom sheet "Budget Auto-Dialokasikan"
   - result.alreadyAllocatedWarning != null → tampilkan banner kuning
3. Refresh: accounts, transactions, budget, statistics
```

**Contoh UI feedback alokasi:**
```dart
if (result.allocations.isNotEmpty) {
  showModalBottomSheet(
    context: context,
    builder: (_) => AllocationResultSheet(allocations: result.allocations),
  );
}
```

---

### 7.4 Tambah Pengeluaran (Expense)

**API:** `POST /expenses`

**Form fields:**
- Jumlah (amount) — wajib
- Deskripsi — wajib
- Kategori (category_id) — dropdown — wajib
- Akun sumber (account_id) — dropdown — **WAJIB** (Phase 4)
- Tanggal

**HARD error handling (saldo tidak cukup):**
```dart
try {
  final result = await ref.read(expenseNotifierProvider.notifier).createExpense(...);

  // SOFT warning — transaksi tetap berhasil
  if (result.budgetWarning != null) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(result.budgetWarning!),
        backgroundColor: Colors.orange,
      ),
    );
  }
  Navigator.pop(context);  // kembali ke home

} on ServerFailure catch (e) {
  // HARD error — transaksi GAGAL (saldo tidak cukup)
  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      content: Text(e.message),    // "insufficient account balance"
      backgroundColor: Colors.red,
    ),
  );
}
```

---

### 7.5 Transfer Antar Akun

**API:** `POST /transfers`

**Form fields:**
- Dari akun (from_account_id) — dropdown
- Ke akun (to_account_id) — dropdown (exclude from_account)
- Jumlah
- Catatan (opsional)
- Tanggal

**Cancel transfer:** `DELETE /transfers/:id` — akan reverse balance kedua akun.

---

### 7.6 Budget Screen

**API yang dipakai:**
- `GET /budgets` → list semua budget aktif
- `GET /budgets/summary` → ringkasan budget bulan ini
- `POST /budgets/reallocate` → pindahkan sisa budget
- `GET /budgets/reallocations` → history realokasi

**BudgetProgressBar widget:**
```dart
class BudgetProgressBar extends StatelessWidget {
  final String categoryName;
  final double allocated;
  final double spent;

  double get percentage => spent / allocated;
  Color get barColor {
    if (percentage >= 1.0) return Colors.red;
    if (percentage >= 0.8) return Colors.orange;
    return Colors.green;
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(children: [
          Text(categoryName),
          const Spacer(),
          Text('${CurrencyFormatter.format(spent)} / ${CurrencyFormatter.format(allocated)}'),
        ]),
        LinearProgressIndicator(
          value: percentage.clamp(0.0, 1.0),
          color: barColor,
        ),
      ],
    );
  }
}
```

**Request realokasi:**
```dart
// POST /budgets/reallocate
{
  "from_category_id": "uuid-xxx",
  "to_category_id": "uuid-yyy",
  "amount": 200000,
  "reason": "Sisa budget makan dipindah ke tabungan"
}
```

---

### 7.7 Analytics Screen

**API yang dipakai:**
- `GET /analytics/spending-pattern?month=X&year=Y`
- `GET /analytics/category-comparison?month=X&year=Y`
- `GET /analytics/top-spending?month=X&year=Y&limit=5`
- `GET /analytics/budget-performance?month=X&year=Y`

**Charts menggunakan `fl_chart`:**

```dart
// Pie chart untuk category-comparison
PieChart(
  PieChartData(
    sections: categoryData.map((c) => PieChartSectionData(
      value: c.totalSpent,
      title: c.categoryName,
      color: Color(int.parse(c.color.replaceAll('#', '0xFF'))),
    )).toList(),
  ),
)

// Bar chart untuk budget-performance
BarChart(
  BarChartData(
    barGroups: performanceData.map((p) => BarChartGroupData(
      x: index,
      barRods: [
        BarChartRodData(toY: p.budgetAmount, color: Colors.grey),
        BarChartRodData(toY: p.spentAmount, color: p.isOverBudget ? Colors.red : Colors.blue),
      ],
    )).toList(),
  ),
)
```

---

### 7.8 Reports Screen

**API yang dipakai:**
- `GET /reports/monthly?month=X&year=Y`
- `GET /reports/monthly/export?format=pdf` → download PDF
- `GET /reports/monthly/export?format=excel` → download Excel
- `GET /reports/yearly?year=Y`

**Export + share:**
```dart
Future<void> exportAndShare(int month, int year, String format) async {
  final bytes = await ref.read(reportRepositoryProvider)
      .exportMonthlyReport(month, year, format);

  final dir = await getTemporaryDirectory();
  final ext = format == 'pdf' ? 'pdf' : 'xlsx';
  final file = File('${dir.path}/laporan_${year}_${month}.$ext');
  await file.writeAsBytes(bytes);

  await Share.shareXFiles([XFile(file.path)],
    text: 'Laporan Keuangan ${DateFormatter.monthYear(month, year)}');
}
```

---

### 7.9 Categories Screen

**API yang dipakai:**
- `GET /categories`
- `POST /categories`
- `PATCH /categories/:id`
- `DELETE /categories/:id`

**Form create category:**
```dart
// Type: DAILY_CONTINUOUS → tampilkan field daily_amount
// Type: USAGE_BASED / SUBSCRIPTION → tampilkan field monthly_budget

if (type == CategoryType.dailyContinuous) {
  // POST body:
  // { "name": "Makan", "type": "DAILY_CONTINUOUS", "daily_amount": 40000 }
  // monthly_budget akan auto-dihitung: 40000 * 30 = 1.200.000
} else {
  // POST body:
  // { "name": "Netflix", "type": "SUBSCRIPTION", "monthly_budget": 186000 }
}
```

---

### 7.10 Scheduled Funds Screen

**API yang dipakai:**
- `GET /scheduled-funds`
- `POST /scheduled-funds`
- `PATCH /scheduled-funds/:id`
- `DELETE /scheduled-funds/:id`

**Contoh create:**
```dart
// POST /scheduled-funds
{
  "name": "Transfer ke Cash Bulanan",
  "from_account_id": "uuid-bca",
  "to_account_id": "uuid-cash",
  "amount": 2000000,
  "schedule_day": 1,        // hari ke-1 setiap bulan
  "is_active": true
}
```

---

### 7.11 Notification Settings Screen

**API yang dipakai:**
- `GET /notifications/settings`
- `POST /notifications/settings`
- `PATCH /notifications/settings/:id`
- `DELETE /notifications/settings/:id`
- `POST /notifications/test` → kirim notif test

**Contoh create:**
```dart
// POST /notifications/settings
{
  "title": "Catat Pengeluaran",
  "body": "Jangan lupa catat pengeluaran hari ini!",
  "hour": 21,
  "minute": 0,
  "is_active": true,
  "repeat_type": "DAILY"     // DAILY | MONTHLY
}
```

---

## 8. Progress Checklist

### Phase 1 — Foundation (Setup)
- [ ] Buat project Flutter baru
- [ ] Setup `pubspec.yaml` dengan semua dependencies
- [ ] Setup folder struktur
- [ ] Buat `AppConfig` dengan base URL
- [ ] Buat Dio client + interceptors
- [ ] Buat `Failure` sealed class
- [ ] Setup `CurrencyFormatter` untuk IDR (Rp 1.500.000)
- [ ] Setup `DateFormatter` untuk format Indonesia
- [ ] Jalankan `dart pub get` + `build_runner`

### Phase 2 — Models & Repositories
- [ ] `Account` + `AccountSummary` model
- [ ] `Income` + `IncomeCreateResult` model
- [ ] `Expense` + `ExpenseCreateResult` model
- [ ] `Category` model
- [ ] `Transfer` model
- [ ] `Budget` + `BudgetSummary` model
- [ ] `Statistics` (Monthly + Overview) model
- [ ] `Analytics` models (4 types)
- [ ] `Report` models
- [ ] `ScheduledFund` model
- [ ] `NotificationSetting` model
- [ ] Semua repository dengan implementasi lengkap
- [ ] Jalankan `build_runner` untuk generate `.g.dart` + `.freezed.dart`

### Phase 3 — State Management
- [ ] `AccountProvider` + `AccountNotifier`
- [ ] `IncomeProvider` + `IncomeNotifier`
- [ ] `ExpenseProvider` + `ExpenseNotifier` (dengan HARD/SOFT handling)
- [ ] `TransferProvider` + `TransferNotifier`
- [ ] `BudgetProvider`
- [ ] `StatisticsProvider`
- [ ] `AnalyticsProvider`
- [ ] `CategoryProvider`

### Phase 4 — Core Screens
- [ ] **Home/Dashboard** (BalanceCard, QuickStats, RecentTransactions)
- [ ] **Accounts list + AccountCard**
- [ ] **Tambah Income** (dengan feedback alokasi)
- [ ] **Tambah Expense** (dengan HARD error + SOFT warning)
- [ ] **Transfer antar akun**
- [ ] **Transactions feed** (income + expense + transfer digabung)

### Phase 5 — Feature Screens
- [ ] **Budget overview + progress bar**
- [ ] **Budget reallocation form**
- [ ] **Categories CRUD** (dengan form DAILY_CONTINUOUS vs lainnya)
- [ ] **Account detail + summary**
- [ ] **Scheduled funds CRUD**
- [ ] **Notification settings CRUD**

### Phase 6 — Analytics & Reports
- [ ] **Analytics screen** (4 charts dengan fl_chart)
- [ ] **Monthly report screen**
- [ ] **Yearly report screen**
- [ ] **Export PDF + share**
- [ ] **Export Excel + share**

### Phase 7 — Polish
- [ ] GoRouter navigation setup
- [ ] Bottom navigation bar (Home, Akun, Transaksi, Anggaran, Analitik)
- [ ] Dark mode support
- [ ] Empty state widgets
- [ ] Error state widgets
- [ ] Pull-to-refresh di semua list
- [ ] Loading skeleton (Shimmer)
- [ ] Test di emulator Android + iOS Simulator

---

## Catatan Penting

### Handling HARD vs SOFT error di UI

| Kondisi | HTTP Status | UI Response |
|---|---|---|
| Saldo akun tidak cukup | 400 | SnackBar merah, tetap di form |
| Budget terlampaui | 200 + `budget_warning` | SnackBar oranye, berhasil masuk |
| Budget tidak ada | 200 + `budget_warning` | SnackBar oranye, berhasil masuk |
| Server error | 500 | Dialog error + retry button |

### Format tanggal ke API

Selalu kirim dalam format **RFC3339 / ISO8601**:
```dart
date.toUtc().toIso8601String()
// "2026-02-28T14:00:00.000Z"
```

### Format currency IDR

```dart
class CurrencyFormatter {
  static final _fmt = NumberFormat.currency(
    locale: 'id_ID',
    symbol: 'Rp ',
    decimalDigits: 0,
  );

  static String format(double amount) => _fmt.format(amount);
  // output: "Rp 1.500.000"
}
```

### SALARY auto-allocation

Auto-alokasi **hanya terjadi** jika:
- Akun yang dipilih bertipe `CARD`
- Akun memiliki `income_type = "SALARY"`
- Ada kategori aktif dengan budget

Tampilkan informasi ini di form tambah pemasukan sebagai tooltip/info di sebelah dropdown akun.
