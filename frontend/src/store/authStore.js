import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';

/**
 * Modern Auth Store with Zustand
 * - 60% smaller than Redux
 * - Better TypeScript support
 * - Automatic persistence
 * - Immer integration for immutability
 */
export const useAuthStore = create(
  persist(
    immer((set, get) => ({
      // State
      user: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
      
      // Actions
      setUser: (user) => set((state) => {
        state.user = user;
        state.isAuthenticated = !!user;
        state.error = null;
      }),
      
      setLoading: (loading) => set((state) => {
        state.isLoading = loading;
      }),
      
      setError: (error) => set((state) => {
        state.error = error;
        state.isLoading = false;
      }),
      
      login: async (credentials) => {
        set((state) => {
          state.isLoading = true;
          state.error = null;
        });
        
        try {
          // Login logic here
          const response = await loginUser(credentials);
          set((state) => {
            state.user = response.user;
            state.isAuthenticated = true;
            state.isLoading = false;
          });
          return response;
        } catch (error) {
          set((state) => {
            state.error = error.message;
            state.isLoading = false;
          });
          throw error;
        }
      },
      
      logout: () => set((state) => {
        state.user = null;
        state.isAuthenticated = false;
        state.error = null;
      }),
      
      clearError: () => set((state) => {
        state.error = null;
      }),
      
      // Selectors (memoized)
      getUser: () => get().user,
      getIsAuthenticated: () => get().isAuthenticated,
    })),
    {
      name: 'auth-store',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);

// Optimized selectors to prevent unnecessary re-renders
export const useUser = () => useAuthStore((state) => state.user);
export const useIsAuthenticated = () => useAuthStore((state) => state.isAuthenticated);
export const useAuthLoading = () => useAuthStore((state) => state.isLoading);
export const useAuthError = () => useAuthStore((state) => state.error);

// Action selectors
export const useAuthActions = () => useAuthStore((state) => ({
  setUser: state.setUser,
  login: state.login,
  logout: state.logout,
  setLoading: state.setLoading,
  setError: state.setError,
  clearError: state.clearError,
})); 