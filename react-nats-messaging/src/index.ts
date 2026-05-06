// Core exports
export { NatsProvider } from './core/NatsProvider';
export type { NatsProviderProps } from './core/NatsProvider';
export { NatsConnectionManager } from './core/NatsConnection';
export { MessageEncoder, createMessageType } from './core/MessageEncoder';
export type { EncoderConfig } from './core/MessageEncoder';

// Type exports
export type {
  NatsConfig,
  NatsContextValue,
  ConnectionStatus,
  PublishOptions,
  RequestOptions,
  MessageHandler,
  ProtobufConfig,
  MessageEnvelope,
  SubscriptionOptions,
  UseSubscriptionResult,
  UsePublishResult,
  MessageType,
  StreamMessage,
  WebsocketMessageData
} from './core/types';

// Hook exports
export { useNats } from './hooks/useNats';
export { useSubscription } from './hooks/useSubscription';
export { usePublish } from './hooks/usePublish';
export { useConnectionStatus } from './hooks/useConnectionStatus';
export { useMessageHandler } from './hooks/useMessageHandler';
export { useCommentsActions } from './hooks/useCommentsActions';
export type { CommentPayload, UseCommentsActionsOptions } from './hooks/useCommentsActions';
export { useChatHistory } from './hooks/useChatHistory';
export type { ChatMessage, UseChatHistoryOptions } from './hooks/useChatHistory';

// Component exports
export { ConnectionIndicator } from './components/ConnectionIndicator';
export type { ConnectionIndicatorProps } from './components/ConnectionIndicator';
export { MessageDebugger } from './components/MessageDebugger';
export type { MessageDebuggerProps } from './components/MessageDebugger';

// API exports
export { createCommentsApiClient } from './api/comments';
export type { Comment, CommentsApiClient } from './api/comments';
export { createMessagingApiClient } from './api/messaging';
export type { Message, Conversation, MessagingApiClient } from './api/messaging';

// Utility exports
export { MessageDeduplicator } from './utils/deduplication';
export { withRetry, sleep, ExponentialBackoff } from './utils/retry';
export type { RetryOptions } from './utils/retry';

// Proto type exports (for convenience)
export { ProtoHelpers } from './proto/types';
export type { 
  StreamMessage as ProtoStreamMessage,
  WebsocketMessageData as ProtoWebsocketMessageData,
  SendMessage as ProtoSendMessage,
  AddComment as ProtoAddComment
} from './proto/types';