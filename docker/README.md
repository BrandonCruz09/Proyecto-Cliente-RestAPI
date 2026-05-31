# Docker

Esta carpeta se deja para notas o archivos Docker adicionales del proyecto.

La configuracion principal esta en:

- `../docker-compose.yml`: levanta PostgreSQL y la API.
- `../backend/Dockerfile`: construye la imagen del backend Go.

El cliente de escritorio no tiene Dockerfile porque es una aplicacion Windows con ventana grafica. Se ejecuta con:

```cmd
cd frontend
go run .
```

