<script lang="ts">
  import { onMount } from 'svelte';
  import { messages, isConnected, connect, sendMessage } from "$lib/stores/chat";

  let newMessage = $state("");
  let displayName = $state("anon");

  onMount(() => {
    connect(displayName);
  })

  function handleSend() {
    if (newMessage.trim()) {
      sendMessage(newMessage);
      newMessage = "";
    }
  }
</script>

<div class="flex flex-col h-screen max-w-2xl-mx-auto p-4">

  <div class="bg-base-100 shadow-xl rounded-box mb-4">
    <div class="flex-1">
      <span class="text-xl font-bold px-4">Echo</span>
    </div>
    <div class="flex-none px-4">
      <div class="badge {$isConnected ? 'badge-success' : 'badge-error'} gap-2">
        {$isConnected ? 'Online' : 'Offline'}
      </div>
    </div>
  </div>

  <div class="flex-1 overflow-y-auto p-4 space-y-4 bg-base-200 rounded-box">
    {#each $messages as msg}
      <div class="chat {msg.displayName === userName ? 'chat-end' : 'chat-start'}">
        <div class="chat-header">
          {msg.displayName}
          <time class="text-xs opacity-50">{new Date(msg.timestamp).toLocaleTimeString()}</time>
        </div>
        <div class="chat-bubble chat-bubble-primary">{msg.text}</div>
      </div>
    {/each}
  </div>

  <div class="mt-4 flex gap-2">
    <input
      type="text"
      bind:value={newMessage}
      placeholder="Send a message"
      class="input input-bordered flex-1"
      onkeydown={(e) => e.key === 'Enter' && handleSend()}
    />
    <button class="btn btn-primary" onclick={handleSend}>Send</button>
  </div>
 </div>
