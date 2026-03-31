import { messages, nearbyCount } from '$lib/stores/chat';
import type { ChatMessageRecieved, SystemMessage } from '$lib/types/types';

export function handleChatMessage(data: any) {
  const receivedMessage: ChatMessageRecieved = {
    type: 'chat',
    id: data.id,
    text: data.text,
    displayName: data.displayName || data.display_name,
    avatarUrl: data.avatarUrl || data.avatar_url,
    timestamp: data.timestamp,
    status: 'sent'
  };

  messages.update((prev) => {
    // Update existing message or add new one
    const existingIndex = prev.findIndex((m) => m.id === receivedMessage.id);
    if (existingIndex !== -1) {
      const updated = [...prev];
      updated[existingIndex] = receivedMessage;
      return updated;
    }
    return [...prev, receivedMessage];
  });
}

export function handleSystemMessage(data: any) {
  console.log(data);
  const systemMessage: SystemMessage = {
    type: 'system',
    id: data.id || crypto.randomUUID(),
    code: data.code,
    text: data.text,
    timestamp: data.timestamp || Date.now(),
    isConsoleMessage: data.is_console_message
  };

  if (systemMessage.isConsoleMessage) {
    console.log("System:", systemMessage.text)
    return
  }

  messages.update((prev) => [...prev, systemMessage]);
}

export function handleNearbyUpdate(data: any) {
  nearbyCount.update(count => Math.max(0, count + data.delta));
}
