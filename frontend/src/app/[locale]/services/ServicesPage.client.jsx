"use client";

import React, {memo, useState, useEffect} from 'react';
import styles from './ServicesPage.module.css';
import {useIsMobile} from '../../../hooks/useMobileDetection';
import {useSelector, useDispatch} from 'react-redux';
import {useFilteredSearch} from '../../../hooks/useFilteredSearch';
import {setFilters} from '../../../redux/slices/listingFiltersSlice';

// Import the actual component, not the page
import ServiceCard from "../../../components/services/ServiceCard";

/**
 * ServicesPageClient - Client component for the services listing page
 *
 * @param {Object} props - Component props
 * @param {Array} props.serverServices - Services data fetched from the server
 * @param {Object} props.serverFilters - Filters data fetched from the server
 * @param {string} props.fetchError - Any error message from server
 * @param {string} props.currentCategory - Current category slug
 * @param {number} props.totalPages - Total number of pages
 * @param {number} props.currentPage - Current page number
 * @param {number} props.totalCount - Total count of items
 * @param {string} props.pageTitle - Page title from translations
 * @param {string} props.emptyMsg - Empty state message from translations
 * @returns {JSX.Element} Rendered component
 */
const ServicesPageClient = memo(function ServicesPageClient({
                                                                serverServices = [],
                                                                serverFilters = {},
                                                                fetchError = null,
                                                                currentCategory = null,
                                                                totalPages = 0,
                                                                currentPage = 1,
                                                                totalCount = 0,
                                                                pageTitle = "Services",
                                                                emptyMsg = "No services found"
                                                            }) {
    // Hooks
    const isMobile = useIsMobile();
    const dispatch = useDispatch();
    const listingFilters = useSelector((state) => state.listingFilters);
    const {displayMode} = listingFilters;
    
    // State for services
    const [services, setServices] = useState(serverServices);
    const [categoryReady, setCategoryReady] = useState(false);
    
    // Sync category from URL with Redux state
    useEffect(() => {
        if (currentCategory) {
            dispatch(setFilters({ 
                ...listingFilters,
                category: currentCategory 
            }));
            setCategoryReady(true);
        } else {
            setCategoryReady(true);
        }
        
        // Clear category filter when component unmounts
        return () => {
            if (currentCategory) {
                dispatch(setFilters({ 
                    ...listingFilters,
                    category: '' 
                }));
            }
        };
    }, [currentCategory, dispatch]);
    
    // Use filtered search hook to get services based on Redux filters
    const { 
        data: searchResults, 
        isLoading, 
        error: searchError 
    } = useFilteredSearch({
        entityType: 'service',
        enabled: categoryReady,
        onSuccess: (data) => {
            if (data?.services) {
                setServices(data.services);
            }
        }
    });
    
    // Update services when search results change
    useEffect(() => {
        if (searchResults?.services) {
            setServices(searchResults.services);
        }
    }, [searchResults]);

    // Error state handling
    if (fetchError || searchError) {
        return (
            <div className={styles.errorState}>
                <h2>Error Loading Services</h2>
                <p>{fetchError || searchError?.message || 'An error occurred'}</p>
                <button onClick={() => window.location.reload()}>Retry</button>
            </div>
        );
    }
    
    // Loading state
    if (isLoading && !services.length) {
        return (
            <div className={styles.loadingState}>
                <div className={styles.spinner}></div>
                <p>Loading services...</p>
            </div>
        );
    }

    // Empty state handling
    if (!services.length && !isLoading) {
        return (
            <div className={styles.emptyState}>
                <h2>{emptyMsg}</h2>
                <p>
                    {currentCategory
                        ? `No services available in the ${currentCategory} category.`
                        : 'No services available at the moment.'
                    }
                </p>
            </div>
        );
    }

    return (
        <div className={styles.container}>
            <main className={styles.mainContent}>
                {displayMode === 'list' ? (
                    <ListLayout
                        isMobile={isMobile}
                        services={services}
                    />
                ) : (
                    <GridLayout
                        isMobile={isMobile}
                        services={services}
                    />
                )}
            </main>
        </div>
    );
});

/**
 * ListLayout - Component for list view of services
 */
const ListLayout = memo(function ListLayout({isMobile, services}) {
    return (
        <div className={styles.layoutListGrid}>

            <section className={styles.contentListArea}>
                {services.map((service) => (
                    <ServiceListItem
                        key={service.id}
                        serviceId={service.id}
                        service={service}
                    />
                ))}
            </section>
        </div>
    );
});

/**
 * GridLayout - Component for grid view of services
 */
const GridLayout = memo(function GridLayout({isMobile, services}) {
    return (
        <div className={styles.layoutGrid}>

            <section className={styles.contentArea}>
                <ul className={styles.serviceList}>
                    {services.map((service) => (
                        <li key={service.id}>
                            <ServiceCard service={service}/>
                        </li>
                    ))}
                </ul>
            </section>
        </div>
    );
});

/**
 * ServiceListItem - Component for individual service in list view
 */
const ServiceListItem = memo(function ServiceListItem({service, serviceId}) {
    return (
        <div className={styles.serviceListItem}>
            <ServiceCard service={service} displayMode="list"/>
        </div>
    );
});

export default ServicesPageClient;