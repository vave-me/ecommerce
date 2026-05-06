import {
    formatPrice,
    formatDate,
    formatDateTime,
    formatRelativeTime,
    formatFileSize,
    formatCompactNumber,
    truncateText
} from '../../src/utils/formatters';

describe('formatters', () => {
    describe('formatPrice', () => {
        it('should format price with default EUR currency', () => {
            expect(formatPrice(29.99)).toBe('29,99 €');
            expect(formatPrice(1000)).toBe('1.000,00 €');
            expect(formatPrice(0)).toBe('0,00 €');
        });

        it('should handle different currencies', () => {
            expect(formatPrice(29.99, 'USD', 'en-US')).toBe('$29.99');
            expect(formatPrice(29.99, 'GBP', 'en-GB')).toBe('£29.99');
        });

        it('should handle null/undefined/NaN values', () => {
            expect(formatPrice(null)).toBe('0,00 €');
            expect(formatPrice(undefined)).toBe('0,00 €');
            expect(formatPrice(NaN)).toBe('0,00 €');
            expect(formatPrice('invalid')).toBe('0,00 €');
        });

        it('should format large numbers correctly', () => {
            expect(formatPrice(1234567.89)).toBe('1.234.567,89 €');
        });

        it('should handle negative prices', () => {
            expect(formatPrice(-29.99)).toBe('-29,99 €');
        });
    });

    describe('formatDate', () => {
        it('should format date correctly', () => {
            const date = new Date('2025-08-04');
            const formatted = formatDate(date);
            expect(formatted).toMatch(/4\.\s*Aug\.\s*2025/);
        });

        it('should handle string dates', () => {
            const formatted = formatDate('2025-08-04');
            expect(formatted).toMatch(/4\.\s*Aug\.\s*2025/);
        });

        it('should handle null/undefined dates', () => {
            expect(formatDate(null)).toBe('');
            expect(formatDate(undefined)).toBe('');
            expect(formatDate('')).toBe('');
        });

        it('should handle invalid dates', () => {
            expect(formatDate('invalid-date')).toBe('');
            expect(formatDate(new Date('invalid'))).toBe('');
        });

        it('should accept custom options', () => {
            const date = new Date('2025-08-04');
            const formatted = formatDate(date, { 
                year: 'numeric', 
                month: 'long', 
                day: 'numeric' 
            });
            expect(formatted).toMatch(/4\.\s*August\s*2025/);
        });
    });

    describe('formatDateTime', () => {
        it('should format date with time', () => {
            const date = new Date('2025-08-04T14:30:00');
            const formatted = formatDateTime(date);
            expect(formatted).toMatch(/4\.\s*Aug\.\s*2025/);
            expect(formatted).toMatch(/14:30/);
        });

        it('should handle edge cases', () => {
            expect(formatDateTime(null)).toBe('');
            expect(formatDateTime(undefined)).toBe('');
        });
    });

    describe('formatRelativeTime', () => {
        beforeEach(() => {
            // Mock current time for consistent testing
            jest.useFakeTimers();
            jest.setSystemTime(new Date('2025-08-04T12:00:00'));
        });

        afterEach(() => {
            jest.useRealTimers();
        });

        it('should format "just now" for recent times', () => {
            const date = new Date('2025-08-04T11:59:30');
            expect(formatRelativeTime(date)).toBe('just now');
        });

        it('should format minutes ago', () => {
            const date = new Date('2025-08-04T11:55:00');
            expect(formatRelativeTime(date)).toBe('5 minutes ago');
            
            const date2 = new Date('2025-08-04T11:59:00');
            expect(formatRelativeTime(date2)).toBe('1 minute ago');
        });

        it('should format hours ago', () => {
            const date = new Date('2025-08-04T10:00:00');
            expect(formatRelativeTime(date)).toBe('2 hours ago');
            
            const date2 = new Date('2025-08-04T11:00:00');
            expect(formatRelativeTime(date2)).toBe('1 hour ago');
        });

        it('should format days ago', () => {
            const date = new Date('2025-08-02T12:00:00');
            expect(formatRelativeTime(date)).toBe('2 days ago');
            
            const date2 = new Date('2025-08-03T12:00:00');
            expect(formatRelativeTime(date2)).toBe('1 day ago');
        });

        it('should format full date for older dates', () => {
            const date = new Date('2025-07-20T12:00:00');
            const formatted = formatRelativeTime(date);
            expect(formatted).toMatch(/20\.\s*Juli\s*2025/);
        });

        it('should handle null/undefined', () => {
            expect(formatRelativeTime(null)).toBe('');
            expect(formatRelativeTime(undefined)).toBe('');
        });
    });

    describe('formatFileSize', () => {
        it('should format bytes correctly', () => {
            expect(formatFileSize(0)).toBe('0 Bytes');
            expect(formatFileSize(100)).toBe('100 Bytes');
            expect(formatFileSize(1023)).toBe('1023 Bytes');
        });

        it('should format KB correctly', () => {
            expect(formatFileSize(1024)).toBe('1 KB');
            expect(formatFileSize(1536)).toBe('1.5 KB');
            expect(formatFileSize(2048)).toBe('2 KB');
        });

        it('should format MB correctly', () => {
            expect(formatFileSize(1048576)).toBe('1 MB');
            expect(formatFileSize(1572864)).toBe('1.5 MB');
            expect(formatFileSize(5242880)).toBe('5 MB');
        });

        it('should format GB correctly', () => {
            expect(formatFileSize(1073741824)).toBe('1 GB');
            expect(formatFileSize(2147483648)).toBe('2 GB');
        });

        it('should round to 2 decimal places', () => {
            expect(formatFileSize(1234567)).toBe('1.18 MB');
            expect(formatFileSize(123456789)).toBe('117.74 MB');
        });
    });

    describe('formatCompactNumber', () => {
        it('should not abbreviate small numbers', () => {
            expect(formatCompactNumber(0)).toBe('0');
            expect(formatCompactNumber(999)).toBe('999');
        });

        it('should abbreviate thousands', () => {
            expect(formatCompactNumber(1000)).toMatch(/1K|1k/);
            expect(formatCompactNumber(1500)).toMatch(/1\.5K|1\.5k/);
            expect(formatCompactNumber(999999)).toMatch(/1M|1m/);
        });

        it('should abbreviate millions', () => {
            expect(formatCompactNumber(1000000)).toMatch(/1M|1m/);
            expect(formatCompactNumber(2500000)).toMatch(/2\.5M|2\.5m/);
        });

        it('should handle negative numbers', () => {
            expect(formatCompactNumber(-1500)).toMatch(/-1\.5K|-1\.5k/);
        });
    });

    describe('truncateText', () => {
        it('should not truncate short text', () => {
            const text = 'Short text';
            expect(truncateText(text)).toBe(text);
            expect(truncateText(text, 20)).toBe(text);
        });

        it('should truncate long text with ellipsis', () => {
            const text = 'This is a very long text that needs to be truncated because it exceeds the maximum length allowed for display purposes';
            expect(truncateText(text, 50)).toBe('This is a very long text that needs to be truncat...');
        });

        it('should use default length of 100', () => {
            const text = 'a'.repeat(150);
            const truncated = truncateText(text);
            expect(truncated.length).toBe(103); // 100 + '...'
            expect(truncated.endsWith('...')).toBe(true);
        });

        it('should handle null/undefined text', () => {
            expect(truncateText(null)).toBe(null);
            expect(truncateText(undefined)).toBe(undefined);
            expect(truncateText('')).toBe('');
        });

        it('should trim whitespace before adding ellipsis', () => {
            const text = 'This is a text with trailing spaces     that will be truncated';
            expect(truncateText(text, 40)).toBe('This is a text with trailing spaces...');
        });
    });
});