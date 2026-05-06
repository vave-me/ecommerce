// src/api/config.js - Centralized API configuration
export const API_CONFIG = {
    BASE_URL: process.env.NEXT_PUBLIC_API_BASE_URL || '/api',
    timeout: 30000, // 30 seconds
    headers: {
        'Content-Type': 'application/json',
    },
    // Rate limiting configuration
    rateLimit: {
        maxRequests: 10,
        perTimeWindow: 1000, // 1 second
    }
};
