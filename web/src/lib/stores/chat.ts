import { writable } from 'svelte/store';
import type { DisplayMessage } from '$lib/types/types';

export const messages = writable<DisplayMessage[]>([]);
export const isConnected = writable(false);
export const nearbyCount = writable(0);
