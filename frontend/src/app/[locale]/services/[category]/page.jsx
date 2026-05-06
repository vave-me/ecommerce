import React from 'react';
import {searchServicesWithFilters, searchServicesWithCategory} from '../../../../api/searchApi';
import ServicesPageClient from '../ServicesPage.client';
import {notFound} from 'next/navigation';
import {getTranslations} from 'next-intl/server';
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';
import styles from "../../../../components/services/ServiceCard.module.css";

export const dynamic = 'force-dynamic';

// Helper function to safely resolve props in Next.js 15+
async function resolveProps(props) {
    // In Next.js 15+, we need to await both params and searchParams
    const params = props.params ? await props.params : {};
    const searchParams = props.searchParams ? await props.searchParams : {};

    return { params, searchParams };
}

export async function generateMetadata(props) {
    try {
        // Safely resolve props first
        const { params } = await resolveProps(props);
        const { locale, category } = params;

        // Get translations for better titles
        const t = await getTranslations({locale, namespace: 'ServicesPage'});
        const categoryName = category ?
            t(`category_${category}_label`, {}, {defaultValue: category}) :
            t('allServices');

        return {
            title: t('metaTitleCategory', {categoryName}),
            description: t('metaDescriptionCategory', {categoryName}),
        };
    } catch (e) {
        // Error: `[generateMetadata] Failed:`, e...
        // Fallback metadata if anything fails
        const fallbackCategoryName = (props.params?.category) ? await props.params.category : 'Services';
        return {
            title: `Services in ${fallbackCategoryName} | Vaveme`,
            description: `Find professional service providers in the ${fallbackCategoryName} category.`
        };
    }
}

export default async function ServicesCategoryPage(props) {
    try {
        // 1. Safely resolve props first
        const { params, searchParams } = await resolveProps(props);

        // 2. Now it's safe to destructure params
        const { locale, category } = params;

        // Validate category early
        if (!category) {
            // Error: "ServicesCategoryPage: Category parameter is missi...
            notFound();
        }

        // 3. Get translations
        const t = await getTranslations({locale, namespace: 'ServicesPage'});

        // 4. Extract searchParams values
        const displayMode = searchParams?.displayMode || 'grid';
        const page = parseInt(searchParams?.page || '1', 10);
        const limit = 20;
        const sortBy = searchParams?.sortBy || '';

        // 5. Create listing filters for searchAPI
        const listingFilters = {
            category,
            displayMode,
            page,
            limit,
            sortBy,
        };

        // 6. Fetch services filtered by category using search API (cached proxy)
        let fetchedData = {services: [], totalCount: 0, totalPages: 0};
        let fetchError = null;

        try {
            fetchedData = await searchServicesWithCategory(category, listingFilters);
        } catch (err) {
            // Error: `Error fetching services on server (/services/${ca...:`, err);
            fetchError = t('errorFetchingCategoryServices', {
                category,
                message: err.message || 'Unknown fetch error'
            });
        }

        // 7. Process the fetched data
        const serverServices = fetchedData?.services || [];
        const totalPages = parseInt(fetchedData?.totalPages || '0', 10);
        const totalCount = parseInt(fetchedData?.totalCount || '0', 10);
        const currentPage = parseInt(fetchedData?.currentPage || '1', 10);

        // 8. Get category name for JSON-LD
        const jsonLdCategoryName = category ?
            t(`category_${category}_label`, {}, {defaultValue: category}) :
            t('allServices');

        // 9. Construct JSON-LD for SEO with correct service URLs
        const servicesJsonLd = {
            '@context': 'https://schema.org',
            '@type': 'ItemList',
            'name': t('jsonLdListName', {categoryName: jsonLdCategoryName}),
            'description': `Professional service providers in the ${category} category.`,
            itemListElement: serverServices.map((service, index) => {
                return {
                    '@type': 'ListItem',
                    position: index + 1,
                    url: `https://www.sfx-markt.de/${locale}/services/${category}/${service.id}`,
                    item: {
                        '@type': 'Service',
                        name: service.name,
                        image: service.thumbnail ?? '',
                        description: service.description?.substring(0, 150) || '',
                        identifier: {'@type': 'PropertyValue', 'name': 'serviceId', 'value': service.id},
                        datePosted: service.postedDate || service.createdAt,
                        provider: {
                            '@type': 'Organization',
                            name: service.providerName ?? '',
                        },
                        serviceType: category,
                        areaServed: service.location ?? '',
                        offers: {
                            '@type': 'Offer',
                            price: service.basePrice ?? '0.00',
                            priceCurrency: 'EUR',
                            availability: 'https://schema.org/InStock',
                        },
                    },
                };
            }),
        };

        // 10. Return UI components
        return (
            <>
                <script
                    type="application/ld+json"
                    dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(servicesJsonLd)}}
                />
                <ServicesPageClient
                    serverServices={serverServices}
                    serverFilters={listingFilters}
                    fetchError={fetchError}
                    currentCategory={category}
                    totalPages={totalPages}
                    currentPage={currentPage}
                    totalCount={totalCount}
                />
            </>
        );
    } catch (e) {
        // Global error handler
        // Error: `[ServicesCategoryPage] Failed during page renderi...

        return (
            <div className={styles.errorContainer}>
                <h1 className={styles.errorTitle}>Error Loading Services</h1>
                <p className={styles.errorMessage}>We couldn't load this category. Please try again later.</p>
            </div>
        );
    }
}