// Core exports
export { AxiosClient, type CreateAxiosInstanceOptions } from './core/axios-client';
export { TokenManager, type TokenPayload } from './core/token-manager';
export { ApiErrorHandler, ErrorSeverity, type ApiError } from './core/error-handler';
export { defaultConfig, ssrConfig, publicConfig, type ApiConfig } from './core/config';

// Utils exports
export { Validators, type ValidationResult } from './utils/validators';
export { Encoders } from './utils/encoders';
export { Mappers, type ApiResponse, type PaginatedResponse } from './utils/mappers';

// Client exports
export { BaseApiClient, type RequestOptions } from './clients/base-client';

// Service exports
export { AuthService, type LoginRequest, type RegisterRequest, type AuthResponse } from './services/auth/auth.service';
export { UserService, type User, type UpdateUserRequest, type UserSearchParams } from './services/user/user.service';
export { SearchService, type SearchFilters, type SearchResult, type SearchResponse } from './services/search/search.service';

// Factory function to create API instances
export interface ApiServices {
  auth: AuthService;
  users: UserService;
  search: SearchService;
}

export function createApiServices(axiosInstance?: any): ApiServices {
  const axios = axiosInstance || AxiosClient.getDefault();
  
  return {
    auth: new AuthService(axios),
    users: new UserService(axios),
    search: new SearchService(axios),
  };
}

// Re-export for convenience
export const defaultApiServices = createApiServices();