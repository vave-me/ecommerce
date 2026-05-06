import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { AuthProvider, useAuth } from '../AuthContext';
import { jwtDecode } from 'jwt-decode';

// Mock jwt-decode
jest.mock('jwt-decode');

// Mock the API modules
jest.mock('../../api/client/userApi', () => ({
    getAccessToken: jest.fn(),
    setAccessToken: jest.fn(),
    setRefreshToken: jest.fn(),
    clearTokens: jest.fn(),
    refreshAccessToken: jest.fn(),
    initFromLocalStorage: jest.fn(),
    initializeAuth: jest.fn().mockReturnValue(true),
    logoutUser: jest.fn(),
    clearUserTokens: jest.fn(),
    getRefreshToken: jest.fn()
}));

jest.mock('../../api/userApi', () => ({
    loginUser: jest.fn(),
    registerUser: jest.fn(),
    loginWithGoogle: jest.fn()
}));

// Mock router
jest.mock('next/navigation', () => ({
    useRouter: () => ({
        push: jest.fn(),
        refresh: jest.fn()
    })
}));

// Mock React Query
jest.mock('@tanstack/react-query', () => ({
    useQueryClient: () => ({
        invalidateQueries: jest.fn(),
        removeQueries: jest.fn()
    })
}));

// Test component to access auth context
const TestComponent = () => {
    const { user } = useAuth();
    return (
        <div>
            {user && (
                <>
                    <div data-testid="user-id">{user.userId}</div>
                    <div data-testid="user-email">{user.email}</div>
                    <div data-testid="user-role">{user.role}</div>
                </>
            )}
        </div>
    );
};

describe('AuthContext - Role Extraction', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        localStorage.clear();
    });

    it('should extract role from JWT token', async () => {
        const mockToken = 'mock.jwt.token';
        const mockDecodedToken = {
            userId: 'user123',
            email: 'redacted-email@example.com',
            userName: 'testuser',
            role: 'admin',
            lat: 40.7128,
            lng: -74.0060,
            exp: Math.floor(Date.now() / 1000) + 3600 // 1 hour from now
        };

        // Mock the userApi functions
        const { getAccessToken } = require('../../api/client/userApi');
        getAccessToken.mockReturnValue(mockToken);
        jwtDecode.mockReturnValue(mockDecodedToken);

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('user-role')).toHaveTextContent('admin');
        });

        expect(screen.getByTestId('user-id')).toHaveTextContent('user123');
        expect(screen.getByTestId('user-email')).toHaveTextContent('redacted-email@example.com');
    });

    it('should default to customer role if no role in token', async () => {
        const mockToken = 'mock.jwt.token';
        const mockDecodedToken = {
            userId: 'user456',
            email: 'redacted-email@example.com',
            userName: 'noroleuser',
            // No role field
            lat: 40.7128,
            lng: -74.0060,
            exp: Math.floor(Date.now() / 1000) + 3600
        };

        const { getAccessToken } = require('../../api/client/userApi');
        getAccessToken.mockReturnValue(mockToken);
        jwtDecode.mockReturnValue(mockDecodedToken);

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('user-role')).toHaveTextContent('customer');
        });
    });

    it('should have admin role in dev bypass mode', async () => {
        // Mock dev environment with bypass
        const originalEnv = process.env.NODE_ENV;
        process.env.NODE_ENV = 'development';
        
        // Mock window.location.search
        delete window.location;
        window.location = { search: '?dev-bypass-auth=true' };

        render(
            <AuthProvider>
                <TestComponent />
            </AuthProvider>
        );

        await waitFor(() => {
            expect(screen.getByTestId('user-role')).toHaveTextContent('admin');
        });

        expect(screen.getByTestId('user-id')).toHaveTextContent('dev-user-123');

        // Restore
        process.env.NODE_ENV = originalEnv;
    });
});