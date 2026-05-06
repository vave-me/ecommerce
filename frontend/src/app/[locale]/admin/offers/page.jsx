import { Suspense } from 'react';
import OffersManagement from './OffersManagement.client';
import LoadingSpinner from '@/components/common/LoadingSpinner';

export const metadata = {
  title: 'Offers Management - Admin Dashboard',
  description: 'Manage and monitor all offer types including BuyNow, Lease, Reservation, and BuyBack',
};

export default function OffersPage() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <OffersManagement />
    </Suspense>
  );
} 