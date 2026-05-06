export interface ApiConfig {
  baseUrl: string;
  timeout: number;
  isDevelopment: boolean;
  enableLogging: boolean;
  retryAttempts: number;
  retryDelay: number;
}

export const defaultConfig: ApiConfig = {
  baseUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3001',
  timeout: process.env.NODE_ENV === 'development' ? 600000 : 8000,
  isDevelopment: process.env.NODE_ENV === 'development',
  enableLogging: process.env.NODE_ENV === 'development',
  retryAttempts: 3,
  retryDelay: 1000,
};

export const ssrConfig: ApiConfig = {
  ...defaultConfig,
  timeout: 5000, // Aggressive timeout for SSR
  enableLogging: false,
  retryAttempts: 1,
  retryDelay: 0,
};

export const publicConfig: ApiConfig = {
  ...defaultConfig,
  timeout: 10000,
  retryAttempts: 2,
  retryDelay: 500,
};