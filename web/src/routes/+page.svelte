<script lang="ts">
  import { onMount } from 'svelte';
  import {
    messages,
    isConnected,
    connect,
    sendMessage,
    userDisplayName,
    userId,
  } from '$lib/stores/chat';
  import ChatHeader from '$lib/components/ChatHeader.svelte';
  import MessageList from '$lib/components/MessageList.svelte';
  import ChatInput from '$lib/components/ChatInput.svelte';

  let displayName = $state('anon');
  let nearbyCounter = $state(0);

  onMount(() => {
    connect(displayName);
  });
</script>

<div
  class="mx-auto flex h-screen max-w-2xl flex-col bg-sky-100 p-4 lg:border-r-4 lg:border-l-4 lg:border-solid lg:border-sky-800"
>
  <ChatHeader {isConnected} {nearbyCounter} />
  <MessageList messages={$messages} currentUserId={$userId} />
  <ChatInput onSend={sendMessage} {isConnected} />
</div>
