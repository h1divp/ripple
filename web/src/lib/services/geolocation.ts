import { get } from 'svelte/store';
import { userCoords, isConnected } from '$lib/stores';
import { sendLocationPing } from '$lib/services/websocket';

export interface GeolocationManager {
  start: () => void;
  stop: () => void;
}

export function createGeolocationManager(
  onLocationUpdate: (lat: number, lon: number) => void,
  onError: (error: GeolocationPositionError) => void
): GeolocationManager {
  let watchId: number | null = null;
  let pingInterval: number | null = null;

  const start = () => {
    watchId = navigator.geolocation.watchPosition(
      (pos) => {
        const lat = pos.coords.latitude;
        const lon = pos.coords.longitude;
        userCoords.set({ lat, lon });
        onLocationUpdate(lat, lon);
      },
      onError,
      {
        enableHighAccuracy: false,
        timeout: 15000,
        maximumAge: 60000,
      }
    );

    pingInterval = setInterval(() => {
      const coords = get(userCoords);
      if (coords.lat !== 0 && get(isConnected)) {
        sendLocationPing(coords.lat, coords.lon);
      }
    }, 60000);
  };

  const stop = () => {
    if (watchId !== null) {
      navigator.geolocation.clearWatch(watchId);
      watchId = null;
    }
    if (pingInterval !== null) {
      clearInterval(pingInterval);
      pingInterval = null;
    }
  };

  return { start, stop };
}
