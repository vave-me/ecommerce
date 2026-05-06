"use client";
import React, { forwardRef, memo } from "react";
import { MapPin, X, ChevronDown, Filter } from "@/icons";
import styles from "./SearchBar.module.css";
/**
 * FilterBar Component
 * Contains location input, radius select (desktop) and mobile filter button
 * Extracted from SearchBar to improve component modularity
 */
const FilterBar = memo(({
    // Location props
    locationQuery,
    onLocationInputChange,
    onLocationKeyDown,
    onLocationFocus,
    onLocationClear,
    locationInputRef,
    showLocationSuggestions,
    // Radius props
    radius,
    onRadiusChange,
    // Mobile filter props
    isMobile,
    filtersActive,
    onToggleMobileFilters,
    // Additional props
    ...props
}) => {
    return (
        <>
            {/* Mobile filter button */}
            {isMobile && (
                <button
                    type="button"
                    onClick={onToggleMobileFilters}
                    aria-label={`${filtersActive ? 'Edit active filters' : 'Open filters'}`}
                    className={`${styles.filterButton} ${filtersActive ? styles.filterActive : ''}`}
                >
                    <Filter size={18}/>
                </button>
            )}
            {/* DESKTOP ONLY - LOCATION SEARCH */}
            {!isMobile && (
                <div className={styles.locationSection}>
                    {/* Location icon */}
                    <MapPin size={18} className={styles.locationIcon}/>
                    <input
                        type="text"
                        placeholder="Location..."
                        aria-label="Search by location"
                        aria-autocomplete="list"
                        aria-controls="location-suggestions"
                        aria-expanded={showLocationSuggestions}
                        value={locationQuery}
                        onChange={onLocationInputChange}
                        onFocus={onLocationFocus}
                        onKeyDown={onLocationKeyDown}
                        ref={locationInputRef}
                        className={styles.locationInput}
                        autoComplete="off"
                    />
                    {/* Clear location button */}
                    {locationQuery && (
                        <button
                            type="button"
                            onClick={onLocationClear}
                            aria-label="Clear Location"
                            className={styles.clearButton}
                        >
                            <X size={16}/>
                        </button>
                    )}
                </div>
            )}
            {/* DESKTOP ONLY - RADIUS SELECT */}
            {!isMobile && (
                <div className={styles.radiusSection}>
                    <select
                        className={styles.radiusSelect}
                        value={radius}
                        onChange={onRadiusChange}
                        aria-label="Search Radius"
                    >
                        <option value="5">5km</option>
                        <option value="10">10km</option>
                        <option value="25">25km</option>
                        <option value="50">50km</option>
                        <option value="100">100km</option>
                    </select>
                    <ChevronDown size={16} className={styles.selectIcon}/>
                </div>
            )}
        </>
    );
});
FilterBar.displayName = 'FilterBar';
export default FilterBar; 