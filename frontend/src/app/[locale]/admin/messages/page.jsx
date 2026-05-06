import MessagesManagement from './MessagesManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Messages Management - Admin',
  description: 'Monitor and manage user conversations and messaging system',
};

export default function MessagesManagementPage() {
  return <MessagesManagement />;
} 