// JavaScript pour initialiser la carte
mapboxgl.accessToken = "YOUR_MAPBOX_ACCESS_TOKEN"; // Remplacez par votre jeton d'accès Mapbox
const map = new mapboxgl.Map({
  container: "map",
  style: "mapbox://styles/mapbox/streets-v11", // Style de la carte (vous pouvez choisir un autre style si vous le souhaitez)
  center: [0, 0], // Coordonnées du centre de la carte par défaut
  zoom: 1, // Niveau de zoom par défaut
});

// Fonction pour ajouter un marqueur à la carte pour chaque emplacement
function addMarkers(locations) {
  locations.forEach((location) => {
    // Utilisez la latitude et la longitude de chaque emplacement pour ajouter un marqueur à la carte
    new mapboxgl.Marker()
      .setLngLat([location.longitude, location.latitude])
      .addTo(map);
  });
}

// Récupérez les données de localisation de votre API Go
fetch("/api/locations") // Remplacez cet URL par l'URL de votre endpoint pour récupérer les données de localisation
  .then((response) => response.json())
  .then((data) => {
    const locations = processData(data); // Traitez les données pour obtenir les coordonnées de latitude et de longitude
    addMarkers(locations); // Ajoutez des marqueurs à la carte
  })
  .catch((error) => {
    console.error("Error fetching data:", error);
  });

// Fonction pour traiter les données et récupérer les coordonnées de latitude et de longitude
function processData(data) {
  // Traitez vos données ici pour obtenir les coordonnées de latitude et de longitude
  // Retournez un tableau d'objets avec les coordonnées de chaque emplacement
  return data.map((location) => ({
    latitude: location.latitude,
    longitude: location.longitude,
  }));
}
