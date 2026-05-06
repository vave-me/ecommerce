import ProductManagement from './ProductManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Product Management - Admin',
  description: 'Manage products, inventory, pricing, and stock levels',
};

export default function ProductManagementPage() {
  return <ProductManagement />;
}