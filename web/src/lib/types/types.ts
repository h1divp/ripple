export type MessageType = 'chat' | 'system' | 'location_update';
export type MessageStatus = 'sending' | 'sent' | 'error';

export interface Message {
  id: string;
  senderId: string;
  type: MessageType;
  lat: number;
  lon: number;
  
  text?: string;
  displayName?: string;
  senderId?: string;
  avatarSeed?: string;
  timestamp?: number;
  status?: MessageStatus;
}
