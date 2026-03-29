<script lang="ts">
  import { onMount } from 'svelte';
  import { messages, isConnected } from '$lib/stores/chat';
  import { locationError, isInitialized, initError } from '$lib/stores/app';
  import { userDisplayName, userAvatarUrl } from '$lib/stores/user';
  import { sendMessage } from '$lib/services/websocket';
  import { initializeApp, initializeGeolocation } from '$lib/services/initialization';
  import ChatHeader from '$lib/components/ChatHeader.svelte';
  import MessageList from '$lib/components/MessageList.svelte';
  import ChatInput from '$lib/components/ChatInput.svelte';

  let geolocationManager: ReturnType<typeof initializeGeolocation>;

  onMount(async () => {
    if (geolocationManager) {
      console.warn('Geolocation manager already exists');
      geolocationManager.stop();
    }
    
    try {
      await initializeApp();
      geolocationManager = initializeGeolocation();
      geolocationManager.start();
    } catch (error) {
      console.error('Failed to initialize app:', error);
    }

    return () => {
      geolocationManager?.stop();
    };
  });

  let debugInterval: ReturnType<typeof setInterval>;
    onMount(() => {
    // Count active intervals every 5 seconds
    debugInterval = setInterval(() => {
    }, 5000);
    
    return () => {
      clearInterval(debugInterval);
    };
  });
</script>

{#if $initError}
  <div class="error">Failed to initialize: {$initError}</div>
{:else if !$isInitialized}
  <div class="loading">Loading...</div>
{:else}
  <div
    class="mx-auto flex h-screen max-w-2xl flex-col bg-sky-100 p-4 lg:border-r-4 lg:border-l-4 lg:border-solid lg:border-sky-800"
  >
    <ChatHeader locationError={$locationError} />
    <MessageList messages={$messages} />
    <ChatInput onSend={sendMessage} {isConnected} />
  </div>
{/if}
