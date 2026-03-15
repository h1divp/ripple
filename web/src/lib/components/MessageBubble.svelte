<script lang="ts">
  import { IconCheck, IconLoader2, IconAlertCircle } from '@tabler/icons-svelte';

  let { message, isMe, showDetails } = $props();

  const date = $derived(new Date(message.timestamp));
  const timeString = $derived(date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }));
  const fullDateString = $derived(date.toLocaleString([], { dateStyle: 'short', timeStyle: 'short' }));
  const avatarUrl = `https://api.dicebear.com/7.x/identicon/svg?seed=${message.avatarSeed}`;
</script>

<div class="chat">

  <div class="group flex flex-row items-start px-2 py-0 hover:bg-black/5">

  <div class="w-12 flex-shrink-0">
    {#if showDetails}
      <div class="chat-image avatar py-1">
        <div class="w-10 rounded-full border-sky-900 bg-white overflow-hidden">
          <img 
            src={avatarUrl} 
            alt="User Avatar" 
            class="h-full w-full object-cover"
          />
        </div>
      </div>
    {:else}
      <div class="hidden group-hover:block text-xs leading-tight text-gray-600 mt-1 text-center">
        {showDetails ? fullDateString : timeString}
      </div>
    {/if}
  </div>

  <div class="flex flex-col overflow-hidden">
    {#if showDetails}
      <div class="flex items-baseline gap-2">
        <span class="font-bold text-md text-sky-900 cursor-pointer">
          {message.displayName}
        </span>
        <span class="text-xs text-gray-500">
          {fullDateString}
        </span>
      </div>
    {/if}

    <div class="text-md text-gray-800 leading-tight break-words">
      {message.text}
    </div>
  </div>
</div>

  
</div>
