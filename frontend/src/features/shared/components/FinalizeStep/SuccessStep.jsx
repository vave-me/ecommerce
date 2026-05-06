import React from 'react';
import PropTypes from 'prop-types';
import { Check, ExternalLink, Plus } from "@/icons";
/**
 * Shared SuccessStep Component
 * Used across all creation modals for consistent success messaging
 */
export function SuccessStep({
    itemData,
    itemType = "listing",
    onClose,
    styles,
    // Configurable content
    title,
    message,
    nextSteps,
    primaryActionLabel = "Create Another",
    secondaryActionLabel = "View My Listings",
    onViewListings,
    // Additional actions
    additionalActions = []
}) {
    // Default configurations based on item type
    const getDefaultConfig = () => {
        switch (itemType) {
            case 'deal':
                return {
                    title: "Deal Published!",
                    message: "Your deal listing has been successfully published.",
                    nextSteps: [
                        "Share your listing link on social media or with potential buyers.",
                        "Keep an eye out for messages or offers from interested buyers.",
                        "Prepare your item for potential shipping or pickup."
                    ]
                };
            case 'job':
                return {
                    title: "Job Posted!",
                    message: "Your job posting has been successfully published.",
                    nextSteps: [
                        "Review applications as they come in.",
                        "Respond to qualified candidates promptly.",
                        "Keep your job posting updated with any changes."
                    ]
                };
            case 'service':
                return {
                    title: "Service Listed!",
                    message: "Your service listing has been successfully published.",
                    nextSteps: [
                        "Respond to service inquiries quickly.",
                        "Keep your availability calendar updated.",
                        "Showcase your work with photos and reviews."
                    ]
                };
            case 'product':
                return {
                    title: "Product Listed!",
                    message: "Your product listing has been successfully published.",
                    nextSteps: [
                        "Monitor your inventory levels.",
                        "Respond to customer questions promptly.",
                        "Consider offering promotions to boost sales."
                    ]
                };
            case 'property':
                return {
                    title: "Property Listed!",
                    message: "Your property listing has been successfully published.",
                    nextSteps: [
                        "Schedule viewings with interested parties.",
                        "Keep your listing photos up to date.",
                        "Respond to inquiries about the property."
                    ]
                };
            case 'vehicle':
                return {
                    title: "Vehicle Listed!",
                    message: "Your vehicle listing has been successfully published.",
                    nextSteps: [
                        "Prepare your vehicle for potential test drives.",
                        "Have all necessary documents ready.",
                        "Respond to buyer inquiries promptly."
                    ]
                };
            case 'tweet':
            case 'post':
                return {
                    title: "Post Published!",
                    message: "Your post has been successfully published.",
                    nextSteps: [
                        "Engage with comments and reactions.",
                        "Share your post to increase visibility.",
                        "Monitor the performance of your post."
                    ]
                };
            default:
                return {
                    title: "Published Successfully!",
                    message: "Your listing has been successfully published.",
                    nextSteps: [
                        "Monitor your listing for activity.",
                        "Respond to inquiries promptly.",
                        "Keep your information up to date."
                    ]
                };
        }
    };
    const defaultConfig = getDefaultConfig();
    const finalTitle = title || defaultConfig.title;
    const finalMessage = message || defaultConfig.message;
    const finalNextSteps = nextSteps || defaultConfig.nextSteps;
    const handleViewListings = () => {
        if (onViewListings) {
            onViewListings();
        } else {
            // Default navigation based on item type
            const routes = {
                deal: '/deals/my-listings',
                job: '/jobs/my-listings',
                service: '/services/my-listings',
                product: '/products/my-listings',
                property: '/properties/my-listings',
                vehicle: '/vehicles/my-listings',
                tweet: '/posts/my-posts',
                post: '/posts/my-posts'
            };
            const route = routes[itemType] || '/my-listings';
            window.location.href = route;
        }
    };
    return (
        <div className={styles.successContainer}>
            <div className={styles.successIcon}>
                <Check size={48} aria-hidden="true" />
            </div>
            <h2 className={styles.successTitle}>{finalTitle}</h2>
            <p className={styles.successMessage}>
                {finalMessage}
            </p>
            {/* Optional: Show a small summary of the published item */}
            {itemData && (itemData.name || itemData.thumbnail || (itemData.images && itemData.images[0])) && (
                <div className={styles.itemSummary || styles.dealSummary}>
                    {(itemData.thumbnail || (itemData.images && itemData.images[0])) && (
                        <div className={styles.thumbnailPreview}>
                            <img 
                                src={itemData.thumbnail || itemData.images[0]} 
                                alt={itemData.name || `${itemType} image`}
                            />
                        </div>
                    )}
                    {itemData.name && (
                        <div 
                            className={styles.itemDetails || styles.dealDetails}
                            style={!(itemData.thumbnail || (itemData.images && itemData.images[0])) ? 
                                {textAlign: 'center', width: '100%'} : {}
                            }
                        >
                            <h3>{itemData.name}</h3>
                            {itemData.price && (
                                <p className={styles.itemPrice}>
                                    {typeof itemData.price === 'number' ? 
                                        `€${itemData.price.toFixed(2)}` : 
                                        itemData.price
                                    }
                                </p>
                            )}
                        </div>
                    )}
                </div>
            )}
            {/* Next Steps */}
            <div className={styles.nextStepOptions}>
                <h4>What's Next?</h4>
                <ul className={styles.nextStepsList}>
                    {finalNextSteps.map((step, index) => (
                        <li key={index}>{step}</li>
                    ))}
                </ul>
            </div>
            {/* Additional Actions */}
            {additionalActions.length > 0 && (
                <div className={styles.additionalActions}>
                    {additionalActions.map((action, index) => (
                        <button
                            key={index}
                            className={action.className || styles.tertiaryButton}
                            type="button"
                            onClick={action.onClick}
                        >
                            {action.icon && <action.icon size={16} />}
                            {action.label}
                        </button>
                    ))}
                </div>
            )}
            {/* Main Actions */}
            <div className={styles.successActions}>
                <button
                    className={styles.secondaryButton}
                    type="button"
                    onClick={handleViewListings}
                >
                    <ExternalLink size={16} />
                    {secondaryActionLabel}
                </button>
                <button
                    className={styles.primaryButton}
                    type="button"
                    onClick={onClose}
                >
                    <Plus size={16} />
                    {primaryActionLabel}
                </button>
            </div>
        </div>
    );
}
SuccessStep.propTypes = {
    itemData: PropTypes.shape({
        name: PropTypes.string,
        thumbnail: PropTypes.string,
        images: PropTypes.array,
        price: PropTypes.oneOfType([PropTypes.string, PropTypes.number])
    }),
    itemType: PropTypes.string,
    onClose: PropTypes.func.isRequired,
    styles: PropTypes.object.isRequired,
    title: PropTypes.string,
    message: PropTypes.string,
    nextSteps: PropTypes.array,
    primaryActionLabel: PropTypes.string,
    secondaryActionLabel: PropTypes.string,
    onViewListings: PropTypes.func,
    additionalActions: PropTypes.arrayOf(PropTypes.shape({
        label: PropTypes.string.isRequired,
        onClick: PropTypes.func.isRequired,
        icon: PropTypes.elementType,
        className: PropTypes.string
    }))
}; 