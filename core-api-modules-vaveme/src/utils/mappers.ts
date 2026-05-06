export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
  metadata?: {
    page?: number;
    pageSize?: number;
    total?: number;
    hasMore?: boolean;
  };
}

export interface PaginatedResponse<T> {
  items: T[];
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
  hasNextPage: boolean;
  hasPreviousPage: boolean;
}

export class Mappers {
  /**
   * Map API response to standardized format
   */
  static mapApiResponse<T>(response: any): ApiResponse<T> {
    // Handle various response formats
    if (response.data?.success !== undefined) {
      return response.data;
    }

    // Standard axios response
    if (response.data !== undefined) {
      return {
        success: true,
        data: response.data,
      };
    }

    // Direct response
    return {
      success: true,
      data: response,
    };
  }

  /**
   * Map paginated response to standardized format
   */
  static mapPaginatedResponse<T>(
    response: any,
    itemMapper?: (item: any) => T
  ): PaginatedResponse<T> {
    const data = response.data || response;
    const items = data.items || data.results || data.data || [];
    const mappedItems = itemMapper ? items.map(itemMapper) : items;

    const page = data.page || data.currentPage || 1;
    const pageSize = data.pageSize || data.perPage || items.length;
    const total = data.total || data.totalCount || items.length;
    const totalPages = data.totalPages || Math.ceil(total / pageSize);

    return {
      items: mappedItems,
      page,
      pageSize,
      total,
      totalPages,
      hasNextPage: page < totalPages,
      hasPreviousPage: page > 1,
    };
  }

  /**
   * Extract ID from various response formats
   */
  static extractId(response: any): string | number | null {
    const data = response.data || response;
    
    return data.id || data._id || data.ID || data.uuid || data.identifier || null;
  }

  /**
   * Map date fields to Date objects
   */
  static mapDates<T extends Record<string, any>>(
    obj: T,
    dateFields: (keyof T)[]
  ): T {
    const mapped = { ...obj };

    dateFields.forEach(field => {
      if (mapped[field]) {
        mapped[field] = new Date(mapped[field] as any) as any;
      }
    });

    return mapped;
  }

  /**
   * Map snake_case to camelCase
   */
  static snakeToCamel<T = any>(obj: any): T {
    if (obj === null || typeof obj !== 'object') return obj;
    
    if (Array.isArray(obj)) {
      return obj.map(item => this.snakeToCamel(item)) as any;
    }

    const mapped: any = {};
    
    Object.entries(obj).forEach(([key, value]) => {
      const camelKey = key.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
      mapped[camelKey] = this.snakeToCamel(value);
    });

    return mapped;
  }

  /**
   * Map camelCase to snake_case
   */
  static camelToSnake<T = any>(obj: any): T {
    if (obj === null || typeof obj !== 'object') return obj;
    
    if (Array.isArray(obj)) {
      return obj.map(item => this.camelToSnake(item)) as any;
    }

    const mapped: any = {};
    
    Object.entries(obj).forEach(([key, value]) => {
      const snakeKey = key.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`);
      mapped[snakeKey] = this.camelToSnake(value);
    });

    return mapped;
  }

  /**
   * Create display name from various fields
   */
  static getDisplayName(obj: any): string {
    return obj.displayName || 
           obj.name || 
           obj.title || 
           obj.label || 
           obj.username ||
           obj.email?.split('@')[0] ||
           'Unknown';
  }

  /**
   * Safe JSON parse with fallback
   */
  static safeJsonParse<T = any>(json: string, fallback: T): T {
    try {
      return JSON.parse(json);
    } catch {
      return fallback;
    }
  }

  /**
   * Extract error message from various error formats
   */
  static extractErrorMessage(error: any): string {
    if (typeof error === 'string') return error;
    
    return error.userMessage ||
           error.message ||
           error.error ||
           error.detail ||
           error.data?.message ||
           error.response?.data?.message ||
           'An unknown error occurred';
  }
}