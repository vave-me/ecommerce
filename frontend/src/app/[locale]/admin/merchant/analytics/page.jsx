import MerchantReports from '../reports/MerchantReports.client';
import '../../admin-theme.css';

export const metadata = {
  title: 'Merchant Analytics - Admin',
  description: 'Google Merchant Center analytics and performance insights',
};

export default function MerchantAnalyticsPage() {
  return <MerchantReports />;
} 