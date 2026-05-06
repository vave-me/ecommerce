"use client";
import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../../context/AuthContext';
/**
 * AuthGuard Component - Protects routes from unauthorized access
 * Prevents flash of protected content by showing loading state until auth is confirmed
 */
const AuthGuard = ({ 
    children, 
    fallback = null, 
    redirectTo = '/login',
    showLoading = true,
    loadingComponent = null 
}) => {
    const { user, isLoading, authChecked } = useAuth();
    const router = useRouter();
    // Redirect to login if user is not authenticated (only after auth check is complete)
    useEffect(() => {
        if (authChecked && !isLoading && !user) {
            if (process.env.NODE_ENV === 'development') {
            }
            router.push(redirectTo);
        }
    }, [user, isLoading, authChecked, router, redirectTo]);
    // Show loading state while authentication is being checked
    if (isLoading || !authChecked) {
        if (loadingComponent) {
            return loadingComponent;
        }
        if (showLoading) {
            return (
                <div style={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    minHeight: '200px',
                    padding: '2rem',
                    textAlign: 'center'
                }}>
                    <div style={{
                        width: '32px',
                        height: '32px',
                        border: '3px solid #f3f3f3',
                        borderTop: '3px solid #3498db',
                        borderRadius: '50%',
                        animation: 'spin 1s linear infinite',
                        marginBottom: '1rem'
                    }}></div>
                    <p style={{ color: '#666', fontSize: '14px' }}>
                        Checking authentication...
                    </p>
                    <style jsx>{`
                        @keyframes spin {
                            0% { transform: rotate(0deg); }
                            100% { transform: rotate(360deg); }
                        }
                    `}</style>
                </div>
            );
        }
        return fallback;
    }
    // Show fallback if user is not authenticated
    if (!user) {
        return fallback;
    }
    // User is authenticated, render children
    return children;
};
export default AuthGuard; 