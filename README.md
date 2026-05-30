# 🛡️ DS-Backend 

> Backend para la clase de **Desarrollo Seguro** - Sistema de gestión inmobiliaria seguro con Go y MySQL

## 📖 Descripción

InmoSoft Backend es un sistema de gestión inmobiliaria desarrollado en **Go** que implementa las mejores prácticas de seguridad para aplicaciones web. Permite la gestión de propiedades, propietarios, prospectos, citas y contratos con un sistema robusto de autenticación JWT y autorización basada en roles.

## 🚀 Tecnologías

- **Backend**: Go (Golang) 1.21+
- **Base de Datos**: MySQL 8.0
- **Containerización**: Docker & Docker Compose
- **Database Driver**: database/sql con MySQL driver
- **Router**: Gin Web Framework
- **Autenticación**: JWT (JSON Web Tokens)
- **Hashing**: bcrypt para contraseñas
- **CORS**: gin-contrib/cors

## ✨ Características

- 🔐 **Autenticación JWT** con cookies HTTP-only
- 👥 **Sistema de roles** (Admin, Agente)
- 🏠 **Gestión de propiedades** completa
- 👤 **Administración de propietarios y prospectos**
- 📅 **Sistema de citas**
- 📄 **Gestión de contratos y documentos**
- 🖼️ **Manejo de imágenes**
- 🔒 **Seguridad robusta** con middleware de autorización
- 🐳 **Completamente dockerizado**

## 📋 Requisitos Previos

Antes de comenzar, asegúrate de tener instalado:

### Windows 🪟
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) v4.0+
- [Git](https://git-scm.com/download/win)

### Linux/macOS 🐧🍎
- [Docker](https://docs.docker.com/engine/install/) v20.10+
- [Docker Compose](https://docs.docker.com/compose/install/) v2.0+
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

Crea el archivo de variables de entorno:

**Windows:**
```powershell
copy .env.example .env
notepad .env
```

**Linux/macOS:**
```bash
cp .env.example .env
nano .env  # o vim .env
```

## ⚙️ Configuración

Edita el archivo `.env` con tus configuraciones **sin espacios alrededor del =**:

```env
# Base de datos MySQL
DB_USER=inmosoft_user
DB_PASSWORD=tu_password_muy_segura_123
DB_HOST=localhost
DB_PORT=3306
DB_NAME=inmosoftDB
USER_PASSWORD=tu_password_muy_segura_123

# API
API_PORT=8080

# Seguridad JWT
JWT_SECRET=tu_jwt_secret_super_seguro_con_mas_de_32_caracteres
```

> ⚠️ **Importante**: Usa contraseñas y secretos fuertes en producción

## 🚀 Uso

### Levantar los servicios

**Primera vez o después de cambios:**

**Windows (PowerShell/CMD):**
```powershell
docker-compose up -d --build
```

**Linux/macOS:**
```bash
docker compose up -d --build
```

### Comandos útiles

**Ver logs en tiempo real:**
```bash
# Windows
docker-compose logs -f

# Linux/macOS
docker compose logs -f
```

**Ver logs específicos:**
```bash
# Backend
docker logs ds_backend -f

# Base de datos
docker logs ds_database -f
```

**Parar servicios:**
```bash
# Windows
docker-compose down

# Linux/macOS
docker compose down
```

**Reiniciar completamente:**
```bash
# Windows
docker-compose down
docker-compose up -d --build

# Linux/macOS
docker compose down
docker compose up -d --build
```

**Acceder a la base de datos:**
```bash
docker exec -it ds_database mysql -u root -p inmosoftDB
```

## 🌐 API Endpoints

La API estará disponible en: `http://localhost:8080`

### Autenticación
| Método | Endpoint | Descripción | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/login` | Iniciar sesión | ❌ |
| POST | `/api/v1/auth/logout` | Cerrar sesión | ✅ |

### Usuarios
| Método | Endpoint | Descripción | Roles |
|--------|----------|-------------|-------|
| GET | `/api/v1/users/all` | Listar todos los usuarios | Admin |
| GET | `/api/v1/users/:id` | Obtener usuario específico | Admin, Owner |
| POST | `/api/v1/users/create` | Crear nuevo usuario | Admin |
| PUT | `/api/v1/users/:id` | Actualizar usuario | Admin, Owner |
| DELETE | `/api/v1/users/:id` | Eliminar usuario | Admin |

### Propiedades
| Método | Endpoint | Descripción | Roles |
|--------|----------|-------------|-------|
| GET | `/api/v1/propiedades/all` | Listar propiedades | Todos |
| GET | `/api/v1/propiedades/:id` | Obtener propiedad | Todos |
| POST | `/api/v1/propiedades/create` | Crear propiedad | Admin, Agente |
| PUT | `/api/v1/propiedades/update/:id` | Actualizar propiedad | Admin, Agente |
| DELETE | `/api/v1/propiedades/eliminar/:id` | Eliminar propiedad | Admin |

### Citas
| Método | Endpoint | Descripción | Roles |
|--------|----------|-------------|-------|
| GET | `/api/v1/citas/all/:id` | Obtener citas del usuario | Admin, Owner |
| POST | `/api/v1/citas/create` | Crear cita | Todos |
| PUT | `/api/v1/citas/update/:id` | Actualizar cita | Admin, Owner |
| DELETE | `/api/v1/citas/eliminar/:id` | Eliminar cita | Admin |

### Otros endpoints disponibles:
- **Propietarios**: `/api/v1/propietarios/*`
- **Prospectos**: `/api/v1/prospectos/*`
- **Contratos**: `/api/v1/contratos/*`
- **Tipos de propiedad**: `/api/v1/tipos-propiedad/*`
- **Imágenes**: `/api/v1/imagenes/*`
- **Documentos**: `/api/v1/documentos/*`

## 🔐 Autenticación y Autorización

### Sistema de Roles

**Admin** 👑
- Acceso completo al sistema
- Gestión de usuarios
- Eliminación de recursos
- Visualización de todos los datos

**Agente** 🏠
- Gestión de propiedades
- Creación de citas y contratos
- Acceso solo a sus propios recursos
- No puede eliminar usuarios

### Autenticación JWT

El sistema usa **cookies HTTP-only** para mayor seguridad:

```javascript
// Login (Frontend)
const response = await fetch('/api/v1/auth/login', {
    method: 'POST',
    credentials: 'include', // Importante para cookies
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
});

// Requests autenticados
const data = await fetch('/api/v1/users/profile', {
    credentials: 'include' // Cookie se envía automáticamente
});
```

### Reset completo del proyecto:

```bash
# Parar y eliminar todo
docker-compose down -v
docker system prune -f

# Reconstruir desde cero
docker-compose up -d --build
```

## 📝 Notas de Seguridad

### En Desarrollo:
- ✅ JWT tokens en cookies HTTP-only
- ✅ Hashing de contraseñas con bcrypt
- ✅ Validación de entrada
- ✅ CORS configurado correctamente
