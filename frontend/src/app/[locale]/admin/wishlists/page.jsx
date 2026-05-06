import WishlistsManagement from './WishlistsManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Wishlists Management - Admin',
  description: 'Manage and analyze user wishlists',
};

export default function WishlistsManagementPage() {
  return <WishlistsManagement />;
}