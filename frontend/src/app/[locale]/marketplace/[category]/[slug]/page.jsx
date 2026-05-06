import React from 'react';
import {notFound} from 'next/navigation';
import {getProduct} from "../../../../../api/searchApi";
import ClassifiedCard from "../../../../../components/classified/ClassifiedCard";
import styles from "../../../../../components/classified/ClassifiedCard.module.css";
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
    const {slug} = params;
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
    return {
        title: product.name || 'Product Details',
        description: product.description?.substring(0, 160) || `Details for product ${slug}`,
        // Add OpenGraph, canonical URL etc.
        // alternates: {
        //    canonical: `/marketplace/${params.category}/${slug}`,
        // },
    };
}
// Define the Page component (Server Component)
export default async function ProductDetailPage(props) {
    // First safely resolve the props
    const {params} = await resolveProps(props);
    // Now safely destructure params
    const {slug, locale, category} = params;
    let productData = null;
    try {
        // Fetch the specific product data using the slug
        productData = await getProduct(slug);
        // Handle case where product is not found by the API
        if (!productData) {
            notFound();
        }
    } catch (err) {
        notFound();
    }
    // Check if styles object exists before accessing properties
    const containerClass = styles?.detailPageContainer || 'container mx-auto p-4';
    return (
        <div className={containerClass}>
            <ClassifiedCard product={productData}/>
        </div>
    );
}