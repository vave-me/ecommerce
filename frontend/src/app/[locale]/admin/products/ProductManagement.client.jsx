"use client";

import React, { useState, useCallback, useMemo, lazy, Suspense, useRef, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Search,
  Filter,
  Plus,
  Edit2,
  Trash2,
  Archive,
  Eye,
  EyeOff,
  Package,
  DollarSign,
  TrendingUp,
  TrendingDown,
  AlertCircle,
  Download,
  Upload,
  RefreshCw,
  MoreVertical,
  Tag,
  Box,
  Calendar,
  Image as ImageIcon,
  CheckCircle,
  XCircle,
  Clock,
  ShoppingCart
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  getProducts,
  getProductById,
  addProduct,
  updateProduct,
  deleteProduct,
  archiveProduct,
  updateProductPrice,
  adjustProductStock,
  markProductAsSold,
  markProductAsLeased,
  markProductAsPawned,
  bulkUploadProducts,
  getCategories
} from '@/api/adminApi';
import { batchSyncProducts } from '@/api/client/admin/merchantApi';
import { syncProducts as syncErpProducts, listConnectors } from '@/api/client/admin/erpApi';
import { exportProducts } from '@/api/client/admin/productsApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './ProductManagement.module.css';

// Lazy load the CreateProductModal
const CreateProductModal = lazy(() => import('@/features/CreateProductModal'));

// Product status badge component
const StatusBadge = ({ status }) => {
  const getStatusStyle = () => {
    switch (status) {
      case 'active':
        return styles.statusActive;
      case 'sold':
        return styles.statusSold;
      case 'leased':
        return styles.statusLeased;
      case 'pawned':
        return styles.statusPawned;
      case 'archived':
        return styles.statusArchived;
      case 'out_of_stock':
        return styles.statusOutOfStock;
      default:
        return styles.statusInactive;
    }
  };

  const getStatusIcon = () => {
    switch (status) {
      case 'active':
        return <CheckCircle size={14} />;
      case 'sold':
        return <ShoppingCart size={14} />;
      case 'leased':
        return <Calendar size={14} />;
      case 'pawned':
        return <DollarSign size={14} />;
      case 'archived':
        return <Archive size={14} />;
      case 'out_of_stock':
        return <XCircle size={14} />;
      default:
        return <Clock size={14} />;
    }
  };

  return (
    <span className={`${styles.statusBadge} ${getStatusStyle()}`}>
      {getStatusIcon()}
      {status.replace('_', ' ')}
    </span>
  );
};

// Product row component
const ProductRow = ({ product, onAction, categories, selected, onSelect }) => {
  const [showMenu, setShowMenu] = useState(false);
  const [menuPosition, setMenuPosition] = useState({ top: 0, right: 0 });
  const menuButtonRef = React.useRef(null);
  const category = categories.find(c => c.id === product.categoryId);
  
  const stockStatus = () => {
    if (product.stock === 0) return 'out_of_stock';
    if (product.stock < 10) return 'low_stock';
    return 'in_stock';
  };

  const formatPrice = (price) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: product.currency || 'USD'
    }).format(price);
  };

  return (
    <tr className={`${styles.productRow} ${selected ? styles.selectedRow : ''}`}>
      <td className={styles.checkboxCell}>
        <input
          type="checkbox"
          checked={selected}
          onChange={(e) => onSelect(e.target.checked)}
          className={styles.checkbox}
        />
      </td>
      <td className={styles.productCell}>
        <div className={styles.productInfo}>
          <div className={styles.productImage}>
            {product.thumbnail ? (
              <img src={product.thumbnail} alt={product.name} />
            ) : (
              <ImageIcon size={24} />
            )}
          </div>
          <div className={styles.productDetails}>
            <div className={styles.productName}>{product.name}</div>
            <div className={styles.productMeta}>
              <span className={styles.productSku}>SKU: {product.sku || 'N/A'}</span>
              {product.brand && <span className={styles.productBrand}>{product.brand}</span>}
              {product.condition && <span className={styles.productCondition}>{product.condition}</span>}
            </div>
          </div>
        </div>
      </td>
      <td className={styles.categoryCell}>
        <div className={styles.categoryInfo}>
          <span className={styles.categoryTag}>
            <Tag size={12} />
            {category?.name || 'Uncategorized'}
          </span>
          {product.tags && product.tags.length > 0 && (
            <div className={styles.tagsList}>
              {product.tags.slice(0, 2).map((tag, index) => (
                <span key={index} className={styles.tag}>
                  {tag}
                </span>
              ))}
              {product.tags.length > 2 && (
                <span className={styles.tag}>+{product.tags.length - 2}</span>
              )}
            </div>
          )}
        </div>
      </td>
      <td className={styles.priceCell}>
        <div className={styles.priceInfo}>
          <span className={styles.currentPrice}>
            {formatPrice(product.basePrice || product.price)}
          </span>
          {product.originalPrice && product.originalPrice > (product.basePrice || product.price) && (
            <span className={styles.originalPrice}>{formatPrice(product.originalPrice)}</span>
          )}
          {product.negotiable && (
            <span className={styles.priceType}>Negotiable</span>
          )}
        </div>
      </td>
      <td className={styles.stockCell}>
        <div className={styles.stockInfo}>
          <span className={styles.stockQuantity}>{product.stock || 0}</span>
          <span className={`${styles.stockLevel} ${
            product.stock === 0 ? styles.stockOut : 
            product.stock < 10 ? styles.stockLow : 
            product.stock < 50 ? styles.stockMedium : 
            styles.stockHigh
          }`}>
            {product.stock === 0 ? 'Out' : product.stock < 10 ? 'Low' : product.stock < 50 ? 'Medium' : 'High'}
          </span>
        </div>
      </td>
      <td className={styles.statusCell}>
        <StatusBadge status={product.status || 'active'} />
        {product.merchantStatus === 'synced' && (
          <span className={styles.merchantBadge} title="Synced to Merchant Center">
            <ExternalLink size={12} />
          </span>
        )}
      </td>
      <td className={styles.actionsCell}>
        <div className={styles.actionMenu}>
          <button
            ref={menuButtonRef}
            className={styles.actionButton}
            onClick={(e) => {
              if (!showMenu) {
                const rect = e.currentTarget.getBoundingClientRect();
                setMenuPosition({
                  top: rect.bottom + window.scrollY,
                  right: window.innerWidth - rect.right
                });
              }
              setShowMenu(!showMenu);
            }}
            aria-label="Product actions"
          >
            <MoreVertical size={14} />
          </button>
          {showMenu && typeof document !== 'undefined' && createPortal(
            <>
              <div
                className={styles.menuOverlay}
                onClick={() => setShowMenu(false)}
              />
              <div 
                className={styles.menuDropdown} 
                style={{ 
                  position: 'fixed',
                  top: `${menuPosition.top}px`,
                  right: `${menuPosition.right}px`,
                  zIndex: 1000
                }}
              >
                <button className={styles.menuItem} onClick={() => { onAction('view', product); setShowMenu(false); }}>
                  <Eye size={14} /> View
                </button>
                <button className={styles.menuItem} onClick={() => { onAction('edit', product); setShowMenu(false); }}>
                  <Edit2 size={14} /> Edit
                </button>
                {product.hasVariants && (
                  <button className={styles.menuItem} onClick={() => { onAction('variants', product); setShowMenu(false); }}>
                    <Package size={14} /> Variants
                  </button>
                )}
                <button className={styles.menuItem} onClick={() => { onAction('price', product); setShowMenu(false); }}>
                  <DollarSign size={14} /> Price
                </button>
                <button className={styles.menuItem} onClick={() => { onAction('stock', product); setShowMenu(false); }}>
                  <Box size={14} /> Stock
                </button>
                {product.status !== 'sold' && product.status !== 'leased' && product.status !== 'pawned' && (
                  <>
                    <button className={styles.menuItem} onClick={() => { onAction('sold', product); setShowMenu(false); }}>
                      <ShoppingCart size={14} /> Sold
                    </button>
                    <button className={styles.menuItem} onClick={() => { onAction('lease', product); setShowMenu(false); }}>
                      <Calendar size={14} /> Lease
                    </button>
                    <button className={styles.menuItem} onClick={() => { onAction('pawn', product); setShowMenu(false); }}>
                      <DollarSign size={14} /> Pawn
                    </button>
                  </>
                )}
                <button className={styles.menuItem} onClick={() => { onAction('archive', product); setShowMenu(false); }}>
                  <Archive size={14} /> {product.status === 'archived' ? 'Unarchive' : 'Archive'}
                </button>
                {product.middlemanService && (
                  <button className={styles.menuItem} onClick={() => { onAction('middleman', product); setShowMenu(false); }}>
                    <CheckCircle size={14} /> Middleman
                  </button>
                )}
                <button 
                  className={`${styles.menuItem} ${styles.menuItemDanger}`}
                  onClick={() => { onAction('delete', product); setShowMenu(false); }}
                >
                  <Trash2 size={14} /> Delete
                </button>
              </div>
            </>,
            document.body
          )}
        </div>
      </td>
    </tr>
  );
};

// Quick action modals
const PriceUpdateModal = ({ product, onClose, onSave }) => {
  const [price, setPrice] = useState(product.price);
  const [originalPrice, setOriginalPrice] = useState(product.originalPrice || '');

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave({ 
      price: parseFloat(price),
      originalPrice: originalPrice ? parseFloat(originalPrice) : null 
    });
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h3>Update Price - {product.name}</h3>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGroup}>
            <label className={styles.filterLabel}>Current Price *</label>
            <input
              type="number"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              min="0"
              step="0.01"
              required
            />
          </div>
          <div className={styles.formGroup}>
            <label className={styles.filterLabel}>Original Price (for discount display)</label>
            <input
              type="number"
              value={originalPrice}
              onChange={(e) => setOriginalPrice(e.target.value)}
              min="0"
              step="0.01"
            />
          </div>
          <div className={styles.modalActions}>
            <button type="button" onClick={onClose} className={styles.cancelButton}>
              Cancel
            </button>
            <button type="submit" className={styles.saveButton}>
              Update Price
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const StockAdjustmentModal = ({ product, onClose, onSave }) => {
  const [adjustment, setAdjustment] = useState(0);
  const [operation, setOperation] = useState('add');

  const handleSubmit = (e) => {
    e.preventDefault();
    const finalAdjustment = operation === 'add' ? adjustment : -adjustment;
    onSave(finalAdjustment);
  };

  const newStock = product.stock + (operation === 'add' ? adjustment : -adjustment);

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h3>Adjust Stock - {product.name}</h3>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.stockStatus}>
            <span>Current Stock: <strong>{product.stock}</strong></span>
            <span>New Stock: <strong className={newStock < 0 ? styles.stockDanger : ''}>{newStock}</strong></span>
          </div>
          <div className={styles.formRow}>
            <div className={styles.formGroup}>
              <label className={styles.filterLabel}>Operation</label>
              <select value={operation} onChange={(e) => setOperation(e.target.value)}>
                <option value="add">Add Stock</option>
                <option value="remove">Remove Stock</option>
              </select>
            </div>
            <div className={styles.formGroup}>
              <label className={styles.filterLabel}>Quantity</label>
              <input
                type="number"
                value={adjustment}
                onChange={(e) => setAdjustment(parseInt(e.target.value) || 0)}
                min="0"
                required
              />
            </div>
          </div>
          {newStock < 0 && (
            <div className={styles.errorMessage}>
              <AlertCircle size={16} />
              Warning: This will result in negative stock!
            </div>
          )}
          <div className={styles.modalActions}>
            <button type="button" onClick={onClose} className={styles.cancelButton}>
              Cancel
            </button>
            <button type="submit" className={styles.saveButton} disabled={newStock < 0}>
              Adjust Stock
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// Lease modal
const LeaseModal = ({ product, onClose, onSave }) => {
  const [monthlyPrice, setMonthlyPrice] = useState('');
  const [leaseTermMonths, setLeaseTermMonths] = useState('12');

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave({
      monthlyPrice: parseFloat(monthlyPrice),
      leaseTermMonths: parseInt(leaseTermMonths)
    });
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h3>Lease Product - {product.name}</h3>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGroup}>
            <label className={styles.filterLabel}>Monthly Lease Price *</label>
            <input
              type="number"
              value={monthlyPrice}
              onChange={(e) => setMonthlyPrice(e.target.value)}
              min="0"
              step="0.01"
              required
            />
          </div>
          <div className={styles.formGroup}>
            <label className={styles.filterLabel}>Lease Term (months) *</label>
            <select
              value={leaseTermMonths}
              onChange={(e) => setLeaseTermMonths(e.target.value)}
            >
              <option value="3">3 months</option>
              <option value="6">6 months</option>
              <option value="12">12 months</option>
              <option value="24">24 months</option>
              <option value="36">36 months</option>
            </select>
          </div>
          <div className={styles.modalActions}>
            <button type="button" onClick={onClose} className={styles.cancelButton}>
              Cancel
            </button>
            <button type="submit" className={styles.saveButton}>
              Mark as Leased
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// Pawn modal
const PawnModal = ({ product, onClose, onSave }) => {
  const [lockedPrice, setLockedPrice] = useState('');
  const [redemptionFee, setRedemptionFee] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave({
      lockedPrice: parseFloat(lockedPrice),
      redemptionFee: parseFloat(redemptionFee)
    });
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h3>Pawn Product - {product.name}</h3>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGroup}>
            <label className={styles.filterLabel}>Locked Price *</label>
            <input
              type="number"
              value={lockedPrice}
              onChange={(e) => setLockedPrice(e.target.value)}
              min="0"
              step="0.01"
              required
            />
          </div>
          <div className={styles.formGroup}>
            <label className={styles.filterLabel}>Redemption Fee *</label>
            <input
              type="number"
              value={redemptionFee}
              onChange={(e) => setRedemptionFee(e.target.value)}
              min="0"
              step="0.01"
              required
            />
          </div>
          <div className={styles.modalActions}>
            <button type="button" onClick={onClose} className={styles.cancelButton}>
              Cancel
            </button>
            <button type="submit" className={styles.saveButton}>
              Mark as Pawned
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// Main component
const ProductManagement = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('ProductManagement');
  } catch (e) {
    // Fallback function for missing translations
    t = (key, options) => options?.defaultValue || key;
  }
  
  const router = useRouter();
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const selectAllRef = useRef(null);
  
  // Store auth context for debugging
  React.useEffect(() => {
    if (typeof window !== 'undefined' && process.env.NODE_ENV === 'development') {
      window.__REACT_AUTH_CONTEXT__ = { user };
    }
  }, [user]);

  const [searchTerm, setSearchTerm] = useState('');
  const [filterCategory, setFilterCategory] = useState('all');
  const [filterStatus, setFilterStatus] = useState('all');
  const [filterBrand, setFilterBrand] = useState('all');
  const [filterCondition, setFilterCondition] = useState('all');
  const [filterNegotiable, setFilterNegotiable] = useState('all');
  const [filterMiddleman, setFilterMiddleman] = useState('all');
  const [sortBy, setSortBy] = useState('createdAt');
  const [sortOrder, setSortOrder] = useState('desc');
  const [currentPage, setCurrentPage] = useState(1);
  const [showFilters, setShowFilters] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [modalType, setModalType] = useState(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingProduct, setEditingProduct] = useState(null);
  const [selectedProducts, setSelectedProducts] = useState([]);
  const [showBulkActions, setShowBulkActions] = useState(false);

  const itemsPerPage = 20;

  // Build query params for API
  const queryParams = useMemo(() => {
    const params = {
      page: currentPage,
      pageSize: itemsPerPage,
      sortBy,
      sortOrder
    };

    if (searchTerm) params.search = searchTerm;
    if (filterCategory !== 'all') params.category = filterCategory;
    if (filterStatus !== 'all') params.status = filterStatus;
    if (filterBrand !== 'all') params.brand = filterBrand;
    if (filterCondition !== 'all') params.condition = filterCondition;
    if (filterNegotiable !== 'all') params.negotiable = filterNegotiable === 'yes';
    if (filterMiddleman !== 'all') params.middlemanService = filterMiddleman === 'yes';

    return params;
  }, [currentPage, itemsPerPage, sortBy, sortOrder, searchTerm, filterCategory, filterStatus, filterBrand, filterCondition, filterNegotiable, filterMiddleman]);

  // Fetch products
  const { data: productsData, isLoading, error, refetch } = useQuery({
    queryKey: ['adminProducts', queryParams],
    queryFn: () => getProducts(queryParams),
    staleTime: 60000,
    keepPreviousData: true
  });

  // Fetch categories
  const { data: categoriesData } = useQuery({
    queryKey: ['categories'],
    queryFn: getCategories,
    staleTime: 300000,
  });

  const categories = categoriesData?.categories || [];

  // Fetch ERP connectors
  const { data: connectorsData } = useQuery({
    queryKey: ['erpConnectors'],
    queryFn: () => listConnectors({ status: 'active' }),
    staleTime: 300000,
  });

  const activeConnectors = connectorsData?.connectors || [];

  // Mutations
  const deleteMutation = useMutation({
    mutationFn: deleteProduct,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminProducts']);
    },
  });

  const archiveMutation = useMutation({
    mutationFn: archiveProduct,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminProducts']);
    },
  });

  const priceMutation = useMutation({
    mutationFn: ({ id, price }) => updateProductPrice(id, price),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminProducts']);
      setModalType(null);
      setSelectedProduct(null);
    },
  });

  const stockMutation = useMutation({
    mutationFn: ({ id, adjustment }) => adjustProductStock(id, adjustment),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminProducts']);
      setModalType(null);
      setSelectedProduct(null);
    },
  });

  const soldMutation = useMutation({
    mutationFn: markProductAsSold,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminProducts']);
    },
  });

  const leaseMutation = useMutation({
    mutationFn: ({ id, ...leaseData }) => markProductAsLeased(id, leaseData),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminProducts']);
      setModalType(null);
      setSelectedProduct(null);
    },
  });

  const pawnMutation = useMutation({
    mutationFn: ({ id, ...pawnData }) => markProductAsPawned(id, pawnData),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminProducts']);
      setModalType(null);
      setSelectedProduct(null);
    },
  });

  // Products are already filtered by the API
  const filteredProducts = productsData?.products || [];
  const totalProducts = productsData?.total || 0;
  const totalPages = Math.ceil(totalProducts / itemsPerPage);

  // Extract unique brands from products
  const uniqueBrands = useMemo(() => {
    const brands = new Set();
    filteredProducts.forEach(p => {
      if (p.brand) brands.add(p.brand);
    });
    return Array.from(brands).sort();
  }, [filteredProducts]);

  // Selection state calculations
  const isAllSelected = filteredProducts.length > 0 && selectedProducts.length === filteredProducts.length;
  const isSomeSelected = selectedProducts.length > 0 && selectedProducts.length < filteredProducts.length;

  // Handle indeterminate state for checkbox
  useEffect(() => {
    if (selectAllRef.current) {
      selectAllRef.current.indeterminate = isSomeSelected;
    }
  }, [isSomeSelected]);

  // Handle modal close
  const handleModalClose = useCallback(() => {
    setShowCreateModal(false);
    setEditingProduct(null);
    queryClient.invalidateQueries(['adminProducts']);
  }, [queryClient]);

  // Selection handlers
  const handleSelectAll = useCallback((checked) => {
    if (checked) {
      setSelectedProducts(filteredProducts.map(p => p.id));
    } else {
      setSelectedProducts([]);
    }
  }, [filteredProducts]);

  const handleSelectProduct = useCallback((productId, checked) => {
    if (checked) {
      setSelectedProducts(prev => [...prev, productId]);
    } else {
      setSelectedProducts(prev => prev.filter(id => id !== productId));
    }
  }, []);

  // Bulk operation handlers
  const handleBulkArchive = useCallback(async () => {
    if (!selectedProducts.length) return;
    if (confirm(`Archive ${selectedProducts.length} products?`)) {
      try {
        await Promise.all(selectedProducts.map(id => archiveProduct(id)));
        queryClient.invalidateQueries(['adminProducts']);
        setSelectedProducts([]);
      } catch (error) {
        // Error: 'Bulk archive failed:', error...
        alert('Failed to archive some products');
      }
    }
  }, [selectedProducts, queryClient]);

  const handleBulkDelete = useCallback(async () => {
    if (!selectedProducts.length) return;
    if (confirm(`Delete ${selectedProducts.length} products? This cannot be undone.`)) {
      try {
        await Promise.all(selectedProducts.map(id => deleteProduct(id)));
        queryClient.invalidateQueries(['adminProducts']);
        setSelectedProducts([]);
      } catch (error) {
        // Error: 'Bulk delete failed:', error...
        alert('Failed to delete some products');
      }
    }
  }, [selectedProducts, queryClient]);

  const handleBulkMerchantSync = useCallback(async () => {
    if (!selectedProducts.length) return;
    try {
      await batchSyncProducts({ productIds: selectedProducts });
      queryClient.invalidateQueries(['adminProducts']);
      setSelectedProducts([]);
      alert('Products synced to merchant center successfully');
    } catch (error) {
      // Error: 'Merchant sync failed:', error...
      alert('Failed to sync products to merchant center');
    }
  }, [selectedProducts, queryClient]);

  const handleBulkErpSync = useCallback(async (connectorId) => {
    if (!connectorId) return;
    try {
      await syncErpProducts(connectorId, {
        batchSize: 100,
        filters: queryParams
      });
      queryClient.invalidateQueries(['adminProducts']);
      alert('ERP sync initiated successfully');
    } catch (error) {
      // Error: 'ERP sync failed:', error...
      alert('Failed to sync products to ERP');
    }
  }, [queryParams, queryClient]);

  const handleExport = useCallback(async () => {
    try {
      const result = await exportProducts(queryParams);
      // Create download link
      const blob = new Blob([result], { type: 'text/csv' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `products-export-${dayjs().format('YYYY-MM-DD')}.csv`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (error) {
      // Error: 'Export failed:', error...
      alert('Failed to export products');
    }
  }, [queryParams]);

  // Handle actions
  const handleProductAction = useCallback((action, product) => {
    switch (action) {
      case 'view':
        router.push(`/admin/products/${product.id}`);
        break;
      case 'edit':
        setEditingProduct(product);
        setShowCreateModal(true);
        break;
      case 'variants':
        router.push(`/admin/products/${product.id}/variants`);
        break;
      case 'price':
        setSelectedProduct(product);
        setModalType('price');
        break;
      case 'stock':
        setSelectedProduct(product);
        setModalType('stock');
        break;
      case 'archive':
        if (confirm(`Are you sure you want to ${product.status === 'archived' ? 'unarchive' : 'archive'} "${product.name}"?`)) {
          archiveMutation.mutate(product.id);
        }
        break;
      case 'sold':
        if (confirm(`Mark "${product.name}" as sold? This action cannot be undone.`)) {
          soldMutation.mutate(product.id);
        }
        break;
      case 'lease':
        setSelectedProduct(product);
        setModalType('lease');
        break;
      case 'pawn':
        setSelectedProduct(product);
        setModalType('pawn');
        break;
      case 'middleman':
        router.push(`/admin/middleman/products/${product.id}`);
        break;
      case 'delete':
        if (confirm(`Are you sure you want to delete "${product.name}"? This action cannot be undone.`)) {
          deleteMutation.mutate(product.id);
        }
        break;
    }
  }, [router, archiveMutation, soldMutation, deleteMutation]);

  // Handle bulk upload
  const handleBulkUpload = useCallback(async (file) => {
    try {
      const formData = new FormData();
      formData.append('file', file);
      await bulkUploadProducts(formData);
      queryClient.invalidateQueries(['adminProducts']);
      alert('Products uploaded successfully!');
    } catch (error) {
      // Error: 'Bulk upload failed:', error...
      alert('Failed to upload products. Please check the file format.');
    }
  }, [queryClient]);

  // Calculate stats
  const stats = useMemo(() => {
    const products = productsData?.products || [];
    return {
      total: productsData?.total || 0,
      active: products.filter(p => p.status === 'active').length,
      sold: products.filter(p => p.status === 'sold').length,
      archived: products.filter(p => p.status === 'archived').length,
      outOfStock: products.filter(p => p.stock === 0).length,
      totalValue: products.reduce((sum, p) => sum + (p.price * (p.stock || 0)), 0)
    };
  }, [productsData]);

  if (!isAdmin) {
    
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertCircle size={48} className={styles.errorIcon} />
          <h2 className={styles.errorTitle}>Access Denied</h2>
          <p className={styles.errorDescription}>You need admin privileges to access this page.</p>
          <p className={styles.errorDescription}>Current role: {user?.role || 'Not logged in'}</p>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
      </div>
    );
  }

  if (error) {
    // Error: 'Product loading error:', error...
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertCircle size={48} className={styles.errorIcon} />
          <h2 className={styles.errorTitle}>Failed to load products</h2>
          <p className={styles.errorDescription}>{error.message || 'An error occurred while fetching products'}</p>
          <button 
            className={styles.retryButton} 
            onClick={() => refetch()}
          >
            <RefreshCw size={16} />
            Try Again
          </button>
        </div>
      </div>
    );
  }

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        {/* Header */}
        <div className={styles.header}>
          <div>
            <h1 className={styles.title}>
              {t('title', { defaultValue: 'Product Management' })}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', { defaultValue: 'Manage your product inventory and pricing' })}
            </p>
          </div>
          <div className={styles.headerActions}>
            {selectedProducts.length > 0 && (
              <>
                <span className={styles.selectedCount}>
                  {selectedProducts.length} selected
                </span>
                <button
                  className={styles.bulkButton}
                  onClick={handleBulkArchive}
                >
                  <Archive size={14} />
                  Archive
                </button>
                <button
                  className={styles.bulkButton}
                  onClick={handleBulkMerchantSync}
                >
                  <Upload size={14} />
                  Sync to Merchant
                </button>
                {activeConnectors.length > 0 && (
                  <select
                    className={styles.bulkSelect}
                    onChange={(e) => e.target.value && handleBulkErpSync(e.target.value)}
                    defaultValue=""
                  >
                    <option value="">ERP Sync...</option>
                    {activeConnectors.map(c => (
                      <option key={c.id} value={c.id}>{c.name} ({c.type})</option>
                    ))}
                  </select>
                )}
                <button
                  className={`${styles.bulkButton} ${styles.bulkDanger}`}
                  onClick={handleBulkDelete}
                >
                  <Trash2 size={14} />
                  Delete
                </button>
              </>
            )}
            <button
              className={styles.exportButton}
              onClick={handleExport}
            >
              <Download size={14} />
              Export
            </button>
            <button
              className={styles.addButton}
              onClick={() => setShowCreateModal(true)}
            >
              <Plus size={14} />
              {t('add', { defaultValue: 'Add Product' })}
            </button>
            <button
              className={styles.refreshButton}
              onClick={() => refetch()}
              disabled={isLoading}
            >
              <RefreshCw size={14} className={isLoading ? styles.spinning : ''} />
            </button>
          </div>
        </div>

        {/* Stats */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <Package size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.total}</span>
              <span className={styles.statLabel}>Total Products</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(34, 197, 94, 0.1)', color: '#22c55e' }}>
              <CheckCircle size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.active}</span>
              <span className={styles.statLabel}>Active</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(239, 68, 68, 0.1)', color: '#ef4444' }}>
              <ShoppingCart size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.sold}</span>
              <span className={styles.statLabel}>Sold</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(245, 158, 11, 0.1)', color: '#f59e0b' }}>
              <AlertCircle size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>{stats.outOfStock}</span>
              <span className={styles.statLabel}>Out of Stock</span>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon} style={{ background: 'rgba(59, 130, 246, 0.1)', color: '#3b82f6' }}>
              <DollarSign size={16} />
            </div>
            <div className={styles.statContent}>
              <span className={styles.statValue}>
                ${stats.totalValue.toLocaleString()}
              </span>
              <span className={styles.statLabel}>Total Value</span>
            </div>
          </div>
        </div>

        {/* Controls */}
        <div className={styles.controls}>
          <div className={styles.searchBox}>
            <Search size={16} />
            <input
              type="text"
              placeholder={t('searchPlaceholder', { defaultValue: 'Search products...' })}
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <button
            className={`${styles.filterButton} ${showFilters ? styles.filterActive : ''}`}
            onClick={() => setShowFilters(!showFilters)}
          >
            <Filter size={16} />
            {t('filters', { defaultValue: 'Filters' })}
          </button>
        </div>

        {/* Filter Panel */}
        {showFilters && (
          <div className={styles.filterPanel}>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('filterByCategory', { defaultValue: 'Category' })}</label>
              <select
                value={filterCategory}
                onChange={(e) => setFilterCategory(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Categories</option>
                {categories.map(cat => (
                  <option key={cat.id} value={cat.id}>{cat.name}</option>
                ))}
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('filterByStatus', { defaultValue: 'Status' })}</label>
              <select
                value={filterStatus}
                onChange={(e) => setFilterStatus(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Statuses</option>
                <option value="active">Active</option>
                <option value="sold">Sold</option>
                <option value="leased">Leased</option>
                <option value="pawned">Pawned</option>
                <option value="archived">Archived</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('filterByBrand', { defaultValue: 'Brand' })}</label>
              <select
                value={filterBrand}
                onChange={(e) => setFilterBrand(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Brands</option>
                {Array.from(new Set(productsData?.products?.map(p => p.brand).filter(Boolean))).map(brand => (
                  <option key={brand} value={brand}>{brand}</option>
                ))}
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('filterByCondition', { defaultValue: 'Condition' })}</label>
              <select
                value={filterCondition}
                onChange={(e) => setFilterCondition(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All Conditions</option>
                <option value="new">New</option>
                <option value="like_new">Like New</option>
                <option value="excellent">Excellent</option>
                <option value="good">Good</option>
                <option value="fair">Fair</option>
                <option value="poor">Poor</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('negotiable', { defaultValue: 'Negotiable' })}</label>
              <select
                value={filterNegotiable}
                onChange={(e) => setFilterNegotiable(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All</option>
                <option value="yes">Yes</option>
                <option value="no">No</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('middleman', { defaultValue: 'Middleman' })}</label>
              <select
                value={filterMiddleman}
                onChange={(e) => setFilterMiddleman(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="all">All</option>
                <option value="yes">Available</option>
                <option value="no">Not Available</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('sortBy', { defaultValue: 'Sort By' })}</label>
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="createdAt">Date Added</option>
                <option value="name">Name</option>
                <option value="price">Price</option>
                <option value="stock">Stock</option>
                <option value="brand">Brand</option>
              </select>
            </div>
            <div className={styles.filterGroup}>
              <label className={styles.filterLabel}>{t('sortOrder', { defaultValue: 'Order' })}</label>
              <select
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value)}
                className={styles.filterSelect}
              >
                <option value="desc">Descending</option>
                <option value="asc">Ascending</option>
              </select>
            </div>
          </div>
        )}

        {/* Products Table */}
        <div className={styles.tableContainer}>
          <table className={styles.productsTable}>
            <thead className={styles.tableHeader}>
              <tr>
                <th style={{ width: '40px' }}>
                  <input
                    ref={selectAllRef}
                    type="checkbox"
                    checked={isAllSelected}
                    onChange={(e) => handleSelectAll(e.target.checked)}
                    className={styles.checkbox}
                  />
                </th>
                <th>Product</th>
                <th>Category</th>
                <th>Price</th>
                <th>Stock</th>
                <th>Status</th>
                <th style={{ width: '60px' }}></th>
              </tr>
            </thead>
            <tbody>
              {filteredProducts.length === 0 ? (
                <tr>
                  <td colSpan="7" className={styles.emptyState}>
                    <Package size={48} className={styles.emptyIcon} />
                    <h3 className={styles.emptyTitle}>No products found</h3>
                    <p className={styles.emptyDescription}>Try adjusting your filters or add a new product</p>
                  </td>
                </tr>
              ) : (
                filteredProducts.map(product => (
                  <ProductRow
                    key={product.id}
                    product={product}
                    onAction={handleProductAction}
                    categories={categories}
                    selected={selectedProducts.includes(product.id)}
                    onSelect={(checked) => handleSelectProduct(product.id, checked)}
                  />
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className={styles.pagination}>
            <button
              className={styles.pageButton}
              onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
              disabled={currentPage === 1}
            >
              Previous
            </button>
            <span className={styles.pageInfo}>
              Page {currentPage} of {totalPages} ({totalProducts} products)
            </span>
            <button
              className={styles.pageButton}
              onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
              disabled={currentPage === totalPages}
            >
              Next
            </button>
          </div>
        )}

        {/* Modals */}
        {modalType === 'price' && selectedProduct && (
          <PriceUpdateModal
            product={selectedProduct}
            onClose={() => {
              setModalType(null);
              setSelectedProduct(null);
            }}
            onSave={(priceData) => {
              priceMutation.mutate({ id: selectedProduct.id, ...priceData });
            }}
          />
        )}

        {modalType === 'stock' && selectedProduct && (
          <StockAdjustmentModal
            product={selectedProduct}
            onClose={() => {
              setModalType(null);
              setSelectedProduct(null);
            }}
            onSave={(adjustment) => {
              stockMutation.mutate({ id: selectedProduct.id, adjustment });
            }}
          />
        )}

        {modalType === 'lease' && selectedProduct && (
          <LeaseModal
            product={selectedProduct}
            onClose={() => {
              setModalType(null);
              setSelectedProduct(null);
            }}
            onSave={(leaseData) => {
              leaseMutation.mutate({ id: selectedProduct.id, ...leaseData });
            }}
          />
        )}

        {modalType === 'pawn' && selectedProduct && (
          <PawnModal
            product={selectedProduct}
            onClose={() => {
              setModalType(null);
              setSelectedProduct(null);
            }}
            onSave={(pawnData) => {
              pawnMutation.mutate({ id: selectedProduct.id, ...pawnData });
            }}
          />
        )}

        {/* Create/Edit Product Modal */}
        {showCreateModal && (
          <Suspense fallback={<LoadingSpinner />}>
            <CreateProductModal
              onClose={handleModalClose}
              editMode={!!editingProduct}
              initialProductData={editingProduct}
            />
          </Suspense>
        )}
      </div>
    </ErrorBoundary>
  );
};

export default ProductManagement;