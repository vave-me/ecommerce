"use client";
import React, { memo } from 'react';
import Script from 'next/script';
const AdSenseScript = memo(function AdSenseScript() {
    return (
        <Script
            async
            src="https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-YOUR_PUBLISHER_ID"
            crossOrigin="anonymous"
            strategy="afterInteractive"
        />
    );
});
export default AdSenseScript; 