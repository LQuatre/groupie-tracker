package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

type Filter struct {
	Name string `json:"name"`
	NumberOfMembers int `json:"numberOfMembers"`
	Location string `json:"location"`
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	Locations []string `json:"locations"`
	ConcertDates []string `json:"concertDates"`
	Relations []string `json:"relations"`
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

func (a *API) GetBandFromSearch(search string) ([]Band, error) {
	var bands, err = a.GetAllBands()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de la liste des groupes: %v", err)
	}
	var bandsFound []Band
	for _, band := range bands {
		if band.Name == search {
			bandsFound = append(bandsFound, band)
		}
	}
	return bandsFound, nil
}

func (a *API) GetBandFromFilter(filter Filter) ([]Band, error) {
	var bands, err = a.GetAllBands()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de la liste des groupes: %v", err)
	}
	var bandsFound []Band
	for _, band := range bands {
		if filter.Name != "" && band.Name != filter.Name {
			continue
		}
		if filter.NumberOfMembers != 0 && len(band.Members) != filter.NumberOfMembers {
			continue
		}

		if filter.Location != "" && band.Locations != filter.Location {
			continue
		}
		if filter.StartDate != "" {
			startDate, err := strconv.Atoi(filter.StartDate)
			if err != nil {
				return nil, fmt.Errorf("failed to convert start date to integer: %w", err)
			}
			if band.CreationDate < startDate {
				continue
			}
		}
		if filter.EndDate != "" {
			endDate, err := strconv.Atoi(filter.EndDate)
			if err != nil {
				return nil, fmt.Errorf("failed to convert end date to integer: %w", err)
			}
			if band.CreationDate > endDate {
				continue
			}
		}
		if len(filter.Locations) != 0 {
			found := false
			for _, location := range filter.Locations {
				if location == band.Locations {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(filter.ConcertDates) != 0 {
			found := false
			for _, concertDate := range filter.ConcertDates {
				if concertDate == band.ConcertDates {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(filter.Relations) != 0 {
			found := false
			for _, relation := range filter.Relations {
				if relation == band.Relations {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		bandsFound = append(bandsFound, band)
	}
	return bandsFound, nil
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
