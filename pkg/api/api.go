package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// API est la structure principale de notre package
type API struct {
	BaseURL string
}

// NewAPI crée une nouvelle instance de l'API
func NewAPI(baseURL string) *API {
	return &API{BaseURL: baseURL}
}

// ShowAPI affiche l'URL de base de l'API
func (a *API) ShowAPI() {
	fmt.Printf("API BaseURL: %s\n", a.BaseURL)
}

// Band représente un groupe de musique
type Band struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members,omitempty"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Relationship struct {
    ID            int                    `json:"id"`
    DatesLocations map[string][]string `json:"datesLocations"`
}

// GetAllBands récupère la liste complète de tous les groupes de musique
func (a *API) GetAllBands() ([]Band, error) {
    // Effectuer la requête à l'API pour récupérer tous les groupes de musique
    resp, err := http.Get(fmt.Sprintf("%s/artists", a.BaseURL))
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    // Vérifier le statut de la réponse
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
    }

    // Décode la réponse JSON en une liste de groupes de musique
    var bands []Band
    err = json.NewDecoder(resp.Body).Decode(&bands)
    if err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return bands, nil
}

func (a *API) GetBand(bandID int) (*Band, error) {
	url := fmt.Sprintf("%s/artists/%d", a.BaseURL, bandID)

	// Envoyer une requête GET à l'API
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'envoi de la requête: %v", err)
	}
	defer resp.Body.Close()

	// Vérifier le statut de la réponse
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a renvoyé un statut non-OK: %d", resp.StatusCode)
	}

	// Décode la réponse JSON en une structure Band
	var band Band
	if err := json.NewDecoder(resp.Body).Decode(&band); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de la réponse JSON: %v", err)
	}

	return &band, nil
}

func GetRelationshipData(relationID int) (*Relationship, error) {
    url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", relationID)

    // Envoyer une requête GET à l'API
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de l'envoi de la requête: %v", err)
    }
    defer resp.Body.Close()

    // Vérifier le statut de la réponse
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("l'API a renvoyé un statut non-OK: %d", resp.StatusCode)
    }

    // Décode la réponse JSON en une structure Relationship
    var relationship Relationship
    if err := json.NewDecoder(resp.Body).Decode(&relationship); err != nil {
        return nil, fmt.Errorf("erreur lors du décodage de la réponse JSON: %v", err)
    }

    return &relationship, nil
}
