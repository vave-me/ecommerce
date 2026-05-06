"use client";
import React, { useReducer, useEffect, useRef, useCallback, useMemo, memo, useId } from "react";
import { useDispatch } from "react-redux";
import { useRouter } from "next/navigation";
import { debounce } from '../../utils/debounce';
import { Search, X, Filter } from "@/icons";
// Import extracted components
import SuggestionsList from "./SuggestionsList";
import RecentSearches from "./RecentSearches";
import MobileFiltersModal from "./MobileFiltersModal";
// Import optimized hooks
import { useListingFilters, useActiveFilters } from "./hooks/useOptimizedSelectors";
import { useIsMobile } from "../../hooks/useMobileDetection";
import { useClickOutside } from "./hooks/useEventListeners";
import { updateFilter, setFilters } from "../../redux/slices/listingFiltersSlice";
import { suggestProducts } from "../../api/searchApi";
// Location features not used in this application
import { useCategories } from "../../hooks/useCategories";
import { useSearchSuggestions } from "../../hooks/useSearch";
import styles from "./SearchBar.module.css";
// Action types for reducer
const ACTIONS = {
    SET_FIELD: 'SET_FIELD',
    SET_MULTIPLE: 'SET_MULTIPLE',
    RESET_PRODUCT_SUGGESTIONS: 'RESET_PRODUCT_SUGGESTIONS',
    RESET_ALL: 'RESET_ALL',
};

// Initial state - consolidated from all versions
const initialState = {
    // Core search
    searchQuery: "",
    selectedCategory: "",
    // Product suggestions
    productSuggestions: [],
    showProductSuggestions: false,
    productActiveIndex: -1,
    productLoading: false,
    productError: "",
    // Location features removed - not used in this application
    // UI state
    showMobileFilters: false,
    filtersActive: false,
    filterApplied: false,
    announceText: "",
    suggestionSelected: false,
    isFocused: false,
};

// Reducer function - optimized state management
function searchReducer(state, action) {
    switch (action.type) {
        case ACTIONS.SET_FIELD:
            return { ...state, [action.field]: action.value };
        case ACTIONS.SET_MULTIPLE:
            return { ...state, ...action.updates };
        case ACTIONS.RESET_PRODUCT_SUGGESTIONS:
            return {
                ...state,
                productSuggestions: [],
                showProductSuggestions: false,
                productLoading: false,
                productError: "",
                productActiveIndex: -1
            };
        // Location features removed
        case ACTIONS.RESET_ALL:
            return { ...initialState, ...action.keepValues };
        default:
            return state;
    }
}

/**
 * Unified SearchBar Component
 * Combines features from SearchBar, OptimizedSearchBar, and CatalogSearchBar
 * Configurable via props for different use cases
 */
const SearchBar = memo(({
    // Feature flags
    variant = 'header', // 'header' | 'catalog' | 'minimal'
    showCategories = true,
    showSuggestions = true,
    showMobileFilters = true,
    showClearButton = true,
    showRecentSearches = true,
    
    // Behavior
    placeholder = 'Search for products, services, posts...',
    debounceMs = 300,
    autoFocus = false,
    onSearch = null, // Custom search handler
    controlledMode = false,
    
    // Controlled mode props (for catalog variant)
    value = '',
    onChange = null,
    onClear = null,
    
    // Styling
    className = '',
    theme = 'default', // 'default' | 'blue' | 'minimal'
    
    // Other props
    disabled = false,
    isLoading = false,
    ...props
}) => {
    // Use reducer for optimized state management
    const [state, localDispatch] = useReducer(searchReducer, {
        ...initialState,
        searchQuery: controlledMode ? value : initialState.searchQuery,
    });
    // Refs
    const inputRef = useRef(null);
    const suggestionsRef = useRef(null);
    const searchFormRef = useRef(null);
    const announceRef = useRef(null);
    const debounceRef = useRef(null);
    const reduxUpdateRef = useRef(null);
    // Location not used in this application
    
    // Hooks
    const dispatch = useDispatch();
    const router = useRouter();
    const searchLabelId = useId();
    const locationLabelId = useId();
    const isMobile = useIsMobile();
    // Selectors
    const filters = useListingFilters();
    const { hasActiveFilters, activeFilterCount } = useActiveFilters();
    // IDs for accessibility
    const searchControlId = useId();
    const productSuggestionsId = useId();
    // Fetch real categories from API
    const { data: categories, isLoading: categoriesLoading } = useCategories('marketplace');
    
    // Category options from API with fallback
    const categoryOptions = useMemo(() => {
        const apiOptions = categories?.map(cat => ({
            value: cat.slug || cat.id,
            label: cat.name || cat.description || 'Unknown Category'
        })) || [];
        
        // Always include "All Categories" option first
        return [
            { value: '', label: 'All Categories' },
            ...apiOptions
        ];
    }, [categories]);
    
    // Get search value for controlled/uncontrolled mode
    const searchValue = controlledMode ? value : state.searchQuery;
    
    // Sync state with Redux filters and controlled mode
    useEffect(() => {
        if (!controlledMode && filters) {
            localDispatch({ 
                type: ACTIONS.SET_MULTIPLE, 
                updates: {
                    selectedCategory: filters.category || "",
                    searchQuery: filters.searchText || ""
                }
            });
        }
    }, [filters?.category, filters?.searchText, controlledMode]);
    
    // Sync controlled value
    useEffect(() => {
        if (controlledMode) {
            localDispatch({ type: ACTIONS.SET_FIELD, field: 'searchQuery', value });
        }
    }, [value, controlledMode]);
    
    // Load recent searches from localStorage
    const recentSearches = useMemo(() => {
        if (!showRecentSearches || typeof window === 'undefined') return [];
        try {
            const saved = localStorage.getItem('recentSearches');
            return saved ? JSON.parse(saved).slice(0, 5) : [];
        } catch (error) {
            return [];
        }
    }, [showRecentSearches]);
    // Save search to recent searches
    const saveToRecentSearches = useCallback((query) => {
        // Ensure query is a string
        const queryStr = typeof query === 'string' ? query : String(query || '');
        if (!queryStr.trim()) return;
        try {
            const current = JSON.parse(localStorage.getItem('recentSearches') || '[]');
            const updated = [queryStr, ...current.filter(item => item !== queryStr)].slice(0, 10);
            localStorage.setItem('recentSearches', JSON.stringify(updated));
            // Trigger re-render by updating state
            localDispatch({ type: ACTIONS.SET_FIELD, field: 'announceText', value: 'Search saved' });
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
    }, []);
    // Auto-focus
    useEffect(() => {
        if (autoFocus && inputRef.current) {
            inputRef.current.focus();
        }
    }, [autoFocus]);
    
    // Use new search hook for suggestions if enabled
    const { data: apiSuggestions = [], isLoading: suggestionsLoading } = useSearchSuggestions(
        searchValue,
        {
            enabled: showSuggestions && searchValue.length > 2 && !controlledMode,
            debounceMs,
            entityType: variant === 'catalog' ? 'products' : 'all',
        }
    );
    
    // Debounced search functions
    const debouncedProductSearch = useMemo(
        () => debounce(async (query) => {
            const queryStr = typeof query === 'string' ? query : String(query || '');
            if (!queryStr.trim()) {
                localDispatch({ type: ACTIONS.RESET_PRODUCT_SUGGESTIONS });
                return;
            }
            
            localDispatch({ type: ACTIONS.SET_FIELD, field: 'productLoading', value: true });
            
            try {
                const suggestions = await suggestProducts(queryStr);
                const normalizedSuggestions = Array.isArray(suggestions) 
                    ? suggestions.map(s => typeof s === 'string' ? s : (s?.name || s?.text || String(s)))
                    : [];
                    
                localDispatch({ 
                    type: ACTIONS.SET_MULTIPLE,
                    updates: {
                        productSuggestions: normalizedSuggestions,
                        productLoading: false,
                        productError: ''
                    }
                });
            } catch (error) {
                localDispatch({ 
                    type: ACTIONS.SET_MULTIPLE,
                    updates: {
                        productError: "Failed to load suggestions",
                        productSuggestions: [],
                        productLoading: false
                    }
                });
            }
        }, debounceMs),
        [debounceMs]
    );
    // Dispatch filter update
    const dispatchFilterUpdate = useCallback((key, value) => {
        dispatch(updateFilter({ key, value }));
        localDispatch({ 
            type: ACTIONS.SET_MULTIPLE,
            updates: {
                filtersActive: true,
                filterApplied: true
            }
        });
        setTimeout(() => {
            localDispatch({ type: ACTIONS.SET_FIELD, field: 'filterApplied', value: false });
        }, 2000);
    }, [dispatch]);
    // Handle search input change
    const handleProductInputChange = useCallback((e) => {
        const value = e.target.value;
        
        if (controlledMode && onChange) {
            onChange(value);
        } else {
            localDispatch({ 
                type: ACTIONS.SET_MULTIPLE,
                updates: {
                    searchQuery: value,
                    suggestionSelected: false,
                    productActiveIndex: -1
                }
            });
            
            if (value.trim() && showSuggestions) {
                debouncedProductSearch(value);
                localDispatch({ type: ACTIONS.SET_FIELD, field: 'showProductSuggestions', value: true });
            } else {
                localDispatch({ type: ACTIONS.RESET_PRODUCT_SUGGESTIONS });
            }
        }
        
        // Don't update Redux on every keystroke - only on submit
        // This prevents unnecessary API calls and re-renders
    }, [controlledMode, onChange, showSuggestions, debouncedProductSearch, dispatch]);
    // Handle category change
    const handleCategoryChange = useCallback((category) => {
        localDispatch({ 
            type: ACTIONS.SET_MULTIPLE,
            updates: {
                selectedCategory: category,
                announceText: category ? `Category changed to ${category}` : "Category cleared"
            }
        });
        dispatchFilterUpdate("category", category);
    }, [dispatchFilterUpdate]);
    // Handle search/clear for controlled mode
    const handleClear = useCallback(() => {
        if (controlledMode) {
            onChange?.('');
            onClear?.();
        } else {
            localDispatch({ type: ACTIONS.RESET_ALL });
            dispatch(setFilters({ searchText: '', category: '' }));
        }
        inputRef.current?.focus();
    }, [controlledMode, onChange, onClear, dispatch]);
    
    // Handle form submission
    const handleSubmit = useCallback((e) => {
        if (e) e.preventDefault();
        
        // Clear any pending Redux updates
        if (reduxUpdateRef.current) {
            clearTimeout(reduxUpdateRef.current);
        }
        
        const queryStr = typeof searchValue === 'string' ? searchValue : String(searchValue || '');
        
        // Custom search handler
        if (onSearch) {
            onSearch(queryStr);
            return;
        }
        
        // Default behavior
        if (queryStr.trim()) {
            saveToRecentSearches(queryStr.trim());
        }
        
        // Always update Redux on submit - this triggers the actual search
        dispatch(setFilters({ 
            searchText: queryStr.trim(),
            category: state.selectedCategory || ''
        }));
        
        localDispatch({ 
            type: ACTIONS.SET_MULTIPLE,
            updates: {
                showProductSuggestions: false,
                announceText: queryStr.trim() ? `Searching for ${queryStr}` : 'Showing all results'
            }
        });
        
        // Don't navigate to a separate search page - the feed updates via Redux
        inputRef.current?.blur(); // Remove focus to close mobile keyboard
    }, [searchValue, state, onSearch, saveToRecentSearches, dispatch, router]);
    // Keyboard navigation for suggestions
    const handleKeyDown = useCallback((e) => {
        const suggestions = state.productSuggestions;
        const showSuggestions = state.showProductSuggestions;
        const activeIndex = state.productActiveIndex;
        
        if (!showSuggestions || suggestions.length === 0) {
            if (e.key === 'Escape' && controlledMode) {
                inputRef.current?.blur();
            } else if (e.key === 'Enter') {
                handleSubmit(e);
            }
            return;
        }
        switch (e.key) {
            case 'ArrowDown':
                e.preventDefault();
                localDispatch({ 
                    type: ACTIONS.SET_FIELD, 
                    field: 'productActiveIndex', 
                    value: activeIndex < suggestions.length - 1 ? activeIndex + 1 : 0
                });
                break;
            case 'ArrowUp':
                e.preventDefault();
                localDispatch({ 
                    type: ACTIONS.SET_FIELD, 
                    field: 'productActiveIndex', 
                    value: activeIndex > 0 ? activeIndex - 1 : suggestions.length - 1
                });
                break;
            case 'Enter':
                e.preventDefault();
                if (activeIndex >= 0) {
                    handleProductSuggestionClick(suggestions[activeIndex]);
                } else {
                    // Close suggestions and submit
                    localDispatch({ type: ACTIONS.SET_FIELD, field: 'showProductSuggestions', value: false });
                    handleSubmit(e);
                }
                break;
            case 'Escape':
                localDispatch({ type: ACTIONS.SET_FIELD, field: 'showProductSuggestions', value: false });
                localDispatch({ type: ACTIONS.RESET_PRODUCT_SUGGESTIONS });
                break;
        }
    }, [state, handleSubmit, controlledMode]);
    // Handle suggestion click
    const handleProductSuggestionClick = useCallback((suggestion) => {
        // Extract text from suggestion (could be string or object)
        const suggestionText = typeof suggestion === 'string' 
            ? suggestion 
            : (suggestion?.name || suggestion?.text || String(suggestion));
        
        // Update local state
        localDispatch({ type: ACTIONS.SET_FIELD, field: 'searchQuery', value: suggestionText });
        localDispatch({ type: ACTIONS.RESET_PRODUCT_SUGGESTIONS });
        
        // Clear any pending updates
        if (reduxUpdateRef.current) {
            clearTimeout(reduxUpdateRef.current);
        }
        
        // Save to recent searches
        saveToRecentSearches(suggestionText);
        
        // Update Redux to trigger search - this is the actual search execution
        dispatch(setFilters({ 
            searchText: suggestionText,
            category: state.selectedCategory || ''
        }));
        
        // Close suggestions and blur input
        inputRef.current?.blur();
        
        localDispatch({ type: ACTIONS.SET_FIELD, field: 'announceText', value: `Searching for: ${suggestionText}` });
    }, [dispatch, saveToRecentSearches, router, state.selectedCategory]);
    // Highlight matching text in suggestions
    const highlightMatch = useCallback((text, query) => {
        if (!query) return text;
        const regex = new RegExp(`(${query})`, 'gi');
        return text.replace(regex, '<mark class="highlight">$1</mark>');
    }, []);
    // Mobile filters handlers
    const toggleMobileFilters = useCallback(() => {
        localDispatch({ type: ACTIONS.SET_FIELD, field: 'showMobileFilters', value: !state.showMobileFilters });
    }, []);
    const handleMobileClearAll = useCallback(() => {
        localDispatch({ type: ACTIONS.SET_FIELD, field: 'searchQuery', value: '' });
        localDispatch({ type: ACTIONS.SET_FIELD, field: 'selectedCategory', value: '' });
        localDispatch({ type: ACTIONS.RESET_PRODUCT_SUGGESTIONS });
        localDispatch({ type: ACTIONS.SET_FIELD, field: 'filtersActive', value: false });
        // Clear Redux filters
        dispatchFilterUpdate("searchText", "");
        dispatchFilterUpdate("category", "");
    }, [dispatchFilterUpdate]);
    // Outside click handlers
    useClickOutside(suggestionsRef, () => {
        if (state.showProductSuggestions) {
            localDispatch({ type: ACTIONS.RESET_PRODUCT_SUGGESTIONS });
        }
    });
    
    // Update focus state
    const handleFocus = useCallback(() => {
        localDispatch({ type: ACTIONS.SET_FIELD, field: 'isFocused', value: true });
        if (searchValue && state.productSuggestions.length > 0) {
            localDispatch({ type: ACTIONS.SET_FIELD, field: 'showProductSuggestions', value: true });
        }
    }, [searchValue, state.productSuggestions]);
    
    const handleBlur = useCallback(() => {
        localDispatch({ type: ACTIONS.SET_FIELD, field: 'isFocused', value: false });
    }, []);
    
    // Get container classes based on variant and state
    const containerClasses = useMemo(() => {
        const classes = [styles.container];
        if (variant === 'catalog') classes.push(styles.catalogVariant);
        if (variant === 'minimal') classes.push(styles.minimalVariant);
        if (theme === 'blue') classes.push(styles.blueTheme);
        if (state.isFocused) classes.push(styles.focused);
        if (disabled) classes.push(styles.disabled);
        if (className) classes.push(className);
        return classes.filter(Boolean).join(' ');
    }, [variant, theme, state.isFocused, disabled, className]);
    return (
        <div className={styles.container}>
            {/* Accessibility announcement region */}
            <div
                ref={announceRef}
                aria-live="polite"
                aria-atomic="true"
                className="sr-only"
            >
                {state.announceText}
            </div>
            <form
                className={`${styles.searchForm} ${state.filterApplied ? styles.feedbackPulse : ''}`}
                onSubmit={handleSubmit}
                ref={searchFormRef}
                role="search"
                aria-label="Product search for marketplace"
            >
                {/* Toast notification for filter feedback */}
                {state.filterApplied && (
                    <div className={styles.feedbackToast}>
                        Search applied successfully
                    </div>
                )}
                {/* Main search bar */}
                <div className={`${styles.searchBar} ${state.filtersActive ? styles.activeFilters : ''}`}>
                    {/* Search input section */}
                    <div className={styles.searchSection}>
                        <Search size={18} className={styles.searchIcon} />
                        <input
                            ref={inputRef}
                            id={searchControlId}
                            type="text"
                            value={searchValue}
                            onChange={handleProductInputChange}
                            onKeyDown={handleKeyDown}
                            onFocus={handleFocus}
                            onBlur={handleBlur}
                            placeholder={placeholder}
                            className={styles.searchInput}
                            aria-label="Search for products"
                            aria-autocomplete="list"
                            aria-controls={productSuggestionsId}
                            aria-expanded={state.showProductSuggestions}
                            autoComplete="off"
                        />
                        {/* Clear search button */}
                        {searchValue && showClearButton && (
                            <button
                                type="button"
                                onClick={handleClear}
                                aria-label="Clear search"
                                className={styles.clearButton}
                            >
                                <X size={16} />
                            </button>
                        )}
                    </div>
                    {/* Mobile filter button */}
                    {isMobile && showMobileFilters && variant === 'header' && (
                        <button
                            type="button"
                            onClick={toggleMobileFilters}
                            aria-label={`${state.filtersActive ? 'Edit active filters' : 'Open filters'}`}
                            className={`${styles.filterButton} ${state.filtersActive ? styles.filterActive : ''}`}
                        >
                            <Filter size={18}/>
                            {hasActiveFilters && (
                                <span className={styles.filterBadge}>
                                    {activeFilterCount}
                                </span>
                            )}
                        </button>
                    )}
                    {/* Desktop category filter */}
                    {!isMobile && showCategories && variant === 'header' && (
                        <div className={styles.categorySection}>
                            <select
                                value={state.selectedCategory}
                                onChange={(e) => handleCategoryChange(e.target.value)}
                                className={styles.categorySelect}
                                aria-label="Select category"
                                disabled={categoriesLoading || disabled}
                            >
                                {categoriesLoading ? (
                                    <option value="">Loading categories...</option>
                                ) : (
                                    categoryOptions.map(option => (
                                        <option key={option.value} value={option.value}>
                                            {option.label}
                                        </option>
                                    ))
                                )}
                            </select>
                        </div>
                    )}
                    {/* Search button removed - Enter key only search */}
                </div>
                {/* Product Suggestions Dropdown */}
                {state.showProductSuggestions && showSuggestions && (
                    <SuggestionsList
                        suggestions={state.productSuggestions}
                        query={searchValue}
                        loading={state.productLoading || isLoading}
                        error={state.productError}
                        activeIndex={state.productActiveIndex}
                        onItemClick={handleProductSuggestionClick}
                        onDismiss={() => {
                            localDispatch({ type: ACTIONS.RESET_PRODUCT_SUGGESTIONS });
                        }}
                        suggestionsRef={suggestionsRef}
                        highlightMatch={highlightMatch}
                        type="product"
                        isMobile={isMobile}
                        anchorElement={searchFormRef.current}
                    />
                )}
                {/* Enhanced Mobile Filters Modal */}
                {showMobileFilters && variant === 'header' && (
                    <MobileFiltersModal
                        showMobileFilters={state.showMobileFilters}
                        onOpenChange={(value) => localDispatch({ type: ACTIONS.SET_FIELD, field: 'showMobileFilters', value })}
                        // Category props
                        selectedCategory={state.selectedCategory}
                        onCategoryChange={handleCategoryChange}
                        categoryOptions={categoryOptions}
                        // Search props
                        searchQuery={searchValue}
                        onSearchQueryChange={(value) => {
                            if (controlledMode && onChange) {
                                onChange(value);
                            } else {
                                localDispatch({ type: ACTIONS.SET_FIELD, field: 'searchQuery', value });
                            }
                        }}
                        // Actions
                        onSubmit={handleSubmit}
                        onClearAll={handleMobileClearAll}
                    />
                )}
            </form>
        </div>
    );
});
SearchBar.displayName = 'SearchBar';
export default SearchBar;