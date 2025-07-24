# Frontend Integration Guide

## Alur Bisnis (User Flow)

1. **User Login/Register**
   - Endpoint: `POST /api/auth/login`, `POST /api/auth/register`
   - Data: email, password (dan data user lain saat register)
   - Response: JWT token, user info

2. **Lihat Daftar Film & Jadwal**
   - Endpoint: `GET /api/films`, `GET /api/jadwals/film/:filmId`, `GET /api/jadwals/:jadwal_id/kursi-kosong`
   - Data: -
   - Response: List film, jadwal, kursi kosong

3. **Booking Tiket**
   - Endpoint: `POST /api/tikets`
   - Data yang dikirim:
     ```json
     {
       "jadwal_id": "string",
       "kursi": ["A1", "A2"],
       "user_id": "string"
     }
     ```
   - Response: Tiket info (status: belum dibayar)

4. **Lihat Tiket**
   - Endpoint: `GET /api/tikets/me`
   - Data: -
   - Response: List tiket user

5. **Pembayaran**
   - Endpoint: `POST /api/pembayarans`
   - Data yang dikirim:
     ```json
     {
       "tiket_id": "string",
       "user_id": "string",
       "metode_pembayaran": "qris|ovo|dana|dll",
       "jumlah": 50000
     }
     ```
   - Response:
     ```json
     {
       "message": "Terima kasih telah melakukan pembayaran!",
       "status": "paid",
       "pembayaran": { ... }
     }
     ```

6. **(Opsional) Upload Receipt**
   - Frontend generate file JPG (struk digital) dan upload ke Supabase Storage langsung dari frontend.
   - Simpan URL receipt di frontend atau update ke backend via `PUT /api/pembayarans/:id` jika ingin.

7. **Lihat Status Pembayaran**
   - Endpoint: `GET /api/pembayarans/user/:user_id`
   - Response: List pembayaran user, status, dan (opsional) URL receipt.

---

## Kebutuhan Visual Frontend

- **Login/Register Page**: Form email, password, register data.
- **Film List Page**: List film, klik untuk detail/jadwal.
- **Jadwal & Kursi Page**: Pilih jadwal, tampilkan kursi kosong, pilih kursi.
- **Tiket Saya Page**: List tiket user, status pembayaran.
- **Pembayaran Page**: Form konfirmasi pembayaran, tombol "Bayar".
- **(Opsional) Receipt Page**: Tampilkan struk digital (JPG) jika sudah di-upload.
- **Profile Page**: Lihat & edit data user.

---

## Credentials & Konfigurasi

- **JWT Token**: Diperlukan untuk semua endpoint yang butuh autentikasi. Simpan di localStorage/cookie dan kirim di header `Authorization: Bearer <token>`.
- **Supabase Storage**: Untuk upload receipt JPG, frontend perlu akses ke bucket Supabase (gunakan public bucket atau token yang sesuai, tergantung konfigurasi Supabase).
- **API Base URL**: Sesuaikan dengan alamat backend (misal: `http://localhost:8080/api`).

---

## Catatan Penting
- Tidak ada upload bukti pembayaran ke backend, hanya ke Supabase Storage (opsional).
- Semua proses pembayaran dianggap sukses setelah user menekan tombol "Bayar".
- Jika ingin menyimpan URL receipt ke backend, gunakan endpoint update pembayaran.
- Pastikan validasi data di frontend sebelum request ke backend.

---

Silakan sesuaikan kebutuhan visual dan data sesuai flow di atas. Jika ada perubahan API, cek dokumentasi backend terbaru.
