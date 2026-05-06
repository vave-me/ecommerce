"use client";
import React, { useState, useEffect } from 'react';
import { useAuth } from '../../context/AuthContext';
import { getActivity, createActivity, getInteractions } from '../../api/client/activityApi';
/**
 * ActivityDebug Component - Comprehensive debugging for activity issues
 * Only renders in development mode
 */
const ActivityDebug = () => {
    const { user } = useAuth();
    const [debugInfo, setDebugInfo] = useState({
        user: null,
        userId: null,
        apiTests: {},
        errors: [],
        loading: false
    });
    // Only show in development
    if (process.env.NODE_ENV !== 'development') {
        return null;
    }
    const runAPITests = async () => {
        if (!user) return;
        setDebugInfo(prev => ({ ...prev, loading: true, errors: [] }));
        const userId = user.userId || user.id;
        const tests = {};
        const errors = [];
        try {
            // Test 1: Get Interactions (need to get activity first)
            try {
                // First get the user's activity
                const activityResponse = await getActivity(userId);
                const activityId = activityResponse?.activityId;
                
                if (activityId) {
                    // Then get interactions for that activity
                    const interactionsResponse = await getInteractions(activityId);
                    tests.getInteractions = {
                        success: true,
                        response: interactionsResponse,
                        interactionsCount: interactionsResponse?.interactions?.length || 0
                    };
                } else {
                    tests.getInteractions = { 
                        success: false, 
                        error: 'No activityId found for user' 
                    };
                }
            } catch (error) {
                tests.getInteractions = { success: false, error: error.message };
                errors.push(`getInteractions: ${error.message}`);
            }
            // Test 2: Get Activity (by userId)
            try {
                const activityResponse = await getActivity(userId);
                tests.getActivity = {
                    success: true,
                    response: activityResponse
                };
            } catch (error) {
                tests.getActivity = { success: false, error: error.message };
                errors.push(`getActivity: ${error.message}`);
            }
            // Test 3: Create Activity (if needed)
            try {
                const createResponse = await createActivity(userId);
                tests.createActivity = {
                    success: true,
                    response: createResponse
                };
            } catch (error) {
                tests.createActivity = { success: false, error: error.message };
                errors.push(`createActivity: ${error.message}`);
            }
        } catch (globalError) {
            errors.push(`Global error: ${globalError.message}`);
        }
        setDebugInfo(prev => ({
            ...prev,
            user,
            userId,
            apiTests: tests,
            errors,
            loading: false
        }));
    };
    useEffect(() => {
        if (user) {
            runAPITests();
        }
    }, [user]);
    const { user: debugUser, userId, apiTests, errors, loading } = debugInfo;
    return (
        <div style={{
            position: 'fixed',
            bottom: '10px',
            left: '10px',
            background: 'rgba(0, 0, 0, 0.9)',
            color: 'white',
            padding: '15px',
            borderRadius: '8px',
            fontFamily: 'monospace',
            fontSize: '12px',
            maxWidth: '400px',
            maxHeight: '300px',
            overflow: 'auto',
            zIndex: 999999,
            border: '2px solid #333'
        }}>
            <div style={{ marginBottom: '10px', fontWeight: 'bold', color: '#4CAF50' }}>
                🔧 Activity Debug Panel
            </div>
            <div style={{ marginBottom: '8px' }}>
                <strong>User:</strong> {debugUser ? `${debugUser.userId || debugUser.id} (${debugUser.email})` : 'Not logged in'}
            </div>
            <div style={{ marginBottom: '8px' }}>
                <strong>User ID:</strong> {userId || 'N/A'}
            </div>
            {loading && (
                <div style={{ color: '#FFC107' }}>⏳ Running API tests...</div>
            )}
            {errors.length > 0 && (
                <div style={{ marginBottom: '8px' }}>
                    <strong style={{ color: '#F44336' }}>Errors:</strong>
                    {errors.map((error, index) => (
                        <div key={index} style={{ color: '#F44336', fontSize: '11px' }}>
                            • {error}
                        </div>
                    ))}
                </div>
            )}
            <div style={{ marginBottom: '8px' }}>
                <strong>API Tests:</strong>
                {Object.entries(apiTests).map(([testName, result]) => (
                    <div key={testName} style={{ marginLeft: '10px', fontSize: '11px' }}>
                        <span style={{ color: result.success ? '#4CAF50' : '#F44336' }}>
                            {result.success ? '✅' : '❌'} {testName}
                        </span>
                        {result.success && result.response && (
                            <div style={{ marginLeft: '15px', color: '#9E9E9E' }}>
                                {testName === 'getInteractions' && `(${result.interactionsCount} interactions)`}
                                {testName === 'getActivity' && result.response.activityId && `(ID: ${result.response.activityId})`}
                                {testName === 'createActivity' && result.response.id && `(ID: ${result.response.id})`}
                            </div>
                        )}
                    </div>
                ))}
            </div>
            <button 
                onClick={runAPITests}
                disabled={loading || !user}
                style={{
                    background: '#2196F3',
                    color: 'white',
                    border: 'none',
                    padding: '5px 10px',
                    borderRadius: '4px',
                    cursor: loading || !user ? 'not-allowed' : 'pointer',
                    fontSize: '11px',
                    opacity: loading || !user ? 0.5 : 1
                }}
            >
                {loading ? 'Testing...' : 'Run API Tests'}
            </button>
        </div>
    );
};
export default ActivityDebug; 