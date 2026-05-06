// File: src/app/videos/page.jsx
import React from "react";
import {getAllVideos} from "../../../api/mediaApi";
import VideosPageClient from "./VideosPage.client";
import { safeSerializeJsonLd } from '@/utils/secureJsonLd';

export const revalidate = 60;

export const metadata = {
    title: "Short Video Feed",
    description: "A vertical feed of short videos.",
    // etc.
};

export default async function VideosIndexPage() {
    let videoData = {videos: []};
    let errorMessage = "";

    try {
        videoData = await getAllVideos();
    } catch (err) {
        // Error: "Error fetching videos on server:", err...
        errorMessage = "Error fetching videos. Please try again.";
    }

    const serverVideos = videoData?.videos || [];

    // If needed, build your JSON-LD...
    const videosJsonLd = { /* ... same logic as before ... */};

    return (
        <>
            <script
                type="application/ld+json"
                dangerouslySetInnerHTML={{__html: safeSerializeJsonLd(videosJsonLd)}}
            />
            <VideosPageClient
                serverVideos={serverVideos}
                errorMessage={errorMessage}
            />
        </>
    );
}
