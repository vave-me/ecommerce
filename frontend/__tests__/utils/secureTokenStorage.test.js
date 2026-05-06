import { jest } from '@jest/globals';
import { secureTokenStorage } from '../../src/utils/secureTokenStorage';

// Mock crypto for consistent testing
const mockCrypto = {
    getRandomValues: (arr) => {
        for (let i = 0; i < arr.length; i++) {
            arr[i] = Math.floor(Math.random() * 256);
        }
        return arr;
    },
    subtle: {
        encrypt: jest.fn().mockResolvedValue(new ArrayBuffer(16)),
        decrypt: jest.fn().mockResolvedValue(new ArrayBuffer(16))
    }
};

// Mock TextEncoder/TextDecoder
global.TextEncoder = jest.fn().mockImplementation(() => ({
    encode: jest.fn().mockReturnValue(new Uint8Array([1, 2, 3]))
}));

global.TextDecoder = jest.fn().mockImplementation(() => ({
    decode: jest.fn().mockReturnValue('test-token')
}));

describe('secureTokenStorage', () => {
    beforeEach(() => {
        // Clear storage and mocks
        localStorage.clear();
        sessionStorage.clear();
        jest.clearAllMocks();
        
        // Reset internal state
        secureTokenStorage.clearAll();
        
        // Mock window.crypto
        Object.defineProperty(window, 'crypto', {
            value: mockCrypto,
            writable: true
        });
    });

    describe('setAccessToken', () => {
        it('should store access token in memory', () => {
            const token = 'test-access-token';
            secureTokenStorage.setAccessToken(token);
            
            expect(secureTokenStorage.getAccessToken()).toBe(token);
        });

        it('should set token expiry time', () => {
            const token = 'test-access-token';
            const expiresIn = 3600; // 1 hour
            
            const now = Date.now();
            secureTokenStorage.setAccessToken(token, expiresIn);
            
            // Token should not be expired immediately
            expect(secureTokenStorage.isTokenExpired()).toBe(false);
        });

        it('should fallback to encrypted sessionStorage', () => {
            const token = 'test-access-token';
            secureTokenStorage.setAccessToken(token);
            
            // Verify encryption was attempted
            expect(sessionStorage.setItem).toHaveBeenCalled();
        });
    });

    describe('setRefreshToken', () => {
        it('should store refresh token in memory', () => {
            const token = 'test-refresh-token';
            secureTokenStorage.setRefreshToken(token);
            
            expect(secureTokenStorage.getRefreshToken()).toBe(token);
        });

        it('should encrypt refresh token for storage', () => {
            const token = 'test-refresh-token';
            secureTokenStorage.setRefreshToken(token);
            
            expect(sessionStorage.setItem).toHaveBeenCalled();
        });
    });

    describe('getAccessToken', () => {
        it('should return null if no token is set', () => {
            expect(secureTokenStorage.getAccessToken()).toBeNull();
        });

        it('should return null if token is expired', () => {
            const token = 'test-access-token';
            secureTokenStorage.setAccessToken(token, -1); // Already expired
            
            expect(secureTokenStorage.getAccessToken()).toBeNull();
        });

        it('should return valid token if not expired', () => {
            const token = 'test-access-token';
            secureTokenStorage.setAccessToken(token, 3600);
            
            expect(secureTokenStorage.getAccessToken()).toBe(token);
        });
    });

    describe('isTokenExpired', () => {
        it('should return true if no token exists', () => {
            expect(secureTokenStorage.isTokenExpired()).toBe(true);
        });

        it('should return true if token is expired', () => {
            secureTokenStorage.setAccessToken('token', -1);
            expect(secureTokenStorage.isTokenExpired()).toBe(true);
        });

        it('should return false if token is valid', () => {
            secureTokenStorage.setAccessToken('token', 3600);
            expect(secureTokenStorage.isTokenExpired()).toBe(false);
        });
    });

    describe('clearAll', () => {
        it('should clear all tokens from memory', () => {
            secureTokenStorage.setAccessToken('access-token');
            secureTokenStorage.setRefreshToken('refresh-token');
            
            secureTokenStorage.clearAll();
            
            expect(secureTokenStorage.getAccessToken()).toBeNull();
            expect(secureTokenStorage.getRefreshToken()).toBeNull();
        });

        it('should clear encrypted storage', () => {
            secureTokenStorage.setAccessToken('access-token');
            secureTokenStorage.clearAll();
            
            expect(sessionStorage.removeItem).toHaveBeenCalledWith('_sat');
            expect(sessionStorage.removeItem).toHaveBeenCalledWith('_srt');
        });
    });

    describe('rotateTokens', () => {
        it('should update both access and refresh tokens', () => {
            const newAccess = 'new-access-token';
            const newRefresh = 'new-refresh-token';
            
            secureTokenStorage.rotateTokens(newAccess, newRefresh);
            
            expect(secureTokenStorage.getAccessToken()).toBe(newAccess);
            expect(secureTokenStorage.getRefreshToken()).toBe(newRefresh);
        });

        it('should clear old tokens before setting new ones', () => {
            secureTokenStorage.setAccessToken('old-access');
            secureTokenStorage.setRefreshToken('old-refresh');
            
            const clearSpy = jest.spyOn(secureTokenStorage, 'clearAll');
            secureTokenStorage.rotateTokens('new-access', 'new-refresh');
            
            expect(clearSpy).toHaveBeenCalled();
        });
    });

    describe('encryption fallback', () => {
        it('should handle encryption errors gracefully', async () => {
            mockCrypto.subtle.encrypt.mockRejectedValueOnce(new Error('Encryption failed'));
            
            const token = 'test-token';
            // Should not throw
            expect(() => {
                secureTokenStorage.setAccessToken(token);
            }).not.toThrow();
            
            // Should still store in memory
            expect(secureTokenStorage.getAccessToken()).toBe(token);
        });

        it('should handle decryption errors gracefully', async () => {
            mockCrypto.subtle.decrypt.mockRejectedValueOnce(new Error('Decryption failed'));
            
            // Set encrypted value in storage
            sessionStorage.setItem('_sat', 'encrypted-data');
            
            // Should return null instead of throwing
            expect(secureTokenStorage.getAccessToken()).toBeNull();
        });
    });

    describe('security features', () => {
        it('should not expose tokens in console or debugging', () => {
            const token = 'sensitive-token';
            secureTokenStorage.setAccessToken(token);
            
            // Converting to string should not reveal token
            const stringified = JSON.stringify(secureTokenStorage);
            expect(stringified).not.toContain(token);
        });

        it('should use secure random values for encryption', () => {
            const spy = jest.spyOn(mockCrypto, 'getRandomValues');
            
            secureTokenStorage.setAccessToken('token');
            
            expect(spy).toHaveBeenCalled();
        });
    });
});