package main

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"
	access "github.com/jakemakesstuff/cf-access-fiber"
	"invoicer/db"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed scripts/*.js
var scripts embed.FS

//go:embed schema.sql
var schema []byte

// The route that is hit when a unauthenticated/invalid Access request is made.
func unauthorized(ctx *fiber.Ctx) error {
	ctx.Status(400)
	_, err := ctx.WriteString("Invalid Cloudflare Access request.")
	return err
}

func main() {
	// Load the schema.
	db.LoadSchema(schema)

	// Create the application.
	app := fiber.New(fiber.Config{
		Views: html.NewFileSystem(http.FS(templates), ".html"),
	})

	// Serve the bundled JS.
	dir, err := fs.Sub(scripts, "scripts")
	if err != nil {
		panic(err)
	}
	app.Use("/scripts", filesystem.New(filesystem.Config{
		Root:   http.FS(dir),
		Browse: true,
	}))

	// Apply the application wide Access middleware if it is not blank.
	teamDomain := os.Getenv("TEAM_DOMAIN")
	applicationAud := os.Getenv("APPLICATION_AUD")
	if !(teamDomain == "" || applicationAud == "") {
		app.Use(access.Validate(teamDomain, applicationAud, unauthorized))
	}

	// Load the template routes.
	loadTemplateRoutes(app)

	// Start the invoice loop.
	invoiceLoop()

	// Runs the server.
	log.Fatal(app.Listen(":3000"))
}
