const fs = require('fs');
const path = require('path');
const { JSDOM } = require('jsdom');

let dom;
let container;

// Setup before each test
beforeEach(() => {
  const html = fs.readFileSync(path.resolve(__dirname, 'artist.html'), 'utf8');
  dom = new JSDOM(html);
  container = dom.window.document;
});

test('it has a title of "Artist"', () => {
  expect(container.querySelector('title').textContent).toBe('Artist');
});

test('it has a heading of "Artist Profile"', () => {
  expect(container.querySelector('h1').textContent).toBe('Artist Profile');
});

test('it has a div with class "artist-card"', () => {
  expect(container.querySelector('.artist-card')).not.toBeNull();
});

test('it has a script tag for Google Maps API', () => {
  const scriptTag = container.querySelector('script[src^="https://maps.googleapis.com/maps/api/js"]');
  expect(scriptTag).not.toBeNull();
  expect(scriptTag.getAttribute('async')).not.toBeNull();
});

test('it has a script tag for artist.js', () => {
  expect(container.querySelector('script[src="/static/assets/artist.js"]')).not.toBeNull();
});