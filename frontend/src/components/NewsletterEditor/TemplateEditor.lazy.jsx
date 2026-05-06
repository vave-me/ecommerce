import React, { lazy, Suspense } from 'react';

const TemplateEditorLazy = lazy(() => import('./TemplateEditor.jsx'));

export default function TemplateEditor(props) {
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
            <TemplateEditorLazy {...props} />
        </Suspense>
    );
}
