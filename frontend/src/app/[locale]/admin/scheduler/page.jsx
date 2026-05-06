import { Suspense } from 'react';
import SchedulerManagement from './SchedulerManagement.client';
import LoadingSpinner from '@/components/common/LoadingSpinner';

export const metadata = {
  title: 'Scheduler Management - Admin Dashboard',
  description: 'Manage and monitor automated tasks and scheduler activities',
};

export default function SchedulerPage() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <SchedulerManagement />
    </Suspense>
  );
} 