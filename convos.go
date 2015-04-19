package main

import (
	"fmt"
	"log"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/nt3rp/convos/handlers"
)

var (
	httpPort int = 8080
)

func main() {
	m := martini.Classic()

	// Add additional middleware
	m.Use(render.Renderer())
	// TODO: Add "Authorization" middleware that restricts user access to their convos

	// Define Routes
	// TODO: Fix trailing slashes
	m.Group("/convos", func(r martini.Router) {
		r.Get("/", handlers.GetConvos)
		r.Get("/:id/", handlers.GetConvo)
		r.Post("/", handlers.CreateConvo)
		r.Patch("/:id/", handlers.UpdateConvo)  // Generally, edits aren't allowed. Restrict this to updating the status...
		r.Delete("/:id/", handlers.DeleteConvo)
	})

	log.Printf("listening on %v\n", httpPort)
	httpAddr := fmt.Sprintf(":%d", httpPort)
	m.RunOnAddr(httpAddr)
}
