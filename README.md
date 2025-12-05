# 👤 Service Customer - Desa Murni Batik

Perkhidmatan pelanggan untuk platform **Desa Murni Batik**.

## 🚀 Ciri-ciri

- 👤 **Profiles** - Profil pelanggan
- 📍 **Addresses** - Alamat penghantaran
- ❤️ **Wishlist** - Senarai hajat
- 📋 **Order History** - Sejarah pesanan

## 🛠️ Tech Stack

- Go 1.21+
- Gin Framework
- GORM
- PostgreSQL

## 📦 Setup

```bash
go mod download
go run cmd/server/main.go
```

Server: http://localhost:8084

## 🔗 Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/customers/me` | My profile |
| PUT | `/api/v1/customers/me` | Update profile |
| GET | `/api/v1/customers/addresses` | Addresses |
| GET | `/api/v1/customers/wishlist` | Wishlist |

---

**© 2024 Desa Murni Batik** | [KilangDesaMurniBatik](https://github.com/KilangDesaMurniBatik)
