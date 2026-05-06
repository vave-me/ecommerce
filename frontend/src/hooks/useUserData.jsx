"use client";
import { useState, useEffect, useMemo } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import axios from '../api/axiosInstance';
import { getBaseUserById } from '../api/client/userApi';
import { useAuth } from '../context/AuthContext';
import { QUERY_KEYS } from '../lib/reactQuery';

// Simple in-memory cache for user data (for backward compatibility)
const userDataCache = new Map();

// Mock user data for development/fallback
const generateMockUserData = (userId) => {
    const mockNames = [
        'Alex Johnson', 'Sarah Chen', 'Mike Rodriguez', 'Emma Wilson', 'David Kim',
        'Lisa Thompson', 'John Davis', 'Maria Garcia', 'Chris Lee', 'Anna Brown'
    ];
    const mockAvatars = [
        '/images/user-user.webp', '/images/psyche.jpg'
    ];
    
    // Generate consistent mock data based on userId
    const hash = userId.split('').reduce((a, b) => {
        a = ((a << 5) - a) + b.charCodeAt(0);
        return a & a;
    }, 0);
    
    const nameIndex = Math.abs(hash) % mockNames.length;
    const avatarIndex = Math.abs(hash) % mockAvatars.length;
    
    return {
        userName: mockNames[nameIndex],
        firstName: mockNames[nameIndex].split(' ')[0],
        lastName: mockNames[nameIndex].split(' ')[1],
        avatar: mockAvatars[avatarIndex],
        thumbnail: mockAvatars[avatarIndex],
        email: '',
        isOnline: Math.abs(hash) % 2 === 0,
        lastSeen: new Date(Date.now() - Math.abs(hash) % 86400000).toISOString()
    };
};

/**
 * Unified hook for fetching user data
 * Combines functionality from useUserData and useUserProfile
 * Supports both React Query and legacy in-memory caching
 * 
 * @param {string} userId - The user ID to fetch data for
 * @param {Object} options - Options for the hook
 * @param {boolean} options.enabled - Whether to fetch data (default: true)
 * @param {Object} options.fallback - Fallback user data
 * @param {boolean} options.useCache - Whether to use in-memory cache (default: false, uses React Query)
 * @param {boolean} options.useReactQuery - Whether to use React Query (default: true)
 * @returns {Object} { userData, isLoading, error, refetch }
 */
export const useUserData = (userId, options = {}) => {
    const { user: currentUser } = useAuth();
    const targetUserId = userId || currentUser?.userId;
    
    const { 
        enabled = true, 
        useCache = false, // Legacy in-memory cache
        useReactQuery = true,
        fallback = { 
            userName: 'Anonymous', 
            avatar: '/images/user-user.webp' 
        },
        staleTime = 5 * 60 * 1000, // 5 minutes
        cacheTime = 30 * 60 * 1000, // 30 minutes
    } = options;
    
    // React Query implementation
    const queryResult = useQuery({
        queryKey: targetUserId ? QUERY_KEYS.user.byId(targetUserId) : ['user', 'anonymous'],
        queryFn: async () => {
            if (!targetUserId) {
                return fallback;
            }
            
            try {
                // Try modern API endpoint first
                const response = await axios.get(`/users/${targetUserId}`);
                return processUserData(response.data, targetUserId);
            } catch (error) {
                // Fallback to legacy API
                const response = await getBaseUserById(targetUserId);
                const user = response.user || response;
                return processUserData(user, targetUserId);
            }
        },
        enabled: enabled && useReactQuery && !!targetUserId,
        staleTime,
        cacheTime,
        refetchOnWindowFocus: false,
        onError: (error) => {
            // Error: 'Error fetching user data:', error...
        },
        // Provide fallback data
        placeholderData: generateMockUserData(targetUserId || 'anonymous'),
    });
    
    // Legacy state management (for backward compatibility)
    const [legacyUserData, setLegacyUserData] = useState(null);
    const [legacyIsLoading, setLegacyIsLoading] = useState(enabled && !!targetUserId && !useReactQuery);
    const [legacyError, setLegacyError] = useState(null);
    
    // Legacy implementation
    useEffect(() => {
        if (useReactQuery || !enabled || !targetUserId) {
            if (!targetUserId) {
                setLegacyUserData(fallback);
                setLegacyIsLoading(false);
            }
            return;
        }
        
        // Check cache first
        if (useCache && userDataCache.has(targetUserId)) {
            const cachedData = userDataCache.get(targetUserId);
            setLegacyUserData(cachedData);
            setLegacyIsLoading(false);
            return;
        }
        
        // Fetch user data
        const fetchUserData = async () => {
            try {
                setLegacyIsLoading(true);
                setLegacyError(null);
                
                const response = await getBaseUserById(targetUserId);
                const user = response.user || response;
                const processedData = processUserData(user, targetUserId);
                
                // Cache the result
                if (useCache) {
                    userDataCache.set(targetUserId, processedData);
                }
                
                setLegacyUserData(processedData);
            } catch (err) {
                setLegacyError(err);
                // Use mock data as fallback
                const mockData = generateMockUserData(targetUserId);
                
                if (useCache) {
                    userDataCache.set(targetUserId, mockData);
                }
                
                setLegacyUserData(mockData);
            } finally {
                setLegacyIsLoading(false);
            }
        };
        
        fetchUserData();
    }, [targetUserId, enabled, useCache, useReactQuery]);
    
    // Return unified result
    if (useReactQuery) {
        return {
            userData: queryResult.data,
            isLoading: queryResult.isLoading,
            error: queryResult.error,
            refetch: queryResult.refetch,
            isError: queryResult.isError,
            isFetching: queryResult.isFetching,
        };
    }
    
    return {
        userData: legacyUserData,
        isLoading: legacyIsLoading,
        error: legacyError,
        refetch: () => {
            // Clear cache and refetch for legacy mode
            if (targetUserId) {
                userDataCache.delete(targetUserId);
                setLegacyUserData(null);
                setLegacyIsLoading(true);
            }
        }
    };
};

/**
 * Process user data to ensure consistent format
 */
function processUserData(user, userId) {
    return {
        id: user.id || userId,
        userName: user.userName || user.username || 'Anonymous',
        firstName: user.firstName || '',
        lastName: user.lastName || '',
        avatar: user.thumbnail || user.avatar || user.profileImage || '/images/user-user.webp',
        thumbnail: user.thumbnail || user.avatar || user.profileImage || '/images/user-user.webp',
        location: user.location || '',
        bio: user.bio || '',
        privacy: user.privacy || '',
        background: user.background || '',
        lat: user.lat || 0,
        lng: user.lng || 0,
        email: user.email || '',
        isOnline: user.isOnline || false,
        lastSeen: user.lastSeen || null,
        // Include any additional fields
        ...user
    };
}

/**
 * Hook for fetching multiple users data
 * @param {string[]} userIds - Array of user IDs to fetch
 * @param {Object} options - Options for the hook
 * @returns {Object} { usersData, isLoading, errors }
 */
export const useUsersData = (userIds, options = {}) => {
    const { enabled = true, useCache = true, useReactQuery = true } = options;
    const queryClient = useQueryClient();
    
    // React Query implementation for multiple users
    const queries = useQuery({
        queryKey: ['users', 'multiple', userIds?.join(',') || ''],
        queryFn: async () => {
            if (!userIds?.length) return {};
            
            const results = {};
            const errors = {};
            
            await Promise.allSettled(
                userIds.map(async (userId) => {
                    try {
                        // Check React Query cache first
                        const cachedData = queryClient.getQueryData(QUERY_KEYS.user.byId(userId));
                        if (cachedData) {
                            results[userId] = cachedData;
                            return;
                        }
                        
                        // Fetch from API
                        const response = await getBaseUserById(userId);
                        const user = response.user || response;
                        results[userId] = processUserData(user, userId);
                        
                        // Update individual user cache
                        queryClient.setQueryData(
                            QUERY_KEYS.user.byId(userId),
                            results[userId]
                        );
                    } catch (err) {
                        errors[userId] = err;
                        results[userId] = generateMockUserData(userId);
                    }
                })
            );
            
            return { results, errors };
        },
        enabled: enabled && useReactQuery && !!userIds?.length,
        staleTime: 5 * 60 * 1000,
        cacheTime: 30 * 60 * 1000,
    });
    
    if (useReactQuery) {
        return {
            usersData: queries.data?.results || {},
            isLoading: queries.isLoading,
            errors: queries.data?.errors || {},
            refetch: queries.refetch,
        };
    }
    
    // Legacy implementation
    const [usersData, setUsersData] = useState({});
    const [isLoading, setIsLoading] = useState(false);
    const [errors, setErrors] = useState({});
    
    useEffect(() => {
        if (!enabled || !userIds?.length || useReactQuery) return;
        
        const fetchUsersData = async () => {
            setIsLoading(true);
            const newUsersData = {};
            const newErrors = {};
            
            // Check cache first for each user
            const userIdsToFetch = [];
            userIds.forEach(userId => {
                if (useCache && userDataCache.has(userId)) {
                    newUsersData[userId] = userDataCache.get(userId);
                } else {
                    userIdsToFetch.push(userId);
                }
            });
            
            // Fetch only uncached users
            if (userIdsToFetch.length > 0) {
                await Promise.allSettled(
                    userIdsToFetch.map(async (userId) => {
                        try {
                            const response = await getBaseUserById(userId);
                            const user = response.user || response;
                            const processedUserData = processUserData(user, userId);
                            
                            newUsersData[userId] = processedUserData;
                            
                            // Cache the result
                            if (useCache) {
                                userDataCache.set(userId, processedUserData);
                            }
                        } catch (err) {
                            newErrors[userId] = err;
                            // Use mock data as fallback
                            const mockData = generateMockUserData(userId);
                            newUsersData[userId] = mockData;
                            
                            // Cache the mock data
                            if (useCache) {
                                userDataCache.set(userId, mockData);
                            }
                        }
                    })
                );
            }
            
            setUsersData(newUsersData);
            setErrors(newErrors);
            setIsLoading(false);
        };
        
        fetchUsersData();
    }, [userIds?.join(','), enabled, useCache, useReactQuery]);
    
    return { usersData, isLoading, errors };
};

/**
 * Hook for fetching user profile (alias for useUserData with React Query)
 * Maintained for backward compatibility
 */
export function useUserProfile(userId) {
    return useUserData(userId, { useReactQuery: true });
}

/**
 * Clear user data cache (useful for when user data updates)
 * Works for both React Query and legacy cache
 * @param {string} userId - Optional specific user ID to clear
 */
export const clearUserDataCache = (userId = null) => {
    const queryClient = useQueryClient();
    
    if (userId) {
        // Clear React Query cache
        queryClient.invalidateQueries(QUERY_KEYS.user.byId(userId));
        // Clear legacy cache
        userDataCache.delete(userId);
    } else {
        // Clear all user queries
        queryClient.invalidateQueries(['user']);
        // Clear legacy cache
        userDataCache.clear();
    }
};

export default useUserData;