import { Suspense } from 'react';
import BusinessDashboard from './BusinessDashboard.client';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';

export const metadata = {
  title: 'Business Dashboard',
  description: 'Manage your business and products',
};

export default async function BusinessDashboardPage({ params }) {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <BusinessDashboard />
    </Suspense>
  );
}