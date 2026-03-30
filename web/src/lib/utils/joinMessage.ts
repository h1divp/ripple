import type { SystemMessage } from "$lib/types/types";

const now = new Date();

export const joinMessage: JoinMessage = {
  id: crypto.randomUUID(),
  type: 'join',
  text: `joined at ${now.toLocaleString([], { hour: '2-digit', minute: '2-digit' })}`,
} as SystemMessage;
