#!/usr/bin/env bash
# Seed admin user via API (run after backend is up)
set -euo pipefail
API="${API:-http://localhost:8080}"

# First create admin manually in MongoDB or use admin API after bootstrapping.
# Example: register citizen then promote in DB, or POST /api/v1/admin/users with admin token.

echo "Use API examples in backend/docs/API_EXAMPLES.md"
echo "Create police/admin users via authenticated admin endpoint."
