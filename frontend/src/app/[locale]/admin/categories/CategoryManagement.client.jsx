"use client";

import React, {useState, useCallback, useMemo} from 'react';
import {useRouter} from 'next/navigation';
import {useTranslations} from 'next-intl';
import {useQuery, useMutation, useQueryClient} from '@tanstack/react-query';
import {
    Plus,
    Edit2,
    Trash2,
    ChevronRight,
    ChevronDown,
    Folder,
    FolderOpen,
    AlertCircle,
    Search,
    RefreshCw,
    Move,
    Copy,
    Eye,
    EyeOff,
    Package,
    Hash,
    Calendar,
    Archive,
    ShoppingCart,
    Settings,
    FileText
} from 'lucide-react';
import {useAuth} from '@/context/AuthContext';
import {useUserRole} from '@/hooks/useUserRole';
import {
    getCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    archiveCategory,
    reorderCategories,
    getCategoryStats
} from '@/api/adminApi';
import { fetchSubCategories } from '@/api/categories';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './CategoryManagement.module.css';

// Category type configuration
const CATEGORY_TYPES = {
    marketplace: {
        label: 'Marketplace',
        icon: ShoppingCart,
        description: 'Product categories for marketplace'
    },
    service: {
        label: 'Services', 
        icon: Settings,
        description: 'Service categories'
    },
    posts: {
        label: 'Posts',
        icon: FileText,
        description: 'Blog and content categories'
    }
};

// Category tree item component
const CategoryTreeItem = ({
                              category,
                              level = 0,
                              onEdit,
                              onDelete,
                              onToggle,
                              onAddSubcategory,
                              onToggleActive,
                              onArchive,
                              expandedCategories,
                              categoryStats
                          }) => {
    const [showActions, setShowActions] = useState(false);
    const hasSubcategories = category.subcategories && category.subcategories.length > 0;
    const isExpanded = expandedCategories.has(category.id);
    const stats = categoryStats[category.id] || {productCount: 0, activeProducts: 0};
    const isParent = level === 0; // First level categories are parents

    const handleToggle = (e) => {
        e.stopPropagation();
        if (hasSubcategories) {
            onToggle(category.id);
        }
    };

    const handleRowClick = (e) => {
        // Don't toggle if clicking on action buttons
        if (e.target.closest(`.${styles.categoryActions}`)) {
            return;
        }
        if (hasSubcategories) {
            onToggle(category.id);
        }
    };

    return (
        <>
            <div
                className={`${styles.categoryItem} ${styles[`level${Math.min(level, 4)}`]} ${isParent ? styles.parentCategory : styles.childCategory} ${hasSubcategories ? styles.hasChildren : ''}`}
                onMouseEnter={() => setShowActions(true)}
                onMouseLeave={() => setShowActions(false)}
                onClick={handleRowClick}
                style={{cursor: hasSubcategories ? 'pointer' : 'default'}}
            >
                <div className={styles.categoryContent}>
                    <div className={styles.expandButton}>
                        {hasSubcategories ? (
                            isExpanded ? <ChevronDown size={16}/> : <ChevronRight size={16}/>
                        ) : (
                            <span className={styles.noExpand}/>
                        )}
                    </div>

                    <div className={styles.categoryIcon}>
                        {hasSubcategories ? (
                            isExpanded ? <FolderOpen size={16}/> : <Folder size={16}/>
                        ) : (
                            <Package size={16}/>
                        )}
                    </div>

                    <div className={styles.categoryInfo}>
                        <div className={styles.categoryName}>
                            {category.description || category.slug}
                            {isParent && <span className={styles.parentBadge}>Parent</span>}
                        </div>
                        <div className={styles.categoryMeta}>
                            <span className={styles.categorySlug}>{category.slug}</span>
                            <span className={styles.categoryType}>{category.categoryType}</span>
                            <span className={styles.categoryLang}>{category.lang || 'en'}</span>
                        </div>
                    </div>

                    <div className={styles.categoryStats}>
                        <span className={styles.statItem}>
                            {stats.productCount || 0}
                        </span>
                        <span className={styles.statLabel}>items</span>
                    </div>

                    <div className={styles.categoryStatus}>
                        <span className={`${styles.statusBadge} ${category.isActive ? styles.active : styles.inactive}`}>
                            {category.isActive ? 'Active' : 'Inactive'}
                        </span>
                    </div>

                    {showActions && (
                        <div className={styles.categoryActions}>
                            <button
                                className={styles.actionButton}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    onAddSubcategory(category);
                                }}
                                title="Add subcategory"
                            >
                                <Plus size={14}/>
                            </button>
                            <button
                                className={styles.actionButton}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    onEdit(category);
                                }}
                                title="Edit category"
                            >
                                <Edit2 size={14}/>
                            </button>
                            <button
                                className={styles.actionButton}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    onToggleActive(category);
                                }}
                                title={category.isActive ? "Deactivate category" : "Activate category"}
                            >
                                {category.isActive ? <EyeOff size={14}/> : <Eye size={14}/>}
                            </button>
                            <button
                                className={styles.actionButton}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    onDelete(category);
                                }}
                                title="Delete category"
                                disabled={hasSubcategories || stats.productCount > 0}
                            >
                                <Trash2 size={14}/>
                            </button>
                        </div>
                    )}
                </div>
            </div>

            {hasSubcategories && isExpanded && (
                <div className={styles.categoryChildren}>
                    {category.subcategories.map((subcategory) => (
                        <CategoryTreeItem
                            key={subcategory.id}
                            category={subcategory}
                            level={level + 1}
                            onEdit={onEdit}
                            onDelete={onDelete}
                            onToggle={onToggle}
                            onAddSubcategory={onAddSubcategory}
                            onToggleActive={onToggleActive}
                            onArchive={onArchive}
                            expandedCategories={expandedCategories}
                            categoryStats={categoryStats}
                        />
                    ))}
                </div>
            )}
        </>
    );
};

// Category form modal
const CategoryFormModal = ({category, parentCategory, onClose, onSave, selectedCategoryType}) => {
    const t = useTranslations('CategoryManagement');
    const [formData, setFormData] = useState({
        description: category?.description || '',
        slug: category?.slug || '',
        parentId: parentCategory?.id || category?.parentId || '',
        isActive: category?.isActive !== false,
        googleCategoryId: category?.googleCategoryId || '',
        tags: category?.tags || [],
        seoTitle: category?.seoTitle || '',
        seoKeywords: category?.seoKeywords || [],
        seoDesc: category?.seoDesc || '',
        categoryType: category?.categoryType || selectedCategoryType,
        lang: category?.lang || 'en'
    });
    const [errors, setErrors] = useState({});

    const validateForm = () => {
        const newErrors = {};
        if (!formData.description.trim()) {
            newErrors.description = 'Category description is required';
        }
        if (!formData.slug.trim()) {
            newErrors.slug = 'Category slug is required';
        } else if (!/^[a-z0-9-]+$/.test(formData.slug)) {
            newErrors.slug = 'Slug must contain only lowercase letters, numbers, and hyphens';
        }
        if (!formData.categoryType) {
            newErrors.categoryType = 'Category type is required';
        }
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        if (validateForm()) {
            onSave(formData);
        }
    };

    const generateSlug = (name) => {
        return name
            .toLowerCase()
            .replace(/[^a-z0-9]+/g, '-')
            .replace(/^-+|-+$/g, '');
    };

    return (
        <div className={styles.modalOverlay} onClick={onClose}>
            <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
                <div className={styles.modalHeader}>
                    <h3>{category ? 'Edit Category' : 'Create Category'}</h3>
                    <button className={styles.closeButton} onClick={onClose}>×</button>
                </div>

                <form onSubmit={handleSubmit} className={styles.categoryForm}>
                    {parentCategory && (
                        <div className={styles.formGroup}>
                            <label>Parent Category</label>
                            <div className={styles.parentInfo}>
                                <Folder size={16}/>
                                <span>{parentCategory.description}</span>
                            </div>
                        </div>
                    )}

                    <div className={styles.formGroup}>
                        <label htmlFor="description">Description *</label>
                        <input
                            type="text"
                            id="description"
                            value={formData.description}
                            onChange={(e) => {
                                setFormData({
                                    ...formData,
                                    description: e.target.value,
                                    slug: generateSlug(e.target.value)
                                });
                            }}
                            className={errors.description ? styles.errorInput : ''}
                            placeholder="e.g., Electronics"
                        />
                        {errors.description && <span className={styles.errorText}>{errors.description}</span>}
                    </div>

                    <div className={styles.formGroup}>
                        <label htmlFor="slug">Slug *</label>
                        <input
                            type="text"
                            id="slug"
                            value={formData.slug}
                            onChange={(e) => setFormData({...formData, slug: e.target.value})}
                            className={errors.slug ? styles.errorInput : ''}
                            placeholder="e.g., electronics"
                        />
                        {errors.slug && <span className={styles.errorText}>{errors.slug}</span>}
                    </div>

                    <div className={styles.formRow}>
                        <div className={styles.formGroup}>
                            <label htmlFor="categoryType">Category Type *</label>
                            <select
                                id="categoryType"
                                value={formData.categoryType}
                                onChange={(e) => setFormData({...formData, categoryType: e.target.value})}
                                className={errors.categoryType ? styles.errorInput : ''}
                            >
                                <option value="">Select type</option>
                                {Object.entries(CATEGORY_TYPES).map(([key, config]) => (
                                    <option key={key} value={key}>{config.label}</option>
                                ))}
                            </select>
                            {errors.categoryType && <span className={styles.errorText}>{errors.categoryType}</span>}
                        </div>

                        <div className={styles.formGroup}>
                            <label htmlFor="lang">Language</label>
                            <select
                                id="lang"
                                value={formData.lang}
                                onChange={(e) => setFormData({...formData, lang: e.target.value})}
                            >
                                <option value="en">English</option>
                                <option value="de">German</option>
                                <option value="it">Italian</option>
                                <option value="pl">Polish</option>
                            </select>
                        </div>
                    </div>

                    <div className={styles.formGroup}>
                        <label htmlFor="googleCategoryId">Google Category ID</label>
                        <input
                            type="text"
                            id="googleCategoryId"
                            value={formData.googleCategoryId}
                            onChange={(e) => setFormData({...formData, googleCategoryId: e.target.value})}
                            placeholder="Google Merchant Center category ID"
                        />
                    </div>

                    <div className={styles.formGroup}>
                        <label htmlFor="tags">Tags (comma separated)</label>
                        <input
                            type="text"
                            id="tags"
                            value={formData.tags.join(', ')}
                            onChange={(e) => setFormData({...formData, tags: e.target.value.split(',').map(tag => tag.trim()).filter(Boolean)})}
                            placeholder="tag1, tag2, tag3"
                        />
                    </div>

                    <div className={styles.formGroup}>
                        <label htmlFor="seoTitle">SEO Title</label>
                        <input
                            type="text"
                            id="seoTitle"
                            value={formData.seoTitle}
                            onChange={(e) => setFormData({...formData, seoTitle: e.target.value})}
                            placeholder="SEO optimized title"
                        />
                    </div>

                    <div className={styles.formGroup}>
                        <label htmlFor="seoKeywords">SEO Keywords (comma separated)</label>
                        <input
                            type="text"
                            id="seoKeywords"
                            value={formData.seoKeywords.join(', ')}
                            onChange={(e) => setFormData({...formData, seoKeywords: e.target.value.split(',').map(kw => kw.trim()).filter(Boolean)})}
                            placeholder="keyword1, keyword2, keyword3"
                        />
                    </div>

                    <div className={styles.formGroup}>
                        <label htmlFor="seoDesc">SEO Description</label>
                        <textarea
                            id="seoDesc"
                            value={formData.seoDesc}
                            onChange={(e) => setFormData({...formData, seoDesc: e.target.value})}
                            rows={3}
                            placeholder="SEO optimized description"
                        />
                    </div>

                    <div className={styles.formGroup}>
                        <label className={styles.checkboxLabel}>
                            <input
                                type="checkbox"
                                checked={formData.isActive}
                                onChange={(e) => setFormData({...formData, isActive: e.target.checked})}
                            />
                            <span>Active</span>
                        </label>
                    </div>

                    <div className={styles.modalActions}>
                        <button type="button" onClick={onClose} className={styles.cancelButton}>
                            Cancel
                        </button>
                        <button type="submit" className={styles.saveButton}>
                            {category ? 'Update' : 'Create'} Category
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

// Main component
const CategoryManagement = () => {
    // Handle missing translations gracefully
    let t;
    try {
        t = useTranslations('CategoryManagement');
    } catch (e) {
        // Fallback function for missing translations
        t = (key, options) => options?.defaultValue || key;
    }
    
    const router = useRouter();
    const queryClient = useQueryClient();
    const {user} = useAuth();
    const {isAdmin} = useUserRole();

    const [searchTerm, setSearchTerm] = useState('');
    const [expandedCategories, setExpandedCategories] = useState(new Set());
    const [showInactive, setShowInactive] = useState(false);
    const [editingCategory, setEditingCategory] = useState(null);
    const [parentCategory, setParentCategory] = useState(null);
    const [showForm, setShowForm] = useState(false);
    const [mounted, setMounted] = useState(false);
    const [selectedLanguage, setSelectedLanguage] = useState('all');
    const [selectedCategoryType, setSelectedCategoryType] = useState('marketplace'); // Default to marketplace

    // Ensure client-side only rendering for initial load
    React.useEffect(() => {
        setMounted(true);
    }, []);

    // Create a refetch function for both queries
    const refetch = useCallback(async () => {
        await Promise.all([
            queryClient.invalidateQueries(['adminCategories', selectedCategoryType, selectedLanguage])
        ]);
    }, [queryClient, selectedCategoryType, selectedLanguage]);

    // Fetch all categories in a single call
    const {data: categoriesData, isLoading, error} = useQuery({
        queryKey: ['adminCategories', selectedCategoryType, selectedLanguage],
        queryFn: async () => {
            const result = await getCategories({
                page: 1,
                pageSize: 1000, // Get all categories for tree view
                sortBy: 'slug',
                sortOrder: 'asc',
                categoryType: selectedCategoryType,
                ...(selectedLanguage !== 'all' && { lang: selectedLanguage })
            });
            return result;
        },
        staleTime: 300000,
        retry: 1,
        enabled: mounted
    });

    // Add console logging to debug the data
    React.useEffect(() => {

        if (categoriesData) {

            if (categoriesData.categories) {
                const categories = categoriesData.categories;
                
                // Debug: Log the raw categories with detailed info

                // Log parentId analysis with more detail
                const parentAnalysis = categories.map(cat => ({
                    id: cat.id,
                    description: cat.description,
                    parentId: cat.parentId,
                    parentIdType: typeof cat.parentId,
                    parentIdLength: cat.parentId?.length,
                    hasParentId: !!(cat.parentId && cat.parentId.trim() !== ''),
                    categoryType: cat.categoryType,
                    isActive: cat.isActive
                }));

                // Show unique parentId values
                const uniqueParentIds = [...new Set(categories.map(c => c.parentId))];

                // Check if parentId values are actually category IDs that exist in our dataset
                const categoryIds = new Set(categories.map(c => c.id));
                const validParentIds = uniqueParentIds.filter(pid => pid && categoryIds.has(pid));
                const invalidParentIds = uniqueParentIds.filter(pid => pid && !categoryIds.has(pid));
                
                // Valid parent IDs logged for debugging
                // Invalid parent IDs logged for debugging
                
                // Count different types with better logic
                const realParents = categories.filter(c => {
                    // A real parent either has no parentId, or its parentId doesn't exist in current dataset
                    return !c.parentId || c.parentId.trim() === '' || !categoryIds.has(c.parentId);
                });
                
                const realChildren = categories.filter(c => {
                    // A real child has a parentId that exists in the current dataset
                    return c.parentId && c.parentId.trim() !== '' && categoryIds.has(c.parentId);
                });

                // Category analysis logged for debugging
            }
        } else if (error) {
            // Error: 'API Error:', error...
        } else if (isLoading) {
            
        } else {
            
        }
    }, [categoriesData, isLoading, error, selectedCategoryType, selectedLanguage]);

    // Fetch category stats
    const {data: categoryStats = {}} = useQuery({
        queryKey: ['categoryStats', selectedCategoryType],
        queryFn: async () => {
            try {
                const stats = await getCategoryStats();
                return stats.reduce((acc, stat) => {
                    acc[stat.categoryId] = stat;
                    return acc;
                }, {});
            } catch (error) {
                // Error: 'Error fetching category stats:', error...
                return {};
            }
        },
        staleTime: 300000,
        enabled: mounted && !!categoriesData // Only run after mounted and categories loaded
    });

    // Mutations
    const createMutation = useMutation({
        mutationFn: createCategory,
        onSuccess: () => {
            queryClient.invalidateQueries(['adminCategories']);
            queryClient.invalidateQueries(['categoryStats']);
            setShowForm(false);
            setEditingCategory(null);
            setParentCategory(null);
        },
    });

    const updateMutation = useMutation({
        mutationFn: ({id, ...data}) => updateCategory(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries(['adminCategories']);
            queryClient.invalidateQueries(['categoryStats']);
            setShowForm(false);
            setEditingCategory(null);
        },
    });

    const deleteMutation = useMutation({
        mutationFn: ({id, userId}) => deleteCategory(id, userId),
        onSuccess: () => {
            queryClient.invalidateQueries(['adminCategories']);
            queryClient.invalidateQueries(['categoryStats']);
        },
    });

    const archiveMutation = useMutation({
        mutationFn: ({id, userId}) => archiveCategory(id, userId),
        onSuccess: () => {
            queryClient.invalidateQueries(['adminCategories']);
            queryClient.invalidateQueries(['categoryStats']);
        },
    });

    const toggleActiveMutation = useMutation({
        mutationFn: ({id, isActive}) => updateCategory(id, { isActive }),
        onSuccess: () => {
            queryClient.invalidateQueries(['adminCategories']);
        },
    });

    // Build category tree using the application's standard pattern
    const categoryTree = useMemo(() => {
        if (!categoriesData?.categories) {
            
            return [];
        }

        const categories = categoriesData.categories;
        // Processing categories data

        // Use the same pattern as the rest of the application:
        // 1. Main categories come from the API
        // 2. Subcategories are fetched separately using fetchSubCategories
        // 3. Each category has a subcategories array (not children)
        
        // Map categories to match application format
        const mappedCategories = categories.map(cat => ({
            ...cat,
            name: cat.name || cat.description || 'Unknown Category',
            subcategories: cat.subcategories || [], // Use subcategories, not children
            // Mark if subcategories have been fetched (for lazy loading)
            subcategoriesFetched: Array.isArray(cat.subcategories) && cat.subcategories.length >= 0
        }));

        // Sort categories like other components do
        const sortedCategories = mappedCategories.sort((a, b) => {
            const aName = a.name || a.description || a.slug || '';
            const bName = b.name || b.description || b.slug || '';
            return aName.localeCompare(bName);
        });

        return sortedCategories;
    }, [categoriesData]);

    // Filter categories based on search and active status
    const filteredTree = useMemo(() => {
        if (!searchTerm && showInactive) return categoryTree;

        const filterCategory = (category) => {
            const matchesSearch = !searchTerm ||
                category.slug?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                category.description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                category.name?.toLowerCase().includes(searchTerm.toLowerCase());

            const matchesActive = showInactive || category.isActive;

            const filteredSubcategories = category.subcategories
                ?.map(subcategory => filterCategory(subcategory))
                .filter(Boolean) || [];

            if (matchesSearch && matchesActive) {
                return {...category, subcategories: filteredSubcategories};
            } else if (filteredSubcategories.length > 0) {
                return {...category, subcategories: filteredSubcategories};
            }

            return null;
        };

        return categoryTree
            .map(cat => filterCategory(cat))
            .filter(Boolean);
    }, [categoryTree, searchTerm, showInactive]);

    // Handlers
    const handleToggleCategory = useCallback((categoryId) => {
        setExpandedCategories(prev => {
            const newSet = new Set(prev);
            if (newSet.has(categoryId)) {
                newSet.delete(categoryId);
            } else {
                newSet.add(categoryId);
            }
            return newSet;
        });
    }, []);

    const handleExpandAll = useCallback(() => {
        const allIds = new Set();
        const collectIds = (categories) => {
            categories.forEach(cat => {
                if (cat.subcategories && cat.subcategories.length > 0) {
                    allIds.add(cat.id);
                    collectIds(cat.subcategories);
                }
            });
        };
        collectIds(categoryTree);
        setExpandedCategories(allIds);
    }, [categoryTree]);

    const handleCollapseAll = useCallback(() => {
        setExpandedCategories(new Set());
    }, []);

    const handleEdit = useCallback((category) => {
        setEditingCategory(category);
        setShowForm(true);
    }, []);

    const handleAddSubcategory = useCallback((parent) => {
        setParentCategory(parent);
        setEditingCategory(null);
        setShowForm(true);
    }, []);

    const handleDelete = useCallback((category) => {
        if (confirm(`Are you sure you want to delete "${category.description}"?`)) {
            const userId = user?.id;
            deleteMutation.mutate({id: category.id, userId});
        }
    }, [deleteMutation, user]);

    const handleArchive = useCallback((category) => {
        if (confirm(`Are you sure you want to archive "${category.description}"?`)) {
            const userId = user?.id;
            archiveMutation.mutate({id: category.id, userId});
        }
    }, [archiveMutation, user]);

    const handleToggleActive = useCallback((category) => {
        toggleActiveMutation.mutate({
            id: category.id, 
            isActive: !category.isActive
        });
    }, [toggleActiveMutation]);

    const handleSave = useCallback((formData) => {
        if (editingCategory) {
            updateMutation.mutate({id: editingCategory.id, ...formData});
        } else {
            createMutation.mutate(formData);
        }
    }, [editingCategory, createMutation, updateMutation]);

    // Calculate stats
    const totalCategories = categoriesData?.categories?.length || 0;
    const activeCategories = categoriesData?.categories?.filter(c => c.isActive).length || 0;
    const parentCategories = categoriesData?.categories?.filter(c => !c.parentId || c.parentId.trim() === '').length || 0;
    const childCategories = totalCategories - parentCategories;

    // Show loading state only on client side to avoid hydration mismatch
    if (!mounted || isLoading) {
        return (
            <div className={styles.loadingContainer}>
                <LoadingSpinner/>
            </div>
        );
    }

    // Handle error state
    if (error) {
        return (
            <div className={styles.container}>
                <div className={styles.errorState}>
                    <AlertCircle size={48} className={styles.errorIcon} />
                    <h2 className={styles.errorTitle}>Failed to load categories</h2>
                    <p className={styles.errorMessage}>{error.message || 'An error occurred while fetching categories'}</p>
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

    const currentTypeConfig = CATEGORY_TYPES[selectedCategoryType];

    return (
        <ErrorBoundary>
            <div className={styles.container}>
                {/* Header */}
                <div className={styles.header}>
                    <div>
                        <h1 className={styles.title}>
                            {t('title', {defaultValue: 'Category Management'})}
                        </h1>
                        <p className={styles.subtitle}>
                            {t('subtitle', {defaultValue: `Manage ${currentTypeConfig?.description || 'categories'}`})}
                        </p>
                    </div>
                    <div className={styles.headerActions}>
                        <button
                            className={styles.createButton}
                            onClick={() => {
                                setEditingCategory(null);
                                setParentCategory(null);
                                setShowForm(true);
                            }}
                        >
                            <Plus size={16}/>
                            {t('createCategory', {defaultValue: 'Create Category'})}
                        </button>
                        <button
                            className={styles.refreshButton}
                            onClick={() => refetch()}
                        >
                            <RefreshCw size={16}/>
                        </button>

                    </div>
                </div>

                {/* Category Type Switcher */}
                <div className={styles.typeSwitcher}>
                    {Object.entries(CATEGORY_TYPES).map(([key, config]) => {
                        const IconComponent = config.icon;
                        return (
                            <button
                                key={key}
                                className={`${styles.typeButton} ${selectedCategoryType === key ? styles.active : ''}`}
                                onClick={() => setSelectedCategoryType(key)}
                            >
                                <IconComponent size={18} />
                                <span>{config.label}</span>
                            </button>
                        );
                    })}
                </div>

                {/* Stats */}
                <div className={styles.statsRow}>
                    <div className={styles.statCard}>
                        <div className={styles.statIcon}>
                            <Folder size={20}/>
                        </div>
                        <div className={styles.statContent}>
                            <span className={styles.statValue}>{totalCategories}</span>
                            <span className={styles.statLabel}>Total Categories</span>
                        </div>
                    </div>
                    <div className={styles.statCard}>
                        <div className={styles.statIcon}>
                            <FolderOpen size={20}/>
                        </div>
                        <div className={styles.statContent}>
                            <span className={styles.statValue}>{parentCategories}</span>
                            <span className={styles.statLabel}>Parent Categories</span>
                        </div>
                    </div>
                    <div className={styles.statCard}>
                        <div className={styles.statIcon}>
                            <Package size={20}/>
                        </div>
                        <div className={styles.statContent}>
                            <span className={styles.statValue}>{childCategories}</span>
                            <span className={styles.statLabel}>Child Categories</span>
                        </div>
                    </div>
                    <div className={styles.statCard}>
                        <div className={styles.statIcon}>
                            <Eye size={20}/>
                        </div>
                        <div className={styles.statContent}>
                            <span className={styles.statValue}>{activeCategories}</span>
                            <span className={styles.statLabel}>Active Categories</span>
                        </div>
                    </div>
                </div>

                {/* Controls */}
                <div className={styles.controls}>
                    <div className={styles.searchBox}>
                        <Search size={20} className={styles.searchIcon}/>
                        <input
                            type="text"
                            placeholder={t('searchPlaceholder', {defaultValue: 'Search categories...'})}
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            className={styles.searchInput}
                        />
                    </div>
                    <div className={styles.controlButtons}>
                        <select 
                            className={styles.languageFilter}
                            value={selectedLanguage}
                            onChange={(e) => setSelectedLanguage(e.target.value)}
                        >
                            <option value="all">All Languages</option>
                            <option value="en">English (EN)</option>
                            <option value="de">Deutsch (DE)</option>
                            <option value="it">Italiano (IT)</option>
                            <option value="pl">Polski (PL)</option>
                        </select>
                        <button
                            className={styles.toggleButton}
                            onClick={() => setShowInactive(!showInactive)}
                            title={showInactive ? 'Hide inactive categories' : 'Show inactive categories'}
                        >
                            {showInactive ? <EyeOff size={16}/> : <Eye size={16}/>}
                        </button>
                        <button
                            className={styles.toggleButton}
                            onClick={handleExpandAll}
                            title="Expand all categories"
                        >
                            <ChevronDown size={16}/>
                        </button>
                        <button
                            className={styles.toggleButton}
                            onClick={handleCollapseAll}
                            title="Collapse all categories"
                        >
                            <ChevronRight size={16}/>
                        </button>
                    </div>
                </div>

                {/* Category Tree */}
                <div className={styles.categoryTree}>
                    {filteredTree.length === 0 ? (
                        <div className={styles.emptyState}>
                            {searchTerm ? (
                                <>
                                    <AlertCircle size={48} className={styles.emptyIcon}/>
                                    <h3 className={styles.emptyTitle}>{t('noCategoriesFound', {defaultValue: 'No categories found'})}</h3>
                                    <p className={styles.emptyDescription}>{t('adjustSearchTerm', {defaultValue: 'Try adjusting your search term'})}</p>
                                </>
                            ) : (
                                <>
                                    {React.createElement(currentTypeConfig?.icon || Folder, { size: 48, className: styles.emptyIcon })}
                                    <h3 className={styles.emptyTitle}>{t('noCategories', {defaultValue: `No ${currentTypeConfig?.label.toLowerCase()} categories yet`})}</h3>
                                    <p className={styles.emptyDescription}>{t('noCategoriesDesc', {defaultValue: 'Create your first category to get started'})}</p>
                                    <button
                                        className={styles.createButton}
                                        onClick={() => {
                                            setEditingCategory(null);
                                            setParentCategory(null);
                                            setShowForm(true);
                                        }}
                                    >
                                        <Plus size={16}/>
                                        {t('createCategory', {defaultValue: 'Create Category'})}
                                    </button>
                                </>
                            )}
                        </div>
                    ) : (
                        filteredTree.map(category => (
                            <CategoryTreeItem
                                key={category.id}
                                category={category}
                                level={0}
                                onEdit={handleEdit}
                                onDelete={handleDelete}
                                onToggle={handleToggleCategory}
                                onAddSubcategory={handleAddSubcategory}
                                onToggleActive={handleToggleActive}
                                onArchive={handleArchive}
                                expandedCategories={expandedCategories}
                                categoryStats={categoryStats}
                            />
                        ))
                    )}
                </div>

                {/* Category Form Modal */}
                {showForm && (
                    <CategoryFormModal
                        category={editingCategory}
                        parentCategory={parentCategory}
                        selectedCategoryType={selectedCategoryType}
                        onClose={() => {
                            setShowForm(false);
                            setEditingCategory(null);
                            setParentCategory(null);
                        }}
                        onSave={handleSave}
                    />
                )}
            </div>
        </ErrorBoundary>
    );
};

export default CategoryManagement;