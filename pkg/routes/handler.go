package routes

import (
	"fmt"
	"net/http"
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
	if len(parts) < 3 {
		if len(parts) == 2 && parts[1] == "api" {
			onlySendData(w, myAPI.BaseApi)
		} else {
			w.Header().Set("Content-Type", "text/html")
			handleError(w, fmt.Errorf("Invalid endpoint"))
		}
		return
	}

	endpoint := ""
	if len(parts) >= 3 {
		endpoint = parts[2]
	}

	var endpoints map[string]func()

	// if endpoint == nil {
	// 	onlySendData(w, myAPI.BaseApi)
	// } else {
	endpoints = map[string]func(){
		"":          func() { onlySendData(w, myAPI.BaseApi) },
		"artists":   func() { onlySendData(w, myAPI.Artists) },
		"locations": func() { onlySendData(w, myAPI.Locations) },
		"dates":     func() { onlySendData(w, myAPI.Dates) },
		"relation":  func() { onlySendData(w, myAPI.Relation) },
	}
	// }

	handleEndpoint, ok := endpoints[endpoint]
	if !ok {
		w.Header().Set("Content-Type", "text/html")
		handleError(w, fmt.Errorf("Invalid endpoint"))
		return
	}

	handleEndpoint()
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
