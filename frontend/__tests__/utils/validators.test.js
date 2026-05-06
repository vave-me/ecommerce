import {
    isEmpty,
    isValidEmail,
    isValidPhone,
    isValidUrl,
    isValidPrice,
    validateRequiredFields,
    isValidLength,
    sanitizeInput,
    isValidFileType,
    isValidFileSize
} from '../../src/utils/validators';

describe('validators', () => {
    describe('isEmpty', () => {
        it('should return true for null and undefined', () => {
            expect(isEmpty(null)).toBe(true);
            expect(isEmpty(undefined)).toBe(true);
        });

        it('should return true for empty string', () => {
            expect(isEmpty('')).toBe(true);
        });

        it('should return true for empty array', () => {
            expect(isEmpty([])).toBe(true);
        });

        it('should return true for empty object', () => {
            expect(isEmpty({})).toBe(true);
        });

        it('should return false for non-empty values', () => {
            expect(isEmpty('text')).toBe(false);
            expect(isEmpty(' ')).toBe(false);
            expect(isEmpty(0)).toBe(false);
            expect(isEmpty(false)).toBe(false);
            expect(isEmpty([1])).toBe(false);
            expect(isEmpty({ a: 1 })).toBe(false);
        });
    });

    describe('isValidEmail', () => {
        it('should validate correct email formats', () => {
            expect(isValidEmail('redacted-email@example.com')).toBe(true);
            expect(isValidEmail('redacted-email@example.com')).toBe(true);
            expect(isValidEmail('redacted-email@example.com')).toBe(true);
            expect(isValidEmail('redacted-email@example.com')).toBe(true);
        });

        it('should reject invalid email formats', () => {
            expect(isValidEmail('invalid')).toBe(false);
            expect(isValidEmail('@example.com')).toBe(false);
            expect(isValidEmail('user@')).toBe(false);
            expect(isValidEmail('user @example.com')).toBe(false);
            expect(isValidEmail('user@example')).toBe(false);
            expect(isValidEmail('')).toBe(false);
        });
    });

    describe('isValidPhone', () => {
        it('should validate correct phone formats', () => {
            expect(isValidPhone('+49123456789')).toBe(true);
            expect(isValidPhone('0123456789')).toBe(true);
            expect(isValidPhone('+1-555-123-4567')).toBe(true);
            expect(isValidPhone('(555) 123-4567')).toBe(true);
            expect(isValidPhone('555.123.4567')).toBe(true);
        });

        it('should reject invalid phone formats', () => {
            expect(isValidPhone('abc')).toBe(false);
            expect(isValidPhone('123')).toBe(false);
            expect(isValidPhone('')).toBe(false);
            expect(isValidPhone('++123456789')).toBe(false);
        });

        it('should handle phone with spaces and dashes', () => {
            expect(isValidPhone('555 123 4567')).toBe(true);
            expect(isValidPhone('555-123-4567')).toBe(true);
        });
    });

    describe('isValidUrl', () => {
        it('should validate correct URLs', () => {
            expect(isValidUrl('https://example.com')).toBe(true);
            expect(isValidUrl('http://example.com')).toBe(true);
            expect(isValidUrl('https://example.com/path')).toBe(true);
            expect(isValidUrl('https://example.com:8080')).toBe(true);
            expect(isValidUrl('https://example.com?query=value')).toBe(true);
            expect(isValidUrl('ftp://example.com')).toBe(true);
        });

        it('should reject invalid URLs', () => {
            expect(isValidUrl('not a url')).toBe(false);
            expect(isValidUrl('example.com')).toBe(false);
            expect(isValidUrl('//example.com')).toBe(false);
            expect(isValidUrl('')).toBe(false);
            expect(isValidUrl('javascript:alert(1)')).toBe(false);
        });
    });

    describe('isValidPrice', () => {
        it('should validate correct prices', () => {
            expect(isValidPrice(0)).toBe(true);
            expect(isValidPrice(29.99)).toBe(true);
            expect(isValidPrice('29.99')).toBe(true);
            expect(isValidPrice(100)).toBe(true);
            expect(isValidPrice('100.00')).toBe(true);
            expect(isValidPrice('0.01')).toBe(true);
        });

        it('should reject invalid prices', () => {
            expect(isValidPrice(-1)).toBe(false);
            expect(isValidPrice('abc')).toBe(false);
            expect(isValidPrice('29.999')).toBe(false); // More than 2 decimals
            expect(isValidPrice('')).toBe(false);
            expect(isValidPrice(null)).toBe(false);
            expect(isValidPrice(undefined)).toBe(false);
        });
    });

    describe('validateRequiredFields', () => {
        it('should validate all required fields are present', () => {
            const obj = {
                name: 'John',
                email: 'redacted-email@example.com',
                age: 25
            };
            const result = validateRequiredFields(obj, ['name', 'email']);
            
            expect(result.isValid).toBe(true);
            expect(result.missingFields).toEqual([]);
        });

        it('should identify missing fields', () => {
            const obj = {
                name: 'John',
                email: ''
            };
            const result = validateRequiredFields(obj, ['name', 'email', 'phone']);
            
            expect(result.isValid).toBe(false);
            expect(result.missingFields).toEqual(['email', 'phone']);
        });

        it('should handle null and undefined as missing', () => {
            const obj = {
                name: null,
                email: undefined,
                phone: 'valid'
            };
            const result = validateRequiredFields(obj, ['name', 'email', 'phone']);
            
            expect(result.isValid).toBe(false);
            expect(result.missingFields).toEqual(['name', 'email']);
        });

        it('should handle empty arrays and objects as missing', () => {
            const obj = {
                tags: [],
                metadata: {},
                name: 'Valid'
            };
            const result = validateRequiredFields(obj, ['tags', 'metadata', 'name']);
            
            expect(result.isValid).toBe(false);
            expect(result.missingFields).toEqual(['tags', 'metadata']);
        });
    });

    describe('isValidLength', () => {
        it('should validate string within length range', () => {
            expect(isValidLength('hello', 1, 10)).toBe(true);
            expect(isValidLength('hello', 5, 5)).toBe(true);
            expect(isValidLength('', 0, 10)).toBe(true);
        });

        it('should reject strings outside length range', () => {
            expect(isValidLength('hello', 10, 20)).toBe(false);
            expect(isValidLength('hello world', 1, 5)).toBe(false);
            expect(isValidLength('', 1, 10)).toBe(false);
        });

        it('should handle default parameters', () => {
            expect(isValidLength('any length string')).toBe(true);
            expect(isValidLength('hello', 3)).toBe(true);
        });

        it('should handle null/undefined', () => {
            expect(isValidLength(null, 0)).toBe(true);
            expect(isValidLength(null, 1)).toBe(false);
            expect(isValidLength(undefined, 0)).toBe(true);
        });
    });

    describe('sanitizeInput', () => {
        it('should escape HTML special characters', () => {
            expect(sanitizeInput('<script>alert("XSS")</script>'))
                .toBe('&lt;script&gt;alert(&quot;XSS&quot;)&lt;/script&gt;');
            
            expect(sanitizeInput('Hello & goodbye'))
                .toBe('Hello &amp; goodbye');
            
            expect(sanitizeInput("It's a test"))
                .toBe('It&#039;s a test');
        });

        it('should handle empty and null inputs', () => {
            expect(sanitizeInput('')).toBe('');
            expect(sanitizeInput(null)).toBe('');
            expect(sanitizeInput(undefined)).toBe('');
        });

        it('should not modify safe text', () => {
            const safeText = 'Hello World 123';
            expect(sanitizeInput(safeText)).toBe(safeText);
        });
    });

    describe('isValidFileType', () => {
        it('should validate specific file types', () => {
            const pngFile = new File([''], 'test.png', { type: 'image/png' });
            const jpgFile = new File([''], 'test.jpg', { type: 'image/jpeg' });
            const pdfFile = new File([''], 'test.pdf', { type: 'application/pdf' });
            
            expect(isValidFileType(pngFile, ['image/png'])).toBe(true);
            expect(isValidFileType(jpgFile, ['image/jpeg', 'image/png'])).toBe(true);
            expect(isValidFileType(pdfFile, ['application/pdf'])).toBe(true);
        });

        it('should handle wildcard types', () => {
            const pngFile = new File([''], 'test.png', { type: 'image/png' });
            const jpgFile = new File([''], 'test.jpg', { type: 'image/jpeg' });
            const pdfFile = new File([''], 'test.pdf', { type: 'application/pdf' });
            
            expect(isValidFileType(pngFile, ['image/*'])).toBe(true);
            expect(isValidFileType(jpgFile, ['image/*'])).toBe(true);
            expect(isValidFileType(pdfFile, ['image/*'])).toBe(false);
        });

        it('should reject invalid file types', () => {
            const exeFile = new File([''], 'test.exe', { type: 'application/x-msdownload' });
            
            expect(isValidFileType(exeFile, ['image/*', 'application/pdf'])).toBe(false);
        });

        it('should handle null/undefined files', () => {
            expect(isValidFileType(null, ['image/*'])).toBe(false);
            expect(isValidFileType(undefined, ['image/*'])).toBe(false);
        });
    });

    describe('isValidFileSize', () => {
        it('should validate file size within limit', () => {
            const smallFile = new File(['a'.repeat(1024 * 1024)], 'test.txt'); // 1MB
            const mediumFile = new File(['a'.repeat(5 * 1024 * 1024)], 'test.txt'); // 5MB
            
            expect(isValidFileSize(smallFile, 2)).toBe(true);
            expect(isValidFileSize(smallFile, 1)).toBe(true);
            expect(isValidFileSize(mediumFile, 10)).toBe(true);
        });

        it('should reject files exceeding size limit', () => {
            const largeFile = new File(['a'.repeat(10 * 1024 * 1024)], 'test.txt'); // 10MB
            
            expect(isValidFileSize(largeFile, 5)).toBe(false);
            expect(isValidFileSize(largeFile, 9)).toBe(false);
        });

        it('should handle edge cases', () => {
            const emptyFile = new File([''], 'empty.txt');
            
            expect(isValidFileSize(emptyFile, 1)).toBe(true);
            expect(isValidFileSize(null, 1)).toBe(false);
            expect(isValidFileSize(undefined, 1)).toBe(false);
        });
    });
});