let map;

async function initMap(locations) {
  if (!Array.isArray(locations)) {
    locations = [locations];
  }

  console.log(locations);
  const latLngDateTriples = locations[0].match(
    /\{-?\d+\.\d+ -?\d+\.\d+\s\[\d{2}-\d{2}-\d{4}(?:\s\d{2}-\d{2}-\d{4})*\]\}/g // {lat lng [date1 date2 ...]} || oui je fais des trucs bizzare des des fois
  );
  const positions = [];

  for (const latLngDate of latLngDateTriples) {
    const [lat, lng, ...dates] = latLngDate.match(
      /-?\d+\.\d+|-?\d{2}-\d{2}-\d{4}/g // lat, lng, and dates || j'ai du le faire une deuxième fois merdeeeeeeeeeeeeeeeee
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

  const { Map } = await google.maps.importLibrary("maps");
  const { AdvancedMarkerView } = await google.maps.importLibrary("marker");

  map = new Map(document.getElementById("map"), {
    zoom: 4,
    center: positions[0],
    mapId: "DEMO_MAP_ID",
  });

  for (let i = 0; i < positions.length; i++) {
    dateString = "Date : ";
    for (let j = 0; j < positions[i].dates.length; j++) {
      dateString += positions[i].dates[j];
      if (j < positions[i].dates.length - 1) {
        dateString += "\nDate : ";
      }
    }
    new AdvancedMarkerView({
      position: positions[i],
      map: map,
      title: dateString,
    });
  }
}

function showMap() {
  var mapDiv = document.getElementById("map");
  var locations = mapDiv.getAttribute("data-locations");

  initMap(locations);
}

showMap();
