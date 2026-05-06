import React from 'react';
import {notFound} from 'next/navigation';
import {getTranslations} from 'next-intl/server';
import {getProduct} from "../../../../../api/searchApi";
import {fetchMainCategories} from "../../../../../api/categories";
import DetailedProductView from "../../../../../components/classified/DetailedProductView.lazy";
// Helper function to safely resolve props in Next.js 15+
async function resolveProps(props) {
    // In Next.js 15+, we need to await both params and searchParams
    const params = props.params ? await props.params : {};
    const searchParams = props.searchParams ? await props.searchParams : {};
    return {params, searchParams};
}
// Generate dynamic metadata for SEO
export async function generateMetadata(props) {
    // First safely resolve the props
    const {params} = await resolveProps(props);
    const {slug, locale, category} = params;
    let product = null;
    try {
        product = await getProduct(slug);
    } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }
    if (!product) {
        return {title: 'Product Not Found'};
    }
    const t = await getTranslations({locale, namespace: 'ProductsPage'});
    // Get category name for translations
    let categoryName = category;
    try {
        const categories = await fetchMainCategories({ categoryType: 'marketplace', lang: locale });
        const foundCategory = categories?.categories?.find(cat => 
            cat.slug === category || cat.name?.toLowerCase() === category.toLowerCase()
        );
        categoryName = foundCategory?.name || category;
    } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }
    return {
        title: product.name || t('productDetailTitle', {
            productName: product.name || slug,
            defaultValue: 'Product Details'
        }),
        description: product.description?.substring(0, 160) || t('productDetailDescription', {
            productName: product.name || slug,
            categoryName: categoryName,
            slug,
            defaultValue: `Details for product ${slug}`
        }),
        alternates: {
            canonical: `/products/${category}/${slug}`,
        },
    };
}
// Define the Page component (Server Component)
export default async function ProductDetailPage(props) {
    // First safely resolve the props
    const {params} = await resolveProps(props);
    // Now safely destructure params
    const {slug, locale, category} = params;
    let productData = null;
    let categories = [];
    let validCategory = null;
    try {
        // Fetch both product data and categories for validation
        const [productResponse, categoriesResponse] = await Promise.all([
            getProduct(slug),
            fetchMainCategories({ categoryType: 'marketplace', lang: locale })
        ]);
        productData = productResponse;
        categories = categoriesResponse?.categories || [];
        // Validate category exists
        validCategory = categories.find(cat => 
            cat.slug === category || cat.name?.toLowerCase() === category.toLowerCase()
        );
        if (!validCategory) {
        }
        // Handle case where product is not found by the API
        if (!productData) {
            notFound();
        }
        // Optional: Validate that product belongs to the specified category
        if (productData.categorySlug && productData.categorySlug !== category) {
        }
    } catch (err) {
        notFound();
    }
    return (
        <DetailedProductView 
            product={productData} 
            locale={locale} 
            category={validCategory}
            availableCategories={categories}
        />
    );
}