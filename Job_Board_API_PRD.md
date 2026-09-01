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
Menyimpan data autentikasi dan peran pengguna.
- `id` (UUID/Int, Primary Key)
- `name` (String)
- `email` (String, Unique)
- `password` (String, Hashed)
- `role` (Enum/String: 'RECRUITER' atau 'CANDIDATE')
- `created_at`, `updated_at`

### 3.2. Company
Menyimpan profil perusahaan (Dikelola oleh Recruiter).
- `id` (UUID/Int, Primary Key)
- `recruiter_id` (Foreign Key ke User)
- `name` (String)
- `description` (Text)
- `location` (String)
- `created_at`, `updated_at`

### 3.3. Job
Menyimpan detail lowongan pekerjaan.
- `id` (UUID/Int, Primary Key)
- `company_id` (Foreign Key ke Company)
- `title` (String)
- `description` (Text)
- `salary_range` (String)
- `status` (Enum/String: 'OPEN', 'CLOSED')
- `created_at`, `updated_at`

### 3.4. Application
Tabel transaksional untuk mencatat lamaran (Candidate melamar Job).
- `id` (UUID/Int, Primary Key)
- `job_id` (Foreign Key ke Job)
- `candidate_id` (Foreign Key ke User)
- `resume_url` (String)
- `status` (Enum/String: 'PENDING', 'REVIEWED', 'ACCEPTED', 'REJECTED')
- `created_at`, `updated_at`

## 4. Spesifikasi API Endpoint

### 4.1. Authentication (Public)
| Method | Endpoint | Keterangan | Validasi (Validator) |
|---|---|---|---|
| POST | `/api/auth/register` | Mendaftarkan user baru. | `email: required, email`, `password: required, min=6`, `role: required, oneof=RECRUITER CANDIDATE` |
| POST | `/api/auth/login` | Login user, mengembalikan token JWT. | `email: required, email`, `password: required` |

### 4.2. Company Management (Recruiter Only)
| Method | Endpoint | Keterangan | Validasi (Validator) |
|---|---|---|---|
| POST | `/api/companies` | Membuat profil perusahaan. | `name: required`, `location: required` |
| GET | `/api/companies` | Melihat perusahaan miliknya. | - |

### 4.3. Job Management (Recruiter & Public)
| Method | Endpoint | Keterangan | Validasi (Validator) |
|---|---|---|---|
| POST | `/api/jobs` | Membuat lowongan baru (Recruiter). | `title: required`, `company_id: required` |
| GET | `/api/jobs` | Melihat semua lowongan aktif (Public/Pagination). | - |
| GET | `/api/jobs/:id` | Melihat detail lowongan (Public). | - |

### 4.4. Application Management (Candidate & Recruiter)
| Method | Endpoint | Keterangan | Validasi (Validator) |
|---|---|---|---|
| POST | `/api/jobs/:id/apply` | Melamar pekerjaan (Candidate). | `resume_url: required, url` |
| GET | `/api/jobs/:id/applications` | Melihat daftar pelamar (Recruiter). | - |
| PATCH | `/api/applications/:id/status`| Update status lamaran (Recruiter).| `status: required, oneof=REVIEWED ACCEPTED REJECTED` |

## 5. Kriteria Penerimaan (Acceptance Criteria)

1. **Konfigurasi (Viper)**: Aplikasi tidak boleh memiliki hardcoded secret/port. Semua harus dibaca dari file konfigurasi.
2. **Database (GORM)**: Saat aplikasi berjalan pertama kali, GORM harus melakukan `AutoMigrate` untuk memastikan tabel tersedia tanpa error. Relasi tabel berjalan dengan baik (One-to-Many).
3. **Validasi (Validator)**: Jika payload tidak sesuai spesifikasi, API harus merespons dengan HTTP Status 400 (Bad Request) beserta detail error tiap field.
4. **Arsitektur & Dependency (Wire)**: File `wire_gen.go` berhasil digenerate, memastikan Inversion of Control (IoC) antara Handler, Service, dan Repository berjalan baik.
5. **Routing (GoFiber)**: Struktur endpoint terkelompok (Grouped routes) menggunakan middleware autentikasi (JWT) untuk endpoint yang membutuhkan hak akses khusus.

## 6. Fase Implementasi yang Disarankan
1. **Fase 1: Setup & Config** -> Inisialisasi modul Go, setup Viper, koneksi DB dengan GORM.
2. **Fase 2: Domain Layer** -> Pembuatan Struct entity model dan DTOs (Data Transfer Object).
3. **Fase 3: Repository & Service Layer** -> Implementasi interaksi DB, logika bisnis, dan integrasi Validator pada Service/Handler.
4. **Fase 4: Dependency Injection** -> Setup Google Wire dan generate binding.
5. **Fase 5: Delivery Layer** -> Setup Fiber Controller, Routing, dan Middleware.
