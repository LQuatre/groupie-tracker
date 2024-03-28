package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"groupietracker.com/m/pkg/api"
	userGestion "groupietracker.com/m/pkg/user"
)

var staticDir = os.Getenv("STATIC_DIR")

type DataToWeb struct {
	Bands          []api.Band
	UserIsLoggedIn bool
	Username       string
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

// Configuration des routes API et autres routes
func setupRoutes(apiUrl string, myApi *api.API) {
	routes := []struct {
		route func(string) error
	}{
		{func(s string) error { return Setup404Route() }},
		{func(s string) error { return SetupErrorRoute() }},

		{func(s string) error { return SetAPIRoutes(myApi) }},
		{func(s string) error { return SetSearchRoutes(myApi) }},
		{func(s string) error { return SetArtistsRoutes(myApi) }},
		{func(s string) error { return SetLoginRoutes(myApi) }},
		{func(s string) error { return SetRegisterRoutes(myApi) }},
		{func(s string) error { return SetLogoutRoutes(myApi) }},
		{func(s string) error { return SetProfileRoutes(myApi) }},
		{func(s string) error { return SetGetArtistNamesRoute(myApi) }},
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
		//	handleAPIRequest(w, myApi, r.URL.Path)
	})

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		//	handleAPIEndpointRequest(w, r, myApi)
	})
	return nil
}

var redirectNeeded bool = false
var redirected bool = false

func SetSearchRoutes(myapi *api.API) error {
	if myapi == nil {
		return fmt.Errorf("API is required")
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Filtrer les paramètres de requête vides et le paramètre submit
		filteredParams := make(url.Values)
		for key, values := range r.URL.Query() {
			if key != "submit" && len(values) > 0 {
				filteredParams[key] = values
			}
		}

		// Générer l'URL sans les paramètres de requête vides et le paramètre submit
		cleanURL := "/search"
		if len(filteredParams) > 0 {
			cleanURL += "?" + filteredParams.Encode()
		}

		query := r.URL.Query().Get("query")
		submit := r.URL.Query().Has("submit")

		if submit {
			// fmt.Println("Redirection vers", cleanURL)
			redirectNeeded = true
			http.Redirect(w, r, cleanURL, http.StatusFound)
			return
		} else {
			redirected = false
		}

		dataToWeb := DataToWeb{}

		hasFilter := len(filteredParams) > 0

		if hasFilter {
			// Extrait les valeurs des filtres
			members := r.FormValue("members")
			numberOfMembers := r.FormValue("numberofmember")
			location := r.FormValue("location")
			createDate := r.FormValue("creation-date")
			firstAlbum := r.FormValue("first-album")
			concertDate := r.FormValue("concert-date")

			// Convertit les valeurs nécessaires en entiers
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

			// Filtre les groupes en fonction des paramètres
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

			dataToWeb.Bands = filteredBands
		}

		strToInt := func(s string) int {
			i, err := strconv.Atoi(s)
			if err != nil {
				return 0
			}
			return i
		}

		if query == "" && !hasFilter {
			bands, err := myapi.GetAllBands()
			if err != nil {
				handleError(w, err)
				return
			}
			dataToWeb.Bands = bands
		} else if query != "" && hasFilter {
			filteredBands, err := myapi.FilterBands(api.Filter{
				Members:         r.FormValue("members"),
				NumberOfMembers: strToInt(r.FormValue("numberofmember")),
				Location:        r.FormValue("location"),
				CreationDate:    strToInt(r.FormValue("creation-date")),
				FirstAlbum:      r.FormValue("first-album"),
				ConcertDate:     r.FormValue("concert-date"),
			})
			bands, err := myapi.GetBandFromSearchWithBands(query, filteredBands)
			if err != nil {
				handleError(w, err)
				return
			}
			dataToWeb.Bands = bands
		} else if query != "" && !hasFilter {
			bands, err := myapi.FilterBands(api.Filter{})
			if err != nil {
				handleError(w, err)
				return
			}
			dataToWeb.Bands = bands
		}

		// Gestion de l'authentification
		if _, err := r.Cookie("loggedIn"); err == nil {
			dataToWeb.UserIsLoggedIn = true
			if cookie, err := r.Cookie("username"); err == nil {
				dataToWeb.Username = cookie.Value
			}
		}

		// Affiche la page appropriée en fonction de la nécessité de redirection
		if redirectNeeded || redirected {
			// fmt.Println("Afficahge de la page de redirection")
			redirectNeeded = false
			redirected = true
			renderTemplate(w, "web/template/galery.html", dataToWeb)
			return
		}

		// fmt.Println("Affichage des recherche")
		renderTemplate(w, "web/template/search.html", dataToWeb)
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
	Username        string `json:"username"`
	Password        string `json:"password"`
	Mail            string `json:"mail"`
	Starred         string `json:"starred"`
	Grade           string `json:"grade"`
	ErrUserNotFound string `json:"errUserNotFound"`
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
			user, errTxt := userGestion.Login(w, username, password)
			if errTxt != "" {
				renderTemplate(w, "web/template/login.html", UserStruct{ErrUserNotFound: errTxt})
			}
			if user != (userGestion.UserStruct{}) {
				// Ajouter un cookie indiquant la connexion réussie
				http.SetCookie(w, &http.Cookie{
					Name:  "loggedIn",
					Value: "true",
					Path:  "/",
				})
				http.SetCookie(w, &http.Cookie{
					Name:  "username",
					Value: user.Username,
					Path:  "/",
				})
				http.Redirect(w, r, "/artists", http.StatusFound)
				// fmt.Println("User logged in successfully.")
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
			mail := r.FormValue("email")
			thisuser, errTxt := userGestion.Register(username, password, mail)
			if errTxt != "" {
				renderTemplate(w, "web/template/register.html", UserStruct{ErrUserNotFound: errTxt})
				return
			}
			if thisuser != (userGestion.UserStruct{}) {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			http.Redirect(w, r, "/register", http.StatusFound)
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
			user, err := userGestion.GetUser(username)
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

func SetupAdminRoutes(myapi *api.API) error {
	if myapi == nil {
		return fmt.Errorf("API is required")
	}
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
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
			user, err := userGestion.GetUser(username)
			if err != nil {
				handleError(w, err)
				return
			}

			// Vérifier si l'utilisateur est un administrateur
			if user.Grade != "admin" {
				// Rediriger l'utilisateur vers la page 404
				http.Redirect(w, r, "/404", http.StatusFound)
				return
			}

			// Récuperer les informations de tous les utilisateurs
			users, err := userGestion.GetAllUsers()
			if err != nil {
				handleError(w, err)
				return
			}

			artists, err := myapi.GetAllBands()
			if err != nil {
				handleError(w, err)
				return
			}

			var dataAdmin struct {
				Users   []userGestion.UserStruct
				Artists []api.Band
			}
			dataAdmin.Users = users
			dataAdmin.Artists = artists

			// Afficher la page d'administration
			renderTemplate(w, "web/template/admin.html", dataAdmin)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	return nil
}

func SetGetArtistNamesRoute(myApi *api.API) error {
	http.HandleFunc("/get-artist-names", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Récupérer les noms d'artistes depuis votre API Go
			bands, err := myApi.GetAllBands()
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Extraire les noms des artistes
			var artistNames []string
			for _, band := range bands {
				artistNames = append(artistNames, band.Name)
			}

			// Convertir les noms d'artistes en JSON
			jsonResponse, err := json.Marshal(artistNames)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Renvoyer la réponse JSON
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	return nil
}
