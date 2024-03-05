package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"groupietracker.com/m/pkg/api"
)

var staticDir = os.Getenv("STATIC_DIR")

func Setup(indexPath string, apiUrl string, myApi *api.API) {
	fileServer := http.FileServer(http.Dir(staticDir + "web/static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			tmpl, err := template.ParseFiles(indexPath)
			if err != nil {
				handleError(w, err)
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				handleError(w, err)
				return
			}
		}
	})

	err := SetupAPIRoutes(apiUrl)
	if err != nil {
		fmt.Printf("Erreur lors de la configuration des routes de l'API: %v\n", err)
		return
	}
	err = SetSearchRoutes(myApi)
	if err != nil {

		return
	}
	err = SetFilterRoutes(myApi)
	if err != nil {
		// handle error
		return
	}
	err = SetArtistsRoutes(myApi)
	if err != nil {
		// handle error
		return
	}

	fmt.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erreur lors du démarrage du serveur: %v\n", err)
	}
}

func SetupAPIRoutes(apiUrl string) error {
	if apiUrl == "" {
		return fmt.Errorf("API URL is required")
	}
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		handleAPIRequest(w, apiUrl)
	})

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		handleAPIEndpointRequest(w, r, apiUrl)
	})
	return nil
}

func handleAPIRequest(w http.ResponseWriter, apiUrl string) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := http.Get(apiUrl)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		handleError(w, err)
		return
	}
}

func handleAPIEndpointRequest(w http.ResponseWriter, r *http.Request, apiUrl string) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	endpoint := parts[2]
	url := apiUrl

	if len(parts) == 3 {
		endpoints := map[string]string{
			"":          url,
			"artists":   url + "/artists",
			"locations": url + "/locations",
			"dates":     url + "/dates",
			"relation":  url + "/relation",
		}

		url, ok := endpoints[endpoint]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		handleAPIRequest(w, url)
	} else if len(parts) == 4 {
		endpoints := map[string]string{
			"artists":   url + "/artists",
			"locations": url + "/locations",
			"dates":     url + "/dates",
			"relation":  url + "/relation",
		}

		url, ok := endpoints[endpoint]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		id := parts[3]

		resp, err := http.Get(url + "/" + id)
		if err != nil {
			handleError(w, err)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			handleError(w, err)
			return
		}
	}
}

func SetSearchRoutes(api *api.API) error {
	if api == nil {
		return fmt.Errorf("API is required")
	}
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			query := r.URL.Query().Get("query")
			bands, err := api.GetBandFromSearch(query)
			if err != nil {
				handleError(w, err)
				return
			}

			renderTemplate(w, "web/template/galery.html", bands)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	return nil
}

func SetFilterRoutes(myapi *api.API) error {
	if myapi == nil {
		return fmt.Errorf("API is required")
	}
	http.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			members := r.FormValue("members")
			numberOfMembers := r.FormValue("numberofmember")
			location := r.FormValue("location")
			createDate := r.FormValue("creation-date")
			firstAlbum := r.FormValue("first-album")
			concertDate := r.FormValue("concert-date")

			var err error

			var numberOfMembersInt int
			if numberOfMembers != "" {
				numberOfMembersInt, err = strconv.Atoi(numberOfMembers)
				if err != nil {
					handleError(w, err)
					return
				}
			}

			var createDateInt int
			if createDate != "" {
				createDateInt, err = strconv.Atoi(createDate)
				if err != nil {
					handleError(w, err)
					return
				}
			}

			filteredBands, err := myapi.FilterBands(api.Filter{
				Members:         members,
				NumberOfMembers: numberOfMembersInt,
				Location:        location,
				CreationDate:    createDateInt,
				FirstAlbum:      firstAlbum,
				ConcertDate:     concertDate,
			})
			if err != nil {
				handleError(w, err)
				return
			}

			renderTemplate(w, "web/template/galery.html", filteredBands)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	return nil
}

func SetArtistsRoutes(myapi *api.API) error {
	if myapi == nil {
		return fmt.Errorf("API is required")
	}
	http.HandleFunc("/artists/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 2 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if parts[2] == "" {
			bands, err := myapi.GetAllBands()
			if err != nil {
				handleError(w, err)
				return
			}
			renderTemplate(w, "web/template/galery.html", bands)
		}
		if parts[2] != "" {
			id := parts[2]
			idInt, err := strconv.Atoi(id)
			if err != nil {
				handleError(w, err)
				return
			}
			band, err := myapi.GetBand(idInt)
			if err != nil {
				handleError(w, err)
				return
			}
			stringListLocations, err := myapi.GetRelation(idInt)
			if err != nil {
				handleError(w, err)
				return
			}

			band.LocationsCoordinates = []api.Location{}

			for key, value := range stringListLocations.DatesLocations {
				lat, lng := GeocodeAddress(key)
				thisLocation := api.Location{Lat: lat, Lng: lng, Dates: value}
				band.LocationsCoordinates = append(band.LocationsCoordinates, thisLocation)
			}
			renderTemplate(w, "web/template/artist.html", band)
		}
	})
	return nil
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Printf("Error: %v\n", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		handleError(w, err)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		handleError(w, err)
		return
	}
}

func GeocodeAddress(placeName string) (float64, float64) {
	apiKey := "AIzaSyAX7_r2A6VAL2v8gKKnZmXmD1Z2bEdov2o"

	// Construire l'URL de l'API de géocodage de Google
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", placeName, apiKey)

	// Effectuer la requête HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return 0.0, 0.0
	}
	defer resp.Body.Close()

	// Vérifier le code de statut HTTP
	if resp.StatusCode != http.StatusOK {
		fmt.Println("La requête HTTP a retourné un code d'état non-OK:", resp.StatusCode)
		return 0.0, 0.0
	}

	// Décodez la réponse JSON
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Erreur lors du décodage de la réponse JSON:", err)
		return 0.0, 0.0
	}

	// Vérifiez le statut de la réponse
	if status, ok := response["status"].(string); !ok || status != "OK" {
		fmt.Println("La réponse de l'API n'est pas OK:", status)
		return 0.0, 0.0
	}

	lat := 0.0
	lng := 0.0

	// fmt.Println(response)

	// Récupérer les coordonnées géographiques (latitude et longitude)
	results := response["results"].([]interface{})
	if len(results) > 0 {
		geometry := results[0].(map[string]interface{})["geometry"].(map[string]interface{})
		location := geometry["location"].(map[string]interface{})
		lat = location["lat"].(float64)
		lng = location["lng"].(float64)
		fmt.Printf("Coordonnées de %s: Latitude %f, Longitude %f\n", placeName, lat, lng)
	} else {
		fmt.Println("Aucun résultat trouvé pour", placeName)
	}

	return lat, lng
}
