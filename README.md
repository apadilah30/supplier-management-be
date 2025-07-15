# Supplier Management API

Ini adalah backend API sederhana yang dibangun dengan Go untuk mengelola data supplier, termasuk informasi detail seperti alamat, kontak, dan grup. Proyek ini juga dilengkapi dengan sistem migrasi database otomatis menggunakan `golang-migrate`.

## Fitur

* **Manajemen Supplier**: Buat dan lihat daftar supplier.
* **Data Lengkap**: Menyimpan data alamat, kontak, dan grup untuk setiap supplier.
* **Transaksi Atomik**: Proses pembuatan supplier baru dijamin aman menggunakan transaksi database.
* **Migrasi Otomatis**: Skema database dikelola melalui file SQL dan diterapkan secara otomatis saat aplikasi dimulai.

## Teknologi

* **Bahasa**: Go
* **Database**: PostgreSQL
* **Migrasi**: `golang-migrate/migrate`
* **Driver DB**: `lib/pq`

## Prasyarat

* Go (versi 1.18 atau lebih baru)
* PostgreSQL

## Instalasi & Setup

1.  **Clone Repositori**
    ```bash
    git clone https://github.com/apadilah30/supplier-management-be
    cd https://github.com/apadilah30/supplier-management-be
    ```

2.  **Buat Database**
    Buat database baru di PostgreSQL dengan nama `suppliers_db`.

3.  **Konfigurasi Database**
    Pastikan DB memiliki credentials yang sama dengan di file `main.go`

4.  **Instalasi Dependencies**
    ```bash
    go mod tidy
    ```

5.  **Jalankan Aplikasi**
    Aplikasi akan secara otomatis menjalankan migrasi database saat pertama kali dimulai.
    ```bash
    go run main.go
    ```
    Server akan berjalan di `http://localhost:8080`.

## API Endpoints

### 1. Membuat Supplier Baru

* **Endpoint**: `POST /suppliers`
* **Deskripsi**: Membuat supplier baru beserta alamat, kontak, dan grupnya.
* **Contoh Request**:
    ```bash
    curl -X POST http://localhost:8080/suppliers \
    -H "Content-Type: application/json" \
    -d '{
      "supplier_name": "PT Setroom Indonesia",
      "nick_name": "Setroom",
      "addresses": [{"name": "Head Office", "address": "Fatmawati Raya St, 123", "is_main": true}],
      "contacts": [{"name": "Albert", "job_position": "CEO", "email": "einstein@gmail.com", "phone": "021123456", "mobile": "0811234567", "is_main": true}],
      "groups": [{"group_name": "Industry", "value": "Manufacture", "is_active": true}]
    }'
    ```

### 2. Melihat Daftar Supplier

* **Endpoint**: `GET /suppliers`
* **Deskripsi**: Mengambil daftar semua supplier yang ada.
* **Contoh Request**:
    ```bash
    curl http://localhost:8080/suppliers
    ```