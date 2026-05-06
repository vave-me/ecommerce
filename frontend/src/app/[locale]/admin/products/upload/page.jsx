import { Metadata } from 'next';
import ProductsUpload from './ProductsUpload.client';

export const metadata = {
  title: 'Bulk Upload - Products Management',
  description: 'Upload and manage products in bulk',
};

export default function ProductsUploadPage() {
  return <ProductsUpload />;
} 