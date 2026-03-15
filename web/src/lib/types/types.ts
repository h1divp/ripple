export type MessageType = 'chat' | 'system';
export type MessageStatus = 'sending' | 'sent' | 'error';

export interface Message {
  id: string;
  text: string;
  displayName: string;
  timestamp: number;
  status: MessageStatus;
  senderId: string;
  avatarSeed: string;
  type?: MessageType;
}
