// src/utils/categoryUtils.js
/**
 * Recursively flatten a category tree, respecting an expanded map if provided.
 * If expandedMap[cat.id] = true, we include subcategories in the flattened result.
 *
 * @param {Array} categories - Array of category objects
 * @param {Object} expandedMap - Object mapping category IDs to boolean expand state
 * @returns {Array} Flattened array of categories
 */
export function flattenTree(categories, expandedMap = {}) {
    if (!categories || !Array.isArray(categories)) {
        return [];
    }
    const result = [];
    function traverse(list, depth = 0) {
        list.forEach((cat) => {
            // Add current category to results
            result.push({
                ...cat,
                depth // Add depth information for more context
            });
            // If category is expanded and has subcategories, traverse those
            if (expandedMap[cat.id] && Array.isArray(cat.subcategories) && cat.subcategories.length > 0) {
                traverse(cat.subcategories, depth + 1);
            }
        });
    }
    traverse(categories);
    return result;
}
/**
 * Recursively filter categories by a search query (case-insensitive).
 * Returns a new array with only matching categories (and their subtrees if relevant).
 *
 * @param {Array} categories - Array of category objects
 * @param {string} query - Search query string
 * @returns {Array} Filtered array of categories
 */
export function filterCategories(categories, query) {
    if (!categories || !Array.isArray(categories)) {
        return [];
    }
    if (!query || !query.trim()) {
        return categories; // No filter needed
    }
    const lowerQuery = query.toLowerCase().trim();
    return categories.reduce((acc, cat) => {
        // Skip if category doesn't have a name
        if (!cat) return acc;
        // Create a normalized version of the category for matching
        const normalizedCat = {
            ...cat,
            name: cat.name || cat.description || 'Unnamed Category'
        };
        // Check for matches in name, description or id
        const nameMatches = normalizedCat.name.toLowerCase().includes(lowerQuery);
        const descMatches = normalizedCat.description && normalizedCat.description.toLowerCase().includes(lowerQuery);
        const idMatches = normalizedCat.id && normalizedCat.id.toString().toLowerCase().includes(lowerQuery);
        // Filter subcategories recursively
        let filteredSubs = [];
        if (Array.isArray(cat.subcategories) && cat.subcategories.length > 0) {
            filteredSubs = filterCategories(cat.subcategories, query);
        }
        // If this category matches, or any of its subcategories match, keep it
        if (nameMatches || descMatches || idMatches || filteredSubs.length > 0) {
            acc.push({
                ...normalizedCat,
                subcategories: filteredSubs,
            });
        }
        return acc;
    }, []);
}
/**
 * Ensures each category has the required properties for display and interaction
 *
 * @param {Array} categories - Array of category objects from API
 * @returns {Array} Processed categories with normalized properties
 */
export function processCategoryData(categories) {
    if (!categories || !Array.isArray(categories)) {
        return [];
    }
    return categories.map(cat => {
        // Skip if the category is invalid
        if (!cat) return null;
        return {
            ...cat,
            // Ensure name exists (use description if name is missing)
            name: cat.name || cat.description || 'Unnamed Category',
            // Initialize empty subcategories array if none exists
            subcategories: cat.subcategories || [],
        };
    }).filter(Boolean); // Remove any null entries
}
/**
 * Build a nested category tree from a flat array of categories with parent relationships
 *
 * @param {Array} flatCategories - Flat array of categories with parentId references
 * @returns {Array} Nested category tree
 */
export function buildCategoryTree(flatCategories) {
    if (!flatCategories || !Array.isArray(flatCategories)) {
        return [];
    }
    const categoryMap = {};
    const rootCategories = [];
    // First pass: populate the map
    flatCategories.forEach(category => {
        if (!category) return;
        // Create a normalized category with subcategories array
        const normalizedCategory = {
            ...category,
            name: category.name || category.description || 'Unnamed Category',
            subcategories: []
        };
        // Store in map by ID for quick lookup
        categoryMap[category.id] = normalizedCategory;
    });
    // Second pass: build the tree structure
    flatCategories.forEach(category => {
        if (!category) return;
        const normalizedCategory = categoryMap[category.id];
        if (category.parentId && categoryMap[category.parentId]) {
            // This category has a parent, add it to parent's subcategories
            categoryMap[category.parentId].subcategories.push(normalizedCategory);
        } else {
            // This is a root category
            rootCategories.push(normalizedCategory);
        }
    });
    return rootCategories;
}
/**
 * Find a category by ID in a nested category tree
 *
 * @param {Array} categories - Nested category tree
 * @param {string|number} categoryId - ID of category to find
 * @returns {Object|null} Found category or null
 */
export function findCategoryById(categories, categoryId) {
    if (!categories || !Array.isArray(categories) || !categoryId) {
        return null;
    }
    // Convert categoryId to string for consistent comparison
    const idToFind = categoryId.toString();
    // Recursive search function
    function search(cats) {
        for (const cat of cats) {
            // Check current category
            if (cat.id.toString() === idToFind) {
                return cat;
            }
            // Check subcategories if they exist
            if (cat.subcategories && cat.subcategories.length > 0) {
                const found = search(cat.subcategories);
                if (found) return found;
            }
        }
        return null;
    }
    return search(categories);
}
/**
 * Get the path from root to a specific category
 *
 * @param {Array} categories - Nested category tree
 * @param {string|number} categoryId - ID of category to find path to
 * @returns {Array} Array of categories from root to target (inclusive)
 */
export function getCategoryPath(categories, categoryId) {
    if (!categories || !Array.isArray(categories) || !categoryId) {
        return [];
    }
    // Convert categoryId to string for consistent comparison
    const idToFind = categoryId.toString();
    // Recursive search function
    function findPath(cats, path = []) {
        for (const cat of cats) {
            // Try this path
            const newPath = [...path, cat];
            // Check if this is the category we're looking for
            if (cat.id.toString() === idToFind) {
                return newPath;
            }
            // Check subcategories if they exist
            if (cat.subcategories && cat.subcategories.length > 0) {
                const subPath = findPath(cat.subcategories, newPath);
                if (subPath.length > 0) return subPath;
            }
        }
        return []; // Not found in this branch
    }
    return findPath(categories);
}