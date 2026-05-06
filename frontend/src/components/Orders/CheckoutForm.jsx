// src/features/Orders/components/CheckoutForm.jsx
"use client";
import React, { useState, memo, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { toast } from 'react-toastify';
import { createOrder } from '@/api/client/orderingApi';
import { getBasket, checkoutBasket } from '@/api/client/basketApi';
import { useAuth } from '@/context/AuthContext';
import styles from './CheckoutForm.module.css';

const CheckoutForm = memo(function CheckoutForm({ basketId, locale = 'en' }) {
    const router = useRouter();
    const t = useTranslations('Checkout');
    const { user } = useAuth();
    const [formData, setFormData] = useState({
        email: '',
        firstName: '',
        lastName: '',
        address: '',
        city: '',
        postalCode: '',
        country: '',
        phone: '',
        paymentMethod: 'card'
    });
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [basketData, setBasketData] = useState(null);
    const [basketLoading, setBasketLoading] = useState(true);

    // Load basket data on mount
    useEffect(() => {
        const loadBasket = async () => {
            if (!basketId) return;
            
            try {
                setBasketLoading(true);
                const response = await getBasket(basketId);
                if (response && response.basket) {
                    setBasketData(response.basket);
                } else {
                    toast.error('Basket not found');
                    router.push(`/${locale}/basket`);
                }
            } catch (error) {
                // Error: 'Failed to load basket:', error...
                toast.error('Failed to load basket');
            } finally {
                setBasketLoading(false);
            }
        };
        
        loadBasket();
    }, [basketId, locale, router]);

    // Pre-fill user email if available
    useEffect(() => {
        if (user?.email) {
            setFormData(prev => ({ ...prev, email: user.email }));
        }
    }, [user]);

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setIsSubmitting(true);
        
        try {
            // Validate form data
            if (!formData.email || !formData.firstName || !formData.lastName || 
                !formData.address || !formData.city || !formData.postalCode || !formData.country) {
                toast.error('Please fill in all required fields');
                setIsSubmitting(false);
                return;
            }

            if (!basketData || !basketData.items || basketData.items.length === 0) {
                toast.error('Your basket is empty');
                setIsSubmitting(false);
                return;
            }

            // Prepare order items from basket
            const orderItems = basketData.items.map(item => ({
                productId: item.product_id,
                productName: item.product_name,
                userSellerId: item.user_seller_id,
                userSellerName: item.user_seller_name || '',
                price: parseFloat(item.price),
                quantity: parseInt(item.quantity, 10)
            }));

            // Create the order
            const orderData = {
                userCustomerId: user?.userId || basketData.user_id,
                paymentMethodId: formData.paymentMethod,
                items: orderItems
            };

            const orderResponse = await createOrder(orderData);
            
            if (orderResponse.success) {
                // Note: We'll checkout the basket after payment is completed
                // For now, just store the order info
                
                // Store order info in sessionStorage for payment page
                sessionStorage.setItem('pendingOrder', JSON.stringify({
                    orderId: orderResponse.id,
                    basketId: basketId,
                    paymentMethod: formData.paymentMethod,
                    shippingAddress: {
                        firstName: formData.firstName,
                        lastName: formData.lastName,
                        address: formData.address,
                        city: formData.city,
                        postalCode: formData.postalCode,
                        country: formData.country,
                        phone: formData.phone
                    },
                    email: formData.email,
                    total: calculateTotal()
                }));
                
                // Navigate to payment page
                router.push(`/${locale}/checkout/payment`);
            } else {
                toast.error(orderResponse.userMessage || 'Failed to create order');
            }
        } catch (error) {
            // Error: 'Checkout error:', error...
            toast.error('Failed to process checkout. Please try again.');
        } finally {
            setIsSubmitting(false);
        }
    };

    const calculateTotal = () => {
        if (!basketData || !basketData.items) return 0;
        return basketData.items.reduce((total, item) => 
            total + (parseFloat(item.price) * parseInt(item.quantity, 10)), 0
        );
    };
    return (
        <form className={styles.checkoutForm} onSubmit={handleSubmit}>
            <div className={styles.section}>
                <h3>{t('contactInformation')}</h3>
                <div className={styles.field}>
                    <label htmlFor="email">{t('email')}</label>
                    <input
                        type="email"
                        id="email"
                        name="email"
                        value={formData.email}
                        onChange={handleInputChange}
                        required
                    />
                </div>
            </div>
            <div className={styles.section}>
                <h3>{t('shippingAddress')}</h3>
                <div className={styles.fieldRow}>
                    <div className={styles.field}>
                        <label htmlFor="firstName">{t('firstName')}</label>
                        <input
                            type="text"
                            id="firstName"
                            name="firstName"
                            value={formData.firstName}
                            onChange={handleInputChange}
                            required
                        />
                    </div>
                    <div className={styles.field}>
                        <label htmlFor="lastName">{t('lastName')}</label>
                        <input
                            type="text"
                            id="lastName"
                            name="lastName"
                            value={formData.lastName}
                            onChange={handleInputChange}
                            required
                        />
                    </div>
                </div>
                <div className={styles.field}>
                    <label htmlFor="address">{t('address')}</label>
                    <input
                        type="text"
                        id="address"
                        name="address"
                        value={formData.address}
                        onChange={handleInputChange}
                        required
                    />
                </div>
                <div className={styles.fieldRow}>
                    <div className={styles.field}>
                        <label htmlFor="city">{t('city')}</label>
                        <input
                            type="text"
                            id="city"
                            name="city"
                            value={formData.city}
                            onChange={handleInputChange}
                            required
                        />
                    </div>
                    <div className={styles.field}>
                        <label htmlFor="postalCode">{t('postalCode')}</label>
                        <input
                            type="text"
                            id="postalCode"
                            name="postalCode"
                            value={formData.postalCode}
                            onChange={handleInputChange}
                            required
                        />
                    </div>
                </div>
                <div className={styles.fieldRow}>
                    <div className={styles.field}>
                        <label htmlFor="country">{t('country')}</label>
                        <input
                            type="text"
                            id="country"
                            name="country"
                            value={formData.country}
                            onChange={handleInputChange}
                            required
                        />
                    </div>
                    <div className={styles.field}>
                        <label htmlFor="phone">{t('phone')}</label>
                        <input
                            type="tel"
                            id="phone"
                            name="phone"
                            value={formData.phone}
                            onChange={handleInputChange}
                        />
                    </div>
                </div>
            </div>
            <div className={styles.section}>
                <h3>{t('payment')}</h3>
                <div className={styles.field}>
                    <label htmlFor="paymentMethod">{t('paymentMethod')}</label>
                    <select
                        id="paymentMethod"
                        name="paymentMethod"
                        value={formData.paymentMethod}
                        onChange={handleInputChange}
                    >
                        <option value="card">{t('creditCard')}</option>
                        <option value="paypal">{t('paypal')}</option>
                        <option value="bank">{t('bankTransfer')}</option>
                    </select>
                </div>
            </div>
            {basketLoading ? (
                <div className={styles.loading}>
                    <p>{t('loadingBasket')}</p>
                </div>
            ) : basketData && basketData.items && basketData.items.length > 0 ? (
                <div className={styles.orderSummary}>
                    <h4>{t('orderSummary')}</h4>
                    <p>{t('items')}: {basketData.items.length}</p>
                    <p>{t('total')}: ${calculateTotal().toFixed(2)}</p>
                </div>
            ) : null}
            
            <button
                type="submit"
                className={styles.submitButton}
                disabled={isSubmitting || basketLoading || !basketData || !basketData.items || basketData.items.length === 0}
            >
                {isSubmitting ? t('processing') : t('proceedToPayment')}
            </button>
        </form>
    );
});
export default CheckoutForm;
