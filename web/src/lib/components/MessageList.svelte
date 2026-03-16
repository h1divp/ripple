<script lang="ts">
  import ChatMessage from './ChatMessage.svelte';
  import SystemMessage from './SystemMessage.svelte';

  let { messages, currentUserId } = $props();
  let scrollContainer: HTMLDivElement | undefined = $state();

  // Auto-scroll to bottom when new messages arrive
  $effect(() => {
    if (messages && scrollContainer) {
      scrollContainer.scrollTo({
        top: scrollContainer.scrollHeight,
        behavior: 'instant',
      });
    }
  });
</script>

<div
  bind:this={scrollContainer}
  class="bg-base-200 flex-1 space-y-1 overflow-y-auto
         lg:p-4"
>
  {#each messages as msg, i (msg.id)}
    {#if msg.type === 'system'}
      <SystemMessage message={msg.text} />
    {:else}
      {@const isFirstInGroup =
        i === 0 || messages[i - 1].senderId !== msg.senderId || messages[i - 1].type === 'system'}
      <ChatMessage
        message={msg}
        isMe={msg.senderId === currentUserId}
        showDetails={isFirstInGroup}
      />
    {/if}
  {/each}
</div>
