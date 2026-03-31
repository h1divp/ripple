<script lang="ts">
  import { PUBLIC_MAX_MESSAGE_CHARACTERS } from '$env/static/public';
  import { isConnected, rateLimitEndTime } from '$lib/stores/chat';
  import { hasLocation } from '$lib/stores/location';
  import { IconSend2, IconCloudOff } from '@tabler/icons-svelte';
  import { onMount } from 'svelte';
  
  let { onSend }: { onSend: (text: string) => void } = $props();
  
  let newMessage = $state('');
  let textarea: HTMLTextAreaElement;
  let timeRemaining = $state(0);

  const isRateLimited = $derived($rateLimitEndTime > 0)
  const canSend = $derived($isConnected && $hasLocation && !isRateLimited);
  const remainingChars = $derived(PUBLIC_MAX_MESSAGE_CHARACTERS - newMessage.length);
  const isNearLimit = $derived(remainingChars <= 50);
  const isOverLimit = $derived(newMessage.length > PUBLIC_MAX_MESSAGE_CHARACTERS );

  onMount(() => {
    const interval = setInterval(() => {
      if ($rateLimitEndTime.getTime() > 0) {
        const now = new Date().getTime();
        const remaining = Math.max(0, Math.ceil(($rateLimitEndTime.getTime() - now) / 1000));
        timeRemaining = remaining;
        
        // Clear rate limit when time expires
        if (remaining === 0) {
          rateLimitEndTime.set(new Date(0));
        }
      } else {
        const timeRemaining = 0;
      }
    }, 100); // Update every 100ms for smooth countdown

    return () => clearInterval(interval);
  });

  function formatTime(seconds: number): string {
    if (seconds <= 0) {
      return "";
    }
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return mins > 0 ? `${mins}:${secs.toString().padStart(2, '0')}` : `${secs}s`;
  }
  
  function handleSend() {
    const trimmed = newMessage.trim();
    if (trimmed && $isConnected) {
      onSend(trimmed);
      newMessage = '';
      if (textarea) textarea.style.height = 'auto';
    }
  }
  
  function autoResize() {
    if (!textarea) return;
    textarea.style.height = 'auto';
    const offset = textarea.offsetHeight - textarea.clientHeight;
    const newHeight = Math.max(40, Math.min(textarea.scrollHeight + offset, 150));
    textarea.style.height = newHeight + 'px';
  }
  
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      const trimmed = newMessage.trim();
      if (trimmed && trimmed.length <= PUBLIC_MAX_MESSAGE_CHARACTERS && canSend) {
        handleSend();
        if (textarea) textarea.style.height = 'auto';
      }
    }
  }
</script>

<div class="mt-4 flex gap-2">

  <textarea
    bind:this={textarea}
    bind:value={newMessage}
    oninput={autoResize}
    onkeydown={handleKeydown}
    disabled={!canSend}
    placeholder={!$isConnected ? "Disconnected" : 
                 !$hasLocation ? "Location required" :
                 isRateLimited ? `Rate limited - ${formatTime(timeRemaining)}` :
                 "Send a message"}
    rows="1"
    class="textarea textarea-bordered scrollbar-none flex-1 resize-none overflow-hidden overflow-y-auto rounded-lg border-gray-300 py-2.25 leading-5 focus:border-sky-500 focus:ring-1 focus:ring-gray-500 focus:outline-none disabled:cursor-not-allowed disabled:bg-gray-100"
    style="height: 40px; min-height: 40px;"
  ></textarea>
  <div class="flex flex-col justify-between">
  <button
    class={`btn btn-primary cursor-pointer mt-1 hover:[&>svg]:fill-sky-50 active:text-sky-700 text-sky-900 disabled:text-sky-900/50 disabled:cursor-not-allowed`}
    onclick={handleSend}
    disabled={!canSend || isOverLimit}
  >
    {#if $isConnected}
      <IconSend2 size={30} />
    {:else}
      <IconCloudOff size={30} />
    {/if}
  </button>
  {#if isNearLimit || isOverLimit}
    <div class="text-right text-sm font-bold" class:text-grey={isNearLimit} class:text-red-600={isOverLimit}>
      {remainingChars}
    </div>
  {/if}
  </div>
</div>
