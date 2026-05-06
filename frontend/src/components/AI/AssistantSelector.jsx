"use client";
import React, {useState, useEffect, useCallback, useMemo, memo} from 'react';
import {useTranslations} from 'next-intl';
import {
    Bot,
    ChevronDown,
    Plus,
    Image as ImageIcon,
    Mic,
    MessageSquare,
    Zap,
    Check
} from '@/icons';
import { getAssistantService } from '../../services/ai/AssistantService';
import {useAuth} from '../../context/AuthContext';
import styles from './AssistantSelector.module.css';
/**
 * Compact Assistant Selector - Modern ChatGPT style
 * OPTIMIZED: React.memo with custom comparison for performance
 */
const AssistantSelector = memo(function AssistantSelector({
                                                              selectedAssistantId = null,
                                                              onAssistantSelect = () => {
                                                              },
                                                              className = ''
                                                          }) {
    const t = useTranslations('AI');
    const {user} = useAuth();
    const [assistants, setAssistants] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isOpen, setIsOpen] = useState(false);
    const [error, setError] = useState(null);
    /**
     * Load active assistants
     */
    const loadAssistants = useCallback(async () => {
        if (!user) return;
        setIsLoading(true);
        setError(null);
        try {
            const assistantService = getAssistantService();
            const response = await assistantService.getAssistants();
            if (response.success) {
                const assistantsData = response.data?.assistants || [];
                if (Array.isArray(assistantsData)) {
                    // No transformation - use API response directly
                    setAssistants(assistantsData);
                    // Only auto-select first assistant on initial load, not on every re-render
                    if (!selectedAssistantId && assistantsData.length > 0) {
                        onAssistantSelect(assistantsData[0]);
                    }
                } else {
                    setAssistants([]);
                }
            } else {
                setError(response.error || 'Failed to load assistants');
            }
        } catch (err) {
            setError('Failed to load assistants');
        } finally {
            setIsLoading(false);
        }
    }, [user]); // REMOVED selectedAssistantId and onAssistantSelect to prevent infinite loop
    // Load assistants on mount
    useEffect(() => {
        loadAssistants();
    }, [loadAssistants]);
    // Handle assistant selection
    const handleSelect = useCallback((assistant) => {
        onAssistantSelect(assistant);
        setIsOpen(false);
    }, [onAssistantSelect]);
    // Toggle dropdown
    const toggleDropdown = useCallback((e) => {
        e.stopPropagation();
        setIsOpen(!isOpen);
    }, [isOpen]);
    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (isOpen) {
                setIsOpen(false);
            }
        };
        if (isOpen) {
            document.addEventListener('click', handleClickOutside);
            return () => document.removeEventListener('click', handleClickOutside);
        }
    }, [isOpen]);
    // Render capability indicators
    const renderCapabilities = useCallback((assistant) => {
        const icons = [];
        // Simple checks based on backend response
        if (assistant?.supports_vision) {
            icons.push(<ImageIcon key="image" size={14} title="Image analysis"/>);
        }
        if (assistant?.supports_audio) {
            icons.push(<Mic key="speech" size={14} title="Speech processing"/>);
        }
        return icons.length > 0 ? (
            <div className={styles.capabilities}>
                {icons}
            </div>
        ) : null;
    }, []);
    // Find selected assistant
    const selectedAssistant = assistants.find(a => a.id === selectedAssistantId) || assistants[0];
    if (!user) {
        return null;
    }
    if (error) {
        return (
            <div className={`${styles.selector} ${styles.error} ${className}`}>
                <span>Error loading assistants</span>
            </div>
        );
    }
    if (isLoading || assistants.length === 0) {
        return (
            <div className={`${styles.selector} ${styles.loading} ${className}`}>
                <Bot size={16}/>
                <span>{isLoading ? 'Loading...' : 'No assistants'}</span>
            </div>
        );
    }
    return (
        <div className={`${styles.container} ${className}`}>
            <button
                onClick={toggleDropdown}
                className={`${styles.selector} ${isOpen ? styles.open : ''}`}
                disabled={isLoading}
            >
                <div className={styles.selectedAssistant}>
                    <Bot size={16}/>
                    <div className={styles.assistantInfo}>
                        <span className={styles.assistantName}>
                            {selectedAssistant?.name || 'Select Assistant'}
                        </span>
                        {selectedAssistant && renderCapabilities(selectedAssistant)}
                    </div>
                </div>
                <ChevronDown size={16} className={styles.chevron}/>
            </button>
            {isOpen && (
                <div className={styles.dropdown}>
                    <div className={styles.dropdownContent}>
                        {assistants.map((assistant) => (
                            <button
                                key={assistant.id}
                                onClick={() => handleSelect(assistant)}
                                className={`${styles.assistantOption} ${
                                    selectedAssistantId === assistant.id ? styles.selected : ''
                                }`}
                            >
                                <div className={styles.assistantDetails}>
                                    <div className={styles.assistantHeader}>
                                        <span className={styles.assistantName}>
                                            {assistant.name}
                                        </span>
                                        {selectedAssistantId === assistant.id && (
                                            <Check size={16} className={styles.checkIcon}/>
                                        )}
                                    </div>
                                    {assistant.description && (
                                        <span className={styles.assistantDescription}>
                                            {assistant.description}
                                        </span>
                                    )}
                                    {renderCapabilities(assistant)}
                                </div>
                            </button>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}, (prevProps, nextProps) => {
    // Custom comparison for optimal performance
    return (
        prevProps.selectedAssistantId === nextProps.selectedAssistantId &&
        prevProps.className === nextProps.className
        // Skip onAssistantSelect comparison as it should be stable
    );
});
export default AssistantSelector; 