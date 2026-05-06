import MediaManagement from './MediaManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Media Management - Admin',
  description: 'Manage uploaded media files, images, and videos',
};

export default function MediaManagementPage() {
  return <MediaManagement />;
}