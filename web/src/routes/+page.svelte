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
    userId,
  } from '$lib/stores/chat';
  import ChatHeader from '$lib/components/ChatHeader.svelte';
  import MessageList from '$lib/components/MessageList.svelte';
  import ChatInput from '$lib/components/ChatInput.svelte';

  let nearbyCounter = $state(0);
  let locationError = $state(true);

  onMount(() => {
    connect($userId);

    const watchId = navigator.geolocation.watchPosition(
      (pos) => {
        locationError = false;
        const lat = pos.coords.latitude;
        const lon = pos.coords.longitude;
        userCoords.set({ lat, lon });
        sendLocationPing(lat, lon);
      },
      (err) => {
        console.error("Geolocation error:", err);
        if (err.code === 1) {
          locationError = true;
        }
      },
      { enableHighAccuracy: true }
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
  <ChatHeader {isConnected} {nearbyCounter} {locationError} />
  <MessageList messages={$messages} currentUserId={$userId} />
  <ChatInput onSend={sendMessage} {isConnected} />
</div>
