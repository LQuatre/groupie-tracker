package routes

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
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
