import { Metadata } from 'next';
import DatabaseTools from './DatabaseTools.client';

export const metadata = {
  title: 'Database Tools - Admin Dashboard',
  description: 'Database maintenance and optimization tools',
};

export default function DatabasePage() {
  return <DatabaseTools />;
} 