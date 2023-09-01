# Panadería API REST con Fiber

Este es un proyecto de ejemplo de una API REST para una panadería, construido con Go y el framework Fiber.

## Pre-requisitos

- [Go](https://golang.org/dl/) (versión 1.16 o superior)
- [SQLite](https://sqlite.org/download.html) (opcional, si deseas configurar la base de datos manualmente)

## Instalación

### Clonar el repositorio

#### En Linux

```bash
git clone https://github.com/tu-usuario/panaderia-api-rest-fiber.git
cd panaderia-api-rest-fiber
```

#### En Windows

Abre la terminal y ejecuta:

```cmd
git clone https://github.com/tu-usuario/panaderia-api-rest-fiber.git
cd panaderia-api-rest-fiber
```

### Configuración de la base de datos

Para el **primer inicio**, asegúrate de descomentar las líneas 39 y 40 en el archivo `main.go`:

```go
// Descomentar las 2 líneas de código para el primer uso
//database.SetupDB()
//database.SetupProductAndCartTables()
```

Este paso es crucial para la inicialización de la base de datos SQLite.

### Ejecutar la aplicación

#### En Linux

```bash
go run main.go
```

#### En Windows

```cmd
go run main.go
```

## Uso

Una vez que la aplicación esté en funcionamiento, podrás acceder a ella a través de `http://localhost:3000`.
