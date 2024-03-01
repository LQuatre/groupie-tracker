package server

import (
	"groupietracker.com/m/pkg/api"
	"groupietracker.com/m/pkg/routes"
)

func StartServer() {
	apiUrl := "https://groupietrackers.herokuapp.com/api"
	localApiUrl := apiUrl

	myAPI := api.NewAPI(localApiUrl)

	myAPI.ShowAPI()

	routes.Setup("web/template/index.html", apiUrl, myAPI)

	routes.Run()
}
