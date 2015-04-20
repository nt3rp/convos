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

	// Define Routes
	m.Group("/convos", func(r martini.Router) {
		r.Get("/", handlers.GetConvos)
		r.Post("/", handlers.CreateConvo)
		r.Get("/:id/", handlers.GetConvo)
		r.Patch("/:id/", handlers.UpdateConvo)
		r.Delete("/:id/", handlers.DeleteConvo)
		r.Post("/:id/reply/", handlers.CreateConvo)
	}, handlers.UserAuthorizationMiddleware)

	log.Printf("listening on %v\n", httpPort)
	httpAddr := fmt.Sprintf(":%d", httpPort)
	m.RunOnAddr(httpAddr)
}
