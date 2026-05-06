// TypeScript interfaces for the protobuf messages
// This avoids the need for protobuf compilation during library usage

export interface StreamMessage {
  id: string;
  name: string;
  data: Uint8Array;
  metadata?: Record<string, any>;
  sentAt: Date | string;
}

export interface WebsocketMessageData {
  payload: Uint8Array;
  occurredAt: Date | string;
}

export interface SendMessage {
  id: string;
  conversationId: string;
  userId: string;
  messageText: string;
  createdAt: string;
}

export interface AddComment {
  id: string;
  itemId: string;
  userId: string;
  parentId?: string;
  text: string;
  createdAt: string;
}

// Helper functions for encoding/decoding without protobuf
export const ProtoHelpers = {
  encodeStreamMessage(msg: StreamMessage): Uint8Array {
    // Simple JSON encoding for the library
    const json = JSON.stringify({
      ...msg,
      sentAt: msg.sentAt instanceof Date ? msg.sentAt.toISOString() : msg.sentAt
    });
    return new TextEncoder().encode(json);
  },

  decodeStreamMessage(data: Uint8Array): StreamMessage {
    const json = new TextDecoder().decode(data);
    const parsed = JSON.parse(json);
    return {
      ...parsed,
      sentAt: new Date(parsed.sentAt)
    };
  },

  encodeWebsocketMessageData(msg: WebsocketMessageData): Uint8Array {
    const json = JSON.stringify({
      ...msg,
      occurredAt: msg.occurredAt instanceof Date ? msg.occurredAt.toISOString() : msg.occurredAt
    });
    return new TextEncoder().encode(json);
  },

  decodeWebsocketMessageData(data: Uint8Array): WebsocketMessageData {
    const json = new TextDecoder().decode(data);
    const parsed = JSON.parse(json);
    return {
      ...parsed,
      occurredAt: new Date(parsed.occurredAt)
    };
  }
};