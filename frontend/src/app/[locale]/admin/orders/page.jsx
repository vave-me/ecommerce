import OrdersManagement from './OrdersManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Orders & Payments Management - Admin',
  description: 'Manage orders and payment transactions',
};

export default function OrdersManagementPage() {
  return <OrdersManagement />;
}