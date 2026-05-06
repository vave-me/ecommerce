import PaymentManagement from './PaymentManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Payment Management - Admin',
  description: 'Manage payment transactions, refunds, and financial operations',
};

export default function PaymentManagementPage() {
  return <PaymentManagement />;
} 