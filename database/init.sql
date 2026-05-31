CREATE TABLE employees (
    emp_no SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    birth_date DATE NOT NULL,
    gender CHAR(1) NOT NULL,
    hire_date DATE NOT NULL
);

CREATE TABLE titles (
    id SERIAL PRIMARY KEY,
    emp_no INT NOT NULL REFERENCES employees(emp_no) ON DELETE CASCADE,
    title VARCHAR(50) NOT NULL,
    from_date DATE NOT NULL,
    to_date DATE
);

CREATE TABLE salaries (
    id SERIAL PRIMARY KEY,
    emp_no INT NOT NULL REFERENCES employees(emp_no) ON DELETE CASCADE,
    salary INT NOT NULL,
    from_date DATE NOT NULL,
    to_date DATE
);

CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    correo VARCHAR(100) UNIQUE NOT NULL,
    contrasena VARCHAR(100) NOT NULL,
    rol VARCHAR(20) NOT NULL,
    empleado_id INT UNIQUE NOT NULL REFERENCES employees(emp_no) ON DELETE CASCADE
);

CREATE TABLE clientes (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(80) NOT NULL,
    apellido VARCHAR(80) NOT NULL,
    correo VARCHAR(120) UNIQUE NOT NULL,
    telefono VARCHAR(30),
    direccion VARCHAR(200),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO employees (first_name, last_name, birth_date, gender, hire_date) VALUES
('Juan', 'Perez', '1990-01-01', 'M', '2020-05-15'),
('Ana', 'Lopez', '1994-09-18', 'F', '2021-03-10'),
('Carlos', 'Martinez', '1998-05-12', 'M', '2024-02-01');

INSERT INTO titles (emp_no, title, from_date, to_date) VALUES
(1, 'Ingeniero de Software', '2020-05-15', '2022-12-31'),
(1, 'Lider Tecnico', '2023-01-01', NULL),
(2, 'Analista de Datos', '2021-03-10', NULL),
(3, 'Soporte Tecnico', '2024-02-01', NULL);

INSERT INTO salaries (emp_no, salary, from_date, to_date) VALUES
(1, 60000, '2020-05-15', '2022-12-31'),
(1, 78000, '2023-01-01', NULL),
(2, 52000, '2021-03-10', NULL),
(3, 38000, '2024-02-01', NULL);

INSERT INTO usuarios (correo, contrasena, rol, empleado_id) VALUES
('juan@empresa.com', '123456', 'admin', 1),
('ana@empresa.com', '123456', 'empleado', 2),
('carlos@empresa.com', '123456', 'empleado', 3);

INSERT INTO clientes (nombre, apellido, correo, telefono, direccion) VALUES
('Maria', 'Gomez', 'maria.gomez@cliente.com', '555-100-2000', 'Av. Central 123'),
('Roberto', 'Sanchez', 'roberto.sanchez@cliente.com', '555-200-3000', 'Calle Norte 45'),
('Lucia', 'Hernandez', 'lucia.hernandez@cliente.com', '555-300-4000', 'Boulevard Sur 789');
