import ServiceManagement from './ServiceManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Service Management - Admin',
  description: 'Manage services, providers, and availability',
};

export default function ServiceManagementPage() {
  return <ServiceManagement />;
}