"use client";
import React, { useCallback, useEffect, useMemo, useRef, useState, memo } from 'react';
import {
    Calendar, Clock, Eye, MapPin, MessageCircle, Star, Tag, Heart, ThumbsUp, ThumbsDown, 
    Share2, Check, AlertCircle, Bookmark, Shield, Users, User, Phone, Mail, Globe, 
    ChevronLeft, ChevronRight, Camera, Award, CreditCard, Info, Layers, CheckCircle,
    Briefcase, Target, Timer, DollarSign, BarChart3, TrendingUp, Zap, Home, ArrowRight
} from '@/icons';
import Link from 'next/link';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { useAuth } from '../../context/AuthContext';
import { useDispatch } from 'react-redux';
import { openMessageModal } from '../../redux/slices/modalsSlice';
import useActivityApi from '../../hooks/useActivityApi';
import { getMediaByItem } from '../../api/mediaApi';
import CommentsSetup from '../../features/Comments/CommentsSetup';
import styles from './DetailedServiceView.module.css';

// Extend dayjs with relative time plugin
dayjs.extend(relativeTime);

/**
 * StarRating - Service rating component
 */
const StarRating = memo(({ rating, reviewCount }) => {
    const stars = [];
    const fullStars = Math.floor(rating);
    const hasHalfStar = rating - fullStars >= 0.5;

    for (let i = 0; i < 5; i++) {
        if (i < fullStars) {
            stars.push(<Star key={i} size={20} className={styles.starFilled} />);
        } else if (i === fullStars && hasHalfStar) {
            stars.push(<Star key={i} size={20} className={styles.starHalf} />);
        } else {
            stars.push(<Star key={i} size={20} className={styles.starEmpty} />);
        }
    }

    return (
        <div className={styles.starRating} aria-label={`${rating} stars out of 5`}>
            {stars}
            {reviewCount > 0 && (
                <span className={styles.reviewCount}>({reviewCount} reviews)</span>
            )}
        </div>
    );
});

StarRating.displayName = 'StarRating';

/**
 * DetailedServiceView - Full service detail page component
 * Shows comprehensive service information with enhanced layout
 */
const DetailedServiceView = memo(({ service, locale, category, availableCategories }) => {
    const { user } = useAuth();
    const userId = user?.userId;
    const dispatch = useDispatch();
    const { handleLike, handleDislike } = useActivityApi();
    
    // Process service data
    const serviceData = useMemo(() => {
        const s = service?.service || service || {};
        const metrics = s.metrics || {};
        
        // Calculate rating from metrics
        const calculateRating = () => {
            const likes = parseInt(metrics?.likesCount || '0', 10);
            const dislikes = parseInt(metrics?.dislikesCount || '0', 10);
            const total = likes + dislikes;
            
            if (total === 0) return 4.5; // Default rating
            
            const ratio = likes / total;
            return Math.max(3.0, Math.min(5.0, 3.0 + (ratio * 2))); // Scale to 3.0-5.0
        };
        
        return {
            // Core data
            id: s.id || '',
            name: s.name || s.title || 'Unnamed Service',
            description: s.description || 'No description available.',
            basePrice: s.basePrice || s.price || '0',
            hourlyRate: s.hourlyRate || s.basePrice || '0',
            created: s.created || s.createdAt || new Date().toISOString(),
            availability: s.availability || 'Contact for availability',
            pricingModel: s.pricingModel || 'hourly',
            
            // Service details
            providerName: s.providerName || s.merchantName || 'Service Provider',
            serviceType: s.serviceType || 'Professional Service',
            duration: s.duration || 'Varies',
            experience: s.experience || '5+ years',
            certifications: s.certifications || [],
            languages: s.languages || ['English', 'German'],
            serviceArea: s.serviceArea || '50km radius',
            responseTime: s.responseTime || '< 2 hours',
            
            // Business info
            negotiable: s.negotiable !== false,
            instantBooking: s.instantBooking || false,
            onlineService: s.onlineService || false,
            onsiteService: s.onsiteService !== false,
            emergencyService: s.emergencyService || false,
            verified: s.verified !== false,
            featured: s.featured || false,
            categorySlug: s.categorySlug || 'services',
            categoryId: s.categoryId || '',
            tags: s.tags || [],
            
            // Location
            location: s.location || s.city || 'Germany',
            lat: s.lat || null,
            lng: s.lng || null,
            address: s.address || '',
            
            // Metrics
            metrics: {
                likes: parseInt(metrics?.likesCount || '0', 10),
                dislikes: parseInt(metrics?.dislikesCount || '0', 10),
                comments: parseInt(metrics?.commentsCount || '0', 10),
                shares: parseInt(metrics?.sharedCount || '0', 10),
                views: parseInt(metrics?.visitedCount || '0', 10),
                saved: parseInt(metrics?.addedToWishlistCount || '0', 10),
                bookings: parseInt(metrics?.bookingsCount || '0', 10),
            },
            
            // Reviews and ratings
            rating: calculateRating(),
            reviewCount: parseInt(metrics?.commentsCount || '0', 10),
            completedJobs: s.completedJobs || 342,
            
            // Service features
            features: s.features || [
                'Licensed and insured professional',
                'Free consultation included',
                'Satisfaction guarantee',
                'Available weekends',
                'Emergency service available'
            ],
            
            // What's included
            included: s.included || [
                'Initial consultation',
                'Professional assessment',
                'Detailed service plan',
                'All necessary equipment',
                'Clean-up after service',
                'Follow-up support'
            ],
            
            // Terms
            cancellationPolicy: s.cancellationPolicy || 'Free cancellation up to 24 hours before',
            paymentTerms: s.paymentTerms || 'Payment due upon completion',
            warranty: s.warranty || 'Service warranty for 30 days',
            
            // Provider info
            provider: {
                name: s.providerName || 'Professional Service Provider',
                rating: s.providerRating || 4.8,
                responseTime: s.responseTime || '< 2 hours',
                totalJobs: s.totalJobs || 1543,
                memberSince: s.memberSince || '2020',
                verificationBadges: s.verificationBadges || ['Identity Verified', 'Background Check', 'Insurance Verified'],
                specializations: s.specializations || ['Residential', 'Commercial', 'Emergency Services']
            }
        };
    }, [service]);
    
    // State management
    const [state, setState] = useState({
        currentImageIndex: 0,
        activeTab: 'description',
        showFullDescription: false,
        liked: false,
        disliked: false,
        favorite: false,
        showComments: false,
        zoomedImage: null,
        selectedDate: null,
        selectedTime: null,
        bookingNote: ''
    });
    
    const [mediaItems, setMediaItems] = useState([]);
    const [isLoadingMedia, setIsLoadingMedia] = useState(false);
    
    // Load media
    useEffect(() => {
        const loadMedia = async () => {
            if (!serviceData.id) return;
            
            setIsLoadingMedia(true);
            try {
                const mediaResponse = await getMediaByItem(serviceData.id);
                if (mediaResponse?.media?.mediaOrder?.length > 0) {
                    const formattedMedia = mediaResponse.media.mediaOrder.map(item => ({
                        id: item.mediaItemId || item.id,
                        url: item.url,
                        type: item.type || 'image',
                        alt: item.altText || serviceData.name || 'Service image'
                    }));
                    setMediaItems(formattedMedia);
                } else if (service.thumbnail) {
                    setMediaItems([{ url: service.thumbnail, type: 'image', alt: serviceData.name }]);
                }
            } catch (error) {
                if (service.thumbnail) {
                    setMediaItems([{ url: service.thumbnail, type: 'image', alt: serviceData.name }]);
                }
            } finally {
                setIsLoadingMedia(false);
            }
        };
        loadMedia();
    }, [serviceData.id, serviceData.name, service.thumbnail]);
    
    // Price formatting
    const formattedBasePrice = useMemo(() => {
        const price = parseFloat(serviceData.basePrice);
        return price.toLocaleString('de-DE', { style: 'currency', currency: 'EUR' });
    }, [serviceData.basePrice]);
    
    const formattedHourlyRate = useMemo(() => {
        const rate = parseFloat(serviceData.hourlyRate);
        return rate.toLocaleString('de-DE', { style: 'currency', currency: 'EUR' });
    }, [serviceData.hourlyRate]);
    
    // Handlers
    const handleImageNavigation = useCallback((direction) => {
        setState(prev => ({
            ...prev,
            currentImageIndex: direction === 'next'
                ? (prev.currentImageIndex + 1) % mediaItems.length
                : (prev.currentImageIndex - 1 + mediaItems.length) % mediaItems.length
        }));
    }, [mediaItems.length]);
    
    const handleBookService = useCallback(() => {
        if (!userId) {
            return;
        }
        
        // Here you would implement the booking logic
    }, [userId]);
    
    const handleLikeClick = useCallback(() => {
        if (!userId) {
            return;
        }
        handleLike(serviceData.id, userId);
        setState(prev => ({ ...prev, liked: true, disliked: false }));
    }, [serviceData.id, userId, handleLike]);
    
    const handleDislikeClick = useCallback(() => {
        if (!userId) {
            return;
        }
        handleDislike(serviceData.id, userId);
        setState(prev => ({ ...prev, liked: false, disliked: true }));
    }, [serviceData.id, userId, handleDislike]);
    
    const handleShare = useCallback(() => {
        if (navigator.share && typeof window !== 'undefined') {
            navigator.share({
                title: serviceData.name,
                text: serviceData.description,
                url: window.location.href
            }).catch(() => {});
        } else {
            navigator.clipboard.writeText(window.location.href);
        }
    }, [serviceData.name, serviceData.description]);
    
    const handleContactProvider = useCallback(() => {
        dispatch(openMessageModal({
            itemId: serviceData.id,
            recipientId: serviceData.provider.name
        }));
    }, [serviceData.id, serviceData.provider.name, dispatch]);
    
    return (
        <div className={styles.container}>
            <div className={styles.detailView}>
                {/* Breadcrumb Navigation */}
                <nav className={styles.breadcrumb}>
                    <Link href="/services" className={styles.breadcrumbLink}>Services</Link>
                    <ChevronRight size={16} />
                    <Link href={`/services/${serviceData.categorySlug}`} className={styles.breadcrumbLink}>
                        {category?.name || serviceData.categorySlug}
                    </Link>
                    <ChevronRight size={16} />
                    <span className={styles.breadcrumbCurrent}>{serviceData.name}</span>
                </nav>
                
                {/* Main Content Grid */}
                <div className={styles.mainGrid}>
                    {/* Full Width Hero Image Section */}
                    <div className={styles.imageSection}>
                        {/* Main Image Display */}
                        <div className={styles.mainImageContainer}>
                            {isLoadingMedia ? (
                                <div className={styles.imagePlaceholder}>
                                    <div className={styles.spinner} />
                                </div>
                            ) : mediaItems.length > 0 ? (
                                <>
                                    <div className={styles.mainImage}>
                                        <img
                                            src={mediaItems[state.currentImageIndex]?.url}
                                            alt={mediaItems[state.currentImageIndex]?.alt || serviceData.name}
                                            onClick={() => setState(prev => ({ ...prev, zoomedImage: mediaItems[state.currentImageIndex]?.url }))}
                                            className={styles.serviceImage}
                                        />
                                        
                                        {/* Navigation Arrows */}
                                        {mediaItems.length > 1 && (
                                            <>
                                                <button
                                                    className={`${styles.imageNav} ${styles.prevNav}`}
                                                    onClick={() => handleImageNavigation('prev')}
                                                    aria-label="Previous image"
                                                >
                                                    <ChevronLeft size={24} />
                                                </button>
                                                <button
                                                    className={`${styles.imageNav} ${styles.nextNav}`}
                                                    onClick={() => handleImageNavigation('next')}
                                                    aria-label="Next image"
                                                >
                                                    <ChevronRight size={24} />
                                                </button>
                                            </>
                                        )}
                                    </div>
                                    
                                    {/* Image Counter */}
                                    <div className={styles.imageCounter}>
                                        {state.currentImageIndex + 1} / {mediaItems.length}
                                    </div>
                                    
                                    {/* Thumbnail Gallery - Inside Image Container */}
                                    {mediaItems.length > 1 && (
                                        <div className={styles.thumbnailGallery}>
                                            {mediaItems.map((media, index) => (
                                                <button
                                                    key={index}
                                                    className={`${styles.thumbnail} ${index === state.currentImageIndex ? styles.activeThumbnail : ''}`}
                                                    onClick={() => setState(prev => ({ ...prev, currentImageIndex: index }))}
                                                    aria-label={`View image ${index + 1}`}
                                                >
                                                    <img src={media.url} alt={`Service image ${index + 1}`} />
                                                </button>
                                            ))}
                                        </div>
                                    )}
                                </>
                            ) : (
                                <div className={styles.noImage}>
                                    <Camera size={64} />
                                    <p>No images available</p>
                                </div>
                            )}
                        </div>
                    </div>
                    
                    {/* Content Section with Grid Layout */}
                    <div className={styles.infoSection}>
                        {/* Service Header Card */}
                        <div className={styles.serviceHeader}>
                            <h1 className={styles.serviceTitle}>{serviceData.name}</h1>
                            <div className={styles.serviceType}>
                                <Briefcase size={16} />
                                <span>{serviceData.serviceType}</span>
                            </div>
                            
                            {/* Rating and Reviews */}
                            <div className={styles.ratingSection}>
                                <StarRating rating={serviceData.rating} reviewCount={serviceData.reviewCount} />
                                <span className={styles.separator}>•</span>
                                <span className={styles.completedJobs}>
                                    <CheckCircle size={16} />
                                    {serviceData.completedJobs} completed
                                </span>
                                <span className={styles.separator}>•</span>
                                <span className={styles.viewCount}>{serviceData.metrics.views} views</span>
                            </div>
                        </div>
                        
                        {/* Price Section - Sticky Sidebar */}
                        <div className={styles.priceSection}>
                            <div className={styles.priceMain}>
                                <span className={styles.price}>
                                    {serviceData.pricingModel === 'hourly' 
                                        ? `${formattedHourlyRate}/hr` 
                                        : formattedBasePrice}
                                </span>
                                {serviceData.negotiable && (
                                    <span className={styles.negotiableBadge}>
                                        <Tag size={16} />
                                        Negotiable
                                    </span>
                                )}
                            </div>
                            <div className={styles.priceInfo}>
                                <span className={styles.pricingModel}>
                                    {serviceData.pricingModel === 'hourly' ? 'Hourly Rate' : 
                                     serviceData.pricingModel === 'fixed' ? 'Fixed Price' : 'Quote on Request'}
                                </span>
                            </div>
                            
                            {/* Booking Section - Moved inside price card */}
                            <div className={styles.bookingSection}>
                                <div className={styles.actionButtons}>
                                    <button
                                        className={styles.bookServiceBtn}
                                        onClick={handleBookService}
                                    >
                                        <Calendar size={20} />
                                        <span>Book Service</span>
                                    </button>
                                    
                                    <button
                                        className={`${styles.favoriteBtn} ${state.favorite ? styles.active : ''}`}
                                        onClick={() => setState(prev => ({ ...prev, favorite: !prev.favorite }))}
                                        aria-label="Save service"
                                    >
                                        <Heart size={20} />
                                    </button>
                                </div>
                                
                                <div className={styles.secondaryActions}>
                                    <button onClick={handleContactProvider} className={styles.secondaryBtn}>
                                        <MessageCircle size={18} />
                                        Message Provider
                                    </button>
                                    <button onClick={handleShare} className={styles.secondaryBtn}>
                                        <Share2 size={18} />
                                        Share
                                    </button>
                                </div>
                            </div>
                        </div>
                        
                        {/* Quick Info Grid */}
                        <div className={styles.quickInfo}>
                            <div className={styles.infoItem}>
                                <MapPin size={20} />
                                <div>
                                    <span className={styles.infoLabel}>Service Area</span>
                                    <span className={styles.infoValue}>{serviceData.location} • {serviceData.serviceArea}</span>
                                </div>
                            </div>
                            <div className={styles.infoItem}>
                                <Timer size={20} />
                                <div>
                                    <span className={styles.infoLabel}>Response Time</span>
                                    <span className={styles.infoValue}>{serviceData.responseTime}</span>
                                </div>
                            </div>
                            <div className={styles.infoItem}>
                                <Calendar size={20} />
                                <div>
                                    <span className={styles.infoLabel}>Availability</span>
                                    <span className={styles.infoValue}>{serviceData.availability}</span>
                                </div>
                            </div>
                            <div className={styles.infoItem}>
                                <Globe size={20} />
                                <div>
                                    <span className={styles.infoLabel}>Languages</span>
                                    <span className={styles.infoValue}>{serviceData.languages.join(', ')}</span>
                                </div>
                            </div>
                        </div>
                        
                        {/* Service Badges */}
                        <div className={styles.serviceBadges}>
                            {serviceData.verified && (
                                <span className={styles.badge}>
                                    <Shield size={16} />
                                    Verified
                                </span>
                            )}
                            {serviceData.instantBooking && (
                                <span className={styles.badge}>
                                    <Zap size={16} />
                                    Instant Booking
                                </span>
                            )}
                            {serviceData.onlineService && (
                                <span className={styles.badge}>
                                    <Globe size={16} />
                                    Online Service
                                </span>
                            )}
                            {serviceData.emergencyService && (
                                <span className={styles.badge}>
                                    <AlertCircle size={16} />
                                    Emergency Available
                                </span>
                            )}
                        </div>
                        
                        {/* Provider Info Box */}
                        <div className={styles.providerBox}>
                            <h3>Service Provider</h3>
                            <div className={styles.providerInfo}>
                                <div className={styles.providerHeader}>
                                    <div className={styles.providerAvatar}>
                                        <User size={32} />
                                    </div>
                                    <div className={styles.providerDetails}>
                                        <h4>{serviceData.provider.name}</h4>
                                        <div className={styles.providerStats}>
                                            <div className={styles.providerRating}>
                                                <Star size={16} className={styles.starFilled} />
                                                <span>{serviceData.provider.rating}</span>
                                            </div>
                                            <span>•</span>
                                            <span>{serviceData.provider.totalJobs} jobs</span>
                                            <span>•</span>
                                            <span>Since {serviceData.provider.memberSince}</span>
                                        </div>
                                    </div>
                                </div>
                                
                                {/* Verification Badges */}
                                <div className={styles.verificationBadges}>
                                    {serviceData.provider.verificationBadges.map((badge, index) => (
                                        <span key={index} className={styles.verifyBadge}>
                                            <CheckCircle size={14} />
                                            {badge}
                                        </span>
                                    ))}
                                </div>
                                
                                <div className={styles.providerActions}>
                                    <span className={styles.responseTime}>
                                        Usually responds within {serviceData.provider.responseTime}
                                    </span>
                                    <Link href={`/providers/${serviceData.provider.name}`} className={styles.viewProfileBtn}>
                                        View Profile
                                        <ArrowRight size={16} />
                                    </Link>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                
                {/* Bottom Section - Tabs */}
                <div className={styles.bottomSection}>
                    {/* Tab Navigation */}
                    <div className={styles.tabNav}>
                        <button
                            className={`${styles.tabBtn} ${state.activeTab === 'description' ? styles.activeTab : ''}`}
                            onClick={() => setState(prev => ({ ...prev, activeTab: 'description' }))}
                        >
                            <Info size={18} />
                            Description
                        </button>
                        <button
                            className={`${styles.tabBtn} ${state.activeTab === 'included' ? styles.activeTab : ''}`}
                            onClick={() => setState(prev => ({ ...prev, activeTab: 'included' }))}
                        >
                            <CheckCircle size={18} />
                            What's Included
                        </button>
                        <button
                            className={`${styles.tabBtn} ${state.activeTab === 'reviews' ? styles.activeTab : ''}`}
                            onClick={() => setState(prev => ({ ...prev, activeTab: 'reviews' }))}
                            id="reviews"
                        >
                            <Star size={18} />
                            Reviews ({serviceData.reviewCount})
                        </button>
                        <button
                            className={`${styles.tabBtn} ${state.activeTab === 'terms' ? styles.activeTab : ''}`}
                            onClick={() => setState(prev => ({ ...prev, activeTab: 'terms' }))}
                        >
                            <Shield size={18} />
                            Terms & Policies
                        </button>
                    </div>
                    
                    {/* Tab Content */}
                    <div className={styles.tabContent}>
                        {/* Description Tab */}
                        {state.activeTab === 'description' && (
                            <div className={styles.descriptionTab}>
                                <p className={styles.description}>
                                    {serviceData.description}
                                </p>
                                
                                {serviceData.features.length > 0 && (
                                    <div className={styles.features}>
                                        <h3>Service Highlights</h3>
                                        <ul>
                                            {serviceData.features.map((feature, index) => (
                                                <li key={index}>
                                                    <CheckCircle size={16} />
                                                    {feature}
                                                </li>
                                            ))}
                                        </ul>
                                    </div>
                                )}
                                
                                {serviceData.provider.specializations.length > 0 && (
                                    <div className={styles.specializations}>
                                        <h3>Specializations</h3>
                                        <div className={styles.specializationList}>
                                            {serviceData.provider.specializations.map((spec, index) => (
                                                <span key={index} className={styles.specializationTag}>
                                                    <Target size={14} />
                                                    {spec}
                                                </span>
                                            ))}
                                        </div>
                                    </div>
                                )}
                            </div>
                        )}
                        
                        {/* What's Included Tab */}
                        {state.activeTab === 'included' && (
                            <div className={styles.includedTab}>
                                <h3>This service includes:</h3>
                                <ul className={styles.includedList}>
                                    {serviceData.included.map((item, index) => (
                                        <li key={index}>
                                            <Check size={20} />
                                            <span>{item}</span>
                                        </li>
                                    ))}
                                </ul>
                                
                                <div className={styles.duration}>
                                    <h4>Service Duration</h4>
                                    <p>{serviceData.duration}</p>
                                </div>
                            </div>
                        )}
                        
                        {/* Reviews Tab */}
                        {state.activeTab === 'reviews' && (
                            <div className={styles.reviewsTab}>
                                <div className={styles.reviewsSummary}>
                                    <div className={styles.overallRating}>
                                        <h2>{serviceData.rating.toFixed(1)}</h2>
                                        <StarRating rating={serviceData.rating} reviewCount={0} />
                                        <p>Based on {serviceData.reviewCount} reviews</p>
                                    </div>
                                    
                                    <div className={styles.ratingBars}>
                                        {[5, 4, 3, 2, 1].map(stars => (
                                            <div key={stars} className={styles.ratingBar}>
                                                <span>{stars}</span>
                                                <Star size={14} className={styles.starFilled} />
                                                <div className={styles.barContainer}>
                                                    <div 
                                                        className={styles.barFill} 
                                                        style={{ width: `${Math.random() * 100}%` }}
                                                    />
                                                </div>
                                                <span>{Math.floor(Math.random() * 100)}%</span>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                                
                                <button className={styles.writeReviewBtn}>
                                    Write a Review
                                </button>
                                
                                {/* Comments Component */}
                                <div className={styles.reviewsList}>
                                    <CommentsSetup
                                        userId={userId}
                                        itemId={serviceData.id}
                                        itemType="service"
                                        toggleCommentsList={() => {}}
                                        categoryId={serviceData.categoryId}
                                    />
                                </div>
                            </div>
                        )}
                        
                        {/* Terms Tab */}
                        {state.activeTab === 'terms' && (
                            <div className={styles.termsTab}>
                                <div className={styles.termsSection}>
                                    <h3>Cancellation Policy</h3>
                                    <p>{serviceData.cancellationPolicy}</p>
                                </div>
                                
                                <div className={styles.termsSection}>
                                    <h3>Payment Terms</h3>
                                    <p>{serviceData.paymentTerms}</p>
                                </div>
                                
                                <div className={styles.termsSection}>
                                    <h3>Service Warranty</h3>
                                    <p>{serviceData.warranty}</p>
                                </div>
                                
                                <div className={styles.termsSection}>
                                    <h3>Important Information</h3>
                                    <ul>
                                        <li>
                                            <Info size={18} />
                                            <span>All prices include VAT</span>
                                        </li>
                                        <li>
                                            <CreditCard size={18} />
                                            <span>Secure payment processing</span>
                                        </li>
                                        <li>
                                            <Shield size={18} />
                                            <span>Service satisfaction guaranteed</span>
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
                
                {/* Engagement Bar */}
                <div className={styles.engagementBar}>
                    <div className={styles.engagementStats}>
                        <span><Eye size={18} /> {serviceData.metrics.views} views</span>
                        <span><Heart size={18} /> {serviceData.metrics.saved} saved</span>
                        <span><Calendar size={18} /> {serviceData.metrics.bookings} bookings</span>
                    </div>
                    
                    <div className={styles.engagementActions}>
                        <button
                            className={`${styles.engageBtn} ${state.liked ? styles.active : ''}`}
                            onClick={handleLikeClick}
                        >
                            <ThumbsUp size={18} />
                            <span>{serviceData.metrics.likes}</span>
                        </button>
                        <button
                            className={`${styles.engageBtn} ${state.disliked ? styles.active : ''}`}
                            onClick={handleDislikeClick}
                        >
                            <ThumbsDown size={18} />
                            <span>{serviceData.metrics.dislikes}</span>
                        </button>
                        <button
                            className={styles.engageBtn}
                            onClick={() => setState(prev => ({ ...prev, showComments: !prev.showComments }))}
                        >
                            <MessageCircle size={18} />
                            <span>{serviceData.metrics.comments}</span>
                        </button>
                        <button className={styles.engageBtn} onClick={handleShare}>
                            <Share2 size={18} />
                            <span>Share</span>
                        </button>
                    </div>
                </div>
                
                {/* Tags */}
                {serviceData.tags.length > 0 && (
                    <div className={styles.tagsSection}>
                        <h3>Related Tags</h3>
                        <div className={styles.tagsList}>
                            {serviceData.tags.map((tag, index) => (
                                <Link key={index} href={`/search?tag=${encodeURIComponent(tag)}`} className={styles.tag}>
                                    <Tag size={14} />
                                    {tag}
                                </Link>
                            ))}
                        </div>
                    </div>
                )}
                
                {/* Image Zoom Modal */}
                {state.zoomedImage && (
                    <div className={styles.zoomModal} onClick={() => setState(prev => ({ ...prev, zoomedImage: null }))}>
                        <img src={state.zoomedImage} alt="Zoomed service image" />
                        <button className={styles.closeZoom} aria-label="Close zoom">×</button>
                    </div>
                )}
            </div>
        </div>
    );
});

DetailedServiceView.displayName = 'DetailedServiceView';

export default DetailedServiceView;