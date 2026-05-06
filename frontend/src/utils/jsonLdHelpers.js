import { safeSerializeJsonLd as secureSafeSerializeJsonLd } from './secureJsonLd.js';
/**
 * Safe JSON-LD serialization wrapper
 * @param {Object} jsonLdObject - JSON-LD object to serialize
 * @returns {string} - Safely serialized JSON-LD string
 */
export function safeSerializeJsonLd(jsonLdObject) {
    return secureSafeSerializeJsonLd(jsonLdObject);
}
/**
 * JSON-LD Helpers for generating structured data
 * Uses secure serialization to prevent XSS vulnerabilities
 */
/**
 * Generate JSON-LD for a listing/product
 * @param {Object} listing - Listing data
 * @param {string} locale - Current locale
 * @returns {Object} - JSON-LD structured data
 */
export function generateListingJsonLd(listing, locale = 'en') {
    if (!listing) return null;
    const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || 'https://example.com';
    const currency = locale === 'pl' ? 'PLN' : 'EUR';
    return {
        '@context': 'https://schema.org',
        '@type': 'Product',
        '@id': `${baseUrl}/${locale}/catalog/${listing.id}`,
        name: listing.title,
        description: listing.description,
        url: `${baseUrl}/${locale}/listings/${listing.id}`,
        image: listing.images?.map(img => img.url) || [],
        brand: listing.brand || undefined,
        category: listing.category,
        offers: {
            '@type': 'Offer',
            price: listing.price,
            priceCurrency: currency,
            availability: listing.available ? 'https://schema.org/InStock' : 'https://schema.org/OutOfStock',
            seller: {
                '@type': 'Person',
                name: listing.seller?.name || 'Anonymous',
                url: listing.seller?.profileUrl || undefined
            },
            validFrom: listing.createdAt,
            validThrough: listing.expiresAt || undefined
        },
        aggregateRating: listing.rating ? {
            '@type': 'AggregateRating',
            ratingValue: listing.rating.average,
            reviewCount: listing.rating.count,
            bestRating: 5,
            worstRating: 1
        } : undefined,
        location: listing.location ? {
            '@type': 'Place',
            name: listing.location,
            address: {
                '@type': 'PostalAddress',
                addressLocality: listing.location
            }
        } : undefined,
        datePublished: listing.createdAt,
        dateModified: listing.updatedAt
    };
}
/**
 * Generate JSON-LD for listings page
 * @param {Array} listings - Array of listings
 * @param {string} locale - Current locale
 * @param {Object} pageInfo - Page information
 * @returns {Object} - JSON-LD structured data
 */
export function generateListingsPageJsonLd(listings = [], locale = 'en', pageInfo = {}) {
    const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || 'https://example.com';
    return {
        '@context': 'https://schema.org',
        '@type': 'CollectionPage',
        '@id': `${baseUrl}/${locale}/catalog`,
        name: pageInfo.title || 'Listings',
        description: pageInfo.description || 'Browse our collection of listings',
        url: `${baseUrl}/${locale}/catalog`,
        mainEntity: {
            '@type': 'ItemList',
            numberOfItems: listings.length,
            itemListElement: listings.map((listing, index) => ({
                '@type': 'ListItem',
                position: index + 1,
                item: generateListingJsonLd(listing, locale)
            }))
        },
        breadcrumb: {
            '@type': 'BreadcrumbList',
            itemListElement: [
                {
                    '@type': 'ListItem',
                    position: 1,
                    name: 'Home',
                    item: `${baseUrl}/${locale}`
                },
                {
                    '@type': 'ListItem',
                    position: 2,
                    name: 'Listings',
                    item: `${baseUrl}/${locale}/catalog`
                }
            ]
        }
    };
}
/**
 * Generate JSON-LD for organization/website
 * @param {string} locale - Current locale
 * @returns {Object} - JSON-LD structured data
 */
export function generateOrganizationJsonLd(locale = 'en') {
    const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || 'https://example.com';
    return {
        '@context': 'https://schema.org',
        '@type': 'Organization',
        '@id': `${baseUrl}/#organization`,
        name: 'Classified Marketplace',
        url: baseUrl,
        logo: `${baseUrl}/logo.png`,
        description: 'Your trusted marketplace for buying and selling',
        contactPoint: {
            '@type': 'ContactPoint',
            telephone: '+1-555-123-4567',
            contactType: 'customer service',
            availableLanguage: ['en', 'pl']
        },
        sameAs: [
            'https://facebook.com/yourpage',
            'https://twitter.com/yourhandle',
            'https://instagram.com/yourhandle'
        ]
    };
}
/**
 * Generate JSON-LD for website
 * @param {string} locale - Current locale
 * @returns {Object} - JSON-LD structured data
 */
export function generateWebsiteJsonLd(locale = 'en') {
    const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || 'https://example.com';
    return {
        '@context': 'https://schema.org',
        '@type': 'WebSite',
        '@id': `${baseUrl}/#website`,
        name: 'Classified Marketplace',
        url: baseUrl,
        description: 'Your trusted marketplace for buying and selling',
        publisher: {
            '@id': `${baseUrl}/#organization`
        },
        potentialAction: {
            '@type': 'SearchAction',
            target: {
                '@type': 'EntryPoint',
                urlTemplate: `${baseUrl}/${locale}/search?q={search_term_string}`
            },
            'query-input': 'required name=search_term_string'
        },
        inLanguage: locale
    };
}
/**
 * Generate JSON-LD for breadcrumbs
 * @param {Array} breadcrumbs - Array of breadcrumb items
 * @param {string} locale - Current locale
 * @returns {Object} - JSON-LD structured data
 */
export function generateBreadcrumbJsonLd(breadcrumbs = [], locale = 'en') {
    const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || 'https://example.com';
    return {
        '@context': 'https://schema.org',
        '@type': 'BreadcrumbList',
        itemListElement: breadcrumbs.map((crumb, index) => ({
            '@type': 'ListItem',
            position: index + 1,
            name: crumb.name,
            item: crumb.url.startsWith('http') ? crumb.url : `${baseUrl}${crumb.url}`
        }))
    };
}
/**
 * Generate JSON-LD for FAQ page
 * @param {Array} faqs - Array of FAQ items
 * @returns {Object} - JSON-LD structured data
 */
export function generateFaqJsonLd(faqs = []) {
    return {
        '@context': 'https://schema.org',
        '@type': 'FAQPage',
        mainEntity: faqs.map(faq => ({
            '@type': 'Question',
            name: faq.question,
            acceptedAnswer: {
                '@type': 'Answer',
                text: faq.answer
            }
        }))
    };
}
/**
 * Generate JSON-LD for local business
 * @param {Object} business - Business information
 * @returns {Object} - JSON-LD structured data
 */
export function generateLocalBusinessJsonLd(business = {}) {
    return {
        '@context': 'https://schema.org',
        '@type': 'LocalBusiness',
        name: business.name,
        description: business.description,
        url: business.url,
        telephone: business.phone,
        email: business.email,
        address: {
            '@type': 'PostalAddress',
            streetAddress: business.address?.street,
            addressLocality: business.address?.city,
            addressRegion: business.address?.region,
            postalCode: business.address?.postalCode,
            addressCountry: business.address?.country
        },
        geo: business.coordinates ? {
            '@type': 'GeoCoordinates',
            latitude: business.coordinates.lat,
            longitude: business.coordinates.lng
        } : undefined,
        openingHours: business.openingHours,
        priceRange: business.priceRange
    };
}
/**
 * Safely serialize and return JSON-LD string
 * @param {Object} jsonLdData - JSON-LD object
 * @returns {string} - Safe JSON-LD string
 */
export function safeJsonLdString(jsonLdData) {
    return safeSerializeJsonLd(jsonLdData);
}
/**
 * Combine multiple JSON-LD objects into an array
 * @param {...Object} jsonLdObjects - Multiple JSON-LD objects
 * @returns {Array} - Array of JSON-LD objects
 */
export function combineJsonLd(...jsonLdObjects) {
    return jsonLdObjects.filter(obj => obj && typeof obj === 'object');
}
/**
 * Generate complete page JSON-LD with multiple schemas
 * @param {Object} options - Options for generating JSON-LD
 * @returns {Array} - Array of JSON-LD objects
 */
export function generatePageJsonLd(options = {}) {
    const {
        type = 'webpage',
        listing = null,
        listings = [],
        breadcrumbs = [],
        faqs = [],
        business = null,
        locale = 'en',
        pageInfo = {}
    } = options;
    const jsonLdObjects = [];
    // Always include organization and website
    jsonLdObjects.push(generateOrganizationJsonLd(locale));
    jsonLdObjects.push(generateWebsiteJsonLd(locale));
    // Add breadcrumbs if provided
    if (breadcrumbs.length > 0) {
        jsonLdObjects.push(generateBreadcrumbJsonLd(breadcrumbs, locale));
    }
    // Add specific content based on type
    switch (type) {
        case 'listing':
            if (listing) {
                jsonLdObjects.push(generateListingJsonLd(listing, locale));
            }
            break;
        case 'listings':
            jsonLdObjects.push(generateListingsPageJsonLd(listings, locale, pageInfo));
            break;
        case 'faq':
            if (faqs.length > 0) {
                jsonLdObjects.push(generateFaqJsonLd(faqs));
            }
            break;
        case 'business':
            if (business) {
                jsonLdObjects.push(generateLocalBusinessJsonLd(business));
            }
            break;
    }
    return jsonLdObjects.filter(obj => obj !== null);
} 