# Nutrition Platform

A comprehensive nutrition and fitness platform with backend API and frontend dashboard.

## 🚀 Quick Start

### Backend
```bash
cd backend
go build -o bin/server .
./bin/server
```

### Frontend
```bash
cd frontend-nextjs
npm install
npm run dev
```

## 📋 Prerequisites

- Go 1.24+
- Node.js 18+
- PostgreSQL (optional, SQLite supported)
- Redis (optional, in-memory fallback)

## 🔐 Environment Variables

Copy `.env.example` files and configure:
- `backend/.env.example` → `backend/.env`
- `frontend-nextjs/.env.example` → `frontend-nextjs/.env.local`

## 📚 Documentation

- [Deployment Guide](30_MINUTE_DEPLOYMENT_GUIDE.md)
- [Production Status](PRODUCTION_STATUS.md)
- [API Documentation](docs/API_REFERENCE.md)

## 🔒 Security

- Never commit `.env` files
- Use environment variables for secrets
- See [SECURITY.md](.github/SECURITY.md) for details

## 📝 License

[Your License Here]
