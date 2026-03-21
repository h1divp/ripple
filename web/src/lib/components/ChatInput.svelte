<script lang="ts">
  let { onSend, isConnected } = $props();
  import { IconSend2, IconCloudOff } from '@tabler/icons-svelte';
  import { userDisplayName, userAvatarSeed  } from '$lib/stores/chat';

  let newMessage = $state('');
  function handleSend() {
    const trimmed = newMessage.trim();
    if (trimmed && $isConnected) {
      onSend(trimmed);
      newMessage = '';
      if (textarea) textarea.style.height = 'auto';
    }
  }

  let textarea: HTMLTextAreaElement;
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
    placeholder={$isConnected ? "Send a message" : "Disconnected"}
    rows="1"
    class="textarea textarea-bordered
        scrollbar-none flex-1
         resize-none
         overflow-hidden overflow-y-auto rounded-lg border-gray-300 py-[9px] leading-5 focus:border-sky-500 focus:ring-1 focus:ring-gray-500 focus:outline-none disabled:cursor-not-allowed
         disabled:bg-gray-100"
    style="height: 40px; min-height: 40px;"
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
