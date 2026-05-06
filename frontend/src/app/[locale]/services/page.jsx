/* ──────────────────────────────────────────────────────────────
   Services page – locale–aware server component
────────────────────────────────────────────────────────────── */

import React from "react";
import {getTranslations} from "next-intl/server";
import {notFound} from "next/navigation";

import {searchServicesWithFilters} from "../../../api/searchApi";
import ServicesPageClient from "./ServicesPage.client";
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
    // Safely resolve props first
    const { params } = await resolveProps(props);
    const { locale } = params;

    // Get translations
    const t = await getTranslations({locale, namespace: "ServicesPage"});

    return {
        title: t("metaTitle", { defaultValue: "Professional Services | Vaveme" }),
        description: t("metaDescription", { defaultValue: "Find trusted service providers for all your needs on our platform." }),
    };
}

export default async function ServicesIndexPage(props) {
    // First safely resolve the props
    const { params, searchParams } = await resolveProps(props);

    // Now it's safe to destructure params
    const { locale } = params;

    /* 1) i18n helper ------------------------------------------ */
    const t = await getTranslations({locale, namespace: "ServicesPage"});

    /* 2) optional filters (extract from searchParams if needed) */
    const displayMode = searchParams?.displayMode || "grid";
    const page = parseInt(searchParams?.page || "1", 10);
    const limit = 20;
    const sortBy = searchParams?.sortBy || "";
    const sortOrder = searchParams?.sortOrder || "asc";

    const listingFilters = {
        displayMode,
        page,
        limit,
        sortBy,
        sortOrder,
        // Add other filters as needed
    };

    /* 3) fetch services using search API (cached proxy) */
    let servicesData = {services: [], totalCount: 0, totalPages: 0};
    let fetchError = null;

    try {
        servicesData = await searchServicesWithFilters(listingFilters);
    } catch (err) {
        // Error: "Error fetching services:", err...
        fetchError = t("errorFetchingServices", {
            defaultValue: "Error loading services. Please try again."
        });
    }

    const services = servicesData.services ?? [];
    const totalPages = parseInt(servicesData?.totalPages || '0', 10);
    const totalCount = parseInt(servicesData?.totalCount || '0', 10);
    const currentPage = parseInt(servicesData?.currentPage || '1', 10);

    /* 4) 404 if nothing found - optional based on requirements */
    // Note: You might want to show an empty state instead of 404
    // if (services.length === 0 && !fetchError) notFound();

    /* 5) JSON-LD for SEO with correct service URLs */
    const jsonLd = {
        "@context": "https://schema.org",
        "@type": "ItemList",
        "name": "Professional Services",
        "description": "Find trusted service providers for all your needs.",
        itemListElement: services.map((s, i) => ({
            "@type": "ListItem",
            position: i + 1,
            url: `https://www.sfx-markt.de/${locale}/services/${s.categorySlug || 'general'}/${s.id}`,
            name: s.name,
            item: {
                "@type": "Service",
                name: s.name,
                image: s.thumbnail ?? "",
                description: s.description?.substring(0, 150) || "",
                provider: {
                    "@type": "Organization",
                    name: s.providerName ?? "",
                },
                serviceType: s.category ?? "",
                areaServed: s.location ?? "",
                offers: {
                    "@type": "Offer",
                    price: s.basePrice ?? "0.00",
                    priceCurrency: "EUR",
                    availability: "https://schema.org/InStock"
                }
            }
        }))
    };

    /* 6) render */
    return (
        <>
            <script
                type="application/ld+json"
                dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(jsonLd)}}
            />

            <ServicesPageClient
                serverServices={services}
                serverFilters={listingFilters}
                pageTitle={t("title")}
                emptyMsg={t("empty")}
                fetchError={fetchError}
                totalPages={totalPages}
                totalCount={totalCount}
                currentPage={currentPage}
            />
        </>
    );
}