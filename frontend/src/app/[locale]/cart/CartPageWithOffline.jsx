"use client";
import React, {useCallback, useEffect, useRef, useState} from "react";
import {useTranslations} from "next-intl";
import {useRouter} from "next/navigation";
import {
    checkoutBasket,
    getBasket,
    getOrCreateBasket,
    removeItemFromBasket,
    updateItemQuantity
} from "../../../api/client/basketApi";
import {useAuth} from "../../../context/AuthContext";
import { useBasket } from "../../../hooks/useBasket";
import {offlineBasketStorage} from "../../../utils/offlineStorage";
import styles from "./BasketPage.module.css";

export default function CartPageWithOffline() {
    const t = useTranslations("BasketPage");
    const {user} = useAuth();
    const router = useRouter();
    const userId = user?.userId;
    const isLoggedIn = !!userId;
    
    // Use the unified basket hook
    const {
        basket,
        isLoading: loading,
        error,
        refetchBasket: fetchBasket,
        removeItem,
        updateQuantity,
        initiateCheckout,
        isSyncing: syncInProgress
    } = useBasket();
    
    const [offlineItems, setOfflineItems] = useState([]);
    const hasFetched = useRef(false);
    
    // Load basket data on mount
    useEffect(() => {
        if (isLoggedIn && !hasFetched.current) {
            hasFetched.current = true;
            fetchBasket();
        } else if (!isLoggedIn) {
            // Load offline items
            const items = offlineBasketStorage.getItems();
            setOfflineItems(items);
        }
    }, [isLoggedIn, fetchBasket]);
    
    // Update offline items when they change
    useEffect(() => {
        if (!isLoggedIn) {
            const handleStorageChange = () => {
                const items = offlineBasketStorage.getItems();
                setOfflineItems(items);
            };
            
            // Listen for storage changes from other tabs
            window.addEventListener('storage', handleStorageChange);
            return () => window.removeEventListener('storage', handleStorageChange);
        }
    }, [isLoggedIn]);
    
    // Handle remove item
    const handleRemoveItem = useCallback(async (itemId) => {
        // The new hook handles both online and offline automatically
        await removeItem(itemId);
        
        // Update offline items display if not logged in
        if (!isLoggedIn) {
            const items = offlineBasketStorage.getItems();
            setOfflineItems(items);
        }
    }, [isLoggedIn, removeItem]);
    
    // Handle quantity update
    const handleUpdateQuantity = useCallback(async (itemId, newQuantity) => {
        if (newQuantity < 1) {
            // Quantity must be at least 1
            return;
        }
        
        // The new hook handles both online and offline automatically
        await updateQuantity(itemId, newQuantity);
        
        // Update offline items display if not logged in
        if (!isLoggedIn) {
            const items = offlineBasketStorage.getItems();
            setOfflineItems(items);
        }
    }, [isLoggedIn, updateQuantity, t]);
    
    // Get items to display - handle both camelCase and snake_case
    const displayItems = isLoggedIn ? (basket?.items || []) : offlineItems;
    const hasItems = displayItems.length > 0;
    
    // Calculate totals
    const subtotal = displayItems.reduce((sum, item) => {
        // Handle different field names from backend
        const price = parseFloat(item.product_price || item.productPrice || item.price || 0);
        const quantity = item.quantity || 1;
        return sum + (price * quantity);
    }, 0);
    
    const shipping = isLoggedIn ? (basket?.shippingCost || 2.99) : 2.99;
    const discount = isLoggedIn ? (basket?.discount || 0) : 0;
    const total = subtotal + shipping - discount;
    
    // Handle checkout
    const handleCheckoutClick = useCallback(async () => {
        if (!isLoggedIn) {
            // Redirect to login with return URL
            router.push('/login?returnUrl=/cart');
            
            return;
        }
        
        if (basket?.id && displayItems.length > 0) {
            // Calculate total in cents for payment
            const totalInCents = Math.round(total * 100);
            
            // Navigate to payment page with necessary data
            const params = new URLSearchParams({
                basketId: basket.id,
                userCustomerId: userId,
                amount: totalInCents.toString()
            });
            router.push(`/payments?${params.toString()}`);
        }
    }, [isLoggedIn, basket, displayItems.length, total, userId, router]);
    
    // Render states
    if (loading || syncInProgress) {
        return (
            <div className={styles.loadingContainer}>
                {syncInProgress ? t("syncingItems", "Syncing your saved items...") : t("loading")}
            </div>
        );
    }
    
    if (error) {
        return <div className={styles.errorMessage}>{error?.message || error}</div>;
    }
    
    if (!hasItems) {
        return (
            <div className={styles.emptyBasketContainer}>
                <div className={styles.emptyBasketIcon}>🛒</div>
                <h2 className={styles.emptyBasketTitle}>{t("empty")}</h2>
                <p className={styles.emptyBasketDescription}>
                    {isLoggedIn ? t("emptyDescription") : t("emptyDescriptionOffline", "Your cart is empty. Items you add will be saved locally until you log in.")}
                </p>
                <button 
                    className={styles.continueShoppingButton}
                    onClick={() => router.push('/')}
                >
                    {t("continueShopping")}
                </button>
            </div>
        );
    }
    
    return (
        <div className={styles.basketContainer}>
            <h2 className={styles.basketHeading}>
                {t("heading")}
                {!isLoggedIn && (
                    <span className={styles.offlineIndicator}>
                        {" "}({t("offlineMode", "Offline Mode")})
                    </span>
                )}
            </h2>
            
            {!isLoggedIn && (
                <div className={styles.offlineNotice}>
                    <p>{t("offlineNotice", "Your items are saved locally. Log in to complete your purchase.")}</p>
                </div>
            )}
            
            <div className={styles.basketItems}>
                {displayItems.map((item) => (
                    <BasketItemRow
                        key={item.productId}
                        item={item}
                        onRemove={handleRemoveItem}
                        onUpdateQuantity={handleUpdateQuantity}
                        isOffline={!isLoggedIn}
                    />
                ))}
            </div>
            
            <div className={styles.summarySection}>
                <Row label={t("labels.subtotal")} value={`€${subtotal.toFixed(2)}`}/>
                <Row label={t("labels.shipping")} value={`€${shipping.toFixed(2)}`}/>
                {discount > 0 && <Row label={t("labels.discount")} value={`-€${discount.toFixed(2)}`}/>}
                <hr className={styles.divider}/>
                <Row label={t("labels.total")} value={`€${total.toFixed(2)}`} total/>
            </div>
            
            <div className={styles.actionsRow}>
                {isLoggedIn && (
                    <button 
                        className={`${styles.baseButton} ${styles.refreshButton}`} 
                        onClick={fetchBasket}
                    >
                        {t("actions.refresh")}
                    </button>
                )}
                <button 
                    className={`${styles.baseButton} ${styles.checkoutButton}`} 
                    onClick={handleCheckoutClick}
                >
                    {isLoggedIn ? t("actions.checkout") : t("actions.loginToCheckout", "Login to Checkout")}
                </button>
            </div>
        </div>
    );
}

/* ---------- helpers ------------------------------- */
function Row({label, value, total}) {
    return (
        <div className={styles.row}>
            <span className={total ? styles.totalLabel : styles.rowLabel}>{label}</span>
            <span className={total ? styles.totalPrice : styles.rowValue}>{value}</span>
        </div>
    );
}

function BasketItemRow({item, onRemove, onUpdateQuantity, isOffline}) {
    const t = useTranslations("BasketPage");
    // Handle both camelCase and snake_case field names from backend
    const price = parseFloat(item.product_price || item.productPrice || item.price || 0);
    const productId = item.product_id || item.productId;
    const itemId = item.item_id || item.id || productId; // For remove/update operations
    const productName = item.product_name || item.productName || item.name || "Product";
    const quantity = item.quantity || 1;
    
    return (
        <div className={`${styles.basketItem} ${styles.transitionAll}`}>
            <div className={styles.itemImageWrapper}>
                <img 
                    src={item.imageUrl || item.image || "/images/default-product.webp"} 
                    alt={productName} 
                    className={styles.itemImage}
                    onError={(e) => {
                        e.target.src = "/images/default-product.webp";
                    }}
                />
            </div>
            <div className={styles.itemDetails}>
                <div className={styles.itemName}>{productName}</div>
                {item.description && (
                    <div className={styles.itemDescription}>{item.description}</div>
                )}
                {item.sku && (
                    <div className={styles.skuLabel}>{t("sku", {sku: item.sku})}</div>
                )}
            </div>
            <div className={styles.itemQuantity}>
                {t("quantity")}
                <input
                    className={styles.quantityInput}
                    type="number"
                    value={quantity}
                    min="1"
                    onChange={e => onUpdateQuantity(itemId, parseInt(e.target.value, 10))}
                    aria-label={t("aria.quantityFor", {name: productName})}
                />
            </div>
            <div className={styles.itemPrice}>€{price.toFixed(2)}</div>
            <button 
                className={styles.removeButton} 
                onClick={() => onRemove(itemId)}
            >
                {t("actions.remove")}
            </button>
        </div>
    );
}