export type MessageType =
  | 'chat'
  | 'system'
  | 'location_update'
  | 'username_update'
  | 'icon_update';

export type MessageStatus = 'sending' | 'sent' | 'error';

interface BaseMessage {
  type: MessageType;
  id?: string;
  timestamp?: number;
}

export interface ChatMessageSend extends BaseMessage {
  type: 'chat';
  id: string;
  text: string;
  status: MessageStatus;
}
export interface ChatMessageRecieved extends BaseMessage {
  type: 'chat';
  id: string;
  text: string;
  displayName: string;
  avatarUrl: string;
  status: MessageStatus;
}
export type ChatMessage = ChatMessageSend | ChatMessageRecieved;

export interface SystemMessage extends BaseMessage {
  type: 'system';
  id: string;
  text: string;
}

export interface LocationUpdate extends BaseMessage {
  type: 'location_update';
  lat: number;
  lon: number;
}

export interface UsernameUpdate extends BaseMessage {
  type: 'username_update';
  username: string;
}

export interface IconUpdate extends BaseMessage {
  type: 'icon_update';
  icon: string;
}

