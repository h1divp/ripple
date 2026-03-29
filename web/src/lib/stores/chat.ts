import { writable, get } from 'svelte/store';
import type { LocationUpdate, DisplayMessage, ChatMessageSend, ChatMessageRecieved } from '$lib/types/types';
import { PUBLIC_API_WS_URL, PUBLIC_API_URL } from '$env/static/public';
import { joinMessage } from '$lib/utils/joinMessage';

export const messages = writable<DisplayMessage[]>([]);
export const isConnected = writable(false);
export const userCoords = writable({ lat: 0, lon: 0 });
export const nearbyCount = writable(0);

export const userDisplayName = writable("Anonymous Bat");
export const userAvatarUrl = writable("https://api.dicebear.com/9.x/icons/svg?backgroundColor=ffab91&seed=Adrian");
export const userId = writable("");
export const userTheme = writable("default");
export const userPreferences = writable({});
const storeMap = {
  displayName: userDisplayName,
  avatarUrl: userAvatarUrl,
  userId: userId,
  theme: userTheme,
  preferences: userPreferences
};

let socket: WebSocket;

export async function getProfile() {
  const url = `${PUBLIC_API_URL}/profile`;
  try {
    const response = await fetch(url, {
      method: 'GET',
      credentials: 'include'
    });
    if (!response.ok) {
      console.error('Profile retrieval failed with status:', response.status);
      return;
    }
    const data = await response.json();
    Object.entries(storeMap).forEach(([key, store]) => {
      if (data[key] !== undefined) {
        store.set(data[key]);
      }
    });
  } catch (err) {
    console.error('Profile retrieval failed:', err);
  }
}

export async function getSession() {
  const url = `${PUBLIC_API_URL}/register`;

  // Even if we have a cookie, we cannot check for it here
  // because of the httpOnly flag. however the browser will
  // send cookies in requests using credentials: 'include',
  // so we can check for it in the api.
  // TODO: return as promise and handle appropriately
  try {
    const response = await fetch(url, {
      method: 'POST',
      credentials: 'include'
    });
    if (!response.ok) {
      console.error('Session registration failed with status:', response.status);
    }
  } catch (err) {
    console.error('Session registration failed:', err);
  }
}

export function connect() {
  // cannot send displayName over as json body since this is not http, so we use query params
  const connUrl = `${PUBLIC_API_WS_URL}/chat/ws`;
  socket = new WebSocket(connUrl);

  socket.onopen = () => {
    isConnected.set(true);

    setTimeout(() => {
      messages.update((prev) => [...prev, joinMessage]);
    }, 500);

    const coords = get(userCoords);
    // TODO: retrieve displayName
    if (coords.lat !== 0) {
      sendLocationPing(coords.lat, coords.lon);
    }
  };
  socket.onclose = () => isConnected.set(false);

  // TODO: decompose
  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (data.type === 'chat') {
      const receivedMessage: ChatMessageRecieved = {
        type: 'chat',
        id: data.id,
        text: data.text,
        displayName: data.displayName || data.display_name, // Handle both camelCase and snake_case
        avatarUrl: data.avatarUrl || data.avatar_url,       // Handle both camelCase and snake_case
        timestamp: data.timestamp,
        status: 'sent'
      };

      messages.update((prev) => {
        // Make message change color when user recieves their own message back from the server
        const existingIndex = prev.findIndex((m) => m.id === receivedMessage.id);
        if (existingIndex !== -1) {
          const updated = [...prev];
          updated[existingIndex] = receivedMessage;
          return updated;
        }

        // Add new message from another user
        return [...prev, receivedMessage];
      });
    } else if (data.type === 'nearby_update') {
      nearbyCount.update(count => Math.max(0, count + data.delta));
      return;
    } else {
      return;
    }
  };
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

// TODO: turn socket check into gaurd clause
export function sendMessage(text: string) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    const newMessage: ChatMessageSend = {
      type: 'chat',
      id: crypto.randomUUID(),
      text: text,
      timestamp: Date.now(),
      status: 'sending',
    };

    const displayMessage: ChatMessageRecieved = {
      ...sendMessage,
      displayName: get(userDisplayName),
      avatarUrl: get(userAvatarUrl),
    } as ChatMessageRecieved;

    messages.update((prev) => [...prev, displayMessage]);
    socket.send(JSON.stringify(newMessage));
  }
}
