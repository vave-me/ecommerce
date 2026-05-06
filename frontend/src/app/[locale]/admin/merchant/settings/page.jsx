import MerchantSettings from './MerchantSettings.client';
import '../../admin-theme.css';

export const metadata = {
  title: 'Merchant Settings - Admin',
  description: 'Configure Google Merchant Center integration settings',
};

export default function MerchantSettingsPage() {
  return <MerchantSettings />;
} 