// src/components/LazyNATSProvider.jsx
"use client"
import React, {lazy, Suspense, memo} from 'react';
const NATSProviderLazy = lazy(() =>
    import('../../context/NATSContext').then((module) => ({default: module.NATSProvider}))
);
const LazyNATSProvider = memo(function LazyNATSProvider({children}) {
    return (
        <Suspense fallback={<div>Loading NATS...</div>}>
            <NATSProviderLazy>
                {children}
            </NATSProviderLazy>
        </Suspense>
    );
});
export default LazyNATSProvider;
