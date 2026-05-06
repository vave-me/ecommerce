import MerchantProducts from './MerchantProducts.client';
import '../../admin-theme.css';

export const metadata = {
  title: 'Merchant Products - Admin',
  description: 'Manage products in Google Merchant Center',
};

export default function MerchantProductsPage() {
  return <MerchantProducts />;
} 