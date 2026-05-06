export interface Message {
  id: string;
  conversationId: string;
  senderId: string;
  recipientId: string;
  itemId?: string;
  body: string;
  createdAt: string | number;
  isRead?: boolean;
}

export interface Conversation {
  id: string;
  senderId: string;
  recipientId: string;
  itemId?: string;
  lastMessage?: string;
  lastMessageTime?: string;
  unreadCount?: number;
}

export interface MessagingApiClient {
  getConversations: (userId: string) => Promise<Conversation[]>;
  startConversation: (senderId: string, recipientId: string, itemId?: string) => Promise<{ id: string }>;
  getMessagesByConversation: (conversationId: string) => Promise<Message[]>;
  getConversationByRecipientAndProduct: (recipientId: string, itemId: string) => Promise<Conversation[]>;
}

/**
 * Factory function to create a messaging API client
 * @param axiosInstance - Your configured axios instance
 * @param baseUrl - Base URL for messaging API (default: '/messages')
 */
export function createMessagingApiClient(
  axiosInstance: any,
  baseUrl: string = '/messages'
): MessagingApiClient {
  return {
    async getConversations(userId: string): Promise<Conversation[]> {
      if (!userId) return [];
      try {
        const response = await axiosInstance.get(
          `${baseUrl}/conversations?userId=${userId}`
        );
        return response.data.conversations;
      } catch (err) {
        console.error(`Error retrieving conversations for userId=${userId}`, err);
        throw err;
      }
    },

    async startConversation(
      senderId: string,
      recipientId: string,
      itemId?: string
    ): Promise<{ id: string }> {
      const body = { senderId, recipientId, itemId };
      try {
        const response = await axiosInstance.post(
          `${baseUrl}/conversations`,
          body
        );
        return response.data;
      } catch (err) {
        console.error(
          `Error starting conversation (sender=${senderId}, recipient=${recipientId}, item=${itemId}):`,
          err
        );
        throw err;
      }
    },

    async getMessagesByConversation(conversationId: string): Promise<Message[]> {
      if (!conversationId) return [];
      try {
        const response = await axiosInstance.get(
          `${baseUrl}/conversations/${conversationId}/message?page=1&limit=100`
        );
        return response.data.messages;
      } catch (err) {
        console.error(`Error retrieving messages for conversationId=${conversationId}`, err);
        throw err;
      }
    },

    async getConversationByRecipientAndProduct(
      recipientId: string,
      itemId: string
    ): Promise<Conversation[]> {
      try {
        const response = await axiosInstance.get(
          `${baseUrl}/conversations/recipient/${recipientId}/item/${itemId}`
        );
        const convo = response.data.conversation;
        return convo ? [convo] : [];
      } catch (err: any) {
        if (err.response && err.response.status === 404) {
          console.warn(`No conversation found for recipientId=${recipientId}, itemId=${itemId}`);
          return [];
        }
        console.error(`Error retrieving conversation:`, err);
        throw err;
      }
    }
  };
}