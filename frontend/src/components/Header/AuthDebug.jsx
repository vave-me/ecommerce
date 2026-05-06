"use client";
import React from 'react';
import { useAuth } from '../../context/AuthContext';
import { getAccessToken, isTokenValid } from '../../api/client/userApi';
/**
 * AuthDebug Component - Shows authentication state in development
 * Only renders in development mode
 */
const AuthDebug = () => {
    const { user, isLoading, authChecked } = useAuth();
    // Only show in development
    if (process.env.NODE_ENV !== 'development') {
        return null;
    }
    const token = getAccessToken();
    const tokenValid = token ? isTokenValid(token) : false;
    
    // Log role information when user changes
    React.useEffect(() => {
        if (user) {
            
        } else {
            
        }
    }, [user]);
    return (
        <div style={{
            position: 'fixed',
            top: '10px',
            right: '10px',
            background: 'rgba(0, 0, 0, 0.8)',
            color: 'white',
            padding: '10px',
            borderRadius: '5px',
            fontSize: '12px',
            zIndex: 10000,
            maxWidth: '300px',
            fontFamily: 'monospace'
        }}>
            <div style={{ fontWeight: 'bold', marginBottom: '5px' }}>
                🔐 Auth Debug
            </div>
            <div>Auth Checked: {authChecked ? '✅' : '❌'}</div>
            <div>Loading: {isLoading ? '⏳' : '✅'}</div>
            <div>User: {user ? `✅ ${user.userId}` : '❌ None'}</div>
            <div>Token: {token ? '✅ Present' : '❌ Missing'}</div>
            <div>Token Valid: {tokenValid ? '✅' : '❌'}</div>
            {user && (
                <div style={{ marginTop: '5px', fontSize: '11px' }}>
                    <div>Email: {user.email || 'N/A'}</div>
                    <div>Username: {user.username || 'N/A'}</div>
                    <div>Role: {user.role || 'N/A'}</div>
                </div>
            )}
        </div>
    );
};
export default AuthDebug; 