# Finance Tracking App — Struktur Project

Aplikasi finance tracking pribadi — manajemen akun, alokasi budget otomatis, transfer, dan notifikasi via OneSignal.

---

## Backend — `be_finance_tracking_app/`

```
be_finance_tracking_app/
├── main.go                        # Entry point, routing, dependency wiring
├── go.mod / go.sum
├── config/
│   └── config.go                  # Env & DB config loader
├── database/
│   └── database.go                # GORM init, AutoMigrate, preMigrate helpers
├── models/
│   ├── account.go                 # Account, Income, Expense, Transfer structs
│   ├── transaction.go             # Transaction (unified read model)
│   └── notification.go            # NotificationSetting struct
├── repositories/
│   ├── account_repository.go
│   ├── income_repository.go
│   ├── expense_repository.go
│   ├── transfer_repository.go
│   ├── budget_repository.go
│   ├── category_repository.go
│   ├── scheduled_fund_repository.go
│   ├── notification_repository.go
│   ├── alert_repository.go
│   ├── allocation_repository.go
│   └── reallocation_repository.go
├── services/
│   ├── account_service.go
│   ├── income_service.go
│   ├── expense_service.go
│   ├── transfer_service.go
│   ├── budget_service.go
│   ├── category_service.go
│   ├── scheduled_fund_service.go
│   ├── scheduler_service.go
│   ├── notification_service.go
│   ├── alert_service.go
│   ├── analytics_service.go
│   ├── report_service.go
│   ├── statistics_service.go
│   └── transaction_service.go
├── handlers/
│   ├── account_handler.go
│   ├── transaction_handler.go
│   ├── transactions_handler.go
│   ├── transfer_handler.go
│   ├── budget_handler.go
│   ├── category_handler.go
│   ├── expense_handler.go
│   ├── scheduled_fund_handler.go
│   ├── notification_handler.go
│   ├── alert_handler.go
│   ├── analytics_handler.go
│   ├── report_handler.go
│   └── statistics_handler.go
├── utils/
│   ├── helpers.go
│   ├── export_excel.go
│   └── export_pdf.go
├── docs/                          # Swagger (auto-generated via swag init)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── scripts/
│   ├── seed/
│   │   └── main.go               # Seed demo data via HTTP API (11 steps)
│   └── cleanup/
│       └── main.go               # Hapus semua data seed
└── test/
    ├── architecture/
    │   └── clean_architecture_test.go
    ├── handlers/
    │   └── category_handler_test.go
    ├── integration/
    │   ├── api_integration_test.go
    │   ├── analytics_test.go
    │   ├── report_export_test.go
    │   └── transaction_test.go
    ├── repositories/
    │   ├── category_repository_test.go
    │   └── income_repository_test.go
    ├── services/
    │   ├── expense_service_test.go
    │   └── income_service_test.go
    └── helpers/
        ├── database_helper.go
        └── http_helper.go
```

---

## Frontend — `fe_finance_tracking_app/` *(Flutter — planned)*

```
fe_finance_tracking_app/
├── pubspec.yaml
├── lib/
│   ├── main.dart
│   ├── app.dart
│   ├── core/
│   │   ├── constants/
│   │   │   ├── api_constants.dart
│   │   │   └── app_colors.dart
│   │   ├── utils/
│   │   │   ├── currency_formatter.dart
│   │   │   └── date_helper.dart
│   │   └── widgets/
│   ├── models/
│   │   ├── account.dart
│   │   ├── transaction.dart
│   │   ├── budget.dart
│   │   ├── category.dart
│   │   ├── transfer.dart
│   │   ├── alert.dart
│   │   ├── notification_setting.dart
│   │   └── scheduled_fund.dart
│   ├── services/
│   │   ├── api_service.dart
│   │   ├── account_service.dart
│   │   ├── transaction_service.dart
│   │   ├── budget_service.dart
│   │   ├── transfer_service.dart
│   │   └── notification_service.dart
│   ├── providers/
│   │   ├── account_provider.dart
│   │   ├── budget_provider.dart
│   │   └── transaction_provider.dart
│   └── screens/
│       ├── home/
│       │   └── home_screen.dart
│       ├── accounts/
│       │   ├── accounts_screen.dart
│       │   ├── account_detail_screen.dart
│       │   ├── topup_screen.dart
│       │   ├── spent_screen.dart
│       │   └── transfer_screen.dart
│       ├── budget/
│       │   ├── budget_screen.dart
│       │   └── reallocate_screen.dart
│       ├── transactions/
│       │   └── transactions_screen.dart
│       ├── analytics/
│       │   └── analytics_screen.dart
│       ├── reports/
│       │   └── reports_screen.dart
│       ├── notifications/
│       │   └── notification_settings_screen.dart
│       └── scheduled_funds/
│           └── scheduled_funds_screen.dart
└── test/
```
