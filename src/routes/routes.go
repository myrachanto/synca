package routes

import (
	"sync"

	"github.com/gin-gonic/gin"
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
		{DatabaseA: "single", DatabaseAUrl: "mongodb://localhost:27017", Databaseb: "syncab", DatabasebURL: "mongodb://localhost:27017", CollectionName: "product"},
	}

	loader := load.NewloadController(load.NewloadService(load.NewloadRepo(synca[0].DatabaseA, synca[0].DatabaseAUrl, synca[0].Databaseb, synca[0].DatabasebURL, synca[0].CollectionName)))
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"https://foo.com"},
	// 	AllowMethods:     []string{"PUT", "PATCH"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return origin == "https://github.com"
	// 	},
	// 	MaxAge: 12 * time.Hour,
	// }))
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

	api := router.Group("/api")

	api.GET("/getsyncs", loader.GetAll)

	router.Run(":3500")
}

// func ApiLoader() {
// 
	// synca := []loadrepository{
	// 	{DatabaseA: "single", DatabaseAUrl: "mongodb://localhost:27017", Databaseb: "syncab", DatabasebURL: "mongodb://localhost:27017", CollectionName: "product"},
	// }

// 	loader := load.NewloadController(load.NewloadService(load.NewloadRepo(synca[0].DatabaseA, synca[0].DatabaseAUrl, synca[0].Databaseb, synca[0].DatabasebURL, synca[0].CollectionName)))
// 	e := echo.New()

// 	e.Use(middleware.Logger())
// 	e.Use(middleware.Recover())
// 	e.Use(middleware.CORS())
// 	l := len(synca)
// 	wg.Add(l)
// 	for _, sita := range synca {
// 		go func(sita loadrepository) {
// 			dg := load.NewloadRepo(sita.DatabaseA, sita.DatabaseAUrl, sita.Databaseb, sita.DatabasebURL, sita.CollectionName)
// 			dg.StartSychronization()
// 		}(sita)
// 		wg.Done()
// 	}
// 	wg.Wait()
// 	api := e.Group("/api")

// 	api.GET("/getsyncs", loader.GetAll).Name = "get-all"

//		e.Logger.Fatal(e.Start(":3500"))
//	}
