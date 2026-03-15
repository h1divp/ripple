<script lang="ts">
  import { onMount } from 'svelte';
  import { messages, isConnected, connect, sendMessage } from "$lib/stores/chat";
  import { IconSend2, IconCloudOff, IconUsers } from '@tabler/icons-svelte';

  let newMessage = $state("");
  let displayName = $state("anon");
  let nearbyCounter = $state(0);

  onMount(() => {
    connect(displayName);
  })

  function handleSend() {
    if (newMessage.trim()) {
      sendMessage(newMessage);
      newMessage = "";
    }
  }

  let textarea: HTMLTextAreaElement;
  function autoResize() {
    if (!textarea) return;
    textarea.style.height = 'auto';
    const newHeight = Math.min(textarea.scrollHeight, 150);
    textarea.style.height = newHeight + 'px';
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
      if (textarea) textarea.style.height = 'auto';
    }
  }
</script>

<div class="flex flex-col h-screen max-w-2xl mx-auto p-4 bg-sky-100 lg:border-l-4 lg:border-r-4 lg:border-solid lg:border-sky-800">
  <div class="mb-4 flex flex-row flex-wrap items-center justify-center gap-2 sm:justify-start">
    <div class="w-fit flex flex-row rounded-lg border-4 border-solid border-sky-800 p-1 gap-3 text-sky-900 bg-sky-100">

    {#if $isConnected}
    <div class="flex items-center justify-center rounded-lgp-2">
      <IconUsers size={20} class="mx-1" />
      <span class="text-xl font-bold pr-2">
        {nearbyCounter}
      </span>
    </div>
    {:else}
    <div class="flex items-center justify-center rounded-lgp-2">
      <span class="text-xl font-bold px-2">Loading</span>
    </div>
    {/if}
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
    <textarea
      bind:this={textarea}
      bind:value={newMessage}
      oninput={autoResize}
      onkeydown={handleKeydown}
      disabled={!$isConnected}
      type="text"
      placeholder="Send a message"
      rows="1"
      class="textarea textarea-bordered resize-none flex-1 rounded-lg min-h-[2.5rem] border-gray-300 focus:border-sky-500 focus:outline-none focus:ring-1 focus:ring-gray-500 disabled:bg-gray-100 disabled:cursor-not-allowed disabled:placeholder-gray-400"
    ></textarea>
    
    <button 
      class="btn btn-primary cursor-pointer disabled:cursor-not-allowed text-sky-900" 
      onclick={handleSend} 
      disabled={!$isConnected}
    >
      {#if $isConnected}
        <IconSend2 size={30} />
      {:else}
        <IconCloudOff size={30} />
      {/if}
    </button>

  </div>
 </div>
