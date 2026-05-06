import ReviewManagement from './ReviewManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Reviews Management - Admin',
  description: 'Manage product reviews and ratings',
};

export default function ReviewManagementPage() {
  return <ReviewManagement />;
}