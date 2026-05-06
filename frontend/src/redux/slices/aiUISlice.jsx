import { createSlice } from '@reduxjs/toolkit';

/**
 * AI UI State Slice
 * Manages UI-only state for AI features
 * Server state is handled by React Query
 */

const initialState = {
    // UI selections
    selectedAssistantId: null,
    activeConversationId: null,
    
    // UI states
    isTyping: false,
    viewMode: 'chat', // 'chat' | 'split' | 'results'
    sidebarCollapsed: false,
    
    // UI preferences
    showTimestamps: true,
    showAvatars: true,
    compactMode: false,
    
    // Temporary UI states
    isCreatingConversation: false,
    isSwitchingAssistant: false,
    
    // Error states for UI
    lastError: null,
    showErrorBanner: false
};

const aiUISlice = createSlice({
    name: 'aiUI',
    initialState,
    reducers: {
        // Assistant selection
        setSelectedAssistant: (state, action) => {
            state.selectedAssistantId = action.payload;
            // Clear conversation when switching assistants
            state.activeConversationId = null;
            state.isSwitchingAssistant = false;
        },
        
        startAssistantSwitch: (state) => {
            state.isSwitchingAssistant = true;
        },
        
        // Conversation management
        setActiveConversation: (state, action) => {
            state.activeConversationId = action.payload;
            state.isCreatingConversation = false;
        },
        
        startConversationCreation: (state) => {
            state.isCreatingConversation = true;
        },
        
        clearActiveConversation: (state) => {
            state.activeConversationId = null;
        },
        
        // Typing indicator
        setTyping: (state, action) => {
            state.isTyping = action.payload;
        },
        
        // View mode
        setViewMode: (state, action) => {
            if (['chat', 'split', 'results'].includes(action.payload)) {
                state.viewMode = action.payload;
            }
        },
        
        toggleViewMode: (state) => {
            const modes = ['chat', 'split', 'results'];
            const currentIndex = modes.indexOf(state.viewMode);
            state.viewMode = modes[(currentIndex + 1) % modes.length];
        },
        
        // Sidebar
        toggleSidebar: (state) => {
            state.sidebarCollapsed = !state.sidebarCollapsed;
        },
        
        setSidebarCollapsed: (state, action) => {
            state.sidebarCollapsed = action.payload;
        },
        
        // UI preferences
        updateUIPreferences: (state, action) => {
            const { showTimestamps, showAvatars, compactMode } = action.payload;
            if (showTimestamps !== undefined) state.showTimestamps = showTimestamps;
            if (showAvatars !== undefined) state.showAvatars = showAvatars;
            if (compactMode !== undefined) state.compactMode = compactMode;
        },
        
        // Error handling
        setUIError: (state, action) => {
            state.lastError = action.payload;
            state.showErrorBanner = true;
        },
        
        clearUIError: (state) => {
            state.lastError = null;
            state.showErrorBanner = false;
        },
        
        // Reset UI state
        resetAIUIState: () => initialState
    }
});

// Export actions
export const {
    setSelectedAssistant,
    startAssistantSwitch,
    setActiveConversation,
    startConversationCreation,
    clearActiveConversation,
    setTyping,
    setViewMode,
    toggleViewMode,
    toggleSidebar,
    setSidebarCollapsed,
    updateUIPreferences,
    setUIError,
    clearUIError,
    resetAIUIState
} = aiUISlice.actions;

// Export selectors
export const selectSelectedAssistantId = (state) => state.aiUI.selectedAssistantId;
export const selectActiveConversationId = (state) => state.aiUI.activeConversationId;
export const selectIsTyping = (state) => state.aiUI.isTyping;
export const selectViewMode = (state) => state.aiUI.viewMode;
export const selectSidebarCollapsed = (state) => state.aiUI.sidebarCollapsed;
export const selectUIPreferences = (state) => ({
    showTimestamps: state.aiUI.showTimestamps,
    showAvatars: state.aiUI.showAvatars,
    compactMode: state.aiUI.compactMode
});
export const selectIsCreatingConversation = (state) => state.aiUI.isCreatingConversation;
export const selectIsSwitchingAssistant = (state) => state.aiUI.isSwitchingAssistant;
export const selectUIError = (state) => state.aiUI.lastError;
export const selectShowErrorBanner = (state) => state.aiUI.showErrorBanner;

// Export reducer
export default aiUISlice.reducer;