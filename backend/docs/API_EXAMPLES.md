# API Testing Examples

## Register

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Citizen",
    "email": "jane@example.com",
    "password": "SecurePass123!",
    "role": "citizen",
    "phone": "+1234567890"
  }'
```

## Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"jane@example.com","password":"SecurePass123!"}'
```

## Create FIR

```bash
TOKEN="your_access_token"
curl -X POST http://localhost:8080/api/v1/fir/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Theft Report",
    "description": "Mobile phone stolen at central market.",
    "crimeType": "theft",
    "location": "Downtown Market"
  }'
```

## List FIRs with filters

```bash
curl "http://localhost:8080/api/v1/fir/all?status=pending&page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

## Upload evidence

```bash
curl -X POST http://localhost:8080/api/v1/evidence/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "firId=FIR-abc12345" \
  -F "file=@/path/to/evidence.jpg"
```

## Analytics (police/admin)

```bash
curl http://localhost:8080/api/v1/fir/analytics \
  -H "Authorization: Bearer $TOKEN"
```
