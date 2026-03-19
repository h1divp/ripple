<script lang="ts">
  import { onMount } from 'svelte';
  import {
    messages,
    isConnected,
    connect,
    sendMessage,
    userDisplayName,
    sendLocationPing,
    userCoords,
  } from '$lib/stores/chat';
  import ChatHeader from '$lib/components/ChatHeader.svelte';
  import MessageList from '$lib/components/MessageList.svelte';
  import ChatInput from '$lib/components/ChatInput.svelte';

  let nearbyCounter = $state(0);
  let locationError = $state(true);

  const getCookie = (name: string) => {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()?.split(';').shift();
    return null;
  };

  onMount(async () => {
    // get session cookie if none found
    let sessionCookie = getCookie('session_id');
    if (!sessionCookie) {
      try {
        const response = await fetch('/register', { method: 'POST' });
        if (!response.ok) {
          console.error('Session registration failed with status:', response.status);
        }
      } catch (err) {
        console.error('Session registration failed:', err);
      }
    }

    // create websocket connection
    connect();

    const watchId = navigator.geolocation.watchPosition(
      (pos) => {
        locationError = false;
        const lat = pos.coords.latitude;
        const lon = pos.coords.longitude;
        userCoords.set({ lat, lon });
        sendLocationPing(lat, lon);
      },
      (err) => {
        console.error('Geolocation error:', err);
        if (err.code === 1) {
          locationError = true;
        }
      },
      {
        enableHighAccuracy: true,
        timeout: 15000,
        maximumAge: 60000,
      }
    );

    const pingInterval = setInterval(() => {
      const coords = $userCoords;
      if (coords.lat !== 0 && $isConnected) {
        sendLocationPing(coords.lat, coords.lon);
      }
    }, 15000);

    return () => {
      navigator.geolocation.clearWatch(watchId);
      clearInterval(pingInterval);
    };
  });
</script>

<div
  class="mx-auto flex h-screen max-w-2xl flex-col bg-sky-100 p-4 lg:border-r-4 lg:border-l-4 lg:border-solid lg:border-sky-800"
>
  <ChatHeader {$isConnected} {nearbyCounter} {locationError} />
  <MessageList messages={$messages}/>
  <ChatInput onSend={sendMessage} {isConnected} />
</div>
