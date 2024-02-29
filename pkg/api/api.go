package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// API est la structure principale de notre package
type API struct {
	BaseURL       string
	Bands         []Band
	Relationships []Relationship
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

type Filter struct {
	Members         string   `json:"members,omitempty"`
	NumberOfMembers int      `json:"numberOfMembers"`
	Location        string   `json:"location"`
	CreationDate    int   `json:"creationDate"`
	FirstAlbum      string   `json:"firstAlbum"`
	ConcertDate    	string `json:"concertDate"`
}

type Relationship struct {
	ID              int                    `json:"id"`
	DatesLocations  map[string][]string    `json:"datesLocations"`
}

// GetAllBands récupère la liste complète de tous les groupes de musique
func (a *API) GetAllBands() ([]Band, error) {
	if len(a.Bands) > 0 {
		return a.Bands, nil
	}

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

	a.Bands = bands

	return bands, nil
}

func (a *API) GetBandFromSearch(search string) ([]Band, error) {
	var bands, err = a.GetAllBands()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de la liste des groupes: %v", err)
	}
	var bandsFound []Band
	for _, band := range bands {
		if search != "" && strings.Contains(strings.ToLower(band.Name), strings.ToLower(search)) {
			bandsFound = append(bandsFound, band)
		}
	}
	return bandsFound, nil
}

func (a *API) FilterBands(filter Filter) ([]Band, error) {
	var bands, err = a.GetAllBands()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de la liste des groupes: %v", err)
	}
	var bandsFound []Band
	for _, band := range bands {
		if filter.Members != "" {
			members := strings.Split(filter.Members, ",")
			for _, member := range members {
				if !strings.Contains(strings.ToLower(strings.Join(band.Members, " ")), strings.ToLower(member)) {
					continue
				}
			}
		}
		if filter.NumberOfMembers != 0 && len(band.Members) != filter.NumberOfMembers {
			continue
		}

		if filter.Location != "" && band.Locations != filter.Location {
			continue
		}
		if filter.CreationDate != 0 {
			if band.CreationDate != filter.CreationDate {
				continue
			}
		}
		if filter.FirstAlbum != "" && band.FirstAlbum != filter.FirstAlbum {
			continue
		}
		if filter.ConcertDate != "" {
			if !strings.Contains(filter.ConcertDate, band.ConcertDates) {
				continue
			}
		}
		bandsFound = append(bandsFound, band)
	}
	return bandsFound, nil
}

func (a *API) GetBand(bandID int) (*Band, error) {
	var band *Band
	for _, b := range a.Bands {
		if b.ID == bandID {
			band = &b
			break
		}
	}

	if band != nil {
		return band, nil
	}

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
	if err := json.NewDecoder(resp.Body).Decode(&band); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de la réponse JSON: %v", err)
	}

	a.Bands = append(a.Bands, *band)

	return band, nil
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
