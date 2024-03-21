// Fonction pour envoyer une requête GET avec les filtres et la recherche
function sendSearchRequest() {
  var query = document.getElementById("searchInput").value;
  var members = document.getElementById("members").value;
  var numberOfMembers = document.getElementById("numberofmember").value;
  var createDate = document.getElementById("creation-date").value;
  var firstAlbum = document.getElementById("first-album").value;
  var concertDate = document.getElementById("concert-date").value;

  var xhr = new XMLHttpRequest();
  var url =
    "/search?query=" +
    query +
    "&members=" +
    members +
    "&numberofmember=" +
    numberOfMembers +
    "&creation-date=" +
    createDate +
    "&first-album=" +
    firstAlbum +
    "&concert-date=" +
    concertDate;

  xhr.open("GET", url, true);
  xhr.onreadystatechange = function () {
    if (xhr.readyState == XMLHttpRequest.DONE) {
      if (xhr.status == 200) {
        document.getElementById("artists").innerHTML = xhr.responseText;
      } else {
        console.error("Error:", xhr.status);
      }
    }
  };
  xhr.send();
}

// Ajouter des écouteurs d'événements pour l'input de recherche et chaque élément de filtre
document
  .getElementById("searchInput")
  .addEventListener("input", sendSearchRequest);
document.getElementById("members").addEventListener("input", sendSearchRequest);
document
  .getElementById("numberofmember")
  .addEventListener("change", sendSearchRequest);
document
  .getElementById("creation-date")
  .addEventListener("input", sendSearchRequest);
document
  .getElementById("first-album")
  .addEventListener("input", sendSearchRequest);
document
  .getElementById("concert-date")
  .addEventListener("input", sendSearchRequest);
