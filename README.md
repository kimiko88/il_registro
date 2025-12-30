# RegistroV2 - Electronic School Register

RegistroV2 is a comprehensive, modern solution for managing Italian school activities, including grades, attendance, digital documents, and communications.

## 🚀 Key Features
- **Role-Based Access**: Specialized views for Teachers, Students, Parents, and Administrators.
- **Advanced Grading**: Supports numeric grades, judgments, credits, and weighted averages.
- **Real-time Attendance**: Track presence, delays, and justifications.
- **Digital Class Register**: Seamless management of daily activities.
- **Documents & Workflows**: Digital signature and circular management.
- **Secure**: JWT Authentication, MFA, and Audit Logging.

## 📁 Repository Structure
- **[Back-end](./registro-backend/)**: Go (Golang) API server.
- **[Front-end](./registro-frontend/)**: Vue 3 + Quasar PWA/SPA.
- **[Docs](./docs/)**: Detailed project documentation.

## 📚 Documentation
- **[Architecture & Design](./docs/ARCHITECTURE.md)**: Detailed breakdown of code structure and modules.
- **[Setup Guide](./docs/SETUP_GUIDE.md)**: Instructions for installation, configuration, and running locally.
- **[API Documentation]**: (See Postman collection or Swagger if available).

## 🛠 Quick Start
```bash
# 1. Start Backend
cd registro-backend
go run cmd/api-server/main.go

# 2. Start Frontend (in new terminal)
cd registro-frontend
npm install && npm run dev
```

## 🧪 Testing
We maintain a high level of test coverage (~80%).

**Run Backend Tests:**
```bash
cd registro-backend && make test-unit
```

**Run Frontend Tests:**
```bash
cd registro-frontend && npm run test:unit
```

## 📜 License
Private / Proprietary.
