import type { SystemMessage } from "$lib/types/types";

const now = new Date();

export const joinMessage: SystemMessage = {
  id: crypto.randomUUID(),
  type: 'system',
  text: `joined at ${now.toLocaleString([], { hour: '2-digit', minute: '2-digit' })}`,
} as SystemMessage;
