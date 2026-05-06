// File: src/redux/slices/appModeSlice.jsx
import { createSlice } from "@reduxjs/toolkit";
/**
 * Application modes enum
 */
export const APP_MODES = {
    CLASSIC: 'classic',
    AI: 'ai'
};
/**
 * Initial state for application mode
 */
const initialState = {
    currentMode: APP_MODES.CLASSIC,
    transitions: {
        isTransitioning: false,
        previousMode: null
    }
};
/**
 * App Mode Redux Slice
 * Manages switching between Classic and AI application modes
 */
const appModeSlice = createSlice({
    name: "appMode",
    initialState,
    reducers: {
        /**
         * Switch to Classic mode
         */
        switchToClassicMode: (state) => {
            state.transitions.previousMode = state.currentMode;
            state.currentMode = APP_MODES.CLASSIC;
            // Don't clear transitioning here - let the component handle it
        },
        /**
         * Switch to AI mode
         */
        switchToAiMode: (state) => {
            state.transitions.previousMode = state.currentMode;
            state.currentMode = APP_MODES.AI;
            // Don't clear transitioning here - let the component handle it
        },
        /**
         * Set transitioning state for smooth mode switches
         */
        setModeTransitioning: (state, action) => {
            state.transitions.isTransitioning = action.payload;
        },
        /**
         * Reset app mode state to initial
         */
        resetAppModeState: (state) => {
            return initialState;
        }
    },
});
// Export actions
export const {
    switchToClassicMode,
    switchToAiMode,
    setModeTransitioning,
    resetAppModeState
} = appModeSlice.actions;
// Export selectors for easy component usage
export const selectCurrentMode = (state) => state.appMode.currentMode;
export const selectIsAiMode = (state) => state.appMode.currentMode === APP_MODES.AI;
export const selectIsClassicMode = (state) => state.appMode.currentMode === APP_MODES.CLASSIC;
export const selectIsTransitioning = (state) => state.appMode.transitions.isTransitioning;
export default appModeSlice.reducer; 