import React from 'react';
import {notFound} from 'next/navigation';
import {getTranslations} from 'next-intl/server'; // Import
import {getProductsBySlug} from "../../../../../api/productsApi"; // Adjust path
import ProductCard from "./ProductCard"; // Fix import path to local ProductCard
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';
import styles from './page.module.css'; // Optional styles
// Helper function (assume defined)
async function resolveProps(props) {
    const params = props.params ? await props.params : {};
    const searchParams = props.searchParams ? await props.searchParams : {};
    return {params, searchParams};
}
export async function generateMetadata(props) {
    let t;
    let slug;
    let category;
    let product = null;
    try {
        const {params} = await resolveProps(props);
        const {locale} = params;
        slug = params.slug;
        category = params.category;
        t = await getTranslations({locale, namespace: 'ProductDetailPage'});
        product = await getProductsBySlug(slug);
        if (!product) {
            return {title: t('metaTitleNotFound')};
        }
        const productName = product.name || t('metaTitleDefault');
        // Attempt to get translated category name
        const categoryName = category ? t(`ProductsPage.category_${category}_label`, {}, {defaultValue: category}) : t('fallbackCategoryName');
        const description = product.description?.substring(0, 160) || t('metaDescriptionDefault', {
            productName: product.name,
            categoryName
        });
        // Map availability for OpenGraph
        const ogAvailability = product.inStock ? t('ogAvailabilityInStock') : t('ogAvailabilityOutOfStock');
        return {
            title: productName,
            description: description,
            openGraph: {
                title: productName,
                description: description,
                type: 'product', // Correct OG type
                images: product.images?.map(img => ({url: img})) || (product.thumbnail ? [{url: product.thumbnail}] : []),
                // Add price and currency for OG
                price: {
                    amount: String(product.basePrice || '0.00'),
                    currency: t('common.currencyCode')
                },
                availability: ogAvailability, // Use mapped availability
                // Potentially add brand, category, etc.
            },
            // 'other' metadata is less standard, ensure consumers support it
            other: {
                'product:price:amount': String(product.basePrice || '0.00'),
                'product:price:currency': t('common.currencyCode'),
                'product:availability': ogAvailability, // Use mapped availability
            },
            alternates: {
                canonical: `/${locale}/products/${category}/${slug}`,
            },
        };
    } catch (err) {
        const errorTitle = t ? t('metaTitleError') : "Error Loading Product";
        return {title: errorTitle};
    }
}
export default async function ProductDetailPage(props) {
    let t; // For potential use in catch block
    let slug; // For potential use in catch block
    let category; // For potential use in catch block
    let locale; // For potential use in catch block
    try {
        const {params} = await resolveProps(props);
        locale = params.locale; // Assign locale
        slug = params.slug; // Assign slug
        category = params.category; // Assign category
        t = await getTranslations({locale, namespace: 'ProductDetailPage'});
        let productData = null;
        try {
            productData = await getProductsBySlug(slug);
            if (!productData) {
                notFound();
            }
        } catch (err) {
            notFound();
        }
        const containerClass = styles?.detailPageContainer || 'container mx-auto p-4'; // Use optional chaining
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
        const itemCondition = conditionMap[productData.condition] || conditionMap['new'];
        const availability = availabilityMap[productData.availability] || availabilityMap['inStock'];
        const productUrl = `https://www.sfx-market.de/${locale}/products/${category}/${slug}`;
        const productJsonLd = {
            '@context': 'https://schema.org',
            '@type': 'Product',
            'name': productData.name || t('fallbackProductName'), // Add fallback name
            'image': productData.images?.[0] || productData.thumbnail || '', // Prefer first image or thumbnail
            'description': productData.description || '',
            'sku': productData.sku || undefined, // Use undefined if empty
            'mpn': productData.mpn || undefined, // Use undefined if empty
            'brand': {
                '@type': 'Brand',
                'name': productData.brand || t('common.defaultBrand') // Use common fallback
            },
            'offers': {
                '@type': 'Offer',
                'url': productUrl,
                'price': String(productData.basePrice || productData.price || '0.00'),
                'priceCurrency': t('common.currencyCode'),
                'availability': availability, // Use mapped availability
                'itemCondition': itemCondition, // Use mapped condition
                // Optionally add seller information if available
                // 'seller': {
                //    '@type': 'Organization', // or 'Person'
                //    'name': productData.seller?.name || 'Vaveme Marketplace'
                // }
            },
            // Add reviews if available
            // 'aggregateRating': {
            //   '@type': 'AggregateRating',
            //   'ratingValue': productData.rating || '4.5',
            //   'reviewCount': productData.reviewCount || '10'
            // }
        };
        return (
            <div className={containerClass}>
                            <script 
              type="application/ld+json"
              dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(productJsonLd)}}
            />
                {/* Pass product data and locale to the client card component */}
                <ProductCard product={productData} locale={locale}/>
            </div>
        );
    } catch (e) {
        // Attempt to get translations for error page
        let errorTitle = "Error Loading Product";
        let errorMessage = "Could not load the product details. Please try again later.";
        try {
            if (!t) t = await getTranslations({locale: locale || 'en', namespace: "ProductDetailPage"});
            errorTitle = t('globalErrorTitle');
            errorMessage = t('globalErrorMessage');
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