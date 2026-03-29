import { get } from 'svelte/store';
import { PUBLIC_API_WS_URL } from '$env/static/public';
import { messages, isConnected } from '$lib/stores/chat';
import { userDisplayName, userAvatarUrl } from '$lib/stores/user';
import { userCoords } from '$lib/stores/location';
import { joinMessage } from '$lib/utils/joinMessage';
import type { LocationUpdate, ChatMessageSend, ChatMessageRecieved } from '$lib/types/types';
import { handleChatMessage, handleNearbyUpdate } from './messageHandlers';

let socket: WebSocket;

export function connect() {
  const connUrl = `${PUBLIC_API_WS_URL}/chat/ws`;
  socket = new WebSocket(connUrl);

  socket.onopen = handleSocketOpen;
  socket.onclose = handleSocketClose;
  socket.onmessage = handleSocketMessage;
}

function handleSocketOpen() {
  isConnected.set(true);
  setTimeout(() => {
    messages.update((prev) => [...prev, joinMessage]);
  }, 500);

  const coords = get(userCoords);
  if (coords.lat !== 0) {
    sendLocationPing(coords.lat, coords.lon);
  }
}

function handleSocketClose() {
  isConnected.set(false);
}

function handleSocketMessage(event: MessageEvent) {
  const data = JSON.parse(event.data);

  switch (data.type) {
    case 'chat':
      handleChatMessage(data);
      break;
    case 'nearby_update':
      handleNearbyUpdate(data);
      break;
    default:
      console.log('Unknown message type:', data.type);
  }
}

export function sendLocationPing(lat: number, lon: number) {
  if (socket?.readyState === WebSocket.OPEN) {
    const locationMsg: LocationUpdate = {
      type: 'location_update',
      lat: lat,
      lon: lon
    };
    socket.send(JSON.stringify(locationMsg));
  }
}

export function sendMessage(text: string) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    const currentDisplayName = get(userDisplayName);
    const currentAvatarUrl = get(userAvatarUrl);

    if (!currentDisplayName || currentDisplayName === "Anonymous Bat") {
      console.warn("Profile not loaded yet, cannot send message");
      return;
    }

    const newMessage: ChatMessageSend = {
      type: 'chat',
      id: crypto.randomUUID(),
      text: text,
      timestamp: Date.now(),
      status: 'sending',
    };

    const displayMessage: ChatMessageRecieved = {
      type: 'chat',
      id: newMessage.id,
      text: newMessage.text,
      timestamp: newMessage.timestamp,
      status: 'sending',
      displayName: currentDisplayName,
      avatarUrl: currentAvatarUrl,
    };

    messages.update((prev) => [...prev, displayMessage]);
    socket.send(JSON.stringify(newMessage));
  }
}
