import React from 'react';
import {getTranslations} from 'next-intl/server'; // Import
import {fetchProductsByFilters} from '../../../api/productsApi'; // Adjust path
import ProductsPageClient from "./ProductsPage.client"; // Adjust path
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';
// Helper function (assume defined)
async function resolveProps(props) {
    const params = props.params ? await props.params : {};
    const searchParams = props.searchParams ? await props.searchParams : {};
    return {params, searchParams};
}
export const dynamic = 'force-dynamic';
export async function generateMetadata(props) {
    let t;
    try {
        const {params} = await resolveProps(props);
        const {locale} = params;
        t = await getTranslations({locale, namespace: "ProductsPage"}); // Use namespace
        return {
            title: t('metaTitle'), // Translated
            description: t('metaDescription'), // Translated
        };
    } catch (e) {
        return { // Hardcoded fallback
            title: 'Products | Vaveme',
            description: 'Browse our selection of high-quality products.',
        };
    }
}
export default async function ProductsPageServerEntry(props) {
    let t; // For potential use in catch block
    let locale; // For potential use in catch block
    try {
        const {params} = await resolveProps(props);
        locale = params.locale; // Assign locale
        t = await getTranslations({locale, namespace: "ProductsPage"}); // Get translations
        let serverFetchedProducts = [];
        try {
            const seoFilters = {limit: 10, locale}; // Example SEO fetch
            const featuredData = await fetchProductsByFilters(seoFilters);
            serverFetchedProducts = featuredData?.products || [];
        } catch (err) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', err);
        }
        throw err; // Re-throw for caller to handle
    }
        // Use common translations for schema types/URLs
        const conditionMap = {
            new: t('common.conditionSchemaNew', {}, { fallbackNamespace: 'ProductsPage' }),
            used: t('common.conditionSchemaUsed', {}, { fallbackNamespace: 'ProductsPage' }),
            refurbished: t('common.conditionSchemaRefurbished', {}, { fallbackNamespace: 'ProductsPage' }),
            // Add other conditions if needed
        };
        const availabilityMap = {
            inStock: t('common.availabilitySchemaInStock', {}, { fallbackNamespace: 'ProductsPage' }),
            outOfStock: t('common.availabilitySchemaOutOfStock', {}, { fallbackNamespace: 'ProductsPage' }),
            // Add other availabilities if needed
        };
        const productsJsonLd = {
            '@context': 'https://schema.org',
            '@type': 'ItemList',
            'name': t('jsonLdListName'), // Translated list name
            'description': t('jsonLdListDescription'), // Translated list description
            itemListElement: serverFetchedProducts.map((prod, index) => {
                if (!prod || !prod.id || !prod.name) {
                    return null;
                }
                // Use mapped condition, default to New if unknown
                const itemCondition = conditionMap[prod.condition] || conditionMap['new'];
                // Use mapped availability, default to InStock if unknown
                const availability = availabilityMap[prod.availability] || availabilityMap['inStock'];
                return {
                    '@type': 'ListItem',
                    position: index + 1,
                    url: `https://www.sfx-market.de/${locale}/product/${prod.id}`, // Use product ID if no slug
                    name: prod.name, // Use product name directly
                    item: {
                        '@type': 'Product',
                        name: prod.name,
                        image: prod.thumbnail || prod.images?.[0] || undefined, // Use thumbnail or first image
                        description: prod.description?.substring(0, 160) || undefined, // Truncate description
                        brand: prod.brand ? {'@type': 'Brand', name: prod.brand} : {
                            '@type': 'Brand',
                            name: t('common.defaultBrand')
                        }, // Add fallback brand
                        mpn: prod.mpn || prod.sku || prod.model || undefined,
                        itemCondition: itemCondition,
                        offers: {
                            '@type': 'Offer',
                            price: String(prod.basePrice ?? '0.00'),
                            priceCurrency: t('common.currencyCode'), // Use common currency code
                            availability: availability,
                            url: `https://www.sfx-market.de/${locale}/product/${prod.id}` // Offer URL
                        },
                        // Add SKU if available
                        ...(prod.sku && {sku: prod.sku}),
                    },
                };
            }).filter(Boolean),
        };
        return (
            <>
                {productsJsonLd.itemListElement.length > 0 && (
                                <script 
              type="application/ld+json"
              dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(productsJsonLd)}}
            />
                )}
                {/* Pass necessary translated labels to the client */}
                <ProductsPageClient
                    locale={locale}
                    labels={{
                        loading: t('loading'),
                        errorTitle: t('errorTitle'),
                        errorMessage: t('errorMessage'),
                        errorDetailPrefix: t('errorDetailPrefix'),
                        empty: t('empty'),
                        refetching: t('refetching')
                    }}
                />
            </>
        );
    } catch (e) {
        // Attempt to get translations for error page
        let errorTitle = "Error Loading Products";
        let errorMessage = "Could not load products at this time. Please try again later.";
        try {
            if (!t) t = await getTranslations({locale: locale || 'en', namespace: "ProductsPage"});
            errorTitle = t('globalErrorTitle'); // Use a specific global error key
            errorMessage = t('globalErrorMessage'); // Use a specific global error key
        } catch (transErr) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', transErr);
        }
    }
        return (
            <div className="container mx-auto p-8 text-center">
                <h1 className="text-2xl font-bold mb-4">{errorTitle}</h1>
                <p>{errorMessage}</p>
            </div>
        );
    }
}