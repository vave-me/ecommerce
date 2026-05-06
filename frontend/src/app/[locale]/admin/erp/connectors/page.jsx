import { Suspense } from 'react';
import { getTranslations } from 'next-intl/server';
import ConnectorsManagement from './ConnectorsManagement.client';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';

/**
 * ERP Connectors Management Page
 * Manage ERP system connectors and configurations
 */
export default async function ConnectorsPage() {
  const t = await getTranslations('ERPConnectors');

  return (
    <div className="admin-page">
      <div className="page-header">
        <h1>{t('title', { defaultValue: 'ERP Connectors' })}</h1>
        <p>{t('subtitle', { defaultValue: 'Manage ERP system connectors and integrations' })}</p>
      </div>

      <Suspense 
        fallback={
          <div className="loading-container">
            <LoadingSpinner />
            <p>{t('loading', { defaultValue: 'Loading connectors...' })}</p>
          </div>
        }
      >
        <ConnectorsManagement />
      </Suspense>
    </div>
  );
}

export async function generateMetadata({ params }) {
  const t = await getTranslations('ERPConnectors');
  
  return {
    title: t('metaTitle', { defaultValue: 'ERP Connectors Management' }),
    description: t('metaDescription', { defaultValue: 'Manage ERP system connectors, configurations, and integrations' }),
  };
} 