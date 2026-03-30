export type MessageType =
  | 'chat'
  | 'system'
  | 'join'
  | 'location_update'
  | 'nearby_update'
  | 'username_update'
  | 'icon_update';

export type MessageStatus = 'sending' | 'sent' | 'error';
export const SystemMessageErrorCode = 'too_many_messages';

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
export interface SystemMessage extends BaseMessage {
  type: 'system';
  code: string;
  id: string;
  text: string;
}
export interface JoinMessage extends BaseMessage {
  type: 'join';
  id: string;
  text: string;
}
export type ChatMessage = ChatMessageSend | ChatMessageRecieved;
export type DisplayMessage = ChatMessageRecieved | SystemMessage | JoinMessage;

export interface LocationUpdate extends BaseMessage {
  type: 'location_update';
  lat: number;
  lon: number;
}

export interface NearbyUpdate extends BaseMessage {
  type: 'nearby_update';
  delta: number;
}

export interface UsernameUpdate extends BaseMessage {
  type: 'username_update';
  username: string;
}

export interface IconUpdate extends BaseMessage {
  type: 'icon_update';
  icon: string;
}

