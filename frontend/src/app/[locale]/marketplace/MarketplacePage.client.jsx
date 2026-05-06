"use client";
import React from 'react';
import styles from './MarketplacePage.module.css';
// Leftside removed - using HorizontalFilters mobile modal instead
import { useIsMobile } from '../../../hooks/useMobileDetection';
import {useSelector} from 'react-redux';
import useMasonryLayout from '../../../hooks/useMasonryLayout';
// Import the actual component, not the page
import ClassifiedCard from "../../../components/classified/ClassifiedCard";
/**
 * MarketplacePageClient - Client component for the marketplace listing page
 *
 * @param {Object} props - Component props
 * @param {Array} props.serverProducts - Products data fetched from the server
 * @param {Object} props.serverFilters - Filters data fetched from the server
 * @returns {JSX.Element} Rendered component
 */
export default function MarketplacePageClient({serverProducts = [], serverFilters = {}}) {
    // Hooks
    const isMobile = useIsMobile();
    const listingFilters = useSelector((state) => state.listingFilters);
    const {displayMode} = listingFilters;
    // Empty state handling
    if (!serverProducts.length) {
        return <div className={styles.emptyState}>No products found.</div>;
    }
    return (
        <div className={styles.container}>
            <main className={styles.mainContent}>
                    <GridLayout
                        isMobile={isMobile}
                        products={serverProducts}
                    />
            </main>
        </div>
    );
}
/**
 * GridLayout - Component for grid view of products
 */
function GridLayout({isMobile, products}) {
    // Use masonry layout for desktop only
    const masonryRef = useMasonryLayout([products, !isMobile]);
    
    return (
        <div className={styles.layoutGrid}>
            {/* Leftside removed - using HorizontalFilters mobile modal instead */}
            <section className={styles.contentArea}>
                <ul className={styles.productList} ref={!isMobile ? masonryRef : null}>
                    {products.map((prod) => (
                        <li key={prod.id}>
                            <ClassifiedCard product={prod}/>
                        </li>
                    ))}
                </ul>
            </section>
        </div>
    );
} 