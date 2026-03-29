<script lang="ts">
  import { IconCheck, IconLoader2, IconAlertCircle } from '@tabler/icons-svelte';
  import type { ChatMessageRecieved } from '$lib/types/types';
  
  let { message, isMe, showDetails }: {
    message: ChatMessageRecieved;
    isMe: boolean;
    showDetails: boolean;
  } = $props();
  
  const date = $derived(new Date(message.timestamp ?? Date.now()));
  const timeString = $derived(date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }));
  const fullDateString = $derived(
    date.toLocaleString([], { dateStyle: 'short', timeStyle: 'short' })
  );
  const isSending = $derived(message.status === 'sending');
</script>

<div class="chat">
  <div class="group flex flex-row items-start py-0 hover:bg-black/5">
    <div class="w-14 shrink-0">
      {#if showDetails}
        <div class="chat-image avatar py-1">
          <div class="w-10 overflow-hidden rounded-full border-sky-900 bg-white">
            <img src={message.avatarUrl} alt="User Avatar" class="h-full w-full object-cover" />
          </div>
        </div>
      {:else}
        <div
          class="mt-1 hidden text-center text-xs leading-tight text-gray-600 select-none group-hover:block text-nowrap"
        >
          {showDetails ? fullDateString : timeString}
        </div>
      {/if}
    </div>
    <div class="flex flex-col overflow-hidden">
      {#if showDetails}
        <div class="flex items-baseline gap-2">
          <span class="text-md font-bold text-sky-900">
            {message.displayName ?? 'Anonymous'}
          </span>
          <span class="text-xs text-gray-500">
            {fullDateString}
          </span>
        </div>
      {/if}
      <div
        class="text-md leading-tight wrap-break-word duration-0"
        class:text-sky-900={isSending}
        class:opacity-70={isSending}
        class:text-gray-800={!isSending}
      >
        {message.text}
      </div>
    </div>
  </div>
</div>
