import { userCoords } from '$lib/stores';

export interface GeolocationManager {
  start: () => void;
  stop: () => void;
}

export function createGeolocationManager(
  onLocationUpdate: (lat: number, lon: number) => void,
  onError: (error: GeolocationPositionError) => void
): GeolocationManager {
  let intervalId: ReturnType<typeof setInterval> | null = null;
  const LOCATION_CHECK_INTERVAL = 30000; // 30 seconds

  const getCurrentLocation = () => {
    navigator.geolocation.getCurrentPosition(
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
        maximumAge: 30000,
      }
    );
  };

  const start = () => {
    getCurrentLocation();
    intervalId = setInterval(getCurrentLocation, LOCATION_CHECK_INTERVAL);
  };

  const stop = () => {
    console.log('Stopping geolocation checks');
    if (intervalId !== null) {
      clearInterval(intervalId);
      intervalId = null;
    }
  };

  return { start, stop };
}
