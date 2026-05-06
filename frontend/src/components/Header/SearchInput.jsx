"use client";
import React, { forwardRef } from "react";
import { Search, X } from "@/icons";
import styles from "./SearchBar.module.css";
/**
 * SearchInput Component
 * Main search input with search icon and clear button
 * Extracted from SearchBar to improve component modularity
 */
const SearchInput = forwardRef(({
    searchQuery,
    onInputChange,
    onKeyDown,
    onFocus,
    onClear,
    showSuggestions,
    placeholder = "Search...",
    onSubmit,
    isMobile = false,
    ...props
}, ref) => {
    return (
        <div className={styles.searchSection}>
            {/* Search icon / Submit button */}
            <button
                type="submit"
                aria-label="Submit search"
                className={styles.searchIcon}
                onClick={onSubmit}
            >
                <Search size={18}/>
            </button>
            {/* Main search input */}
            <input
                type="text"
                placeholder={placeholder}
                aria-label="Search Products"
                aria-autocomplete="list"
                aria-controls="search-suggestions"
                aria-expanded={showSuggestions}
                value={searchQuery}
                onChange={onInputChange}
                onFocus={onFocus}
                onKeyDown={onKeyDown}
                ref={ref}
                className={styles.searchInput}
                autoComplete="off"
                {...props}
            />
            {/* Clear button */}
            {searchQuery && (
                <button
                    type="button"
                    onClick={onClear}
                    aria-label="Clear Search"
                    className={styles.clearButton}
                >
                    <X size={16}/>
                </button>
            )}
        </div>
    );
});
SearchInput.displayName = "SearchInput";
export default SearchInput; 