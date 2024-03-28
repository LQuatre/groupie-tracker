package routes

import (
	"encoding/json"
	"net/http"
	"text/template"
)

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

func onlySendData(w http.ResponseWriter, data interface{}) {
	// Convertir les données en JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		handleError(w, err)
		return
	}

	// Définir le type de contenu de la réponse comme JSON
	w.Header().Set("Content-Type", "application/json")
	
	w.Write(jsonData)
}
