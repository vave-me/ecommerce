export interface ValidationResult {
  isValid: boolean;
  errors: string[];
}

export class Validators {
  static email(email: string): ValidationResult {
    const errors: string[] = [];
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    if (!email) {
      errors.push('Email is required');
    } else if (!emailRegex.test(email)) {
      errors.push('Invalid email format');
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  static password(password: string): ValidationResult {
    const errors: string[] = [];

    if (!password) {
      errors.push('Password is required');
    } else {
      if (password.length < 8) {
        errors.push('Password must be at least 8 characters long');
      }
      if (!/[A-Z]/.test(password)) {
        errors.push('Password must contain at least one uppercase letter');
      }
      if (!/[a-z]/.test(password)) {
        errors.push('Password must contain at least one lowercase letter');
      }
      if (!/[0-9]/.test(password)) {
        errors.push('Password must contain at least one number');
      }
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  static phoneNumber(phone: string): ValidationResult {
    const errors: string[] = [];
    const phoneRegex = /^\+?[1-9]\d{1,14}$/;

    if (!phone) {
      errors.push('Phone number is required');
    } else if (!phoneRegex.test(phone.replace(/[\s-()]/g, ''))) {
      errors.push('Invalid phone number format');
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  static url(url: string): ValidationResult {
    const errors: string[] = [];

    try {
      new URL(url);
    } catch {
      errors.push('Invalid URL format');
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  static required<T>(value: T, fieldName: string): ValidationResult {
    const errors: string[] = [];

    if (value === null || value === undefined || value === '') {
      errors.push(`${fieldName} is required`);
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  static minLength(value: string, min: number, fieldName: string): ValidationResult {
    const errors: string[] = [];

    if (value && value.length < min) {
      errors.push(`${fieldName} must be at least ${min} characters long`);
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  static maxLength(value: string, max: number, fieldName: string): ValidationResult {
    const errors: string[] = [];

    if (value && value.length > max) {
      errors.push(`${fieldName} must be no more than ${max} characters long`);
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  static range(value: number, min: number, max: number, fieldName: string): ValidationResult {
    const errors: string[] = [];

    if (value < min || value > max) {
      errors.push(`${fieldName} must be between ${min} and ${max}`);
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  static cleanParams<T extends Record<string, any>>(params: T): Partial<T> {
    const cleaned: Partial<T> = {};

    Object.entries(params).forEach(([key, value]) => {
      if (value !== null && value !== undefined && value !== '') {
        cleaned[key as keyof T] = value;
      }
    });

    return cleaned;
  }

  static validateBatch(validations: ValidationResult[]): ValidationResult {
    const allErrors = validations.flatMap(v => v.errors);
    
    return {
      isValid: allErrors.length === 0,
      errors: allErrors,
    };
  }
}