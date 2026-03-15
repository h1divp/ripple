import { writable, get } from 'svelte/store';
import type { Message } from '$lib/types/types';
import { PUBLIC_WS_URL } from '$env/static/public';
import { generateRandomName, generateSessionSeed } from '$lib/utils/identity';

export const messages = writable<Message[]>([]);
export const isConnected = writable(false);
export const userDisplayName = writable(generateRandomName());
export const userAvatarSeed = writable(generateSessionSeed());
export const userId = writable(crypto.randomUUID());

let socket: WebSocket;

export function connect(displayName: string) {
  // cannot send displayName over as json body since this is not http, so we use query params
  const connUrl = `${PUBLIC_WS_URL}?displayName=${displayName}`;
  socket = new WebSocket(connUrl);

  socket.onopen = () => isConnected.set(true);
  socket.onclose = () => isConnected.set(false);

  socket.onmessage = (event) => {
    const incomingData = JSON.parse(event.data);

    messages.update((prev) => {
      const existingIndex = prev.findIndex(m => m.id === incomingData.id);

      if (existingIndex !== -1) {
        const updated = [...prev];
        updated[existingIndex] = { ...incomingData, status: 'sent' };
        return updated;
      }

      return [...prev, { incomingData, status: 'sent' }];
    })
  }
}

export function sendMessage(text: string) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    const newMessage: Message = {
      id: crypto.randomUUID(),
      text: text,
      displayName: get(userDisplayName),
      avatarSeed: get(userAvatarSeed),
      senderId: get(userId),
      timestamp: Date.now(),
      status: 'sending'
    };

    messages.update((prev) => [...prev, newMessage]);
    socket.send(JSON.stringify(newMessage))
  }
}

