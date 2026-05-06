/* ──────────────────────────────────────────────────────────────
   Marketplace page – locale–aware server component
────────────────────────────────────────────────────────────── */
import React from "react";
import {getTranslations} from "next-intl/server";
import {notFound} from "next/navigation";
import {fetchProductsByFilters} from "../../../api/productsApi";
import MarketplacePageClient from "./MarketplacePage.client";
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
    const t = await getTranslations({locale, namespace: "MarketplacePage"});
    return {
        title: t("metaTitle", { defaultValue: "Marketplace | Vaveme" }),
        description: t("metaDescription", { defaultValue: "Browse our selection of high-quality products." }),
    };
}
export default async function MarketplaceIndexPage(props) {
    // First safely resolve the props
    const { params, searchParams } = await resolveProps(props);
    // Now it's safe to destructure params
    const { locale } = params;
    /* 1) i18n helper ------------------------------------------ */
    const t = await getTranslations({locale, namespace: "MarketplacePage"});
    /* 2) optional filters (extract from searchParams if needed) */
    const displayMode = searchParams?.displayMode || "list";
    const page = parseInt(searchParams?.page || "1", 10);
    const listingFilters = {
        displayMode,
        page,
        // Add other filters as needed
    };
    /* 3) fetch products */
    let productsData = {products: [], totalCount: 0, totalPages: 0};
    let fetchError = null;
    try {
        productsData = await fetchProductsByFilters(listingFilters);
    } catch (err) {
        fetchError = t("errorFetchingProducts", {
            defaultValue: "Error loading products. Please try again."
        });
    }
    const products = productsData.products ?? [];
    const totalPages = parseInt(productsData?.totalPages || '0', 10);
    const totalCount = parseInt(productsData?.totalCount || '0', 10);
    const currentPage = parseInt(productsData?.currentPage || '1', 10);
    /* 4) 404 if nothing found - optional based on requirements */
    // Note: You might want to show an empty state instead of 404
    // if (products.length === 0 && !fetchError) notFound();
    /* 5) JSON-LD for SEO */
    const jsonLd = {
        "@context": "https://schema.org",
        "@type": "ItemList",
        itemListElement: products.map((p, i) => ({
            "@type": "ListItem",
            position: i + 1,
            url: `https://www.sfx-market.de/product/${p.id}`,
            name: p.name,
            item: {
                "@type": "Product",
                name: p.name,
                image: p.thumbnail ?? "",
                description: p.description,
                brand: p.brand ?? "",
                mpn: p.mpn || p.sku || p.model || "",
                itemCondition:
                    p.condition === "used"
                        ? "https://schema.org/UsedCondition"
                        : "https://schema.org/NewCondition",
                offers: {
                    "@type": "Offer",
                    price: p.basePrice ?? "0.00",
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
            <MarketplacePageClient
                serverProducts={products}
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