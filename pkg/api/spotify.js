const clientID = '029d534a35a84eafb27b44aee587d3a6D';
const clientSecret = '52ddbf2f6adf44aeba247cbc75473f8d';

// Encodez le client ID et le client secret en base64
const basicAuth = btoa(`${clientID}:${clientSecret}`);

// Paramètres de la demande
const requestOptions = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded',
    'Authorization': `Basic ${basicAuth}`
  },
  body: 'grant_type=client_credentials'
};

// Envoyer une demande pour obtenir un jeton d'accès
fetch('https://accounts.spotify.com/api/token', requestOptions)
  .then(response => response.json())
  .then(data => {
    const accessToken = data.access_token;
    console.log('Jetoon d\'accès Spotify :', accessToken);
    // Utilisez maintenant accessToken pour authentifier vos requêtes à l'API Spotify
  })
  .catch(error => {
    console.error('Erreur lors de la récupération du jeton d\'accès :', error);
  });
