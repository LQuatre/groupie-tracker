document.getElementById("searchInput").addEventListener("input", function () {
  var query = this.value;
  var xhr = new XMLHttpRequest();
  xhr.open("GET", "/search?query=" + query, true);
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
});
