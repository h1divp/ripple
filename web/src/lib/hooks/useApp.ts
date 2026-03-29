import { onMount } from 'svelte';
import { initializeApp, initializeGeolocation } from '$lib/services/initialization';

export function useApp() {
  const geolocation = initializeGeolocation();
  let isInitialized = $state(false);
  let initError = $state<string | null>(null);

  onMount(async () => {
    try {
      await initializeApp();
      geolocation.start();
      isInitialized = true;
    } catch (error) {
      console.error('Failed to initialize app:', error);
      initError = error instanceof Error ? error.message : 'Unknown error';
    }

    return () => {
      geolocation.stop();
    };
  });

  return {
    isInitialized: () => isInitialized,
    initError: () => initError,
    locationError: geolocation.locationError
  };
}
