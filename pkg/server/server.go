package server

import (
	"groupietracker.com/m/pkg/api"
	"groupietracker.com/m/pkg/routes"
	userGestion "groupietracker.com/m/pkg/user"

	_ "github.com/go-sql-driver/mysql"
)

func StartServer() {
	apiUrl := "https://groupietrackers.herokuapp.com/api"

	myAPI := api.NewAPI(apiUrl)
	myAPI.ShowAPI()

	userGestion.SetMySQL()

	routes.Setup("web/template/index.html", apiUrl, myAPI)
}
