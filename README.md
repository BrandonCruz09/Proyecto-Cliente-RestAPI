# Sistema CRUD de Empleados y Clientes

Proyecto escolar completo con API REST en Go, cliente de escritorio en Go, base de datos PostgreSQL, Docker Compose y GitHub Actions.

La idea principal es sencilla:

- El cliente de escritorio muestra botones y formularios en una ventana de Windows.
- El cliente envia peticiones HTTP con datos JSON.
- La API REST recibe esas peticiones.
- La API consulta o modifica datos en PostgreSQL.
- Docker levanta la API y la base de datos sin instalar PostgreSQL manualmente.

## 1. Explicacion del proyecto

### Que es una REST API

Una REST API es un programa servidor que recibe solicitudes por internet o por red local. En este proyecto la API esta hecha en Go y responde usando JSON.

Ejemplo:

- El cliente quiere listar empleados.
- El cliente llama a `GET http://localhost:8090/empleados`.
- La API consulta la tabla `employees`.
- La API devuelve una lista en JSON.

JSON es un formato de texto para enviar datos. Ejemplo:

```json
{
  "first_name": "Diana",
  "last_name": "Ramirez",
  "birth_date": "1996-07-20",
  "gender": "F",
  "hire_date": "2025-01-15"
}
```

### Que hace el cliente

El cliente esta en `frontend/main.go`. Es una aplicacion de escritorio para Windows hecha en Go. Tiene campos de texto y botones para:

- Registrar empleado.
- Modificar empleado.
- Eliminar empleado.
- Buscar empleado por ID.
- Buscar empleado por nombre.
- Listar empleados.
- Consultar titulos, salarios y usuario relacionado.
- Hacer CRUD de clientes.

El cliente no se conecta directo a PostgreSQL. Se conecta a la API usando HTTP y JSON.

### Que hace PostgreSQL

PostgreSQL guarda los datos de forma relacional. Aqui se usan tablas con llaves primarias y foraneas:

- `employees`: empleados.
- `titles`: puestos o titulos laborales del empleado.
- `salaries`: salarios del empleado.
- `usuarios`: usuario relacionado con el empleado.
- `clientes`: clientes extra del sistema.

### Que hace Docker

Docker permite ejecutar la API y PostgreSQL en contenedores. Asi todos los integrantes pueden levantar el mismo proyecto con los mismos comandos.

### Como se conectan

```text
cliente de escritorio Go
        |
        | HTTP + JSON
        v
API REST Go
        |
        | SQL
        v
PostgreSQL
```

## 2. Arquitectura general

```text
proyecto_Servidor/
|
|-- frontend/                 cliente de escritorio local en Go
|   |-- main.go               Ventana, formularios, botones y consumo HTTP
|   |-- go.mod                Modulo Go del cliente
|
|-- backend/                  Servidor API REST en Go
|   |-- main.go               Conexion DB, rutas, handlers, modelos y queries
|   |-- Dockerfile            Imagen Docker del backend
|   |-- go.mod                Dependencias del backend
|
|-- database/
|   |-- init.sql              CREATE TABLE e INSERTS de prueba
|
|-- docker/
|   |-- README.md             Notas de Docker del proyecto
|
|-- .github/
|   |-- workflows/
|       |-- ci.yml            GitHub Actions
|
|-- docker-compose.yml        Levanta PostgreSQL y API
|-- .env                      Variables de entorno
|-- pruebas.http              Pruebas manuales de endpoints
|-- README.md                 Guia completa del proyecto
```

Flujo de datos:

```text
1. Usuario escribe datos en el cliente.
2. Cliente convierte datos a JSON.
3. Cliente envia HTTP a la API.
4. API valida el JSON.
5. API ejecuta SQL en PostgreSQL.
6. PostgreSQL devuelve datos.
7. API responde JSON.
8. Cliente muestra la respuesta.
```

## 3. Instalacion paso a paso para principiantes

### Instalar Go

1. Entra a `https://go.dev/dl/`.
2. Descarga el instalador de Windows.
3. Ejecuta el instalador.
4. Da clic en Next hasta terminar.
5. Abre CMD:
   - Presiona `Windows + R`.
   - Escribe `cmd`.
   - Presiona Enter.
6. Verifica Go:

```cmd
go version
```

Debe aparecer algo parecido a:

```cmd
go version go1.26.3 windows/amd64
```

### Instalar Docker Desktop

1. Entra a `https://www.docker.com/products/docker-desktop/`.
2. Descarga Docker Desktop para Windows.
3. Ejecuta el instalador.
4. Reinicia la computadora si Windows lo pide.
5. Abre Docker Desktop.
6. Espera a que diga que Docker esta corriendo.
7. Abre CMD y verifica:

```cmd
docker --version
docker compose version
```

### Instalar Git

1. Entra a `https://git-scm.com/download/win`.
2. Descarga Git para Windows.
3. Instala con las opciones por defecto.
4. Abre CMD o Git Bash.
5. Verifica:

```cmd
git --version
```

### Editores recomendados

Puedes usar cualquiera:

- VSCode.
- GoLand.
- Android Studio, aunque para Go es mas comun usar VSCode o GoLand.

Para abrir VSCode desde CMD:

```cmd
code .
```

## 4. Crear y abrir el proyecto

Desde CMD:

```cmd
cd Desktop
mkdir proyecto_Servidor
cd proyecto_Servidor
code .
```

En este workspace los archivos ya estan creados. Si lo haces desde cero, crea estas carpetas:

```cmd
mkdir backend
mkdir frontend
mkdir database
mkdir docker
mkdir .github
mkdir .github\workflows
```

Para crear archivos en VSCode:

1. Clic derecho en la carpeta.
2. New File.
3. Escribe el nombre, por ejemplo `main.go`.
4. Pega el codigo.
5. Guarda con `Ctrl + S`.

## 5. Base de datos

El script completo esta en `database/init.sql`.

Incluye:

- `CREATE TABLE employees`.
- `CREATE TABLE titles`.
- `CREATE TABLE salaries`.
- `CREATE TABLE usuarios`.
- `CREATE TABLE clientes`.
- Primary keys.
- Foreign keys.
- Inserts de prueba.

PostgreSQL se levanta por Docker, por eso no necesitas crear la base manualmente para ejecutar el proyecto.

Para revisar tablas dentro del contenedor:

```cmd
docker compose exec db psql -U postgres -d empresa
```

Dentro de `psql`:

```sql
\dt
SELECT * FROM employees;
SELECT * FROM titles;
SELECT * FROM salaries;
SELECT * FROM usuarios;
SELECT * FROM clientes;
```

Para salir:

```sql
\q
```

Si quieres recrear toda la base desde cero:

```cmd
docker compose down -v
docker compose up --build
```

## 6. Backend en Go

El backend esta en `backend/main.go`.

Ese archivo contiene:

- Modelos Go:
  - `Empleado`
  - `Titulo`
  - `Salario`
  - `Usuario`
  - `Cliente`
- Conexion a PostgreSQL con `database/sql`.
- Driver PostgreSQL `github.com/lib/pq`.
- Rutas HTTP.
- Handlers o controladores.
- Queries SQL.
- Respuestas JSON.

La conexion se lee desde `DATABASE_URL`. En Docker Compose se configura asi:

```text
postgres://postgres:admin123@db:5432/empresa?sslmode=disable
```

Si ejecutas el backend fuera de Docker, usa por defecto:

```text
postgres://postgres:admin123@localhost:5442/empresa?sslmode=disable
```

### Ejecutar backend sin Docker

Primero deja PostgreSQL corriendo con Docker:

```cmd
docker compose up db
```

En otra CMD:

```cmd
cd backend
go run .
```

La API queda en:

```text
http://localhost:8090
```

### Compilar backend

```cmd
cd backend
go build -o server.exe .
server.exe
```

## 7. Cliente de escritorio

El cliente esta en `frontend/main.go`.

Es una aplicacion de escritorio para Windows hecha en Go. Usa la libreria walk para mostrar una ventana real con formularios, botones y respuesta JSON.

### Ejecutar cliente

Primero levanta la API:

```cmd
docker compose up --build
```

En otra CMD:

```cmd
cd frontend
go run .
```

### Compilar cliente

```cmd
cd frontend
go build -o cliente-sistema.exe .
cliente-sistema.exe
```

### Cambiar URL de la API

Por defecto usa:

```text
http://localhost:8090
```

Si quieres cambiarla:

```cmd
cd frontend
set API_URL=http://localhost:8090
go run .
```

## 8. Docker

### docker-compose.yml

El archivo `docker-compose.yml` levanta dos servicios:

- `db`: PostgreSQL 17.9.
- `api`: backend Go.

Puertos:

- PostgreSQL en Windows: `localhost:5442`.
- PostgreSQL dentro de Docker: `db:5432`.
- API en Windows: `http://localhost:8090`.
### Levantar todo

```cmd
docker compose build
docker compose up
```

O en segundo plano:

```cmd
docker compose up -d --build
```

### Ver contenedores

```cmd
docker compose ps
```

### Ver logs

```cmd
docker compose logs api
docker compose logs db
```

### Apagar

```cmd
docker compose down
```

### Apagar y borrar la base

```cmd
docker compose down -v
```

## 9. Variables de entorno

Archivo `.env`:

```env
DB_USER=postgres
DB_PASSWORD=admin123
DB_NAME=empresa
DB_PORT=5442
API_PORT=8090
```

Explicacion:

- `DB_USER`: usuario de PostgreSQL.
- `DB_PASSWORD`: contrasena de PostgreSQL.
- `DB_NAME`: nombre de la base de datos.
- `DB_PORT`: puerto publicado en Windows.
- `API_PORT`: puerto publicado para la API.

## 10. Endpoints obligatorios

Todos responden JSON.

| Metodo | Endpoint | Funcion |
|---|---|---|
| POST | `/empleados` | Crear empleado |
| GET | `/empleados` | Listar empleados |
| GET | `/empleados/{id}` | Buscar empleado por ID |
| PUT | `/empleados/{id}` | Modificar empleado |
| DELETE | `/empleados/{id}` | Eliminar empleado |
| GET | `/empleados/buscar?nombre=` | Buscar empleado por nombre |
| GET | `/empleados/{id}/titulos` | Consultar titulos de empleado |
| GET | `/empleados/{id}/salarios` | Consultar salarios de empleado |
| GET | `/empleados/{id}/usuario` | Consultar usuario relacionado |

Tambien hay CRUD de clientes:

| Metodo | Endpoint | Funcion |
|---|---|---|
| POST | `/clientes` | Crear cliente |
| GET | `/clientes` | Listar clientes |
| GET | `/clientes/{id}` | Buscar cliente por ID |
| PUT | `/clientes/{id}` | Modificar cliente |
| DELETE | `/clientes/{id}` | Eliminar cliente |
| GET | `/clientes/buscar?nombre=` | Buscar cliente por nombre |

## 11. Pruebas con curl

Primero levanta el proyecto:

```cmd
docker compose up --build
```

### GET listar empleados

```cmd
curl http://localhost:8090/empleados
```

### POST crear empleado

```cmd
curl -X POST http://localhost:8090/empleados -H "Content-Type: application/json" -d "{\"first_name\":\"Diana\",\"last_name\":\"Ramirez\",\"birth_date\":\"1996-07-20\",\"gender\":\"F\",\"hire_date\":\"2025-01-15\"}"
```

### GET buscar por ID

```cmd
curl http://localhost:8090/empleados/1
```

### GET buscar por nombre

```cmd
curl "http://localhost:8090/empleados/buscar?nombre=Juan"
```

### PUT actualizar empleado

```cmd
curl -X PUT http://localhost:8090/empleados/1 -H "Content-Type: application/json" -d "{\"first_name\":\"Juan\",\"last_name\":\"Perez Gomez\",\"birth_date\":\"1990-01-01\",\"gender\":\"M\",\"hire_date\":\"2020-05-15\"}"
```

### DELETE eliminar empleado

```cmd
curl -X DELETE http://localhost:8090/empleados/3
```

### Consultar titulos

```cmd
curl http://localhost:8090/empleados/1/titulos
```

### Consultar salarios

```cmd
curl http://localhost:8090/empleados/1/salarios
```

### Consultar usuario relacionado

```cmd
curl http://localhost:8090/empleados/1/usuario
```

Tambien puedes usar `pruebas.http` con la extension REST Client de VSCode.

## 12. Casos de uso obligatorios

### 1. Alta de nuevo empleado

Que hace:

Registra un empleado nuevo en la tabla `employees`.

Como funciona:

1. El usuario escribe nombre, apellido, fecha de nacimiento, genero y fecha de contratacion.
2. El cliente envia JSON a la API.
3. La API ejecuta `INSERT INTO employees`.
4. PostgreSQL guarda el empleado.
5. La API responde el empleado creado.

Tablas:

- `employees`.

Endpoint:

- `POST /empleados`.

### 2. Cambio de puesto

Que hace:

Permite revisar el historial de puestos de un empleado. En una version ampliada, se agregaria un nuevo registro a `titles` para representar el nuevo puesto.

Como funciona en este proyecto:

1. El usuario escribe el ID del empleado.
2. El cliente consulta la API.
3. La API busca los puestos en `titles`.
4. La API responde el historial.

Tablas:

- `employees`.
- `titles`.

Endpoint:

- `GET /empleados/{id}/titulos`.

### 3. Cambio de salario

Que hace:

Permite revisar el historial salarial de un empleado. En una version ampliada, se agregaria un nuevo registro a `salaries`.

Como funciona en este proyecto:

1. El usuario escribe el ID del empleado.
2. El cliente llama a la API.
3. La API consulta `salaries`.
4. La API responde los salarios.

Tablas:

- `employees`.
- `salaries`.

Endpoint:

- `GET /empleados/{id}/salarios`.

### 4. Baja de empleado

Que hace:

Elimina un empleado del sistema.

Como funciona:

1. El usuario escribe el ID.
2. El cliente envia `DELETE`.
3. La API ejecuta `DELETE FROM employees`.
4. Como las tablas relacionadas usan `ON DELETE CASCADE`, tambien se eliminan titulos, salarios y usuario relacionados.

Tablas:

- `employees`.
- `titles`.
- `salaries`.
- `usuarios`.

Endpoint:

- `DELETE /empleados/{id}`.

### 5. Consulta de historial laboral

Que hace:

Muestra informacion del empleado, sus puestos, salarios y usuario relacionado.

Como funciona:

1. Se consulta el empleado por ID.
2. Se consultan sus titulos.
3. Se consultan sus salarios.
4. Se consulta su usuario.

Tablas:

- `employees`.
- `titles`.
- `salaries`.
- `usuarios`.

Endpoints:

- `GET /empleados/{id}`.
- `GET /empleados/{id}/titulos`.
- `GET /empleados/{id}/salarios`.
- `GET /empleados/{id}/usuario`.

## 13. GitHub paso a paso

### Crear cuenta

1. Entra a `https://github.com`.
2. Clic en Sign up.
3. Escribe correo, usuario y contrasena.
4. Confirma la cuenta.

### Crear repositorio

1. En GitHub, clic en New repository.
2. Nombre sugerido: `proyecto-servidor-crud`.
3. Elige Public o Private.
4. No marques README si ya tienes este proyecto local.
5. Clic en Create repository.

### Subir proyecto

Desde CMD en la carpeta del proyecto:

```cmd
git init
git add .
git commit -m "primer avance"
git branch -M main
git remote add origin URL_DEL_REPOSITORIO
git push -u origin main
```

Explicacion:

- `git init`: crea el repositorio Git local.
- `git add .`: prepara todos los archivos.
- `git commit -m "primer avance"`: guarda una version del proyecto.
- `git branch -M main`: renombra la rama principal a `main`.
- `git remote add origin URL_DEL_REPOSITORIO`: conecta el proyecto local con GitHub.
- `git push -u origin main`: sube los archivos a GitHub.

## 14. Ramas Git para 3 integrantes

Una rama es una linea de trabajo separada. Sirve para que cada integrante trabaje sin romper lo de los demas.

Ramas sugeridas:

- `main`: version estable.
- `backend`: API Go.
- `database`: SQL y Docker.
- `frontend`: cliente de escritorio.

Crear rama backend:

```cmd
git checkout -b backend
```

Subir rama:

```cmd
git push -u origin backend
```

Volver a main:

```cmd
git checkout main
```

Unir cambios:

```cmd
git merge backend
```

## 15. Trabajo para 3 integrantes

### Integrante 1

- Crear backend.
- Crear endpoints.
- Conectar PostgreSQL.
- Probar API.

Archivos principales:

- `backend/main.go`.
- `backend/Dockerfile`.
- `backend/go.mod`.

### Integrante 2

- Crear base de datos.
- Crear script SQL.
- Configurar Docker.
- Configurar Docker Compose.
- Agregar datos de prueba.

Archivos principales:

- `database/init.sql`.
- `docker-compose.yml`.
- `.env`.

### Integrante 3

- Crear cliente de escritorio.
- Conectar cliente con API.
- Crear formularios CRUD.
- Hacer documentacion.
- Subir a GitHub.

Archivos principales:

- `frontend/main.go`.
- `README.md`.
- `pruebas.http`.
- `.github/workflows/ci.yml`.

## 16. GitHub Actions

El workflow esta en:

```text
.github/workflows/ci.yml
```

Que hace:

- Instala Go.
- Descarga dependencias.
- Ejecuta `go test`.
- Compila backend.
- Compila frontend en Windows.
- Valida `docker compose config`.

Cuando se ejecuta:

- En cada `push`.
- En cada `pull_request`.

## 17. Comandos CMD exactos

Abrir proyecto:

```cmd
cd C:\Users\baapa\Desktop\proyecto_Servidor
```

Levantar Docker:

```cmd
docker compose up --build
```

Levantar en segundo plano:

```cmd
docker compose up -d --build
```

Ver contenedores:

```cmd
docker compose ps
```

Apagar:

```cmd
docker compose down
```

Apagar y borrar datos:

```cmd
docker compose down -v
```

Ejecutar backend local:

```cmd
cd backend
go run .
```

Compilar backend:

```cmd
cd backend
go build -o server.exe .
```

Ejecutar cliente:

```cmd
cd frontend
go run .
```

Compilar cliente:

```cmd
cd frontend
go build -o cliente-sistema.exe .
```

Ejecutar pruebas Go:

```cmd
cd backend
go test ./...
```

```cmd
cd frontend
go test ./...
```

## 18. Errores comunes

### Docker no abre

Solucion:

1. Abre Docker Desktop.
2. Espera a que termine de iniciar.
3. Vuelve a ejecutar:

```cmd
docker compose up --build
```

### Puerto ocupado

Si el puerto `8090` o `5442` esta ocupado, cambia `.env`:

```env
DB_PORT=5443
API_PORT=8091
```

Luego:

```cmd
docker compose down
docker compose up --build
```

### La base no actualiza cambios del SQL

Docker conserva datos en el volumen. Para recargar `init.sql`:

```cmd
docker compose down -v
docker compose up --build
```

### El cliente no conecta

Verifica que la API responda:

```cmd
curl http://localhost:8090/empleados
```

Si no responde, revisa logs:

```cmd
docker compose logs api
```

Verifica tambien que el cliente local responda:

```cmd
```

### Error de JSON

Revisa:

- Comillas dobles.
- Fechas con formato `YYYY-MM-DD`.
- `gender` con `M` o `F`.
- Header `Content-Type: application/json`.

## 19. Checklist final

- [x] API en Go creada.
- [x] PostgreSQL configurado.
- [x] Docker Compose configurado.
- [x] CRUD de empleados funcionando.
- [x] Busqueda por ID funcionando.
- [x] Busqueda por nombre funcionando.
- [x] Consulta de titulos funcionando.
- [x] Consulta de salarios funcionando.
- [x] Consulta de usuario relacionado funcionando.
- [x] cliente de escritorio creado.
- [x] Cliente consume HTTP y JSON.
- [x] CRUD de clientes agregado.
- [x] Variables de entorno creadas.
- [x] Pruebas manuales documentadas.
- [x] GitHub Actions agregado.
- [x] Casos de uso explicados.
- [x] Proyecto listo para entregar.

## 20. Archivos de codigo completos

El codigo completo esta guardado en estos archivos:

- `backend/main.go`: API REST completa.
- `frontend/main.go`: cliente de escritorio completo.
- `database/init.sql`: script SQL completo.
- `backend/Dockerfile`: Dockerfile backend.
- `docker-compose.yml`: Docker Compose completo.
- `.github/workflows/ci.yml`: GitHub Actions completo.

