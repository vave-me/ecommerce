"use client";
import React, {useState, useCallback, useMemo, memo, useReducer, useEffect} from 'react';
import {useQuery} from '@tanstack/react-query';
import {useSelector, useDispatch} from 'react-redux';
import {useRouter, useSearchParams} from 'next/navigation';
import {Grid3X3, List, Map, Filter, Search, SlidersHorizontal, RefreshCw, Play, ChevronDown} from '@/icons';
import { useIsMobile } from '../../../hooks/useMobileDetection';
import useMasonryLayout from '../../../hooks/useMasonryLayout';
import {ServiceCard} from '../../../components/services';

import {getUnifiedCatalog} from '../../../api/client/searchApi';
import LoadingPlaceholder from '../../../components/Utils/LoadingPlaceholder';
import ErrorPlaceholder from '../../../components/Utils/ErrorPlaceholder';
import EmptyPlaceholder from '../../../components/Utils/EmptyPlaceholder';
// Import entity-specific cards
import {ClassifiedCard} from '../../../components/classified';
import ImprovedPostCard from '../../../components/ImprovedPostCard';
import CatalogsFilters from '../../../components/Listings/ListingsFilters';
import CatalogsPagination from '../../../components/Listings/ListingsPagination';
import SearchBar from '../../../components/Header/SearchBar';
import styles from './CatalogPage.module.css';
import {setFilters} from "@/redux/slices/listingFiltersSlice";
// Constants for better maintainability
const ENTITY_TYPES = {
    PRODUCT: 'product',
    SERVICE: 'service',
    POST: 'post'
};
const VIEW_MODES = {
    GRID: 'grid',
    LIST: 'list',
    MAP: 'map'
};
const DEFAULT_FILTERS = {
    contentType: 'all',
    displayMode: VIEW_MODES.GRID,
    page: 1,
    limit: 20,
    sortBy: 'createdAt',
    sortOrder: 'desc'
};
/**
 * Transform and normalize unified catalog response
 * Following Header's data processing patterns
 */
const normalizeUnifiedCatalogResponse = (data) => {
    if (!data) {
        return {catalogs: [], totalCount: 0, totalPages: 0, currentPage: 1};
    }
    let catalogs = [];
    if (data.items && Array.isArray(data.items)) {
        catalogs = data.items;
    } else if (Array.isArray(data)) {
        catalogs = data;
    } else if (data.data && Array.isArray(data.data)) {
        catalogs = data.data;
    }
    // Ensure proper entity type assignment and data consistency
    catalogs = catalogs.map((item, index) => ({
        ...item,
        entityType: item.entityType || item.type || 'unknown',
        type: item.entityType || item.type || 'unknown',
        createdAt: item.createdAt || item.added || item.created || new Date().toISOString(),
        id: item.id || `${item.entityType || 'item'}-${Math.random().toString(36).substr(2, 9)}`
    }));
    return {
        catalogs,
        totalCount: data.totalCount || data.total || catalogs.length,
        totalPages: data.totalPages || Math.ceil((data.totalCount || data.total || catalogs.length) / 20),
        currentPage: data.currentPage || data.page || 1
    };
};
/**
 * Catalog state reducer - Following Header's state management patterns
 */
const catalogReducer = (state, action) => {
    switch (action.type) {
        case 'SET_LOAD_REQUESTED':
            return {...state, loadCatalogRequested: action.payload};
        case 'SET_USER_INITIATED':
            return {...state, userInitiatedLoad: action.payload};
        case 'TOGGLE_FILTERS':
            return {...state, showFilters: !state.showFilters};
        case 'SET_FILTERS_VISIBILITY':
            return {...state, showFilters: action.payload};
        case 'RESET_STATE':
            return {loadCatalogRequested: false, userInitiatedLoad: false, showFilters: false};
        default:
            return state;
    }
};
/**
 * CatalogPageClient - Refactored to match Header's architecture
 * Clean, performance-optimized, and maintainable
 */
const CatalogPageClient = memo(function CatalogPageClient({
    userId,
    searchParams,
    locale = 'en',
    labels = {}
}) {
    // Hooks - Following Header's hook organization
    const isMobile = useIsMobile();
    const dispatch = useDispatch();
    const router = useRouter();
    const currentSearchParams = useSearchParams();
    // Redux state
    const catalogFilters = useSelector((state) => state.listingsFilters);
    // Local state with reducer pattern (Header-inspired)
    const [catalogState, catalogDispatch] = useReducer(catalogReducer, {
        loadCatalogRequested: false,
        userInitiatedLoad: false,
        showFilters: false
    });
    // Auto-trigger catalog loading - Performance optimized
    useEffect(() => {
        if (userId && !catalogState.loadCatalogRequested) {
            catalogDispatch({type: 'SET_LOAD_REQUESTED', payload: true});
            catalogDispatch({type: 'SET_USER_INITIATED', payload: true});
        }
    }, [userId, catalogState.loadCatalogRequested]);
    // Process search params - searchParams is already resolved by server component
    // No need to use React.use() since it's passed as a regular prop
    const unifiedFilters = useMemo(() => ({
        userId: userId,
        ...DEFAULT_FILTERS,
        contentType: searchParams?.type || DEFAULT_FILTERS.contentType,
        displayMode: searchParams?.displayMode || DEFAULT_FILTERS.displayMode,
        page: parseInt(searchParams?.page || '1', 10),
        limit: parseInt(searchParams?.limit || '20', 10),
        searchText: searchParams?.q,
        category: searchParams?.category,
        location: searchParams?.location,
        minPrice: searchParams?.minPrice ? parseFloat(searchParams.minPrice) : undefined,
        maxPrice: searchParams?.maxPrice ? parseFloat(searchParams.maxPrice) : undefined,
        sortBy: searchParams?.sortBy || DEFAULT_FILTERS.sortBy,
        sortOrder: searchParams?.sortOrder || DEFAULT_FILTERS.sortOrder,
        locale
    }), [userId, searchParams, locale]);
    const currentFilters = useMemo(() => ({
        ...unifiedFilters,
        ...catalogFilters,
    }), [unifiedFilters, catalogFilters]);
    // Helper functions - Following Header's utility patterns
    const cleanParams = useCallback((params) => {
        const cleaned = {};
        for (const key in params) {
            const value = params[key];
            if (value !== null && typeof value !== 'undefined' && value !== '') {
                cleaned[key] = value;
            }
        }
        return cleaned;
    }, []);
    // Security filter function - Enhanced for better security
    const applySecurityFilter = useCallback((items, expectedUserId) => {
        if (!Array.isArray(items)) return [];
        return items.filter(item => {
            const itemUserId = item.post?.userId ||
                item.vehicle?.userId ||
                item.job?.userId ||
                item.service?.userId ||
                item.property?.userId ||
                item.deal?.userId ||
                item.userId ||
                item.ownerId ||
                item.createdBy ||
                item.authorId;
            return itemUserId === expectedUserId;
        });
    }, []);
    // Catalog query - Optimized with better error handling
    const {
        data: catalogData,
        isLoading,
        isError,
        error,
        refetch,
        isFetching
    } = useQuery({
        queryKey: ['userCatalog', userId, cleanParams(currentFilters)],
        queryFn: async () => {
            const cleanedParams = cleanParams(currentFilters);
            const response = await getUnifiedCatalog(userId, cleanedParams);
            // Apply client-side security filtering
            if (response?.items && Array.isArray(response.items)) {
                response.items = applySecurityFilter(response.items, userId);
            }
            return normalizeUnifiedCatalogResponse(response);
        },
        enabled: catalogState.loadCatalogRequested && !!userId,
        retry: (failureCount, error) => {
            if (error?.response?.status === 401 || error?.response?.status === 403) {
                return false;
            }
            return failureCount < 2;
        },
        staleTime: 1000 * 60 * 5, // 5 minutes
        cacheTime: 1000 * 60 * 10, // 10 minutes
    });
    // Event handlers - Following Header's callback patterns
    const handleLoadCatalog = useCallback(() => {
        catalogDispatch({type: 'SET_LOAD_REQUESTED', payload: true});
        catalogDispatch({type: 'SET_USER_INITIATED', payload: true});
    }, []);
    const handleFilterChange = useCallback((newFilters) => {
        dispatch(setFilters(newFilters));
        // Update URL params
        const params = new URLSearchParams(currentSearchParams);
        Object.entries(newFilters).forEach(([key, value]) => {
            if (value !== undefined && value !== null && value !== '') {
                params.set(key, value);
            } else {
                params.delete(key);
            }
        });
        router.push(`/${locale}/catalog?${params.toString()}`);
        // Refetch if catalog was already loaded
        if (catalogState.loadCatalogRequested) {
            setTimeout(() => refetch(), 100);
        }
    }, [dispatch, router, locale, currentSearchParams, catalogState.loadCatalogRequested, refetch]);
    const handleDisplayModeChange = useCallback((mode) => {
        handleFilterChange({displayMode: mode});
    }, [handleFilterChange]);
    const handleSearch = useCallback((searchText) => {
        handleFilterChange({searchText, page: 1});
    }, [handleFilterChange]);
    const handlePageChange = useCallback((page) => {
        handleFilterChange({page});
        window.scrollTo({top: 0, behavior: 'smooth'});
    }, [handleFilterChange]);
    const toggleFilters = useCallback(() => {
        catalogDispatch({type: 'TOGGLE_FILTERS'});
    }, []);
    const handleRefresh = useCallback(() => {
        refetch();
    }, [refetch]);
    // Entity card renderer - Optimized with better error handling
    const renderEntityCard = useCallback((item, index) => {
        const entityType = item.entityType || item.type;
        const uniqueKey = `catalog-${entityType}-${item.id || index}`;
        try {
            const cardProps = {
                key: uniqueKey,
                className: styles.entityCard
            };
            switch (entityType) {
                case ENTITY_TYPES.PRODUCT:
                    return <div {...cardProps}><ClassifiedCard product={item.product || item}/></div>;
                case ENTITY_TYPES.SERVICE:
                    return <div {...cardProps}><ServiceCard service={item.service || item}/></div>;
                case ENTITY_TYPES.POST:
                    return <div {...cardProps}><ImprovedPostCard post={item}/></div>;
                default:
                    return (
                        <div {...cardProps}>
                            <div className={styles.unsupportedType}>
                                <p>Unsupported content type: {entityType}</p>
                            </div>
                        </div>
                    );
            }
        } catch (error) {
            return (
                <div key={uniqueKey} className={styles.entityCard}>
                    <div className={styles.cardError}>
                        <p>Error rendering {entityType}</p>
                    </div>
                </div>
            );
        }
    }, []);
    // Process catalog data - Optimized with proper sorting
    const processedCatalogData = useMemo(() => {
        if (!catalogData) return {catalogs: [], totalCount: 0, totalPages: 0, currentPage: 1};
        let catalogs = [...catalogData.catalogs];
        // Apply client-side filtering
        if (currentFilters.contentType && currentFilters.contentType !== 'all') {
            catalogs = catalogs.filter(item =>
                item.entityType === currentFilters.contentType ||
                item.type === currentFilters.contentType
            );
        }
        // Apply sorting
        if (currentFilters.sortBy && currentFilters.sortOrder) {
            catalogs.sort((a, b) => {
                const aValue = a[currentFilters.sortBy] || '';
                const bValue = b[currentFilters.sortBy] || '';
                if (currentFilters.sortBy === 'createdAt') {
                    const dateA = new Date(aValue);
                    const dateB = new Date(bValue);
                    return currentFilters.sortOrder === 'asc' ? dateA - dateB : dateB - dateA;
                }
                return currentFilters.sortOrder === 'asc'
                    ? aValue.localeCompare(bValue)
                    : bValue.localeCompare(aValue);
            });
        }
        return {
            catalogs,
            totalCount: catalogs.length,
            totalPages: Math.ceil(catalogs.length / (currentFilters.limit || 20)),
            currentPage: currentFilters.page || 1
        };
    }, [catalogData, currentFilters]);
    // Render states - Following Header's rendering patterns
    if (!catalogState.loadCatalogRequested) {
        return (
            <div className={styles.container}>
                <header className={styles.header}>
                    <div className={styles.headerContent}>
                        <div className={styles.titleSection}>
                            <h1 className={styles.title}>{labels.pageTitle}</h1>
                            <p className={styles.subtitle}>Your catalog is ready to load</p>
                        </div>
                    </div>
                </header>
                <div className={styles.loadPrompt}>
                    <div className={styles.loadPromptContent}>
                        <Play className={styles.loadIcon}/>
                        <h2>Load My Catalog</h2>
                        <p>Click below to view your personalized catalog</p>
                        <button
                            onClick={handleLoadCatalog}
                            className={styles.loadButton}
                            type="button"
                        >
                            {labels.loadCatalog || 'Load Catalog'}
                        </button>
                    </div>
                </div>
            </div>
        );
    }
    if (isLoading && catalogState.userInitiatedLoad) {
        return (
            <div className={styles.container}>
                <header className={styles.header}>
                    <div className={styles.headerContent}>
                        <div className={styles.titleSection}>
                            <h1 className={styles.title}>{labels.pageTitle}</h1>
                        </div>
                    </div>
                </header>
                <LoadingPlaceholder message={labels.loading}/>
            </div>
        );
    }
    if (isError) {
        return (
            <div className={styles.container}>
                <header className={styles.header}>
                    <div className={styles.headerContent}>
                        <div className={styles.titleSection}>
                            <h1 className={styles.title}>{labels.pageTitle}</h1>
                        </div>
                    </div>
                </header>
                <ErrorPlaceholder
                    message={error?.message || labels.errorFetching}
                    onRetry={handleRefresh}
                />
            </div>
        );
    }
    const {catalogs, totalCount, totalPages, currentPage} = processedCatalogData;
    const displayMode = currentFilters.displayMode || VIEW_MODES.GRID;
    
    // Use masonry layout for desktop grid view only
    const masonryRef = useMasonryLayout([catalogs, !isMobile && displayMode === VIEW_MODES.GRID]);
    
    return (
        <div className={styles.container}>
            {/* Header - Matching Header component patterns with enhanced accessibility */}
            <header className={styles.header} role="banner">
                <div className={styles.headerContent}>
                    <div className={styles.topRow}>
                        <div className={styles.titleSection}>
                            <h1 className={styles.title} id="catalog-title">{labels.pageTitle}</h1>
                            <p className={styles.subtitle} aria-live="polite" aria-atomic="true">
                                {totalCount} {totalCount === 1 ? 'item' : 'items'} in your catalog
                            </p>
                        </div>
                        <div className={styles.headerControls} role="group" aria-label="Catalog controls">
                            <button
                                onClick={handleRefresh}
                                disabled={isFetching}
                                className={styles.refreshButton}
                                title="Refresh catalog"
                                aria-label={isFetching ? "Refreshing catalog..." : "Refresh catalog"}
                                aria-describedby="refresh-status"
                            >
                                <RefreshCw className={isFetching ? styles.spinning : ''} aria-hidden="true"/>
                                <span id="refresh-status" className="sr-only">
                                    {isFetching ? "Refreshing catalog..." : ""}
                                </span>
                            </button>
                            {!isMobile && (
                                <div 
                                    className={styles.viewModeButtons} 
                                    role="group" 
                                    aria-label="View mode selection"
                                >
                                    <button
                                        onClick={() => handleDisplayModeChange(VIEW_MODES.GRID)}
                                        className={`${styles.viewModeButton} ${displayMode === VIEW_MODES.GRID ? styles.active : ''}`}
                                        title={labels.viewModes?.grid || "Grid view"}
                                        aria-label="Grid view"
                                        aria-pressed={displayMode === VIEW_MODES.GRID}
                                    >
                                        <Grid3X3 aria-hidden="true"/>
                                    </button>
                                    <button
                                        onClick={() => handleDisplayModeChange(VIEW_MODES.LIST)}
                                        className={`${styles.viewModeButton} ${displayMode === VIEW_MODES.LIST ? styles.active : ''}`}
                                        title={labels.viewModes?.list || "List view"}
                                        aria-label="List view"
                                        aria-pressed={displayMode === VIEW_MODES.LIST}
                                    >
                                        <List aria-hidden="true"/>
                                    </button>
                                </div>
                            )}
                            <button
                                onClick={toggleFilters}
                                className={`${styles.filtersButton} ${catalogState.showFilters ? styles.active : ''}`}
                                aria-label={catalogState.showFilters ? "Hide filters" : "Show filters"}
                                aria-expanded={catalogState.showFilters}
                                aria-controls="filters-panel"
                            >
                                <SlidersHorizontal aria-hidden="true"/>
                                {labels.filterBy}
                            </button>
                        </div>
                    </div>
                    {/* Search Section - Using sophisticated SearchBar component */}
                    <div className={styles.searchSection}>
                        <SearchBar
                            variant="catalog"
                            value={currentFilters.searchText || ''}
                            placeholder={labels.searchPlaceholder}
                            onChange={handleSearch}
                            controlledMode={true}
                            showCategories={false}
                            showMobileFilters={false}
                            theme="blue"
                            onClear={() => handleSearch('')}
                            debounceMs={300}
                            disabled={isLoading}
                            className={styles.searchContainer}
                        />
                    </div>
                </div>
            </header>
            {/* Filters Panel - Optimized visibility with accessibility */}
            {catalogState.showFilters && (
                <div 
                    className={styles.filtersContainer}
                    id="filters-panel"
                    role="region"
                    aria-labelledby="filters-title"
                >
                    <h2 id="filters-title" className="sr-only">Catalog Filters</h2>
                    <CatalogsFilters
                        filters={currentFilters}
                        onFilterChange={handleFilterChange}
                        labels={labels}
                        className={styles.filtersPanel}
                    />
                </div>
            )}
            {/* Main Content */}
            <main className={styles.main} role="main">
                {catalogs.length === 0 ? (
                    <div className={styles.emptyState}>
                        <h3>Your catalog is empty</h3>
                        <p>You haven't created any listings yet. Start by adding your first item!</p>
                        <div className={styles.emptyActions}>
                            <button 
                                onClick={handleRefresh} 
                                className={styles.refreshCatalogButton}
                                aria-label="Refresh catalog"
                            >
                                Refresh Catalog
                            </button>
                            <a href="/create" className={styles.createFirstItemButton}>
                                Create Your First Item
                            </a>
                        </div>
                    </div>
                ) : (
                    <>
                        <div 
                            className={`${styles.catalogGrid} ${styles[displayMode]}`}
                            ref={(!isMobile && displayMode === VIEW_MODES.GRID) ? masonryRef : null}
                        >
                            {catalogs.map((item, index) => renderEntityCard(item, index))}
                        </div>
                        {totalPages > 1 && (
                            <CatalogsPagination
                                currentPage={currentPage}
                                totalPages={totalPages}
                                totalCount={totalCount}
                                onPageChange={handlePageChange}
                                className={styles.pagination}
                            />
                        )}
                    </>
                )}
            </main>
            {/* Loading overlay - Performance optimized */}
            {isFetching && !catalogState.userInitiatedLoad && (
                <div className={styles.loadingOverlay}>
                    <RefreshCw className={styles.spinning}/>
                </div>
            )}
        </div>
    );
});
export default CatalogPageClient; 