"use client";
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import { Plus, Minus } from '@/icons';
import { useItems } from '../../hooks/queries/useItemsQuery';
import styles from './ItemSelector.module.css';
/**
 * ItemSelector Component
 * Fetches and displays items for newsletter selection using React Query
 */
const ItemSelector = memo(({ onItemSelect, selectedItems = [] }) => {
    const { data: items = [], isLoading: loading, error } = useItems();
    const handleItemToggle = (item) => {
        if (onItemSelect) {
            onItemSelect(item);
        }
    };
    if (loading) {
        return <div className={styles.loading}>Loading items...</div>;
    }
    if (error) {
        return <div className={styles.error}>Error loading items: {error.message}</div>;
    }
    if (items.length === 0) {
        return <div className={styles.empty}>No items available for selection.</div>;
    }
    return (
        <div className={styles.itemSelector}>
            <h3>Select Items for Newsletter</h3>
            <div className={styles.itemsGrid}>
                {items.map((item) => {
                    const isSelected = selectedItems.some(selected => selected.id === item.id);
                    return (
                        <div 
                            key={item.id} 
                            className={`${styles.itemCard} ${isSelected ? styles.selected : ''}`}
                            onClick={() => handleItemToggle(item)}
                        >
                            <h4>{item.title || item.name}</h4>
                            <p>{item.description}</p>
                            <div className={styles.itemMeta}>
                                <span className={styles.category}>{item.category}</span>
                                {item.price && <span className={styles.price}>{item.price}</span>}
                            </div>
                        </div>
                    );
                })}
            </div>
        </div>
    );
});
ItemSelector.displayName = 'ItemSelector';
ItemSelector.propTypes = {
    onItemSelect: PropTypes.func,
    selectedItems: PropTypes.arrayOf(PropTypes.object),
};
export default ItemSelector;