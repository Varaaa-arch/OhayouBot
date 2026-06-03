# Database Layer — OhayouBot

Dokumentasi ini menjelaskan implementasi database layer pada OhayouBot API, mencakup migrasi, konfigurasi SQLC, dan repository pattern yang digunakan.

---

## Daftar Isi

1. [Gambaran Umum](#gambaran-umum)
2. [Teknologi yang Digunakan](#teknologi-yang-digunakan)
3. [Struktur Folder](#struktur-folder)
4. [Part 1 — Migrasi Database](#part-1--migrasi-database)
5. [Part 2 — SQLC Query Generator](#part-2--sqlc-query-generator)
6. [Part 3 — Repository Pattern](#part-3--repository-pattern)
7. [Koneksi Database](#koneksi-database)
8. [Cara Menjalankan](#cara-menjalankan)

---

## Gambaran Umum

Database layer OhayouBot menggunakan pendekatan **code generation** lewat [SQLC](https://sqlc.dev/). Alih-alih menulis query manual atau menggunakan ORM berat, kita tulis SQL murni lalu SQLC otomatis generate-kan kode Go yang type-safe.

Alurnya seperti ini:

```
SQL files (sql/*.sql)
        ↓
   sqlc generate
        ↓
Generated Go code (generated/db/)
        ↓
  Repository layer
  (internal/*/repository.go)
```

---

## Teknologi yang Digunakan

| Teknologi | Kegunaan |
|---|---|
| **PostgreSQL** | Database utama |
| **pgx/v5** | Driver PostgreSQL untuk Go (via `database/sql` adapter) |
| **golang-migrate** | Manajemen migrasi skema database |
| **SQLC** | Generate kode Go dari SQL query |
| **google/uuid** | Tipe data UUID untuk primary key |

---

## Struktur Folder

```
apps/api/
├── migrations/              # File migrasi SQL (up & down)
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_contacts.up.sql
│   ├── 000002_create_contacts.down.sql
│   ├── 000003_create_schedules.up.sql
│   ├── 000003_create_schedules.down.sql
│   ├── 000004_create_message_logs.up.sql
│   ├── 000004_create_message_logs.down.sql
│   └── 000005_create_whatsapp_sessions.up.sql
│
├── sql/                     # Query SQL untuk SQLC
│   ├── contacts.sql
│   ├── schedules.sql
│   ├── users.sql
│   ├── message_logs.sql
│   └── whatsapp_sessions.sql
│
├── generated/db/            # Kode Go hasil generate SQLC (jangan diedit manual)
│   ├── db.go
│   ├── models.go
│   ├── contacts.sql.go
│   ├── schedules.sql.go
│   ├── users.sql.go
│   ├── message_logs.sql.go
│   └── whatsapp_sessions.sql.go
│
├── configs/
│   └── database.go          # Fungsi koneksi database
│
├── internal/
│   ├── contact/repository.go
│   ├── schedule/repository.go
│   └── message/repository.go
│
└── sqlc.yaml                # Konfigurasi SQLC
```

---

## Part 1 — Migrasi Database

Migrasi dikelola dengan **golang-migrate**. Setiap perubahan skema punya dua file: `up` (terapkan) dan `down` (batalkan/rollback).

### Tabel-tabel yang Dibuat

#### `users`
Menyimpan data pengguna yang login ke sistem.

```sql
CREATE TABLE users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
```

#### `contacts`
Menyimpan daftar kontak WhatsApp yang akan dikirim pesan.

```sql
CREATE TABLE contacts (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    phone      TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
```

#### `schedules`
Menyimpan jadwal pengiriman pesan untuk tiap kontak.

```sql
CREATE TABLE schedules (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_id UUID REFERENCES contacts(id) ON DELETE CASCADE,
    message    TEXT NOT NULL,
    time       TEXT NOT NULL,       -- format: "07:00"
    enabled    BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now()
);
```

> `ON DELETE CASCADE` → kalau kontak dihapus, semua jadwalnya ikut terhapus otomatis.

#### `message_logs`
Mencatat riwayat pengiriman pesan.

```sql
CREATE TABLE message_logs (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_id UUID REFERENCES contacts(id) ON DELETE CASCADE,
    message    TEXT NOT NULL,
    status     TEXT NOT NULL,       -- contoh: "sent", "failed"
    sent_at    TIMESTAMPTZ DEFAULT now()
);
```

#### `whatsapp_sessions`
Menyimpan data sesi WhatsApp agar tidak perlu scan QR ulang setiap restart.

```sql
CREATE TABLE whatsapp_sessions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_data TEXT NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT now()
);
```

### Konvensi Penamaan File Migrasi

Format: `{nomor urut}_{nama_tabel}.{up|down}.sql`

- Nomor urut dimulai dari `000001`, naik satu per satu
- Urutan penting — tabel yang di-*reference* harus dibuat lebih dulu (misal: `contacts` harus ada sebelum `schedules`)

---

## Part 2 — SQLC Query Generator

### Konfigurasi (`sqlc.yaml`)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "sql/"          # lokasi file .sql berisi query
    schema: "migrations/"    # lokasi file migrasi sebagai schema
    gen:
      go:
        package: "db"
        out: "generated/db"  # output kode Go
```

### Query yang Tersedia

#### `sql/contacts.sql`

| Query | Fungsi |
|---|---|
| `CreateContact` | Insert kontak baru, return semua field |
| `GetContacts` | Ambil semua kontak, urut terbaru dulu |

#### `sql/schedules.sql`

| Query | Fungsi |
|---|---|
| `CreateSchedule` | Insert jadwal baru, return semua field |
| `GetSchedules` | Ambil semua jadwal, urut terbaru dulu |

#### `sql/users.sql`

| Query | Fungsi |
|---|---|
| `CreateUser` | Insert user baru, return semua field |

#### `sql/message_logs.sql`

| Query | Fungsi |
|---|---|
| `CreateMessageLog` | Insert log pesan baru, return semua field |
| `GetMessageLogsByContact` | Ambil log berdasarkan `contact_id`, urut terbaru dulu |

#### `sql/whatsapp_sessions.sql`

| Query | Fungsi |
|---|---|
| `CreateSession` | Insert sesi baru, return semua field |
| `GetLatestSession` | Ambil sesi paling baru (limit 1) |

### Hasil Generate

SQLC menghasilkan beberapa file di `generated/db/`:

- **`models.go`** — struct Go untuk setiap tabel (misal: `Contact`, `Schedule`, `MessageLog`)
- **`db.go`** — interface `DBTX` dan struct `Queries` sebagai entry point
- **`*.sql.go`** — implementasi fungsi query per domain

Contoh struct hasil generate:

```go
// models.go
type Contact struct {
    ID        uuid.UUID
    Name      string
    Phone     string
    CreatedAt sql.NullTime
}

type Schedule struct {
    ID        uuid.UUID
    ContactID uuid.NullUUID
    Message   string
    Time      string
    Enabled   sql.NullBool
    CreatedAt sql.NullTime
}
```

> **Jangan edit file di `generated/db/` secara manual.** Semua perubahan harus lewat file `.sql` lalu jalankan ulang `sqlc generate`.

---

## Part 3 — Repository Pattern

Setiap domain punya `repository.go` sendiri yang membungkus SQLC queries. Ini memisahkan logika akses data dari logika bisnis di layer service.

### Contact Repository

```go
// internal/contact/repository.go
type Repository struct {
    q *db.Queries
}

func NewRepository(q *db.Queries) *Repository
func (r *Repository) Create(ctx, db.CreateContactParams) (db.Contact, error)
func (r *Repository) List(ctx) ([]db.Contact, error)
```

### Schedule Repository

```go
// internal/schedule/repository.go
type Repository struct {
    q *db.Queries
}

func NewRepository(q *db.Queries) *Repository
func (r *Repository) Create(ctx, db.CreateScheduleParams) (db.Schedule, error)
func (r *Repository) List(ctx) ([]db.Schedule, error)
```

### Message Repository

```go
// internal/message/repository.go
type Repository struct {
    q *db.Queries
}

func NewRepository(q *db.Queries) *Repository
func (r *Repository) Create(ctx, db.CreateMessageLogParams) (db.MessageLog, error)
func (r *Repository) ListByContact(ctx, uuid.NullUUID) ([]db.MessageLog, error)
```

### Cara Pakai di Service

```go
// Contoh inisialisasi di main.go atau wire
queries, _ := configs.ConnectDB(os.Getenv("DATABASE_URL"))

contactRepo  := contact.NewRepository(queries)
scheduleRepo := schedule.NewRepository(queries)
messageRepo  := message.NewRepository(queries)
```

---

## Koneksi Database

File `configs/database.go` bertanggung jawab membuka koneksi dan menginisialisasi `*db.Queries`.

```go
func ConnectDB(dsn string) (*db.Queries, error) {
    conn, err := sql.Open("pgx", dsn)  // pgx/v5 via stdlib adapter
    if err != nil {
        return nil, err
    }
    return db.New(conn), nil
}
```

**Kenapa pakai `pgx/v5/stdlib` bukan native pgx?**

SQLC generate interface `DBTX` yang berbasis `database/sql` (`QueryRowContext`, `QueryContext`, dll). Paket `pgx/v5/stdlib` adalah adapter yang membungkus pgx agar kompatibel dengan `database/sql`, sehingga kita tetap dapat performa pgx tanpa harus ganti interface.

Format DSN:

```
postgres://user:password@host:port/dbname?sslmode=disable
```

---

## Cara Menjalankan

### 1. Set environment variable

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/ohayoubot?sslmode=disable"
```

### 2. Jalankan migrasi

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

Untuk rollback satu langkah:

```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

### 3. Generate ulang kode SQLC (setelah ada perubahan query/skema)

```bash
sqlc generate
```

> Jalankan dari direktori `apps/api/`.

### 4. Install tools (jika belum ada)

```bash
# golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
