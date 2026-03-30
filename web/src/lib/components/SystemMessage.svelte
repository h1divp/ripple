<script lang="ts">
  let { message, isError }: { message: string, isError: boolean } = $props();
</script>

<div class="my-4 flex items-center gap-4">
    <span 
      class="italic text-center"
      class:text-red-800={isError}
      class:text-sky-900={!isError}
    >
      System: {message}
    </span>
</div>
