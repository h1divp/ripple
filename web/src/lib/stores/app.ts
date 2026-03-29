import { writable } from 'svelte/store';

export const locationError = writable(true);
export const isInitialized = writable(false);
export const initError = writable<string | null>(null);
