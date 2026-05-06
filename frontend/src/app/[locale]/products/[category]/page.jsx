/* ──────────────────────────────────────────────────────────────
   Products category page – locale–aware server component
────────────────────────────────────────────────────────────── */
import React from "react";
import {getTranslations} from "next-intl/server";
import {notFound} from "next/navigation";
import {searchProductsWithFilters, searchProductsWithCategorySlug} from "../../../../api/searchApi";
import {fetchMainCategories} from "../../../../api/categories";
import ProductsPageClient from "../ProductsPage.client";
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';
import styles from "../../../../components/classified/ClassifiedCard.module.css";
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
        // Get translations for better titles
        const t = await getTranslations({locale, namespace: 'ProductsPage'});
        return {
            title: t('categoryMetaTitle', { categoryName: category, defaultValue: `${category} Products` }),
            description: t('categoryMetaDescription', { categoryName: category, defaultValue: `Find the best products in the ${category} category.` })
        };
    } catch (e) {
        return {
            title: 'Category Products',
            description: 'Find products by category.'
        };
    }
}
export default async function ProductsCategoryPage(props) {
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
        const pageSize = 20;
        const sortBy = searchParams?.sortBy || '';
        // 4. Create listing filters
        const listingFilters = {
            category,
            displayMode,
            page,
            pageSize,
            sortBy,
        };
        // 5. Load translations and validate category
        const t = await getTranslations({locale, namespace: 'ProductsPage'});
        // 6. Fetch products and categories
        let fetchedData = {products: [], totalCount: 0, totalPages: 0};
        let fetchError = null;
        let categories = [];
        let validCategory = null;
        try {
            // Fetch categories for validation and navigation
            const categoriesResponse = await fetchMainCategories({ 
                categoryType: 'marketplace', 
                lang: locale 
            });
            categories = categoriesResponse?.categories || [];
            // Validate category exists
            validCategory = categories.find(cat => 
                cat.slug === category || cat.name?.toLowerCase() === category.toLowerCase()
            );
            if (!validCategory) {
                // Still try to fetch products, but mark as potentially invalid
            }
            // Try specific category API first, then fallback to filtered search
            try {
                
                fetchedData = await searchProductsWithCategorySlug(category, listingFilters);
                
            } catch (slugError) {

                fetchedData = await searchProductsWithFilters(listingFilters);
                
            }
        } catch (err) {
            fetchError = t('categoryLoadError', {
                defaultValue: "Could not load products for this category."
            });
        }
        // 7. Process the fetched data
        const serverProducts = fetchedData?.products || [];
        const totalPages = parseInt(fetchedData?.totalPages || '0', 10);
        const totalCount = parseInt(fetchedData?.totalCount || '0', 10);
        const currentPage = parseInt(fetchedData?.currentPage || '1', 10);
        // 8. Construct JSON-LD for SEO
        const categoryName = validCategory?.name || category;
        const productsJsonLd = {
            '@context': 'https://schema.org',
            '@type': 'ItemList',
            'name': `Products in ${categoryName}`,
            'description': `Best products in the ${categoryName} category.`,
            itemListElement: serverProducts.map((product, index) => {
                return {
                    '@type': 'ListItem',
                    position: index + 1,
                    url: `https://www.sfx-market.de/products/${category}/${product.slug || product.id}`,
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
            }).filter(Boolean),
        };
        // 9. Return UI components
        return (
            <>
                {productsJsonLd.itemListElement.length > 0 && (
                    <script
                        type="application/ld+json"
                        dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(productsJsonLd)}}
                    />
                )}
                <ProductsPageClient
                    serverProducts={serverProducts}
                    serverFilters={listingFilters}
                    availableCategories={categories}
                    fetchError={fetchError}
                    currentCategory={category}
                    categoryName={categoryName}
                    totalPages={totalPages}
                    currentPage={currentPage}
                    totalCount={totalCount}
                    locale={locale}
                    labels={{
                        empty: t('categoryEmpty', {categoryName, defaultValue: `No products found in ${categoryName}`})
                    }}
                />
            </>
        );
    } catch (e) {
        // Global error handler
        return (
            <div className={styles.errorContainer}>
                <h1 className={styles.errorTitle}>Error Loading Products</h1>
                <p className={styles.errorMessage}>We couldn't load this category. Please try again later.</p>
            </div>
        );
    }
}