import { AuthService } from '../src/services/auth/auth.service';
import { TokenManager } from '../src/core/token-manager';
import { AxiosError } from 'axios';

// Mock TokenManager
jest.mock('../src/core/token-manager');

describe('AuthService', () => {
  let authService: AuthService;
  let mockAxios: any;

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
    
    // Mock axios instance
    mockAxios = {
      post: jest.fn(),
      get: jest.fn(),
      put: jest.fn(),
      delete: jest.fn(),
    };

    authService = new AuthService(mockAxios);
  });

  describe('login', () => {
    it('should login successfully and store tokens', async () => {
      const loginData = {
        email: 'redacted-email@example.com',
        password: 'password123',
      };

      const mockResponse = {
        data: {
          success: true,
          accessToken: 'access-token',
          refreshToken: 'refresh-token',
          user: {
            id: '123',
            email: 'redacted-email@example.com',
            username: 'testuser',
          },
        },
      };

      mockAxios.post.mockResolvedValue(mockResponse);

      const result = await authService.login(loginData);

      expect(mockAxios.post).toHaveBeenCalledWith('/auth/login', loginData, expect.any(Object));
      expect(TokenManager.setTokens).toHaveBeenCalledWith('access-token', 'refresh-token');
      expect(result.success).toBe(true);
      expect(result.data?.user.email).toBe('redacted-email@example.com');
    });

    it('should throw error for invalid email', async () => {
      const loginData = {
        email: 'invalid-email',
        password: 'password123',
      };

      await expect(authService.login(loginData)).rejects.toThrow('Invalid email format');
      expect(mockAxios.post).not.toHaveBeenCalled();
    });

    it('should handle login failure', async () => {
      const loginData = {
        email: 'redacted-email@example.com',
        password: 'wrongpassword',
      };

      const error = new AxiosError('Unauthorized');
      error.response = {
        status: 401,
        data: { error: 'Invalid credentials' },
      } as any;

      mockAxios.post.mockRejectedValue(error);

      await expect(authService.login(loginData)).rejects.toMatchObject({
        statusCode: 401,
        userMessage: expect.stringContaining('logged in'),
      });
    });
  });

  describe('register', () => {
    it('should register successfully', async () => {
      const registerData = {
        email: 'redacted-email@example.com',
        password: 'Password123',
        username: 'newuser',
      };

      const mockResponse = {
        data: {
          success: true,
          accessToken: 'access-token',
          refreshToken: 'refresh-token',
          user: {
            id: '456',
            email: 'redacted-email@example.com',
            username: 'newuser',
          },
        },
      };

      mockAxios.post.mockResolvedValue(mockResponse);

      const result = await authService.register(registerData);

      expect(mockAxios.post).toHaveBeenCalledWith('/auth/register', registerData, expect.any(Object));
      expect(TokenManager.setTokens).toHaveBeenCalledWith('access-token', 'refresh-token');
      expect(result.success).toBe(true);
    });

    it('should validate password requirements', async () => {
      const registerData = {
        email: 'redacted-email@example.com',
        password: 'weak',
        username: 'testuser',
      };

      await expect(authService.register(registerData)).rejects.toThrow('Password must be at least 8 characters long');
    });
  });

  describe('logout', () => {
    it('should logout and clear tokens', async () => {
      mockAxios.post.mockResolvedValue({ data: { success: true } });

      await authService.logout();

      expect(mockAxios.post).toHaveBeenCalledWith('/auth/logout', undefined, expect.any(Object));
      expect(TokenManager.clearTokens).toHaveBeenCalled();
    });

    it('should clear tokens even if logout request fails', async () => {
      mockAxios.post.mockRejectedValue(new Error('Network error'));

      await expect(authService.logout()).rejects.toThrow();
      expect(TokenManager.clearTokens).toHaveBeenCalled();
    });
  });

  describe('refreshToken', () => {
    it('should refresh token successfully', async () => {
      (TokenManager.getRefreshToken as jest.Mock).mockReturnValue('old-refresh-token');

      const mockResponse = {
        data: {
          success: true,
          accessToken: 'new-access-token',
          refreshToken: 'new-refresh-token',
          user: { id: '123' },
        },
      };

      mockAxios.post.mockResolvedValue(mockResponse);

      const result = await authService.refreshToken();

      expect(mockAxios.post).toHaveBeenCalledWith(
        '/auth/refresh',
        { refreshToken: 'old-refresh-token' },
        expect.any(Object)
      );
      expect(TokenManager.setTokens).toHaveBeenCalledWith('new-access-token', 'new-refresh-token');
      expect(result.success).toBe(true);
    });

    it('should throw error if no refresh token available', async () => {
      (TokenManager.getRefreshToken as jest.Mock).mockReturnValue(null);

      await expect(authService.refreshToken()).rejects.toThrow('No refresh token available');
    });
  });

  describe('getCurrentUser', () => {
    it('should fetch current user', async () => {
      const mockUser = {
        id: '123',
        email: 'redacted-email@example.com',
        username: 'testuser',
      };

      mockAxios.get.mockResolvedValue({ data: mockUser });

      const result = await authService.getCurrentUser();

      expect(mockAxios.get).toHaveBeenCalledWith('/auth/me', expect.any(Object));
      expect(result.data).toEqual(mockUser);
    });
  });

  describe('changePassword', () => {
    it('should change password successfully', async () => {
      const changePasswordData = {
        currentPassword: 'oldPassword123',
        newPassword: 'NewPassword123',
      };

      mockAxios.post.mockResolvedValue({ data: { success: true } });

      await authService.changePassword(changePasswordData);

      expect(mockAxios.post).toHaveBeenCalledWith(
        '/auth/password/change',
        changePasswordData,
        expect.any(Object)
      );
    });

    it('should validate new password', async () => {
      const changePasswordData = {
        currentPassword: 'oldPassword123',
        newPassword: 'weak',
      };

      await expect(authService.changePassword(changePasswordData)).rejects.toThrow(
        'Password must be at least 8 characters long'
      );
    });
  });

  describe('utility methods', () => {
    it('should check authentication status', () => {
      (TokenManager.isAccessTokenValid as jest.Mock).mockReturnValue(true);
      expect(authService.isAuthenticated()).toBe(true);

      (TokenManager.isAccessTokenValid as jest.Mock).mockReturnValue(false);
      expect(authService.isAuthenticated()).toBe(false);
    });

    it('should get user ID from token', () => {
      (TokenManager.getAccessToken as jest.Mock).mockReturnValue('mock-token');
      (TokenManager.getUserIdFromToken as jest.Mock).mockReturnValue('user-123');

      expect(authService.getUserId()).toBe('user-123');
    });
  });
});