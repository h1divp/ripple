import { getSession, getProfile } from '$lib/services/api';
import { connect } from '$lib/services/websocket';
import { createGeolocationManager } from './geolocation';
import { sendLocationPing } from '$lib/services/websocket';
import { locationError, isInitialized, initError } from '$lib/stores/app';

export async function initializeApp() {
  try {
    await getSession();
    await getProfile();
    connect();
    isInitialized.set(true);
    initError.set(null);
  } catch (error) {
    console.error('Failed to initialize app:', error);
    initError.set(error instanceof Error ? error.message : 'Unknown error');
    throw error;
  }
}

export function initializeGeolocation() {
  const handleLocationUpdate = (lat: number, lon: number) => {
    locationError.set(false);
    sendLocationPing(lat, lon);
  };

  const handleLocationError = (err: GeolocationPositionError) => {
    console.error('Geolocation error:', err);
    if (err.code === 1) {
      locationError.set(true);
    } else if (err.code === 2) {
      console.log('Geolocation service unavailable');
    }
  };

  const geolocationManager = createGeolocationManager(
    handleLocationUpdate,
    handleLocationError
  );

  return geolocationManager;
}
