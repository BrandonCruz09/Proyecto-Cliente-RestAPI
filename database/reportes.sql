
CREATE OR REPLACE VIEW vista_empleados_resumen AS
SELECT
    e.emp_no,
    e.first_name,
    e.last_name,
    e.first_name || ' ' || e.last_name AS nombre_completo,
    e.birth_date,
    e.gender,
    e.hire_date,
    COALESCE(t.title, 'Sin titulo vigente') AS titulo_actual,
    COALESCE(s.salary, 0) AS salario_actual,
    COALESCE(u.correo, 'Sin usuario') AS correo_usuario,
    COALESCE(u.rol, 'Sin rol') AS rol_usuario
FROM employees e
LEFT JOIN titles t
    ON t.emp_no = e.emp_no
   AND t.to_date IS NULL
LEFT JOIN salaries s
    ON s.emp_no = e.emp_no
   AND s.to_date IS NULL
LEFT JOIN usuarios u
    ON u.empleado_id = e.emp_no;

CREATE OR REPLACE VIEW vista_historial_empleados AS
SELECT
    e.emp_no,
    e.first_name || ' ' || e.last_name AS empleado,
    'TITULO' AS tipo_movimiento,
    t.title AS detalle,
    NULL::INT AS salario,
    t.from_date,
    t.to_date
FROM employees e
INNER JOIN titles t
    ON t.emp_no = e.emp_no

UNION ALL

SELECT
    e.emp_no,
    e.first_name || ' ' || e.last_name AS empleado,
    'SALARIO' AS tipo_movimiento,
    'Cambio de salario' AS detalle,
    s.salary AS salario,
    s.from_date,
    s.to_date
FROM employees e
INNER JOIN salaries s
    ON s.emp_no = e.emp_no;

CREATE OR REPLACE VIEW vista_clientes_directorio AS
SELECT
    id AS cliente_id,
    nombre || ' ' || apellido AS cliente,
    correo,
    telefono,
    direccion,
    created_at
FROM clientes;
