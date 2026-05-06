import axios, { AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios';
import { ApiConfig, defaultConfig, publicConfig, ssrConfig } from './config';
import { TokenManager } from './token-manager';
import { ApiErrorHandler } from './error-handler';

export interface CreateAxiosInstanceOptions {
  config?: Partial<ApiConfig>;
  includeAuth?: boolean;
  isSSR?: boolean;
}

export class AxiosClient {
  private static instances: Map<string, AxiosInstance> = new Map();

  static create(options: CreateAxiosInstanceOptions = {}): AxiosInstance {
    const {
      config: customConfig = {},
      includeAuth = true,
      isSSR = false,
    } = options;

    // Select base config
    let baseConfig: ApiConfig;
    if (isSSR) {
      baseConfig = { ...ssrConfig, ...customConfig };
    } else if (!includeAuth) {
      baseConfig = { ...publicConfig, ...customConfig };
    } else {
      baseConfig = { ...defaultConfig, ...customConfig };
    }

    // Create instance
    const instance = axios.create({
      baseURL: baseConfig.baseUrl,
      timeout: baseConfig.timeout,
      headers: {
        'Content-Type': 'application/json',
        ...(isSSR && { 'User-Agent': 'NextJS-SSR-Bot/1.0' }),
      },
    });

    // Add request interceptor
    instance.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        // Add auth token if needed
        if (includeAuth && !isSSR) {
          const token = TokenManager.getAccessToken();
          if (token && !TokenManager.isTokenExpired(token)) {
            config.headers.Authorization = `Bearer ${token}`;
          }
        }

        // Logging in development
        if (baseConfig.enableLogging) {
          console.log(`[API] ${config.method?.toUpperCase()} ${config.url}`, {
            params: config.params,
            data: config.data,
          });
        }

        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Add response interceptor
    instance.interceptors.response.use(
      (response) => {
        if (baseConfig.enableLogging) {
          console.log(`[API] Response from ${response.config.url}:`, response.data);
        }
        return response;
      },
      async (error) => {
        const originalRequest = error.config;

        // Handle 401 errors with token refresh
        if (
          includeAuth &&
          !isSSR &&
          error.response?.status === 401 &&
          !originalRequest._retry
        ) {
          originalRequest._retry = true;

          try {
            const refreshToken = TokenManager.getRefreshToken();
            if (refreshToken && !TokenManager.isTokenExpired(refreshToken)) {
              // Attempt to refresh token
              const refreshResponse = await this.refreshAccessToken(refreshToken);
              if (refreshResponse.success) {
                TokenManager.setTokens(refreshResponse.accessToken, refreshResponse.refreshToken);
                originalRequest.headers.Authorization = `Bearer ${refreshResponse.accessToken}`;
                return instance(originalRequest);
              }
            }
          } catch (refreshError) {
            // Refresh failed, clear tokens
            TokenManager.clearTokens();
          }
        }

        // Retry logic for transient errors
        if (
          ApiErrorHandler.isRetryable(error) &&
          (!originalRequest._retryCount || originalRequest._retryCount < baseConfig.retryAttempts)
        ) {
          originalRequest._retryCount = (originalRequest._retryCount || 0) + 1;
          
          const delay = ApiErrorHandler.getRetryDelay(
            originalRequest._retryCount,
            baseConfig.retryDelay
          );

          if (baseConfig.enableLogging) {
            console.log(
              `[API] Retrying request (attempt ${originalRequest._retryCount}/${baseConfig.retryAttempts}) after ${delay}ms`
            );
          }

          await new Promise(resolve => setTimeout(resolve, delay));
          return instance(originalRequest);
        }

        // Log error in development
        if (baseConfig.enableLogging) {
          console.error('[API] Request failed:', {
            url: error.config?.url,
            status: error.response?.status,
            data: error.response?.data,
          });
        }

        return Promise.reject(error);
      }
    );

    return instance;
  }

  static getDefault(): AxiosInstance {
    const key = 'default';
    if (!this.instances.has(key)) {
      this.instances.set(key, this.create());
    }
    return this.instances.get(key)!;
  }

  static getPublic(): AxiosInstance {
    const key = 'public';
    if (!this.instances.has(key)) {
      this.instances.set(key, this.create({ includeAuth: false }));
    }
    return this.instances.get(key)!;
  }

  static getSSR(): AxiosInstance {
    const key = 'ssr';
    if (!this.instances.has(key)) {
      this.instances.set(key, this.create({ isSSR: true }));
    }
    return this.instances.get(key)!;
  }

  private static async refreshAccessToken(refreshToken: string): Promise<any> {
    try {
      const response = await this.getPublic().post('/auth/refresh', {
        refreshToken,
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  }

  static clearCache(): void {
    this.instances.clear();
  }
}