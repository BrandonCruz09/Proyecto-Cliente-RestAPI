package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Empleado struct {
	EmpNo     int    `json:"emp_no"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	HireDate  string `json:"hire_date"`
}

type Cliente struct {
	ID        int    `json:"id"`
	Nombre    string `json:"nombre"`
	Apellido  string `json:"apellido"`
	Correo    string `json:"correo"`
	Telefono  string `json:"telefono"`
	Direccion string `json:"direccion"`
	CreatedAt string `json:"created_at"`
}

type Titulo struct {
	ID       int    `json:"id"`
	EmpNo    int    `json:"emp_no"`
	Title    string `json:"title"`
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}

type Salario struct {
	ID       int    `json:"id"`
	EmpNo    int    `json:"emp_no"`
	Salary   int    `json:"salary"`
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
}

type Usuario struct {
	ID         int    `json:"id"`
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena"`
	Rol        string `json:"rol"`
	EmpleadoID int    `json:"empleado_id"`
}

type mensaje struct {
	Message string `json:"message"`
}

var db *sql.DB

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:admin123@localhost:5442/empresa?sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error al preparar conexion a PostgreSQL: ", err)
	}
	defer db.Close()

	if err := esperarBaseDeDatos(); err != nil {
		log.Fatal("No se pudo conectar a PostgreSQL: ", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /empleados", crearEmpleado)
	mux.HandleFunc("GET /empleados", obtenerEmpleados)
	mux.HandleFunc("GET /empleados/buscar", buscarEmpleado)
	mux.HandleFunc("GET /empleados/{id}", obtenerEmpleadoPorID)
	mux.HandleFunc("PUT /empleados/{id}", actualizarEmpleado)
	mux.HandleFunc("DELETE /empleados/{id}", eliminarEmpleado)
	mux.HandleFunc("GET /empleados/{id}/titulos", obtenerTitulos)
	mux.HandleFunc("GET /empleados/{id}/salarios", obtenerSalarios)
	mux.HandleFunc("GET /empleados/{id}/usuario", obtenerUsuario)

	mux.HandleFunc("POST /clientes", crearCliente)
	mux.HandleFunc("GET /clientes", obtenerClientes)
	mux.HandleFunc("GET /clientes/buscar", buscarCliente)
	mux.HandleFunc("GET /clientes/{id}", obtenerClientePorID)
	mux.HandleFunc("PUT /clientes/{id}", actualizarCliente)
	mux.HandleFunc("DELETE /clientes/{id}", eliminarCliente)

	log.Println("API corriendo en http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", mux))
}

func esperarBaseDeDatos() error {
	var err error
	for i := 1; i <= 20; i++ {
		err = db.Ping()
		if err == nil {
			return nil
		}
		log.Printf("Esperando PostgreSQL... intento %d/20", i)
		time.Sleep(2 * time.Second)
	}
	return err
}

func responderJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Error escribiendo JSON:", err)
	}
}

func responderError(w http.ResponseWriter, status int, texto string) {
	responderJSON(w, status, mensaje{Message: texto})
}

func decodificarJSON(r *http.Request, destino any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(destino)
}

func nullString(valor sql.NullString) string {
	if !valor.Valid {
		return ""
	}
	return valor.String
}

func crearEmpleado(w http.ResponseWriter, r *http.Request) {
	var emp Empleado
	if err := decodificarJSON(r, &emp); err != nil {
		responderError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	if err := normalizarEmpleado(&emp); err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}

	query := `
		INSERT INTO employees (first_name, last_name, birth_date, gender, hire_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING emp_no`

	err := db.QueryRow(query, emp.FirstName, emp.LastName, emp.BirthDate, emp.Gender, emp.HireDate).Scan(&emp.EmpNo)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusCreated, emp)
}

func obtenerEmpleados(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT emp_no, first_name, last_name,
		       to_char(birth_date, 'YYYY-MM-DD'),
		       gender,
		       to_char(hire_date, 'YYYY-MM-DD')
		FROM employees
		ORDER BY emp_no`)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	empleados, err := escanearEmpleados(rows)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, empleados)
}

func obtenerEmpleadoPorID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	emp, err := consultarEmpleadoPorID(id)
	if errors.Is(err, sql.ErrNoRows) {
		responderError(w, http.StatusNotFound, "Empleado no encontrado")
		return
	}
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, emp)
}

func actualizarEmpleado(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var emp Empleado
	if err := decodificarJSON(r, &emp); err != nil {
		responderError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	if err := normalizarEmpleado(&emp); err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}

	resultado, err := db.Exec(`
		UPDATE employees
		SET first_name=$1, last_name=$2, birth_date=$3, gender=$4, hire_date=$5
		WHERE emp_no=$6`,
		emp.FirstName, emp.LastName, emp.BirthDate, emp.Gender, emp.HireDate, id)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if filas, _ := resultado.RowsAffected(); filas == 0 {
		responderError(w, http.StatusNotFound, "Empleado no encontrado")
		return
	}

	actualizado, err := consultarEmpleadoPorID(id)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, actualizado)
}

func eliminarEmpleado(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	resultado, err := db.Exec("DELETE FROM employees WHERE emp_no = $1", id)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if filas, _ := resultado.RowsAffected(); filas == 0 {
		responderError(w, http.StatusNotFound, "Empleado no encontrado")
		return
	}
	responderJSON(w, http.StatusOK, mensaje{Message: "Empleado eliminado"})
}

func buscarEmpleado(w http.ResponseWriter, r *http.Request) {
	nombre := r.URL.Query().Get("nombre")
	rows, err := db.Query(`
		SELECT emp_no, first_name, last_name,
		       to_char(birth_date, 'YYYY-MM-DD'),
		       gender,
		       to_char(hire_date, 'YYYY-MM-DD')
		FROM employees
		WHERE first_name ILIKE $1 OR last_name ILIKE $1
		ORDER BY emp_no`, "%"+nombre+"%")
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	empleados, err := escanearEmpleados(rows)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, empleados)
}

func consultarEmpleadoPorID(id string) (Empleado, error) {
	var emp Empleado
	err := db.QueryRow(`
		SELECT emp_no, first_name, last_name,
		       to_char(birth_date, 'YYYY-MM-DD'),
		       gender,
		       to_char(hire_date, 'YYYY-MM-DD')
		FROM employees
		WHERE emp_no = $1`, id).
		Scan(&emp.EmpNo, &emp.FirstName, &emp.LastName, &emp.BirthDate, &emp.Gender, &emp.HireDate)
	return emp, err
}

func escanearEmpleados(rows *sql.Rows) ([]Empleado, error) {
	empleados := []Empleado{}
	for rows.Next() {
		var emp Empleado
		if err := rows.Scan(&emp.EmpNo, &emp.FirstName, &emp.LastName, &emp.BirthDate, &emp.Gender, &emp.HireDate); err != nil {
			return nil, err
		}
		empleados = append(empleados, emp)
	}
	return empleados, rows.Err()
}

func normalizarEmpleado(emp *Empleado) error {
	var err error
	emp.FirstName = strings.TrimSpace(emp.FirstName)
	emp.LastName = strings.TrimSpace(emp.LastName)

	if emp.FirstName == "" {
		return errors.New("El nombre del empleado es obligatorio")
	}
	if emp.LastName == "" {
		return errors.New("El apellido del empleado es obligatorio")
	}

	emp.BirthDate, err = normalizarFecha(emp.BirthDate, false)
	if err != nil {
		return errors.New("Fecha de nacimiento invalida. Usa YYYY-MM-DD, DD-MM-YY o DD-MM-YYYY")
	}

	emp.HireDate, err = normalizarFecha(emp.HireDate, true)
	if err != nil {
		return errors.New("Fecha de contratacion invalida. Usa YYYY-MM-DD, DD-MM-YY, DD-MM-YYYY, Hoy o Inmediata")
	}

	emp.Gender, err = normalizarGenero(emp.Gender)
	if err != nil {
		return err
	}
	return nil
}

func normalizarFecha(valor string, aceptarInmediata bool) (string, error) {
	valor = strings.TrimSpace(valor)
	if valor == "" {
		return "", errors.New("fecha vacia")
	}

	normalizada := strings.ToLower(valor)
	if aceptarInmediata {
		switch normalizada {
		case "hoy", "inmediata", "inmediato", "actual", "now", "today":
			return time.Now().Format("2006-01-02"), nil
		}
	}

	formatos := []string{
		"2006-01-02",
		"02-01-06",
		"2-1-06",
		"02/01/06",
		"2/1/06",
		"02-01-2006",
		"2-1-2006",
		"02/01/2006",
		"2/1/2006",
	}
	for _, formato := range formatos {
		fecha, err := time.Parse(formato, valor)
		if err == nil {
			return fecha.Format("2006-01-02"), nil
		}
	}
	return "", errors.New("fecha invalida")
}

func normalizarGenero(valor string) (string, error) {
	valor = strings.TrimSpace(strings.ToLower(valor))
	switch valor {
	case "f", "femenino", "mujer":
		return "F", nil
	case "m", "masculino", "hombre":
		return "M", nil
	default:
		return "", errors.New("Genero invalido. Usa M, F, Masculino o Femenino")
	}
}

func obtenerTitulos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, emp_no, title, to_char(from_date, 'YYYY-MM-DD'), to_char(to_date, 'YYYY-MM-DD')
		FROM titles
		WHERE emp_no = $1
		ORDER BY from_date DESC`, r.PathValue("id"))
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	titulos := []Titulo{}
	for rows.Next() {
		var titulo Titulo
		var toDate sql.NullString
		if err := rows.Scan(&titulo.ID, &titulo.EmpNo, &titulo.Title, &titulo.FromDate, &toDate); err != nil {
			responderError(w, http.StatusInternalServerError, err.Error())
			return
		}
		titulo.ToDate = nullString(toDate)
		titulos = append(titulos, titulo)
	}
	responderJSON(w, http.StatusOK, titulos)
}

func obtenerSalarios(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, emp_no, salary, to_char(from_date, 'YYYY-MM-DD'), to_char(to_date, 'YYYY-MM-DD')
		FROM salaries
		WHERE emp_no = $1
		ORDER BY from_date DESC`, r.PathValue("id"))
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	salarios := []Salario{}
	for rows.Next() {
		var salario Salario
		var toDate sql.NullString
		if err := rows.Scan(&salario.ID, &salario.EmpNo, &salario.Salary, &salario.FromDate, &toDate); err != nil {
			responderError(w, http.StatusInternalServerError, err.Error())
			return
		}
		salario.ToDate = nullString(toDate)
		salarios = append(salarios, salario)
	}
	responderJSON(w, http.StatusOK, salarios)
}

func obtenerUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario Usuario
	err := db.QueryRow(`
		SELECT id, correo, contrasena, rol, empleado_id
		FROM usuarios
		WHERE empleado_id = $1`, r.PathValue("id")).
		Scan(&usuario.ID, &usuario.Correo, &usuario.Contrasena, &usuario.Rol, &usuario.EmpleadoID)
	if errors.Is(err, sql.ErrNoRows) {
		responderError(w, http.StatusNotFound, "Usuario no encontrado para este empleado")
		return
	}
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, usuario)
}

func crearCliente(w http.ResponseWriter, r *http.Request) {
	var cliente Cliente
	if err := decodificarJSON(r, &cliente); err != nil {
		responderError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}

	err := db.QueryRow(`
		INSERT INTO clientes (nombre, apellido, correo, telefono, direccion)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, to_char(created_at, 'YYYY-MM-DD HH24:MI:SS')`,
		cliente.Nombre, cliente.Apellido, cliente.Correo, cliente.Telefono, cliente.Direccion).
		Scan(&cliente.ID, &cliente.CreatedAt)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusCreated, cliente)
}

func obtenerClientes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, nombre, apellido, correo, telefono, direccion, to_char(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM clientes
		ORDER BY id`)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	clientes, err := escanearClientes(rows)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, clientes)
}

func obtenerClientePorID(w http.ResponseWriter, r *http.Request) {
	cliente, err := consultarClientePorID(r.PathValue("id"))
	if errors.Is(err, sql.ErrNoRows) {
		responderError(w, http.StatusNotFound, "Cliente no encontrado")
		return
	}
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, cliente)
}

func actualizarCliente(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var cliente Cliente
	if err := decodificarJSON(r, &cliente); err != nil {
		responderError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}

	resultado, err := db.Exec(`
		UPDATE clientes
		SET nombre=$1, apellido=$2, correo=$3, telefono=$4, direccion=$5
		WHERE id=$6`,
		cliente.Nombre, cliente.Apellido, cliente.Correo, cliente.Telefono, cliente.Direccion, id)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if filas, _ := resultado.RowsAffected(); filas == 0 {
		responderError(w, http.StatusNotFound, "Cliente no encontrado")
		return
	}

	actualizado, err := consultarClientePorID(id)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, actualizado)
}

func eliminarCliente(w http.ResponseWriter, r *http.Request) {
	resultado, err := db.Exec("DELETE FROM clientes WHERE id = $1", r.PathValue("id"))
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if filas, _ := resultado.RowsAffected(); filas == 0 {
		responderError(w, http.StatusNotFound, "Cliente no encontrado")
		return
	}
	responderJSON(w, http.StatusOK, mensaje{Message: "Cliente eliminado"})
}

func buscarCliente(w http.ResponseWriter, r *http.Request) {
	nombre := r.URL.Query().Get("nombre")
	rows, err := db.Query(`
		SELECT id, nombre, apellido, correo, telefono, direccion, to_char(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM clientes
		WHERE nombre ILIKE $1 OR apellido ILIKE $1
		ORDER BY id`, "%"+nombre+"%")
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	clientes, err := escanearClientes(rows)
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, clientes)
}

func consultarClientePorID(id string) (Cliente, error) {
	var cliente Cliente
	err := db.QueryRow(`
		SELECT id, nombre, apellido, correo, telefono, direccion, to_char(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM clientes
		WHERE id = $1`, id).
		Scan(&cliente.ID, &cliente.Nombre, &cliente.Apellido, &cliente.Correo, &cliente.Telefono, &cliente.Direccion, &cliente.CreatedAt)
	return cliente, err
}

func escanearClientes(rows *sql.Rows) ([]Cliente, error) {
	clientes := []Cliente{}
	for rows.Next() {
		var cliente Cliente
		if err := rows.Scan(&cliente.ID, &cliente.Nombre, &cliente.Apellido, &cliente.Correo, &cliente.Telefono, &cliente.Direccion, &cliente.CreatedAt); err != nil {
			return nil, err
		}
		clientes = append(clientes, cliente)
	}
	return clientes, rows.Err()
}
