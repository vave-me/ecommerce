"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useMutation } from '@tanstack/react-query';
import {
  ArrowLeft,
  User,
  Mail,
  Lock,
  Shield,
  Eye,
  EyeOff,
  Check,
  AlertCircle,
  Save,
  RefreshCw,
  Users,
  Settings,
  Crown,
  UserPlus,
  Calendar,
  Phone,
  MapPin,
  Globe
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from './AddAdminManagement.module.css';

// Admin roles configuration
const adminRoles = [
  {
    id: 'admin',
    label: 'Administrator',
    description: 'Full access to all admin features',
    icon: Crown,
    color: '#3b82f6',
    permissions: [
      'Manage users',
      'Manage products',
      'Manage orders',
      'View analytics',
      'System settings'
    ]
  },
  {
    id: 'moderator',
    label: 'Moderator',
    description: 'Content moderation and user management',
    icon: Shield,
    color: '#059669',
    permissions: [
      'Moderate content',
      'Manage reviews',
      'Basic user management',
      'View reports'
    ]
  },
  {
    id: 'support',
    label: 'Support Agent',
    description: 'Customer support and ticket management',
    icon: Users,
    color: '#7c3aed',
    permissions: [
      'Manage support tickets',
      'View user profiles',
      'Basic reporting',
      'Customer communication'
    ]
  },
  {
    id: 'analyst',
    label: 'Analyst',
    description: 'Data analysis and reporting access',
    icon: Settings,
    color: '#dc2626',
    permissions: [
      'View all analytics',
      'Generate reports',
      'Export data',
      'Performance metrics'
    ]
  }
];

const AddAdminManagement = () => {
  // Handle missing translations gracefully
  let t;
  try {
    t = useTranslations('AddAdminManagement');
  } catch (e) {
    t = (key, options) => options?.defaultValue || key;
  }

  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();

  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    password: '',
    confirmPassword: '',
    role: 'moderator',
    phone: '',
    department: '',
    location: '',
    notes: '',
    sendWelcomeEmail: true,
    requirePasswordChange: true
  });

  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitSuccess, setSubmitSuccess] = useState(false);

  // Form validation
  const validateForm = () => {
    const newErrors = {};

    // Required fields
    if (!formData.firstName.trim()) {
      newErrors.firstName = 'First name is required';
    }

    if (!formData.lastName.trim()) {
      newErrors.lastName = 'Last name is required';
    }

    if (!formData.email.trim()) {
      newErrors.email = 'Email is required';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = 'Please enter a valid email address';
    }

    if (!formData.password) {
      newErrors.password = 'Password is required';
    } else if (formData.password.length < 8) {
      newErrors.password = 'Password must be at least 8 characters';
    } else if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(formData.password)) {
      newErrors.password = 'Password must contain uppercase, lowercase, and number';
    }

    if (!formData.confirmPassword) {
      newErrors.confirmPassword = 'Please confirm your password';
    } else if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match';
    }

    if (!formData.role) {
      newErrors.role = 'Please select a role';
    }

    if (formData.phone && !/^\+?[\d\s\-()]+$/.test(formData.phone)) {
      newErrors.phone = 'Please enter a valid phone number';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Handle form input changes
  const handleInputChange = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    
    // Clear error when user starts typing
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);
    
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 2000));

      setSubmitSuccess(true);
      
      // Reset form after success
      setTimeout(() => {
        router.push('/admin/users');
      }, 2000);
      
    } catch (error) {
      // Error: 'Error creating admin:', error...
      setErrors({ submit: 'Failed to create admin user. Please try again.' });
    } finally {
      setIsSubmitting(false);
    }
  };

  // Get selected role info
  const selectedRole = adminRoles.find(role => role.id === formData.role);

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>Access Denied</h2>
          <p>You need admin privileges to access this page.</p>
          <p>Current role: {user?.role || 'Not logged in'}</p>
        </div>
      </div>
    );
  }

  if (submitSuccess) {
    return (
      <div className={styles.container}>
        <div className={styles.successState}>
          <div className={styles.successIcon}>
            <Check size={48} />
          </div>
          <h2>Admin User Created Successfully!</h2>
          <p>The new administrator account has been created and a welcome email has been sent.</p>
          <div className={styles.successActions}>
            <button 
              className={styles.primaryButton}
              onClick={() => router.push('/admin/users')}
            >
              View All Users
            </button>
            <button 
              className={styles.secondaryButton}
              onClick={() => window.location.reload()}
            >
              Create Another Admin
            </button>
          </div>
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
              onClick={() => router.back()}
            >
              <ArrowLeft size={16} />
              Back
            </button>
            <div>
              <h1 className={styles.title}>
                {t('title', { defaultValue: 'Add Administrator' })}
              </h1>
              <p className={styles.subtitle}>
                {t('subtitle', { defaultValue: 'Create a new administrator account with specific permissions' })}
              </p>
            </div>
          </div>
          <div className={styles.headerActions}>
            <UserPlus size={24} className={styles.headerIcon} />
          </div>
        </div>

        <div className={styles.content}>
          {/* Role Selection Cards */}
          <div className={styles.roleSection}>
            <h3 className={styles.sectionTitle}>Select Administrator Role</h3>
            <div className={styles.roleGrid}>
              {adminRoles.map(role => {
                const IconComponent = role.icon;
                return (
                  <div
                    key={role.id}
                    className={`${styles.roleCard} ${formData.role === role.id ? styles.roleCardSelected : ''}`}
                    onClick={() => handleInputChange('role', role.id)}
                  >
                    <div className={styles.roleCardHeader}>
                      <div 
                        className={styles.roleIcon}
                        style={{ backgroundColor: `${role.color}20`, color: role.color }}
                      >
                        <IconComponent size={20} />
                      </div>
                      <div className={styles.roleInfo}>
                        <h4 className={styles.roleLabel}>{role.label}</h4>
                        <p className={styles.roleDescription}>{role.description}</p>
                      </div>
                    </div>
                    <div className={styles.rolePermissions}>
                      <span className={styles.permissionsLabel}>Permissions:</span>
                      <ul className={styles.permissionsList}>
                        {role.permissions.slice(0, 3).map((permission, index) => (
                          <li key={index}>{permission}</li>
                        ))}
                        {role.permissions.length > 3 && (
                          <li>+{role.permissions.length - 3} more</li>
                        )}
                      </ul>
                    </div>
                  </div>
                );
              })}
            </div>
            {errors.role && <span className={styles.fieldError}>{errors.role}</span>}
          </div>

          {/* Admin Creation Form */}
          <form onSubmit={handleSubmit} className={styles.form}>
            {/* Personal Information */}
            <div className={styles.formSection}>
              <h3 className={styles.sectionTitle}>Personal Information</h3>
              <div className={styles.formGrid}>
                <div className={styles.formGroup}>
                  <label className={styles.label}>
                    <User size={16} />
                    First Name *
                  </label>
                  <input
                    type="text"
                    value={formData.firstName}
                    onChange={(e) => handleInputChange('firstName', e.target.value)}
                    className={`${styles.input} ${errors.firstName ? styles.inputError : ''}`}
                    placeholder="Enter first name"
                  />
                  {errors.firstName && <span className={styles.fieldError}>{errors.firstName}</span>}
                </div>

                <div className={styles.formGroup}>
                  <label className={styles.label}>
                    <User size={16} />
                    Last Name *
                  </label>
                  <input
                    type="text"
                    value={formData.lastName}
                    onChange={(e) => handleInputChange('lastName', e.target.value)}
                    className={`${styles.input} ${errors.lastName ? styles.inputError : ''}`}
                    placeholder="Enter last name"
                  />
                  {errors.lastName && <span className={styles.fieldError}>{errors.lastName}</span>}
                </div>

                <div className={styles.formGroup}>
                  <label className={styles.label}>
                    <Mail size={16} />
                    Email Address *
                  </label>
                  <input
                    type="email"
                    value={formData.email}
                    onChange={(e) => handleInputChange('email', e.target.value)}
                    className={`${styles.input} ${errors.email ? styles.inputError : ''}`}
                    placeholder="Enter email address"
                  />
                  {errors.email && <span className={styles.fieldError}>{errors.email}</span>}
                </div>

                <div className={styles.formGroup}>
                  <label className={styles.label}>
                    <Phone size={16} />
                    Phone Number
                  </label>
                  <input
                    type="tel"
                    value={formData.phone}
                    onChange={(e) => handleInputChange('phone', e.target.value)}
                    className={`${styles.input} ${errors.phone ? styles.inputError : ''}`}
                    placeholder="Enter phone number"
                  />
                  {errors.phone && <span className={styles.fieldError}>{errors.phone}</span>}
                </div>
              </div>
            </div>

            {/* Account Security */}
            <div className={styles.formSection}>
              <h3 className={styles.sectionTitle}>Account Security</h3>
              <div className={styles.formGrid}>
                <div className={styles.formGroup}>
                  <label className={styles.label}>
                    <Lock size={16} />
                    Password *
                  </label>
                  <div className={styles.passwordInput}>
                    <input
                      type={showPassword ? 'text' : 'password'}
                      value={formData.password}
                      onChange={(e) => handleInputChange('password', e.target.value)}
                      className={`${styles.input} ${errors.password ? styles.inputError : ''}`}
                      placeholder="Enter secure password"
                    />
                    <button
                      type="button"
                      className={styles.passwordToggle}
                      onClick={() => setShowPassword(!showPassword)}
                    >
                      {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                    </button>
                  </div>
                  {errors.password && <span className={styles.fieldError}>{errors.password}</span>}
                  <div className={styles.passwordHint}>
                    Must contain at least 8 characters with uppercase, lowercase, and number
                  </div>
                </div>

                <div className={styles.formGroup}>
                  <label className={styles.label}>
                    <Lock size={16} />
                    Confirm Password *
                  </label>
                  <div className={styles.passwordInput}>
                    <input
                      type={showConfirmPassword ? 'text' : 'password'}
                      value={formData.confirmPassword}
                      onChange={(e) => handleInputChange('confirmPassword', e.target.value)}
                      className={`${styles.input} ${errors.confirmPassword ? styles.inputError : ''}`}
                      placeholder="Confirm password"
                    />
                    <button
                      type="button"
                      className={styles.passwordToggle}
                      onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    >
                      {showConfirmPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                    </button>
                  </div>
                  {errors.confirmPassword && <span className={styles.fieldError}>{errors.confirmPassword}</span>}
                </div>
              </div>
            </div>

            {/* Additional Information */}
            <div className={styles.formSection}>
              <h3 className={styles.sectionTitle}>Additional Information</h3>
              <div className={styles.formGrid}>
                <div className={styles.formGroup}>
                  <label className={styles.label}>
                    <Settings size={16} />
                    Department
                  </label>
                  <select
                    value={formData.department}
                    onChange={(e) => handleInputChange('department', e.target.value)}
                    className={styles.input}
                  >
                    <option value="">Select department</option>
                    <option value="operations">Operations</option>
                    <option value="marketing">Marketing</option>
                    <option value="sales">Sales</option>
                    <option value="support">Customer Support</option>
                    <option value="finance">Finance</option>
                    <option value="tech">Technology</option>
                    <option value="hr">Human Resources</option>
                  </select>
                </div>

                <div className={styles.formGroup}>
                  <label className={styles.label}>
                    <MapPin size={16} />
                    Location
                  </label>
                  <input
                    type="text"
                    value={formData.location}
                    onChange={(e) => handleInputChange('location', e.target.value)}
                    className={styles.input}
                    placeholder="Enter location or office"
                  />
                </div>
              </div>

              <div className={styles.formGroup}>
                <label className={styles.label}>
                  Notes
                </label>
                <textarea
                  value={formData.notes}
                  onChange={(e) => handleInputChange('notes', e.target.value)}
                  className={styles.textarea}
                  placeholder="Additional notes about this admin user (optional)"
                  rows={3}
                />
              </div>
            </div>

            {/* Account Options */}
            <div className={styles.formSection}>
              <h3 className={styles.sectionTitle}>Account Options</h3>
              <div className={styles.checkboxGroup}>
                <label className={styles.checkboxLabel}>
                  <input
                    type="checkbox"
                    checked={formData.sendWelcomeEmail}
                    onChange={(e) => handleInputChange('sendWelcomeEmail', e.target.checked)}
                    className={styles.checkbox}
                  />
                  <span className={styles.checkboxText}>
                    <Mail size={16} />
                    Send welcome email with login instructions
                  </span>
                </label>

                <label className={styles.checkboxLabel}>
                  <input
                    type="checkbox"
                    checked={formData.requirePasswordChange}
                    onChange={(e) => handleInputChange('requirePasswordChange', e.target.checked)}
                    className={styles.checkbox}
                  />
                  <span className={styles.checkboxText}>
                    <Shield size={16} />
                    Require password change on first login
                  </span>
                </label>
              </div>
            </div>

            {/* Selected Role Summary */}
            {selectedRole && (
              <div className={styles.rolesummary}>
                <h3 className={styles.sectionTitle}>Selected Role Summary</h3>
                <div className={styles.roleSummaryCard}>
                  <div className={styles.roleSummaryHeader}>
                    <div 
                      className={styles.roleIcon}
                      style={{ backgroundColor: `${selectedRole.color}20`, color: selectedRole.color }}
                    >
                      <selectedRole.icon size={20} />
                    </div>
                    <div>
                      <h4>{selectedRole.label}</h4>
                      <p>{selectedRole.description}</p>
                    </div>
                  </div>
                  <div className={styles.allPermissions}>
                    <span className={styles.permissionsLabel}>Full Permissions:</span>
                    <div className={styles.permissionsGrid}>
                      {selectedRole.permissions.map((permission, index) => (
                        <span key={index} className={styles.permissionTag}>
                          <Check size={12} />
                          {permission}
                        </span>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Submit Errors */}
            {errors.submit && (
              <div className={styles.submitError}>
                <AlertCircle size={16} />
                {errors.submit}
              </div>
            )}

            {/* Form Actions */}
            <div className={styles.formActions}>
              <button
                type="button"
                className={styles.secondaryButton}
                onClick={() => router.back()}
                disabled={isSubmitting}
              >
                Cancel
              </button>
              <button
                type="submit"
                className={styles.primaryButton}
                disabled={isSubmitting}
              >
                {isSubmitting ? (
                  <>
                    <RefreshCw size={16} className={styles.spinning} />
                    Creating Admin...
                  </>
                ) : (
                  <>
                    <Save size={16} />
                    Create Administrator
                  </>
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default AddAdminManagement; 