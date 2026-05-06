import ProductDetail from './ProductDetail.client';

export const metadata = {
  title: 'Product Details - Admin',
  description: 'View and manage product information',
};

export default function ProductDetailPage({ params }) {
  return <ProductDetail locale={params.locale} productId={params.productId} />;
}