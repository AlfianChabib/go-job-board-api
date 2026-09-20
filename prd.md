# Product Requirements Document (PRD) - Job Board REST API

## 1. Ringkasan Proyek

Proyek ini adalah pembuatan RESTful API untuk platform pencarian kerja (Job Board). Sistem ini memungkinkan dua jenis pengguna: **Recruiter** (yang dapat membuat profil perusahaan dan memposting lowongan) dan **Candidate** (yang dapat mencari dan melamar lowongan tersebut).

Proyek ini dirancang sebagai sarana latihan komprehensif untuk mengimplementasikan ekosistem Golang modern dengan arsitektur yang bersih (Clean Architecture).

## 2. Tech Stack & Peran

- **Golang**: Bahasa pemrograman utama.
- **GoFiber**: Web framework untuk routing, middleware, dan penanganan HTTP request/response.
- **GORM**: Object-Relational Mapping (ORM) untuk interaksi database (PostgreSQL/MySQL).
- **Viper**: Manajemen konfigurasi (membaca `.env`, `.yaml`, atau `.json`).
- **Go-Playground/Validator**: Validasi payload/struktur data dari request body.
- **Google Wire**: Dependency Injection (DI) saat compile-time untuk merangkai layer aplikasi.

## 3. Entitas dan Skema Database (GORM)

### 3.1. User

Menyimpan data utama pengguna dan perannya.

- `id` (UUID/Int, Primary Key)
- `name` (String)
- `email` (String, Unique)
- `role` (Enum/String: 'RECRUITER' atau 'CANDIDATE')
- `created_at`, `updated_at`

### 3.2. Auth

Menyimpan data credential / otentikasi kata sandi pengguna (terpisah dari data profil user).

- `id` (UUID/Int, Primary Key)
- `user_id` (Foreign Key ke User, Unique)
- `password` (String, Hashed)
- `created_at`, `updated_at`

### 3.3. Token

Menyimpan data Refresh Token untuk manajemen sesi autentikasi pengguna.

- `id` (UUID/Int, Primary Key)
- `user_id` (Foreign Key ke User)
- `refresh_token` (Text, Unique)
- `expires_at` (Timestamp)
- `is_revoked` (Boolean, default false)
- `created_at`, `updated_at`

### 3.4. Profile

Menyimpan detail profil tambahan untuk pengguna dengan role `CANDIDATE`.

- `id` (UUID/Int, Primary Key)
- `user_id` (Foreign Key ke User, Unique)
- `avatar_url` (String, Nullable) - Foto profil pengguna
- `headline` (String)
- `bio` (Text)
- `phone` (String)
- `resume_url` (String)
- `created_at`, `updated_at`

### 3.5. Skill (Master Data)

Menyimpan master data keahlian/skills yang dapat dipilih oleh pengguna di frontend atau ditambahkan secara manual jika belum ada.

- `id` (UUID/Int, Primary Key)
- `name` (String, Unique) - Nama skill (contoh: "Golang", "PostgreSQL", "React")
- `created_at`, `updated_at`

### 3.6. ProfileSkill (Tabel Penghubung / Many-to-Many)

Tabel perantara relasi Many-to-Many antara `Profile` dan `Skill`.

- `profile_id` (Foreign Key ke Profile)
- `skill_id` (Foreign Key ke Skill)
- Primary Key Composite (`profile_id`, `skill_id`)

### 3.7. Experience (Pengalaman Kerja)

Menyimpan riwayat pengalaman kerja Candidate (Relasi One-to-Many dari `Profile` ke `Experience`).

- `id` (UUID/Int, Primary Key)
- `profile_id` (Foreign Key ke Profile)
- `company_name` (String)
- `position` (String)
- `start_date` (Date/Timestamp)
- `end_date` (Date/Timestamp, Nullable - kosong jika masih bekerja)
- `is_current` (Boolean, default false)
- `description` (Text, Nullable)
- `created_at`, `updated_at`

### 3.8. Company

Menyimpan profil perusahaan (Dikelola oleh Recruiter).

- `id` (UUID/Int, Primary Key)
- `recruiter_id` (Foreign Key ke User)
- `name` (String)
- `logo_url` (String, Nullable) - Visualisasi logo di kartu lowongan.
- `banner_url` (String, Nullable) - Visualisasi cover/banner pada header profil perusahaan.
- `website` (String, Nullable) - Tautan kredibilitas perusahaan.
- `industry` (String) - Kategori sektor (contoh: "Information Technology", "Finance").
- `employee_size` (Enum/String: '1-50', '51-200', '201-500', '500+') - Indikator skala perusahaan.
- `description` (Text)
- `location` (String)
- `created_at`, `updated_at`

### 3.9. Job

Menyimpan detail lowongan pekerjaan.

- `id` (UUID/Int, Primary Key)
- `company_id` (Foreign Key ke Company)
- `title` (String) - Posisi/judul lowongan
- `description` (Text) - Deskripsi pekerjaan & tanggung jawab
- `requirements` (Text, Nullable) - Kualifikasi yang dibutuhkan
- `employment_type` (Enum/String: 'FULL_TIME', 'PART_TIME', 'CONTRACT', 'INTERNSHIP', 'FREELANCE')
- `work_mode` (Enum/String: 'ON_SITE', 'HYBRID', 'REMOTE')
- `location` (String) - Lokasi penempatan kerja
- `min_salary` (BigInt/Int64, Nullable) - Batas bawah gaji
- `max_salary` (BigInt/Int64, Nullable) - Batas atas gaji
- `currency` (String, default: 'IDR')
- `is_salary_negotiable` (Boolean, default: false)
- `status` (Enum/String: 'OPEN', 'CLOSED')
- `created_at`, `updated_at`

### 3.10. Application

Tabel transaksional untuk mencatat lamaran (Candidate melamar Job).

- `id` (UUID/Int, Primary Key)
- `job_id` (Foreign Key ke Job)
- `candidate_id` (Foreign Key ke User)
- `resume_url` (String)
- `status` (Enum/String: 'PENDING', 'REVIEWED', 'ACCEPTED', 'REJECTED')
- `created_at`, `updated_at`

## 4. Spesifikasi API Endpoint

### 4.1. Authentication (Public & Authenticated)

| Method | Endpoint             | Keterangan                                                                          | Validasi (Validator)                                                                                                 |
| ------ | -------------------- | ----------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| POST   | `/api/auth/register` | Mendaftarkan user baru (menyimpan ke `users` & `auths`).                            | `email: required, email`, `password: required, min=6`, `role: required, oneof=RECRUITER CANDIDATE`, `name: required` |
| POST   | `/api/auth/login`    | Login user, mengembalikan access token JWT & refresh token (menyimpan ke `tokens`). | `email: required, email`, `password: required`                                                                       |
| POST   | `/api/auth/refresh`  | Memperbarui access token dengan mengirimkan refresh token yang valid.               | `refresh_token: required`                                                                                            |
| POST   | `/api/auth/logout`   | Logout user dan me-revoke refresh token di tabel `tokens`.                          | `refresh_token: required`                                                                                            |

### 4.2. Master Data Skills (Public & Authenticated)

| Method | Endpoint      | Keterangan                                                                                        | Validasi (Validator) |
| ------ | ------------- | ------------------------------------------------------------------------------------------------- | -------------------- |
| GET    | `/api/skills` | Mendapatkan daftar seluruh master skills yang tersedia (untuk dropdown/autocomplete di frontend). | -                    |

### 4.3. Candidate Profile, Skills & Experience (Candidate Only)

| Method | Endpoint                         | Keterangan                                                                                     | Validasi (Validator)                                                             |
| ------ | -------------------------------- | ---------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| GET    | `/api/candidate/profile`         | Melihat profil pribadi Candidate (termasuk foto profil, skills, dan riwayat pengalaman kerja). | -                                                                                |
| PUT    | `/api/candidate/profile`         | Membuat / memperbarui detail utama profil Candidate.                                           | `headline: required`, `phone: required`                                          |
| PATCH  | `/api/candidate/avatar`          | Memperbarui foto profil Candidate (`avatar_url`).                                              | `avatar_url: required, url`                                                      |
| PUT    | `/api/candidate/skills`          | Memperbarui daftar skills Candidate (bisa memilih skill yang ada atau otomatis menambah baru). | `skills: required, dive, min=1` (Array of string)                                |
| GET    | `/api/candidate/experiences`     | Melihat daftar riwayat pengalaman kerja Candidate.                                             | -                                                                                |
| POST   | `/api/candidate/experiences`     | Menambahkan pengalaman kerja baru.                                                             | `company_name: required`, `position: required`, `start_date: required, datetime` |
| PUT    | `/api/candidate/experiences/:id` | Memperbarui data pengalaman kerja tertentu.                                                    | `company_name: required`, `position: required`, `start_date: required, datetime` |
| DELETE | `/api/candidate/experiences/:id` | Menghapus data pengalaman kerja tertentu.                                                      | -                                                                                |

### 4.4. Company Management (Recruiter & Public)

| Method | Endpoint                        | Keterangan                                                                     | Validasi (Validator)                                                                                              |
| ------ | ------------------------------- | ------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------- |
| POST   | `/api/recruiter/company`        | Membuat profil perusahaan milik recruiter.                                     | `name: required`, `location: required`, `industry: required`, `employee_size: required`                           |
| GET    | `/api/recruiter/company`        | Melihat profil perusahaan milik recruiter yang sedang login.                   | -                                                                                                                 |
| PUT    | `/api/recruiter/company`        | Memperbarui profil perusahaan recruiter.                                       | `name: required`, `location: required`, `industry: required`, `employee_size: required`, `website: omitempty,url` |
| PUT    | `/api/recruiter/company/logo`   | Upload/Replace logo perusahaan (`multipart/form-data`).                        | `file: required, mime=image/png image/jpeg image/webp, max=2MB`                                                   |
| PUT    | `/api/recruiter/company/banner` | Upload/Replace banner perusahaan (`multipart/form-data`).                      | `file: required, mime=image/png image/jpeg image/webp, max=5MB`                                                   |
| GET    | `/api/companies`                | Menampilkan direktori seluruh perusahaan (Public/Pagination).                  | -                                                                                                                 |
| GET    | `/api/companies/:id`            | Melihat detail profil perusahaan spesifik beserta daftar lowongannya (Public). | -                                                                                                                 |

### 4.5. Job Management (Recruiter & Public)

| Method | Endpoint               | Keterangan                                                                                                                                           | Validasi (Validator)                                                                                                                                                                                                                                                                                                                                                                                                   |
| ------ | ---------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| POST   | `/api/jobs`            | Membuat lowongan baru (Recruiter Only).                                                                                                              | `company_id: required`, `title: required`, `description: required`, `requirements: omitempty`, `employment_type: required, oneof=FULL_TIME PART_TIME CONTRACT INTERNSHIP FREELANCE`, `work_mode: required, oneof=ON_SITE HYBRID REMOTE`, `location: required`, `min_salary: omitempty,numeric,min=0`, `max_salary: omitempty,numeric,gtefield=MinSalary`, `currency: omitempty,len=3`, `is_salary_negotiable: boolean` |
| GET    | `/api/jobs`            | Melihat semua lowongan aktif (`OPEN`) dengan pagination & filter query: `search`, `employment_type`, `work_mode`, `location`, `min_salary` (Public). | Query: `page: omitempty,numeric,min=1`, `limit: omitempty,numeric,min=1,max=100`, `employment_type: omitempty,oneof=FULL_TIME PART_TIME CONTRACT INTERNSHIP FREELANCE`, `work_mode: omitempty,oneof=ON_SITE HYBRID REMOTE`, `min_salary: omitempty,numeric`                                                                                                                                                            |
| GET    | `/api/jobs/my-jobs`    | Melihat seluruh daftar lowongan milik recruiter yang sedang login, baik status `OPEN` maupun `CLOSED` (Recruiter Only).                              | Query: `page: omitempty,numeric,min=1`, `limit: omitempty,numeric,min=1,max=100`, `status: omitempty,oneof=OPEN CLOSED`                                                                                                                                                                                                                                                                                                |
| GET    | `/api/jobs/:id`        | Melihat detail lengkap lowongan beserta informasi perusahaan (Public).                                                                               | -                                                                                                                                                                                                                                                                                                                                                                                                                      |
| PUT    | `/api/jobs/:id`        | Memperbarui informasi lowongan kerja (Recruiter Only, wajib verifikasi kepemilikan/IDOR).                                                            | `title: required`, `description: required`, `requirements: omitempty`, `employment_type: required, oneof=FULL_TIME PART_TIME CONTRACT INTERNSHIP FREELANCE`, `work_mode: required, oneof=ON_SITE HYBRID REMOTE`, `location: required`, `min_salary: omitempty,numeric,min=0`, `max_salary: omitempty,numeric,gtefield=MinSalary`, `currency: omitempty,len=3`, `is_salary_negotiable: boolean`                         |
| PATCH  | `/api/jobs/:id/status` | Mengubah status ketersediaan lowongan (`OPEN` / `CLOSED`) (Recruiter Only, wajib verifikasi kepemilikan/IDOR).                                       | `status: required, oneof=OPEN CLOSED`                                                                                                                                                                                                                                                                                                                                                                                  |
| DELETE | `/api/jobs/:id`        | Menghapus lowongan kerja (Recruiter Only, wajib verifikasi kepemilikan/IDOR).                                                                        | -                                                                                                                                                                                                                                                                                                                                                                                                                      |

### 4.6. Application Management (Candidate & Recruiter)

| Method | Endpoint                       | Keterangan                          | Validasi (Validator)                                 |
| ------ | ------------------------------ | ----------------------------------- | ---------------------------------------------------- |
| POST   | `/api/jobs/:id/apply`          | Melamar pekerjaan (Candidate).      | `resume_url: required, url`                          |
| GET    | `/api/jobs/:id/applications`   | Melihat daftar pelamar (Recruiter). | -                                                    |
| PATCH  | `/api/applications/:id/status` | Update status lamaran (Recruiter).  | `status: required, oneof=REVIEWED ACCEPTED REJECTED` |

## 5. Kriteria Penerimaan (Acceptance Criteria)

1. **Konfigurasi (Viper)**: Aplikasi tidak boleh memiliki hardcoded secret/port. Semua harus dibaca dari file konfigurasi.
2. **Database (GORM)**: Saat aplikasi berjalan pertama kali, GORM harus melakukan `AutoMigrate` untuk memastikan semua tabel (`users`, `auths`, `tokens`, `profiles`, `skills`, `profile_skills`, `experiences`, `companies`, `jobs`, `applications`) tersedia tanpa error. Relasi tabel berjalan dengan baik.
3. **Validasi (Validator)**: Jika payload tidak sesuai spesifikasi, API harus merespons dengan HTTP Status 400 (Bad Request) beserta detail error tiap field.
4. **Arsitektur & Dependency (Wire)**: File `wire_gen.go` berhasil digenerate, memastikan Inversion of Control (IoC) antara Handler, Service, dan Repository berjalan baik.
5. **Routing (GoFiber)**: Struktur endpoint terkelompok (Grouped routes) menggunakan middleware autentikasi (JWT) untuk endpoint yang membutuhkan hak akses khusus.
6. **Refresh Token**: Refresh token harus divalidasi ke tabel `tokens` sebelum menerbitkan Access Token baru dan dapat di-revoke saat logout.
7. **Eager Loading**: Endpoint untuk mengambil lowongan kerja (`GET /api/jobs`) harus memuat (Preload) data parsial relasi perusahaan (logo_url, name, location) agar frontend tidak perlu melakukan multi-fetching.

## 6. Fase Implementasi yang Disarankan

1. **Fase 1: Setup & Config** -> Inisialisasi modul Go, setup Viper, koneksi DB dengan GORM.
2. **Fase 2: Domain Layer** -> Pembuatan Struct entity model (`User`, `Auth`, `Token`, `Profile`, `Skill`, `ProfileSkill`, `Experience`, `Company`, `Job`, `Application`) dan DTOs.
3. **Fase 3: Repository & Service Layer** -> Implementasi interaksi DB, logika bisnis (Auth, Refresh Token, Profile), dan integrasi Validator pada Service/Handler.
4. **Fase 4: Dependency Injection** -> Setup Google Wire dan generate binding.
5. **Fase 5: Delivery Layer** -> Setup Fiber Controller, Routing, dan Middleware.
