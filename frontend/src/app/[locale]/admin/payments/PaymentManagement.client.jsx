"use client";

import React, { useState, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  CreditCard,
  DollarSign,
  RotateCcw,
  TrendingUp,
  TrendingDown,
  Search,
  Filter,
  Download,
  RefreshCw,
  Eye,
  AlertCircle,
  CheckCircle,
  XCircle,
  Clock,
  Calendar,
  ArrowUpRight,
  ArrowDownRight,
  MoreVertical,
  Receipt,
  Ban,
  FileText,
  User
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import {
  getPaymentTransactions,
  refundPayment,
  getPaymentAnalytics,
  adjustInvoice,
  createInvoice,
  cancelInvoice,
  payInvoice
} from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './PaymentManagement.module.css';

const PaymentStatusBadge = ({ status }) => {
  const statusConfig = {
    completed: { color: '#22c55e', bg: 'rgba(34, 197, 94, 0.1)', text: 'Completed', icon: CheckCircle },
    pending: { color: '#f59e0b', bg: 'rgba(245, 158, 11, 0.1)', text: 'Pending', icon: Clock },
    failed: { color: '#ef4444', bg: 'rgba(239, 68, 68, 0.1)', text: 'Failed', icon: XCircle },
    refunded: { color: '#6366f1', bg: 'rgba(99, 102, 241, 0.1)', text: 'Refunded', icon: Refund },
    cancelled: { color: '#64748b', bg: 'rgba(100, 116, 139, 0.1)', text: 'Cancelled', icon: Ban }
  };

  const config = statusConfig[status] || statusConfig.pending;
  const Icon = config.icon;

  return (
    <span 
      className={styles.statusBadge}
      style={{ color: config.color, backgroundColor: config.bg }}
    >
      <Icon size={12} />
      {config.text}
    </span>
  );
};

const PaymentTypeIcon = ({ type }) => {
  const icons = {
    card: CreditCard,
    bank: Receipt,
    paypal: DollarSign,
    crypto: TrendingUp
  };
  
  const Icon = icons[type] || CreditCard;
  return <Icon size={16} />;
};

const TransactionRow = ({ transaction, onAction }) => {
  const [showMenu, setShowMenu] = useState(false);

  return (
    <tr className={styles.transactionRow}>
      <td className={styles.transactionCell}>
        <div className={styles.transactionInfo}>
          <div className={styles.transactionId}>
            <PaymentTypeIcon type={transaction.paymentMethod} />
            <span>{transaction.id.slice(-8)}</span>
          </div>
          <div className={styles.transactionMeta}>
            {new Date(transaction.createdAt).toLocaleDateString()}
          </div>
        </div>
      </td>
      <td className={styles.userCell}>
        <div className={styles.userInfo}>
          <User size={16} />
          <span>{transaction.user?.name || 'Unknown User'}</span>
        </div>
      </td>
      <td className={styles.amountCell}>
        <div className={styles.amount}>
          <span className={styles.currency}>{transaction.currency}</span>
          <span className={styles.value}>{transaction.amount.toLocaleString()}</span>
        </div>
      </td>
      <td className={styles.statusCell}>
        <PaymentStatusBadge status={transaction.status} />
      </td>
      <td className={styles.methodCell}>
        <span className={styles.paymentMethod}>
          {transaction.paymentMethod?.toUpperCase()}
        </span>
      </td>
      <td className={styles.actionCell}>
        <div className={styles.actionMenu}>
          <button
            className={styles.menuTrigger}
            onClick={() => setShowMenu(!showMenu)}
          >
            <MoreVertical size={16} />
          </button>
          {showMenu && (
            <div className={styles.actionDropdown}>
              <button onClick={() => onAction('view', transaction)}>
                <Eye size={14} />
                View Details
              </button>
              {transaction.status === 'completed' && (
                <button onClick={() => onAction('refund', transaction)}>
                  <Refund size={14} />
                  Process Refund
                </button>
              )}
              <button onClick={() => onAction('invoice', transaction)}>
                <FileText size={14} />
                View Invoice
              </button>
            </div>
          )}
        </div>
      </td>
    </tr>
  );
};

const PaymentManagement = () => {
  const t = useTranslations('PaymentManagement');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const queryClient = useQueryClient();

  const [filters, setFilters] = useState({
    status: 'all',
    paymentMethod: 'all',
    dateRange: '30d',
    search: ''
  });
  const [selectedTransactions, setSelectedTransactions] = useState([]);
  const [showRefundModal, setShowRefundModal] = useState(false);
  const [refundTransaction, setRefundTransaction] = useState(null);

  // Fetch payment transactions
  const { 
    data: transactionsData, 
    isLoading: transactionsLoading, 
    error: transactionsError,
    refetch: refetchTransactions 
  } = useQuery({
    queryKey: ['adminPayments', filters],
    queryFn: () => getPaymentTransactions(filters),
    enabled: isAdmin
  });

  // Fetch payment analytics
  const { 
    data: analyticsData, 
    isLoading: analyticsLoading 
  } = useQuery({
    queryKey: ['paymentAnalytics', filters.dateRange],
    queryFn: () => getPaymentAnalytics({ dateRange: filters.dateRange }),
    enabled: isAdmin
  });

  // Refund mutation
  const refundMutation = useMutation({
    mutationFn: ({ paymentId, refundData }) => refundPayment(paymentId, refundData),
    onSuccess: () => {
      queryClient.invalidateQueries(['adminPayments']);
      queryClient.invalidateQueries(['paymentAnalytics']);
      setShowRefundModal(false);
      setRefundTransaction(null);
    },
    onError: (error) => {
      
      alert('Failed to process refund. Please try again.');
    }
  });

  const handleTransactionAction = useCallback((action, transaction) => {
    switch (action) {
      case 'view':
        router.push(`/admin/payments/${transaction.id}`);
        break;
      case 'refund':
        setRefundTransaction(transaction);
        setShowRefundModal(true);
        break;
      case 'invoice':
        window.open(`/admin/payments/invoices/${transaction.invoiceId}`, '_blank');
        break;
    }
  }, [router]);

  const handleRefund = useCallback((refundData) => {
    if (refundTransaction) {
      refundMutation.mutate({
        paymentId: refundTransaction.id,
        refundData
      });
    }
  }, [refundTransaction, refundMutation]);

  const handleFilterChange = useCallback((key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  }, []);

  const handleExport = useCallback(() => {
    // Export transactions to CSV
    const csvData = transactions.map(t => ({
      'Transaction ID': t.id,
      'User': t.user?.name || 'Unknown',
      'Amount': t.amount,
      'Currency': t.currency,
      'Status': t.status,
      'Payment Method': t.paymentMethod,
      'Date': new Date(t.createdAt).toLocaleDateString()
    }));
    
    // Convert to CSV and download
    
  }, [transactionsData]);

  // Process data
  const transactions = transactionsData?.transactions || [];
  const analytics = analyticsData || {};

  // Calculate summary stats
  const stats = useMemo(() => {
    const completed = transactions.filter(t => t.status === 'completed');
    const totalRevenue = completed.reduce((sum, t) => sum + t.amount, 0);
    const refunded = transactions.filter(t => t.status === 'refunded');
    const totalRefunded = refunded.reduce((sum, t) => sum + t.amount, 0);

    return {
      totalTransactions: transactions.length,
      totalRevenue,
      totalRefunded,
      successRate: transactions.length > 0 ? (completed.length / transactions.length) * 100 : 0,
      averageTransaction: completed.length > 0 ? totalRevenue / completed.length : 0
    };
  }, [transactions]);

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access payment management.' })}</p>
          <p>{t('currentRole', { defaultValue: 'Current role' })}: {user?.role || t('notLoggedIn', { defaultValue: 'Not logged in' })}</p>
        </div>
      </div>
    );
  }

  if (transactionsLoading && !transactionsData) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <LoadingSpinner />
          <h2>{t('loadingTitle', { defaultValue: 'Loading Payment Data...' })}</h2>
          <p>{t('loadingMessage', { defaultValue: 'Please wait while we fetch transaction information.' })}</p>
        </div>
      </div>
    );
  }

  if (transactionsError) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <AlertCircle size={48} />
          <h2>{t('errorTitle', { defaultValue: 'Failed to Load Payments' })}</h2>
          <p>{transactionsError.message || t('errorMessage', { defaultValue: 'An error occurred while fetching payment data' })}</p>
          <button className={styles.retryButton} onClick={() => refetchTransactions()}>
            <RefreshCw size={16} />
            {t('retry', { defaultValue: 'Try Again' })}
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
          <div className={styles.headerContent}>
            <h1 className={styles.title}>{t('title', { defaultValue: 'Payment Management' })}</h1>
            <p className={styles.subtitle}>{t('subtitle', { defaultValue: 'Monitor and manage payment transactions' })}</p>
          </div>
          <div className={styles.headerActions}>
            <button className={styles.exportButton} onClick={handleExport}>
              <Download size={16} />
              {t('export', { defaultValue: 'Export' })}
            </button>
            <button className={styles.refreshButton} onClick={() => refetchTransactions()}>
              <RefreshCw size={16} />
              {t('refresh', { defaultValue: 'Refresh' })}
            </button>
          </div>
        </div>

        {/* Stats Overview */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <Receipt size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{stats.totalTransactions.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('totalTransactions', { defaultValue: 'Total Transactions' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <DollarSign size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>${stats.totalRevenue.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('totalRevenue', { defaultValue: 'Total Revenue' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <Refund size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>${stats.totalRefunded.toLocaleString()}</div>
              <div className={styles.statLabel}>{t('totalRefunded', { defaultValue: 'Total Refunded' })}</div>
            </div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <TrendingUp size={20} />
            </div>
            <div className={styles.statContent}>
              <div className={styles.statValue}>{stats.successRate.toFixed(1)}%</div>
              <div className={styles.statLabel}>{t('successRate', { defaultValue: 'Success Rate' })}</div>
            </div>
          </div>
        </div>

        {/* Filters */}
        <div className={styles.filtersSection}>
          <div className={styles.searchContainer}>
            <Search size={16} />
            <input
              type="text"
              placeholder={t('searchPlaceholder', { defaultValue: 'Search transactions...' })}
              value={filters.search}
              onChange={(e) => handleFilterChange('search', e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <div className={styles.filterControls}>
            <select
              value={filters.status}
              onChange={(e) => handleFilterChange('status', e.target.value)}
              className={styles.filterSelect}
            >
              <option value="all">{t('allStatuses', { defaultValue: 'All Statuses' })}</option>
              <option value="completed">{t('completed', { defaultValue: 'Completed' })}</option>
              <option value="pending">{t('pending', { defaultValue: 'Pending' })}</option>
              <option value="failed">{t('failed', { defaultValue: 'Failed' })}</option>
              <option value="refunded">{t('refunded', { defaultValue: 'Refunded' })}</option>
            </select>
            <select
              value={filters.paymentMethod}
              onChange={(e) => handleFilterChange('paymentMethod', e.target.value)}
              className={styles.filterSelect}
            >
              <option value="all">{t('allMethods', { defaultValue: 'All Payment Methods' })}</option>
              <option value="card">{t('card', { defaultValue: 'Credit Card' })}</option>
              <option value="bank">{t('bank', { defaultValue: 'Bank Transfer' })}</option>
              <option value="paypal">{t('paypal', { defaultValue: 'PayPal' })}</option>
            </select>
            <select
              value={filters.dateRange}
              onChange={(e) => handleFilterChange('dateRange', e.target.value)}
              className={styles.filterSelect}
            >
              <option value="7d">{t('last7Days', { defaultValue: 'Last 7 Days' })}</option>
              <option value="30d">{t('last30Days', { defaultValue: 'Last 30 Days' })}</option>
              <option value="90d">{t('last90Days', { defaultValue: 'Last 90 Days' })}</option>
              <option value="1y">{t('lastYear', { defaultValue: 'Last Year' })}</option>
            </select>
          </div>
        </div>

        {/* Transactions Table */}
        <div className={styles.tableSection}>
          <div className={styles.tableContainer}>
            <table className={styles.transactionsTable}>
              <thead>
                <tr>
                  <th>{t('transaction', { defaultValue: 'Transaction' })}</th>
                  <th>{t('user', { defaultValue: 'User' })}</th>
                  <th>{t('amount', { defaultValue: 'Amount' })}</th>
                  <th>{t('status', { defaultValue: 'Status' })}</th>
                  <th>{t('method', { defaultValue: 'Method' })}</th>
                  <th>{t('actions', { defaultValue: 'Actions' })}</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((transaction) => (
                  <TransactionRow
                    key={transaction.id}
                    transaction={transaction}
                    onAction={handleTransactionAction}
                  />
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Refund Modal */}
        {showRefundModal && refundTransaction && (
          <div className={styles.modal}>
            <div className={styles.modalContent}>
              <h3>{t('processRefund', { defaultValue: 'Process Refund' })}</h3>
              <p>{t('refundConfirmation', { 
                defaultValue: 'Are you sure you want to refund this transaction?',
                values: { amount: `$${refundTransaction.amount}` }
              })}</p>
              <div className={styles.modalActions}>
                <button 
                  className={styles.cancelButton}
                  onClick={() => setShowRefundModal(false)}
                >
                  {t('cancel', { defaultValue: 'Cancel' })}
                </button>
                <button 
                  className={styles.confirmButton}
                  onClick={() => handleRefund({ reason: 'Admin refund' })}
                  disabled={refundMutation.isLoading}
                >
                  {refundMutation.isLoading ? t('processing', { defaultValue: 'Processing...' }) : t('confirmRefund', { defaultValue: 'Confirm Refund' })}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </ErrorBoundary>
  );
};

export default PaymentManagement; 