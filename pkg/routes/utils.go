package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GeocodeAddress(placeName string) (float64, float64) {
	apiKey := "AIzaSyDKgbLCKstOrqzxFuKjD0-GH4aXAN8CjXM"

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
