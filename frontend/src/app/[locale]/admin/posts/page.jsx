import PostsManagement from './PostsManagement.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Posts Management - Admin',
  description: 'Manage posts and content',
};

export default function PostsManagementPage() {
  return <PostsManagement />;
}