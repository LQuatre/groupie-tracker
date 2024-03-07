package server

import (
	"groupietracker.com/m/pkg/api"
	"groupietracker.com/m/pkg/routes"
	"groupietracker.com/m/pkg/user"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func StartServer() {
	apiUrl := "https://groupietrackers.herokuapp.com/api"
	localApiUrl := apiUrl

	myAPI := api.NewAPI(localApiUrl)

	myAPI.ShowAPI()

	user.SetMySQL()

	routes.Setup("web/template/index.html", apiUrl, myAPI)
}