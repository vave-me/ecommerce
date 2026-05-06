import AssistantService from '../AssistantService';
import AssistantServiceWithRetry, { ConversationNotFoundError, AssistantNotFoundError } from '../AssistantServiceWithRetry';
import axiosInstance from '../../../api/axiosInstance';

// Mock axios instance
jest.mock('../../../api/axiosInstance');

describe('AssistantService', () => {
    let service;

    beforeEach(() => {
        // Clear singleton instance before each test
        AssistantService.clearInstance();
        service = AssistantService.getInstance();
        jest.clearAllMocks();
    });

    describe('Singleton Pattern', () => {
        it('should return the same instance', () => {
            const instance1 = AssistantService.getInstance();
            const instance2 = AssistantService.getInstance();
            expect(instance1).toBe(instance2);
        });

        it('should create new instance after clearing', () => {
            const instance1 = AssistantService.getInstance();
            AssistantService.clearInstance();
            const instance2 = AssistantService.getInstance();
            expect(instance1).not.toBe(instance2);
        });
    });

    describe('Assistant Management', () => {
        it('should fetch assistants successfully', async () => {
            const mockData = {
                assistants: [
                    { id: '1', name: 'Assistant 1' },
                    { id: '2', name: 'Assistant 2' }
                ]
            };
            axiosInstance.get.mockResolvedValueOnce({ data: mockData });

            const result = await service.getAssistants({ page: 1, limit: 20 });

            expect(axiosInstance.get).toHaveBeenCalledWith('/assistants', {
                params: { page: 1, limit: 20 }
            });
            expect(result).toEqual({
                success: true,
                data: mockData
            });
        });

        it('should handle error when fetching assistants', async () => {
            const error = new Error('Network error');
            axiosInstance.get.mockRejectedValueOnce(error);

            const result = await service.getAssistants();

            expect(result).toEqual({
                success: false,
                error: 'Network error',
                status: 500,
                data: null
            });
        });

        it('should fetch single assistant', async () => {
            const mockData = { assistant: { id: '1', name: 'Assistant 1' } };
            axiosInstance.get.mockResolvedValueOnce({ data: mockData });

            const result = await service.getAssistant('1');

            expect(axiosInstance.get).toHaveBeenCalledWith('/assistants/1');
            expect(result).toEqual({
                success: true,
                data: mockData
            });
        });

        it('should update assistant configuration', async () => {
            const config = { temperature: 0.7, maxTokens: 1000 };
            const mockData = { success: true };
            axiosInstance.put.mockResolvedValueOnce({ data: mockData });

            const result = await service.updateAssistantConfig('1', config);

            expect(axiosInstance.put).toHaveBeenCalledWith('/assistants/1', config);
            expect(result).toEqual({
                success: true,
                data: mockData
            });
        });
    });

    describe('Conversation Management', () => {
        it('should create conversation successfully', async () => {
            const mockData = {
                conversationId: 'conv-123',
                conversation: { id: 'conv-123' }
            };
            axiosInstance.post.mockResolvedValueOnce({ data: mockData });

            const result = await service.createConversation('assistant-1', { topic: 'test' });

            expect(axiosInstance.post).toHaveBeenCalledWith('/assistants/conversations', {
                assistantId: 'assistant-1',
                initialContext: {
                    topic: 'test',
                    source: 'web_app',
                    timestamp: expect.any(String)
                }
            });
            expect(result.success).toBe(true);
            expect(result.data.conversationId).toBe('conv-123');
        });

        it('should handle conversation_id field in response', async () => {
            const mockData = {
                conversation_id: 'conv-456',
                conversation: { id: 'conv-456' }
            };
            axiosInstance.post.mockResolvedValueOnce({ data: mockData });

            const result = await service.createConversation('assistant-1');

            expect(result.data.conversationId).toBe('conv-456');
        });

        it('should check if conversation exists', async () => {
            axiosInstance.get.mockResolvedValueOnce({ data: {} });

            const exists = await service.conversationExists('conv-123');

            expect(axiosInstance.get).toHaveBeenCalledWith('/assistants/conversations/conv-123');
            expect(exists).toBe(true);
        });

        it('should return false for non-existent conversation', async () => {
            const error = new Error('Not found');
            error.response = { status: 404 };
            axiosInstance.get.mockRejectedValueOnce(error);

            const exists = await service.conversationExists('conv-999');

            expect(exists).toBe(false);
        });

        it('should throw error for other errors in conversationExists', async () => {
            const error = new Error('Server error');
            error.response = { status: 500 };
            axiosInstance.get.mockRejectedValueOnce(error);

            await expect(service.conversationExists('conv-123')).rejects.toThrow('Server error');
        });
    });

    describe('Messaging', () => {
        it('should send message successfully', async () => {
            const mockExistsResponse = { data: {} };
            const mockSendResponse = {
                data: {
                    userMessageId: 'msg-1',
                    messageId: 'msg-2',
                    response: 'AI response',
                    actions: [],
                    confidence: 0.95
                }
            };

            axiosInstance.get.mockResolvedValueOnce(mockExistsResponse); // conversationExists
            axiosInstance.post.mockResolvedValueOnce({ data: mockSendResponse });

            const result = await service.sendMessage('conv-123', 'Hello AI', { mood: 'friendly' });

            expect(axiosInstance.get).toHaveBeenCalledWith('/assistants/conversations/conv-123');
            expect(axiosInstance.post).toHaveBeenCalledWith(
                '/assistants/conversations/conv-123/chat',
                {
                    message: 'Hello AI',
                    context: {
                        mood: 'friendly',
                        timestamp: expect.any(String)
                    }
                }
            );
            expect(result.success).toBe(true);
            expect(result.data.userMessage.content).toBe('Hello AI');
            expect(result.data.assistantMessage.content).toBe('AI response');
        });

        it('should fail if conversation does not exist', async () => {
            const error = new Error('Not found');
            error.response = { status: 404 };
            axiosInstance.get.mockRejectedValueOnce(error);

            const result = await service.sendMessage('conv-999', 'Hello');

            expect(result.success).toBe(false);
            expect(result.error).toContain('Failed to send message');
        });

        it('should get messages with pagination', async () => {
            const mockData = {
                messages: [
                    { id: '1', content: 'Hello', role: 'USER' },
                    { id: '2', content: 'Hi there', role: 'ASSISTANT' }
                ]
            };
            axiosInstance.get.mockResolvedValueOnce({ data: mockData });

            const result = await service.getMessages('conv-123', { page: 2, limit: 10 });

            expect(axiosInstance.get).toHaveBeenCalledWith(
                '/assistants/conversations/conv-123/messages',
                { params: { page: 2, limit: 10 } }
            );
            expect(result).toEqual({
                success: true,
                data: mockData
            });
        });
    });

    describe('Error Handling', () => {
        it('should handle canceled requests', async () => {
            const error = new Error('Request canceled');
            error.code = 'ERR_CANCELED';
            axiosInstance.get.mockRejectedValueOnce(error);

            const result = await service.getAssistants();

            expect(result).toEqual({
                success: false,
                error: 'Request canceled',
                canceled: true,
                status: 0,
                data: null
            });
        });

        it('should use error response message if available', async () => {
            const error = new Error('Generic error');
            error.response = {
                status: 400,
                data: { message: 'Invalid request parameters' }
            };
            axiosInstance.get.mockRejectedValueOnce(error);

            const result = await service.getAssistants();

            expect(result.error).toBe('Invalid request parameters');
            expect(result.status).toBe(400);
        });

        it('should log errors to console', async () => {
            const consoleSpy = jest.spyOn(console, 'error').mockImplementation();
            const error = new Error('Test error');
            axiosInstance.get.mockRejectedValueOnce(error);

            await service.getAssistants();

            expect(consoleSpy).toHaveBeenCalledWith(
                'AssistantService Error:',
                expect.objectContaining({
                    message: 'Test error',
                    status: 500
                })
            );

            consoleSpy.mockRestore();
        });
    });
});

describe('AssistantServiceWithRetry', () => {
    let service;

    beforeEach(() => {
        AssistantServiceWithRetry.clearInstance();
        service = AssistantServiceWithRetry.getInstance();
        jest.clearAllMocks();
        // Speed up tests by reducing delays
        service.configureRetry({
            maxRetries: 2,
            initialDelay: 10,
            maxDelay: 50
        });
    });

    describe('Retry Logic', () => {
        it('should retry on retryable errors', async () => {
            const error = new Error('Network error');
            error.response = { status: 503 };

            axiosInstance.get
                .mockRejectedValueOnce(error)
                .mockResolvedValueOnce({ data: { assistants: [] } });

            const result = await service.getAssistants();

            expect(axiosInstance.get).toHaveBeenCalledTimes(2);
            expect(result.success).toBe(true);
        });

        it('should not retry on non-retryable errors', async () => {
            const error = new Error('Bad request');
            error.response = { status: 400 };

            axiosInstance.get.mockRejectedValueOnce(error);

            const result = await service.getAssistants();

            expect(axiosInstance.get).toHaveBeenCalledTimes(1);
            expect(result.success).toBe(false);
        });

        it('should respect max retries', async () => {
            const error = new Error('Server error');
            error.response = { status: 500 };

            axiosInstance.get.mockRejectedValue(error);

            const result = await service.getAssistants();

            expect(axiosInstance.get).toHaveBeenCalledTimes(2); // 1 initial + 1 retry (maxRetries = 2)
            expect(result.success).toBe(false);
        });

        it('should not retry conversation creation', async () => {
            const error = new Error('Server error');
            error.response = { status: 500 };

            axiosInstance.post.mockRejectedValueOnce(error);

            const result = await service.createConversation('assistant-1');

            expect(axiosInstance.post).toHaveBeenCalledTimes(1);
            expect(result.success).toBe(false);
        });
    });

    describe('Custom Error Types', () => {
        it('should throw ConversationNotFoundError', async () => {
            const error = new Error('Not found');
            error.response = { status: 404 };
            error.config = { url: '/assistants/conversations/conv-123/chat' };

            axiosInstance.get.mockResolvedValueOnce({}); // conversationExists
            axiosInstance.post.mockRejectedValue(error);

            await expect(service.sendMessage('conv-123', 'Hello')).rejects.toThrow(ConversationNotFoundError);
        });

        it('should throw AssistantNotFoundError', async () => {
            const error = new Error('Not found');
            error.response = { status: 404 };
            error.config = { url: '/assistants/assistant-999' };

            axiosInstance.get.mockRejectedValue(error);

            await expect(service.getAssistant('assistant-999')).rejects.toThrow(AssistantNotFoundError);
        });
    });

    describe('Configuration', () => {
        it('should allow retry configuration', () => {
            service.configureRetry({
                maxRetries: 5,
                initialDelay: 2000
            });

            expect(service.retryConfig.maxRetries).toBe(5);
            expect(service.retryConfig.initialDelay).toBe(2000);
        });
    });
});