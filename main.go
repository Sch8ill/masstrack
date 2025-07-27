package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sch8ill/masstrack/db"
	"github.com/sch8ill/masstrack/handlers"
	"github.com/sch8ill/masstrack/location"

	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload"
)

//go:embed static
var assets embed.FS

//go:embed templates
var templates embed.FS

func main() {
	r := gin.Default()

	if gin.Mode() == gin.DebugMode {
		r.LoadHTMLGlob("templates/*")
		r.Static("/static", "./static")

	} else {
		t, err := template.ParseFS(templates, "templates/*")
		if err != nil {
			panic(err)
		}
		r.SetHTMLTemplate(t)

		r.GET("/static/*path", gin.WrapH(http.FileServer(http.FS(assets))))

		r.SetTrustedProxies([]string{"127.0.0.1", "::1", "172.17.0.0/16"})
	}

	db, err := db.NewDB(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/", handlers.Index)
	api := r.Group("/api/v1")
	api.GET("/locations", handlers.Locations(db))

	if os.Getenv("MASSTRACK_COLLECT") != "false" {
		log.Println("running location collector")
		go locationsService(db)
	}

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}

func locationsService(db *db.DB) {
	for {
		locations, err := location.FetchLocations()
		if err != nil {
			log.Fatal(err)
		}

		for _, l := range locations {
			if err := db.NewLocation(l); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(time.Minute * 10)
	}
}
