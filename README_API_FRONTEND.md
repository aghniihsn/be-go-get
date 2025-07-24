# API & Frontend Integration Guide

Dokumen ini merangkum seluruh endpoint utama dari backend beserta kebutuhan data, response, dan saran tampilan frontend.

---

## AUTHENTICATION

### Register
- **POST /api/auth/register**
- **Body:**
  ```json
  {
    "email": "string",
    "username": "string",
    "password": "string",
    "firstname": "string",
    "lastname": "string",
    "gender": "male|female",
    "phone_number": "string",
    "address": "string"
  }
  ```
- **Response:** User info

### Login
- **POST /api/auth/login**
- **Body:**
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response:** JWT token, user info

### Get Profile
- **GET /api/auth/profile**
- **Header:** Authorization: Bearer <token>
- **Response:** User info

### Update Profile
- **PUT /api/auth/profile**
- **Header:** Authorization: Bearer <token>
- **Body:** multipart/form-data (field yang ingin diubah)
- **Response:** Success message

---

## FILM & JADWAL

### List Film
- **GET /api/films**
- **Response:** List film

### Detail Film
- **GET /api/films/:id**
- **Response:** Film detail

### List Jadwal
- **GET /api/jadwals**
- **Response:** List jadwal

### Jadwal per Film
- **GET /api/jadwals/film/:filmId**
- **Response:** List jadwal film

### Detail Jadwal
- **GET /api/jadwals/:id**
- **Response:** Jadwal detail

### Kursi Kosong
- **GET /api/jadwals/:jadwal_id/kursi-kosong**
- **Response:** List kursi kosong

---

## TIKET

### Booking Tiket
- **POST /api/tikets**
- **Header:** Authorization: Bearer <token>
- **Body:**
  ```json
  {
    "jadwal_id": "string",
    "kursi": ["A1", "A2"]
  }
  ```
- **Response:** Tiket info

### Tiket User (riwayat)
- **GET /api/tikets/me**
- **Header:** Authorization: Bearer <token>
- **Response:** List tiket user

### Tiket User by ID
- **GET /api/tikets/user/:user_id**
- **Header:** Authorization: Bearer <token>
- **Response:** List tiket user

### Update Tiket
- **PUT /api/tikets/:id**
- **Header:** Authorization: Bearer <token>
- **Body:**
  ```json
  {
    "kursi": ["A1", "A2"]
  }
  ```
- **Response:** Tiket info

### Cancel Tiket
- **DELETE /api/tikets/:id**
- **Header:** Authorization: Bearer <token>
- **Response:** Success message

### Tiket Summary
- **GET /api/tikets/:id/summary**
- **Header:** Authorization: Bearer <token>
- **Response:** Tiket detail + summary

---

## PEMBAYARAN

### Buat Pembayaran
- **POST /api/pembayarans**
- **Header:** Authorization: Bearer <token>
- **Body:**
  ```json
  {
    "tiket_id": "string",
    "user_id": "string",
    "metode_pembayaran": "qris|ovo|dana|dll",
    "jumlah": 50000
  }
  ```
- **Response:**
  ```json
  {
    "message": "Terima kasih telah melakukan pembayaran!",
    "status": "paid",
    "pembayaran": { ... }
  }
  ```

### Riwayat Pembayaran User
- **GET /api/pembayarans/user/:user_id**
- **Header:** Authorization: Bearer <token>
- **Response:** List pembayaran user

### Detail Pembayaran
- **GET /api/pembayarans/:id**
- **Header:** Authorization: Bearer <token>
- **Response:** Pembayaran detail

### Update Pembayaran (opsional, update status/receipt)
- **PUT /api/pembayarans/:id**
- **Header:** Authorization: Bearer <token>
- **Body:**
  ```json
  {
    "status": "completed|pending",
    "bukti_pembayaran": "url_jpg"
  }
  ```
- **Response:** Success message

---

## PAYMENT METHODS
- **GET /api/payment-methods**
- **Response:** List metode pembayaran

---

## USER
- **GET /api/users/:id**
- **PUT /api/users/:id**
- **Header:** Authorization: Bearer <token>
- **Body:** Data user yang ingin diubah

---

## KEBUTUHAN VISUAL FRONTEND
- Login/Register Page
- Film List & Detail Page
- Jadwal & Kursi Selection Page
- Tiket Saya (riwayat & detail)
- Pembayaran Page (tombol Bayar, konfirmasi, status)
- Upload & Tampilkan Receipt (JPG)
- Profile Page
- Admin: Manajemen Film, Jadwal, Tiket, Pembayaran (dengan upload image ambil dari file manager)

---

## CREDENTIALS & KONFIGURASI
- JWT Token: simpan di localStorage/cookie, kirim di header Authorization
- Supabase Storage: untuk upload receipt JPG 
- API Base URL: sesuaikan dengan alamat backend
- Gunakan Axios untuk fetch API

---

## CATATAN TAMBAHAN
- Semua endpoint yang butuh autentikasi harus mengirim header Authorization: Bearer <token>
- Tidak ada upload bukti pembayaran ke backend, hanya ke Supabase Storage (opsional)
- Jika ingin menyimpan URL receipt ke backend, gunakan endpoint update pembayaran
- Validasi data di frontend sebelum request ke backend

---

Silakan gunakan dokumen ini sebagai referensi pengembangan frontend dan integrasi API.
