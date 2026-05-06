import NewsletterManagement from './NewsletterManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Newsletter Management - Admin',
  description: 'Manage newsletter subscriptions and campaigns',
};

export default function NewsletterManagementPage() {
  return <NewsletterManagement />;
}