import React from 'react';
import {getTranslations} from 'next-intl/server'; // Import
import {fetchProductsByFilters} from '../../../../api/productsApi'; // Keep the correct API function
import ProductsPageClient from '../ProductsPage.client'; // Adjust path
import {notFound} from 'next/navigation';
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';
export const dynamic = 'force-dynamic';
// Helper function (assume defined)
async function resolveProps(props) {
    const params = props.params ? await props.params : {};
    const searchParams = props.searchParams ? await props.searchParams : {};
    return {params, searchParams};
}
export async function generateMetadata(props) {
    let t;
    let category;
    try {
        const {params} = await resolveProps(props);
        const {locale} = params;
        category = params.category; // Assign category
        t = await getTranslations({locale, namespace: 'ProductsPage'});
        // Get translated category name, fallback to raw category slug
        const categoryName = category;
        return {
            title: t('categoryMetaTitle', {categoryName}), // Use translated category name
            description: t('categoryMetaDescription', {categoryName}) // Use translated category name
        };
    } catch (e) {
        const fallbackCategoryName = category || 'Category';
        // Use translated fallbacks if possible
        const fallbackTitle = t ? t('categoryMetaTitleFallback', {categoryName: fallbackCategoryName}) : `${fallbackCategoryName} Products`;
        const fallbackDescription = t ? t('categoryMetaDescriptionFallback', {categoryName: fallbackCategoryName}) : `Browse products in the ${fallbackCategoryName} category.`;
        return {
            title: fallbackTitle,
            description: fallbackDescription,
        };
    }
}
export default async function ProductsCategoryPage(props) {
    let t; // For potential use in catch block
    let category; // For potential use in catch block
    let locale; // For potential use in catch block
    try {
        const {params, searchParams} = await resolveProps(props);
        locale = params.locale; // Assign locale
        category = params.category; // Assign category
        t = await getTranslations({locale, namespace: 'ProductsPage'}); // Get translations
        if (!category) {
             // Translated log
            notFound();
        }
        // Get translated category name for UI/Error messages
        const categoryName = category;
        const displayMode = searchParams?.displayMode || 'grid';
        const page = parseInt(searchParams?.page || '1', 10);
        const limit = 20;
        const sortBy = searchParams?.sortBy || '';
        const listingFilters = {
            category, // Pass category slug
            displayMode,
            page,
            limit,
            sortBy,
            locale,
        };
        let fetchedData = {products: [], totalCount: 0, totalPages: 0, currentPage: 1};
        let fetchError = null;
        try {
            // Use fetchProductsByFilters which handles category filtering
            fetchedData = await fetchProductsByFilters(listingFilters);
        } catch (err) {
             // Translated log
            fetchError = t('errorFetchingCategoryProducts', {category: categoryName}); // Translated error message
        }
        const serverProducts = fetchedData?.products || [];
        const totalPages = parseInt(fetchedData?.totalPages || '0', 10);
        const totalCount = parseInt(fetchedData?.totalCount || '0', 10);
        const currentPage = parseInt(fetchedData?.currentPage || '1', 10);
        // Use common translations for schema types/URLs
        const conditionMap = {
            new: t('common.conditionSchemaNew', {}, { fallbackNamespace: 'ProductsPage' }),
            used: t('common.conditionSchemaUsed', {}, { fallbackNamespace: 'ProductsPage' }),
            refurbished: t('common.conditionSchemaRefurbished', {}, { fallbackNamespace: 'ProductsPage' }),
        };
        const availabilityMap = {
            inStock: t('common.availabilitySchemaInStock', {}, { fallbackNamespace: 'ProductsPage' }),
            outOfStock: t('common.availabilitySchemaOutOfStock', {}, { fallbackNamespace: 'ProductsPage' }),
        };
        const productsJsonLd = {
            '@context': 'https://schema.org',
            '@type': 'ItemList',
            'name': t('jsonLdListNameCategory', {categoryName}), // Translated name
            'description': t('jsonLdListDescriptionCategory', {categoryName}), // Translated description
            itemListElement: serverProducts.map((product, index) => {
                if (!product || !product.id || !product.name) return null;
                const itemCondition = conditionMap[product.condition] || conditionMap['new'];
                const availability = availabilityMap[product.availability] || availabilityMap['inStock'];
                const productUrl = `https://www.sfx-market.de/${locale}/products/${category}/${product.slug || product.id}`; // Construct URL
                return {
                    '@type': 'ListItem',
                    position: (currentPage - 1) * limit + index + 1, // Calculate position
                    url: productUrl,
                    item: { // Embed Product data
                        '@type': 'Product',
                        name: product.name,
                        image: product.thumbnail || product.images?.[0] || undefined,
                        description: product.description?.substring(0, 160) || undefined,
                        brand: product.brand ? {'@type': 'Brand', name: product.brand} : {
                            '@type': 'Brand',
                            name: t('common.defaultBrand')
                        },
                        mpn: product.mpn || product.sku || product.model || undefined,
                        itemCondition: itemCondition,
                        offers: {
                            '@type': 'Offer',
                            price: String(product.basePrice ?? '0.00'),
                            priceCurrency: t('common.currencyCode'),
                            availability: availability,
                            url: productUrl
                        },
                        ...(product.sku && {sku: product.sku}),
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
                {/* Pass translated labels needed by the client */}
                <ProductsPageClient
                    // serverProducts={serverProducts} // Client fetches its own data
                    serverFilters={listingFilters} // Pass initial filters if needed by client
                    // fetchError={fetchError} // Client handles its own fetch errors
                    currentCategory={category} // Pass category slug
                    // Pass labels, client handles loading/error/empty states
                    labels={{
                        loading: t('loading'),
                        errorTitle: t('errorTitleCategory', {categoryName}), // Category specific error title
                        errorMessage: t('errorMessageCategory', {categoryName}), // Category specific error message
                        errorDetailPrefix: t('errorDetailPrefix'),
                        empty: t('emptyCategory', {categoryName}), // Category specific empty message
                        refetching: t('refetching')
                    }}
                    locale={locale}
                    // Pass pagination from server only if client doesn't recalculate it
                    // totalPages={totalPages}
                    // currentPage={currentPage}
                    // totalCount={totalCount}
                />
            </>
        );
    } catch (e) {
        // Attempt to get translations for error page
        let errorTitle = "Error Loading Category Products";
        let errorMessage = "Could not load products for this category. Please try again later.";
        try {
            if (!t) t = await getTranslations({locale: locale || 'en', namespace: "ProductsPage"});
            const categoryName = category ? t(`category_${category}_label`, {}, {defaultValue: category}) : t('fallbackCategoryName');
            errorTitle = t('globalErrorTitleCategory', {categoryName});
            errorMessage = t('globalErrorMessageCategory', {categoryName});
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