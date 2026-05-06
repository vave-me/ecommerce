import AdminSettings from './AdminSettings.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Admin Settings - Platform Configuration',
  description: 'Configure platform settings, features, and system preferences',
};

export default function AdminSettingsPage() {
  return <AdminSettings />;
}