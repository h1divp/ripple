<script lang="ts">
  import ChatMessage from './ChatMessage.svelte';
  import SystemMessage from './SystemMessage.svelte';
  import { userDisplayName } from '$lib/stores/chat';

  let { messages } = $props();
  let scrollContainer: HTMLDivElement | undefined = $state();

  const displayMessages = $derived(messages.filter((m) => 
    m && m.type && m.type !== 'location_update'
  ));
  // Auto-scroll to bottom when new messages arrive
  $effect(() => {
    if (displayMessages && scrollContainer) {
      scrollContainer.scrollTo({
        top: scrollContainer.scrollHeight,
        behavior: 'instant',
      });
    }
  });
</script>

<div
  bind:this={scrollContainer}
  class="bg-base-200 custom-scrollbar flex-1 space-y-1
         overflow-y-auto lg:p-4"
>
  {#each displayMessages as msg, i (msg.id ?? `fallback-${i}`)}
  {@const isFirstInGroup =
    i === 0 ||
    (msg.type === 'chat' && displayMessages[i-1].displayName !== msg.displayName) ||
    displayMessages[i - 1].type === 'system'
  }
{console.log('Index:', i, 'Message:', msg, 'Type:', msg.type, 'isFirstInGroup:', isFirstInGroup)}
  {#if msg.type === 'system'}
    <SystemMessage message={msg.text} />
  {:else if msg.type === 'chat'}
    <ChatMessage
      message={msg}
      isMe={msg.displayName === $userDisplayName}
      showDetails={isFirstInGroup}
    />
  {/if}
{/each}
</div>
