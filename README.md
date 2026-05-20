# Blockchain FIR Management System

Tamper-proof **First Information Report (FIR)** management with blockchain hash storage, IPFS evidence, role-based access, and audit logging.

## Stack

| Layer | Technology |
|-------|------------|
| Frontend | Flutter, Riverpod, Go Router, Dio, fl_chart |
| Backend | Go, Gin, MongoDB, JWT, bcrypt |
| Blockchain | Solidity ^0.8.20, Hardhat, go-ethereum |
| Storage | IPFS (Pinata) |
| DevOps | Docker, GitHub Actions |

## Project Structure

```
decntrafir/
├── frontend/          # Flutter mobile/web app
├── backend/           # Go REST API (MVC)
├── contracts/         # Solidity + Hardhat
├── scripts/           # Utility scripts
├── test/              # Integration test placeholders
└── docker-compose.yml
```

## Quick Start

### 1. Start infrastructure

```bash
cp .env.example .env
docker compose up -d mongodb
```

### 2. Smart contracts

```bash
cd contracts && npm install
npx hardhat compile
npm test
npx hardhat node          # optional local chain
npm run deploy:local
```

### 3. Backend

```bash
cd backend
cp .env.example .env
# Set CONTRACT_ADDRESS from deploy output
go mod tidy
go run ./cmd/server
```

API: `http://localhost:8080`  
Swagger: `http://localhost:8080/swagger/index.html`  
Health: `http://localhost:8080/health`

### 4. Flutter frontend

```bash
cd frontend
flutter pub get
flutter run -d chrome --dart-define=API_BASE_URL=http://localhost:8080
```

## Roles

| Role | Capabilities |
|------|----------------|
| **Citizen** | Register FIR, upload evidence, view own FIRs |
| **Police** | Review/update FIR status, assign officers, analytics |
| **Admin** | Manage users, system analytics |

## API Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Citizen registration |
| POST | `/login` | JWT login |
| POST | `/api/v1/fir/create` | Create FIR (+ blockchain tx) |
| GET | `/api/v1/fir/:id` | Get FIR |
| PUT | `/api/v1/fir/update/:id` | Update status (police/admin) |
| GET | `/api/v1/fir/all` | List with filters |
| GET | `/api/v1/fir/analytics` | Dashboard stats |
| POST | `/api/v1/evidence/upload` | IPFS evidence upload |

See [backend/docs/API_EXAMPLES.md](backend/docs/API_EXAMPLES.md).

## Blockchain Flow

1. Citizen submits FIR → backend stores metadata in MongoDB.
2. SHA-256 hash of metadata is written to `FIRManagement` smart contract.
3. Transaction hash saved on FIR record for verification.
4. Evidence files uploaded to IPFS; CID stored in MongoDB.

## Environment Variables

- Root: [.env.example](.env.example)
- Backend: [backend/.env.example](backend/.env.example)
- Contracts: [contracts/.env.example](contracts/.env.example)
- Frontend: [frontend/.env.example](frontend/.env.example)

## Docker

```bash
docker compose up --build
```

## Testing

```bash
# Contracts
cd contracts && npm test

# Backend build
cd backend && go build ./...

# Flutter analyze
cd frontend && flutter analyze
```

## Deployment

- **Backend**: Docker image → Render/Railway/Fly.io
- **Frontend**: Flutter web → Vercel/Netlify
- **Contracts**: Polygon Mumbai → see [contracts/docs/DEPLOYMENT.md](contracts/docs/DEPLOYMENT.md)
- **Database**: MongoDB Atlas

## Security Notes

- Change `JWT_SECRET` in production.
- Use Pinata JWT for IPFS uploads.
- Never commit private keys.
- Enable HTTPS and restrict CORS origins.
- Blockchain runs in mock mode if `CONTRACT_ADDRESS` is unset.

## License

MIT
