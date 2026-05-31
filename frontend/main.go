package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Empleado struct {
	EmpNo     int    `json:"emp_no,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	HireDate  string `json:"hire_date"`
}

type Cliente struct {
	ID        int    `json:"id,omitempty"`
	Nombre    string `json:"nombre"`
	Apellido  string `json:"apellido"`
	Correo    string `json:"correo"`
	Telefono  string `json:"telefono"`
	Direccion string `json:"direccion"`
	CreatedAt string `json:"created_at,omitempty"`
}

type empleadoModel struct {
	walk.TableModelBase
	items []Empleado
}

func (m *empleadoModel) RowCount() int {
	return len(m.items)
}

func (m *empleadoModel) Value(row, col int) interface{} {
	e := m.items[row]
	switch col {
	case 0:
		return e.EmpNo
	case 1:
		return e.FirstName
	case 2:
		return e.LastName
	case 3:
		return e.BirthDate
	case 4:
		return e.Gender
	case 5:
		return e.HireDate
	}
	return ""
}

func (m *empleadoModel) setItems(items []Empleado) {
	m.items = items
	m.PublishRowsReset()
}

type clienteModel struct {
	walk.TableModelBase
	items []Cliente
}

func (m *clienteModel) RowCount() int {
	return len(m.items)
}

func (m *clienteModel) Value(row, col int) interface{} {
	c := m.items[row]
	switch col {
	case 0:
		return c.ID
	case 1:
		return c.Nombre
	case 2:
		return c.Apellido
	case 3:
		return c.Correo
	case 4:
		return c.Telefono
	case 5:
		return c.Direccion
	}
	return ""
}

func (m *clienteModel) setItems(items []Cliente) {
	m.items = items
	m.PublishRowsReset()
}

var (
	apiURL = obtenerAPIURL()

	mw     *walk.MainWindow
	status *walk.Label
	output *walk.TextEdit

	empTable                                                                         *walk.TableView
	empModel                                                                         = &empleadoModel{}
	empID, empNombre, empApellido, empNacimiento, empGenero, empIngreso, empBusqueda *walk.LineEdit

	cliTable                                                                         *walk.TableView
	cliModel                                                                         = &clienteModel{}
	cliID, cliNombre, cliApellido, cliCorreo, cliTelefono, cliDireccion, cliBusqueda *walk.LineEdit
)

func main() {
	_, err := MainWindow{
		AssignTo: &mw,
		Title:    "Sistema CRUD - Empleados y Clientes",
		Visible:  true,
		MinSize:  Size{Width: 1180, Height: 740},
		Layout:   VBox{MarginsZero: true, SpacingZero: true},
		Children: []Widget{
			Composite{
				Layout:  VBox{Margins: Margins{Left: 18, Top: 14, Right: 18, Bottom: 10}},
				MaxSize: Size{Height: 100},
				Children: []Widget{
					Label{Text: "Sistema CRUD", Font: Font{PointSize: 15, Bold: true}},
					Label{Text: "Selecciona un registro en la tabla para editarlo, eliminarlo o consultar su historial."},
					Label{Text: "API: " + apiURL + "    PostgreSQL: localhost:5442"},
				},
			},
			TabWidget{
				StretchFactor: 3,
				Pages: []TabPage{
					empleadosPage(),
					clientesPage(),
				},
			},
			GroupBox{
				Title:  "Estado y respuesta",
				Layout: VBox{Margins: Margins{Left: 12, Top: 8, Right: 12, Bottom: 12}},
				MinSize: Size{
					Height: 180,
				},
				Children: []Widget{
					Label{AssignTo: &status, Text: "Listo. Presiona Recargar para consultar la API."},
					TextEdit{
						AssignTo: &output,
						ReadOnly: true,
						VScroll:  true,
						HScroll:  true,
						Text:     "Aqui se mostraran resultados, errores y confirmaciones.",
					},
				},
			},
		},
	}.Run()
	if err != nil {
		walk.MsgBox(nil, "Error", "No se pudo crear la ventana: "+err.Error(), walk.MsgBoxIconError)
	}
}

func empleadosPage() TabPage {
	return TabPage{
		Title:  "Empleados",
		Layout: HBox{Margins: Margins{Left: 12, Top: 10, Right: 12, Bottom: 10}},
		Children: []Widget{
			GroupBox{
				Title:   "Listado de empleados",
				Layout:  VBox{Margins: Margins{Left: 10, Top: 8, Right: 10, Bottom: 10}},
				MinSize: Size{Width: 650},
				Children: []Widget{
					Composite{
						Layout: HBox{},
						Children: []Widget{
							PushButton{Text: "Recargar", OnClicked: cargarEmpleados},
							Label{Text: "Buscar por nombre"},
							LineEdit{AssignTo: &empBusqueda},
							PushButton{Text: "Buscar", OnClicked: buscarEmpleadoNombre},
						},
					},
					TableView{
						AssignTo:            &empTable,
						AlternatingRowBG:    true,
						ColumnsOrderable:    true,
						LastColumnStretched: true,
						Columns: []TableViewColumn{
							{Title: "ID", Width: 55},
							{Title: "Nombre", Width: 120},
							{Title: "Apellido", Width: 120},
							{Title: "Nacimiento", Width: 105},
							{Title: "Genero", Width: 70},
							{Title: "Contratacion", Width: 110},
						},
						Model:                    empModel,
						OnSelectedIndexesChanged: seleccionarEmpleado,
						OnItemActivated:          seleccionarEmpleado,
					},
				},
			},
			GroupBox{
				Title:   "Formulario de empleado",
				Layout:  VBox{Margins: Margins{Left: 14, Top: 10, Right: 14, Bottom: 14}},
				MinSize: Size{Width: 410},
				Children: []Widget{
					Label{Text: "1. Selecciona una fila o llena los datos para crear un empleado."},
					Composite{
						Layout: Grid{Columns: 2, Spacing: 8},
						Children: []Widget{
							Label{Text: "ID"},
							LineEdit{AssignTo: &empID, ReadOnly: true},
							Label{Text: "Nombre"},
							LineEdit{AssignTo: &empNombre},
							Label{Text: "Apellido"},
							LineEdit{AssignTo: &empApellido},
							Label{Text: "Nacimiento"},
							LineEdit{AssignTo: &empNacimiento, Text: "1990-01-01"},
							Label{Text: "Genero"},
							LineEdit{AssignTo: &empGenero, Text: "F"},
							Label{Text: "Contratacion"},
							LineEdit{AssignTo: &empIngreso, Text: "Hoy"},
						},
					},
					Label{Text: "2. Ejecuta una accion."},
					Composite{
						Layout: Grid{Columns: 2, Spacing: 8},
						Children: []Widget{
							PushButton{Text: "Crear nuevo", OnClicked: crearEmpleado},
							PushButton{Text: "Guardar cambios", OnClicked: actualizarEmpleado},
							PushButton{Text: "Buscar ID", OnClicked: buscarEmpleadoID},
							PushButton{Text: "Eliminar seleccionado", OnClicked: eliminarEmpleado},
							PushButton{Text: "Ver titulos", OnClicked: func() { historialEmpleado("titulos") }},
							PushButton{Text: "Ver salarios", OnClicked: func() { historialEmpleado("salarios") }},
							PushButton{Text: "Ver usuario", OnClicked: func() { historialEmpleado("usuario") }},
							PushButton{Text: "Limpiar", OnClicked: limpiarEmpleados},
						},
					},
					Label{Text: "Fechas aceptadas: YYYY-MM-DD, DD-MM-YY, DD-MM-YYYY, Hoy o Inmediata."},
				},
			},
		},
	}
}

func clientesPage() TabPage {
	return TabPage{
		Title:  "Clientes",
		Layout: HBox{Margins: Margins{Left: 12, Top: 10, Right: 12, Bottom: 10}},
		Children: []Widget{
			GroupBox{
				Title:   "Listado de clientes",
				Layout:  VBox{Margins: Margins{Left: 10, Top: 8, Right: 10, Bottom: 10}},
				MinSize: Size{Width: 650},
				Children: []Widget{
					Composite{
						Layout: HBox{},
						Children: []Widget{
							PushButton{Text: "Recargar", OnClicked: cargarClientes},
							Label{Text: "Buscar por nombre"},
							LineEdit{AssignTo: &cliBusqueda},
							PushButton{Text: "Buscar", OnClicked: buscarClienteNombre},
						},
					},
					TableView{
						AssignTo:            &cliTable,
						AlternatingRowBG:    true,
						ColumnsOrderable:    true,
						LastColumnStretched: true,
						Columns: []TableViewColumn{
							{Title: "ID", Width: 55},
							{Title: "Nombre", Width: 120},
							{Title: "Apellido", Width: 120},
							{Title: "Correo", Width: 180},
							{Title: "Telefono", Width: 105},
							{Title: "Direccion", Width: 180},
						},
						Model:                    cliModel,
						OnSelectedIndexesChanged: seleccionarCliente,
						OnItemActivated:          seleccionarCliente,
					},
				},
			},
			GroupBox{
				Title:   "Formulario de cliente",
				Layout:  VBox{Margins: Margins{Left: 14, Top: 10, Right: 14, Bottom: 14}},
				MinSize: Size{Width: 410},
				Children: []Widget{
					Label{Text: "1. Selecciona una fila o llena los datos para crear un cliente."},
					Composite{
						Layout: Grid{Columns: 2, Spacing: 8},
						Children: []Widget{
							Label{Text: "ID"},
							LineEdit{AssignTo: &cliID, ReadOnly: true},
							Label{Text: "Nombre"},
							LineEdit{AssignTo: &cliNombre},
							Label{Text: "Apellido"},
							LineEdit{AssignTo: &cliApellido},
							Label{Text: "Correo"},
							LineEdit{AssignTo: &cliCorreo},
							Label{Text: "Telefono"},
							LineEdit{AssignTo: &cliTelefono},
							Label{Text: "Direccion"},
							LineEdit{AssignTo: &cliDireccion},
						},
					},
					Label{Text: "2. Ejecuta una accion."},
					Composite{
						Layout: Grid{Columns: 2, Spacing: 8},
						Children: []Widget{
							PushButton{Text: "Crear nuevo", OnClicked: crearCliente},
							PushButton{Text: "Guardar cambios", OnClicked: actualizarCliente},
							PushButton{Text: "Buscar ID", OnClicked: buscarClienteID},
							PushButton{Text: "Eliminar seleccionado", OnClicked: eliminarCliente},
							PushButton{Text: "Limpiar", OnClicked: limpiarClientes},
						},
					},
				},
			},
		},
	}
}

func obtenerAPIURL() string {
	valor := strings.TrimSpace(os.Getenv("API_URL"))
	if valor == "" {
		return "http://localhost:8090"
	}
	return strings.TrimRight(valor, "/")
}

func cargarEmpleados() {
	setStatus("Cargando empleados...")
	go func() {
		var empleados []Empleado
		data, err := requestBytes("GET", "/empleados", nil)
		if err == nil {
			err = json.Unmarshal(data, &empleados)
		}
		mw.Synchronize(func() {
			if err != nil {
				mostrarError(err)
				return
			}
			empModel.setItems(empleados)
			setStatus(fmt.Sprintf("Empleados cargados: %d", len(empleados)))
			setOutput(prettyJSON(data))
		})
	}()
}

func buscarEmpleadoNombre() {
	nombre := texto(empBusqueda)
	if nombre == "" {
		mostrarAviso("Escribe un nombre o apellido en Buscar por nombre.")
		return
	}
	setStatus("Buscando empleados...")
	go func() {
		var empleados []Empleado
		data, err := requestBytes("GET", "/empleados/buscar?nombre="+url.QueryEscape(nombre), nil)
		if err == nil {
			err = json.Unmarshal(data, &empleados)
		}
		mw.Synchronize(func() {
			if err != nil {
				mostrarError(err)
				return
			}
			empModel.setItems(empleados)
			setStatus(fmt.Sprintf("Coincidencias encontradas: %d", len(empleados)))
			setOutput(prettyJSON(data))
		})
	}()
}

func buscarEmpleadoID() {
	id := texto(empID)
	if id == "" {
		mostrarAviso("Selecciona una fila o escribe un ID.")
		return
	}
	ejecutar("GET", "/empleados/"+url.PathEscape(id), nil, false)
}

func crearEmpleado() {
	if !validarEmpleado() {
		return
	}
	ejecutar("POST", "/empleados", empleadoDesdeUI(), true)
}

func actualizarEmpleado() {
	id := texto(empID)
	if id == "" {
		mostrarAviso("Selecciona un empleado antes de guardar cambios.")
		return
	}
	if !validarEmpleado() {
		return
	}
	ejecutar("PUT", "/empleados/"+url.PathEscape(id), empleadoDesdeUI(), true)
}

func eliminarEmpleado() {
	id := texto(empID)
	if id == "" {
		mostrarAviso("Selecciona el empleado que deseas eliminar.")
		return
	}
	if walk.MsgBox(mw, "Confirmar", "Eliminar empleado ID "+id+"?", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) != walk.DlgCmdYes {
		return
	}
	ejecutar("DELETE", "/empleados/"+url.PathEscape(id), nil, true)
	limpiarEmpleados()
}

func historialEmpleado(tipo string) {
	id := texto(empID)
	if id == "" {
		mostrarAviso("Selecciona un empleado para consultar historial.")
		return
	}
	ejecutar("GET", "/empleados/"+url.PathEscape(id)+"/"+tipo, nil, false)
}

func seleccionarEmpleado() {
	if empTable == nil {
		return
	}
	idx := empTable.CurrentIndex()
	if idx < 0 || idx >= len(empModel.items) {
		return
	}
	llenarEmpleado(empModel.items[idx])
}

func llenarEmpleado(e Empleado) {
	_ = empID.SetText(fmt.Sprint(e.EmpNo))
	_ = empNombre.SetText(e.FirstName)
	_ = empApellido.SetText(e.LastName)
	_ = empNacimiento.SetText(e.BirthDate)
	_ = empGenero.SetText(e.Gender)
	_ = empIngreso.SetText(e.HireDate)
}

func cargarClientes() {
	setStatus("Cargando clientes...")
	go func() {
		var clientes []Cliente
		data, err := requestBytes("GET", "/clientes", nil)
		if err == nil {
			err = json.Unmarshal(data, &clientes)
		}
		mw.Synchronize(func() {
			if err != nil {
				mostrarError(err)
				return
			}
			cliModel.setItems(clientes)
			setStatus(fmt.Sprintf("Clientes cargados: %d", len(clientes)))
			setOutput(prettyJSON(data))
		})
	}()
}

func buscarClienteNombre() {
	nombre := texto(cliBusqueda)
	if nombre == "" {
		mostrarAviso("Escribe un nombre o apellido en Buscar por nombre.")
		return
	}
	setStatus("Buscando clientes...")
	go func() {
		var clientes []Cliente
		data, err := requestBytes("GET", "/clientes/buscar?nombre="+url.QueryEscape(nombre), nil)
		if err == nil {
			err = json.Unmarshal(data, &clientes)
		}
		mw.Synchronize(func() {
			if err != nil {
				mostrarError(err)
				return
			}
			cliModel.setItems(clientes)
			setStatus(fmt.Sprintf("Coincidencias encontradas: %d", len(clientes)))
			setOutput(prettyJSON(data))
		})
	}()
}

func buscarClienteID() {
	id := texto(cliID)
	if id == "" {
		mostrarAviso("Selecciona una fila o escribe un ID.")
		return
	}
	ejecutar("GET", "/clientes/"+url.PathEscape(id), nil, false)
}

func crearCliente() {
	if !validarCliente() {
		return
	}
	ejecutar("POST", "/clientes", clienteDesdeUI(), true)
}

func actualizarCliente() {
	id := texto(cliID)
	if id == "" {
		mostrarAviso("Selecciona un cliente antes de guardar cambios.")
		return
	}
	if !validarCliente() {
		return
	}
	ejecutar("PUT", "/clientes/"+url.PathEscape(id), clienteDesdeUI(), true)
}

func eliminarCliente() {
	id := texto(cliID)
	if id == "" {
		mostrarAviso("Selecciona el cliente que deseas eliminar.")
		return
	}
	if walk.MsgBox(mw, "Confirmar", "Eliminar cliente ID "+id+"?", walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) != walk.DlgCmdYes {
		return
	}
	ejecutar("DELETE", "/clientes/"+url.PathEscape(id), nil, true)
	limpiarClientes()
}

func seleccionarCliente() {
	if cliTable == nil {
		return
	}
	idx := cliTable.CurrentIndex()
	if idx < 0 || idx >= len(cliModel.items) {
		return
	}
	c := cliModel.items[idx]
	_ = cliID.SetText(fmt.Sprint(c.ID))
	_ = cliNombre.SetText(c.Nombre)
	_ = cliApellido.SetText(c.Apellido)
	_ = cliCorreo.SetText(c.Correo)
	_ = cliTelefono.SetText(c.Telefono)
	_ = cliDireccion.SetText(c.Direccion)
}

func ejecutar(method, path string, body any, recargar bool) {
	setStatus("Consultando API...")
	setOutput("Procesando solicitud...")

	go func() {
		data, err := requestBytes(method, path, body)
		mw.Synchronize(func() {
			if err != nil {
				mostrarError(err)
				return
			}
			setStatus("Operacion completada correctamente.")
			setOutput(prettyJSON(data))
			if recargar {
				cargarEmpleados()
				cargarClientes()
			}
		})
	}()
}

func requestBytes(method, path string, body any) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, apiURL+path, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar con la API en %s: %w", apiURL, err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("respuesta HTTP %d\n%s", res.StatusCode, prettyJSON(data))
	}
	return data, nil
}

func prettyJSON(data []byte) string {
	var valor any
	if err := json.Unmarshal(data, &valor); err != nil {
		return string(data)
	}
	formateado, err := json.MarshalIndent(valor, "", "  ")
	if err != nil {
		return string(data)
	}
	return string(formateado)
}

func empleadoDesdeUI() Empleado {
	return Empleado{
		FirstName: texto(empNombre),
		LastName:  texto(empApellido),
		BirthDate: texto(empNacimiento),
		Gender:    texto(empGenero),
		HireDate:  texto(empIngreso),
	}
}

func clienteDesdeUI() Cliente {
	return Cliente{
		Nombre:    texto(cliNombre),
		Apellido:  texto(cliApellido),
		Correo:    texto(cliCorreo),
		Telefono:  texto(cliTelefono),
		Direccion: texto(cliDireccion),
	}
}

func validarEmpleado() bool {
	if texto(empNombre) == "" {
		mostrarAviso("Escribe el nombre del empleado.")
		return false
	}
	if texto(empApellido) == "" {
		mostrarAviso("Escribe el apellido del empleado.")
		return false
	}
	if texto(empNacimiento) == "" {
		mostrarAviso("Escribe la fecha de nacimiento.")
		return false
	}
	if texto(empGenero) == "" {
		mostrarAviso("Escribe genero: M, F, Masculino o Femenino.")
		return false
	}
	if texto(empIngreso) == "" {
		mostrarAviso("Escribe contratacion: YYYY-MM-DD, Hoy o Inmediata.")
		return false
	}
	return true
}

func validarCliente() bool {
	if texto(cliNombre) == "" {
		mostrarAviso("Escribe el nombre del cliente.")
		return false
	}
	if texto(cliApellido) == "" {
		mostrarAviso("Escribe el apellido del cliente.")
		return false
	}
	if texto(cliCorreo) == "" {
		mostrarAviso("Escribe el correo del cliente.")
		return false
	}
	return true
}

func limpiarEmpleados() {
	limpiar(empID, empNombre, empApellido, empNacimiento, empGenero, empIngreso, empBusqueda)
}

func limpiarClientes() {
	limpiar(cliID, cliNombre, cliApellido, cliCorreo, cliTelefono, cliDireccion, cliBusqueda)
}

func limpiar(campos ...*walk.LineEdit) {
	for _, campo := range campos {
		_ = campo.SetText("")
	}
	setOutput("Formulario limpiado.")
	setStatus("Listo.")
}

func mostrarAviso(mensaje string) {
	setStatus("Falta informacion.")
	setOutput(mensaje)
}

func mostrarError(err error) {
	setStatus("No se pudo completar la operacion.")
	setOutput("Error: " + err.Error())
}

func texto(campo *walk.LineEdit) string {
	return strings.TrimSpace(campo.Text())
}

func setStatus(text string) {
	_ = status.SetText(text)
}

func setOutput(text string) {
	_ = output.SetText(text)
}
