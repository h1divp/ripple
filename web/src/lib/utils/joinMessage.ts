const now = new Date();

export const joinMessage: Message = {
  id: crypto.randomUUID(),
  senderId: crypto.randomUUID(),
  lat: 0,
  lon: 0,
  status: 'sent',
  type: 'system',
  text: `joined at ${now.toLocaleString([], { hour: '2-digit', minute: '2-digit' })}`,
};
