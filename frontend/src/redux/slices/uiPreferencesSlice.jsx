import { createSlice } from "@reduxjs/toolkit";

// Load preferences from localStorage if available
const loadPreferences = () => {
  if (typeof window === 'undefined') {
    return null; // Return null on server
  }
  
  try {
    const stored = localStorage.getItem('uiPreferences');
    return stored ? JSON.parse(stored) : null;
  } catch (error) {
    // Error: 'Failed to load UI preferences:', error...
    return null;
  }
};

// Save preferences to localStorage
const savePreferences = (preferences) => {
  if (typeof window === 'undefined') {
    return;
  }
  
  try {
    localStorage.setItem('uiPreferences', JSON.stringify(preferences));
  } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
};

// Default state values
const defaultState = {
  // Composer visibility preferences
  showUnifiedComposer: true,
  showComposerOnHomepage: true,
  showComposerOnMarketplace: true,
  showComposerOnProducts: true,
  
  // Other UI preferences
  compactMode: false,
  autoPlayVideos: true,
  showTrendingSection: true,
  defaultFeedView: 'grid', // grid, list, masonry
  
  // Accessibility
  reduceAnimations: false,
  highContrast: false,
  
  // AI preferences
  autoEnableAIMode: false,
  showAIResponses: true,
};

// Initialize with default values, will be updated on client
const initialState = defaultState;

const uiPreferencesSlice = createSlice({
  name: "uiPreferences",
  initialState,
  reducers: {
    // Toggle unified composer globally
    toggleUnifiedComposer: (state) => {
      state.showUnifiedComposer = !state.showUnifiedComposer;
      savePreferences(state);
    },
    
    // Hide unified composer globally
    hideUnifiedComposer: (state) => {
      state.showUnifiedComposer = false;
      savePreferences(state);
    },
    
    // Show unified composer globally
    showUnifiedComposer: (state) => {
      state.showUnifiedComposer = true;
      savePreferences(state);
    },
    
    // Toggle composer on specific pages
    toggleComposerOnPage: (state, action) => {
      const { page } = action.payload;
      switch (page) {
        case 'homepage':
          state.showComposerOnHomepage = !state.showComposerOnHomepage;
          break;
        case 'marketplace':
          state.showComposerOnMarketplace = !state.showComposerOnMarketplace;
          break;
        case 'products':
          state.showComposerOnProducts = !state.showComposerOnProducts;
          break;
        default:
          break;
      }
      savePreferences(state);
    },
    
    // Set composer visibility for specific page
    setComposerVisibilityForPage: (state, action) => {
      const { page, visible } = action.payload;
      switch (page) {
        case 'homepage':
          state.showComposerOnHomepage = visible;
          break;
        case 'marketplace':
          state.showComposerOnMarketplace = visible;
          break;
        case 'products':
          state.showComposerOnProducts = visible;
          break;
        default:
          break;
      }
      savePreferences(state);
    },
    
    // Update multiple preferences at once
    updatePreferences: (state, action) => {
      Object.assign(state, action.payload);
      savePreferences(state);
    },
    
    // Reset all preferences to defaults
    resetPreferences: (state) => {
      Object.assign(state, {
        showUnifiedComposer: true,
        showComposerOnHomepage: true,
        showComposerOnMarketplace: true,
        showComposerOnProducts: true,
        compactMode: false,
        autoPlayVideos: true,
        showTrendingSection: true,
        defaultFeedView: 'grid',
        reduceAnimations: false,
        highContrast: false,
        autoEnableAIMode: false,
        showAIResponses: true,
      });
      savePreferences(state);
    },
    
    // Toggle compact mode
    toggleCompactMode: (state) => {
      state.compactMode = !state.compactMode;
      savePreferences(state);
    },
    
    // Toggle auto-play videos
    toggleAutoPlayVideos: (state) => {
      state.autoPlayVideos = !state.autoPlayVideos;
      savePreferences(state);
    },
    
    // Set default feed view
    setDefaultFeedView: (state, action) => {
      state.defaultFeedView = action.payload;
      savePreferences(state);
    },
    
    // Toggle AI preferences
    toggleAutoEnableAIMode: (state) => {
      state.autoEnableAIMode = !state.autoEnableAIMode;
      savePreferences(state);
    },
    
    toggleShowAIResponses: (state) => {
      state.showAIResponses = !state.showAIResponses;
      savePreferences(state);
    },
    
    // Hydrate state from localStorage (client-side only)
    hydrateFromLocalStorage: (state) => {
      const saved = loadPreferences();
      if (saved) {
        return { ...state, ...saved };
      }
      return state;
    },
  },
});

export const {
  toggleUnifiedComposer,
  hideUnifiedComposer,
  showUnifiedComposer,
  toggleComposerOnPage,
  setComposerVisibilityForPage,
  updatePreferences,
  resetPreferences,
  toggleCompactMode,
  toggleAutoPlayVideos,
  setDefaultFeedView,
  toggleAutoEnableAIMode,
  toggleShowAIResponses,
  hydrateFromLocalStorage,
} = uiPreferencesSlice.actions;

// Selectors
export const selectShowUnifiedComposer = (state) => state.uiPreferences?.showUnifiedComposer ?? false;
export const selectShowComposerOnHomepage = (state) => state.uiPreferences.showComposerOnHomepage;
export const selectShowComposerOnMarketplace = (state) => state.uiPreferences.showComposerOnMarketplace;
export const selectShowComposerOnProducts = (state) => state.uiPreferences.showComposerOnProducts;
export const selectCompactMode = (state) => state.uiPreferences.compactMode;
export const selectAutoPlayVideos = (state) => state.uiPreferences.autoPlayVideos;
export const selectShowTrendingSection = (state) => state.uiPreferences.showTrendingSection;
export const selectDefaultFeedView = (state) => state.uiPreferences.defaultFeedView;
export const selectAutoEnableAIMode = (state) => state.uiPreferences.autoEnableAIMode;
export const selectShowAIResponses = (state) => state.uiPreferences.showAIResponses;

export default uiPreferencesSlice.reducer;