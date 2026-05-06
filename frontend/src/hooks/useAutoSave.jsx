"use client";
// CreatePostModal/hooks/useAutosave.js
import { useCallback, useEffect, useRef, useState } from 'react';
import { debounce } from "../utils/debounce.js";
import { UnifiedUtils } from "../utils/duplicateEliminator";
/**
 * Enhanced hook for automatically saving post drafts
 * @param {string} userId - Current user ID
 * @param {Object} postData - Data to be saved
 * @param {string} postId - ID of the post (optional for new posts)
 * @param {Object} options - Configuration options
 * @returns {Object} Last saved timestamp and saving status
 */
export function useAutoSave(userId, postData, postId, options = {}) {
    const {
        saveInterval = 30000,  // 30 seconds
        debounceDelay = 1000,  // 1 second
        shouldSaveEmpty = false,
    } = options;
    const [lastSaved, setLastSaved] = useState(null);
    const [isSaving, setIsSaving] = useState(false);
    const timerRef = useRef(null);
    // Check if content exists and should trigger a save
    const hasContent = shouldSaveEmpty || 
        Boolean(postData?.name?.trim() || postData?.description?.trim());
    // Create a memoized save function that won't unnecessarily re-render components
    const saveData = useCallback(async () => {
        // Don't save if user is not logged in or no content to save
        if (!userId || !hasContent) {
            return;
        }
        try {
            setIsSaving(true);
            // Prepare data for saving
            const tagsArray = postData.tags
                ? postData.tags.split(",").map(t => t.trim()).filter(Boolean)
                : [];
            const draftData = {
                id: postId,
                userId,
                name: postData.name || '',
                description: postData.description || '',
                tags: tagsArray,
                status: "draft",
                thumbnail: postData.thumbnail || "",
                lat: postData.lat || 0,
                lng: postData.lng || 0,
                lastModified: new Date().toISOString()
            };
            // Save to localStorage as backup
            await UnifiedUtils.drafts.saveDraft('auto', userId, draftData);
            // Update last saved timestamp
            setLastSaved(new Date());
        } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    } finally {
            setIsSaving(false);
        }
    }, [userId, postId, postData, hasContent]);
    // Create stable reference to debounced save function
    const debouncedSave = useRef(
        debounce(saveData, debounceDelay)
    ).current;
    // Set up auto-save effect
    useEffect(() => {
        // Cleanup any existing timer
        if (timerRef.current) {
            clearInterval(timerRef.current);
            timerRef.current = null;
        }
        // Only set up autosave if we have content to save
        if (hasContent && userId) {
            // Trigger initial save with debounce to avoid saving empty data
            debouncedSave();
            // Set up interval for periodic saves
            timerRef.current = setInterval(() => {
                if (!isSaving && hasContent) {
                    saveData();
                }
            }, saveInterval);
        }
        // Cleanup on unmount or deps change
        return () => {
            if (timerRef.current) {
                clearInterval(timerRef.current);
                timerRef.current = null;
            }
        };
    }, [userId, hasContent, saveInterval, debouncedSave, saveData, isSaving]);
    // Force a save (non-debounced) - useful for save buttons
    const forceSave = useCallback(() => {
        if (hasContent && !isSaving) {
            saveData();
        }
    }, [hasContent, isSaving, saveData]);
    return {
        lastSaved,
        isSaving,
        forceSave
    };
}