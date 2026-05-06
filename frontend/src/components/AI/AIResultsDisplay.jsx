"use client";
import React, { useMemo, useState, useCallback } from 'react';
import { 
    Grid3X3, 
    List, 
    SlidersHorizontal, 
    TrendingUp,
    Package,
    FileText,
    Wrench,
    Filter,
    ChevronDown,
    Sparkles,
    Bot
} from '@/icons';
import ClassifiedCard from '../classified/ClassifiedCard';
import ServiceCard from '../services/ServiceCard';
import PostCard from '../PostCard/PostCard';
import styles from './AIResultsDisplay.module.css';

/**
 * AIResultsDisplay - Renders assistant search results as cards
 * High-quality component that seamlessly integrates AI results with existing card components
 */
const AIResultsDisplay = ({ 
    assistantResponse,
    displayMode = 'grid', // 'grid' | 'list' | 'cards'
    onItemClick,
    loading = false,
    className = ''
}) => {
    const [currentDisplayMode, setCurrentDisplayMode] = useState(displayMode);
    const [sortBy, setSortBy] = useState('relevance');
    const [filterBy, setFilterBy] = useState('all');

    // Extract items from assistant response
    const items = useMemo(() => {
        if (!assistantResponse?.Data?.items) return [];
        return assistantResponse.Data.items;
    }, [assistantResponse]);

    // Get metadata
    const metadata = useMemo(() => {
        return assistantResponse?.Data?.metadata || {};
    }, [assistantResponse]);

    // Group items by entity type
    const groupedItems = useMemo(() => {
        const groups = {};
        items.forEach(item => {
            const type = item.entityType || 'unknown';
            if (!groups[type]) groups[type] = [];
            groups[type].push(item);
        });
        return groups;
    }, [items]);

    // Sort items based on selected criteria
    const sortedItems = useMemo(() => {
        const itemsCopy = [...items];
        
        switch (sortBy) {
            case 'relevance':
                return itemsCopy.sort((a, b) => (b.relevanceScore || 0) - (a.relevanceScore || 0));
            case 'newest':
                return itemsCopy.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));
            case 'price-low':
                return itemsCopy.sort((a, b) => {
                    const priceA = a.product?.basePrice || a.service?.price || 0;
                    const priceB = b.product?.basePrice || b.service?.price || 0;
                    return priceA - priceB;
                });
            case 'price-high':
                return itemsCopy.sort((a, b) => {
                    const priceA = a.product?.basePrice || a.service?.price || 0;
                    const priceB = b.product?.basePrice || b.service?.price || 0;
                    return priceB - priceA;
                });
            default:
                return itemsCopy;
        }
    }, [items, sortBy]);

    // Filter items
    const filteredItems = useMemo(() => {
        if (filterBy === 'all') return sortedItems;
        return sortedItems.filter(item => item.entityType === filterBy);
    }, [sortedItems, filterBy]);

    // Render individual item based on type
    const renderItem = useCallback((item) => {
        const key = item.id || `${item.entityType}_${Math.random()}`;
        
        // Add AI source indicator
        const aiEnhancedProps = {
            className: styles.aiSourcedCard,
            'data-ai-source': true,
            'data-relevance': item.relevanceScore
        };

        switch (item.entityType) {
            case 'product':
                return (
                    <div key={key} className={styles.cardWrapper}>
                        <div className={styles.aiIndicator}>
                            <Bot size={14} />
                            <span>AI Match</span>
                        </div>
                        <ClassifiedCard 
                            product={item.product} 
                            {...aiEnhancedProps}
                            onClick={() => onItemClick?.(item)}
                        />
                    </div>
                );
            
            case 'service':
                return (
                    <div key={key} className={styles.cardWrapper}>
                        <div className={styles.aiIndicator}>
                            <Bot size={14} />
                            <span>AI Match</span>
                        </div>
                        <ServiceCard 
                            service={item.service} 
                            {...aiEnhancedProps}
                            onClick={() => onItemClick?.(item)}
                        />
                    </div>
                );
            
            case 'post':
                return (
                    <div key={key} className={styles.cardWrapper}>
                        <div className={styles.aiIndicator}>
                            <Bot size={14} />
                            <span>AI Match</span>
                        </div>
                        <PostCard 
                            post={item.post} 
                            {...aiEnhancedProps}
                            onClick={() => onItemClick?.(item)}
                        />
                    </div>
                );
            
            default:
                return null;
        }
    }, [onItemClick]);

    // Loading skeleton
    if (loading) {
        return (
            <div className={`${styles.container} ${className}`}>
                <div className={styles.loadingState}>
                    <div className={styles.loadingSkeleton}>
                        {[1, 2, 3, 4, 5, 6].map(i => (
                            <div key={i} className={styles.skeletonCard} />
                        ))}
                    </div>
                </div>
            </div>
        );
    }

    // No results
    if (!items.length) {
        return (
            <div className={`${styles.container} ${className}`}>
                <div className={styles.emptyState}>
                    <Sparkles size={48} />
                    <h3>No results found</h3>
                    <p>Try refining your search or asking a different question</p>
                </div>
            </div>
        );
    }

    const displayModeClass = {
        grid: styles.gridView,
        list: styles.listView,
        cards: styles.cardsView
    }[currentDisplayMode];

    return (
        <div className={`${styles.container} ${className}`}>
            {/* Results Header */}
            <div className={styles.resultsHeader}>
                <div className={styles.resultsInfo}>
                    <h2>
                        <Sparkles size={20} />
                        AI Search Results
                    </h2>
                    <span className={styles.resultCount}>
                        {metadata.totalCount || items.length} results found
                    </span>
                    {metadata.queryInterpretation && (
                        <span className={styles.queryInterpretation}>
                            for "{metadata.queryInterpretation}"
                        </span>
                    )}
                </div>

                {/* Controls */}
                <div className={styles.controls}>
                    {/* Filter Dropdown */}
                    <div className={styles.filterDropdown}>
                        <button className={styles.dropdownTrigger}>
                            <Filter size={16} />
                            <span>{filterBy === 'all' ? 'All Types' : filterBy}</span>
                            <ChevronDown size={16} />
                        </button>
                        <div className={styles.dropdownMenu}>
                            <button onClick={() => setFilterBy('all')}>All Types</button>
                            {Object.keys(groupedItems).map(type => (
                                <button key={type} onClick={() => setFilterBy(type)}>
                                    {type} ({groupedItems[type].length})
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Sort Dropdown */}
                    <div className={styles.sortDropdown}>
                        <button className={styles.dropdownTrigger}>
                            <SlidersHorizontal size={16} />
                            <span>Sort by {sortBy}</span>
                            <ChevronDown size={16} />
                        </button>
                        <div className={styles.dropdownMenu}>
                            <button onClick={() => setSortBy('relevance')}>
                                <TrendingUp size={14} /> Relevance
                            </button>
                            <button onClick={() => setSortBy('newest')}>
                                Newest First
                            </button>
                            <button onClick={() => setSortBy('price-low')}>
                                Price: Low to High
                            </button>
                            <button onClick={() => setSortBy('price-high')}>
                                Price: High to Low
                            </button>
                        </div>
                    </div>

                    {/* Display Mode Toggle */}
                    <div className={styles.displayModeToggle}>
                        <button 
                            className={currentDisplayMode === 'grid' ? styles.active : ''}
                            onClick={() => setCurrentDisplayMode('grid')}
                            title="Grid view"
                        >
                            <Grid3X3 size={18} />
                        </button>
                        <button 
                            className={currentDisplayMode === 'list' ? styles.active : ''}
                            onClick={() => setCurrentDisplayMode('list')}
                            title="List view"
                        >
                            <List size={18} />
                        </button>
                    </div>
                </div>
            </div>

            {/* Results Grid */}
            <div className={`${styles.resultsGrid} ${displayModeClass}`}>
                {filteredItems.map(renderItem)}
            </div>

            {/* AI Context Footer */}
            {assistantResponse?.Response && (
                <div className={styles.aiContext}>
                    <div className={styles.aiContextHeader}>
                        <Bot size={18} />
                        <span>Assistant's Summary</span>
                    </div>
                    <p>{assistantResponse.Response}</p>
                </div>
            )}
        </div>
    );
};

export default AIResultsDisplay;