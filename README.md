# 🛡️ DS-Backend

> Backend para la clase de **Desarrollo Seguro** - Sistema de gestión seguro con Go y MySQL

## 📋 Tabla de Contenidos

- [Descripción](#descripción)
- [Tecnologías](#tecnologías)
- [Requisitos Previos](#requisitos-previos)
- [Instalación](#instalación)
- [Configuración](#configuración)
- [Uso](#uso)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Contribución](#contribución)

## 📖 Descripción

Este proyecto es un backend desarrollado en **Go** que implementa las mejores prácticas de seguridad para aplicaciones web. Utiliza **MySQL** como base de datos y está completamente dockerizado para facilitar el desarrollo y despliegue.

## 🚀 Tecnologías

- **Backend**: Go (Golang)
- **Base de Datos**: MySQL 8.0
- **Containerización**: Docker & Docker Compose
- **ORM**: GORM
- **Router**: Gin
- **Autenticación**: JWT

## 📋 Requisitos Previos

Antes de comenzar, asegúrate de tener instalado:

### Windows 🪟
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Git](https://git-scm.com/download/win)

### Linux/macOS 🐧🍎
- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- Git

## 🔧 Instalación

### 1. Clonar el repositorio

**Windows (PowerShell/CMD):**
```powershell
git clone https://github.com/tu-usuario/DS-backend.git
cd DS-backend
```

**Linux/macOS:**
```bash
git clone https://github.com/tu-usuario/DS-backend.git
cd DS-backend
```

### 2. Configuración del entorno

Renombra el archivo de ejemplo y configura las variables de entorno:

**Windows:**
```powershell
copy env-example .env
notepad .env  # o tu editor
```

**Linux/macOS:**
```bash
cp env-example .env
nano .env  # o vim .env
```

## ⚙️ Configuración

Edita el archivo `.env` con tus configuraciones **sin espacios alrededor del =**:

```env
# Base de datos
DB_USER=tu_usuario
DB_PASSWORD=tu_password_segura
DB_HOST=localhost
DB_PORT=3306
DB_NAME=ds_database

# API
API_PORT=8080

# Seguridad
JWT_SECRET=tu_jwt_secret_muy_seguro
```

## 🚀 Uso

### Levantar los servicios

**Primera vez o después de cambios:**

**Windows (PowerShell/CMD):**

Puedes utilizar docker desktop con compose

```powershell
docker-compose up -d --build
```

**Linux/macOS:**
```bash
docker compose up -d --build
```

### Comandos útiles

**Ver logs:**
```bash
# Windows
docker-compose logs -f

# Linux/macOS
docker compose logs -f
```

**Parar servicios:**
```bash
# Windows
docker-compose down

# Linux/macOS
docker compose down
```

**Reiniciar completamente (con rebuild):**
```bash
# Windows
docker-compose down
docker-compose up -d --build

# Linux/macOS
docker compose down
docker compose up -d --build
```

**Ver contenedores activos:**
```bash
docker ps
```

## 🌐 API Endpoints

Una vez que el servidor esté corriendo, la API estará disponible en: `http://localhost:3000`

### Ejemplos de endpoints:

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET    | `/health` | Verificar estado del servidor |
| POST   | `/auth/login` | Iniciar sesión |
| POST   | `/auth/register` | Registrar usuario |
| GET    | `/api/users` | Obtener usuarios (requiere auth) |

## 🧪 Testing

### Usando Postman
1. Descarga [Postman Desktop](https://www.postman.com/downloads/)
2. Importa la colección de endpoints (si está disponible)
3. Configura el environment con `base_url: http://localhost:3000`

### Usando curl

**Verificar que el servidor esté corriendo:**
```bash
curl http://localhost:8080/health
```

## 🔍 Troubleshooting

### Problemas comunes:

**❌ "Database is not reachable"**
- Verifica que MySQL esté corriendo: `docker ps`
- Revisa los logs: `docker logs ds_database`

**❌ "Port already in use"**
- Cambia el puerto en `.env` o libera el puerto 8080
- En Windows: `netstat -ano | findstr :8080`
- En Linux: `sudo lsof -i :8080`

**❌ Variables de entorno no cargadas**
- Asegúrate de no tener espacios en el archivo `.env`
- Formato correcto: `DB_USER=valor` (sin espacios)

### Logs útiles:

```bash
# Ver logs del backend
docker logs ds_backend

# Ver logs de MySQL
docker logs ds_database

# Ver logs en tiempo real
docker logs -f ds_backend

## 📝 Notas de Seguridad

- ⚠️ Nunca commitees el archivo `.env` al repositorio
- 🔐 Usa contraseñas seguras para la base de datos
- 🛡️ El JWT secret debe ser único y complejo

