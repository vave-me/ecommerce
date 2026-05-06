"use client";
import React, {createContext, useContext, useState, useEffect, useCallback, useMemo} from 'react';
import {jwtDecode} from 'jwt-decode';
import {useRouter} from 'next/navigation';
import {toast} from 'react-toastify';
import {useQueryClient} from '@tanstack/react-query';
import axios from '../api/axiosInstance';
import {
    getAccessToken,
    setAccessToken,
    setRefreshToken,
    clearTokens,
    refreshAccessToken,
    initFromLocalStorage,
    logoutUser,
    clearUserTokens,
    getRefreshToken
} from '../api/client/userApi';
import {
    loginUser,
    registerUser,
    loginWithGoogle
} from '../api/userApi';

const AuthContext = createContext();
export const useAuth = () => useContext(AuthContext);

// Helper function to validate redirect URLs for security
const validateRedirectUrl = (url) => {
    if (!url || typeof url !== 'string') return '/';

    try {
        // If it's a relative URL, it's safe
        if (url.startsWith('/') && !url.startsWith('//')) {
            return url;
        }

        // If it's an absolute URL, check if it's from the same origin
        const urlObj = new URL(url, window.location.origin);
        if (urlObj.origin === window.location.origin) {
            return urlObj.pathname + urlObj.search + urlObj.hash;
        }

        // If it's an external URL, redirect to home for security
        return '/';
    } catch (error) {
        // If URL parsing fails, redirect to home
        return '/';
    }
};

export const AuthProvider = ({children}) => {
    const router = useRouter();
    const queryClient = useQueryClient();
    const [mounted, setMounted] = useState(false);
    const [user, setUser] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [authChecked, setAuthChecked] = useState(false);

    // Development bypass for testing theme functionality
    const isDevelopment = process.env.NODE_ENV === 'development';
    const devBypassAuth = isDevelopment && typeof window !== 'undefined' && 
        window.location.search.includes('dev-bypass-auth=true');

    // Helper function to clear user-related React Query caches
    const clearUserQueries = useCallback(async () => {
        try {
            // Remove all queries that depend on userId
            await queryClient.removeQueries({
                predicate: (query) => {
                    const queryKey = query.queryKey;
                    if (!Array.isArray(queryKey)) return false;
                    
                    // Remove queries containing userId or userCatalog
                    return queryKey.some(key => 
                        typeof key === 'string' && (
                            key.includes('userCatalog') ||
                            key.includes('userId') ||
                            key.includes('user')
                        )
                    ) || queryKey.includes('userCatalog') || queryKey.includes('publicCatalog');
                }
            });
            
            // Clear specific query patterns
            await queryClient.removeQueries(['userCatalog']);
            await queryClient.removeQueries(['publicCatalog']);
            await queryClient.removeQueries(['user']);
            await queryClient.removeQueries(['listings']);
            
            } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }, [queryClient]);

    // Helper function to invalidate user-related React Query caches (for login)
    const invalidateUserQueries = useCallback(async () => {
        try {
            // Invalidate all user-related queries to refetch with new auth
            await queryClient.invalidateQueries({
                predicate: (query) => {
                    const queryKey = query.queryKey;
                    if (!Array.isArray(queryKey)) return false;
                    
                    // Invalidate queries containing userId or userCatalog
                    return queryKey.some(key => 
                        typeof key === 'string' && (
                            key.includes('userCatalog') ||
                            key.includes('userId') ||
                            key.includes('user')
                        )
                    ) || queryKey.includes('userCatalog') || queryKey.includes('publicCatalog');
                }
            });
            
            // Invalidate specific query patterns
            await queryClient.invalidateQueries(['userCatalog']);
            await queryClient.invalidateQueries(['publicCatalog']);
            await queryClient.invalidateQueries(['user']);
            await queryClient.invalidateQueries(['listings']);
            
            } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }, [queryClient]);

    // Parse user information from JWT token
    const parseUserFromToken = useCallback((token) => {
        try {
            if (!token) {
                if (process.env.NODE_ENV === 'development') {
                    
                }
                return null;
            }

            const decoded = jwtDecode(token);

            if (process.env.NODE_ENV === 'development') {
                // Token info logged for debugging
            }

            // Check token expiration
            if (decoded.exp * 1000 < Date.now()) {
                if (process.env.NODE_ENV === 'development') {
                    
                }
                return null;
            }

            const userId = decoded.userId || decoded.sub;
            const email = decoded.email || '';
            const username = decoded.userName;
            const lat = decoded.lat;
            const lng = decoded.lng;
            const role = decoded.role || 'customer'; // Default to customer if no role

            const userData = {userId, email, username, lat, lng, role};

            if (process.env.NODE_ENV === 'development') {
                
            }

            return userData;
        } catch (error) {
            if (process.env.NODE_ENV === 'development') {
                // Error: '🔐 parseUserFromToken: Error parsing token:', err
            }
            return null;
        }
    }, []);

    // Clear all auth data
    const clearAuthData = useCallback(async () => {
        setUser(null);
        await clearTokens();
        // Clear user-related React Query caches
        await clearUserQueries();
    }, [clearUserQueries]);

    // Refresh token and update user state
    const refreshTokenAndSetUser = useCallback(async () => {
        try {
            const newToken = await refreshAccessToken();
            const userData = parseUserFromToken(newToken);

            if (userData) {
                setUser(userData);
            } else {
                throw new Error('Invalid user data from refreshed token');
            }

            return newToken;
        } catch (error) {
            await clearAuthData();
            throw error;
        }
    }, [parseUserFromToken, clearAuthData]);

    // Only initialize auth state after component mounts on client
    useEffect(() => {
        const initializeAuth = async () => {
            try {
                setIsLoading(true);
                setAuthChecked(false); // Reset auth check state
                
                if (process.env.NODE_ENV === 'development') {
                    
                }

                // Development bypass for testing
                if (devBypassAuth) {
                    
                    setUser({
                        userId: 'dev-user-123',
                        email: 'redacted-email@example.com',
                        username: 'dev-user',
                        lat: null,
                        lng: null,
                        role: 'admin' // Dev mode gets admin role for testing
                    });
                    setIsLoading(false);
                    setMounted(true);
                    setAuthChecked(true);
                    return;
                }

                // For migration period: check legacy localStorage format
                initFromLocalStorage();

                // Try to initialize auth from localStorage tokens
                let success = false;
                try {
                    const userApiModule = await import('../api/client/userApi');
                    success = userApiModule.initializeAuth();
                    
                    if (process.env.NODE_ENV === 'development') {
                        // Auth initialization result logged
                    }
                } catch (importError) {
                    // Error: '🔐 AuthContext: Failed to import userApi:', impor
                    success = false;
                }

                if (success) {
                    // Get token from localStorage
                    const token = getAccessToken();
                    
                    if (process.env.NODE_ENV === 'development') {
                        
                    }

                    if (token) {
                        // Parse user from token
                        const userData = parseUserFromToken(token);
                        
                        if (process.env.NODE_ENV === 'development') {
                            
                        }

                        if (userData) {
                            setUser(userData);
                            
                            if (process.env.NODE_ENV === 'development') {
                                
                            }
                        } else {
                            // Token is invalid, try to refresh
                            if (process.env.NODE_ENV === 'development') {
                                
                            }
                            
                            try {
                                await refreshTokenAndSetUser();
                                
                                if (process.env.NODE_ENV === 'development') {
                                    
                                }
                            } catch (refreshError) {
                                // Clear invalid auth state
                                await clearAuthData();
                                
                                if (process.env.NODE_ENV === 'development') {
                                    // Error: '🔐 AuthContext: Token refresh failed:', refreshEr
                                }
                            }
                        }
                    } else {
                        if (process.env.NODE_ENV === 'development') {
                            
                        }
                    }
                } else {
                    if (process.env.NODE_ENV === 'development') {
                        // Auth initialization returned false
                    }
                }
            } catch (error) {
                if (process.env.NODE_ENV === 'development') {
                    // Error: '🔐 AuthContext: Error during initialization:', er
                }
                
                // Clear any corrupted auth state
                await clearAuthData();
            } finally {
                setIsLoading(false);
                setMounted(true);
                setAuthChecked(true);
                
                if (process.env.NODE_ENV === 'development') {
                    
                }
            }
        };

        initializeAuth();
    }, [devBypassAuth, parseUserFromToken, clearAuthData, refreshTokenAndSetUser]);

    // Set auth tokens (helper function)
    const setAuthTokens = useCallback(async (accessToken, refreshToken) => {
        setAccessToken(accessToken);
        if (refreshToken) {
            setRefreshToken(refreshToken);
        }
    }, []);

    // Set auth data from token
    const setAuthData = useCallback(async (token, refreshToken) => {
        try {
            // Store tokens
            if (refreshToken) {
                await setAuthTokens(token, refreshToken);
            } else {
                setAccessToken(token);
            }

            // Parse and set user data
            const userData = parseUserFromToken(token);
            if (userData) {
                setUser(userData);
                // Invalidate user-related queries to refetch with new auth
                await invalidateUserQueries();
            } else {
                throw new Error('Invalid token data');
            }
        } catch (error) {
            if (process.env.NODE_ENV === 'development') {
                // Error: 'Failed to set auth data:', error
            }
            await clearAuthData();
            throw error;
        }
    }, [parseUserFromToken, clearAuthData, setAuthTokens, invalidateUserQueries]);

    // Login with credentials using public userApi
    const signInWithCredentials = useCallback(async ({email, password}, redirectUrl = '/') => {
        try {
            const result = await loginUser({email, password});

            // Get tokens from response
            const token = result.accessToken || result.token;
            const refreshToken = result.refreshToken;

            if (!token) {
                throw new Error('No access token received from login');
            }

            // Store tokens manually since public API doesn't do this
            setAccessToken(token);
            if (refreshToken) {
                setRefreshToken(refreshToken);
            }

            // Set auth data
            await setAuthData(token, refreshToken);

            toast.success('Successfully logged in!');
            const safeRedirectUrl = validateRedirectUrl(redirectUrl);
            router.push(safeRedirectUrl);
        } catch (error) {
            if (process.env.NODE_ENV === 'development') {
                // Error: 'Login failed:', error
            }
            toast.error(error.response?.data?.message || error.message || 'Login failed. Please try again.');
            throw error;
        }
    }, [setAuthData, router]);

    // Registration using userApi
    const signUpWithCredentials = useCallback(async ({
                                                         email,
                                                         password,
                                                         firstName,
                                                         lastName,
                                                         userName,
                                                         address,
                                                         location,
                                                         role = 'customer'
                                                     }, redirectUrl = '/') => {
        try {
            await registerUser({
                email,
                firstName,
                lastName,
                userName,
                password,
                address: address || location,
                role
            });

            toast.success('Registration successful! Logging you in...');
            await signInWithCredentials({email, password}, redirectUrl);
        } catch (error) {
            if (process.env.NODE_ENV === 'development') {
                // Error: 'Registration failed:', error
            }
            toast.error(error.response?.data?.message || error.message || 'Registration failed. Please try again.');
            throw error;
        }
    }, [signInWithCredentials]);

    // Logout using userApi
    const signOutUser = useCallback(async () => {
        try {
            // Use the logoutUser function from userApi which handles the API call and token clearing
            await logoutUser();

            // Clear local state
            setUser(null);
            
            // Clear user-related React Query caches
            await clearUserQueries();
            
            toast.info('Successfully logged out.');
            router.push('/login');
        } catch (error) {
            if (process.env.NODE_ENV === 'development') {
                // Error: 'Logout failed:', error
            }
            toast.error('Logout failed. Please try again.');

            // Clear local auth data even if server logout fails
            await clearAuthData();
            router.push('/login');
        }
    }, [clearAuthData, clearUserQueries, router]);

    const setUserOnlineStatus = useCallback((online) => {
        setUser((prev) => {
            if (!prev) return null;
            return {...prev, isOnline: online};
        });
    }, []);

    // Google login using public userApi
    const signInWithGoogle = useCallback(async (idToken, redirectUrl = '/') => {
        try {
            const result = await loginWithGoogle(idToken);

            // Get tokens from response
            const token = result.accessToken || result.token;
            const refreshToken = result.refreshToken;

            if (!token) {
                throw new Error('Invalid response: No token received');
            }

            // Store tokens manually since public API doesn't do this
            setAccessToken(token);
            if (refreshToken) {
                setRefreshToken(refreshToken);
            }

            // Set auth data
            await setAuthData(token, refreshToken);

            toast.success('Successfully logged in with Google!');
            const safeRedirectUrl = validateRedirectUrl(redirectUrl);
            router.push(safeRedirectUrl);
        } catch (error) {
            if (process.env.NODE_ENV === 'development') {
                // Error: 'Google login failed:', error
            }
            const errorMessage = error.response?.data?.message || error.message || 'Google login failed. Please try again.';
            toast.error(errorMessage);
            throw error;
        }
    }, [setAuthData, router]);

    // Clear user tokens with reason - uses the clear-tokens API endpoint
    const clearUserTokensWithReason = useCallback(async (reason = "security concern") => {
        try {
            if (!user?.userId) {
                throw new Error('No user ID available for token clearing');
            }

            // Use the clearUserTokens function from userApi (which calls /users/clear-tokens)
            // Pass current access token as tokenId since we want to clear the current session
            const currentToken = getAccessToken();
            await clearUserTokens(user.userId, currentToken, getRefreshToken(), reason);
            
            // Clear local state
            await clearAuthData();
            
            toast.info('Tokens have been invalidated for security.');
            router.push('/login');
        } catch (error) {
            if (process.env.NODE_ENV === 'development') {
                // Error: 'Clear tokens failed:', error
            }
            // Clear local auth data even if server clear fails
            await clearAuthData();
            router.push('/login');
        }
    }, [user?.userId, clearAuthData, router]);

    const contextValue = useMemo(() => ({
        user,
        isLoading,
        authChecked,
        signInWithCredentials,
        signUpWithCredentials,
        signOutUser,
        setUserOnlineStatus,
        signInWithGoogle,
        refreshTokenAndSetUser,
        clearUserTokensWithReason
    }), [
        user,
        isLoading,
        authChecked,
        signInWithCredentials,
        signUpWithCredentials,
        signOutUser,
        setUserOnlineStatus,
        signInWithGoogle,
        refreshTokenAndSetUser,
        clearUserTokensWithReason
    ]);

    return (
        <AuthContext.Provider value={contextValue}>
            {children}
        </AuthContext.Provider>
    );
};
