import { BaseApiClient } from '../../clients/base-client';
import { ApiResponse } from '../../utils/mappers';
import { Validators } from '../../utils/validators';
import { TokenManager } from '../../core/token-manager';

export interface LoginRequest {
  email: string;
  password: string;
  rememberMe?: boolean;
}

export interface RegisterRequest {
  email: string;
  password: string;
  username?: string;
  firstName?: string;
  lastName?: string;
  acceptTerms?: boolean;
}

export interface AuthResponse {
  success: boolean;
  accessToken: string;
  refreshToken: string;
  user: {
    id: string;
    email: string;
    username?: string;
    firstName?: string;
    lastName?: string;
    avatar?: string;
    roles?: string[];
  };
}

export interface PasswordResetRequest {
  email: string;
}

export interface PasswordResetConfirm {
  token: string;
  newPassword: string;
}

export interface ChangePasswordRequest {
  currentPassword: string;
  newPassword: string;
}

export class AuthService extends BaseApiClient {
  constructor(axios: any) {
    super(axios, '/auth');
  }

  /**
   * Login user
   */
  async login(request: LoginRequest): Promise<ApiResponse<AuthResponse>> {
    // Validate input
    const emailValidation = Validators.email(request.email);
    if (!emailValidation.isValid) {
      throw new Error(emailValidation.errors[0]);
    }

    const response = await this.post<AuthResponse>('/login', request);

    if (response.success && response.data) {
      TokenManager.setTokens(response.data.accessToken, response.data.refreshToken);
    }

    return response;
  }

  /**
   * Register new user
   */
  async register(request: RegisterRequest): Promise<ApiResponse<AuthResponse>> {
    // Validate input
    const emailValidation = Validators.email(request.email);
    const passwordValidation = Validators.password(request.password);

    if (!emailValidation.isValid) {
      throw new Error(emailValidation.errors[0]);
    }
    if (!passwordValidation.isValid) {
      throw new Error(passwordValidation.errors[0]);
    }

    const response = await this.post<AuthResponse>('/register', request);

    if (response.success && response.data) {
      TokenManager.setTokens(response.data.accessToken, response.data.refreshToken);
    }

    return response;
  }

  /**
   * Logout user
   */
  async logout(): Promise<ApiResponse<void>> {
    try {
      const response = await this.post<void>('/logout');
      TokenManager.clearTokens();
      return response;
    } catch (error) {
      // Clear tokens even if logout fails
      TokenManager.clearTokens();
      throw error;
    }
  }

  /**
   * Refresh access token
   */
  async refreshToken(): Promise<ApiResponse<AuthResponse>> {
    const refreshToken = TokenManager.getRefreshToken();
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    const response = await this.post<AuthResponse>('/refresh', { refreshToken });

    if (response.success && response.data) {
      TokenManager.setTokens(response.data.accessToken, response.data.refreshToken);
    }

    return response;
  }

  /**
   * Get current user
   */
  async getCurrentUser(): Promise<ApiResponse<AuthResponse['user']>> {
    return this.get<AuthResponse['user']>('/me');
  }

  /**
   * Request password reset
   */
  async requestPasswordReset(request: PasswordResetRequest): Promise<ApiResponse<void>> {
    const emailValidation = Validators.email(request.email);
    if (!emailValidation.isValid) {
      throw new Error(emailValidation.errors[0]);
    }

    return this.post<void>('/password/reset', request);
  }

  /**
   * Confirm password reset
   */
  async confirmPasswordReset(request: PasswordResetConfirm): Promise<ApiResponse<void>> {
    const passwordValidation = Validators.password(request.newPassword);
    if (!passwordValidation.isValid) {
      throw new Error(passwordValidation.errors[0]);
    }

    return this.post<void>('/password/confirm', request);
  }

  /**
   * Change password
   */
  async changePassword(request: ChangePasswordRequest): Promise<ApiResponse<void>> {
    const passwordValidation = Validators.password(request.newPassword);
    if (!passwordValidation.isValid) {
      throw new Error(passwordValidation.errors[0]);
    }

    return this.post<void>('/password/change', request);
  }

  /**
   * Verify email
   */
  async verifyEmail(token: string): Promise<ApiResponse<void>> {
    return this.post<void>('/email/verify', { token });
  }

  /**
   * Resend verification email
   */
  async resendVerificationEmail(): Promise<ApiResponse<void>> {
    return this.post<void>('/email/resend');
  }

  /**
   * Check if user is authenticated
   */
  isAuthenticated(): boolean {
    return TokenManager.isAccessTokenValid();
  }

  /**
   * Get user ID from token
   */
  getUserId(): string | null {
    const token = TokenManager.getAccessToken();
    return TokenManager.getUserIdFromToken(token);
  }
}