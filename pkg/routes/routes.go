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
	fileServer := http.FileServer(http.Dir(staticDir))
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

	if apiUrl != "" {
		SetupAPIRoutes(apiUrl)
		SetSearchRoutes(myApi)
		SetFilterRoutes(myApi)
		SetArtistsRoutes(myApi)
	}

	go http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
	fmt.Println("Server started at https://localhost:443")
}

func Run() {
	fmt.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erreur lors du démarrage du serveur: %v\n", err)
	}
}

func SetupAPIRoutes(apiUrl string) {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		handleAPIRequest(w, apiUrl)
	})

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		handleAPIEndpointRequest(w, r, apiUrl)
	})
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

func SetSearchRoutes(api *api.API) {
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
}

func SetFilterRoutes(myapi *api.API) {
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
}

func SetArtistsRoutes(myapi *api.API) {
	http.HandleFunc("/artists/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 2 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if len(parts) == 3 && parts[2] == "" {
			bands, err := myapi.GetAllBands()
			if err != nil {
				handleError(w, err)
				return
			}

			renderTemplate(w, "web/template/galery.html", bands)
		}

		if len(parts) == 3 && parts[2] != "" {
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

			renderTemplate(w, "web/template/artist.html", band)
		}
	})
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

func handleAPILocationsRequest(w http.ResponseWriter, apiUrl string) {
    w.Header().Set("Content-Type", "application/json")
    resp, err := http.Get(apiUrl)
    if err != nil {
        handleError(w, err)
        return
    }
    defer resp.Body.Close()

    var locations []Location
    if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
        handleError(w, err)
        return
    }

    // Transform the data to extract latitude and longitude
    var markers []Marker
    for _, location := range locations {
        for _, dates := range location.DatesLocations {
            for _, date := range dates {
                lat, lng := extractCoordinates(location.ID)
                markers = append(markers, Marker{
                    Latitude:  lat,
                    Longitude: lng,
                    Date:      date,
                })
            }
        }
    }

    // Send the transformed data back as JSON
    if err := json.NewEncoder(w).Encode(markers); err != nil {
        handleError(w, err)
        return
    }
}

// Define structs to match the JSON data
type Location struct {
    ID             int                    `json:"id"`
    DatesLocations map[string][]string   `json:"datesLocations"`
}

type Marker struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Date      string  `json:"date"`
}

// Function to extract coordinates based on location ID
func extractCoordinates(id int) (float64, float64) {
    // Define coordinates mapping based on location ID
    // You need to define this mapping based on your data
	// Declare latitude and longitude variables
	latitude := 0.0
	longitude := 0.0

	coordinates := map[int][2]float64{
		1:  {latitude, longitude}, // Example: {40.7128, -74.0060}
		2:  {latitude, longitude},
		// Add mappings for other location IDs
	}

	// Retrieve coordinates based on location ID
	if coord, ok := coordinates[id]; ok {
		return coord[0], coord[1]
	}

    // Default to (0, 0) if no coordinates are found
    return 0, 0
}
