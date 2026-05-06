/* ──────────────────────────────────────────────────────────────
   Marketplace category page – locale–aware server component
────────────────────────────────────────────────────────────── */
import React from "react";
import {getTranslations} from "next-intl/server";
import {notFound} from "next/navigation";
import {fetchProductsByFilters} from "../../../../api/productsApi";
import MarketplacePageClient from "../MarketplacePage.client";
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';
export const dynamic = "force-dynamic"; // keep reading params.slug etc.
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
        // Optional: Get translations for better titles
        // const t = await getTranslations({locale, namespace: 'MarketplacePage'});
        // const categoryName = category;
        return {
            title: `Products in ${category}`,
            description: `Find the best products in the ${category} category.`
        };
    } catch (e) {
        return {
            title: 'Category Products',
            description: 'Find products by category.'
        };
    }
}
export default async function MarketplaceCategoryPage(props) {
    try {
        // 1. Safely resolve props first
        const { params, searchParams } = await resolveProps(props);
        // 2. Now it's safe to destructure params
        const { locale, category } = params;
        // Validate category early
        if (!category) {
            notFound();
        }
        // 3. Extract searchParams values
        const displayMode = searchParams?.displayMode || 'grid';
        const page = parseInt(searchParams?.page || '1', 10);
        const limit = 20;
        const sortBy = searchParams?.sortBy || '';
        // 4. Create listing filters
        const listingFilters = {
            category,
            displayMode,
            page,
            limit,
            sortBy,
        };
        // 5. Optional: Load translations if needed
        // const t = await getTranslations({locale, namespace: 'MarketplacePage'});
        // 6. Fetch products filtered by category
        let fetchedData = {products: [], totalCount: 0, totalPages: 0};
        let fetchError = null;
        try {
            // Use an API function that accepts category filter
            fetchedData = await fetchProductsByFilters(listingFilters);
            // Optional: Check if category exists / is valid, call notFound() if not
        } catch (err) {
            fetchError = "Could not load products for this category.";
        }
        // 7. Process the fetched data
        const serverProducts = fetchedData?.products || [];
        const totalPages = parseInt(fetchedData?.totalPages || '0', 10);
        const totalCount = parseInt(fetchedData?.totalCount || '0', 10);
        const currentPage = parseInt(fetchedData?.currentPage || '1', 10);
        // 8. Construct JSON-LD for SEO
        const productsJsonLd = {
            '@context': 'https://schema.org',
            '@type': 'ItemList',
            'name': `Products in ${category}`,
            'description': `Best products in the ${category} category.`,
            itemListElement: serverProducts.map((product, index) => {
                return {
                    '@type': 'ListItem',
                    position: index + 1,
                    url: `https://www.sfx-market.de/marketplace/${category}/${product.id}`,
                    item: {
                        '@type': 'Product',
                        name: product.name,
                        image: product.thumbnail ?? '',
                        description: product.description?.substring(0, 150) || '',
                        brand: product.brand ?? '',
                        mpn: product.mpn || product.sku || product.model || '',
                        itemCondition: product.condition === 'used'
                            ? 'https://schema.org/UsedCondition'
                            : 'https://schema.org/NewCondition',
                        offers: {
                            '@type': 'Offer',
                            price: product.basePrice ?? '0.00',
                            priceCurrency: 'EUR',
                            availability: 'https://schema.org/InStock',
                        },
                    },
                };
            }),
        };
        // 9. Return UI components
        return (
            <>
                <script
                    type="application/ld+json"
                    dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(productsJsonLd)}}
                />
                <MarketplacePageClient
                    serverProducts={serverProducts}
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
        return (
            <div className="container mx-auto p-8 text-center">
                <h1 className="text-2xl font-bold mb-4">Error Loading Products</h1>
                <p>We couldn't load this category. Please try again later.</p>
            </div>
        );
    }
}