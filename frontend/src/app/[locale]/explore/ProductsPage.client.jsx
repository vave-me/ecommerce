"use client";
import React from 'react';
import {useQuery} from '@tanstack/react-query';
import {useSelector} from 'react-redux';
// No need to import useTranslations here if labels are passed as props
// Import Components & Hooks (adjust paths)
// Leftside removed - using HorizontalFilters mobile modal instead
import ClassifiedCard from "../../../components/classified/ClassifiedCard"; // Following DealCard pattern
import {fetchProductsByFilters} from '../../../api/productsApi';
import { useIsMobile } from '../../../hooks/useMobileDetection';
import styles from './ProductsPage.module.css'; // Adjust path
// --- Loading & Error Placeholders ---
// Now receive labels as props
const LoadingPlaceholder = ({label}) => <div className={styles.loadingState}>{label}</div>;
const ErrorPlaceholder = ({error, labels}) => (
    <div className={styles.errorState}>
        <h2>{labels.errorTitle}</h2>
        <p>{labels.errorMessage}</p>
        {process.env.NODE_ENV === 'development' && error?.message && (
            <pre className={styles.devError}>{labels.errorDetailPrefix} {error.message}</pre>
        )}
    </div>
);
const EmptyPlaceholder = ({label}) => <div className={styles.emptyState}>{label}</div>;
// Receive labels prop from the server component
export default function ProductsPageClient({locale, labels}) {
    // --- Hooks ---
    const isMobile = useIsMobile();
    const listingFilters = useSelector((state) => state.listingFilters);
    const {displayMode = 'list'} = listingFilters || {};
    // --- Data Fetching (Client-Side) ---
    const {
        data: featuredData,
        isLoading,
        isError,
        error,
        isFetching,
    } = useQuery({
        // Include locale in queryKey if API uses it and filters change based on it
        queryKey: ['products', listingFilters, locale],
        queryFn: () => fetchProductsByFilters({...listingFilters, locale}), // Pass locale to API
        // keepPreviousData: true, // Optional: uncomment for smoother UX
        // staleTime: 5 * 60 * 1000, // Optional: 5 minutes
    });
    // --- Render Logic ---
    // 1. Loading State (Use passed label)
    if (isLoading) {
        return <LoadingPlaceholder label={labels.loading}/>;
    }
    // 2. Error State (Use passed labels)
    if (isError) {
        return <ErrorPlaceholder error={error} labels={labels}/>;
    }
    // 3. Empty or Invalid Data (Use passed label)
    const products = featuredData?.products;
    if (!Array.isArray(products) || products.length === 0) {
        return <EmptyPlaceholder label={labels.empty}/>;
    }
    const validProducts = products.filter(prod => prod && prod.id);
    if (validProducts.length === 0) {
        return <EmptyPlaceholder label={labels.empty}/>;
    }
    // 4. Render Product List/Grid
    return (
        <div className={styles.container}>
            {isFetching && !isLoading && <div className={styles.refetchingIndicator}>{labels.refetching}</div>}
            <main className={styles.mainContent}>
                <GridLayout isMobile={isMobile} products={validProducts} locale={locale}/>
            </main>
        </div>
    );
}
function GridLayout({isMobile, products, locale}) {
    return (
        <div className={styles.layoutGrid}>
            {!isMobile && <Leftside/>}
            <section className={styles.contentArea}>
                <ul className={styles.productList}>
                    {products.map((prod) => (
                        <li key={prod.id}>
                            {/* Assuming ImprovedClassifiedCard is similar to ProductCard and needs locale */}
                            <ClassifiedCard product={prod} locale={locale}/>
                        </li>
                    ))}
                </ul>
            </section>
        </div>
    );
}