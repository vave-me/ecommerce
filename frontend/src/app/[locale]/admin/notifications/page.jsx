import NotificationsManagement from './NotificationsManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Notifications Management - Admin',
  description: 'Manage system notifications, alerts, and user communications',
};

export default function NotificationsManagementPage() {
  return <NotificationsManagement />;
} 