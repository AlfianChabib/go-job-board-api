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
- `headline` (String)
- `bio` (Text)
- `phone` (String)
- `skills` (Text)
- `resume_url` (String)
- `created_at`, `updated_at`

### 3.5. Company
Menyimpan profil perusahaan (Dikelola oleh Recruiter).
- `id` (UUID/Int, Primary Key)
- `recruiter_id` (Foreign Key ke User)
- `name` (String)
- `description` (Text)
- `location` (String)
- `created_at`, `updated_at`

### 3.6. Job
Menyimpan detail lowongan pekerjaan.
- `id` (UUID/Int, Primary Key)
- `company_id` (Foreign Key ke Company)
- `title` (String)
- `description` (Text)
- `salary_range` (String)
- `status` (Enum/String: 'OPEN', 'CLOSED')
- `created_at`, `updated_at`

### 3.7. Application
Tabel transaksional untuk mencatat lamaran (Candidate melamar Job).
- `id` (UUID/Int, Primary Key)
- `job_id` (Foreign Key ke Job)
- `candidate_id` (Foreign Key ke User)
- `resume_url` (String)
- `status` (Enum/String: 'PENDING', 'REVIEWED', 'ACCEPTED', 'REJECTED')
- `created_at`, `updated_at`

## 4. Spesifikasi API Endpoint

### 4.1. Authentication (Public & Authenticated)
| Method | Endpoint | Keterangan | Validasi (Validator) |
|---|---|---|---|
| POST | `/api/auth/register` | Mendaftarkan user baru (menyimpan ke `users` & `auths`). | `email: required, email`, `password: required, min=6`, `role: required, oneof=RECRUITER CANDIDATE`, `name: required` |
| POST | `/api/auth/login` | Login user, mengembalikan access token JWT & refresh token (menyimpan ke `tokens`). | `email: required, email`, `password: required` |
| POST | `/api/auth/refresh` | Memperbarui access token dengan mengirimkan refresh token yang valid. | `refresh_token: required` |
| POST | `/api/auth/logout` | Logout user dan me-revoke refresh token di tabel `tokens`. | `refresh_token: required` |

### 4.2. Candidate Profile (Candidate Only)
| Method | Endpoint | Keterangan | Validasi (Validator) |
|---|---|---|---|
| GET | `/api/candidate/profile` | Melihat profil pribadi Candidate. | - |
| PUT | `/api/candidate/profile` | Membuat / memperbarui profil Candidate. | `headline: required`, `phone: required` |

### 4.3. Company Management (Recruiter Only)
| Method | Endpoint | Keterangan | Validasi (Validator) |
|---|---|---|---|
| POST | `/api/companies` | Membuat profil perusahaan. | `name: required`, `location: required` |
| GET | `/api/companies` | Melihat perusahaan miliknya. | - |

### 4.4. Job Management (Recruiter & Public)
| Method | Endpoint | Keterangan | Validasi (Validator) |
|---|---|---|---|
| POST | `/api/jobs` | Membuat lowongan baru (Recruiter). | `title: required`, `company_id: required` |
| GET | `/api/jobs` | Melihat semua lowongan aktif (Public/Pagination). | - |
| GET | `/api/jobs/:id` | Melihat detail lowongan (Public). | - |

### 4.5. Application Management (Candidate & Recruiter)
| Method | Endpoint | Keterangan | Validasi (Validator) |
|---|---|---|---|
| POST | `/api/jobs/:id/apply` | Melamar pekerjaan (Candidate). | `resume_url: required, url` |
| GET | `/api/jobs/:id/applications` | Melihat daftar pelamar (Recruiter). | - |
| PATCH | `/api/applications/:id/status`| Update status lamaran (Recruiter).| `status: required, oneof=REVIEWED ACCEPTED REJECTED` |

## 5. Kriteria Penerimaan (Acceptance Criteria)

1. **Konfigurasi (Viper)**: Aplikasi tidak boleh memiliki hardcoded secret/port. Semua harus dibaca dari file konfigurasi.
2. **Database (GORM)**: Saat aplikasi berjalan pertama kali, GORM harus melakukan `AutoMigrate` untuk memastikan semua tabel (`users`, `auths`, `tokens`, `profiles`, `companies`, `jobs`, `applications`) tersedia tanpa error. Relasi tabel berjalan dengan baik.
3. **Validasi (Validator)**: Jika payload tidak sesuai spesifikasi, API harus merespons dengan HTTP Status 400 (Bad Request) beserta detail error tiap field.
4. **Arsitektur & Dependency (Wire)**: File `wire_gen.go` berhasil digenerate, memastikan Inversion of Control (IoC) antara Handler, Service, dan Repository berjalan baik.
5. **Routing (GoFiber)**: Struktur endpoint terkelompok (Grouped routes) menggunakan middleware autentikasi (JWT) untuk endpoint yang membutuhkan hak akses khusus.
6. **Refresh Token**: Refresh token harus divalidasi ke tabel `tokens` sebelum menerbitkan Access Token baru dan dapat di-revoke saat logout.

## 6. Fase Implementasi yang Disarankan
1. **Fase 1: Setup & Config** -> Inisialisasi modul Go, setup Viper, koneksi DB dengan GORM.
2. **Fase 2: Domain Layer** -> Pembuatan Struct entity model (`User`, `Auth`, `Token`, `Profile`, `Company`, `Job`, `Application`) dan DTOs.
3. **Fase 3: Repository & Service Layer** -> Implementasi interaksi DB, logika bisnis (Auth, Refresh Token, Profile), dan integrasi Validator pada Service/Handler.
4. **Fase 4: Dependency Injection** -> Setup Google Wire dan generate binding.
5. **Fase 5: Delivery Layer** -> Setup Fiber Controller, Routing, dan Middleware.

