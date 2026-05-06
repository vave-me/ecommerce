import { Suspense } from 'react';
import { getTranslations } from 'next-intl/server';
import ERPDashboard from './ERPDashboard.client';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import '../admin-theme.css';

/**
 * ERP Integration Management - Main Page
 * Enterprise-grade ERP system integrations dashboard
 */
export default async function ERPManagementPage() {
  const t = await getTranslations('ERPManagement');

  return (
    <div className="admin-page">
      <div className="page-header">
        <h1>{t('title', { defaultValue: 'ERP Integration Management' })}</h1>
        <p>{t('subtitle', { defaultValue: 'Manage enterprise resource planning system integrations, connectors, and data synchronization' })}</p>
      </div>

      <Suspense 
        fallback={
          <div className="loading-container">
            <LoadingSpinner />
            <p>{t('loading', { defaultValue: 'Loading ERP dashboard...' })}</p>
          </div>
        }
      >
        <ERPDashboard />
      </Suspense>
    </div>
  );
}

export async function generateMetadata({ params }) {
  const t = await getTranslations('ERPManagement');
  
  return {
    title: t('metaTitle', { defaultValue: 'ERP Integration Management' }),
    description: t('metaDescription', { defaultValue: 'Manage ERP system integrations, connectors, and data synchronization for enterprise resource planning' }),
  };
} 