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

	if len(parts) < 2 || parts[1] != "api" {
		handleError(w, fmt.Errorf("Invalid API path: %v", Path))
		return
	}

	endpoint := ""
	if len(parts) >= 3 {
		endpoint = parts[2]
	}

	switch endpoint {
	case "":
		onlySendData(w, myAPI.BaseApi)
	case "artists", "locations", "dates", "relation":
		if len(id) == 0 {
			sendDataByEndpoint(w, myAPI, endpoint)
		} else {
			sendDataByID(w, myAPI, endpoint, id[0])
		}
	default:
		handleError(w, fmt.Errorf("Invalid API endpoint: %v", endpoint))
	}
}

func sendDataByEndpoint(w http.ResponseWriter, myAPI *api.API, endpoint string) {
	switch endpoint {
	case "artists":
		onlySendData(w, myAPI.Artists)
	case "locations":
		onlySendData(w, myAPI.Locations)
	case "dates":
		onlySendData(w, myAPI.Dates)
	case "relation":
		onlySendData(w, myAPI.Relation)
	}
}

func sendDataByID(w http.ResponseWriter, myAPI *api.API, endpoint, id string) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, fmt.Errorf("Invalid ID: %v", id))
		return
	}

	switch endpoint {
	case "artists":
		for _, band := range myAPI.Artists {
			if band.ID == idInt {
				onlySendData(w, band)
				return
			}
		}
	case "locations":
		for _, location := range myAPI.Locations {
			if location.ID == idInt {
				onlySendData(w, location)
				return
			}
		}
	case "dates":
		for _, date := range myAPI.Dates {
			if date.ID == idInt {
				onlySendData(w, date)
				return
			}
		}
	case "relation":
		for _, relation := range myAPI.Relation {
			if relation.ID == idInt {
				onlySendData(w, relation)
				return
			}
		}
	}

	handleError(w, fmt.Errorf("ID not found: %v", id))
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
