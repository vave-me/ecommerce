import React, { useState, memo } from 'react';
import { Filter, X, Search, MapPin, Euro, Calendar } from '@/icons';
import styles from './ListingsFilters.module.css';
/**
 * ListingsFilters - Atomic Design Component
 * Advanced filtering system for listings
 * 
 * @param {Object} props - Component props
 * @param {Object} props.filters - Current filter values
 * @param {Function} props.onFiltersChange - Filter change handler
 * @param {Array} props.categories - Available categories
 * @param {boolean} props.isOpen - Filter panel open state
 * @param {Function} props.onToggle - Toggle filter panel
 * @returns {JSX.Element} Rendered listings filters
 */
const ListingsFilters = memo(({ 
    filters = {}, 
    onFiltersChange = () => {}, 
    categories = [],
    isOpen = false,
    onToggle = () => {}
}) => {
    const [localFilters, setLocalFilters] = useState(filters);
    const handleFilterChange = (key, value) => {
        const newFilters = { ...localFilters, [key]: value };
        setLocalFilters(newFilters);
        onFiltersChange(newFilters);
    };
    const handleClearFilters = () => {
        const clearedFilters = {};
        setLocalFilters(clearedFilters);
        onFiltersChange(clearedFilters);
    };
    const hasActiveFilters = Object.keys(localFilters).some(key => 
        localFilters[key] && localFilters[key] !== ''
    );
    return (
        <div className={styles.container}>
            {/* Filter Toggle Button */}
            <button 
                className={`${styles.toggleButton} ${isOpen ? styles.active : ''}`}
                onClick={onToggle}
                aria-label="Toggle filters"
            >
                <Filter size={16} />
                <span>Filters</span>
                {hasActiveFilters && <div className={styles.activeIndicator} />}
            </button>
            {/* Filter Panel */}
            <div className={`${styles.panel} ${isOpen ? styles.open : ''}`}>
                <div className={styles.panelHeader}>
                    <h3 className={styles.panelTitle}>Filter Listings</h3>
                    <button 
                        className={styles.closeButton}
                        onClick={onToggle}
                        aria-label="Close filters"
                    >
                        <X size={16} />
                    </button>
                </div>
                <div className={styles.panelContent}>
                    {/* Search Filter */}
                    <div className={styles.filterGroup}>
                        <label className={styles.filterLabel}>
                            <Search size={14} />
                            <span>Search</span>
                        </label>
                        <input
                            type="text"
                            className={styles.filterInput}
                            placeholder="Search listings..."
                            value={localFilters.search || ''}
                            onChange={(e) => handleFilterChange('search', e.target.value)}
                        />
                    </div>
                    {/* Category Filter */}
                    <div className={styles.filterGroup}>
                        <label className={styles.filterLabel}>
                            <span>Category</span>
                        </label>
                        <select
                            className={styles.filterSelect}
                            value={localFilters.category || ''}
                            onChange={(e) => handleFilterChange('category', e.target.value)}
                        >
                            <option value="">All Categories</option>
                            {categories.map(category => (
                                <option key={category.id} value={category.id}>
                                    {category.name}
                                </option>
                            ))}
                        </select>
                    </div>
                    {/* Location Filter */}
                    <div className={styles.filterGroup}>
                        <label className={styles.filterLabel}>
                            <MapPin size={14} />
                            <span>Location</span>
                        </label>
                        <input
                            type="text"
                            className={styles.filterInput}
                            placeholder="Enter location..."
                            value={localFilters.location || ''}
                            onChange={(e) => handleFilterChange('location', e.target.value)}
                        />
                    </div>
                    {/* Price Range Filter */}
                    <div className={styles.filterGroup}>
                        <label className={styles.filterLabel}>
                                                            <Euro size={14} />
                            <span>Price Range</span>
                        </label>
                        <div className={styles.priceRange}>
                            <input
                                type="number"
                                className={styles.filterInput}
                                placeholder="Min price"
                                value={localFilters.minPrice || ''}
                                onChange={(e) => handleFilterChange('minPrice', e.target.value)}
                            />
                            <span className={styles.priceSeparator}>to</span>
                            <input
                                type="number"
                                className={styles.filterInput}
                                placeholder="Max price"
                                value={localFilters.maxPrice || ''}
                                onChange={(e) => handleFilterChange('maxPrice', e.target.value)}
                            />
                        </div>
                    </div>
                    {/* Date Filter */}
                    <div className={styles.filterGroup}>
                        <label className={styles.filterLabel}>
                            <Calendar size={14} />
                            <span>Date Posted</span>
                        </label>
                        <select
                            className={styles.filterSelect}
                            value={localFilters.dateRange || ''}
                            onChange={(e) => handleFilterChange('dateRange', e.target.value)}
                        >
                            <option value="">Any time</option>
                            <option value="today">Today</option>
                            <option value="week">This week</option>
                            <option value="month">This month</option>
                            <option value="3months">Last 3 months</option>
                        </select>
                    </div>
                    {/* Sort Filter */}
                    <div className={styles.filterGroup}>
                        <label className={styles.filterLabel}>
                            <span>Sort by</span>
                        </label>
                        <select
                            className={styles.filterSelect}
                            value={localFilters.sortBy || 'newest'}
                            onChange={(e) => handleFilterChange('sortBy', e.target.value)}
                        >
                            <option value="newest">Newest first</option>
                            <option value="oldest">Oldest first</option>
                            <option value="price-low">Price: Low to High</option>
                            <option value="price-high">Price: High to Low</option>
                            <option value="popular">Most popular</option>
                        </select>
                    </div>
                </div>
                {/* Panel Footer */}
                <div className={styles.panelFooter}>
                    <button 
                        className={styles.clearButton}
                        onClick={handleClearFilters}
                        disabled={!hasActiveFilters}
                    >
                        Clear All
                    </button>
                </div>
            </div>
            {/* Overlay */}
            {isOpen && <div className={styles.overlay} onClick={onToggle} />}
        </div>
    );
});
ListingsFilters.displayName = 'ListingsFilters';
export default ListingsFilters; 