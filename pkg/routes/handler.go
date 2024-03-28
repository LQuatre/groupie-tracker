package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"groupietracker.com/m/pkg/api"
)

// Fonction pour gérer l'index et les requêtes vers /index.html
func handleIndex(indexPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		} else {
			// rediriger vers la page 404 "/404"
			http.Redirect(w, r, "/404", http.StatusFound)
			return
		}
	}
}

func handleAPIRequest(w http.ResponseWriter, myAPI *api.API, Path string, id ...string) {
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(Path, "/")

	endpoint := parts[2]
	url := myAPI.BaseURL

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

	// recupérer les données de l'API
	data := struct {
		Artists   []api.Band
		Locations []api.IndexLocations
		Dates     []api.IndexDates
		Relation  []api.Relation
	}{
		Artists:   myAPI.Artists,
		Locations: myAPI.Locations,
		Dates:     myAPI.Dates,
		Relation:  myAPI.Relation,
	}
	
	// Afficher les données de l'API en fonction de l'endpoint
	if id == nil {
		switch endpoint {
		case "artists":
			onlySendData(w, data.Artists)
		case "locations":
			onlySendData(w, data.Locations)
		case "dates":
			onlySendData(w, data.Dates)
		case "relation":
			onlySendData(w, data.Relation)
		}
	} else {
		switch endpoint {
		case "artists":
			for _, band := range data.Artists {
				id, _ := strconv.Atoi(id[0])
				if band.ID == id {
					onlySendData(w, band)
					return
				}
			}
		case "locations":
			for _, location := range data.Locations {
				locationID, _ := strconv.Atoi(id[0])
				if location.ID == locationID {
					onlySendData(w, location)
					return
				}
			}
		case "dates":
			for _, date := range data.Dates {
				dateID, _ := strconv.Atoi(id[0])
				if date.ID == dateID {
					onlySendData(w, date)
					return
				}
			}
		case "relation":
			for _, relation := range data.Relation {
				relationID, _ := strconv.Atoi(id[0])
				if relation.ID == relationID {
					onlySendData(w, relation)
					return
				}
			}
		}
	}
}

func handleAPIEndpointRequest(w http.ResponseWriter, r *http.Request, myAPI *api.API) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(parts) == 3 {
		handleAPIRequest(w, myAPI, r.URL.Path)
	} else if len(parts) == 4 {
		id := parts[3]
		handleAPIRequest(w, myAPI, r.URL.Path, id)

	}
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Printf("Error: %v\n", err)
	
	tmpl, tmplErr := template.ParseFiles("web/template/error.html")
	if tmplErr != nil {
		http.Error(w, "Internal Server Error, You should not see this message", http.StatusInternalServerError)
		return
	}

	data := struct {
		Err string
	}{
		Err: err.Error(),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error, You should not see this message", http.StatusInternalServerError)
		return
	}
}
