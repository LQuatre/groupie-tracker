package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"

	"groupietracker.com/m/pkg/api"
)

var staticDir = os.Getenv("STATIC_DIR")

func Setup(indexPath string, apiUrl string, bands []api.Band) {
    fileServer := http.FileServer(http.Dir(staticDir))
    http.Handle("/static/", http.StripPrefix("/static", fileServer))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" || r.URL.Path == "/index.html" {
            tmpl, err := template.ParseFiles(indexPath)
            if err != nil {
                handleError(w, err)
                return
            }

            bands, err := api.NewAPI(apiUrl).GetAllBands()
            if err != nil {
                handleError(w, err)
                return
            }

            if err := tmpl.Execute(w, bands); err != nil {
                handleError(w, err)
                return
            }
        }
    })

    if apiUrl != "" {
        SetupAPIRoutes(apiUrl)
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
    })

    http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
        parts := strings.Split(r.URL.Path, "/")
        if len(parts) < 3 {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        if len(parts) == 3 {
            endpoint := parts[2]
            endpoints := map[string]string{
				"": 		 apiUrl,
                "artists":   apiUrl + "/artists",
                "locations": apiUrl + "/locations",
                "dates":     apiUrl + "/dates",
                "relation":  apiUrl + "/relation",
            }

            url, ok := endpoints[endpoint]
            if !ok {
                w.WriteHeader(http.StatusNotFound)
                return
            }

            resp, err := http.Get(url)
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
        } else if len(parts) == 4 {
            endpoint := parts[2]
            endpoints := map[string]string{
                "artists":   apiUrl + "/artists",
                "locations": apiUrl + "/locations",
                "dates":     apiUrl + "/dates",
                "relation":  apiUrl + "/relation",
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
    })
}

func handleError(w http.ResponseWriter, err error) {
    fmt.Printf("Error: %v\n", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
