const client_id = '029d534a35a84eafb27b44aee587d3a6'; 
const client_secret = '52ddbf2f6adf44aeba247cbc75473f8d';

async function getToken() {
  const response = await fetch('https://accounts.spotify.com/api/token', {
    method: 'POST',
    body: new URLSearchParams({
      'grant_type': 'client_credentials',
    }),
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      'Authorization': 'Basic ' + (Buffer.from(client_id + ':' + client_secret).toString('base64')),
    },
  });

  return await response.json();
}

async function getArtistId(access_token, artistName) {
  const response = await fetch(`https://api.spotify.com/v1/search?q=${encodeURIComponent(artistName)}&type=artist&limit=1`, {
    method: 'GET',
    headers: { 'Authorization': 'Bearer ' + access_token },
  });

  const data = await response.json();
  return data.artists.items[0].id; // Return the ID of the first artist in the search results
}

async function updateIframeLink(artistName) {
  try {
    // Récupérer le token d'accès
    const tokenResponse = await getToken();
    const access_token = tokenResponse.access_token;

    // Récupérer l'ID Spotify de l'artiste
    const artistId = await getArtistId(access_token, artistName);

    // Mettre à jour le lien iframe avec l'ID Spotify de l'artiste
    const iframeLink = `https://open.spotify.com/embed/artist/${artistId}?utm_source=generator`;
    const iframe = document.getElementById('spotify-iframe');
    iframe.src = iframeLink;

    console.log(`Lien mis à jour pour l'artiste "${artistName}": ${iframeLink}`);
  } catch (error) {
    console.error('Erreur lors de la mise à jour du lien iframe:', error);
  }
}

artistNames.forEach(artistName => {
  updateIframeLink(artistName);
});
