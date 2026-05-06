import React from 'react';
import { notFound } from 'next/navigation';
import { getTranslations } from 'next-intl/server';
import AdminOverview from './AdminOverview.client';

export async function generateMetadata({ params }) {
  const t = await getTranslations('admin');
  
  return {
    title: `${t('overview.title')} | Admin`,
    description: t('overview.description')
  };
}

export default async function AdminPage({ params }) {
  try {
    return <AdminOverview />;
  } catch (error) {
    // Error: 'Admin page error:', error...
    notFound();
  }
} 