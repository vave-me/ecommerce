export class Encoders {
  /**
   * Safely encode a URI component, handling null/undefined values
   */
  static encodeURIComponent(component: any): string {
    if (component === null || component === undefined) {
      return '';
    }
    return encodeURIComponent(String(component));
  }

  /**
   * Encode an object's values for use in query parameters
   */
  static encodeQueryParams(params: Record<string, any>): string {
    const entries = Object.entries(params)
      .filter(([_, value]) => value !== null && value !== undefined && value !== '')
      .map(([key, value]) => {
        if (Array.isArray(value)) {
          return value
            .map(v => `${this.encodeURIComponent(key)}=${this.encodeURIComponent(v)}`)
            .join('&');
        }
        return `${this.encodeURIComponent(key)}=${this.encodeURIComponent(value)}`;
      });

    return entries.join('&');
  }

  /**
   * Build a complete URL with query parameters
   */
  static buildUrl(baseUrl: string, params?: Record<string, any>): string {
    if (!params || Object.keys(params).length === 0) {
      return baseUrl;
    }

    const queryString = this.encodeQueryParams(params);
    const separator = baseUrl.includes('?') ? '&' : '?';
    
    return `${baseUrl}${separator}${queryString}`;
  }

  /**
   * Encode a path parameter, ensuring proper escaping
   */
  static encodePathParam(param: any): string {
    const encoded = this.encodeURIComponent(param);
    // Additional encoding for special characters that might break URLs
    return encoded.replace(/[!'()*]/g, (c) => {
      return '%' + c.charCodeAt(0).toString(16).toUpperCase();
    });
  }

  /**
   * Build a path with parameters replaced
   */
  static buildPath(template: string, params: Record<string, any>): string {
    let path = template;

    Object.entries(params).forEach(([key, value]) => {
      const placeholder = `:${key}`;
      if (path.includes(placeholder)) {
        path = path.replace(placeholder, this.encodePathParam(value));
      }
    });

    return path;
  }

  /**
   * Sanitize filename for safe storage/display
   */
  static sanitizeFilename(filename: string): string {
    // Remove or replace unsafe characters
    return filename
      .replace(/[<>:"/\\|?*\x00-\x1F]/g, '_')
      .replace(/\.{2,}/g, '_')
      .trim();
  }

  /**
   * Generate a slug from text
   */
  static slugify(text: string): string {
    return text
      .toLowerCase()
      .trim()
      .replace(/[^\w\s-]/g, '')
      .replace(/[\s_-]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }
}