import { BaseApiClient } from '../../clients/base-client';
import { ApiResponse, PaginatedResponse } from '../../utils/mappers';
import { Encoders } from '../../utils/encoders';

export interface SearchFilters {
  query?: string;
  category?: string;
  subcategory?: string;
  priceMin?: number;
  priceMax?: number;
  condition?: string;
  location?: {
    latitude?: number;
    longitude?: number;
    radius?: number;
    city?: string;
    country?: string;
  };
  seller?: string;
  tags?: string[];
  brand?: string;
  model?: string;
  year?: number;
  color?: string;
  size?: string;
  material?: string;
  features?: string[];
  availability?: 'in_stock' | 'out_of_stock' | 'all';
  shipping?: 'free' | 'paid' | 'all';
  rating?: number;
  isVerifiedSeller?: boolean;
  hasImages?: boolean;
  hasVideo?: boolean;
  postedAfter?: string;
  postedBefore?: string;
  sortBy?: 'relevance' | 'price_asc' | 'price_desc' | 'date_asc' | 'date_desc' | 'rating' | 'distance';
  page?: number;
  pageSize?: number;
}

export interface SearchResult {
  id: string;
  title: string;
  description?: string;
  price: number;
  currency?: string;
  images?: string[];
  thumbnail?: string;
  category?: string;
  subcategory?: string;
  condition?: string;
  location?: {
    city?: string;
    country?: string;
    distance?: number;
  };
  seller?: {
    id: string;
    name: string;
    avatar?: string;
    rating?: number;
    isVerified?: boolean;
  };
  stats?: {
    views?: number;
    likes?: number;
    shares?: number;
  };
  createdAt?: string;
  updatedAt?: string;
  score?: number; // Search relevance score
}

export interface SearchSuggestion {
  text: string;
  type: 'query' | 'category' | 'brand' | 'seller';
  count?: number;
}

export interface SearchFacets {
  categories?: Array<{ name: string; count: number }>;
  brands?: Array<{ name: string; count: number }>;
  priceRanges?: Array<{ min: number; max: number; count: number }>;
  conditions?: Array<{ name: string; count: number }>;
  locations?: Array<{ name: string; count: number }>;
  features?: Array<{ name: string; count: number }>;
}

export interface SearchResponse<T> extends PaginatedResponse<T> {
  facets?: SearchFacets;
  suggestions?: SearchSuggestion[];
  executionTime?: number;
}

export interface SavedSearch {
  id: string;
  name?: string;
  filters: SearchFilters;
  alertEnabled?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface SearchHistory {
  id: string;
  query: string;
  filters?: SearchFilters;
  resultsCount?: number;
  timestamp: string;
}

export class SearchService extends BaseApiClient {
  constructor(axios: any) {
    super(axios, '/search');
  }

  /**
   * Perform advanced search with filters
   */
  async search(filters: SearchFilters): Promise<ApiResponse<SearchResponse<SearchResult>>> {
    return this.post<SearchResponse<SearchResult>>('/', filters);
  }

  /**
   * Quick search with just a query
   */
  async quickSearch(query: string, limit: number = 20): Promise<ApiResponse<SearchResult[]>> {
    return this.get<SearchResult[]>('/quick', { q: query, limit });
  }

  /**
   * Get search suggestions/autocomplete
   */
  async getSuggestions(query: string, limit: number = 10): Promise<ApiResponse<SearchSuggestion[]>> {
    return this.get<SearchSuggestion[]>('/suggestions', { q: query, limit });
  }

  /**
   * Get trending searches
   */
  async getTrending(limit: number = 10): Promise<ApiResponse<string[]>> {
    return this.get<string[]>('/trending', { limit });
  }

  /**
   * Get search facets for filters
   */
  async getFacets(filters?: Partial<SearchFilters>): Promise<ApiResponse<SearchFacets>> {
    return this.post<SearchFacets>('/facets', filters || {});
  }

  /**
   * Search within a specific category
   */
  async searchByCategory(
    category: string,
    filters: Omit<SearchFilters, 'category'>
  ): Promise<ApiResponse<SearchResponse<SearchResult>>> {
    return this.post<SearchResponse<SearchResult>>('/category', {
      ...filters,
      category
    });
  }

  /**
   * Search by location
   */
  async searchByLocation(
    latitude: number,
    longitude: number,
    radius: number,
    filters?: Omit<SearchFilters, 'location'>
  ): Promise<ApiResponse<SearchResponse<SearchResult>>> {
    return this.post<SearchResponse<SearchResult>>('/location', {
      ...filters,
      location: { latitude, longitude, radius }
    });
  }

  /**
   * Search similar items
   */
  async searchSimilar(itemId: string, limit: number = 10): Promise<ApiResponse<SearchResult[]>> {
    const encodedId = Encoders.encodePathParam(itemId);
    return this.get<SearchResult[]>(`/similar/${encodedId}`, { limit });
  }

  /**
   * Save a search
   */
  async saveSearch(search: Omit<SavedSearch, 'id' | 'createdAt' | 'updatedAt'>): Promise<ApiResponse<SavedSearch>> {
    return this.post<SavedSearch>('/saved', search);
  }

  /**
   * Get saved searches
   */
  async getSavedSearches(): Promise<ApiResponse<SavedSearch[]>> {
    return this.get<SavedSearch[]>('/saved');
  }

  /**
   * Update saved search
   */
  async updateSavedSearch(
    searchId: string,
    updates: Partial<Omit<SavedSearch, 'id' | 'createdAt' | 'updatedAt'>>
  ): Promise<ApiResponse<SavedSearch>> {
    const encodedId = Encoders.encodePathParam(searchId);
    return this.put<SavedSearch>(`/saved/${encodedId}`, updates);
  }

  /**
   * Delete saved search
   */
  async deleteSavedSearch(searchId: string): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(searchId);
    return this.delete<void>(`/saved/${encodedId}`);
  }

  /**
   * Get search history
   */
  async getSearchHistory(limit: number = 20): Promise<ApiResponse<SearchHistory[]>> {
    return this.get<SearchHistory[]>('/history', { limit });
  }

  /**
   * Clear search history
   */
  async clearSearchHistory(): Promise<ApiResponse<void>> {
    return this.delete<void>('/history');
  }

  /**
   * Delete specific search from history
   */
  async deleteFromHistory(historyId: string): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(historyId);
    return this.delete<void>(`/history/${encodedId}`);
  }

  /**
   * Enable/disable search alerts for saved search
   */
  async toggleSearchAlert(searchId: string, enabled: boolean): Promise<ApiResponse<void>> {
    const encodedId = Encoders.encodePathParam(searchId);
    return this.post<void>(`/saved/${encodedId}/alert`, { enabled });
  }

  /**
   * Get popular searches in category
   */
  async getPopularInCategory(category: string, limit: number = 10): Promise<ApiResponse<string[]>> {
    const encodedCategory = Encoders.encodePathParam(category);
    return this.get<string[]>(`/popular/${encodedCategory}`, { limit });
  }

  /**
   * Advanced text search with NLP
   */
  async semanticSearch(query: string, filters?: Omit<SearchFilters, 'query'>): Promise<ApiResponse<SearchResponse<SearchResult>>> {
    return this.post<SearchResponse<SearchResult>>('/semantic', {
      ...filters,
      query
    });
  }

  /**
   * Barcode/QR code search
   */
  async searchByCode(code: string, type: 'barcode' | 'qr' = 'barcode'): Promise<ApiResponse<SearchResult[]>> {
    return this.get<SearchResult[]>('/code', { code, type });
  }

  /**
   * Image-based search
   */
  async searchByImage(imageFile: File, filters?: Omit<SearchFilters, 'query'>): Promise<ApiResponse<SearchResponse<SearchResult>>> {
    const formData = this.createFormData({
      image: imageFile,
      filters: filters ? JSON.stringify(filters) : undefined
    });
    return this.upload<SearchResponse<SearchResult>>('/image', formData);
  }
}