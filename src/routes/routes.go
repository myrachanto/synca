package routes

import (
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myrachanto/synca/src/api/load"
)

var wg = &sync.WaitGroup{}

type loadrepository struct {
	DatabaseA      string
	DatabaseAUrl   string
	Databaseb      string
	DatabasebURL   string
	CollectionName string
}

func ApiLoader() {
	synca := []loadrepository{
		{"single", "mongodb://localhost:27017", "syncab", "mongodb://localhost:27017", "product"},
	}

	loader := load.NewloadController(load.NewloadService(load.NewloadRepo(synca[0].DatabaseA, synca[0].DatabaseAUrl, synca[0].Databaseb, synca[0].DatabasebURL, synca[0].CollectionName)))
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	l := len(synca)
	wg.Add(l)
	for _, sita := range synca {
		go func(sita loadrepository) {
			dg := load.NewloadRepo(sita.DatabaseA, sita.DatabaseAUrl, sita.Databaseb, sita.DatabasebURL, sita.CollectionName)
			dg.StartSychronization()
		}(sita)
		wg.Done()
	}
	wg.Wait()

	// _ = time.AfterFunc(time.Second*3, load.Loadrepository.StartSychronization)
	api := e.Group("/api")

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file in routes")
	// }
	api.GET("/getURL", loader.Synca).Name = "get-url"

	// PORT := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":3500"))
}
