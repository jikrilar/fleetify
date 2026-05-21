# Fleetify

Fleetify adalah aplikasi internal untuk mengelola laporan maintenance kendaraan. Workflow utama dibuat untuk dua role:

- **SA (Service Advisor)** membuat laporan awal, mengisi estimasi item, lalu menyelesaikan laporan yang sudah disetujui.
- **Approval** meninjau laporan yang masuk dan menyetujui laporan dengan status menunggu approval.

Frontend dibuat seperti single-page app sederhana memakai Vanilla JavaScript dan Bootstrap 5. Backend menyediakan REST API memakai Go Fiber, GORM, dan MySQL.

## Tech Stack

- Backend: Go 1.26.3, Go Fiber v3.2.0, GORM v1.31.1
- Database: MySQL 8.0 dengan tabel InnoDB
- Frontend: HTML, Bootstrap 5.3.8, Vanilla JavaScript
- DevOps: Docker Compose

## Cara Menjalankan

```bash
git clone https://github.com/jikrilar/fleetify.git
cd fleetify
cp .env.example .env
docker-compose up --build
```

Aplikasi berjalan di:

```text
http://localhost:8080
```

Health check:

```text
GET http://localhost:8080/health
```

## Seeder

Seeder berjalan otomatis dari folder `docker/mysql/init` saat container MySQL pertama kali dibuat.

Data awal:

| Jenis | Data |
|---|---|
| User | `sa_user`, `approval_user` |
| Kendaraan | `B 1234 FTY`, `B 5678 FTY`, `B 9012 FTY` |
| Item Master | Engine Oil, Oil Filter, Brake Pad, General Service, Brake Inspection |

Seeder dibuat idempotent dengan `ON DUPLICATE KEY UPDATE`.

## Environment Variables

| Variable | Keterangan |
|---|---|
| `APP_PORT` | Port aplikasi Go, default `8080` |
| `DB_HOST` | Host MySQL dari container app, default `mysql` |
| `DB_PORT` | Port MySQL, default `3306` |
| `DB_USER` | Username database |
| `DB_PASSWORD` | Password database |
| `DB_NAME` | Nama database |
| `DB_ROOT_PASSWORD` | Password root MySQL |
| `WEBHOOK_URL` | URL webhook opsional untuk event approve dan complete |
| `UPLOAD_DIR` | Folder upload/simulasi file |
| `FRONTEND_DIR` | Lokasi frontend di container |

## Akun Testing

| Username | Role | X-User-ID |
|---|---|---:|
| `sa_user` | SA | 1 |
| `approval_user` | APPROVAL | 2 |

Frontend memiliki user switcher. API protected wajib mengirim header:

```http
X-User-ID: 1
```

## Workflow

```text
PENDING_APPROVAL -> APPROVED -> COMPLETED
```

Aturan utama:

- SA hanya bisa membuat laporan dan menyelesaikan laporan yang sudah `APPROVED`.
- Approval hanya bisa menyetujui laporan yang masih `PENDING_APPROVAL`.
- Status tidak bisa mundur atau melompat langsung ke status lain.

## Dokumentasi API

### Health Check

```http
GET /health
```

### Get Testing Users

```http
GET /api/users
```

### Get Vehicles

Role: SA, Approval

```http
GET /api/vehicles
X-User-ID: 1
```

### Get Master Items

Role: SA, Approval

```http
GET /api/master-items
X-User-ID: 1
```

### Create Maintenance Report

Role: SA

```http
POST /api/reports
X-User-ID: 1
Content-Type: application/json
```

```json
{
  "vehicle_id": 1,
  "odometer": 120000,
  "complaint": "Rem berbunyi dan mesin bergetar",
  "initial_photo": "foto-awal.jpg",
  "items": [
    {
      "item_id": 1,
      "quantity": 1
    },
    {
      "item_id": 4,
      "quantity": 1
    }
  ]
}
```

Backend selalu mengatur status menjadi `PENDING_APPROVAL` dan mengambil harga dari `master_items.price` ke `report_items.price_snapshot`.

### Get Reports

Role: SA, Approval

```http
GET /api/reports
X-User-ID: 1
```

Filter status:

```http
GET /api/reports?status=PENDING_APPROVAL
X-User-ID: 2
```

### Get Report Detail

Role: SA, Approval

```http
GET /api/reports/1
X-User-ID: 1
```

### Approve Report

Role: Approval

```http
PATCH /api/reports/1/approve
X-User-ID: 2
```

### Complete Report

Role: SA

```http
PATCH /api/reports/1/complete
X-User-ID: 1
Content-Type: application/json
```

```json
{
  "proof_photo": "foto-bukti.jpg"
}
```

## Keputusan Teknis

- **Transaction saat create report**: header laporan dan detail item disimpan dalam satu transaksi agar tidak ada laporan tanpa item jika insert detail gagal.
- **Price snapshot**: harga item disalin dari master item supaya estimasi laporan lama tidak berubah saat harga master item berubah.
- **RBAC sederhana dengan `X-User-ID`**: cukup untuk kebutuhan testing internal tanpa membuat sistem login penuh.
- **Frontend tanpa `.innerHTML`**: data API dirender memakai `createElement`, `textContent`, `DocumentFragment`, dan `replaceChildren`.
- **Retry koneksi database**: backend menunggu MySQL siap sebelum aplikasi gagal start.

## Bonus

- Export CSV dari halaman riwayat memakai Native JavaScript, `Blob`, `URL.createObjectURL()`, dan elemen `<a>` sementara.
- Webhook async memakai goroutine untuk event `REPORT_APPROVED` dan `REPORT_COMPLETED`. Kegagalan webhook hanya dicatat di log dan tidak menggagalkan update status.

## Skenario Test Manual

1. Jalankan `docker-compose up --build`.
2. Buka `http://localhost:8080`.
3. Pilih `sa_user`, buat laporan dengan minimal dua item.
4. Pilih `approval_user`, buka antrean approval, lalu setujui laporan.
5. Pilih `sa_user`, buka menu selesaikan laporan, isi foto bukti, lalu selesaikan.
6. Buka riwayat dan export CSV.
