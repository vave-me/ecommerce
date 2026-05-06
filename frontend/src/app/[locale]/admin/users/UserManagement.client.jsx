"use client";

import React, {useState, useEffect, useCallback, useMemo} from 'react';
import {useRouter} from 'next/navigation';
import {useTranslations} from 'next-intl';
import {useQuery, useMutation, useQueryClient} from '@tanstack/react-query';
import {
    Search,
    Filter,
    MoreVertical,
    UserCheck,
    UserX,
    Shield,
    AlertCircle,
    Download,
    RefreshCw,
    ChevronLeft,
    ChevronRight,
    Mail,
    Calendar,
    Building,
    User,
    UserPlus,
    Settings,
    TrendingUp,
    Activity,
    Clock,
    CheckCircle,
    XCircle,
    Eye,
    Edit,
    BarChart3,
    LogOut,
    MoreHorizontal,
    ArrowUpDown,
    Check,
    X,
    ChevronsLeft,
    ChevronsRight,
    SlidersHorizontal,
    Users
} from 'lucide-react';
import {useAuth} from '@/context/AuthContext';
import {useUserRole} from '@/hooks/useUserRole';
import {
    listUsers,
    enableUser,
    disableUser,
    updateUser,
    clearUserTokens,
    createAdminUser
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './UserManagement.module.css';

const UserRow = ({user, onAction, selected, onSelect}) => {
    const [showMenu, setShowMenu] = useState(false);
    const router = useRouter();

    const getRoleBadgeClass = (role) => {
        const normalizedRole = (role || 'customer').toLowerCase();
        switch (normalizedRole) {
            case 'admin':
                return styles.roleBadgeAdmin;
            case 'business':
                return styles.roleBadgeBusiness;
            case 'customer':
            default:
                return styles.roleBadgeCustomer;
        }
    };

    const getStatusBadgeClass = (status) => {
        return status === 'active' ? styles.statusActive : styles.statusInactive;
    };

    return (
        <tr className={`${styles.userRow} ${selected ? styles.selectedRow : ''}`} onClick={() => onAction('view', user)}>
            <td className={styles.checkboxCell} onClick={(e) => e.stopPropagation()}>
                <input
                    type="checkbox"
                    checked={selected}
                    onChange={(e) => onSelect(user.id, e.target.checked)}
                    className={styles.checkbox}
                />
            </td>
            <td className={styles.userCell}>
                <div className={styles.userInfo}>
                    <div className={styles.userAvatar}>
                        {user.thumbnail ? (
                            <img src={user.thumbnail} alt={user.userName || user.email} onError={(e) => {
                                e.target.style.display = 'none';
                                e.target.nextSibling.style.display = 'flex';
                            }}/>
                        ) : null}
                        <div style={{display: user.thumbnail ? 'none' : 'flex', alignItems: 'center', justifyContent: 'center', width: '100%', height: '100%'}}>
                            <User size={20}/>
                        </div>
                    </div>
                    <div>
                        <div className={styles.userName}>
                            {user.userName || `${user.firstName || ''} ${user.lastName || ''}`.trim() || user.email?.split('@')[0] || 'Unknown User'}
                        </div>
                        <div className={styles.userEmail}>{user.email}</div>
                    </div>
                </div>
            </td>
            <td className={styles.roleCell}>
        <span className={`${styles.roleBadge} ${getRoleBadgeClass(user.role)}`}>
          {(user.role || 'customer').charAt(0).toUpperCase() + (user.role || 'customer').slice(1).toLowerCase()}
        </span>
            </td>
            <td className={styles.statusCell}>
        <span className={`${styles.statusBadge} ${getStatusBadgeClass(user.enabled ? 'active' : 'inactive')}`}>
          {user.enabled ? 'active' : 'inactive'}
        </span>
            </td>
            <td className={styles.dateCell}>
                <div className={styles.dateInfo}>
                    <Calendar size={14}/>
                    <span>{new Date(user.createdAt).toLocaleDateString()}</span>
                </div>
            </td>
            <td className={styles.dateCell}>
                <div className={styles.dateInfo}>
                    <Calendar size={14}/>
                    <span>{user.lastActive ? new Date(user.lastActive).toLocaleDateString() : 'Never'}</span>
                </div>
            </td>
            <td className={styles.actionCell} onClick={(e) => e.stopPropagation()}>
                <div className={styles.actionButtons}>
                    <button
                        className={styles.actionIconButton}
                        onClick={(e) => {
                            e.stopPropagation();
                            onAction('view', user);
                        }}
                        title="View Details"
                    >
                        <Eye size={16}/>
                    </button>
                    <button
                        className={styles.actionIconButton}
                        onClick={(e) => {
                            e.stopPropagation();
                            onAction('edit', user);
                        }}
                        title="Edit User"
                    >
                        <Edit size={16}/>
                    </button>
                    <button
                        className={styles.actionIconButton}
                        onClick={(e) => {
                            e.stopPropagation();
                            onAction('viewMetrics', user);
                        }}
                        title="View Metrics"
                    >
                        <BarChart3 size={16}/>
                    </button>
                    <div className={styles.actionMenu}>
                        <button
                            className={styles.actionIconButton}
                            onClick={(e) => {
                                e.stopPropagation();
                                setShowMenu(!showMenu);
                            }}
                            aria-label="More actions"
                        >
                            <MoreHorizontal size={16}/>
                        </button>
                        {showMenu && (
                            <>
                                <div
                                    className={styles.menuOverlay}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        setShowMenu(false);
                                    }}
                                />
                                <div className={styles.menuDropdown}>
                                    <button onClick={(e) => {
                                        e.stopPropagation();
                                        onAction('toggle', user);
                                        setShowMenu(false);
                                    }}>
                                        {user.enabled ? (
                                            <><UserX size={14} /> Deactivate</>
                                        ) : (
                                            <><UserCheck size={14} /> Activate</>
                                        )}
                                    </button>
                                    <button onClick={(e) => {
                                        e.stopPropagation();
                                        onAction('changeRole', user);
                                        setShowMenu(false);
                                    }}>
                                        <Shield size={14} />
                                        Change Role
                                    </button>
                                    <button onClick={(e) => {
                                        e.stopPropagation();
                                        onAction('forceLogout', user);
                                        setShowMenu(false);
                                    }}>
                                        <LogOut size={14} />
                                        Force Logout
                                    </button>
                                </div>
                            </>
                        )}
                    </div>
                </div>
            </td>
        </tr>
    );
};

const UserManagement = ({ locale = 'en' }) => {
    const t = useTranslations('UserManagement');
    const router = useRouter();
    const queryClient = useQueryClient();
    const {isAdmin} = useUserRole();

    const [searchTerm, setSearchTerm] = useState('');
    const [filterRole, setFilterRole] = useState('all');
    const [filterStatus, setFilterStatus] = useState('all');
    const [currentPage, setCurrentPage] = useState(1);
    const [showFilters, setShowFilters] = useState(false);
    const [selectedUsers, setSelectedUsers] = useState(new Set());
    const [sortBy, setSortBy] = useState('createdAt');
    const [sortOrder, setSortOrder] = useState('desc');
    const [showAddUser, setShowAddUser] = useState(false);
    const [filterDateRange, setFilterDateRange] = useState('all');
    const [viewMode, setViewMode] = useState('table'); // table or grid

    const itemsPerPage = 10;

    // Fetch users data
    const {data: usersData, isLoading, refetch} = useQuery({
        queryKey: ['adminUsers'],
        queryFn: async () => {
            const response = await listUsers();
            let users = response.users || [];

            // Client-side filtering since API doesn't support these filters
            if (searchTerm) {
                users = users.filter(user =>
                    user.userName?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                    user.email?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                    `${user.firstName} ${user.lastName}`.toLowerCase().includes(searchTerm.toLowerCase())
                );
            }

            if (filterRole !== 'all') {
                users = users.filter(user => (user.role || 'customer') === filterRole);
            }

            if (filterStatus !== 'all') {
                users = users.filter(user =>
                    filterStatus === 'active' ? user.enabled : !user.enabled
                );
            }

            // Date range filtering
            if (filterDateRange !== 'all') {
                const now = new Date();
                let startDate;
                switch (filterDateRange) {
                    case '7days':
                        startDate = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
                        break;
                    case '30days':
                        startDate = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
                        break;
                    case '90days':
                        startDate = new Date(now.getTime() - 90 * 24 * 60 * 60 * 1000);
                        break;
                }
                if (startDate) {
                    users = users.filter(user => new Date(user.createdAt) >= startDate);
                }
            }

            // Sorting
            users.sort((a, b) => {
                let aVal = a[sortBy];
                let bVal = b[sortBy];
                
                if (sortBy === 'name') {
                    aVal = a.userName || `${a.firstName} ${a.lastName}` || a.email;
                    bVal = b.userName || `${b.firstName} ${b.lastName}` || b.email;
                }
                
                if (aVal < bVal) return sortOrder === 'asc' ? -1 : 1;
                if (aVal > bVal) return sortOrder === 'asc' ? 1 : -1;
                return 0;
            });

            // Calculate pagination
            const startIndex = (currentPage - 1) * itemsPerPage;
            const endIndex = startIndex + itemsPerPage;
            const paginatedUsers = users.slice(startIndex, endIndex);

            return {
                users: paginatedUsers,
                total: users.length,
                page: currentPage,
                totalPages: Math.ceil(users.length / itemsPerPage),
            };
        },
        staleTime: 60000,
    });

    // Mutations
    const toggleUserMutation = useMutation({
        mutationFn: async ({userId, enabled}) => {
            if (enabled) {
                return await enableUser(userId);
            } else {
                return await disableUser(userId);
            }
        },
        onSuccess: () => {
            queryClient.invalidateQueries(['adminUsers']);
        },
    });

    const updateRoleMutation = useMutation({
        mutationFn: ({userId, role}) => updateUser(userId, {role}),
        onSuccess: () => {
            queryClient.invalidateQueries(['adminUsers']);
        },
    });

    const forceLogoutMutation = useMutation({
        mutationFn: (userId) => clearUserTokens({
            userId,
            reason: 'Admin forced logout'
        }),
        onSuccess: () => {
            queryClient.invalidateQueries(['adminUsers']);
        },
    });

    // Handle user selection
    const handleSelectUser = useCallback((userId, selected) => {
        setSelectedUsers(prev => {
            const newSet = new Set(prev);
            if (selected) {
                newSet.add(userId);
            } else {
                newSet.delete(userId);
            }
            return newSet;
        });
    }, []);

    const handleSelectAll = useCallback((checked) => {
        if (checked) {
            setSelectedUsers(new Set(usersData?.users.map(u => u.id) || []));
        } else {
            setSelectedUsers(new Set());
        }
    }, [usersData]);

    // Handle user actions
    const handleUserAction = useCallback((action, user) => {
        switch (action) {
            case 'view':
                router.push(`/${locale}/admin/users/${user.id}`);
                break;
            case 'edit':
                router.push(`/${locale}/admin/users/${user.id}/edit`);
                break;
            case 'toggle':
                toggleUserMutation.mutate({userId: user.id, enabled: !user.enabled});
                break;
            case 'changeRole':
                const newRole = prompt('Enter new role (customer, business, admin):', user.role);
                if (newRole && ['customer', 'business', 'admin'].includes(newRole)) {
                    updateRoleMutation.mutate({userId: user.id, role: newRole});
                }
                break;
            case 'forceLogout':
                if (confirm('Force logout this user from all sessions?')) {
                    forceLogoutMutation.mutate(user.id);
                }
                break;
            case 'viewMetrics':
                router.push(`/${locale}/admin/users/${user.id}/metrics`);
                break;
        }
    }, [router, toggleUserMutation, updateRoleMutation, forceLogoutMutation, locale]);

    // Handle bulk actions
    const handleBulkAction = useCallback((action) => {
        const userIds = Array.from(selectedUsers);
        if (userIds.length === 0) return;

        switch (action) {
            case 'activate':
                if (confirm(`Activate ${userIds.length} users?`)) {
                    userIds.forEach(id => {
                        const user = usersData?.users.find(u => u.id === id);
                        if (user && !user.enabled) {
                            toggleUserMutation.mutate({userId: id, enabled: true});
                        }
                    });
                    setSelectedUsers(new Set());
                }
                break;
            case 'deactivate':
                if (confirm(`Deactivate ${userIds.length} users?`)) {
                    userIds.forEach(id => {
                        const user = usersData?.users.find(u => u.id === id);
                        if (user && user.enabled) {
                            toggleUserMutation.mutate({userId: id, enabled: false});
                        }
                    });
                    setSelectedUsers(new Set());
                }
                break;
            case 'export':
                handleExportSelected();
                break;
        }
    }, [selectedUsers, usersData, toggleUserMutation]);

    // Export selected users
    const handleExportSelected = useCallback(() => {
        const selectedUserData = usersData?.users.filter(u => selectedUsers.has(u.id)) || [];
        if (selectedUserData.length === 0) return;

        const headers = ['ID', 'Email', 'Username', 'First Name', 'Last Name', 'Role', 'Status', 'Created At'];
        const rows = selectedUserData.map(user => [
            user.id,
            user.email,
            user.userName || '',
            user.firstName || '',
            user.lastName || '',
            (user.role || 'customer').charAt(0).toUpperCase() + (user.role || 'customer').slice(1),
            user.enabled ? 'Active' : 'Inactive',
            new Date(user.createdAt).toLocaleDateString()
        ]);
        
        const csvContent = [
            headers.join(','),
            ...rows.map(row => row.map(cell => `"${String(cell).replace(/"/g, '""')}"`).join(','))
        ].join('\n');
        
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `selected-users-${new Date().toISOString().split('T')[0]}.csv`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);
    }, [selectedUsers, usersData]);

    // Handle sorting
    const handleSort = useCallback((field) => {
        if (sortBy === field) {
            setSortOrder(prev => prev === 'asc' ? 'desc' : 'asc');
        } else {
            setSortBy(field);
            setSortOrder('asc');
        }
    }, [sortBy]);

    // Export users - client-side CSV generation since API doesn't support it
    const handleExport = useCallback(async () => {
        try {
            // Get all users (not just current page)
            const response = await listUsers();
            let allUsers = response.users || [];
            
            // Apply filters
            if (searchTerm) {
                allUsers = allUsers.filter(user =>
                    user.userName?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                    user.email?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                    `${user.firstName} ${user.lastName}`.toLowerCase().includes(searchTerm.toLowerCase())
                );
            }
            
            if (filterRole !== 'all') {
                allUsers = allUsers.filter(user => (user.role || 'customer').toLowerCase() === filterRole);
            }
            
            if (filterStatus !== 'all') {
                allUsers = allUsers.filter(user =>
                    filterStatus === 'active' ? user.enabled : !user.enabled
                );
            }
            
            // Generate CSV
            const headers = ['ID', 'Email', 'Username', 'First Name', 'Last Name', 'Role', 'Status', 'Created At'];
            const rows = allUsers.map(user => [
                user.id,
                user.email,
                user.userName || '',
                user.firstName || '',
                user.lastName || '',
                (user.role || 'customer').charAt(0).toUpperCase() + (user.role || 'customer').slice(1),
                user.enabled ? 'Active' : 'Inactive',
                new Date(user.createdAt).toLocaleDateString()
            ]);
            
            const csvContent = [
                headers.join(','),
                ...rows.map(row => row.map(cell => `"${String(cell).replace(/"/g, '""')}"`).join(','))
            ].join('\n');
            
            const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.download = `users-export-${new Date().toISOString().split('T')[0]}.csv`;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
            window.URL.revokeObjectURL(url);
        } catch (error) {
            // Error: 'Export failed:', error...
            alert('Failed to export users. Please try again.');
        }
    }, [filterRole, filterStatus, searchTerm]);

    const stats = useMemo(() => {
        if (!usersData) return {total: 0, active: 0, inactive: 0, newThisWeek: 0};

        const oneWeekAgo = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000);
        return {
            total: usersData.total || 0,
            active: usersData.users?.filter(u => u.enabled).length || 0,
            inactive: usersData.users?.filter(u => !u.enabled).length || 0,
            newThisWeek: usersData.users?.filter(u => new Date(u.createdAt) >= oneWeekAgo).length || 0,
        };
    }, [usersData]);

    if (!isAdmin) {
        return null;
    }

    if (isLoading) {
        return (
            <div className={styles.loadingContainer}>
                <LoadingSpinner/>
            </div>
        );
    }

    return (
        <ErrorBoundary>
            <div className={styles.container}>
                {/* Header */}
                <div className={styles.header}>
                    <div className={styles.headerLeft}>
                        <h1 className={styles.title}>{t('title', {defaultValue: 'User Management'})}</h1>
                        <p className={styles.subtitle}>
                            {t('subtitle', {defaultValue: 'Manage user accounts and permissions'})}
                        </p>
                    </div>
                    <div className={styles.headerActions}>
                        <button
                            className={styles.secondaryButton}
                            onClick={() => setShowAddUser(true)}
                        >
                            <UserPlus size={18}/>
                            <span>Add User</span>
                        </button>
                        <button
                            className={styles.secondaryButton}
                            onClick={handleExport}
                        >
                            <Download size={18}/>
                            <span>{t('export', {defaultValue: 'Export'})}</span>
                        </button>
                        <button
                            className={styles.iconButton}
                            onClick={() => refetch()}
                            title="Refresh"
                        >
                            <RefreshCw size={18}/>
                        </button>
                        <button
                            className={styles.iconButton}
                            onClick={() => setShowFilters(!showFilters)}
                            title="Settings"
                        >
                            <Settings size={18}/>
                        </button>
                    </div>
                </div>

                {/* Stats */}
                <div className={styles.statsRow}>
                    <div className={styles.statCard}>
                        <div className={styles.statIcon}>
                            <Users size={20} />
                        </div>
                        <div className={styles.statContent}>
                            <span className={styles.statValue}>{stats.total}</span>
                            <span className={styles.statLabel}>{t('totalUsers', {defaultValue: 'Total Users'})}</span>
                        </div>
                    </div>
                    <div className={styles.statCard}>
                        <div className={`${styles.statIcon} ${styles.statIconActive}`}>
                            <CheckCircle size={20} />
                        </div>
                        <div className={styles.statContent}>
                            <span className={styles.statValue}>{stats.active}</span>
                            <span className={styles.statLabel}>{t('activeUsers', {defaultValue: 'Active Users'})}</span>
                        </div>
                    </div>
                    <div className={styles.statCard}>
                        <div className={`${styles.statIcon} ${styles.statIconInactive}`}>
                            <XCircle size={20} />
                        </div>
                        <div className={styles.statContent}>
                            <span className={styles.statValue}>{stats.inactive}</span>
                            <span className={styles.statLabel}>{t('inactiveUsers', {defaultValue: 'Inactive Users'})}</span>
                        </div>
                    </div>
                    <div className={styles.statCard}>
                        <div className={`${styles.statIcon} ${styles.statIconNew}`}>
                            <TrendingUp size={20} />
                        </div>
                        <div className={styles.statContent}>
                            <span className={styles.statValue}>{stats.newThisWeek}</span>
                            <span className={styles.statLabel}>New This Week</span>
                        </div>
                    </div>
                </div>

                {/* Search and Filters */}
                <div className={styles.controls}>
                    <div className={styles.controlsLeft}>
                        <div className={styles.searchBox}>
                            <Search size={18}/>
                            <input
                                type="text"
                                placeholder={t('searchPlaceholder', {defaultValue: 'Search by name or email...'})}
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                                className={styles.searchInput}
                            />
                        </div>
                        <div className={styles.quickFilters}>
                            <button
                                className={`${styles.quickFilterButton} ${filterStatus === 'all' ? styles.quickFilterActive : ''}`}
                                onClick={() => setFilterStatus('all')}
                            >
                                All
                            </button>
                            <button
                                className={`${styles.quickFilterButton} ${filterStatus === 'active' ? styles.quickFilterActive : ''}`}
                                onClick={() => setFilterStatus('active')}
                            >
                                Active
                            </button>
                            <button
                                className={`${styles.quickFilterButton} ${filterStatus === 'inactive' ? styles.quickFilterActive : ''}`}
                                onClick={() => setFilterStatus('inactive')}
                            >
                                Inactive
                            </button>
                        </div>
                    </div>
                    <button
                        className={`${styles.filterButton} ${showFilters ? styles.filterActive : ''}`}
                        onClick={() => setShowFilters(!showFilters)}
                    >
                        <SlidersHorizontal size={16}/>
                        {t('filters', {defaultValue: 'Advanced Filters'})}
                        {(filterRole !== 'all' || filterDateRange !== 'all') && (
                            <span className={styles.filterBadge}>
                                {[filterRole !== 'all', filterDateRange !== 'all'].filter(Boolean).length}
                            </span>
                        )}
                    </button>
                </div>

                {/* Filter Panel */}
                {showFilters && (
                    <div className={styles.filterPanel}>
                        <div className={styles.filterPanelContent}>
                            <div className={styles.filterGroup}>
                                <label className={styles.filterLabel}>{t('filterByRole', {defaultValue: 'Role'})}</label>
                                <select
                                    value={filterRole}
                                    onChange={(e) => setFilterRole(e.target.value)}
                                    className={styles.filterSelect}
                                >
                                    <option value="all">{t('allRoles', {defaultValue: 'All Roles'})}</option>
                                    <option value="customer">{t('customer', {defaultValue: 'Customer'})}</option>
                                    <option value="business">{t('business', {defaultValue: 'Business'})}</option>
                                    <option value="admin">{t('admin', {defaultValue: 'Admin'})}</option>
                                </select>
                            </div>
                            <div className={styles.filterGroup}>
                                <label className={styles.filterLabel}>Date Range</label>
                                <select
                                    value={filterDateRange}
                                    onChange={(e) => setFilterDateRange(e.target.value)}
                                    className={styles.filterSelect}
                                >
                                    <option value="all">All Time</option>
                                    <option value="7days">Last 7 Days</option>
                                    <option value="30days">Last 30 Days</option>
                                    <option value="90days">Last 90 Days</option>
                                </select>
                            </div>
                            <div className={styles.filterGroup}>
                                <label className={styles.filterLabel}>Sort By</label>
                                <select
                                    value={sortBy}
                                    onChange={(e) => setSortBy(e.target.value)}
                                    className={styles.filterSelect}
                                >
                                    <option value="createdAt">Join Date</option>
                                    <option value="name">Name</option>
                                    <option value="email">Email</option>
                                    <option value="lastActive">Last Active</option>
                                </select>
                            </div>
                            <div className={styles.filterGroup}>
                                <label className={styles.filterLabel}>Order</label>
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
                        <button
                            className={styles.clearFiltersButton}
                            onClick={() => {
                                setFilterRole('all');
                                setFilterStatus('all');
                                setFilterDateRange('all');
                                setSortBy('createdAt');
                                setSortOrder('desc');
                                setSearchTerm('');
                            }}
                        >
                            Clear All Filters
                        </button>
                    </div>
                )}

                {/* Bulk Actions Bar */}
                {selectedUsers.size > 0 && (
                    <div className={styles.bulkActionsBar}>
                        <div className={styles.bulkActionsLeft}>
                            <span className={styles.selectedCount}>
                                {selectedUsers.size} user{selectedUsers.size !== 1 ? 's' : ''} selected
                            </span>
                            <button
                                className={styles.clearSelectionButton}
                                onClick={() => setSelectedUsers(new Set())}
                            >
                                Clear selection
                            </button>
                        </div>
                        <div className={styles.bulkActionsRight}>
                            <button
                                className={styles.bulkActionButton}
                                onClick={() => handleBulkAction('activate')}
                            >
                                <CheckCircle size={16} />
                                Activate
                            </button>
                            <button
                                className={styles.bulkActionButton}
                                onClick={() => handleBulkAction('deactivate')}
                            >
                                <XCircle size={16} />
                                Deactivate
                            </button>
                            <button
                                className={styles.bulkActionButton}
                                onClick={() => handleBulkAction('export')}
                            >
                                <Download size={16} />
                                Export Selected
                            </button>
                        </div>
                    </div>
                )}

                {/* Users Table */}
                <div className={styles.tableContainer}>
                    <table className={styles.usersTable}>
                        <thead>
                        <tr>
                            <th className={styles.checkboxCell}>
                                <input
                                    type="checkbox"
                                    checked={selectedUsers.size === usersData?.users.length && usersData?.users.length > 0}
                                    onChange={(e) => handleSelectAll(e.target.checked)}
                                    className={styles.checkbox}
                                />
                            </th>
                            <th className={styles.sortableHeader} onClick={() => handleSort('name')}>
                                <span>{t('user', {defaultValue: 'User'})}</span>
                                {sortBy === 'name' && (
                                    <ArrowUpDown size={14} className={sortOrder === 'asc' ? styles.sortAsc : styles.sortDesc} />
                                )}
                            </th>
                            <th>{t('role', {defaultValue: 'Role'})}</th>
                            <th>{t('status', {defaultValue: 'Status'})}</th>
                            <th className={styles.sortableHeader} onClick={() => handleSort('createdAt')}>
                                <span>{t('joined', {defaultValue: 'Joined'})}</span>
                                {sortBy === 'createdAt' && (
                                    <ArrowUpDown size={14} className={sortOrder === 'asc' ? styles.sortAsc : styles.sortDesc} />
                                )}
                            </th>
                            <th className={styles.sortableHeader} onClick={() => handleSort('lastActive')}>
                                <span>{t('lastActive', {defaultValue: 'Last Active'})}</span>
                                {sortBy === 'lastActive' && (
                                    <ArrowUpDown size={14} className={sortOrder === 'asc' ? styles.sortAsc : styles.sortDesc} />
                                )}
                            </th>
                            <th>{t('actions', {defaultValue: 'Actions'})}</th>
                        </tr>
                        </thead>
                        <tbody>
                        {usersData?.users.length === 0 ? (
                            <tr>
                                <td colSpan="8" className={styles.emptyState}>
                                    <div className={styles.emptyStateContent}>
                                        <Users size={48} />
                                        <h3>No users found</h3>
                                        <p>Try adjusting your filters or search term</p>
                                    </div>
                                </td>
                            </tr>
                        ) : (
                            usersData?.users.map(user => (
                                <UserRow
                                    key={user.id}
                                    user={user}
                                    onAction={handleUserAction}
                                    selected={selectedUsers.has(user.id)}
                                    onSelect={handleSelectUser}
                                />
                            ))
                        )}
                        </tbody>
                    </table>
                </div>

                {/* Pagination */}
                {usersData && usersData.totalPages > 1 && (
                    <div className={styles.pagination}>
                        <div className={styles.paginationLeft}>
                            <span className={styles.paginationInfo}>
                                Showing {((currentPage - 1) * itemsPerPage) + 1} to {Math.min(currentPage * itemsPerPage, usersData.total)} of {usersData.total} users
                            </span>
                        </div>
                        <div className={styles.paginationCenter}>
                            <button
                                className={styles.paginationButton}
                                onClick={() => setCurrentPage(1)}
                                disabled={currentPage === 1}
                                title="First page"
                            >
                                <ChevronsLeft size={16}/>
                            </button>
                            <button
                                className={styles.paginationButton}
                                onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                                disabled={currentPage === 1}
                                title="Previous page"
                            >
                                <ChevronLeft size={16}/>
                            </button>
                            <div className={styles.pageNumbers}>
                                {[...Array(Math.min(5, usersData.totalPages))].map((_, idx) => {
                                    let pageNum;
                                    if (usersData.totalPages <= 5) {
                                        pageNum = idx + 1;
                                    } else if (currentPage <= 3) {
                                        pageNum = idx + 1;
                                    } else if (currentPage >= usersData.totalPages - 2) {
                                        pageNum = usersData.totalPages - 4 + idx;
                                    } else {
                                        pageNum = currentPage - 2 + idx;
                                    }
                                    
                                    if (pageNum > 0 && pageNum <= usersData.totalPages) {
                                        return (
                                            <button
                                                key={pageNum}
                                                className={`${styles.pageButton} ${currentPage === pageNum ? styles.pageButtonActive : ''}`}
                                                onClick={() => setCurrentPage(pageNum)}
                                            >
                                                {pageNum}
                                            </button>
                                        );
                                    }
                                    return null;
                                })}
                            </div>
                            <button
                                className={styles.paginationButton}
                                onClick={() => setCurrentPage(prev => Math.min(usersData.totalPages, prev + 1))}
                                disabled={currentPage === usersData.totalPages}
                                title="Next page"
                            >
                                <ChevronRight size={16}/>
                            </button>
                            <button
                                className={styles.paginationButton}
                                onClick={() => setCurrentPage(usersData.totalPages)}
                                disabled={currentPage === usersData.totalPages}
                                title="Last page"
                            >
                                <ChevronsRight size={16}/>
                            </button>
                        </div>
                        <div className={styles.paginationRight}>
                            <select
                                className={styles.itemsPerPageSelect}
                                value={itemsPerPage}
                                onChange={(e) => {
                                    // This would need implementation
                                    
                                }}
                            >
                                <option value="10">10 / page</option>
                                <option value="25">25 / page</option>
                                <option value="50">50 / page</option>
                                <option value="100">100 / page</option>
                            </select>
                        </div>
                    </div>
                )}
            </div>
        </ErrorBoundary>
    );
};

export default UserManagement;