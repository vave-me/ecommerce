import { Suspense } from 'react';
import { getTranslations } from 'next-intl/server';
import ConnectorDetails from './ConnectorDetails.client';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';

/**
 * ERP Connector Details Page
 * View and manage a specific ERP connector
 */
export default async function ConnectorDetailsPage({ params }) {
  const t = await getTranslations('ERPConnectors');

  return (
    <div className="admin-page">
      <div className="page-header">
        <h1>{t('details.title', { defaultValue: 'Connector Details' })}</h1>
        <p>{t('details.subtitle', { defaultValue: 'View and manage ERP connector configuration' })}</p>
      </div>

      <Suspense 
        fallback={
          <div className="loading-container">
            <LoadingSpinner />
            <p>{t('loading', { defaultValue: 'Loading connector details...' })}</p>
          </div>
        }
      >
        <ConnectorDetails connectorId={params.id} />
      </Suspense>
    </div>
  );
}

export async function generateMetadata({ params }) {
  const t = await getTranslations('ERPConnectors');
  
  return {
    title: t('details.metaTitle', { defaultValue: 'ERP Connector Details' }),
    description: t('details.metaDescription', { defaultValue: 'View and manage ERP connector configuration and status' }),
  };
}