package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type IndexLocations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type IndexDates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type API struct {
	BaseURL   string
	BaseApi   map[string]string
	Artists   []Band
	Locations []IndexLocations
	Dates     []IndexDates
	Relation  []Relation
}

func NewAPI(baseURL string) *API {
	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("L'API a renvoyé un statut non-OK:", resp.StatusCode)
		return nil
	}

	var resp2 map[string]string
	err = json.NewDecoder(resp.Body).Decode(&resp2)
	if err != nil {
		fmt.Println("Failed to decode response:", err)
		return nil
	}

	// for key, value := range resp2 {
	// 	fmt.Printf("%s: %s\n", key, value)
	// }

	artists, err := http.Get(resp2["artists"])
	if err != nil {
		fmt.Println("Error getting artists:", err)
		return nil
	}
	defer artists.Body.Close()

	var bands []Band
	err = json.NewDecoder(artists.Body).Decode(&bands)
	if err != nil {
		fmt.Println("Failed to decode artists:", err)
		return nil
	}

	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("data.json does not exist, ignoring...")
		} else {
			fmt.Println("Failed to read data.json:", err)
			return nil
		}
	}

	var bands2 []Band

	if len(data) > 0 {
		err = json.Unmarshal(data, &bands2)
		if err != nil {
			fmt.Println("Failed to decode data.json:", err)
			return nil
		}
	}

	bands = append(bands, bands2...)

	var locs []IndexLocations
	var dts []IndexDates
	var rels []Relation

	for _, band := range bands {
		fmt.Println("Getting info for band:", band.Name)

		locURL := fmt.Sprintf("%s/%d", resp2["locations"], band.ID)
		locResp, err := http.Get(locURL)
		if err != nil {
			fmt.Println("Error getting location:", err)
			continue
		}
		defer locResp.Body.Close()

		var loc IndexLocations
		err = json.NewDecoder(locResp.Body).Decode(&loc)
		if err != nil {
			fmt.Println("Failed to decode location:", err)
			continue
		}
		locs = append(locs, loc)

		datesURL := fmt.Sprintf("%s/%d", resp2["dates"], band.ID)
		datesResp, err := http.Get(datesURL)
		if err != nil {
			fmt.Println("Error getting dates:", err)
			continue
		}
		defer datesResp.Body.Close()

		var dt IndexDates
		err = json.NewDecoder(datesResp.Body).Decode(&dt)
		if err != nil {
			fmt.Println("Failed to decode dates:", err)
			continue
		}
		dts = append(dts, dt)

		relURL := fmt.Sprintf("%s/%d", resp2["relation"], band.ID)
		relResp, err := http.Get(relURL)
		if err != nil {
			fmt.Println("Error getting relation:", err)
			continue
		}
		defer relResp.Body.Close()

		var rel Relation
		err = json.NewDecoder(relResp.Body).Decode(&rel)
		if err != nil {
			fmt.Println("Failed to decode relation:", err)
			continue
		}
		rels = append(rels, rel)
	}

	fmt.Println("API initialisée avec succès")

	return &API{
		BaseURL:   baseURL,
		BaseApi:   resp2,
		Artists:   bands,
		Locations: locs,
		Dates:     dts,
		Relation:  rels,
	}
}

func (a *API) ShowAPI() {
	fmt.Printf("API BaseURL: %s\n", a.BaseURL)
}

type Location struct {
	Lat   float64  `json:"lat"`
	Lng   float64  `json:"lng"`
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
	RelationExists       bool       `json:"relationExists"`
}

type Filter struct {
	Members         string `json:"members,omitempty"`
	NumberOfMembers int    `json:"numberOfMembers"`
	Location        string `json:"location"`
	CreationDate    int    `json:"creationDate"`
	FirstAlbum      string `json:"firstAlbum"`
	ConcertDate     string `json:"concertDate"`
}

func (a *API) GetAllBands() ([]Band, error) {
	return a.Artists, nil
}

func (a *API) GetBandFromSearch(search string) ([]Band, error) {
	var bands, err = a.GetAllBands()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de la liste des groupes: %v", err)
	}
	var bandsFound []Band
	for _, band := range bands {
		if strings.Contains(strings.ToLower(band.Name), strings.ToLower(search)) {
			bandsFound = append(bandsFound, band)
		}
	}
	return bandsFound, nil
}

func (a *API) GetBandFromSearchWithBands(search string, bands []Band) ([]Band, error) {
	var bandsFound []Band
	for _, band := range bands {
		if strings.Contains(strings.ToLower(band.Name), strings.ToLower(search)) {
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
			allMembersPresent := true
			for _, member := range members {
				if !strings.Contains(strings.ToLower(strings.Join(band.Members, " ")), strings.ToLower(member)) {
					allMembersPresent = false
					break
				}
			}
			if !allMembersPresent {
				continue
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
			if !strings.Contains(band.ConcertDates, filter.ConcertDate) {
				continue
			}
		}
		// fmt.Println(band.Name)
		bandsFound = append(bandsFound, band)
	}
	return bandsFound, nil
}

func (a *API) GetBand(bandID int) (*Band, error) {
	var band *Band
	for _, b := range a.Artists {
		if b.ID == bandID {
			band = &b
			break
		}
	}
	if band == nil {
		return nil, fmt.Errorf("groupe non trouvé")
	}
	return band, nil
}

func (a *API) GetAllRelations() ([]Relation, error) {
	return a.Relation, nil
}

func (a *API) GetRelation(relationshipID int) (*Relation, error) {
	var relationship *Relation
	for _, r := range a.Relation {
		if r.ID == relationshipID {
			relationship = &r
			break
		}
	}
	if relationship == nil {
		return nil, fmt.Errorf("relation non trouvée")
	}
	return relationship, nil
}

func (a *API) AddBand(band Band) error {
	band.ID = len(a.Artists) + 1

	imageURL, err := GetUrlOfTheFirstSquareImage(band.Name)

	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Affecter l'URL à la propriété Image du groupe
	band.Image = imageURL

	// Ajout du groupe à la liste d'artistes de l'API
	a.Artists = append(a.Artists, band)

	var bands []Band
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("data.json does not exist, ignoring...")
		} else {
			fmt.Println("Failed to read data.json:", err)
			return err
		}
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, &bands)
		if err != nil {
			fmt.Println("Failed to decode data.json:", err)
			return err
		}
	}

	bands = append(bands, band)

	// ajouter dans un fichier json
	jsonFile, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonData, err := json.Marshal(bands)
	if err != nil {
		fmt.Println(err)
	}
	jsonFile.Write(jsonData)

	return nil
}

type BingImageResponse struct {
	Value []struct {
		ContentURL string `json:"contentUrl"`
	} `json:"value"`
}

func GetUrlOfTheFirstSquareImage(search string) (string, error) {
	// Requête à l'API Bing Image Search
	bingURL := fmt.Sprintf("https://api.bing.microsoft.com/v7.0/images/search?q=%s&count=1&aspect=square", search)
	client := &http.Client{}
	req, err := http.NewRequest("GET", bingURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", "e7cf8887a3db45448c45698acc40936c")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Traitement de la réponse JSON
	var bingResponse BingImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&bingResponse); err != nil {
		return "", err
	}

	// Vérification s'il y a une réponse valide
	if len(bingResponse.Value) == 0 {
		return "", fmt.Errorf("aucune image trouvée")
	}

	// Récupération de l'URL de l'image
	imageURL := bingResponse.Value[0].ContentURL
	if imageURL == "" {
		return "", fmt.Errorf("aucune URL d'image trouvée")
	}

	return imageURL, nil
}
