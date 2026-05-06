export interface SupportTicket {
  id: string;
  userId: string;
  subject: string;
  description: string;
  category: SupportCategory;
  priority: SupportPriority;
  status: SupportStatus;
  assignedTo?: string;
  messages?: SupportMessage[];
  attachments?: string[];
  createdAt: string;
  updatedAt: string;
  resolvedAt?: string;
}

export interface SupportMessage {
  id: string;
  ticketId: string;
  userId: string;
  message: string;
  attachments?: string[];
  isInternal?: boolean;
  createdAt: string;
}

export type SupportCategory = 
  | 'technical'
  | 'billing'
  | 'account'
  | 'order'
  | 'product'
  | 'shipping'
  | 'refund'
  | 'other';

export type SupportPriority = 'low' | 'medium' | 'high' | 'urgent';
export type SupportStatus = 'open' | 'in_progress' | 'waiting_customer' | 'waiting_support' | 'resolved' | 'closed';

export interface CreateTicketData {
  subject: string;
  description: string;
  category: SupportCategory;
  priority?: SupportPriority;
  attachments?: File[];
}

export interface SupportSession {
  id: string;
  userId: string;
  status: 'active' | 'ended';
  startedAt: string;
  endedAt?: string;
  rating?: number;
  feedback?: string;
}

export interface SupportApiClient {
  // Ticket management
  createTicket: (data: CreateTicketData) => Promise<SupportTicket>;
  listTickets: (filters?: { status?: SupportStatus; category?: SupportCategory; page?: number; limit?: number }) => Promise<{ tickets: SupportTicket[]; total: number }>;
  getTicket: (ticketId: string) => Promise<SupportTicket>;
  updateTicket: (ticketId: string, data: Partial<SupportTicket>) => Promise<SupportTicket>;
  closeTicket: (ticketId: string, reason?: string) => Promise<void>;
  reopenTicket: (ticketId: string) => Promise<SupportTicket>;
  
  // Messages
  addMessage: (ticketId: string, message: string, attachments?: File[]) => Promise<SupportMessage>;
  getMessages: (ticketId: string) => Promise<SupportMessage[]>;
  
  // Support session
  startSupport: () => Promise<SupportSession>;
  endSupport: (sessionId: string, rating?: number, feedback?: string) => Promise<void>;
  getActiveSession: () => Promise<SupportSession | null>;
  
  // FAQ and knowledge base
  searchFAQ: (query: string) => Promise<any[]>;
  getFAQCategories: () => Promise<string[]>;
  getFAQByCategory: (category: string) => Promise<any[]>;
  
  // Utils
  uploadAttachment: (file: File) => Promise<string>;
  getSupportStats: () => Promise<{ openTickets: number; avgResponseTime: number; satisfaction: number }>;
}

export function createSupportApiClient(
  axiosInstance: any,
  baseUrl: string = '/api/support'
): SupportApiClient {
  return {
    // Ticket management
    async createTicket(data: CreateTicketData): Promise<SupportTicket> {
      const formData = new FormData();
      formData.append('subject', data.subject);
      formData.append('description', data.description);
      formData.append('category', data.category);
      if (data.priority) formData.append('priority', data.priority);
      
      if (data.attachments) {
        data.attachments.forEach((file, index) => {
          formData.append(`attachments[${index}]`, file);
        });
      }

      const response = await axiosInstance.post(`${baseUrl}/tickets`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
      return response.data;
    },

    async listTickets(filters) {
      const response = await axiosInstance.get(`${baseUrl}/tickets`, { params: filters });
      return response.data;
    },

    async getTicket(ticketId: string) {
      const response = await axiosInstance.get(`${baseUrl}/tickets/${ticketId}`);
      return response.data;
    },

    async updateTicket(ticketId: string, data: Partial<SupportTicket>) {
      const response = await axiosInstance.put(`${baseUrl}/tickets/${ticketId}`, data);
      return response.data;
    },

    async closeTicket(ticketId: string, reason?: string) {
      await axiosInstance.post(`${baseUrl}/tickets/${ticketId}/close`, { reason });
    },

    async reopenTicket(ticketId: string) {
      const response = await axiosInstance.post(`${baseUrl}/tickets/${ticketId}/reopen`);
      return response.data;
    },

    // Messages
    async addMessage(ticketId: string, message: string, attachments?: File[]) {
      const formData = new FormData();
      formData.append('message', message);
      
      if (attachments) {
        attachments.forEach((file, index) => {
          formData.append(`attachments[${index}]`, file);
        });
      }

      const response = await axiosInstance.post(
        `${baseUrl}/tickets/${ticketId}/messages`, 
        formData,
        { headers: { 'Content-Type': 'multipart/form-data' } }
      );
      return response.data;
    },

    async getMessages(ticketId: string) {
      const response = await axiosInstance.get(`${baseUrl}/tickets/${ticketId}/messages`);
      return response.data;
    },

    // Support session
    async startSupport() {
      const response = await axiosInstance.post(`${baseUrl}/session/start`);
      return response.data;
    },

    async endSupport(sessionId: string, rating?: number, feedback?: string) {
      await axiosInstance.post(`${baseUrl}/session/${sessionId}/end`, { rating, feedback });
    },

    async getActiveSession() {
      try {
        const response = await axiosInstance.get(`${baseUrl}/session/active`);
        return response.data;
      } catch (error: any) {
        if (error.response?.status === 404) {
          return null;
        }
        throw error;
      }
    },

    // FAQ and knowledge base
    async searchFAQ(query: string) {
      const response = await axiosInstance.get(`${baseUrl}/faq/search`, { params: { q: query } });
      return response.data;
    },

    async getFAQCategories() {
      const response = await axiosInstance.get(`${baseUrl}/faq/categories`);
      return response.data;
    },

    async getFAQByCategory(category: string) {
      const response = await axiosInstance.get(`${baseUrl}/faq/category/${category}`);
      return response.data;
    },

    // Utils
    async uploadAttachment(file: File) {
      const formData = new FormData();
      formData.append('file', file);
      
      const response = await axiosInstance.post(`${baseUrl}/attachments`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
      return response.data.url;
    },

    async getSupportStats() {
      const response = await axiosInstance.get(`${baseUrl}/stats`);
      return response.data;
    }
  };
}