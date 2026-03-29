import { writable, derived } from 'svelte/store';

export const userCoords = writable({ lat: 0, lon: 0 });

// Derived store to check if location is available
export const hasLocation = derived(
  userCoords,
  ($coords) => $coords.lat !== 0 || $coords.lon !== 0
);
