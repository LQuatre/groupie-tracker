package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type API struct {
	BaseURL       string
	Bands         []Band
	Relationships []Relationship
}

func NewAPI(baseURL string) *API {
	return &API{BaseURL: baseURL}
}

func (a *API) ShowAPI() {
	fmt.Printf("API BaseURL: %s\n", a.BaseURL)
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	Dates []string `json:"dates"`
}

type Band struct {
	ID                   int        `json:"id"`
	Image                string     `json:"image"`
	Name                 string     `json:"name"`
	Members              []string   `json:"members,omitempty"`
	CreationDate         int        `json:"creationDate"`
	FirstAlbum           string     `json:"firstAlbum"`
	Locations            string     `json:"locations"`
	ConcertDates         string     `json:"concertDates"`
	Relations            string     `json:"relations"`
	LocationsCoordinates []Location `json:"locationsCoordinates"`
}

type Filter struct {
	Members         string `json:"members,omitempty"`
	NumberOfMembers int    `json:"numberOfMembers"`
	Location        string `json:"location"`
	CreationDate    int    `json:"creationDate"`
	FirstAlbum      string `json:"firstAlbum"`
	ConcertDate     string `json:"concertDate"`
}

type Relationship struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func (a *API) GetAllBands() ([]Band, error) {
	if len(a.Bands) > 1 {
		return a.Bands, nil
	}

	resp, err := http.Get(fmt.Sprintf("%s/artists", a.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

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

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'envoi de la requête: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a renvoyé un statut non-OK: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&band); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de la réponse JSON: %v", err)
	}

	a.Bands = append(a.Bands, *band)

	return band, nil
}

func (a *API) GetAllRelations() ([]Relationship, error) {
	if len(a.Relationships) > 1 {
		return a.Relationships, nil
	}

	url := fmt.Sprintf("%s/relation", a.BaseURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'envoi de la requête: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a renvoyé un statut non-OK: %d", resp.StatusCode)
	}

	var relationships []Relationship
	if err := json.NewDecoder(resp.Body).Decode(&relationships); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de la réponse JSON: %v", err)
	}

	a.Relationships = relationships

	return relationships, nil
}

func (a *API) GetRelation(relationshipID int) (*Relationship, error) {
	var relationship *Relationship
	for _, r := range a.Relationships {
		if r.ID == relationshipID {
			relationship = &r
			break
		}
	}

	if relationship != nil {
		return relationship, nil
	}

	url := fmt.Sprintf("%s/relation/%d", a.BaseURL, relationshipID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'envoi de la requête: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a renvoyé un statut non-OK: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&relationship); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de la réponse JSON: %v", err)
	}

	a.Relationships = append(a.Relationships, *relationship)

	return relationship, nil
}
