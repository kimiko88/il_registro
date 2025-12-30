# Registro Elettronico - Frontend (Vue3/Quasar)

This is the frontend application for the Registro Elettronico, built with Vue 3, Quasar Framework, and Vite.

## Features
- **Multi-Role Support**: Admin, Teacher, Student, Parent, Director.
- **PWA Support**: Installable on mobile devices.
- **State Management**: Powered by Pinia.
- **Responsive Design**: Material Design via Quasar.

## Project Structure
- `src/components`: UI components organized by feature/role.
- `src/pages`: Application views (routes).
- `src/layouts`: Main layouts (Sidebar, Header).
- `src/stores`: Pinia state stores.
- `src/composables`: Reusable logic hooks.
- `docker`: Docker configuration.

## Setup
### Install Dependencies
```bash
npm install
```

### Development
```bash
npm run dev
```

### Build
```bash
npm run build
```

### Docker
```bash
docker build -t registro-frontend -f docker/Dockerfile .
docker run -p 8080:80 registro-frontend
```

## Testing
```bash
npm run test
```
