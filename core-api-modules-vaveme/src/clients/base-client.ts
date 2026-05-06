import { AxiosInstance, AxiosRequestConfig } from 'axios';
import { ApiError, ApiErrorHandler } from '../core/error-handler';
import { Mappers, ApiResponse } from '../utils/mappers';
import { Validators } from '../utils/validators';
import { Encoders } from '../utils/encoders';

export interface RequestOptions extends AxiosRequestConfig {
  skipValidation?: boolean;
  mapResponse?: boolean;
}

export abstract class BaseApiClient {
  protected axios: AxiosInstance;
  protected basePath: string;

  constructor(axios: AxiosInstance, basePath: string = '') {
    this.axios = axios;
    this.basePath = basePath;
  }

  /**
   * Make a GET request
   */
  protected async get<T>(
    path: string,
    params?: Record<string, any>,
    options: RequestOptions = {}
  ): Promise<ApiResponse<T>> {
    try {
      const cleanParams = params ? Validators.cleanParams(params) : undefined;
      const url = this.buildUrl(path, cleanParams);

      const response = await this.axios.get<T>(url, {
        ...options,
        params: undefined, // Params are in URL
      });

      return this.handleResponse<T>(response, options);
    } catch (error) {
      throw this.handleError(error, path, 'GET');
    }
  }

  /**
   * Make a POST request
   */
  protected async post<T>(
    path: string,
    data?: any,
    options: RequestOptions = {}
  ): Promise<ApiResponse<T>> {
    try {
      const url = this.buildUrl(path);
      const cleanData = data && !options.skipValidation ? Validators.cleanParams(data) : data;

      const response = await this.axios.post<T>(url, cleanData, options);
      return this.handleResponse<T>(response, options);
    } catch (error) {
      throw this.handleError(error, path, 'POST');
    }
  }

  /**
   * Make a PUT request
   */
  protected async put<T>(
    path: string,
    data?: any,
    options: RequestOptions = {}
  ): Promise<ApiResponse<T>> {
    try {
      const url = this.buildUrl(path);
      const cleanData = data && !options.skipValidation ? Validators.cleanParams(data) : data;

      const response = await this.axios.put<T>(url, cleanData, options);
      return this.handleResponse<T>(response, options);
    } catch (error) {
      throw this.handleError(error, path, 'PUT');
    }
  }

  /**
   * Make a PATCH request
   */
  protected async patch<T>(
    path: string,
    data?: any,
    options: RequestOptions = {}
  ): Promise<ApiResponse<T>> {
    try {
      const url = this.buildUrl(path);
      const cleanData = data && !options.skipValidation ? Validators.cleanParams(data) : data;

      const response = await this.axios.patch<T>(url, cleanData, options);
      return this.handleResponse<T>(response, options);
    } catch (error) {
      throw this.handleError(error, path, 'PATCH');
    }
  }

  /**
   * Make a DELETE request
   */
  protected async delete<T>(
    path: string,
    options: RequestOptions = {}
  ): Promise<ApiResponse<T>> {
    try {
      const url = this.buildUrl(path);
      const response = await this.axios.delete<T>(url, options);
      return this.handleResponse<T>(response, options);
    } catch (error) {
      throw this.handleError(error, path, 'DELETE');
    }
  }

  /**
   * Upload file(s)
   */
  protected async upload<T>(
    path: string,
    formData: FormData,
    options: RequestOptions = {}
  ): Promise<ApiResponse<T>> {
    try {
      const url = this.buildUrl(path);
      const response = await this.axios.post<T>(url, formData, {
        ...options,
        headers: {
          ...options.headers,
          'Content-Type': 'multipart/form-data',
        },
      });
      return this.handleResponse<T>(response, options);
    } catch (error) {
      throw this.handleError(error, path, 'UPLOAD');
    }
  }

  /**
   * Build full URL with base path
   */
  protected buildUrl(path: string, params?: Record<string, any>): string {
    const fullPath = this.basePath ? `${this.basePath}${path}` : path;
    return params ? Encoders.buildUrl(fullPath, params) : fullPath;
  }

  /**
   * Build path with parameters
   */
  protected buildPath(template: string, params: Record<string, any>): string {
    return Encoders.buildPath(template, params);
  }

  /**
   * Handle successful response
   */
  protected handleResponse<T>(response: any, options: RequestOptions): ApiResponse<T> {
    if (options.mapResponse !== false) {
      return Mappers.mapApiResponse<T>(response);
    }
    return response.data;
  }

  /**
   * Handle errors
   */
  protected handleError(error: unknown, endpoint: string, operation: string): ApiError {
    return ApiErrorHandler.handle(error, endpoint, operation);
  }

  /**
   * Create FormData from object
   */
  protected createFormData(data: Record<string, any>): FormData {
    const formData = new FormData();

    Object.entries(data).forEach(([key, value]) => {
      if (value === null || value === undefined) return;

      if (value instanceof File || value instanceof Blob) {
        formData.append(key, value);
      } else if (Array.isArray(value)) {
        value.forEach((item, index) => {
          if (item instanceof File || item instanceof Blob) {
            formData.append(`${key}[${index}]`, item);
          } else {
            formData.append(`${key}[${index}]`, String(item));
          }
        });
      } else if (typeof value === 'object') {
        formData.append(key, JSON.stringify(value));
      } else {
        formData.append(key, String(value));
      }
    });

    return formData;
  }
}