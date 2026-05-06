export interface Comment {
  id: string;
  senderId: string;
  itemId: string;
  content: string;
  categoryId?: string;
  parentId?: string;
  createdAt: string;
  replies?: Comment[];
}

export interface CommentsApiClient {
  getCommentsByItem: (itemId: string) => Promise<Comment[]>;
  getCommentsBySender: (senderId: string) => Promise<Comment[]>;
}

/**
 * Factory function to create a comments API client
 * @param axiosInstance - Your configured axios instance
 * @param baseUrl - Base URL for comments API (default: '/comments')
 */
export function createCommentsApiClient(
  axiosInstance: any,
  baseUrl: string = '/comments'
): CommentsApiClient {
  return {
    async getCommentsByItem(itemId: string): Promise<Comment[]> {
      const response = await axiosInstance.get(`${baseUrl}/${itemId}/all`);
      return response.data;
    },

    async getCommentsBySender(senderId: string): Promise<Comment[]> {
      const response = await axiosInstance.get(`${baseUrl}/sender/${senderId}`);
      return response.data;
    }
  };
}