"use client";

import React, { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import {
    ArrowLeft,
    Package,
    DollarSign,
    BarChart3,
    Image,
    Calendar,
    MapPin,
    Tag,
    ShoppingCart,
    AlertCircle,
    CheckCircle,
    XCircle,
    Edit,
    Trash2,
    Archive,
    TrendingUp,
    TrendingDown,
    Clock,
    Eye,
    Heart,
    MessageCircle,
    Share2,
    MoreVertical,
    Upload,
    Download,
    RefreshCw,
    ShieldCheck,
    Truck,
    CreditCard
} from 'lucide-react';
import { 
    getAdminProductById, 
    updateProductPrice, 
    adjustProductStock,
    archiveProduct,
    markProductSold,
    markProductLeased,
    markProductPawned,
    deleteProduct,
    getProductVariants,
    updateProductThumbnail
} from '@/api/client/admin/productsApi';
import { getProductMerchantStatus, syncProductToMerchant } from '@/api/client/admin/merchantApi';
import { getReviewsForItem } from '@/api/client/admin/reviewsApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import { toast } from 'react-toastify';
import dayjs from 'dayjs';
import styles from './ProductDetail.module.css';

const ProductDetail = ({ locale, productId }) => {
    const t = useTranslations('ProductDetail');
    const router = useRouter();
    const queryClient = useQueryClient();
    const [activeTab, setActiveTab] = useState('overview');
    const [showPriceModal, setShowPriceModal] = useState(false);
    const [showStockModal, setShowStockModal] = useState(false);
    const [showStatusModal, setShowStatusModal] = useState(false);
    const [selectedImage, setSelectedImage] = useState(null);

    // Fetch product details
    const { data: productData, isLoading, error } = useQuery({
        queryKey: ['adminProduct', productId],
        queryFn: () => getAdminProductById(productId),
        staleTime: 60000,
    });

    const product = productData?.product || productData;

    // Fetch product variants
    const { data: variantsData } = useQuery({
        queryKey: ['productVariants', productId],
        queryFn: () => getProductVariants(productId),
        enabled: !!product?.hasVariants,
    });

    // Fetch merchant status
    const { data: merchantStatus } = useQuery({
        queryKey: ['productMerchantStatus', productId],
        queryFn: () => getProductMerchantStatus(productId),
        enabled: !!product,
        retry: false,
    });

    // Fetch reviews
    const { data: reviewsData } = useQuery({
        queryKey: ['productReviews', productId],
        queryFn: () => getReviewsForItem(productId),
        enabled: !!product,
    });

    // Mutations
    const updatePriceMutation = useMutation({
        mutationFn: (priceData) => updateProductPrice(productId, priceData),
        onSuccess: () => {
            queryClient.invalidateQueries(['adminProduct', productId]);
            toast.success('Price updated successfully');
            setShowPriceModal(false);
        },
        onError: (error) => {
            toast.error(`Failed to update price: ${error.message}`);
        }
    });

    const updateStockMutation = useMutation({
        mutationFn: (stockData) => adjustProductStock(productId, stockData),
        onSuccess: () => {
            queryClient.invalidateQueries(['adminProduct', productId]);
            toast.success('Stock updated successfully');
            setShowStockModal(false);
        },
        onError: (error) => {
            toast.error(`Failed to update stock: ${error.message}`);
        }
    });

    const updateStatusMutation = useMutation({
        mutationFn: async ({ status, data }) => {
            switch (status) {
                case 'archived':
                    return archiveProduct(productId);
                case 'sold':
                    return markProductSold(productId);
                case 'leased':
                    return markProductLeased(productId, data);
                case 'pawned':
                    return markProductPawned(productId, data);
                default:
                    throw new Error('Invalid status');
            }
        },
        onSuccess: () => {
            queryClient.invalidateQueries(['adminProduct', productId]);
            queryClient.invalidateQueries(['adminProducts']);
            toast.success('Product status updated');
            setShowStatusModal(false);
        },
        onError: (error) => {
            toast.error(`Failed to update status: ${error.message}`);
        }
    });

    const deleteMutation = useMutation({
        mutationFn: () => deleteProduct(productId),
        onSuccess: () => {
            toast.success('Product deleted successfully');
            router.push(`/${locale}/admin/products`);
        },
        onError: (error) => {
            toast.error(`Failed to delete product: ${error.message}`);
        }
    });

    const syncMerchantMutation = useMutation({
        mutationFn: () => syncProductToMerchant({
            productId,
            product: {
                id: product.id,
                name: product.name,
                description: product.description,
                price: product.basePrice,
                currency: 'USD',
                availability: product.stock > 0 ? 'in stock' : 'out of stock',
                condition: product.condition || 'new',
                brand: product.brand,
                googleProductCategory: product.categoryId,
                imageUrl: product.thumbnail,
                link: `${window.location.origin}/products/${productId}`,
                stock: product.stock,
                sku: product.sku
            }
        }),
        onSuccess: () => {
            queryClient.invalidateQueries(['productMerchantStatus', productId]);
            toast.success('Product synced to merchant center');
        },
        onError: (error) => {
            toast.error(`Failed to sync: ${error.message}`);
        }
    });

    const handleDeleteProduct = () => {
        if (confirm('Are you sure you want to delete this product? This action cannot be undone.')) {
            deleteMutation.mutate();
        }
    };

    const formatCurrency = (amount) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
        }).format(amount / 100);
    };

    const getStatusBadgeClass = (status) => {
        const statusMap = {
            active: styles.statusActive,
            sold: styles.statusSold,
            leased: styles.statusLeased,
            pawned: styles.statusPawned,
            archived: styles.statusArchived,
        };
        return statusMap[status?.toLowerCase()] || styles.statusDefault;
    };

    if (isLoading) {
        return (
            <div className={styles.loadingContainer}>
                <LoadingSpinner />
            </div>
        );
    }

    if (error || !product) {
        return (
            <div className={styles.container}>
                <div className={styles.errorContainer}>
                    <AlertCircle size={48} className={styles.errorIcon} />
                    <h2>Product Not Found</h2>
                    <p>The product you're looking for doesn't exist or has been removed.</p>
                    <button 
                        className={styles.primaryButton}
                        onClick={() => router.push(`/${locale}/admin/products`)}
                    >
                        Back to Products
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
                    <div className={styles.headerLeft}>
                        <button 
                            className={styles.backButton}
                            onClick={() => router.push(`/${locale}/admin/products`)}
                        >
                            <ArrowLeft size={20} />
                            Back to Products
                        </button>
                        <div className={styles.headerInfo}>
                            <h1 className={styles.title}>{product.name}</h1>
                            <div className={styles.headerMeta}>
                                <span className={styles.sku}>SKU: {product.sku || 'N/A'}</span>
                                <span className={`${styles.statusBadge} ${getStatusBadgeClass(product.status)}`}>
                                    {product.status || 'Active'}
                                </span>
                            </div>
                        </div>
                    </div>
                    <div className={styles.headerActions}>
                        <button
                            className={styles.secondaryButton}
                            onClick={() => router.push(`/${locale}/admin/products/${productId}/edit`)}
                        >
                            <Edit size={16} />
                            Edit
                        </button>
                        <div className={styles.dropdown}>
                            <button className={styles.dropdownTrigger}>
                                <MoreVertical size={16} />
                            </button>
                            <div className={styles.dropdownMenu}>
                                <button onClick={() => setShowStatusModal(true)}>
                                    <Archive size={16} />
                                    Change Status
                                </button>
                                <button onClick={() => syncMerchantMutation.mutate()}>
                                    <RefreshCw size={16} />
                                    Sync to Merchant
                                </button>
                                <button onClick={() => router.push(`/${locale}/admin/products/${productId}/variants`)}>
                                    <Package size={16} />
                                    Manage Variants
                                </button>
                                <button onClick={() => router.push(`/${locale}/admin/products/${productId}/analytics`)}>
                                    <BarChart3 size={16} />
                                    View Analytics
                                </button>
                                <hr />
                                <button 
                                    onClick={handleDeleteProduct}
                                    className={styles.dangerAction}
                                >
                                    <Trash2 size={16} />
                                    Delete Product
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Tabs */}
                <div className={styles.tabs}>
                    <button
                        className={`${styles.tab} ${activeTab === 'overview' ? styles.tabActive : ''}`}
                        onClick={() => setActiveTab('overview')}
                    >
                        Overview
                    </button>
                    <button
                        className={`${styles.tab} ${activeTab === 'details' ? styles.tabActive : ''}`}
                        onClick={() => setActiveTab('details')}
                    >
                        Details
                    </button>
                    <button
                        className={`${styles.tab} ${activeTab === 'inventory' ? styles.tabActive : ''}`}
                        onClick={() => setActiveTab('inventory')}
                    >
                        Inventory
                    </button>
                    <button
                        className={`${styles.tab} ${activeTab === 'pricing' ? styles.tabActive : ''}`}
                        onClick={() => setActiveTab('pricing')}
                    >
                        Pricing
                    </button>
                    {product.hasVariants && (
                        <button
                            className={`${styles.tab} ${activeTab === 'variants' ? styles.tabActive : ''}`}
                            onClick={() => setActiveTab('variants')}
                        >
                            Variants ({variantsData?.totalCount || 0})
                        </button>
                    )}
                    <button
                        className={`${styles.tab} ${activeTab === 'reviews' ? styles.tabActive : ''}`}
                        onClick={() => setActiveTab('reviews')}
                    >
                        Reviews ({reviewsData?.reviews?.length || 0})
                    </button>
                </div>

                {/* Tab Content */}
                <div className={styles.tabContent}>
                    {activeTab === 'overview' && (
                        <div className={styles.overviewGrid}>
                            {/* Product Image */}
                            <div className={styles.imageSection}>
                                <div className={styles.mainImage}>
                                    {product.thumbnail ? (
                                        <img 
                                            src={product.thumbnail} 
                                            alt={product.name}
                                            onClick={() => setSelectedImage(product.thumbnail)}
                                        />
                                    ) : (
                                        <div className={styles.noImage}>
                                            <Image size={48} />
                                            <span>No image available</span>
                                        </div>
                                    )}
                                </div>
                                <button className={styles.uploadButton}>
                                    <Upload size={16} />
                                    Update Image
                                </button>
                            </div>

                            {/* Key Metrics */}
                            <div className={styles.metricsGrid}>
                                <div className={styles.metricCard}>
                                    <div className={styles.metricIcon}>
                                        <DollarSign size={24} />
                                    </div>
                                    <div className={styles.metricContent}>
                                        <h3>{formatCurrency(product.basePrice)}</h3>
                                        <p>Base Price</p>
                                    </div>
                                    <button 
                                        className={styles.metricAction}
                                        onClick={() => setShowPriceModal(true)}
                                    >
                                        <Edit size={16} />
                                    </button>
                                </div>

                                <div className={styles.metricCard}>
                                    <div className={styles.metricIcon}>
                                        <Package size={24} />
                                    </div>
                                    <div className={styles.metricContent}>
                                        <h3>{product.stock || 0}</h3>
                                        <p>In Stock</p>
                                    </div>
                                    <button 
                                        className={styles.metricAction}
                                        onClick={() => setShowStockModal(true)}
                                    >
                                        <Edit size={16} />
                                    </button>
                                </div>

                                <div className={styles.metricCard}>
                                    <div className={styles.metricIcon}>
                                        <Eye size={24} />
                                    </div>
                                    <div className={styles.metricContent}>
                                        <h3>0</h3>
                                        <p>Views</p>
                                    </div>
                                </div>

                                <div className={styles.metricCard}>
                                    <div className={styles.metricIcon}>
                                        <ShoppingCart size={24} />
                                    </div>
                                    <div className={styles.metricContent}>
                                        <h3>0</h3>
                                        <p>Orders</p>
                                    </div>
                                </div>
                            </div>

                            {/* Product Info */}
                            <div className={styles.infoSection}>
                                <h2>Product Information</h2>
                                <div className={styles.infoGrid}>
                                    <div className={styles.infoItem}>
                                        <span className={styles.infoLabel}>Category</span>
                                        <span className={styles.infoValue}>{product.categoryId || 'Uncategorized'}</span>
                                    </div>
                                    <div className={styles.infoItem}>
                                        <span className={styles.infoLabel}>Brand</span>
                                        <span className={styles.infoValue}>{product.brand || 'No brand'}</span>
                                    </div>
                                    <div className={styles.infoItem}>
                                        <span className={styles.infoLabel}>Condition</span>
                                        <span className={styles.infoValue}>{product.condition || 'New'}</span>
                                    </div>
                                    <div className={styles.infoItem}>
                                        <span className={styles.infoLabel}>Model</span>
                                        <span className={styles.infoValue}>{product.model || 'N/A'}</span>
                                    </div>
                                    <div className={styles.infoItem}>
                                        <span className={styles.infoLabel}>Seller</span>
                                        <span className={styles.infoValue}>{product.userSellerId || 'Unknown'}</span>
                                    </div>
                                    <div className={styles.infoItem}>
                                        <span className={styles.infoLabel}>Location</span>
                                        <span className={styles.infoValue}>
                                            {product.lat && product.lng ? (
                                                <span className={styles.location}>
                                                    <MapPin size={14} />
                                                    {product.lat.toFixed(2)}, {product.lng.toFixed(2)}
                                                </span>
                                            ) : (
                                                'Not specified'
                                            )}
                                        </span>
                                    </div>
                                </div>
                            </div>

                            {/* Description */}
                            <div className={styles.descriptionSection}>
                                <h2>Description</h2>
                                <p>{product.description || 'No description available.'}</p>
                            </div>

                            {/* Tags */}
                            {product.tags && product.tags.length > 0 && (
                                <div className={styles.tagsSection}>
                                    <h2>Tags</h2>
                                    <div className={styles.tags}>
                                        {product.tags.map((tag, index) => (
                                            <span key={index} className={styles.tag}>
                                                <Tag size={14} />
                                                {tag}
                                            </span>
                                        ))}
                                    </div>
                                </div>
                            )}

                            {/* Merchant Status */}
                            <div className={styles.merchantSection}>
                                <h2>Merchant Center</h2>
                                <div className={styles.merchantStatus}>
                                    {merchantStatus?.status ? (
                                        <>
                                            <div className={styles.merchantBadge}>
                                                <CheckCircle size={16} />
                                                Synced
                                            </div>
                                            <span className={styles.merchantDate}>
                                                Last synced: {dayjs(merchantStatus.lastSyncedAt).format('MMM D, YYYY HH:mm')}
                                            </span>
                                        </>
                                    ) : (
                                        <>
                                            <div className={styles.merchantBadge}>
                                                <XCircle size={16} />
                                                Not Synced
                                            </div>
                                            <button 
                                                className={styles.syncButton}
                                                onClick={() => syncMerchantMutation.mutate()}
                                                disabled={syncMerchantMutation.isLoading}
                                            >
                                                {syncMerchantMutation.isLoading ? 'Syncing...' : 'Sync Now'}
                                            </button>
                                        </>
                                    )}
                                </div>
                            </div>
                        </div>
                    )}

                    {activeTab === 'details' && (
                        <div className={styles.detailsContent}>
                            <div className={styles.detailsGrid}>
                                <div className={styles.detailCard}>
                                    <h3>Product Specifications</h3>
                                    <div className={styles.detailsList}>
                                        <div className={styles.detailItem}>
                                            <span>Weight</span>
                                            <span>{product.weight ? `${product.weight}g` : 'Not specified'}</span>
                                        </div>
                                        <div className={styles.detailItem}>
                                            <span>Dimensions</span>
                                            <span>
                                                {product.width && product.height && product.depth
                                                    ? `${product.width} × ${product.height} × ${product.depth} mm`
                                                    : 'Not specified'}
                                            </span>
                                        </div>
                                        <div className={styles.detailItem}>
                                            <span>Shipping Cost</span>
                                            <span>{product.shippingCost ? formatCurrency(product.shippingCost) : 'Free'}</span>
                                        </div>
                                        <div className={styles.detailItem}>
                                            <span>Negotiable</span>
                                            <span>{product.negotiable ? 'Yes' : 'No'}</span>
                                        </div>
                                        <div className={styles.detailItem}>
                                            <span>Middleman Service</span>
                                            <span>{product.middlemanService ? 'Available' : 'Not available'}</span>
                                        </div>
                                        <div className={styles.detailItem}>
                                            <span>User Type</span>
                                            <span>{product.userType || 'Individual'}</span>
                                        </div>
                                    </div>
                                </div>

                                {product.attributes && product.attributes.length > 0 && (
                                    <div className={styles.detailCard}>
                                        <h3>Custom Attributes</h3>
                                        <div className={styles.detailsList}>
                                            {product.attributes.map((attr, index) => (
                                                <div key={index} className={styles.detailItem}>
                                                    <span>{attr.key}</span>
                                                    <span>{attr.value}</span>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}

                                {product.options && product.options.length > 0 && (
                                    <div className={styles.detailCard}>
                                        <h3>Product Options</h3>
                                        <div className={styles.optionsList}>
                                            {product.options.map((option, index) => (
                                                <div key={index} className={styles.optionItem}>
                                                    <span className={styles.optionName}>{option.name}</span>
                                                    <span className={styles.optionValue}>{option.value}</span>
                                                    {option.price && (
                                                        <span className={styles.optionPrice}>
                                                            +{formatCurrency(option.price)}
                                                        </span>
                                                    )}
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}
                            </div>
                        </div>
                    )}

                    {activeTab === 'inventory' && (
                        <div className={styles.inventoryContent}>
                            <div className={styles.inventoryHeader}>
                                <h2>Stock Management</h2>
                                <button 
                                    className={styles.primaryButton}
                                    onClick={() => setShowStockModal(true)}
                                >
                                    <Edit size={16} />
                                    Adjust Stock
                                </button>
                            </div>

                            <div className={styles.stockInfo}>
                                <div className={styles.stockCard}>
                                    <h3>Current Stock</h3>
                                    <div className={styles.stockValue}>
                                        <Package size={32} />
                                        <span className={styles.stockNumber}>{product.stock || 0}</span>
                                        <span className={styles.stockUnit}>units</span>
                                    </div>
                                    {product.manageStocks && (
                                        <div className={styles.stockStatus}>
                                            {product.stock > 10 ? (
                                                <span className={styles.inStock}>
                                                    <CheckCircle size={16} />
                                                    In Stock
                                                </span>
                                            ) : product.stock > 0 ? (
                                                <span className={styles.lowStock}>
                                                    <AlertCircle size={16} />
                                                    Low Stock
                                                </span>
                                            ) : (
                                                <span className={styles.outOfStock}>
                                                    <XCircle size={16} />
                                                    Out of Stock
                                                </span>
                                            )}
                                        </div>
                                    )}
                                </div>

                                <div className={styles.stockSettings}>
                                    <h3>Stock Settings</h3>
                                    <div className={styles.settingsList}>
                                        <div className={styles.settingItem}>
                                            <span>Stock Management</span>
                                            <span>{product.manageStocks ? 'Enabled' : 'Disabled'}</span>
                                        </div>
                                        <div className={styles.settingItem}>
                                            <span>Low Stock Threshold</span>
                                            <span>10 units</span>
                                        </div>
                                        <div className={styles.settingItem}>
                                            <span>Allow Backorders</span>
                                            <span>No</span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}

                    {activeTab === 'pricing' && (
                        <div className={styles.pricingContent}>
                            <div className={styles.pricingHeader}>
                                <h2>Pricing Information</h2>
                                <button 
                                    className={styles.primaryButton}
                                    onClick={() => setShowPriceModal(true)}
                                >
                                    <Edit size={16} />
                                    Update Price
                                </button>
                            </div>

                            <div className={styles.pricingGrid}>
                                <div className={styles.priceCard}>
                                    <h3>Base Price</h3>
                                    <div className={styles.priceValue}>
                                        <DollarSign size={32} />
                                        <span className={styles.priceAmount}>{formatCurrency(product.basePrice)}</span>
                                    </div>
                                    <div className={styles.priceDetails}>
                                        <div className={styles.priceDetail}>
                                            <span>Currency</span>
                                            <span>USD</span>
                                        </div>
                                        <div className={styles.priceDetail}>
                                            <span>Tax Included</span>
                                            <span>No</span>
                                        </div>
                                    </div>
                                </div>

                                <div className={styles.priceHistory}>
                                    <h3>Price History</h3>
                                    <div className={styles.historyList}>
                                        <div className={styles.historyItem}>
                                            <span className={styles.historyDate}>Current</span>
                                            <span className={styles.historyPrice}>{formatCurrency(product.basePrice)}</span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}

                    {activeTab === 'variants' && product.hasVariants && (
                        <div className={styles.variantsContent}>
                            <div className={styles.variantsHeader}>
                                <h2>Product Variants</h2>
                                <button 
                                    className={styles.primaryButton}
                                    onClick={() => router.push(`/${locale}/admin/products/${productId}/variants`)}
                                >
                                    <Package size={16} />
                                    Manage Variants
                                </button>
                            </div>

                            {variantsData?.variants?.length > 0 ? (
                                <div className={styles.variantsGrid}>
                                    {variantsData.variants.map((variant) => (
                                        <div key={variant.id} className={styles.variantCard}>
                                            <div className={styles.variantHeader}>
                                                <h4>{variant.sku}</h4>
                                                <span className={`${styles.variantStatus} ${variant.isAvailable ? styles.available : styles.unavailable}`}>
                                                    {variant.isAvailable ? 'Available' : 'Unavailable'}
                                                </span>
                                            </div>
                                            <div className={styles.variantDetails}>
                                                <div className={styles.variantDetail}>
                                                    <span>Price</span>
                                                    <span>{formatCurrency(variant.variantPrice)}</span>
                                                </div>
                                                <div className={styles.variantDetail}>
                                                    <span>Stock</span>
                                                    <span>{variant.stock} units</span>
                                                </div>
                                                {variant.attributes?.map((attr, index) => (
                                                    <div key={index} className={styles.variantDetail}>
                                                        <span>{attr.key}</span>
                                                        <span>{attr.value}</span>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <div className={styles.emptyState}>
                                    <Package size={48} />
                                    <p>No variants found</p>
                                </div>
                            )}
                        </div>
                    )}

                    {activeTab === 'reviews' && (
                        <div className={styles.reviewsContent}>
                            <div className={styles.reviewsHeader}>
                                <h2>Customer Reviews</h2>
                            </div>

                            {reviewsData?.reviews?.length > 0 ? (
                                <div className={styles.reviewsList}>
                                    {reviewsData.reviews.map((review) => (
                                        <div key={review.id} className={styles.reviewCard}>
                                            <div className={styles.reviewHeader}>
                                                <span className={styles.reviewAuthor}>{review.senderName || 'Anonymous'}</span>
                                                <span className={styles.reviewDate}>
                                                    {dayjs(review.createdAt).format('MMM D, YYYY')}
                                                </span>
                                            </div>
                                            <p className={styles.reviewContent}>{review.content}</p>
                                            <div className={styles.reviewStatus}>
                                                <span className={`${styles.reviewBadge} ${styles[review.reviewStatus]}`}>
                                                    {review.reviewStatus}
                                                </span>
                                                {review.flagged && (
                                                    <span className={styles.flagged}>
                                                        <AlertCircle size={14} />
                                                        Flagged
                                                    </span>
                                                )}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <div className={styles.emptyState}>
                                    <MessageCircle size={48} />
                                    <p>No reviews yet</p>
                                </div>
                            )}
                        </div>
                    )}
                </div>

                {/* Price Update Modal */}
                {showPriceModal && (
                    <div className={styles.modal}>
                        <div className={styles.modalContent}>
                            <h2>Update Product Price</h2>
                            <form onSubmit={(e) => {
                                e.preventDefault();
                                const formData = new FormData(e.target);
                                updatePriceMutation.mutate({
                                    newPrice: parseInt(formData.get('newPrice')) * 100,
                                    oldPrice: product.basePrice
                                });
                            }}>
                                <div className={styles.formGroup}>
                                    <label>Current Price</label>
                                    <input 
                                        type="text" 
                                        value={formatCurrency(product.basePrice)}
                                        disabled
                                    />
                                </div>
                                <div className={styles.formGroup}>
                                    <label>New Price ($)</label>
                                    <input 
                                        type="number" 
                                        name="newPrice"
                                        min="0"
                                        step="1"
                                        defaultValue={product.basePrice / 100}
                                        required
                                    />
                                </div>
                                <div className={styles.modalActions}>
                                    <button 
                                        type="button" 
                                        className={styles.cancelButton}
                                        onClick={() => setShowPriceModal(false)}
                                    >
                                        Cancel
                                    </button>
                                    <button 
                                        type="submit" 
                                        className={styles.primaryButton}
                                        disabled={updatePriceMutation.isLoading}
                                    >
                                        {updatePriceMutation.isLoading ? 'Updating...' : 'Update Price'}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}

                {/* Stock Update Modal */}
                {showStockModal && (
                    <div className={styles.modal}>
                        <div className={styles.modalContent}>
                            <h2>Adjust Stock Level</h2>
                            <form onSubmit={(e) => {
                                e.preventDefault();
                                const formData = new FormData(e.target);
                                updateStockMutation.mutate({
                                    newStock: parseInt(formData.get('newStock'))
                                });
                            }}>
                                <div className={styles.formGroup}>
                                    <label>Current Stock</label>
                                    <input 
                                        type="text" 
                                        value={`${product.stock || 0} units`}
                                        disabled
                                    />
                                </div>
                                <div className={styles.formGroup}>
                                    <label>New Stock Level</label>
                                    <input 
                                        type="number" 
                                        name="newStock"
                                        min="0"
                                        defaultValue={product.stock || 0}
                                        required
                                    />
                                </div>
                                <div className={styles.modalActions}>
                                    <button 
                                        type="button" 
                                        className={styles.cancelButton}
                                        onClick={() => setShowStockModal(false)}
                                    >
                                        Cancel
                                    </button>
                                    <button 
                                        type="submit" 
                                        className={styles.primaryButton}
                                        disabled={updateStockMutation.isLoading}
                                    >
                                        {updateStockMutation.isLoading ? 'Updating...' : 'Update Stock'}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}

                {/* Status Update Modal */}
                {showStatusModal && (
                    <div className={styles.modal}>
                        <div className={styles.modalContent}>
                            <h2>Change Product Status</h2>
                            <form onSubmit={(e) => {
                                e.preventDefault();
                                const formData = new FormData(e.target);
                                const status = formData.get('status');
                                const data = {};
                                
                                if (status === 'leased') {
                                    data.monthlyPrice = parseInt(formData.get('monthlyPrice')) * 100;
                                    data.leaseTermMonths = parseInt(formData.get('leaseTermMonths'));
                                } else if (status === 'pawned') {
                                    data.lockedPrice = parseInt(formData.get('lockedPrice')) * 100;
                                    data.redemptionFee = parseInt(formData.get('redemptionFee')) * 100;
                                }
                                
                                updateStatusMutation.mutate({ status, data });
                            }}>
                                <div className={styles.formGroup}>
                                    <label>Status</label>
                                    <select name="status" defaultValue={product.status?.toLowerCase()} required>
                                        <option value="active">Active</option>
                                        <option value="archived">Archived</option>
                                        <option value="sold">Sold</option>
                                        <option value="leased">Leased</option>
                                        <option value="pawned">Pawned</option>
                                    </select>
                                </div>
                                
                                <div className={styles.conditionalFields}>
                                    <div data-status="leased">
                                        <div className={styles.formGroup}>
                                            <label>Monthly Lease Price ($)</label>
                                            <input type="number" name="monthlyPrice" min="0" step="1" />
                                        </div>
                                        <div className={styles.formGroup}>
                                            <label>Lease Term (months)</label>
                                            <input type="number" name="leaseTermMonths" min="1" />
                                        </div>
                                    </div>
                                    
                                    <div data-status="pawned">
                                        <div className={styles.formGroup}>
                                            <label>Locked Price ($)</label>
                                            <input type="number" name="lockedPrice" min="0" step="1" />
                                        </div>
                                        <div className={styles.formGroup}>
                                            <label>Redemption Fee ($)</label>
                                            <input type="number" name="redemptionFee" min="0" step="1" />
                                        </div>
                                    </div>
                                </div>
                                
                                <div className={styles.modalActions}>
                                    <button 
                                        type="button" 
                                        className={styles.cancelButton}
                                        onClick={() => setShowStatusModal(false)}
                                    >
                                        Cancel
                                    </button>
                                    <button 
                                        type="submit" 
                                        className={styles.primaryButton}
                                        disabled={updateStatusMutation.isLoading}
                                    >
                                        {updateStatusMutation.isLoading ? 'Updating...' : 'Update Status'}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}
            </div>
        </ErrorBoundary>
    );
};

export default ProductDetail;