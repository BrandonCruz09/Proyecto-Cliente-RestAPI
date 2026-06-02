
## Archivo agregado

`database/reportes.sql`

## Vistas creadas

### vista_empleados_resumen

Muestra:

- ID del empleado
- nombre completo
- fecha de nacimiento
- genero
- fecha de contratacion
- titulo actual
- salario actual
- correo de usuario
- rol

Consulta:

```sql
SELECT * FROM vista_empleados_resumen ORDER BY emp_no;
```

### vista_historial_empleados

Muestra titulos y salarios como historial laboral.

Consulta:

```sql
SELECT * FROM vista_historial_empleados WHERE emp_no = 1 ORDER BY from_date;
```

### vista_clientes_directorio

Muestra un directorio simple de clientes.

Consulta:

```sql
SELECT * FROM vista_clientes_directorio ORDER BY cliente_id;
```
