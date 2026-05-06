import React, { lazy, Suspense } from 'react';

const OffersModalFormEditableLazy = lazy(() => import('./OffersModalFormEditable.jsx'));

export default function OffersModalFormEditable(props) {
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
            <OffersModalFormEditableLazy {...props} />
        </Suspense>
    );
}
