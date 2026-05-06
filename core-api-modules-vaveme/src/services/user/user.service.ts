import { BaseApiClient } from '../../clients/base-client';
import { ApiResponse, PaginatedResponse } from '../../utils/mappers';
import { Validators } from '../../utils/validators';
import { Encoders } from '../../utils/encoders';

export interface User {
  id: string;
  email: string;
  username?: string;
  firstName?: string;
  lastName?: string;
  displayName?: string;
  avatar?: string;
  bio?: string;
  phone?: string;
  dateOfBirth?: string;
  address?: UserAddress;
  preferences?: UserPreferences;
  stats?: UserStats;
  roles?: string[];
  isActive?: boolean;
  isVerified?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface UserAddress {
  street?: string;
  city?: string;
  state?: string;
  country?: string;
  postalCode?: string;
}

export interface UserPreferences {
  language?: string;
  currency?: string;
  timezone?: string;
  emailNotifications?: boolean;
  pushNotifications?: boolean;
  marketingEmails?: boolean;
}

export interface UserStats {
  listingsCount?: number;
  purchasesCount?: number;
  reviewsCount?: number;
  averageRating?: number;
  followersCount?: number;
  followingCount?: number;
}

export interface UpdateUserRequest {
  username?: string;
  firstName?: string;
  lastName?: string;
  bio?: string;
  phone?: string;
  dateOfBirth?: string;
  address?: UserAddress;
}

export interface UserSearchParams {
  query?: string;
  role?: string;
  isActive?: boolean;
  isVerified?: boolean;
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface FollowRequest {
  userId: string;
}

export class UserService extends BaseApiClient {
  constructor(axios: any) {
    super(axios, '/users');
  }

  /**
   * Get user by ID
   */
  async getUser(userId: string): Promise<ApiResponse<User>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.get<User>(`/${encodedId}`);
  }

  /**
   * Get current user profile
   */
  async getCurrentUser(): Promise<ApiResponse<User>> {
    return this.get<User>('/me');
  }

  /**
   * Update current user profile
   */
  async updateCurrentUser(data: UpdateUserRequest): Promise<ApiResponse<User>> {
    // Validate phone if provided
    if (data.phone) {
      const phoneValidation = Validators.phoneNumber(data.phone);
      if (!phoneValidation.isValid) {
        throw new Error(phoneValidation.errors[0]);
      }
    }

    return this.put<User>('/me', data);
  }

  /**
   * Update user by ID (admin only)
   */
  async updateUser(userId: string, data: UpdateUserRequest): Promise<ApiResponse<User>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.put<User>(`/${encodedId}`, data);
  }

  /**
   * Delete user account
   */
  async deleteAccount(password: string): Promise<ApiResponse<void>> {
    return this.delete<void>('/me', {
      data: { password }
    });
  }

  /**
   * Delete user by ID (admin only)
   */
  async deleteUser(userId: string): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.delete<void>(`/${encodedId}`);
  }

  /**
   * Search users
   */
  async searchUsers(params: UserSearchParams): Promise<ApiResponse<PaginatedResponse<User>>> {
    return this.post<PaginatedResponse<User>>('/search', params);
  }

  /**
   * Get user suggestions for autocomplete
   */
  async getUserSuggestions(query: string, limit: number = 10): Promise<ApiResponse<User[]>> {
    return this.get<User[]>('/suggestions', { query, limit });
  }

  /**
   * Update user preferences
   */
  async updatePreferences(preferences: UserPreferences): Promise<ApiResponse<UserPreferences>> {
    return this.put<UserPreferences>('/me/preferences', preferences);
  }

  /**
   * Upload user avatar
   */
  async uploadAvatar(file: File): Promise<ApiResponse<{ url: string }>> {
    const formData = this.createFormData({ avatar: file });
    return this.upload<{ url: string }>('/me/avatar', formData);
  }

  /**
   * Delete user avatar
   */
  async deleteAvatar(): Promise<ApiResponse<void>> {
    return this.delete<void>('/me/avatar');
  }

  /**
   * Follow user
   */
  async followUser(userId: string): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.post<void>(`/${encodedId}/follow`);
  }

  /**
   * Unfollow user
   */
  async unfollowUser(userId: string): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.delete<void>(`/${encodedId}/follow`);
  }

  /**
   * Get user followers
   */
  async getFollowers(userId: string, params?: { page?: number; pageSize?: number }): Promise<ApiResponse<PaginatedResponse<User>>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.get<PaginatedResponse<User>>(`/${encodedId}/followers`, params);
  }

  /**
   * Get users that user follows
   */
  async getFollowing(userId: string, params?: { page?: number; pageSize?: number }): Promise<ApiResponse<PaginatedResponse<User>>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.get<PaginatedResponse<User>>(`/${encodedId}/following`, params);
  }

  /**
   * Block user
   */
  async blockUser(userId: string): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.post<void>(`/${encodedId}/block`);
  }

  /**
   * Unblock user
   */
  async unblockUser(userId: string): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.delete<void>(`/${encodedId}/block`);
  }

  /**
   * Get blocked users
   */
  async getBlockedUsers(params?: { page?: number; pageSize?: number }): Promise<ApiResponse<PaginatedResponse<User>>> {
    return this.get<PaginatedResponse<User>>('/me/blocked', params);
  }

  /**
   * Report user
   */
  async reportUser(userId: string, reason: string, details?: string): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.post<void>(`/${encodedId}/report`, { reason, details });
  }

  /**
   * Get user activity/stats
   */
  async getUserStats(userId: string): Promise<ApiResponse<UserStats>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.get<UserStats>(`/${encodedId}/stats`);
  }

  /**
   * Verify user (admin only)
   */
  async verifyUser(userId: string): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.post<void>(`/${encodedId}/verify`);
  }

  /**
   * Suspend user (admin only)
   */
  async suspendUser(userId: string, reason: string, duration?: number): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.post<void>(`/${encodedId}/suspend`, { reason, duration });
  }

  /**
   * Reactivate user (admin only)
   */
  async reactivateUser(userId: string): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(userId);
    return this.post<void>(`/${encodedId}/reactivate`);
  }
}