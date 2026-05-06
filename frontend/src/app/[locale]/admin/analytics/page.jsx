import PlatformAnalytics from './PlatformAnalytics.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Platform Analytics - Admin',
  description: 'Platform performance and analytics dashboard',
};

export default function PlatformAnalyticsPage() {
  return <PlatformAnalytics />;
}