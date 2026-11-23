# Booking Room API

RESTful API untuk pemesanan ruang rapat menggunakan Go (Gin) dan PostgreSQL dengan pencegahan *double booking* memakai pessimistic locking.

## Struktur Proyek

```text
├── cmd/booking-room/main.go
├── internal/
│   ├── config/
│   ├── domain/
│   ├── handler/
│   ├── middleware/
│   ├── repository/
│   └── usecase/
├── migration/
│   └── schema.sql
├── tests/
│   └── concurrency_test.go
├── Dockerfile
└── docker-compose.yml
```

## Menjalankan Aplikasi

Pastikan Docker dan Docker Compose tersedia, lalu jalankan:

```bash
docker-compose up -d --build
```

Aplikasi berjalan di `http://localhost:2254`. Saat pertama kali menyala, PostgreSQL otomatis menjalankan `schema.sql` beserta akun seeder:

- `admin@test.com` / `password123` (`ADMIN`)
- `user1@test.com` / `password123` (`USER`)
- `user2@test.com` / `password123` (`USER`)
- `user3@test.com` / `password123` (`USER`)
- `user4@test.com` / `password123` (`USER`)
- `user5@test.com` / `password123` (`USER`)
- `user6@test.com` / `password123` (`USER`)
- `user7@test.com` / `password123` (`USER`)
- `user8@test.com` / `password123` (`USER`)
- `user9@test.com` / `password123` (`USER`)
- `user10@test.com` / `password123` (`USER`)

## Menjalankan Automated Test

Pastikan container sudah berjalan (`http://localhost:2254`), lalu jalankan:

```bash
go test -v ./tests/...
```

## Contoh Pengujian Endpoint (`curl`)

### 1. Login

```bash
curl -X POST http://localhost:2254/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"password123"}'
```

### 2. Tambah Ruang Meeting (ADMIN)

```bash
curl -X POST http://localhost:2254/api/v1/rooms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN_ADMIN>" \
  -d '{"name":"Ruang Garuda","capacity":10}'
```

### 3. Cek Ketersediaan Ruangan

```bash
curl -X GET "http://localhost:2254/api/v1/rooms/available?start=2026-12-01T10:00:00Z&end=2026-12-01T12:00:00Z" \
  -H "Authorization: Bearer <TOKEN_USER>"
```

### 4. Buat Pemesanan Ruangan

```bash
curl -X POST http://localhost:2254/api/v1/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN_USER>" \
  -d '{
    "room_id": 1,
    "title": "Sprint Planning",
    "start_time": "2026-12-01T10:00:00Z",
    "end_time": "2026-12-01T12:00:00Z"
  }'
```

### 5. Batalkan Pemesanan

```bash
curl -X PATCH http://localhost:2254/api/v1/bookings/1/cancel \
  -H "Authorization: Bearer <TOKEN_USER>"
```
