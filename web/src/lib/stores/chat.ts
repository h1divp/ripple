import { writable } from 'svelte/store';
import type { Message } from '$lib/types/types';
import { PUBLIC_WS_URL } from '$env/static/public';

export const messages = writable<Message[]>([]);
export const isConnected = writable(false);

let socket: WebSocket;

export function connect(displayName: string) {
  // cannot send displayName over as json body since this is not http, so we use query params
  const connUrl = `${PUBLIC_WS_URL}?displayName=${displayName}`;
  socket = new WebSocket(connUrl);

  socket.onopen = () => isConnected.set(true);
  socket.onclose = () => isConnected.set(false);

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    messages.update((prev) => [...prev, data]);
  };
}

export function sendMessage(text: string) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({ text }));
  }
}
