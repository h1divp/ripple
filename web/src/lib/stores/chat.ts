import { writable, get } from 'svelte/store';
import type { Message } from '$lib/types/types';
import { PUBLIC_WS_URL } from '$env/static/public';
import { generateRandomName, generateSessionSeed } from '$lib/utils/identity';
import { joinMessage } from '$lib/utils/joinMessage';

export const messages = writable<Message[]>([]);
export const isConnected = writable(false);
export const userDisplayName = writable(generateRandomName());
export const userAvatarSeed = writable(generateSessionSeed());
export const userId = writable(crypto.randomUUID());
export const userCoords = writable({ lat: 0, lon: 0 });

let socket: WebSocket;

export function connect(userId: string) {
  // cannot send displayName over as json body since this is not http, so we use query params
  const connUrl = `${PUBLIC_WS_URL}?userId=${userId}`;
  socket = new WebSocket(connUrl);

  socket.onopen = () => {
    isConnected.set(true);
    const coords = get(userCoords);
    if (coords.lat !== 0) {
      sendLocationPing(coords.lat, coords.lon);
    }
  };
  socket.onclose = () => isConnected.set(false);

  messages.update((prev) => [...prev, joinMessage]);

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    // TODO: add nearby_update
    // if (data.type === "nearby_update") {
    //   nearbyCount.set(data.nearbyCount);
    //   return
    // }
    if (data.type === "location_update") return;

    messages.update((prev) => {
      const existingIndex = prev.findIndex((m) => m.id === data.id);
      if (existingIndex !== -1) {
        // Update existing optimistic message
        const updated = [...prev];
        updated[existingIndex] = { ...data, status: 'sent' };
        return updated;
      }
      // Add new message from another user
      return [...prev, { ...data, status: 'sent' }];
    });
 };
}

export function sendLocationPing(lat: number, lon: number) {
  if (socket?.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({
      senderID: get(userID),
      type: 'location_update',
      lat,
      lon,
    }));
  }
}

export function sendMessage(text: string) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    const coords = get(userCoords);

    const newMessage: Message = {
      id: crypto.randomUUID(),
      type: 'chat',
      text: text,
      displayName: get(userDisplayName),
      avatarSeed: get(userAvatarSeed),
      senderId: get(userId),
      timestamp: Date.now(),
      status: 'sending',
      lat: coords.lat,
      lon: coords.lon
    };

    messages.update((prev) => [...prev, newMessage]);
    socket.send(JSON.stringify(newMessage));
  }
}
