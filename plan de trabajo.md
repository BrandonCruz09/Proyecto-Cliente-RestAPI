# Plan de trabajo

## Objetivo

Actualizar el proyecto para que el servidor funcione en Docker y el cliente funcione como aplicacion de escritorio con:

- Base de datos PostgreSQL version 17.9.
- API REST en Go conectada a PostgreSQL.
- Aplicacion cliente de escritorio para Windows que consuma la API por HTTP/JSON.
- CRUD completo para empleados y clientes.
- Puertos no default, usando la regla de sumar 10 al puerto convencional.
- Memoria de contexto en `contexto.md` para poder continuar el trabajo en otro chat.

## Estado actual detectado

El proyecto ya tiene esta estructura:

- `backend/`: API en Go.
- `frontend/`: cliente de escritorio en Go para Windows.
- `database/init.sql`: estructura inicial de PostgreSQL.
- `docker-compose.yml`: servicios `db` y `api`.
- `.env`: variables basicas de base de datos.
- `pruebas.http`: pruebas manuales contra la API.

Puntos importantes encontrados:

- `docker-compose.yml` usa actualmente `postgres:15`; se debe cambiar a `postgres:17.9`.
- PostgreSQL expone actualmente `5432:5432`; se debe evitar el puerto default.
- La API expone actualmente `8080:8080`; se debe cambiar a `8090`.
- La API tiene endpoints de empleados, pero algunos endpoints extra responden datos fijos y deben conectarse realmente a la base de datos.
- El cliente debe ser de escritorio, no web.
- El cliente de escritorio debe ejecutarse como una ventana normal de Windows y comunicarse con la API que estara en Docker.

## Convencion de puertos propuesta

Para evitar puertos default y respetar la regla de sumar 10:

- PostgreSQL: host `5442` -> contenedor `5432` (`5432 + 10 = 5442` en el host).
- API Go: host `8090` -> contenedor `8090` (`8080 + 10 = 8090`).
- Cliente de escritorio: no publica puerto porque no es servidor web.

Nota: PostgreSQL puede seguir escuchando internamente en `5432` dentro del contenedor porque es el puerto propio de la imagen. Lo importante para no chocar con otras apps del equipo es que el puerto publicado en la maquina sea `5442`.

## Alcance CRUD

El usuario aclaro que no se refiere a "cross", sino a CRUD.

Por lo tanto, el sistema debera tener operaciones completas para:

- Empleados.
- Clientes.

CRUD significa:

- Crear registros.
- Consultar/listar registros.
- Buscar registros por ID y, cuando aplique, por nombre.
- Actualizar registros.
- Eliminar registros.

Nota tecnica: al ser aplicacion de escritorio, normalmente no se necesita CORS porque CORS aplica a navegadores web. El cliente consumira la API directamente con peticiones HTTP.

## Fase 1: Preparacion y respaldo de contexto

1. Mantener actualizado `contexto.md` con:
   - Requisitos del usuario.
   - Decisiones tecnicas.
   - Puertos usados.
   - Cambios realizados.
   - Comandos importantes.
   - Pendientes.
2. Revisar que no haya cambios del usuario que no deban tocarse.
3. Confirmar estructura final de carpetas antes de editar codigo.

## Fase 2: Docker y variables de entorno

1. Actualizar `docker-compose.yml`:
   - `db` con `postgres:17.9`.
   - Puerto de PostgreSQL `5442:5432`.
   - API en puerto `8090`.
   - Red interna para que API y DB se comuniquen.
   - Volumen persistente para PostgreSQL.
2. Ajustar `.env`:
   - `DB_USER`
   - `DB_PASSWORD`
   - `DB_NAME`
   - `DB_PORT=5442`
   - `API_PORT=8090`
   - `DATABASE_URL` si conviene centralizarla.
3. Verificar que los contenedores no usen puertos default en el host.

## Fase 3: Base de datos PostgreSQL 17.9

1. Revisar `database/init.sql`.
2. Mantener las tablas requeridas:
   - `employees`
   - `titles`
   - `salaries`
   - `usuarios`
3. Confirmar llaves primarias y foraneas.
4. Agregar datos de prueba suficientes para:
   - Listado de empleados.
   - Titulos por empleado.
   - Salarios por empleado.
   - Usuario relacionado con empleado.
5. Si realmente se necesita una tabla `clientes`, definirla antes de implementarla porque no aparece en la estructura original.

## Fase 4: Backend API en Go

1. Ajustar el servidor para escuchar en `8090`.
2. Agregar CRUD completo de empleados.
3. Mejorar manejo de errores JSON.
4. Validar metodos y respuestas HTTP.
5. Implementar endpoints reales con consultas SQL:
   - `POST /empleados`
   - `GET /empleados`
   - `GET /empleados/{id}`
   - `PUT /empleados/{id}`
   - `DELETE /empleados/{id}`
   - `GET /empleados/buscar?nombre=`
   - `GET /empleados/{id}/titulos`
   - `GET /empleados/{id}/salarios`
   - `GET /empleados/{id}/usuario`
6. Agregar CRUD completo de clientes:
   - `POST /clientes`
   - `GET /clientes`
   - `GET /clientes/{id}`
   - `PUT /clientes/{id}`
   - `DELETE /clientes/{id}`
   - `GET /clientes/buscar?nombre=`
7. Ajustar `backend/Dockerfile` para exponer `8090`.

## Fase 5: Aplicacion cliente de escritorio para Windows

1. Mantener el cliente como aplicacion de escritorio en Go.
2. Usar una libreria grafica nativa de Windows validable sin errores; se eligio `github.com/lxn/walk`.
3. Crear una ventana normal de Windows para:
   - Listar empleados.
   - Crear empleado.
   - Buscar empleado por ID.
   - Buscar empleado por nombre.
   - Editar empleado.
   - Eliminar empleado.
   - Ver titulos, salarios y usuario relacionado.
   - Listar clientes.
   - Crear cliente.
   - Buscar cliente por ID o nombre.
   - Editar cliente.
   - Eliminar cliente.
4. Consumir la API por HTTP/JSON usando el paquete `net/http` de Go.
5. Configurar la URL de API como `http://localhost:8090`, porque la API estara publicada desde Docker hacia Windows.
6. Agregar mensajes visibles de exito/error en la ventana.
7. Documentar como ejecutar el cliente desde Windows:
   - `cd frontend`
   - `go run .`
   - o compilar con `go build -o cliente-empleados.exe .`

## Fase 6: Pruebas

1. Probar Docker:
   - `docker compose build`
   - `docker compose up`
   - `docker compose ps`
2. Probar PostgreSQL:
   - Contenedor arriba.
   - Tablas creadas.
   - Datos iniciales cargados.
3. Probar API:
   - `GET http://localhost:8090/empleados`
   - `POST http://localhost:8090/empleados`
   - `PUT http://localhost:8090/empleados/{id}`
   - `DELETE http://localhost:8090/empleados/{id}`
   - Endpoints de titulos, salarios y usuario.
4. Probar CRUD de clientes.
5. Probar que el cliente de escritorio se conecta a `http://localhost:8090`.
6. Compilar el cliente para generar un `.exe`.
7. Actualizar `pruebas.http` con los nuevos puertos.

## Fase 7: Documentacion

1. Crear o actualizar `README.md` con:
   - Explicacion del proyecto.
   - Arquitectura.
   - Puertos.
   - Comandos Docker.
   - Endpoints.
   - Pruebas.
   - Errores comunes.
2. Documentar casos de uso:
   - Alta de nuevo empleado.
   - Cambio de puesto.
   - Cambio de salario.
   - Baja de empleado.
   - Consulta de historial laboral.
3. Si se requiere GitHub Actions, crear workflow en `.github/workflows/`.

## Orden de ejecucion propuesto

1. Actualizar memoria `contexto.md`.
2. Ajustar Docker Compose y `.env`.
3. Actualizar base de datos.
4. Corregir y completar backend.
5. Completar cliente de escritorio.
6. Probar con Docker.
7. Actualizar pruebas y documentacion.

## Pendientes antes de desarrollar

- Ninguno. El desarrollo ya inicio y se eligieron campos para `clientes`: `id`, `nombre`, `apellido`, `correo`, `telefono`, `direccion`, `created_at`.

Decision aplicada: usar cliente de escritorio nativo con `walk`. No crear frontend web.
