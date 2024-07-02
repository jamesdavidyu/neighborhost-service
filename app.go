package main

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jamesdavidyu/neighborhost-service/routes"
)

// go:embed templates/*
// var resources embed.FS

// var t = template.Must(template.ParseFS(resources, "templates/*"))

func main() {
	routes.Routes()
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	data := map[string]string{
	// 		"Region": os.Getenv("FLY_REGION"),
	// 	}

	// 	t.ExecuteTemplate(w, "index.html.tmpl", data)
	// })
}
