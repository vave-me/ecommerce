import { MessageType } from './types';

export interface EncoderConfig {
  streamMessageType?: any; // Protobuf type for StreamMessage
  websocketMessageType?: any; // Protobuf type for WebsocketMessageData
  useStreamWrapper?: boolean;
  useWebsocketWrapper?: boolean;
}

export class MessageEncoder {
  private config: EncoderConfig;

  constructor(config: EncoderConfig = {}) {
    this.config = {
      useStreamWrapper: true,
      useWebsocketWrapper: true,
      ...config
    };
  }

  encode<T>(messageType: MessageType<T>, data: T, metadata?: Record<string, any>): Uint8Array {
    try {
      // Validate messageType
      if (!messageType || typeof messageType.encode !== 'function') {
        throw new Error('Invalid messageType: missing encode function');
      }

      // Encode the domain message
      let encoded = messageType.encode(data);

      // Wrap in WebsocketMessageData if configured
      if (this.config.useWebsocketWrapper && this.config.websocketMessageType) {
        const wsData = this.config.websocketMessageType.create({
          payload: encoded,
          occurredAt: new Date()
        });
        encoded = this.config.websocketMessageType.encode(wsData).finish();
      }

      // Wrap in StreamMessage if configured
      if (this.config.useStreamWrapper && this.config.streamMessageType) {
        const streamMsg = this.config.streamMessageType.create({
          id: this.generateId(),
          name: messageType.subject,
          data: encoded,
          metadata: metadata,
          sentAt: new Date()
        });
        encoded = this.config.streamMessageType.encode(streamMsg).finish();
      }

      return encoded;
    } catch (error) {
      throw new Error(`Failed to encode message: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  decode<T>(messageType: MessageType<T>, data: Uint8Array): T {
    let decoded = data;

    // Unwrap StreamMessage if configured
    if (this.config.useStreamWrapper && this.config.streamMessageType) {
      const streamMsg = this.config.streamMessageType.decode(decoded);
      decoded = streamMsg.data;
    }

    // Unwrap WebsocketMessageData if configured
    if (this.config.useWebsocketWrapper && this.config.websocketMessageType) {
      const wsData = this.config.websocketMessageType.decode(decoded);
      decoded = wsData.payload;
    }

    // Decode the domain message
    return messageType.decode(decoded);
  }

  private generateId(): string {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }
}

// Factory function for creating message types
export function createMessageType<T>(
  subject: string,
  protoType: any
): MessageType<T> {
  return {
    subject,
    encode: (message: T) => {
      const protoMessage = protoType.create(message);
      return protoType.encode(protoMessage).finish();
    },
    decode: (data: Uint8Array) => {
      return protoType.decode(data) as T;
    }
  };
}