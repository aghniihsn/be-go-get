# GoGetCinema - Cinema Booking API Backend

GoGetCinema adalah sistem backend RESTful API untuk layanan pemesanan tiket bioskop online. Dibangun dengan Go dan Fiber framework, menggunakan MongoDB sebagai basis data, dan JWT untuk autentikasi.

## 📋 Daftar Isi
- [Fitur Utama](#fitur-utama)
- [Alur Bisnis](#alur-bisnis)
- [Struktur Proyek](#struktur-proyek)
- [Spesifikasi API](#spesifikasi-api)
- [Teknologi](#teknologi)

## 🚀 Fitur Utama

### 🔐 Autentikasi & Manajemen Pengguna
- Registrasi dan login pengguna
- Autentikasi berbasis JWT
- Manajemen profil pengguna
- Kontrol akses berbasis peran (user/admin)

### 🎬 Manajemen Film
- Pengelolaan data film (CRUD)
- Kategorisasi film berdasarkan genre
- Sistem rating film (Semua Umur, Anak-anak, Remaja, Dewasa)
- Informasi film lengkap dengan deskripsi dan poster

### 📅 Manajemen Jadwal
- Pengelolaan jadwal tayang (CRUD)
- Filter jadwal berdasarkan film, tanggal, dan waktu
- Penempatan studio/ruangan
- Konfigurasi harga untuk setiap jadwal

### 🎟️ Sistem Tiket
- Pemesanan tiket individual
- Pemesanan tiket batch untuk multiple kursi
- Pengecekan ketersediaan kursi
- Manajemen status tiket (confirmed, cancelled, used)
- Informasi tiket detail dengan data film dan jadwal

### 💰 Sistem Pembayaran
- Dukungan berbagai metode pembayaran
- Pelacakan status pembayaran (pending, completed, failed)
- Unggah bukti pembayaran
- Sistem konfirmasi pembayaran
- Update status tiket otomatis setelah pembayaran berhasil

## 🔄 Alur Bisnis

### Alur Pengguna Reguler:
1. **Pendaftaran & Login**
   - Pengguna mendaftar dengan username, email, dan password
   - Login untuk mendapatkan token JWT

2. **Penelusuran Film & Jadwal**
   - Lihat daftar film yang tersedia
   - Pilih film dan lihat jadwal tayang
   - Cek ketersediaan kursi untuk jadwal tertentu

3. **Pemesanan Tiket**
   - Pilih satu atau beberapa kursi
   - Buat tiket untuk jadwal yang dipilih
   - Tiket memiliki status "confirmed" setelah dibuat

4. **Pembayaran**
   - Pilih metode pembayaran
   - Lakukan pembayaran dan upload bukti pembayaran
   - Admin memverifikasi pembayaran
   - Status pembayaran diubah menjadi "completed"

5. **Manajemen Tiket**
   - Lihat tiket yang telah dipesan
   - Lihat detail tiket dengan info film dan jadwal
   - Batal tiket jika diperlukan (sebelum digunakan)

### Alur Admin:
1. **Manajemen Film**
   - Tambah, edit, hapus data film
   - Kelola informasi dan kategori film

2. **Manajemen Jadwal**
   - Buat jadwal baru untuk film
   - Atur studio dan harga tiket
   - Edit atau hapus jadwal yang sudah ada

3. **Verifikasi Pembayaran**
   - Tinjau pembayaran yang masuk
   - Verifikasi bukti pembayaran
   - Update status pembayaran

4. **Manajemen Tiket**
   - Lihat semua tiket dalam sistem
   - Filter tiket berdasarkan status atau pengguna
   - Cancel tiket jika diperlukan

## 📁 Struktur Proyek

```
be-go-get/
├── main.go                     # Entry point aplikasi
├── go.mod                      # Dependensi Go module
├── go.sum                      # Checksum Go module
├── Makefile                    # Script build dan deployment
│
├── config/                     # Konfigurasi aplikasi
│   ├── database.go             # Konfigurasi dan koneksi database
│   ├── cors.go                 # Konfigurasi CORS
│   └── middleware/             # Middleware aplikasi
│       ├── auth.go             # Autentikasi JWT
│       └── encoder.go          # Utilitas encoding
│
├── controllers/                # Handler request HTTP
│   ├── authController.go       # Endpoint autentikasi
│   ├── filmController.go       # Endpoint manajemen film
│   ├── jadwalController.go     # Endpoint manajemen jadwal
│   ├── pembayaranController.go # Endpoint pemrosesan pembayaran
│   ├── tiketController.go      # Endpoint pemesanan tiket
│   ├── batchController.go      # Endpoint untuk operasi batch
│   └── userController.go       # Endpoint manajemen pengguna
│
├── models/                     # Struktur data
│   ├── struct.go               # Struktur utama dengan MongoDB ObjectID
│   └── request.go              # Model untuk request API
│
## 📡 Spesifikasi API

### Rute Publik
- `GET /api/films` - Mendapatkan semua film
- `GET /api/films/:id` - Mendapatkan film berdasarkan ID
- `GET /api/jadwals` - Mendapatkan semua jadwal
- `GET /api/jadwals/detail` - Mendapatkan semua jadwal dengan detail film
- `GET /api/jadwals/:id` - Mendapatkan jadwal berdasarkan ID
- `GET /api/jadwals/film/:filmId` - Mendapatkan jadwal berdasarkan ID film
- `GET /api/jadwals/:jadwal_id/kursi-kosong` - Mendapatkan kursi yang tersedia untuk suatu jadwal
- `GET /api/payment-methods` - Mendapatkan metode pembayaran yang tersedia

### Rute Autentikasi
- `POST /api/auth/register` - Registrasi pengguna baru
- `POST /api/auth/login` - Login pengguna
- `GET /api/auth/profile` - Mendapatkan profil pengguna
- `PUT /api/auth/profile` - Update profil pengguna

### Rute Pengguna Terotentikasi
- `GET /api/tikets/me` - Mendapatkan tiket milik pengguna
- `GET /api/tikets/:id/summary` - Mendapatkan informasi detail tiket
- `POST /api/tikets` - Membuat tiket tunggal
- `POST /api/tikets/batch` - Membuat multiple tiket sekaligus
- `PUT /api/tikets/:id` - Update tiket
- `DELETE /api/tikets/:id` - Membatalkan tiket
- `GET /api/pembayarans/:id` - Mendapatkan pembayaran berdasarkan ID
- `GET /api/pembayarans/user/:user_id` - Mendapatkan pembayaran berdasarkan ID pengguna
- `POST /api/pembayarans` - Membuat pembayaran baru
- `PUT /api/pembayarans/:id` - Update status pembayaran

### Rute Admin
- `POST /api/films` - Membuat film baru
- `PUT /api/films/:id` - Update film
- `DELETE /api/films/:id` - Menghapus film
- `POST /api/jadwals` - Membuat jadwal baru
- `PUT /api/jadwals/:id` - Update jadwal
- `DELETE /api/jadwals/:id` - Menghapus jadwal
- `GET /api/tikets` - Mendapatkan semua tiket
- `GET /api/tikets/:id` - Mendapatkan tiket berdasarkan ID
- `GET /api/pembayarans` - Mendapatkan semua pembayaran
- `DELETE /api/pembayarans/:id` - Menghapus pembayaran

## 🔧 Teknologi

- **Bahasa**: [Go](https://golang.org/)
- **Framework Web**: [Fiber](https://gofiber.io/)
- **Database**: [MongoDB](https://www.mongodb.com/)
- **Autentikasi**: [JWT](https://jwt.io/)
- **Dokumentasi API**: [Swagger](https://swagger.io/)

## 🚀 Development

### Building
```bash
# Build the application
go build -o bin/main main.go

# Run with Makefile
make build
make run
```

### Contributing
1. Create a feature branch
2. Write tests for new functionality
3. Ensure all tests pass
4. Update documentation as needed
5. Submit a pull request


- Use `go test -v` for detailed test output
- Check logs for ObjectID validation errors
- Monitor JWT token expiration and refresh
- Validate MongoDB ObjectID format in API requests

### Collections
- **films**: Film information
- **jadwals**: Movie schedules
- **tikets**: Ticket bookings
- **users**: User accounts
- **pembayarans**: Payment records

