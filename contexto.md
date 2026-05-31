# Contexto del proyecto

## Encargo actual

El usuario pidio preparar el proyecto antes de iniciar desarrollo. Primero se debe crear y mostrar un plan de trabajo en `plan de trabajo.md`.

Tambien pidio empezar una memoria persistente llamada `contexto.md` para poder cambiar de chat sin perder el contexto.

## Requisitos solicitados

- Debe haber CRUD para empleados y clientes. El usuario corrigio que no quiso decir "cross".
- El sistema debe ejecutarse en Docker.
- La base de datos debe ser PostgreSQL version 17.9.
- La estructura de base de datos ya existe y se debe respetar.
- Cambio importante: el cliente no debe ser web.
- El cliente debe ser una aplicacion de escritorio para Windows, con ventana normal.
- El cliente de escritorio debe comunicarse con la API que estara en Docker.
- La API tambien debe estar en un contenedor.
- Usar puertos que no sean default.
- Regla de puertos: sumar 10 al puerto por default porque ya hay otras apps usando puertos como 8000 o 8080.

## Estado actual del repositorio

Carpeta raiz: `C:\Users\baapa\Desktop\proyecto_Servidor`

Archivos/carpetas detectadas:

- `backend/main.go`: API Go existente.
- `backend/Dockerfile`: Dockerfile del backend.
- `database/init.sql`: script SQL inicial.
- `frontend/main.go`: cliente de escritorio para Windows. Originalmente era Fyne; ahora esta migrado a `walk`.
- `docker-compose.yml`: compose actual con `db` y `api`.
- `.env`: variables basicas de DB.
- `pruebas.http`: pruebas HTTP actuales.

## Hallazgos importantes

- `docker-compose.yml` usa `postgres:15`, pero se necesita `postgres:17.9`.
- La DB publica `5432:5432`, pero debe evitarse el puerto default.
- La API publica `8080:8080`, pero debe moverse a `8090`.
- El backend escucha actualmente en `:8080`.
- El frontend actual llama a `http://localhost:8080/empleados`.
- El frontend original era de escritorio con Fyne, pero Fyne requeria CGO/OpenGL para compilar la ventana nativa en este entorno.
- Se migro el cliente a `github.com/lxn/walk`, una libreria Go para interfaces nativas de Windows, para poder compilar y validar el `.exe` sin depender de Fyne.
- Algunos endpoints extra del backend (`titulos`, `salarios`, `usuario`) devuelven datos fijos y deben consultar PostgreSQL.
- La estructura SQL actual no incluye tabla `clientes`; se debe agregar para cumplir CRUD de clientes.

## Cambios realizados

- `docker-compose.yml` ahora levanta:
  - `db` con imagen `postgres:17.9`.
  - `api` con backend Go.
  - PostgreSQL publicado en `5442`.
  - API publicada en `8090`.
- `.env` actualizado con `DB_PORT=5442` y `API_PORT=8090`.
- `database/init.sql` actualizado con:
  - Tablas `employees`, `titles`, `salaries`, `usuarios`.
  - Nueva tabla `clientes`.
  - Datos de prueba para empleados, historial laboral, salarios, usuarios y clientes.
- `backend/main.go` reemplazado por API completa:
  - CRUD empleados.
  - CRUD clientes.
  - Busqueda por nombre.
  - Consultas reales para titulos, salarios y usuario relacionado.
  - Puerto `8090`.
- `backend/Dockerfile` ahora expone `8090`.
- `frontend/main.go` ampliado como app de escritorio nativa Windows:
  - Pestana de empleados.
  - Pestana de clientes.
  - Formularios CRUD.
  - Conexion HTTP/JSON a `http://localhost:8090`.
- `frontend/main.go` redisenado con una interfaz mas pulida:
  - Ventana amplia.
  - Pestanas para empleados y clientes.
  - Encabezado con estado de API.
  - Panel izquierdo con formularios y acciones.
  - Panel derecho para respuestas JSON.
  - Botones agrupados por flujo de trabajo.
  - Secciones separadas para historial laboral de empleados.
- `pruebas.http` actualizado con endpoints en puerto `8090`.
- `README.md` creado con arquitectura, comandos y endpoints.
- `.gitignore` creado para caches locales y ejecutables.

## Puertos propuestos

- PostgreSQL: `5442` en host hacia `5432` en contenedor.
- API Go: `8090`.
- Cliente de escritorio: no usa puerto propio; llama a `http://localhost:8090`.

## Validaciones realizadas

- `go test ./...` en `backend`: correcto.
- `docker compose config`: correcto.
- `docker compose build`: correcto.
- `docker compose up -d`: correcto.
- `GET http://localhost:8090/empleados`: correcto.
- `GET http://localhost:8090/clientes`: correcto.
- `GET http://localhost:8090/empleados/1/titulos`: correcto.
- `POST http://localhost:8090/clientes`: correcto.
- `POST http://localhost:8090/empleados`: correcto.
- `go build -o cliente-sistema.exe .` en `frontend`: correcto.
- `go vet ./...` en `frontend`: correcto.
- Ejecucion real del `frontend/cliente-sistema.exe`: arranco, se mantuvo vivo 3 segundos y cerro con codigo 0.
- Se corrigio un error de arranque de la interfaz:
  - `TTM_ADDTOOL failed` venia de tooltips internos de `walk`.
  - `EM_SETCUEBANNER failed` venia de placeholders `CueBanner` en campos de texto.
  - Se usa una copia local parcheada de `walk` en `third_party/walk`.
  - Se quitaron los `CueBanner`.
  - Validacion posterior: la ventana abre con titulo `Sistema CRUD - Empleados y Clientes`, sin archivo `startup-error.txt`.
- Validacion profunda de puertos:
  - Docker publica API en `0.0.0.0:8090`.
  - Docker publica PostgreSQL en `0.0.0.0:5442`.
  - `Test-NetConnection localhost:8090`: correcto.
  - `Test-NetConnection localhost:5442`: correcto.
  - `GET http://localhost:8090/empleados`: correcto.
  - `GET http://localhost:8090/clientes`: correcto.

## Pendientes tecnicos

- No hay pendientes tecnicos abiertos para compilar el cliente.
- Nota historica: Fyne fue descartado porque no se pudo validar como `.exe` nativo sin errores en este entorno por dependencias CGO/OpenGL.

## Estado actual

El desarrollo principal esta implementado. Servidor y base de datos estan funcionando en Docker. La aplicacion de escritorio nativa de Windows compila y arranca correctamente.
