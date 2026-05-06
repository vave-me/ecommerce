import React, { lazy, Suspense } from 'react';

const PromptInputLazy = lazy(() => import('./PromptInput.jsx'));

export default function PromptInput(props) {
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
            <PromptInputLazy {...props} />
        </Suspense>
    );
}
