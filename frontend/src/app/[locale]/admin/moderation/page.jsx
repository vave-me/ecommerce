import { Metadata } from 'next';
import ContentModeration from './ContentModeration.client';

export const metadata = {
  title: 'Content Moderation - Admin Dashboard',
  description: 'Review and moderate user-generated content',
};

export default function ModerationPage() {
  return <ContentModeration />;
} 