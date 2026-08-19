# PPMS — Project Portfolio Management System
## Technical Documentation
**Version:** 1.1.0 (Phase 0–10 Complete)
**Last Updated:** July 2026

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Technology Stack](#2-technology-stack)
3. [Architecture](#3-architecture)
4. [Project Structure](#4-project-structure)
5. [Environment Setup](#5-environment-setup)
6. [Database](#6-database)
7. [Authentication & Authorization](#7-authentication--authorization)
8. [API Reference](#8-api-reference)
9. [Event Bus](#9-event-bus)
10. [File Storage (MinIO)](#10-file-storage-minio)
11. [Notification System](#11-notification-system)
12. [Reporting](#12-reporting)
13. [Audit Logging](#13-audit-logging)
14. [Frontend Architecture](#14-frontend-architecture)
15. [Deployment](#15-deployment)
16. [Testing](#16-testing)
17. [Known Limitations & Roadmap](#17-known-limitations--roadmap)

---

## 1. Project Overview

PPMS adalah aplikasi internal perusahaan untuk mengelola siklus hidup proyek dari tahap pengajuan hingga penyelesaian. Sistem mencakup:

- **Project Request & Approval Workflow** — pengajuan, review, revisi, approval oleh ADMIN
- **Project Management** — lifecycle project, anggota, milestone, task dengan progress tracking
- **Budget Management** — alokasi anggaran, pencatatan transaksi immutable dengan idempotency
- **File Management** — upload/download dokumen ke MinIO dengan version tracking
- **Handover Management** — pencatatan pengiriman dan penerimaan dokumen fisik
- **Notification System** — notifikasi real-time berbasis event, dapat di-mute per tipe
- **Dashboard & Reporting** — ringkasan metrik, laporan PDF/Excel, global search
- **Portfolio View** — daftar project dengan filter/sorting, health status, ringkasan budget CAPEX/OPEX, pagu tahunan, dan Admin Direct Create
- **Import/Export** — backup & restore penuh seluruh data project dalam format JSON (ADMIN only)
- **Audit Trail** — seluruh aktivitas penting tercatat dengan redaksi field sensitif

---

## 2. Technology Stack

| Layer | Technology |
|---|---|
| Backend Language | Go 1.23 |
| Backend Framework | Gin |
| ORM | GORM |
| Authentication | JWT (golang-jwt/jwt v5) |
| Password Hashing | bcrypt (cost 12) |
| Database | PostgreSQL 16 |
| Object Storage | MinIO |
| PDF Generation | gofpdf (jung-kurt/gofpdf) |
| Excel Generation | excelize v2 (xuri/excelize) |
| Migration Tool | golang-migrate/migrate v4 |
| Logger | zerolog (rs/zerolog) |
| Frontend Framework | React 18 + TypeScript |
| Build Tool | Vite |
| UI Components | Shadcn/UI + Tailwind CSS |
| State Management | TanStack Query v5 |
| Form Validation | React Hook Form + Zod |
| Routing | React Router v6 |
| HTTP Client | Axios |
| Container | Docker + Docker Compose |
| Load Testing | k6 |

---

## 3. Architecture

### 3.1 High-Level Architecture

```
Browser (React SPA)
       │ HTTP/REST
       ▼
  Gin Backend (:8080 internal)
       │
       ├── PostgreSQL (data persistence)
       ├── MinIO (file storage)
       └── In-Memory Event Bus (cross-domain notifications)
```

### 3.2 Architectural Pattern

PPMS menggunakan **Modular Monolith** dengan prinsip **Clean Architecture**:

```
Handler  →  Service  →  Repository  →  Database
  ↑           ↑
(HTTP)    (Business Logic + Event Publishing)
```

Setiap domain (`auth`, `user`, `project`, `task`, dst.) bersifat **mandiri** dalam packagenya sendiri, berkomunikasi antar-domain **hanya melalui interface kecil** (bukan import langsung ke struct konkret) atau **event bus** untuk menghindari circular dependency.

### 3.3 Permission System (Dua Layer)

**Layer 1 — System Role** (level aplikasi):
- `ADMIN` — akses penuh ke semua resource
- `USER` — bisa membuat request, berpartisipasi di project
- `VIEWER` — read-only terbatas

**Layer 2 — Project Role** (level per-project):
- `PROJECT_MANAGER` — kelola project, member, milestone, task, budget, handover
- `MEMBER` — update task yang di-assign ke dirinya
- `OBSERVER` — read-only di dalam project

**Urutan Evaluasi Permission:**
1. ADMIN Override → selalu lolos
2. Ownership Check → apakah resource milik user ini?
3. Project Membership Check → apakah user member aktif di project?
4. Project Role Check → apakah role user di project cukup?
5. Default Deny

---

## 4. Project Structure

### 4.1 Backend (`ppms-backend/`)

```
ppms-backend/
├── cmd/api/
│   └── main.go                    ← entry point, wiring semua dependency
├── configs/
│   └── config.go                  ← load env vars ke struct Config
├── internal/
│   ├── database/
│   │   └── database.go            ← PostgreSQL connection + pool tuning
│   ├── middleware/
│   │   ├── auth.go                ← JWT validation, inject user_id/system_role ke context
│   │   ├── project_context.go     ← validasi project membership, inject project_role
│   │   ├── permission.go          ← RequireSystemRole & RequireProjectRole helpers
│   │   ├── rate_limit.go          ← in-memory rate limiter (IP & user based)
│   │   ├── cors.go                ← CORS configuration
│   │   └── security_headers.go    ← X-Content-Type-Options, X-Frame-Options, dll
│   ├── shared/
│   │   ├── errors/errors.go       ← error codes & AppError struct
│   │   ├── response/response.go   ← standar response envelope JSON
│   │   └── logger/logger.go       ← zerolog initialization
│   ├── events/
│   │   ├── bus.go                 ← in-memory pub/sub event bus
│   │   └── notification_subscriber.go ← base subscribers (approved/rejected/assigned/sent)
│   ├── infrastructure/
│   │   ├── minio/minio.go         ← MinIO client (upload, presigned URL, delete)
│   │   └── membership/            ← MembershipChecker adapter (cross-domain, Phase 7)
│   └── domain/
│       ├── auth/                  ← login, logout, token rotation, session management
│       ├── user/                  ← CRUD user, assign role, soft delete/restore
│       ├── division/              ← CRUD division
│       ├── project_request/       ← draft, submit, review, revision history
│       ├── project/               ← lifecycle, member management, progress computation
│       ├── milestone/             ← CRUD, reorder, status, progress computation
│       ├── task/                  ← CRUD, assignment, comments, progress
│       ├── budget/                ← alokasi budget, immutable transactions
│       ├── attachment/            ← upload/download/version, ownership validation
│       ├── handover/              ← shipment & receipt recording
│       ├── notification/          ← creation, read status, preferences
│       ├── dashboard/             ← aggregation queries
│       ├── audit/                 ← write-once audit log
│       ├── reporting/             ← PDF & Excel generation
│       ├── search/                ← pg_trgm full-text search
│       └── import_export/         ← full backup (export) & restore (import) JSON, ADMIN only
├── migrations/
│   ├── 000001_init_schema.up.sql  ← semua 18 tabel + indexes
│   ├── 000001_init_schema.down.sql
│   ├── 000002_seed_admin.up.sql   ← seed admin awal (Email: admin@ppms.local)
│   ├── 000002_seed_admin.down.sql
│   ├── 000003_add_performance_indexes.up.sql   ← Phase 7: idx_users_system_role, idx_tasks_due_date
│   └── 000003_add_performance_indexes.down.sql
├── loadtest/
│   └── scenario.js                ← k6 load test (NFR-01: 100 concurrent, p95 <3s)
├── scripts/
│   ├── backup.sh                  ← daily PostgreSQL backup script
│   └── restore.sh                 ← disaster recovery restore script
├── docs/
│   └── testing/
│       ├── e2e-checklist.md       ← acceptance criteria test checklist
│       └── soft-delete-checklist.md
├── go.mod
└── go.sum
```

Setiap domain mengikuti struktur internal yang sama:

```
internal/domain/<name>/
├── entity/      ← GORM model (tabel database)
├── dto/         ← request/response data transfer objects
├── errors/      ← domain-specific error variables
├── repository/  ← interface + implementasi GORM (akses database)
├── service/     ← business logic, validasi, event publishing
└── handler/     ← HTTP handler (binding request, panggil service, format response)
```

### 4.2 Frontend (`ppms-frontend/`)

```
ppms-frontend/
├── src/
│   ├── api/
│   │   └── axiosInstance.ts        ← axios dengan token interceptor + refresh logic
│   ├── components/
│   │   ├── ui/                     ← Shadcn/UI components (auto-generated)
│   │   ├── layout/
│   │   │   ├── DashboardLayout.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   ├── Navbar.tsx
│   │   │   ├── NotificationBell.tsx
│   │   │   └── GlobalSearchBar.tsx
│   │   └── shared/
│   │       └── FileUploadCard.tsx   ← reusable component upload/list attachment
│   ├── features/                    ← per-domain feature modules
│   │   ├── auth/                   ← login, auth context, hooks
│   │   ├── users/
│   │   ├── divisions/
│   │   ├── project-requests/
│   │   ├── projects/
│   │   ├── milestones/
│   │   ├── tasks/
│   │   ├── budgets/
│   │   ├── attachments/
│   │   ├── handovers/
│   │   ├── notifications/
│   │   ├── dashboard/
│   │   ├── reporting/
│   │   ├── search/
│   │   └── audit/
│   ├── routes/
│   │   ├── AppRouter.tsx           ← semua route definition
│   │   └── ProtectedRoute.tsx      ← route guard berdasarkan system_role
│   ├── lib/
│   │   ├── queryClient.ts          ← TanStack Query global config
│   │   ├── utils.ts
│   │   └── errorMessages.ts        ← mapping kode error → pesan Bahasa Indonesia
│   ├── types/
│   │   └── index.ts                ← shared TypeScript types
│   └── main.tsx                    ← React entry point
├── Dockerfile                      ← dev: npm run dev
└── Dockerfile.prod                 ← prod: nginx multi-stage build
```

---

## 5. Environment Setup

### 5.1 Prerequisites

- Docker & Docker Compose v2.1+
- Go 1.23+
- Node.js 20+
- golang-migrate CLI: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
- k6 (opsional, untuk load testing)

### 5.2 Environment Variables (Backend)

Salin `.env.example` ke `.env` dan isi nilainya:

```env
# Server
APP_ENV=development          # development | production
APP_PORT=8080
APP_NAME=ppms-backend

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=ppms_user
DB_PASSWORD=change_me
DB_NAME=ppms_db
DB_SSLMODE=disable           # gunakan "require" di production

# JWT
JWT_ACCESS_SECRET=change_me_access_secret    # WAJIB GANTI di production
JWT_REFRESH_SECRET=change_me_refresh_secret  # WAJIB GANTI di production
JWT_ACCESS_EXPIRY_MINUTES=15
JWT_REFRESH_EXPIRY_DAYS=7

# MinIO
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minio_admin
MINIO_SECRET_KEY=change_me_minio_secret
MINIO_BUCKET=ppms-attachments
MINIO_USE_SSL=false

# Bcrypt
BCRYPT_COST=12

# CORS
FRONTEND_URL=http://localhost:5174
```

Generate JWT secret yang aman untuk production:

```bash
openssl rand -base64 32
```

### 5.3 Environment Variables (Frontend)

```env
VITE_API_BASE_URL=http://localhost:8081/api/v1
```

### 5.4 Development Quick Start

```bash
# 1. Clone & masuk ke root project
cd ppms

# 2. Salin environment variables
cp ppms-backend/.env.example ppms-backend/.env
echo "VITE_API_BASE_URL=http://localhost:8081/api/v1" > ppms-frontend/.env

# 3. Jalankan infrastruktur (PostgreSQL + MinIO + auto-create bucket)
docker compose up -d postgres minio minio-init

# 4. Tunggu healthy, lalu jalankan semua migration
migrate -path ppms-backend/migrations \
  -database "postgres://ppms_user:change_me@localhost:5432/ppms_db?sslmode=disable" up

# 5. Build & jalankan backend dan frontend
docker compose up -d --build backend frontend
```

Akses:
- **Frontend:** http://localhost:5174
- **Backend API:** http://localhost:8081/api/v1
- **Health Check:** http://localhost:8081/health
- **Readiness:** http://localhost:8081/ready
- **MinIO Console:** http://localhost:9001

### 5.5 Default Admin Credentials

| Field | Value |
|---|---|
| Email | admin@ppms.local |
| Password | Admin@12345 |

> **Wajib ganti password ini setelah first login!**

---

## 6. Database

### 6.1 Schema Overview

18 tabel utama, dibagi per domain:

| Domain | Tabel |
|---|---|
| User & Auth | `users`, `divisions`, `user_sessions` |
| Project Request | `project_requests`, `project_request_revisions`, `project_request_approvals` |
| Project | `projects`, `project_members` |
| Milestone & Task | `milestones`, `tasks`, `task_assignees`, `task_comments` |
| Budget | `budgets`, `budget_transactions` |
| Storage & Handover | `attachments`, `handovers` |
| Notification | `notifications`, `notification_preferences` |
| Audit | `audit_logs` |

### 6.2 Key Design Decisions

**Soft Delete:** Semua entity utama menggunakan soft delete (`deleted_at`, `deleted_by`). Query default selalu menyertakan `WHERE deleted_at IS NULL`.

**Optimistic Locking:** Semua tabel yang diubah user memiliki kolom `version INTEGER DEFAULT 1`. Update hanya berhasil jika `WHERE id = ? AND version = ?` cocok (FR-17.01). Jika `RowsAffected == 0`, kembalikan `ErrConflict`.

**Audit Log Immutability:** Tabel `audit_logs` tidak memiliki endpoint UPDATE atau DELETE. `old_data`/`new_data` selalu di-redact dari field sensitif (`password_hash`, `refresh_token_hash`, dll) sebelum disimpan.

**Budget Transaction Immutability:** `budget_transactions` tidak memiliki repository method Update/Delete. Kesalahan dikoreksi via transaksi `ADJUSTMENT` baru.

**Polymorphic Attachment:** Tabel `attachments` menggunakan `(entity_type, entity_id)` untuk menghubungkan ke berbagai entitas. Composite index `idx_attachments_entity` memastikan query tetap performant.

**Idempotency:** Kolom `budget_transactions.idempotency_key` (UNIQUE) mencegah transaksi duplikat akibat retry atau double-click. Jika key sudah ada, service mengembalikan transaksi yang ada tanpa insert baru.

### 6.3 Migration

```bash
# Jalankan semua migration yang pending
migrate -path migrations -database "<DSN>" up

# Rollback migration terakhir
migrate -path migrations -database "<DSN>" down 1

# Buat migration baru
migrate create -ext sql -dir migrations -seq <nama_migration>
```

### 6.4 Progress Calculation

- **Task Progress:** Input manual user (0–100). Task CANCELLED dikecualikan dari semua kalkulasi.
- **Milestone Progress:** `AVG(task.progress)` — hanya task yang statusnya bukan CANCELLED.
- **Project Progress:** `AVG(milestone.progress)` — hanya milestone yang statusnya bukan CANCELLED.

Kalkulasi ini dilakukan **on-the-fly** saat endpoint dipanggil (computed, bukan disimpan di database), kecuali `task.progress` yang memang disimpan.

---

## 7. Authentication & Authorization

### 7.1 JWT Flow

```
Login → [access_token (15 menit)] + [refresh_token (7 hari)]
                                           │
                             Disimpan hash-nya di user_sessions
                                           │
                              Saat access_token expired:
                              POST /auth/refresh → token baru (rotation)
                              Refresh token lama langsung direvoke
```

Semua request ke endpoint protected harus menyertakan header:
```
Authorization: Bearer <access_token>
```

### 7.2 Refresh Token Rotation

Setiap kali refresh token digunakan, token tersebut langsung direvoke dan pasangan token baru diterbitkan. Ini mencegah replay attack jika token bocor.

### 7.3 Middleware Stack (per request)

```
Request
  → CORS Middleware
  → Security Headers Middleware
  → Rate Limiter (untuk endpoint tertentu)
  → Auth Middleware (validasi JWT, inject context)
  → Project Context Middleware (inject project_role, hanya untuk route /projects/:id/*)
  → RequireSystemRole / RequireProjectRole
  → Handler
```

### 7.4 Context Keys

Setelah Auth Middleware berjalan, nilai-nilai berikut tersedia di Gin context:

| Key | Type | Deskripsi |
|---|---|---|
| `user_id` | `uint64` | ID user yang sedang login |
| `system_role` | `string` | `ADMIN`, `USER`, atau `VIEWER` |
| `division_id` | `*uint64` | ID divisi user (bisa nil) |

Setelah Project Context Middleware (hanya route `/projects/:id/*`):

| Key | Type | Deskripsi |
|---|---|---|
| `project_id` | `uint64` | ID project dari URL param |
| `project_role` | `string` | `PROJECT_MANAGER`, `MEMBER`, `OBSERVER`, atau `ADMIN_OVERRIDE` |
| `project_member_id` | `uint64` | ID record di tabel project_members |

---

## 8. API Reference

### 8.1 Base URL

```
Development:  http://localhost:8081/api/v1
Production:   https://api.ppms.yourcompany.com/api/v1
```

### 8.2 Standard Response Format

**Success:**
```json
{
  "success": true,
  "data": { ... },
  "message": "optional success message"
}
```

**Paginated:**
```json
{
  "success": true,
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 157
  }
}
```

**Error:**
```json
{
  "success": false,
  "code": "INSUFFICIENT_PROJECT_ROLE",
  "message": "insufficient project role permission"
}
```

### 8.3 Error Codes

| Code | HTTP Status | Deskripsi |
|---|---|---|
| `INSUFFICIENT_SYSTEM_ROLE` | 403 | System role tidak cukup |
| `INSUFFICIENT_PROJECT_ROLE` | 403 | Project role tidak cukup |
| `NOT_PROJECT_MEMBER` | 403 | Bukan member aktif project |
| `RESOURCE_NOT_OWNED` | 403 | Resource bukan milik user ini |
| `PROJECT_LOCKED` | 409 | Project terkunci untuk perubahan |
| `LAST_PM_PROTECTION` | 409 | PM terakhir tidak bisa dihapus/di-demote |
| `INVALID_STATE_TRANSITION` | 400 | Transisi status tidak valid |
| `NOT_FOUND` | 404 | Resource tidak ditemukan |
| `VALIDATION_ERROR` | 400 | Data input tidak valid |
| `UNAUTHORIZED` | 403 | Token tidak valid atau expired |
| `CONFLICT` | 409 | Optimistic locking conflict (version mismatch) |
| `RATE_LIMITED` | 429 | Terlalu banyak request |
| `FILE_TOO_LARGE` | 400 | File melebihi 25MB |
| `UNSUPPORTED_FILE_TYPE` | 400 | MIME type tidak diizinkan |
| `DUPLICATE_ENTRY` | 409 | Idempotency key sudah dipakai |
| `INTERNAL_ERROR` | 500 | Server error |

### 8.4 Endpoint Summary

#### Auth (Public)
| Method | Path | Deskripsi |
|---|---|---|
| POST | `/auth/login` | Login, return access + refresh token |
| POST | `/auth/refresh` | Rotate refresh token |
| POST | `/auth/logout` | Revoke session |
| POST | `/auth/verify-otp` | Verify OTP untuk 2FA |
| POST | `/auth/resend-otp` | Resend OTP |

#### Auth (Protected)
| Method | Path | Deskripsi |
|---|---|---|
| POST | `/auth/change-password` | Ganti password, revoke semua sesi lain |
| POST | `/auth/revoke-sessions` | Revoke semua sesi aktif |

#### Profile & Settings (Protected)
| Method | Path | Deskripsi |
|---|---|---|
| GET | `/me` | Ambil profil user |
| PUT | `/me` | Update profil user |
| POST | `/me/photo` | Upload foto profil |
| POST | `/me/toggle-2fa` | Toggle 2FA |
| POST | `/me/toggle-email-notification` | Toggle email notification |

#### Division
| Method | Path | Role |
|---|---|---|
| GET | `/divisions` | Semua authenticated |
| GET | `/divisions/:id` | Semua authenticated |
| POST | `/divisions` | ADMIN |
| PUT | `/divisions/:id` | ADMIN |
| DELETE | `/divisions/:id` | ADMIN |

#### User Management
| Method | Path | Role |
|---|---|---|
| GET | `/users` | ADMIN |
| GET | `/users/:id` | ADMIN |
| POST | `/users` | ADMIN |
| PUT | `/users/:id` | ADMIN |
| PATCH | `/users/:id/role` | ADMIN |
| DELETE | `/users/:id` | ADMIN |
| POST | `/users/:id/restore` | ADMIN |

#### Project Request
| Method | Path | Role |
|---|---|---|
| POST | `/project-requests` | USER, ADMIN |
| GET | `/project-requests` | USER (own only), ADMIN (all) |
| GET | `/project-requests/:id` | USER (owner), ADMIN |
| PUT | `/project-requests/:id` | USER (owner, DRAFT only) |
| DELETE | `/project-requests/:id` | USER (owner, DRAFT only) |
| POST | `/project-requests/:id/submit` | USER (owner) |
| POST | `/project-requests/:id/review` | ADMIN |
| POST | `/project-requests/:id/revise` | USER (owner) |
| GET | `/project-requests/:id/revisions` | USER (owner), ADMIN |
| GET | `/project-requests/:id/approvals` | USER (owner), ADMIN |

`POST /project-requests/:id/review` menerima body:

```json
{
  "action": "APPROVED",
  "comment": "Approved for execution",
  "project_manager_id": 12
}
```

`project_manager_id` wajib saat `action = APPROVED` dan harus menunjuk user aktif dengan system role `USER` atau `ADMIN`. `REQUEST_REVISION` mengubah request menjadi `REVISION_REQUESTED`, bukan `REJECTED`, sehingga revisi dan penolakan permanen tidak tertukar. Setelah user mengirim revisi, request berstatus `REVISED` bisa direview ulang oleh ADMIN untuk disetujui, ditolak, atau diminta revisi lagi.

#### Project
| Method | Path | Role |
|---|---|---|
| GET | `/projects` | USER (own projects), ADMIN (all) |
| GET | `/projects/:id` | Member / ADMIN |
| PUT | `/projects/:id` | PROJECT_MANAGER / ADMIN |
| PATCH | `/projects/:id/status` | PROJECT_MANAGER / ADMIN |

#### Member Management
| Method | Path | Role |
|---|---|---|
| GET | `/projects/:id/members` | Member / ADMIN |
| POST | `/projects/:id/members` | PROJECT_MANAGER / ADMIN |
| PATCH | `/projects/:id/members/:memberId/role` | PROJECT_MANAGER / ADMIN |
| DELETE | `/projects/:id/members/:memberId` | PROJECT_MANAGER / ADMIN |

#### Milestone
| Method | Path | Role |
|---|---|---|
| GET | `/projects/:id/milestones` | Member / ADMIN |
| POST | `/projects/:id/milestones` | PROJECT_MANAGER / ADMIN |
| PATCH | `/projects/:id/milestones/reorder` | PROJECT_MANAGER / ADMIN |
| PUT | `/projects/:id/milestones/:milestoneId` | PROJECT_MANAGER / ADMIN |
| PATCH | `/projects/:id/milestones/:milestoneId/status` | PROJECT_MANAGER / ADMIN |
| DELETE | `/projects/:id/milestones/:milestoneId` | PROJECT_MANAGER / ADMIN |

#### Task
| Method | Path | Role |
|---|---|---|
| GET | `/projects/:id/tasks` | Member / ADMIN |
| POST | `/projects/:id/tasks` | PROJECT_MANAGER / ADMIN |
| PUT | `/projects/:id/tasks/:taskId` | PROJECT_MANAGER / ADMIN |
| PATCH | `/projects/:id/tasks/:taskId/status` | PM / MEMBER (assigned only) |
| PATCH | `/projects/:id/tasks/:taskId/progress` | PM / MEMBER (assigned only) |
| POST | `/projects/:id/tasks/:taskId/assignees` | PROJECT_MANAGER / ADMIN |
| DELETE | `/projects/:id/tasks/:taskId` | PROJECT_MANAGER / ADMIN |
| GET | `/projects/:id/tasks/:taskId/comments` | Member / ADMIN |
| POST | `/projects/:id/tasks/:taskId/comments` | PM / MEMBER |

#### Budget & Transaction
| Method | Path | Role |
|---|---|---|
| GET | `/projects/:id/budget` | Member / ADMIN |
| POST | `/projects/:id/budget` | PROJECT_MANAGER / ADMIN |
| PUT | `/projects/:id/budget/:budgetId` | PROJECT_MANAGER / ADMIN |
| GET | `/projects/:id/budget/:budgetId/transactions` | Member / ADMIN |
| POST | `/projects/:id/budget/:budgetId/transactions` | PROJECT_MANAGER / ADMIN |

#### Handover
| Method | Path | Role |
|---|---|---|
| GET | `/projects/:id/handovers` | Member / ADMIN |
| POST | `/projects/:id/handovers` | PROJECT_MANAGER / ADMIN |
| PATCH | `/projects/:id/handovers/:handoverId/received` | Member / ADMIN |
| PATCH | `/projects/:id/handovers/:handoverId/returned` | PROJECT_MANAGER / ADMIN |

#### Attachment
| Method | Path | Catatan |
|---|---|---|
| POST | `/attachments` | multipart/form-data; ownership divalidasi di service |
| GET | `/attachments?entity_type=X&entity_id=Y` | Ownership divalidasi di service |
| GET | `/attachments/:id/download` | Mengembalikan presigned URL (15 menit) |
| GET | `/attachments/:id/versions` | Riwayat versi file dengan nama sama |
| DELETE | `/attachments/:id` | Soft delete DB + delete dari MinIO |

#### Notification
| Method | Path | Catatan |
|---|---|---|
| GET | `/notifications` | Query params: `page`, `limit`, `unread_only=true` |
| PATCH | `/notifications/:id/read` | |
| PATCH | `/notifications/read-all` | |
| GET | `/notifications/preferences` | |
| PUT | `/notifications/preferences` | Body: `{"type": "TASK_ASSIGNED", "enabled": false}` |

#### Dashboard, Search, Audit, Reporting
| Method | Path | Role |
|---|---|---|
| GET | `/dashboard` | Semua authenticated (data berbeda per role) |
| GET | `/search?q=<query>` | Semua authenticated (scope sesuai akses) |
| GET | `/audit-logs` | ADMIN |
| POST | `/reports/generate` | ADMIN (system-wide) |
| POST | `/projects/:id/reports/generate` | PROJECT_MANAGER / ADMIN (per project) |

#### Report Request Body
```json
{
  "type": "PROJECT",
  "format": "PDF"
}
```
`type` options: `PROJECT`, `MILESTONE`, `TASK`, `BUDGET`, `HANDOVER`
`format` options: `PDF`, `EXCEL`

#### Admin / Portfolio (ADMIN only)
| Method | Path | Deskripsi |
|---|---|---|
| POST | `/admin/projects` | Buat project langsung tanpa request workflow (auto-generate `project_code`, actor jadi PROJECT_MANAGER, budget opsional). Audit `CREATE_PROJECT_DIRECT`. |
| GET | `/admin/budget-years` | List pagu tahunan CAPEX/OPEX |
| POST | `/admin/budget-years` | Tambah pagu tahunan |
| PUT | `/admin/budget-years/:id` | Update pagu tahunan (optimistic locking) |
| DELETE | `/admin/budget-years/:id` | Hapus pagu tahunan |
| GET | `/admin/export` | Download backup JSON seluruh project + relasi (members, milestones, tasks, budget, transactions). Audit `EXPORT_DATA`. |
| POST | `/admin/import` | Upload backup JSON (multipart `file`) untuk restore. Tiap project dibuat ulang sebagai project baru dengan kode baru. Audit `IMPORT_DATA`. |

#### Deadline / Portfolio (Semua authenticated, scope per role)
| Method | Path | Deskripsi |
|---|---|---|
| GET | `/projects/deadline?window=<overdue\|30\|60\|90>` | Project yang mendekati / melewati deadline |

**Import/Export JSON Backup**

`GET /admin/export` mengembalikan file `ppms-backup-<timestamp>.json` berisi
`{ version, exported_at, exported_by, projects: [...] }`. Setiap project menyertakan
seluruh members, milestones (dengan `ref_id`), tasks (dengan `milestone_ref_id`),
budget, dan transactions.

`POST /admin/import` menerima file JSON tersebut (maks. 20MB). Import bersifat
**non-destruktif**: tiap project dibuat sebagai project baru (`project_code` baru,
`created_by` = admin yang mengimpor). Milestone `ref_id` dipetakan ulang ke id baru
agar relasi task tetap konsisten. Member yang user-nya tidak ada akan dilewati, dan
selalu dipastikan ada minimal satu PROJECT_MANAGER aktif (fallback ke actor).
Response berisi ringkasan `{ total_projects, imported, skipped, errors, imported_project_ids }`.

---

## 9. Event Bus

PPMS menggunakan **in-memory synchronous event bus** untuk komunikasi antar-domain tanpa circular import. Event di-publish oleh service dan dikonsumsi oleh subscriber yang didaftarkan di `main.go`.

### 9.1 Event Catalog

| Event Name | Dipublish Oleh | Subscriber |
|---|---|---|
| `project.request.submitted` | RequestService | Notify semua ADMIN |
| `project.request.approved` | RequestService | Notify requester |
| `project.request.rejected` | RequestService | Notify requester |
| `project.request.revised` | RequestService | Notify semua ADMIN |
| `task.created` | TaskService | (logging) |
| `task.assigned` | TaskService | Notify user yang di-assign |
| `task.completed` | TaskService | Notify semua PM project |
| `task.overdue` | Scheduler (main.go) | Notify assignee |
| `task.due_soon` | Scheduler (main.go) | Notify assignee (H-1 sebelum due date) |
| `budget.warning` | TransactionService | Notify semua PM project (usage ≥ 80%) |
| `budget.over_limit` | TransactionService | Notify semua PM project (usage ≥ 100%) |
| `milestone.completed` | MilestoneService | Notify semua member project |
| `project.completed` | ProjectService | Notify semua member project |
| `handover.created` | HandoverService | Notify receiver + konfirmasi ke sender |
| `handover.sent` | HandoverService | Notify receiver (jika ada) |
| `handover.received` | HandoverService | Notify sender |

### 9.2 Menambahkan Event Baru

```go
// 1. Publish di service layer
s.eventBus.Publish(events.Event{
    Name: "entity.action_happened",
    Data: map[string]interface{}{
        "entity_id": 42,
        "user_id":   7,
    },
})

// 2. Subscribe di main.go (atau di events/notification_subscriber.go)
eventBus.Subscribe("entity.action_happened", func(e events.Event) {
    data := e.Data.(map[string]interface{})
    // lakukan sesuatu
})
```

> **Penting:** Subscriber berjalan di goroutine terpisah (lihat `bus.go`: `go handler(event)`). Pastikan handler aman untuk concurrent access dan tidak panic, karena panic di goroutine akan crash seluruh server.

---

## 10. File Storage (MinIO)

### 10.1 Bucket Structure

Satu bucket `ppms-attachments`, dengan object path:
```
{ENTITY_TYPE}/{ENTITY_ID}/{UUID}_{ORIGINAL_FILENAME}
```
Contoh:
```
TASK/15/a3f9c2d1-..._requirements.pdf
PROJECT/3/7b8e1f0a-..._charter.docx
```

### 10.2 Upload Flow

```
Frontend → POST /attachments (multipart/form-data)
         → AttachmentHandler
         → AttachmentService:
             1. Validasi entity_type
             2. Validasi ownership (project membership / request ownership)
             3. Validasi ukuran file (max 25MB)
             4. Deteksi MIME via magic bytes (512 byte pertama)
             5. Upload ke MinIO
             6. Simpan metadata ke tabel attachments
         → Response: attachment record
```

### 10.3 Download Flow

```
Frontend → GET /attachments/:id/download
         → AttachmentService:
             1. Validasi ownership
             2. Generate presigned URL (15 menit)
         → Frontend langsung redirect ke MinIO URL
```

Presigned URL digunakan agar file besar tidak melewati backend Go (lebih efisien), dan URL expired otomatis setelah 15 menit sehingga tidak bisa dishare sembarangan.

### 10.4 Supported File Types

| MIME Type | Extension |
|---|---|
| `application/pdf` | .pdf |
| `image/jpeg` | .jpg, .jpeg |
| `image/png` | .png |
| `image/webp` | .webp |
| `application/msword` | .doc |
| `application/vnd.openxmlformats-officedocument.wordprocessingml.document` | .docx |
| `application/vnd.ms-excel` | .xls |
| `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet` | .xlsx |
| `text/csv` | .csv |
| `text/plain` | .txt |
| `application/zip` | .zip (juga menangkap .docx/.xlsx yang tidak terdeteksi spesifik) |

---

## 11. Notification System

### 11.1 Notification Types

| Type | Trigger |
|---|---|
| `REQUEST_SUBMITTED` | Project request disubmit |
| `REQUEST_APPROVED` | Project request diapprove |
| `REQUEST_REJECTED` | Project request ditolak permanen |
| `REVISION_REQUESTED` | Admin meminta revisi project request |
| `REQUEST_REVISED` | Project request direvisi oleh requester |
| `TASK_ASSIGNED` | User di-assign ke task |
| `TASK_COMPLETED` | Task ditandai selesai |
| `TASK_OVERDUE` | Task melewati due date dan belum DONE/CANCELLED |
| `TASK_DUE_SOON` | Task akan jatuh tempo dalam 24 jam dan belum DONE/CANCELLED |
| `TASK_OVERDUE` | Task melewati due date dan belum DONE/CANCELLED |
| `MILESTONE_COMPLETED` | Milestone ditandai selesai |
| `PROJECT_COMPLETED` | Project ditandai selesai |
| `BUDGET_WARNING` | Budget usage ≥ 80% |
| `BUDGET_OVER_LIMIT` | Budget usage ≥ 100% |
| `HANDOVER_CREATED` | Handover dibuat (notify receiver + konfirmasi ke sender) |
| `HANDOVER_SENT` | Handover dibuat (notify receiver) |
| `HANDOVER_RECEIVED` | Handover ditandai diterima (notify sender) |

### 11.2 Opt-Out (Mute per Type)

User bisa mute tipe notifikasi tertentu via `PUT /notifications/preferences`. Default **semua notifikasi aktif** (opt-out model). Jika preference diset `enabled: false`, notifikasi tidak dibuat sama sekali (bukan dibuat tapi disembunyikan).

### 11.3 Polling vs WebSocket

Frontend saat ini menggunakan **polling setiap 30 detik** (`refetchInterval: 30000` di TanStack Query). WebSocket akan diimplementasikan di roadmap v2.0 sesuai SDD section 20.

---

## 12. Reporting

### 12.1 Report Types & Scopes

| Type | System-wide (ADMIN) | Per-Project (PM / ADMIN) |
|---|---|---|
| PROJECT | ✅ | ✅ |
| MILESTONE | ✅ | ✅ |
| TASK | ✅ | ✅ |
| BUDGET | ✅ | ✅ |
| HANDOVER | ✅ | ✅ |

### 12.2 Export Formats

- **PDF** — menggunakan `gofpdf`, layout tabel sederhana dengan header kolom
- **Excel (.xlsx)** — menggunakan `excelize v2`, satu sheet per report type

### 12.3 Generate Report

```bash
# System-wide (ADMIN only)
curl -X POST https://api.ppms.yourcompany.com/api/v1/reports/generate \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"type":"BUDGET","format":"EXCEL"}' \
  --output budget_report.xlsx

# Per-project (PROJECT_MANAGER atau ADMIN)
curl -X POST https://api.ppms.yourcompany.com/api/v1/projects/5/reports/generate \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"type":"TASK","format":"PDF"}' \
  --output task_report.pdf
```

---

## 13. Audit Logging

### 13.1 Yang Dicatat

| Module | Action(s) yang Di-log |
|---|---|
| `auth` | LOGIN_SUCCESS, LOGIN_FAILED, PASSWORD_CHANGED, SESSIONS_REVOKED |
| `user` | CREATE_USER, UPDATE_USER, ASSIGN_ROLE, DEACTIVATE_USER, RESTORE_USER |
| `division` | CREATE_DIVISION, UPDATE_DIVISION, DELETE_DIVISION |
| `project_request` | SUBMIT_REQUEST, REVIEW_REQUEST_APPROVED/REJECTED/REQUEST_REVISION, REVISE_REQUEST |
| `project` | CHANGE_PROJECT_STATUS, CREATE_PROJECT_DIRECT |
| `import_export` | EXPORT_DATA, IMPORT_DATA |
| `budget` (portfolio) | CREATE_BUDGET_YEAR, UPDATE_BUDGET_YEAR, DELETE_BUDGET_YEAR |
| `milestone` | CREATE_MILESTONE, UPDATE_MILESTONE, CHANGE_MILESTONE_STATUS, REORDER_MILESTONES, DELETE_MILESTONE |
| `task` | CREATE_TASK, UPDATE_TASK, CHANGE_TASK_STATUS, UPDATE_TASK_PROGRESS, ASSIGN_TASK_USERS, DELETE_TASK |
| `budget` | CREATE_BUDGET, UPDATE_BUDGET_ALLOCATION, CREATE_TRANSACTION_EXPENSE/REFUND/ADJUSTMENT |
| `attachment` | UPLOAD_FILE, DOWNLOAD_FILE, DELETE_FILE |
| `handover` | CREATE_HANDOVER, MARK_HANDOVER_RECEIVED, MARK_HANDOVER_RETURNED |

### 13.2 Redacted Fields

Field berikut **selalu** di-`[REDACTED]` sebelum disimpan ke `old_data`/`new_data`:
- `password_hash`
- `password`
- `refresh_token_hash`
- `access_token`
- `refresh_token`
- `old_password`
- `new_password`

### 13.3 Immutability

Tabel `audit_logs` tidak memiliki endpoint UPDATE atau DELETE. Bersifat append-only by design. Hanya ADMIN yang bisa membaca via `GET /audit-logs`.

---

## 14. Frontend Architecture

### 14.1 State Management

PPMS tidak menggunakan Redux/Zustand untuk server state. Semua state yang berasal dari API dikelola oleh **TanStack Query** (caching, refetching, optimistic updates). State UI lokal yang kecil (form input, toggle) menggunakan `useState`.

### 14.2 Auth Flow (Frontend)

```
Login → AuthContext.login() → simpan access_token + refresh_token di localStorage
                               + simpan user data di localStorage

Request ke API → axiosInstance menambahkan Bearer token otomatis
              → Jika 401 → axiosInstance interceptor panggil POST /auth/refresh
              → Simpan token baru → retry request awal
              → Jika refresh juga gagal → clear localStorage → redirect ke /login
```

> **Security note:** localStorage rentan XSS. Ini diterima untuk MVP internal. Untuk upgrade ke httpOnly cookie, lakukan di Phase hardening berikutnya jika threat model perusahaan mengharuskan.

### 14.3 Route Guard

```typescript
// Semua route di dalam DashboardLayout dibungkus ProtectedRoute
<Route element={<ProtectedRoute />}>
  <Route element={<DashboardLayout />}>
    ...
    // Route dengan role restriction:
    <Route element={<ProtectedRoute allowedRoles={["ADMIN"]} />}>
      <Route path="/audit-logs" element={<AuditLogPage />} />
    </Route>
  </Route>
</Route>
```

### 14.4 Error Handling (Frontend)

Semua error dari API memiliki struktur `{ success: false, code: string, message: string }`. Frontend memetakan `code` ke pesan Bahasa Indonesia menggunakan `src/lib/errorMessages.ts`. Jika code tidak dikenali, fallback ke `message` dari backend.

---

## 15. Deployment

### 15.1 Development

```bash
docker compose up -d --build
```

Akses services:
- Frontend (Vite dev server): http://localhost:5174
- Backend: http://localhost:8081
- MinIO Console: http://localhost:9001

### 15.2 Production

```bash
# 1. Siapkan .env dengan nilai production yang sesungguhnya
# 2. Build & run dengan override file production
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# 3. Jalankan migration (gunakan sslmode=require di production)
migrate -path ppms-backend/migrations \
  -database "postgres://USER:PASS@HOST:5432/DB?sslmode=require" up
```

Perbedaan utama `docker-compose.prod.yml`:
- Port PostgreSQL tidak diekspos ke publik (hanya accessible via Docker network)
- Frontend di-build static + served via nginx (bukan Vite dev server)
- Resource limits (CPU/memory) diset pada service backend
- Sidecar backup container otomatis jalan

### 15.3 Backup & Recovery

```bash
# Backup manual (jika tidak pakai sidecar container)
chmod +x ppms-backend/scripts/backup.sh
DB_CONTAINER=ppms-postgres BACKUP_DIR=/backups ./ppms-backend/scripts/backup.sh

# Restore dari backup
chmod +x ppms-backend/scripts/restore.sh
./ppms-backend/scripts/restore.sh /backups/ppms_db_20260618_020000.sql.gz
```

Backup otomatis via sidecar container berjalan setiap 24 jam dengan retensi 30 hari. File backup tersimpan di folder `./backups/` di Docker host.

### 15.4 Health Monitoring

| Endpoint | Deskripsi |
|---|---|
| `GET /health` | Selalu 200 jika process jalan |
| `GET /ready` | 200 jika database reachable, 503 jika tidak |

---

## 16. Testing

### 16.1 Manual E2E Testing

Ikuti checklist di `ppms-backend/docs/testing/e2e-checklist.md`. Checklist dipetakan langsung ke acceptance criteria SRS section 6:
1. Functional requirements per modul
2. State machine compliance
3. Permission engine (permission matrix)
4. Audit trail completeness
5. NFR performance & security

### 16.2 Load Testing

```bash
# Install k6
brew install k6  # macOS
sudo apt-get install k6  # Debian/Ubuntu

# Jalankan test
k6 run ppms-backend/loadtest/scenario.js

# Dengan custom environment
k6 run \
  -e BASE_URL=http://localhost:8081/api/v1 \
  -e TEST_EMAIL=john@ppms.local \
  -e TEST_PASSWORD=Password123 \
  ppms-backend/loadtest/scenario.js
```

**Threshold yang harus lolos (NFR-01):**
- `http_req_duration` p95 < 3000ms
- `http_req_duration{endpoint:dashboard}` p95 < 5000ms
- `http_req_failed` < 1%

### 16.3 Soft Delete Testing

Ikuti checklist di `ppms-backend/docs/testing/soft-delete-checklist.md` sebelum go-live.

---

## 17. Known Limitations & Roadmap

### 17.1 Known Limitations (MVP v1.0)

| Item | Status | Catatan |
|---|---|---|
| Notification polling (30s interval) | MVP — bukan real-time | WebSocket di roadmap v2.0 |
| Rate limiter in-memory (single instance only) | MVP | Ganti Redis-based jika scale ke multi-instance |
| PROJECT_REQUEST attachment ownership | Dilonggarkan | Validasi hanya "harus login"; full ownership check via `RequestOwnershipAdapter` sudah ada tapi perlu testing lebih lanjut |
| Restore endpoint hanya untuk User | MVP | Division, Milestone, Task, Attachment belum punya restore endpoint |
| Token disimpan di localStorage | MVP | Pertimbangkan httpOnly cookie untuk hardening XSS |
| No automated test suite | MVP | E2E checklist manual; automated test menjadi inisiatif terpisah pasca-MVP |
| Reporting MILESTONE & HANDOVER belum ditest di production | MVP | SQL query sudah ada, belum ada test data representative |

### 17.2 Roadmap

**v1.5 (Near-term)**
- Email notification (channel sudah disiapkan di tabel `notifications.channel`)
- Scheduled reporting (PDF/Excel dikirim otomatis via email per periode)
- Restore endpoints untuk semua entity (Division, Milestone, Task, Attachment)
- Redis-based rate limiter (untuk multi-instance deployment)

**v2.0 (Medium-term)**
- WebSocket notification (menggantikan polling 30 detik)
- KPI Analytics dashboard (trend analysis, burndown chart)
- Automated test suite (Go `httptest` untuk integration tests)
- httpOnly cookie untuk token storage

**v3.0 (Long-term)**
- Historical Project Snapshots
- Multi-Level Approval Workflow
- External Client Portal
- Mobile Application
- Integration dengan ERP eksternal
- AI Recommendation System

---

## Appendix: Useful Commands

```bash
# Masuk ke container database
docker exec -it ppms-postgres psql -U ppms_user -d ppms_db

# Lihat semua migration yang sudah dijalankan
migrate -path ppms-backend/migrations \
  -database "postgres://ppms_user:change_me@localhost:5432/ppms_db?sslmode=disable" version

# Tidy Go modules
cd ppms-backend && go mod tidy

# Build backend binary secara lokal
cd ppms-backend && go build -o ppms-api ./cmd/api

# Rebuild hanya backend container
docker compose up -d --build backend

# Lihat log backend secara realtime
docker logs -f ppms-backend

# Lihat log backup container
docker logs ppms-backup

# Reset semua (hati-hati: menghapus semua data)
docker compose down -v
```

---
