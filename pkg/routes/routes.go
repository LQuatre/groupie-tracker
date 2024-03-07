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
	"groupietracker.com/m/pkg/user"
)

var staticDir = os.Getenv("STATIC_DIR")

type DataToWeb struct {
	Bands []api.Band
	UserIsLoggedIn bool
	Username string
}


func Setup(indexPath string, apiUrl string, myApi *api.API) {
	// Configuration du serveur de fichiers statiques
	fileServer := http.FileServer(http.Dir(staticDir + "web/static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Configuration des routes
	http.HandleFunc("/", handleIndex(indexPath))
	setupRoutes(apiUrl, myApi)

	// Démarrage du serveur
	fmt.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erreur lors du démarrage du serveur: %v\n", err)
	}
}

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
		}
	}
}

// Configuration des routes API et autres routes
func setupRoutes(apiUrl string, myApi *api.API) {
	routes := []struct {
		route func(string) error
	}{
		{func(s string) error { return SetAPIRoutes(apiUrl) }},
		{func(s string) error { return SetSearchRoutes(myApi) }},
		{func(s string) error { return SetFilterRoutes(myApi) }},
		{func(s string) error { return SetArtistsRoutes(myApi) }},
		{func(s string) error { return SetLoginRoutes(myApi) }},
		{func(s string) error { return SetRegisterRoutes(myApi) }},
		{func(s string) error { return SetLogoutRoutes(myApi) }},
		{func(s string) error { return SetProfileRoutes(myApi) }},
	}

	for _, r := range routes {
		if err := r.route(apiUrl); err != nil {
			fmt.Printf("Erreur lors de la configuration des routes: %v\n", err)
			return
		}
	}
}


func SetAPIRoutes(apiUrl string) error {
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

func SetSearchRoutes(myapi *api.API) error {
    if myapi == nil {
        return fmt.Errorf("API is required")
    }
    http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            query := r.URL.Query().Get("query")

            if query != "" {
                bands, err := myapi.GetBandFromSearch(query)
                if err != nil {
                    handleError(w, err)
                    return
                }

                dataToWeb := DataToWeb{Bands: bands}
                _, err = r.Cookie("loggedIn")
                if err != nil {
                    dataToWeb.UserIsLoggedIn = false
                } else {
                    dataToWeb.UserIsLoggedIn = true
                    cookie, err := r.Cookie("username")
                    if err != nil {
                        handleError(w, err)
                        return
                    }
                    dataToWeb.Username = cookie.Value
                }

                // If the request comes from a "submit" button on the search page, redirect to the gallery page
				// verifie s'il y a un search dans l'url et si oui, renvoie la page de recherche
				// if r.FormValue("submit") == "search" { ne fonctionne pas
				fmt.Println(r.URL.Query())
				if r.URL.Query().Has("submit") {
					fmt.Println("galery.html")
					renderTemplate(w, "web/template/galery.html", dataToWeb)
					http.Redirect(w, r, "/search", http.StatusFound)
					return
				}
				fmt.Println("search.html")

                renderTemplate(w, "web/template/search.html", dataToWeb)
                return
            }

            // Si aucune requête de recherche n'est effectuée, renvoyer simplement la page d'accueil
            bands, err := myapi.GetAllBands()
            if err != nil {
                handleError(w, err)
                return
            }
            data := DataToWeb{Bands: bands}

            _, err = r.Cookie("loggedIn")
            if err != nil {
                data.UserIsLoggedIn = false
            } else {
                data.UserIsLoggedIn = true
                cookie, err := r.Cookie("username")
                if err != nil {
                    handleError(w, err)
                    return
                }
                data.Username = cookie.Value
            }

            // Rendre le modèle HTML avec les données
            renderTemplate(w, "web/template/search.html", data)
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

			data := DataToWeb{Bands: filteredBands}
			_, err = r.Cookie("loggedIn")
			if err != nil {
				data.UserIsLoggedIn = false
			} else {
				data.UserIsLoggedIn = true
				cookie, err := r.Cookie("username")
				if err != nil {
					handleError(w, err)
					return
				}
				data.Username = cookie.Value
			}
			renderTemplate(w, "web/template/galery.html", data)
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
			data := DataToWeb{Bands: bands}
			_, err = r.Cookie("loggedIn")
			if err != nil {
				data.UserIsLoggedIn = false
			} else {
				data.UserIsLoggedIn = true
				cookie, err := r.Cookie("username")
				if err != nil {
					handleError(w, err)
					return
				}
				data.Username = cookie.Value
			}
			renderTemplate(w, "web/template/galery.html", data)
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

type UserStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Mail     string `json:"mail"`
	Starred  string `json:"starred"`
	Grade    string `json:"grade"`
}

func SetLoginRoutes(myapi *api.API) error {
	if myapi == nil {
		return fmt.Errorf("API is required")
	}
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			renderTemplate(w, "web/template/login.html", nil)
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				handleError(w, err)
				return
			}
			username := r.FormValue("username")
			password := r.FormValue("password")
			user.SetMySQL()
			thisuser, err := user.Login(w, username, password)
			fmt.Println(thisuser)
			if err != nil {
				handleError(w, err)
				return
			}
			if thisuser != (user.UserStruct{}) {
				// Ajouter un cookie indiquant la connexion réussie
				http.SetCookie(w, &http.Cookie{
					Name:  "loggedIn",
					Value: "true",
					Path:  "/",
				})
				http.SetCookie(w, &http.Cookie{
					Name:  "username",
					Value: thisuser.Username,
					Path:  "/",
				})
				http.Redirect(w, r, "/artists", http.StatusFound)
				fmt.Println("User logged in successfully.")
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	return nil
}

func SetRegisterRoutes(myapi *api.API) error {
	if myapi == nil {
		return fmt.Errorf("API is required")
	}
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			renderTemplate(w, "web/template/register.html", nil)
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				handleError(w, err)
				return
			}
			username := r.FormValue("username")
			password := r.FormValue("password")
			mail := r.FormValue("mail")
			user.SetMySQL()
			thisuser, err := user.Register(username, password, mail)
			fmt.Println(thisuser)
			if err != nil {
				handleError(w, err)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	return nil
}

func SetLogoutRoutes(myapi *api.API) error {
	if myapi == nil {
		return fmt.Errorf("API is required")
	}
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Supprimer le cookie indiquant la connexion réussie
			http.SetCookie(w, &http.Cookie{
				Name:   "loggedIn",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/artists/", http.StatusFound)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	return nil
}

func SetProfileRoutes(myapi *api.API) error {
	if myapi == nil {
		return fmt.Errorf("API is required")
	}
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Vérifier si l'utilisateur est connecté
			_, err := r.Cookie("loggedIn")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Récupérer le nom d'utilisateur à partir du cookie
			cookie, err := r.Cookie("username")
			if err != nil {
				handleError(w, err)
				return
			}
			username := cookie.Value

			// Récupérer les informations de l'utilisateur à partir de l'API
			user, err := user.GetUser(username)
			if err != nil {
				handleError(w, err)
				return
			}

			// Afficher le profil de l'utilisateur
			renderTemplate(w, "web/template/profile.html", user)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

func onlySendData(w http.ResponseWriter, data DataToWeb) {
    // Convertir les données en JSON
    jsonData, err := json.Marshal(data)
    if err != nil {
        handleError(w, err)
        return
    }

    // Définir le type de contenu de la réponse comme JSON
    w.Header().Set("Content-Type", "application/json")

    // Envoyer les données JSON en réponse
    _, err = w.Write(jsonData)
    if err != nil {
        handleError(w, err)
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
