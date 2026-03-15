<script lang="ts">
  import MessageBubble from './MessageBubble.svelte';
  let { messages, currentUserId } = $props();
  let scrollContainer: HTMLDivElement | undefined = $state();

  // Auto-scroll to bottom when new messages arrive
  $effect(() => {
    if (messages && scrollContainer) {
      scrollContainer.scrollTo({ 
        top: scrollContainer.scrollHeight, 
        behavior: 'smooth' 
      });
    }
  });
</script>

<div bind:this={scrollContainer} class="flex-1 overflow-y-auto p-4 space-y-1 bg-base-200">
  {#each messages as msg, i (msg.id)}
    {@const isFirstInGroup = i === 0 || messages[i - 1].senderId !== msg.senderId}
    <MessageBubble 
      message={msg} 
      isMe={msg.senderId === currentUserId} 
      showDetails={isFirstInGroup} 
    />
  {/each}
</div>
