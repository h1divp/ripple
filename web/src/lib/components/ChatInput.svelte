<script lang="ts">
  let { onSend, isConnected } = $props();
  import { IconSend2, IconCloudOff } from '@tabler/icons-svelte';

  let newMessage = $state('');
  function handleSend() {
    if (newMessage.trim()) {
      sendMessage(newMessage);
      newMessage = '';
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
    class="textarea textarea-bordered min-h-[2.5rem] flex-1 resize-none rounded-lg border-gray-300 focus:border-sky-500 focus:ring-1 focus:ring-gray-500 focus:outline-none disabled:cursor-not-allowed disabled:bg-gray-100 disabled:placeholder-gray-400"
  ></textarea>

  <button
    class="btn btn-primary cursor-pointer text-sky-900 disabled:cursor-not-allowed"
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
