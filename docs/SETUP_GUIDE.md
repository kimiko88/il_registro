# Setup & Installation Guide

This guide details how to set up the RegistroV2 development environment.

## Prerequisites
- **Go**: v1.21+
- **Node.js**: v18+ (LTS recommended)
- **PostgreSQL**: v14+
- **Docker** (Optional, for containerized run)

---

## 1. Minimal Setup (Local)

### Backend Setup
1.  Navigate to `registro-backend`:
    ```bash
    cd registro-backend
    ```

2.  **Environment Configuration**:
    Copy the example config (if available) or create `config.yaml` / set Environment Variables:
    ```bash
    # Example Env Vars
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_USER=postgres
    export DB_PASSWORD=secret
    export DB_NAME=registro_db
    export JWT_PRIVATE_KEY_PATH="path/to/private.pem"
    export JWT_PUBLIC_KEY_PATH="path/to/public.pem"
    ```

3.  **Database Migration**:
    Ensure Postgres is running and the database `registro_db` exists.
    Run migrations (assuming `goose` or internal migration tool):
    ```bash
    # If using Makefile
    make migrate-up
    ```

4.  **Run the Server**:
    ```bash
    go run cmd/api-server/main.go
    # API will typically listen on :8080
    ```

### Frontend Setup
1.  Navigate to `registro-frontend`:
    ```bash
    cd registro-frontend
    ```

2.  **Install Dependencies**:
    ```bash
    npm install
    ```

3.  **Configuration**:
    Check `.env` (create if needed) to point to Backend API:
    ```env
    VITE_API_URL=http://localhost:8080
    ```

4.  **Run Development Server**:
    ```bash
    npm run dev
    # UI will be available at http://localhost:9000 (usually)
    ```

---

## 2. Docker Setup (Recommended)
You can run the entire stack using `docker compose` (if file provided in root).

```bash
docker compose up --build
```

If separate Dockerfiles are used:
- **Backend Dockerfile**: `registro-backend/docker/Dockerfile.prod`
- **Frontend Dockerfile**: Needs creation or usage of generic Node image.

---

## 3. Testing

### Running Tests
**Backend**:
```bash
cd registro-backend
make test-unit       # Run logic tests
make test-integration # Run API tests
make test-coverage   # Generate coverage report
```

**Frontend**:
```bash
cd registro-frontend
npm run test:unit    # Run Vitest suite
npm run test:e2e     # Run End-to-End tests (if configured)
```

---

## 4. Common Troubleshooting

- **Database Connection Error**: Verify `DB_HOST` and credentials. If running in Docker, use the service name (e.g., `postgres`) instead of `localhost`.
- **JWT Error**: Ensure RSA keys are generated. You can use `openssl` to generate them:
    ```bash
    openssl genrsa -out private.pem 2048
    openssl rsa -in private.pem -pubout -out public.pem
    ```
- **CORS Issues**: Ensure the Backend has CORS middleware configured to allow requests from the Frontend origin (e.g., `http://localhost:9000`).
