import { Suspense } from 'react';
import { getTranslations } from 'next-intl/server';
import AddConnectorForm from './AddConnectorForm.client';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';

/**
 * Add New ERP Connector Page
 * Create and configure a new ERP system integration
 */
export default async function AddConnectorPage() {
  const t = await getTranslations('ERPConnectors');

  return (
    <div className="admin-page">
      <div className="page-header">
        <h1>{t('addConnector.title', { defaultValue: 'Add New ERP Connector' })}</h1>
        <p>{t('addConnector.subtitle', { defaultValue: 'Configure a new integration with your ERP system' })}</p>
      </div>

      <Suspense 
        fallback={
          <div className="loading-container">
            <LoadingSpinner />
            <p>{t('loading', { defaultValue: 'Loading form...' })}</p>
          </div>
        }
      >
        <AddConnectorForm />
      </Suspense>
    </div>
  );
}

export async function generateMetadata({ params }) {
  const t = await getTranslations('ERPConnectors');
  
  return {
    title: t('addConnector.metaTitle', { defaultValue: 'Add New ERP Connector' }),
    description: t('addConnector.metaDescription', { defaultValue: 'Create and configure a new ERP system integration' }),
  };
}