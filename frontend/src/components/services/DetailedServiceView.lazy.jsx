import React, { lazy, Suspense } from 'react';

const DetailedServiceViewLazy = lazy(() => import('./DetailedServiceView.jsx'));

export default function DetailedServiceView(props) {
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
            <DetailedServiceViewLazy {...props} />
        </Suspense>
    );
}
