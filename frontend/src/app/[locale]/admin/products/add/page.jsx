import AddProduct from './AddProduct.client';
import '../../admin-theme.css';

export const metadata = {
  title: 'Add Product - Admin',
  description: 'Add a new product to the inventory',
};

export default function AddProductPage() {
  return <AddProduct />;
}