"use client";

import React, { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation } from '@tanstack/react-query';
import {
  ArrowLeft,
  Save,
  Upload,
  X,
  Image as ImageIcon,
  Tag,
  DollarSign,
  Package,
  FileText,
  Link,
  AlertCircle,
  Plus,
  Loader2
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import { addProduct, getCategories } from '@/api/adminApi';

// TODO: Implement uploadProductImages in adminApi
const uploadProductImages = async (images) => {
  // Stub implementation
  return { success: true, urls: images.map(img => URL.createObjectURL(img)) };
};
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './AddProduct.module.css';

// Image upload component
const ImageUploader = ({ images, onImagesChange, onRemove }) => {
  const [uploading, setUploading] = useState(false);
  const [dragActive, setDragActive] = useState(false);

  const handleDrag = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  }, []);

  const handleDrop = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFiles(e.dataTransfer.files);
    }
  }, []);

  const handleChange = useCallback((e) => {
    e.preventDefault();
    if (e.target.files && e.target.files[0]) {
      handleFiles(e.target.files);
    }
  }, []);

  const handleFiles = async (files) => {
    setUploading(true);
    const newImages = [];

    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      if (file.type.startsWith('image/')) {
        const reader = new FileReader();
        reader.onload = (e) => {
          newImages.push({
            id: Date.now() + i,
            url: e.target.result,
            file: file,
            name: file.name
          });
          if (newImages.length === files.length) {
            onImagesChange([...images, ...newImages]);
            setUploading(false);
          }
        };
        reader.readAsDataURL(file);
      }
    }
  };

  return (
    <div className={styles.imageUploader}>
      <label className={styles.uploadLabel}>Images</label>
      <div className={styles.imageGrid}>
        {images.map((image, index) => (
          <div key={image.id} className={styles.imageItem}>
            <img src={image.url} alt={`Product ${index + 1}`} />
            <button
              type="button"
              className={styles.removeImage}
              onClick={() => onRemove(image.id)}
            >
              <X size={16} />
            </button>
            {index === 0 && <span className={styles.primaryBadge}>Primary</span>}
          </div>
        ))}
        <label
          className={`${styles.uploadBox} ${dragActive ? styles.dragActive : ''}`}
          onDragEnter={handleDrag}
          onDragLeave={handleDrag}
          onDragOver={handleDrag}
          onDrop={handleDrop}
        >
          <input
            type="file"
            multiple
            accept="image/*"
            onChange={handleChange}
            className={styles.fileInput}
          />
          {uploading ? (
            <Loader2 size={24} className={styles.spinning} />
          ) : (
            <>
              <Upload size={24} />
              <span>Upload Images</span>
              <span className={styles.uploadHint}>or drag and drop</span>
            </>
          )}
        </label>
      </div>
      <p className={styles.imageHelp}>
        First image will be used as the primary product image. Max 10 images.
      </p>
    </div>
  );
};

// Tag input component
const TagInput = ({ tags, onTagsChange }) => {
  const [inputValue, setInputValue] = useState('');

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && inputValue.trim()) {
      e.preventDefault();
      if (!tags.includes(inputValue.trim())) {
        onTagsChange([...tags, inputValue.trim()]);
      }
      setInputValue('');
    }
  };

  const removeTag = (tagToRemove) => {
    onTagsChange(tags.filter(tag => tag !== tagToRemove));
  };

  return (
    <div className={styles.tagInput}>
      <div className={styles.tagList}>
        {tags.map((tag, index) => (
          <span key={index} className={styles.tag}>
            {tag}
            <button
              type="button"
              onClick={() => removeTag(tag)}
              className={styles.removeTag}
            >
              <X size={14} />
            </button>
          </span>
        ))}
      </div>
      <input
        type="text"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Add tags and press Enter"
        className={styles.tagInputField}
      />
    </div>
  );
};

// Main component
const AddProduct = () => {
  const t = useTranslations('AddProduct');
  const router = useRouter();
  const { isAdmin } = useUserRole();

  const [formData, setFormData] = useState({
    name: '',
    description: '',
    price: '',
    originalPrice: '',
    stock: '',
    sku: '',
    categoryId: '',
    tags: [],
    images: [],
    status: 'active',
    featured: false,
    brand: '',
    model: '',
    condition: 'new',
    warranty: '',
    specifications: '',
    shippingInfo: '',
    returnPolicy: '',
    currency: 'USD'
  });

  const [errors, setErrors] = useState({});
  const [submitting, setSubmitting] = useState(false);

  // Fetch categories
  const { data: categoriesData, isLoading: categoriesLoading } = useQuery({
    queryKey: ['categories'],
    queryFn: getCategories,
    staleTime: 300000,
  });

  const categories = categoriesData?.categories || [];

  // Add product mutation
  const addProductMutation = useMutation({
    mutationFn: async (productData) => {
      // First, upload images if any
      let imageUrls = [];
      if (productData.images.length > 0) {
        const formData = new FormData();
        productData.images.forEach((image) => {
          if (image.file) {
            formData.append('images', image.file);
          }
        });
        
        try {
          const uploadResponse = await uploadProductImages(formData);
          imageUrls = uploadResponse.urls || [];
        } catch (error) {
          // Error: 'Image upload failed:', error...
          throw new Error('Failed to upload images');
        }
      }

      // Then create the product
      const productPayload = {
        ...productData,
        images: imageUrls,
        price: parseFloat(productData.price),
        originalPrice: productData.originalPrice ? parseFloat(productData.originalPrice) : null,
        stock: parseInt(productData.stock) || 0,
      };

      return await addProduct(productPayload);
    },
    onSuccess: (data) => {
      router.push('/admin/products');
    },
    onError: (error) => {
      // Error: 'Failed to add product:', error...
      setErrors({ submit: error.message || 'Failed to add product' });
      setSubmitting(false);
    }
  });

  // Form handlers
  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
    // Clear error for this field
    if (errors[name]) {
      setErrors(prev => ({ ...prev, [name]: '' }));
    }
  };

  const handleImagesChange = (images) => {
    setFormData(prev => ({ ...prev, images }));
  };

  const handleRemoveImage = (imageId) => {
    setFormData(prev => ({
      ...prev,
      images: prev.images.filter(img => img.id !== imageId)
    }));
  };

  const handleTagsChange = (tags) => {
    setFormData(prev => ({ ...prev, tags }));
  };

  const validateForm = () => {
    const newErrors = {};

    if (!formData.name.trim()) {
      newErrors.name = 'Product name is required';
    }

    if (!formData.description.trim()) {
      newErrors.description = 'Description is required';
    }

    if (!formData.price || parseFloat(formData.price) <= 0) {
      newErrors.price = 'Valid price is required';
    }

    if (!formData.categoryId) {
      newErrors.categoryId = 'Category is required';
    }

    if (!formData.sku.trim()) {
      newErrors.sku = 'SKU is required';
    }

    if (formData.stock === '' || parseInt(formData.stock) < 0) {
      newErrors.stock = 'Stock must be 0 or greater';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setSubmitting(true);
    addProductMutation.mutate(formData);
  };

  if (!isAdmin) {
    return null;
  }

  if (categoriesLoading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
      </div>
    );
  }

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        {/* Header */}
        <div className={styles.header}>
          <button
            className={styles.backButton}
            onClick={() => router.push('/admin/products')}
          >
            <ArrowLeft size={20} />
            Back to Products
          </button>
          <h1 className={styles.title}>
            {t('title', { defaultValue: 'Add New Product' })}
          </h1>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.formGrid}>
            {/* Basic Information */}
            <div className={styles.section}>
              <h2 className={styles.sectionTitle}>Basic Information</h2>
              
              <div className={styles.formGroup}>
                <label htmlFor="name">
                  Product Name <span className={styles.required}>*</span>
                </label>
                <input
                  type="text"
                  id="name"
                  name="name"
                  value={formData.name}
                  onChange={handleInputChange}
                  className={errors.name ? styles.errorInput : ''}
                  placeholder="Enter product name"
                />
                {errors.name && (
                  <span className={styles.errorMessage}>
                    <AlertCircle size={14} /> {errors.name}
                  </span>
                )}
              </div>

              <div className={styles.formGroup}>
                <label htmlFor="description">
                  Description <span className={styles.required}>*</span>
                </label>
                <textarea
                  id="description"
                  name="description"
                  value={formData.description}
                  onChange={handleInputChange}
                  className={errors.description ? styles.errorInput : ''}
                  placeholder="Enter product description"
                  rows={5}
                />
                {errors.description && (
                  <span className={styles.errorMessage}>
                    <AlertCircle size={14} /> {errors.description}
                  </span>
                )}
              </div>

              <div className={styles.formRow}>
                <div className={styles.formGroup}>
                  <label htmlFor="category">
                    Category <span className={styles.required}>*</span>
                  </label>
                  <select
                    id="category"
                    name="categoryId"
                    value={formData.categoryId}
                    onChange={handleInputChange}
                    className={errors.categoryId ? styles.errorInput : ''}
                  >
                    <option value="">Select a category</option>
                    {categories.map(cat => (
                      <option key={cat.id} value={cat.id}>{cat.name}</option>
                    ))}
                  </select>
                  {errors.categoryId && (
                    <span className={styles.errorMessage}>
                      <AlertCircle size={14} /> {errors.categoryId}
                    </span>
                  )}
                </div>

                <div className={styles.formGroup}>
                  <label htmlFor="sku">
                    SKU <span className={styles.required}>*</span>
                  </label>
                  <input
                    type="text"
                    id="sku"
                    name="sku"
                    value={formData.sku}
                    onChange={handleInputChange}
                    className={errors.sku ? styles.errorInput : ''}
                    placeholder="Product SKU"
                  />
                  {errors.sku && (
                    <span className={styles.errorMessage}>
                      <AlertCircle size={14} /> {errors.sku}
                    </span>
                  )}
                </div>
              </div>

              <div className={styles.formGroup}>
                <label>Tags</label>
                <TagInput
                  tags={formData.tags}
                  onTagsChange={handleTagsChange}
                />
                <p className={styles.fieldHelp}>
                  Add tags to help customers find this product
                </p>
              </div>
            </div>

            {/* Pricing and Inventory */}
            <div className={styles.section}>
              <h2 className={styles.sectionTitle}>Pricing & Inventory</h2>
              
              <div className={styles.formRow}>
                <div className={styles.formGroup}>
                  <label htmlFor="price">
                    Price <span className={styles.required}>*</span>
                  </label>
                  <div className={styles.inputWithIcon}>
                    <DollarSign size={16} />
                    <input
                      type="number"
                      id="price"
                      name="price"
                      value={formData.price}
                      onChange={handleInputChange}
                      className={errors.price ? styles.errorInput : ''}
                      placeholder="0.00"
                      step="0.01"
                      min="0"
                    />
                  </div>
                  {errors.price && (
                    <span className={styles.errorMessage}>
                      <AlertCircle size={14} /> {errors.price}
                    </span>
                  )}
                </div>

                <div className={styles.formGroup}>
                  <label htmlFor="originalPrice">Original Price</label>
                  <div className={styles.inputWithIcon}>
                    <DollarSign size={16} />
                    <input
                      type="number"
                      id="originalPrice"
                      name="originalPrice"
                      value={formData.originalPrice}
                      onChange={handleInputChange}
                      placeholder="0.00"
                      step="0.01"
                      min="0"
                    />
                  </div>
                  <p className={styles.fieldHelp}>For discount display</p>
                </div>
              </div>

              <div className={styles.formRow}>
                <div className={styles.formGroup}>
                  <label htmlFor="stock">
                    Stock Quantity <span className={styles.required}>*</span>
                  </label>
                  <div className={styles.inputWithIcon}>
                    <Package size={16} />
                    <input
                      type="number"
                      id="stock"
                      name="stock"
                      value={formData.stock}
                      onChange={handleInputChange}
                      className={errors.stock ? styles.errorInput : ''}
                      placeholder="0"
                      min="0"
                    />
                  </div>
                  {errors.stock && (
                    <span className={styles.errorMessage}>
                      <AlertCircle size={14} /> {errors.stock}
                    </span>
                  )}
                </div>

                <div className={styles.formGroup}>
                  <label htmlFor="status">Status</label>
                  <select
                    id="status"
                    name="status"
                    value={formData.status}
                    onChange={handleInputChange}
                  >
                    <option value="active">Active</option>
                    <option value="draft">Draft</option>
                    <option value="archived">Archived</option>
                  </select>
                </div>
              </div>

              <div className={styles.formGroup}>
                <label className={styles.checkboxLabel}>
                  <input
                    type="checkbox"
                    name="featured"
                    checked={formData.featured}
                    onChange={handleInputChange}
                  />
                  <span>Featured Product</span>
                </label>
                <p className={styles.fieldHelp}>
                  Featured products appear in special sections
                </p>
              </div>
            </div>

            {/* Product Details */}
            <div className={styles.section}>
              <h2 className={styles.sectionTitle}>Product Details</h2>
              
              <div className={styles.formRow}>
                <div className={styles.formGroup}>
                  <label htmlFor="brand">Brand</label>
                  <input
                    type="text"
                    id="brand"
                    name="brand"
                    value={formData.brand}
                    onChange={handleInputChange}
                    placeholder="Product brand"
                  />
                </div>

                <div className={styles.formGroup}>
                  <label htmlFor="model">Model</label>
                  <input
                    type="text"
                    id="model"
                    name="model"
                    value={formData.model}
                    onChange={handleInputChange}
                    placeholder="Product model"
                  />
                </div>
              </div>

              <div className={styles.formRow}>
                <div className={styles.formGroup}>
                  <label htmlFor="condition">Condition</label>
                  <select
                    id="condition"
                    name="condition"
                    value={formData.condition}
                    onChange={handleInputChange}
                  >
                    <option value="new">New</option>
                    <option value="refurbished">Refurbished</option>
                    <option value="used">Used</option>
                  </select>
                </div>

                <div className={styles.formGroup}>
                  <label htmlFor="warranty">Warranty</label>
                  <input
                    type="text"
                    id="warranty"
                    name="warranty"
                    value={formData.warranty}
                    onChange={handleInputChange}
                    placeholder="e.g., 1 year manufacturer warranty"
                  />
                </div>
              </div>

              <div className={styles.formGroup}>
                <label htmlFor="specifications">Specifications</label>
                <textarea
                  id="specifications"
                  name="specifications"
                  value={formData.specifications}
                  onChange={handleInputChange}
                  placeholder="Enter product specifications"
                  rows={4}
                />
              </div>
            </div>

            {/* Shipping & Returns */}
            <div className={styles.section}>
              <h2 className={styles.sectionTitle}>Shipping & Returns</h2>
              
              <div className={styles.formGroup}>
                <label htmlFor="shippingInfo">Shipping Information</label>
                <textarea
                  id="shippingInfo"
                  name="shippingInfo"
                  value={formData.shippingInfo}
                  onChange={handleInputChange}
                  placeholder="Enter shipping details"
                  rows={3}
                />
              </div>

              <div className={styles.formGroup}>
                <label htmlFor="returnPolicy">Return Policy</label>
                <textarea
                  id="returnPolicy"
                  name="returnPolicy"
                  value={formData.returnPolicy}
                  onChange={handleInputChange}
                  placeholder="Enter return policy"
                  rows={3}
                />
              </div>
            </div>

            {/* Images */}
            <div className={`${styles.section} ${styles.fullWidth}`}>
              <h2 className={styles.sectionTitle}>Product Images</h2>
              <ImageUploader
                images={formData.images}
                onImagesChange={handleImagesChange}
                onRemove={handleRemoveImage}
              />
            </div>
          </div>

          {/* Error message */}
          {errors.submit && (
            <div className={styles.submitError}>
              <AlertCircle size={16} />
              {errors.submit}
            </div>
          )}

          {/* Actions */}
          <div className={styles.formActions}>
            <button
              type="button"
              className={styles.cancelButton}
              onClick={() => router.push('/admin/products')}
            >
              Cancel
            </button>
            <button
              type="submit"
              className={styles.submitButton}
              disabled={submitting}
            >
              {submitting ? (
                <>
                  <Loader2 size={16} className={styles.spinning} />
                  Adding Product...
                </>
              ) : (
                <>
                  <Save size={16} />
                  Add Product
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </ErrorBoundary>
  );
};

export default AddProduct;