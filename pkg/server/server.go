package server

import (
	"log"

	"groupietracker.com/m/pkg/api"
	"groupietracker.com/m/pkg/routes"
)

func StartServer() {
	apiUrl := "https://groupietrackers.herokuapp.com/api"
	localApiUrl := apiUrl

	// Créer une nouvelle instance de l'API avec l'URL de base
	myAPI := api.NewAPI(localApiUrl)

	// Afficher l'URL de base de l'API
	myAPI.ShowAPI()

	// Récupérer la liste complète de tous les groupes de musique
	bands, err := myAPI.GetAllBands()
	if err != nil {
		log.Fatalf("Failed to get all bands: %v", err)
	}

	// Afficher les détails de chaque groupe de musique
	// for _, band := range bands {
	// 	fmt.Printf("ID: %d\n", band.ID)
	// 	fmt.Printf("Nom: %s\n", band.Name)
	// 	fmt.Printf("Image: %s\n", band.Image)
	// 	fmt.Printf("Membres: %v\n", band.Members)
	// 	fmt.Printf("Date de création: %d\n", band.CreationDate)
	// 	fmt.Printf("Premier album: %s\n", band.FirstAlbum)
	// 	fmt.Printf("Locations: %s\n", band.Locations)
	// 	fmt.Printf("Concert Dates: %s\n", band.ConcertDates)
	// 	fmt.Printf("Relations: %s\n", band.Relations)
	// 	fmt.Println("----------------------------------------")
	// }

	routes.Setup("web/template/index.html", apiUrl, bands)

	// Démarrer le serveur
	routes.Run()
}

