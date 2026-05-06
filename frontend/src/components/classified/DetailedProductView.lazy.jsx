import React, { lazy, Suspense } from 'react';

const DetailedProductViewLazy = lazy(() => import('./DetailedProductView.jsx'));

export default function DetailedProductView(props) {
    return (
        <Suspense fallback={
            <div style={{ 
                padding: '2rem', 
                textAlign: 'center',
                minHeight: '200px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
            }}>
                <div className="loading-spinner">Loading...</div>
            </div>
        }>
            <DetailedProductViewLazy {...props} />
        </Suspense>
    );
}
