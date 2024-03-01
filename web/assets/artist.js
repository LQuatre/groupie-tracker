// Initialize and add the map
let map;

async function initMap(locations) {
  if (!Array.isArray(locations)) {
    locations = [locations];
  }

  console.log(locations);
  const latLngPairs = locations[0].match(/-?\d+\.\d+/g);
  const positions = [];

  for (let i = 0; i < latLngPairs.length; i += 2) {
    const lat = parseFloat(latLngPairs[i]);
    const lng = parseFloat(latLngPairs[i + 1]);
    positions.push({ lat, lng });
  }

  console.log(positions);

  for (let i = 0; i < positions.length; i++) {
    if (isNaN(positions[i].lat) || isNaN(positions[i].lng)) {
      alert("Invalid location data");
      return;
    }
  }

  while (window.google === undefined) {
    await new Promise((resolve) => setTimeout(resolve, 100));
  }

  // Request needed libraries.
  //@ts-ignore
  const { Map } = await google.maps.importLibrary("maps");
  const { AdvancedMarkerView } = await google.maps.importLibrary("marker");

  // The map, centered at the first position
  map = new Map(document.getElementById("map"), {
    zoom: 4,
    center: positions[0],
    mapId: "DEMO_MAP_ID",
  });

  // Create waypoints for all positions except the first one
  for (let i = 0; i < positions.length; i++) {
    new AdvancedMarkerView({
      position: positions[i],
      map: map,
      title: "Waypoint " + (i + 1),
    });
  }
}

// recupere les data de la div map et les affiche
function showMap() {
  var mapDiv = document.getElementById("map");
  var locations = mapDiv.getAttribute("data-locations");

  // Initialiser la carte Google Maps
  initMap(locations);
}

showMap();
