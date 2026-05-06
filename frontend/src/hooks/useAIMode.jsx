import { useState, useEffect, useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { switchToAiMode, selectIsAiMode } from '../redux/slices/appModeSlice';
import { useAuth } from '../context/AuthContext';
import { useAIService } from './ai/useAIService';

/**
 * AI Mode Redux Integration Hook
 * 
 * This hook provides Redux integration for AI mode management.
 * It wraps the core useAI hook and adds:
 * - Redux state synchronization
 * - Auto-initialization when switching to AI mode
 * - Initialization status tracking
 * - Retry functionality for error recovery
 * 
 * @returns {Object} AI mode state and operations
 */
const useAIMode = () => {
    const dispatch = useDispatch();
    const { user } = useAuth();
    const isAiMode = useSelector(selectIsAiMode);
    
    // Get all functionality from AI Service hook
    const ai = useAIService();
    
    // Additional state for mode management
    const [isInitializing, setIsInitializing] = useState(false);
    const [initializationError, setInitializationError] = useState(null);
    const [hasInitialized, setHasInitialized] = useState(false);
    
    /**
     * Auto-initialize when switching to AI mode
     */
    useEffect(() => {
        const autoInitialize = async () => {
            if (isAiMode && user?.id && !hasInitialized && !isInitializing) {
                setIsInitializing(true);
                setInitializationError(null);
                
                try {
                    // In the new service, initialization is automatic
                    await ai.loadAssistants();
                    if (ai.assistants.length > 0) {
                        setHasInitialized(true);
                    } else {
                        throw new Error('No assistants available');
                    }
                } catch (error) {
                    setInitializationError(error.message);
                } finally {
                    setIsInitializing(false);
                }
            }
        };
        
        autoInitialize();
    }, [isAiMode, user?.id, hasInitialized, isInitializing, ai.loadAssistants, ai.assistants.length]);
    
    /**
     * Reset initialization state when user changes
     */
    useEffect(() => {
        if (!user?.id) {
            setHasInitialized(false);
            setInitializationError(null);
            setIsInitializing(false);
        }
    }, [user?.id]);
    
    /**
     * Enable AI mode through Redux
     */
    const enableAIMode = useCallback(() => {
        if (!user?.id) {
            
            return;
        }
        dispatch(switchToAiMode());
    }, [dispatch, user?.id]);
    
    /**
     * Force re-initialization (for error recovery)
     */
    const retryInitialization = useCallback(() => {
        if (!user?.id) return;
        
        setHasInitialized(false);
        setInitializationError(null);
        // This will trigger the useEffect to re-initialize
    }, [user?.id]);
    
    /**
     * Get initialization status
     */
    const getInitializationStatus = useCallback(() => {
        if (!user?.id) {
            return { status: 'unauthenticated', message: 'Please sign in' };
        }
        if (!isAiMode) {
            return { status: 'disabled', message: 'AI mode not enabled' };
        }
        if (isInitializing) {
            return { status: 'initializing', message: 'Setting up your AI assistants...' };
        }
        if (initializationError) {
            return { status: 'error', message: initializationError };
        }
        if (!hasInitialized) {
            return { status: 'pending', message: 'Waiting for initialization...' };
        }
        if (ai.assistants.length === 0) {
            return { status: 'no_assistants', message: 'No assistants available' };
        }
        return { status: 'ready', message: 'AI mode ready' };
    }, [user?.id, isAiMode, isInitializing, initializationError, hasInitialized, ai.assistants.length]);
    
    /**
     * Return combined API
     */
    return {
        // All core AI functionality
        ...ai,
        
        // Mode-specific state
        isAiMode,
        isInitializing,
        hasInitialized,
        initializationError,
        
        // Mode-specific actions
        enableAIMode,
        retryInitialization,
        
        // Status helpers
        status: getInitializationStatus(),
        isReady: hasInitialized && ai.assistants.length > 0 && !ai.hasErrors,
        hasDefaultAssistant: ai.assistants.some(a => a.isDefault),
        primaryAssistant: ai.assistants.find(a => a.isDefault) || ai.assistants[0] || null
    };
};

export default useAIMode;