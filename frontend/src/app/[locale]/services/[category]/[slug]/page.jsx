// File: src/app/[locale]/services/[category]/[slug]/page.jsx
// SERVER component responsible for fetching data for a specific service
import React from 'react';
import {notFound} from 'next/navigation';
import {getService} from "../../../../../api/searchApi";
import DetailedServiceView from "../../../../../components/services/DetailedServiceView";

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
    const {slug, category} = params;

    let service = null;
    try {
        service = await getService(slug);
    } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }

    if (!service) {
        return {
            title: 'Service Not Found | Vaveme',
            description: 'The requested service could not be found.'
        };
    }

    return {
        title: `${service.name || 'Service Details'} | Vaveme`,
        description: service.description?.substring(0, 160) || `Details for ${service.name || 'service'} in ${category} category`,
        openGraph: {
            title: service.name || 'Service Details',
            description: service.description?.substring(0, 160) || '',
            images: service.thumbnail ? [service.thumbnail] : [],
        },
        // Add structured data for better SEO
        alternates: {
            canonical: `/services/${category}/${slug}`,
        },
    };
}

// Define the Page component (Server Component)
export default async function ServiceDetailPage(props) {
    // First safely resolve the props
    const {params} = await resolveProps(props);

    // Now safely destructure params
    const {slug, locale, category} = params;

    // Validate required parameters
    if (!slug) {
        // Error: 'ServiceDetailPage: Service slug is missing!'...
        notFound();
    }

    if (!category) {
        // Error: 'ServiceDetailPage: Category is missing!'...
        notFound();
    }

    let serviceData = null;

    try {
        // Fetch the specific service data using the slug from search API (cached proxy)
        serviceData = await getService(slug);

        // Handle case where service is not found by the API
        if (!serviceData) {
            notFound();
        }
    } catch (err) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', err);
        }
        throw err; // Re-throw for caller to handle
    
        notFound();
    }

    return (
        <DetailedServiceView 
            service={serviceData} 
            locale={locale}
            category={category}
        />
    );
}