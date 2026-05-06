"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as Yup from 'yup';
import {
    ArrowLeft,
    Save,
    AlertCircle,
    User,
    Mail,
    Phone,
    MapPin,
    Shield,
    Globe,
    Key
} from 'lucide-react';
import { getUserById, updateUser } from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import { toast } from 'react-toastify';
import styles from './UserEdit.module.css';

const UserEdit = ({ locale, userId }) => {
    const t = useTranslations('UserEdit');
    const router = useRouter();
    const queryClient = useQueryClient();
    const [isSubmitting, setIsSubmitting] = useState(false);

    // Validation schema
    const schema = Yup.object({
        email: Yup.string()
            .email('Invalid email address')
            .required('Email is required'),
        userName: Yup.string()
            .min(3, 'Username must be at least 3 characters')
            .max(50, 'Username must be less than 50 characters'),
        firstName: Yup.string()
            .max(50, 'First name must be less than 50 characters'),
        lastName: Yup.string()
            .max(50, 'Last name must be less than 50 characters'),
        phone: Yup.string()
            .matches(/^[\d\s\-\+\(\)]+$/, 'Invalid phone number format')
            .nullable(),
        address: Yup.string()
            .max(200, 'Address must be less than 200 characters'),
        role: Yup.string()
            .oneOf(['customer', 'business', 'admin'], 'Invalid role'),
        language: Yup.string()
            .oneOf(['en', 'de', 'it', 'pl'], 'Invalid language'),
        privacy: Yup.string()
            .oneOf(['public', 'private', 'friends'], 'Invalid privacy setting')
    });

    // Fetch user details
    const { data: user, isLoading, error } = useQuery({
        queryKey: ['adminUser', userId],
        queryFn: () => getUserById(userId),
        staleTime: 60000,
    });

    // Form setup
    const {
        register,
        handleSubmit,
        formState: { errors, isDirty },
        reset
    } = useForm({
        resolver: yupResolver(schema),
        defaultValues: {
            email: user?.email || '',
            userName: user?.userName || '',
            firstName: user?.firstName || '',
            lastName: user?.lastName || '',
            phone: user?.phone || '',
            address: user?.address || '',
            role: user?.role || 'customer',
            language: user?.language || 'en',
            privacy: user?.privacy || 'public'
        }
    });

    // Update mutation
    const updateUserMutation = useMutation({
        mutationFn: (data) => updateUser(userId, data),
        onSuccess: () => {
            queryClient.invalidateQueries(['adminUser', userId]);
            queryClient.invalidateQueries(['adminUsers']);
            toast.success('User updated successfully');
            router.push(`/${locale}/admin/users/${userId}`);
        },
        onError: (error) => {
            toast.error(error.response?.data?.message || 'Failed to update user');
        }
    });

    const onSubmit = async (data) => {
        setIsSubmitting(true);
        try {
            // Only send changed fields
            const changedFields = {};
            Object.keys(data).forEach(key => {
                if (data[key] !== user[key]) {
                    changedFields[key] = data[key];
                }
            });

            if (Object.keys(changedFields).length > 0) {
                await updateUserMutation.mutateAsync(changedFields);
            } else {
                toast.info('No changes to save');
                router.push(`/${locale}/admin/users/${userId}`);
            }
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    } finally {
            setIsSubmitting(false);
        }
    };

    // Reset form when user data loads
    React.useEffect(() => {
        if (user) {
            reset({
                email: user.email || '',
                userName: user.userName || '',
                firstName: user.firstName || '',
                lastName: user.lastName || '',
                phone: user.phone || '',
                address: user.address || '',
                role: user.role || 'customer',
                language: user.language || 'en',
                privacy: user.privacy || 'public'
            });
        }
    }, [user, reset]);

    if (isLoading) {
        return (
            <div className={styles.loadingContainer}>
                <LoadingSpinner />
            </div>
        );
    }

    if (error || !user) {
        return (
            <div className={styles.container}>
                <div className={styles.errorContainer}>
                    <AlertCircle size={48} className={styles.errorIcon} />
                    <h2>User Not Found</h2>
                    <p>The user you're trying to edit doesn't exist.</p>
                    <button 
                        className={styles.primaryButton}
                        onClick={() => router.push(`/${locale}/admin/users`)}
                    >
                        Back to Users
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
                    <button 
                        className={styles.backButton}
                        onClick={() => router.push(`/${locale}/admin/users/${userId}`)}
                    >
                        <ArrowLeft size={20} />
                        Cancel
                    </button>
                    <h1 className={styles.title}>Edit User</h1>
                </div>

                <form onSubmit={handleSubmit(onSubmit)} className={styles.form}>
                    {/* Basic Information */}
                    <div className={styles.section}>
                        <h2 className={styles.sectionTitle}>
                            <User size={20} />
                            Basic Information
                        </h2>
                        
                        <div className={styles.fieldGrid}>
                            <div className={styles.field}>
                                <label htmlFor="email">Email *</label>
                                <div className={styles.inputWrapper}>
                                    <Mail size={18} className={styles.inputIcon} />
                                    <input
                                        id="email"
                                        type="email"
                                        {...register('email')}
                                        className={errors.email ? styles.inputError : ''}
                                    />
                                </div>
                                {errors.email && (
                                    <span className={styles.errorMessage}>{errors.email.message}</span>
                                )}
                            </div>

                            <div className={styles.field}>
                                <label htmlFor="userName">Username</label>
                                <div className={styles.inputWrapper}>
                                    <User size={18} className={styles.inputIcon} />
                                    <input
                                        id="userName"
                                        type="text"
                                        {...register('userName')}
                                        className={errors.userName ? styles.inputError : ''}
                                    />
                                </div>
                                {errors.userName && (
                                    <span className={styles.errorMessage}>{errors.userName.message}</span>
                                )}
                            </div>

                            <div className={styles.field}>
                                <label htmlFor="firstName">First Name</label>
                                <input
                                    id="firstName"
                                    type="text"
                                    {...register('firstName')}
                                    className={errors.firstName ? styles.inputError : ''}
                                />
                                {errors.firstName && (
                                    <span className={styles.errorMessage}>{errors.firstName.message}</span>
                                )}
                            </div>

                            <div className={styles.field}>
                                <label htmlFor="lastName">Last Name</label>
                                <input
                                    id="lastName"
                                    type="text"
                                    {...register('lastName')}
                                    className={errors.lastName ? styles.inputError : ''}
                                />
                                {errors.lastName && (
                                    <span className={styles.errorMessage}>{errors.lastName.message}</span>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Contact Information */}
                    <div className={styles.section}>
                        <h2 className={styles.sectionTitle}>
                            <Phone size={20} />
                            Contact Information
                        </h2>
                        
                        <div className={styles.fieldGrid}>
                            <div className={styles.field}>
                                <label htmlFor="phone">Phone</label>
                                <div className={styles.inputWrapper}>
                                    <Phone size={18} className={styles.inputIcon} />
                                    <input
                                        id="phone"
                                        type="tel"
                                        {...register('phone')}
                                        className={errors.phone ? styles.inputError : ''}
                                    />
                                </div>
                                {errors.phone && (
                                    <span className={styles.errorMessage}>{errors.phone.message}</span>
                                )}
                            </div>

                            <div className={styles.field}>
                                <label htmlFor="address">Address</label>
                                <div className={styles.inputWrapper}>
                                    <MapPin size={18} className={styles.inputIcon} />
                                    <input
                                        id="address"
                                        type="text"
                                        {...register('address')}
                                        className={errors.address ? styles.inputError : ''}
                                    />
                                </div>
                                {errors.address && (
                                    <span className={styles.errorMessage}>{errors.address.message}</span>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Account Settings */}
                    <div className={styles.section}>
                        <h2 className={styles.sectionTitle}>
                            <Shield size={20} />
                            Account Settings
                        </h2>
                        
                        <div className={styles.fieldGrid}>
                            <div className={styles.field}>
                                <label htmlFor="role">Role</label>
                                <select
                                    id="role"
                                    {...register('role')}
                                    className={errors.role ? styles.inputError : ''}
                                >
                                    <option value="customer">Customer</option>
                                    <option value="business">Business</option>
                                    <option value="admin">Admin</option>
                                </select>
                                {errors.role && (
                                    <span className={styles.errorMessage}>{errors.role.message}</span>
                                )}
                            </div>

                            <div className={styles.field}>
                                <label htmlFor="language">Language</label>
                                <div className={styles.inputWrapper}>
                                    <Globe size={18} className={styles.inputIcon} />
                                    <select
                                        id="language"
                                        {...register('language')}
                                        className={errors.language ? styles.inputError : ''}
                                    >
                                        <option value="en">English</option>
                                        <option value="de">German</option>
                                        <option value="it">Italian</option>
                                        <option value="pl">Polish</option>
                                    </select>
                                </div>
                                {errors.language && (
                                    <span className={styles.errorMessage}>{errors.language.message}</span>
                                )}
                            </div>

                            <div className={styles.field}>
                                <label htmlFor="privacy">Privacy Setting</label>
                                <div className={styles.inputWrapper}>
                                    <Key size={18} className={styles.inputIcon} />
                                    <select
                                        id="privacy"
                                        {...register('privacy')}
                                        className={errors.privacy ? styles.inputError : ''}
                                    >
                                        <option value="public">Public</option>
                                        <option value="private">Private</option>
                                        <option value="friends">Friends Only</option>
                                    </select>
                                </div>
                                {errors.privacy && (
                                    <span className={styles.errorMessage}>{errors.privacy.message}</span>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Read-only Information */}
                    <div className={styles.section}>
                        <h2 className={styles.sectionTitle}>
                            <AlertCircle size={20} />
                            Read-only Information
                        </h2>
                        
                        <div className={styles.readOnlyGrid}>
                            <div className={styles.readOnlyItem}>
                                <span className={styles.readOnlyLabel}>User ID</span>
                                <span className={styles.readOnlyValue}>{user.id}</span>
                            </div>
                            <div className={styles.readOnlyItem}>
                                <span className={styles.readOnlyLabel}>Created</span>
                                <span className={styles.readOnlyValue}>
                                    {new Date(user.createdAt).toLocaleDateString()}
                                </span>
                            </div>
                            <div className={styles.readOnlyItem}>
                                <span className={styles.readOnlyLabel}>Email Verified</span>
                                <span className={styles.readOnlyValue}>
                                    {user.emailVerified ? 'Yes' : 'No'}
                                </span>
                            </div>
                            <div className={styles.readOnlyItem}>
                                <span className={styles.readOnlyLabel}>Account Status</span>
                                <span className={styles.readOnlyValue}>
                                    {user.enabled ? 'Active' : 'Inactive'}
                                </span>
                            </div>
                        </div>
                    </div>

                    {/* Actions */}
                    <div className={styles.actions}>
                        <button
                            type="button"
                            className={styles.cancelButton}
                            onClick={() => router.push(`/${locale}/admin/users/${userId}`)}
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className={styles.saveButton}
                            disabled={!isDirty || isSubmitting}
                        >
                            <Save size={16} />
                            {isSubmitting ? 'Saving...' : 'Save Changes'}
                        </button>
                    </div>
                </form>
            </div>
        </ErrorBoundary>
    );
};

export default UserEdit;