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

  const now = new Date();
  const joinMsg: Message = {
    id: crypto.randomUUID(),
    text: `joined at ${now.toLocaleString([], { hour: '2-digit', minute: '2-digit' })}`,
    displayName: get(userDisplayName),
    avatarSeed: get(userAvatarSeed),
    senderId: get(userId),
    timestamp: Date.now(),
    status: 'sent',
    type: 'system',
  };
  messages.update((prev) => [...prev, joinMsg]);

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);

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

export function sendMessage(text: string) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    const newMessage: Message = {
      id: crypto.randomUUID(),
      text: text,
      displayName: get(userDisplayName),
      avatarSeed: get(userAvatarSeed),
      senderId: get(userId),
      timestamp: Date.now(),
      status: 'sending',
    };

    messages.update((prev) => [...prev, newMessage]);
    socket.send(JSON.stringify(newMessage));
  }
}
