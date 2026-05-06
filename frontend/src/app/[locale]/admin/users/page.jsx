import UserManagement from './UserManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'User Management - Admin',
  description: 'Manage platform users and permissions',
};

export default function UserManagementPage({ params: { locale } }) {
  return <UserManagement locale={locale} />;
}