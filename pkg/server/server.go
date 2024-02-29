package server

import (
	"log"

	"groupietracker.com/m/pkg/api"
	"groupietracker.com/m/pkg/routes"
)

func StartServer() {
	apiUrl := "https://groupietrackers.herokuapp.com/api"
	localApiUrl := apiUrl

	myAPI := api.NewAPI(localApiUrl)

	myAPI.ShowAPI()

	bands, err := myAPI.GetAllBands()
	if err != nil {
		log.Fatalf("Failed to get all bands: %v", err)
	}

	routes.Setup("web/template/index.html", apiUrl, bands, myAPI)

	routes.Run()
}

