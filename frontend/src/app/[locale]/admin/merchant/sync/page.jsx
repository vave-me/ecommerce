import MerchantSync from './MerchantSync.client';
import '../../admin-theme.css';

export const metadata = {
  title: 'Merchant Sync - Admin',
  description: 'Manual Google Merchant Center product synchronization',
};

export default function MerchantSyncPage() {
  return <MerchantSync />;
} 