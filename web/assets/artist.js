// Initialize and add the map
let map;

async function initMap(locations) {
  if (!Array.isArray(locations)) {
    locations = [locations];
  }

  console.log(locations);
  const latLngDateTriples = locations[0].match(
    /\{-?\d+\.\d+ -?\d+\.\d+\s\[\d{2}-\d{2}-\d{4}(?:\s\d{2}-\d{2}-\d{4})*\]\}/g
  );
  const positions = [];

  for (const latLngDate of latLngDateTriples) {
    const [lat, lng, ...dates] = latLngDate.match(
      /-?\d+\.\d+|-?\d{2}-\d{2}-\d{4}/g
    );
    positions.push({ lat: parseFloat(lat), lng: parseFloat(lng), dates });
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
    dateString = "";
    for (let j = 0; j < positions[i].dates.length; j++) {
      dateString += positions[i].dates[j];
      if (j < positions[i].dates.length - 1) {
        dateString += "\n";
      }
    }
    new AdvancedMarkerView({
      position: positions[i],
      map: map,
      title: dateString,
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
