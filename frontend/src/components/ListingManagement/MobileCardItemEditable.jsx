// File: src/components/MobileItemCard/MobileCardItemEditable.jsx
"use client"
import { Edit, Star } from '@/icons';
import React, { useState, useEffect, useMemo, useRef, useCallback, memo } from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import {toast} from 'react-toastify';
import {useDispatch} from 'react-redux';
import {openCommentsFullModal, openMessageModal} from '../../redux/slices/modalsSlice';
import {useAuth} from '../../context/AuthContext';
import useActivityApi from '../../hooks/useActivityApi';
import { useProductBasket } from '../../hooks/useProductBasket';
import useWishlist from '../../hooks/useWishlist';
import {getMediaByItem, isMediaResponseSuccess} from '../../api/mediaApi'; // removed getAllMediaImages, getAllMediaVideos
import {getVariants} from '../../api/postsApi';
import InteractionsMobile from '../../features/ProductCard/InteractionsMobile';
import TimeAgoBadgesRowEditable from './TimeAgoBadgeRowEditable';
import PrivateUserInfoEditable from './PrivateUserInfoEditable';
import CommercialUserInfoEditable from './CommercialUserInfoEditable';
import TitleRowEditable from './TitleRowEditable';
import ProductInfoEditable from './ProductInfoEditable';
import OffersModalFormEditable from './OffersModalFormEditable';
import DescriptionMobileEditable from './DescriptionMobileEditable';
import MediaGalleryMobileEditable from './MediaGalleryMobileEditable';
dayjs.extend(relativeTime);
/* ------------------------------------------------------------------
   Helper: Toast notifications
------------------------------------------------------------------ */
function useToastNotification() {
    return (type, message) => {
        const options = {theme: 'colored'};
        switch (type) {
            case 'success':
                toast.success(message, options);
                break;
            case 'info':
                toast.info(message, options);
                break;
            case 'error':
                toast.error(message, options);
                break;
            case 'warn':
                toast.warn(message, options);
                break;
            default:
                toast(message, options);
        }
    };
}
/* ------------------------------------------------------------------
   Main Component: MobileCardItemEditable
------------------------------------------------------------------ */
const MobileCardItemEditable = memo(({product}) => {
    // 1) Hooks: Auth, Redux, Toast
    const {user} = useAuth();
    const userId = user?.userId;
    const dispatch = useDispatch();
    const showToast = useToastNotification();
    const [fetchedMediaId, setFetchedMediaId] = useState(null);
    // 2) Activity / Basket / Wishlist
    const {handleLike, handleDislike} = useActivityApi();
    const {addToBasket: handleAddToCart} = useProductBasket(product?.id || 'fallback-id', product);
    const {items: wishlistItems, toggleItem: toggleWishlist} = useWishlist();
    // 3) Local ephemeral states
    const [currentMediaIndex, setCurrentMediaIndex] = useState(0);
    const [activeVideoIndex, setActiveVideoIndex] = useState(null);
    const [finalGallery, setFinalGallery] = useState([]);
    const [computedThumbnail, setComputedThumbnail] = useState(
        product?.thumbnail || '/images/palit.png'
    );
    const [fetchedVariants, setFetchedVariants] = useState([]);
    // Inline editing states
    const [title, setTitle] = useState(product?.name || 'Untitled Product');
    const [editingTitle, setEditingTitle] = useState(false);
    const [description, setDescription] = useState(
        product?.description || 'No description available.'
    );
    const [editingDescription, setEditingDescription] = useState(false);
    const [basePrice, setBasePrice] = useState(product?.basePrice || '999');
    const [editingPrice, setEditingPrice] = useState(false);
    const [condition, setCondition] = useState(product?.condition || 'new');
    const [editingCondition, setEditingCondition] = useState(false);
    // Track mounting state
    const isMountedRef = useRef(true);
    useEffect(() => {
        return () => {
            isMountedRef.current = false;
        };
    }, []);
    /* ------------------------------------------------------------------
       4) Media Order Fetch (replacing getAllMediaImages / getAllMediaVideos)
    ------------------------------------------------------------------ */
    useEffect(() => {
        if (!product.id) {
            return;
        }
        (async () => {
            try {
                const mediaResp = await getMediaByItem(product.id);
                // Check if the media API response was successful
                if (!isMediaResponseSuccess(mediaResp)) {
                    return;
                }
                const media = mediaResp?.media;
                if (!media) {
                    return;
                }
                setFetchedMediaId(media.id);
                // mediaOrder is an array:
                // [ {mediaItemId, url}, ... ]
                const orderArray = media.mediaOrder || [];
                // Transform each item in mediaOrder to finalGallery
                const combined = orderArray.map((mItem, index) => {
                    // Guess type by extension
                    let type = 'image';
                    const lowerUrl = (mItem.url || '').toLowerCase();
                    if (lowerUrl.endsWith('.mp4') || lowerUrl.includes('/video/')) {
                        type = 'video';
                    }
                    return {
                        type,
                        src: mItem.url,
                        alt: `Media ${mItem.mediaItemId}`,
                        displayOrder: index, // or actual displayOrder if available
                    };
                });
                setFinalGallery(combined);
                // If there's at least one image, set computedThumbnail
                const firstImage = combined.find(item => item.type === 'image');
                if (firstImage) {
                    setComputedThumbnail(firstImage.src);
                }
            } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }
        })();
    }, [product.id]);
    /* ------------------------------------------------------------------
       5) Merge fallback values => final product data
    ------------------------------------------------------------------ */
    const productWithDefaults = useMemo(() => {
        return {
            ...product,
            id: product?.id || 'fallback-id',
            name: typeof title === 'string' ? title : String(title),
            description: typeof description === 'string' ? description : String(description),
            basePrice: typeof basePrice === 'string' ? basePrice : String(basePrice),
            condition: typeof condition === 'string' ? condition : String(condition),
            added: product?.added || new Date().toISOString(),
            thumbnail: computedThumbnail,
            userType: product?.userType || 'commercial',
            location: product?.location || 'Berlin, 10365',
            likeCount: product?.likeCount ?? 45,
            dislikeCount: product?.dislikeCount ?? 3,
            rating: product?.rating ?? 4.2,
            wishlistCount: product?.wishlistCount ?? 128,
            messages: product?.messages ?? 42,
            hasVariants: product?.hasVariants ?? false,
            shipping: product?.shipping ?? [],
            stock:
                typeof product?.stock === 'number' && product.stock > 0
                    ? Array.from({length: product.stock}, (_, i) => i + 1)
                    : [1, 2, 3],
            productOptions: product?.options ?? [],
            category: product?.category || {},
            lat: 52.514835, // fallback lat/lng
            lng: 13.5055275,
        };
    }, [product, title, description, basePrice, condition, computedThumbnail]);
    const {
        id,
        name: finalName,
        added,
        userType,
        location: productLocation,
        likeCount,
        dislikeCount,
        rating,
        wishlistCount,
        messages,
        middlemanService,
        hasVariants,
        shipping,
        stock,
        productOptions,
    } = productWithDefaults;
    /* ------------------------------------------------------------------
       6) Fetch variants if needed
    ------------------------------------------------------------------ */
    useEffect(() => {
        if (!hasVariants || !id) {
            setFetchedVariants([]);
            return;
        }
        (async () => {
            try {
                const variantData = await getVariants(id);
                if (isMountedRef.current && variantData?.variants) {
                    setFetchedVariants(variantData.variants);
                }
            } catch (err) {
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', err);
        }
        throw err; // Re-throw for caller to handle
    }
        })();
    }, [hasVariants, id]);
    /* ------------------------------------------------------------------
       7) Basic Handlers: Add to Cart, Like/Dislike, Wishlist
    ------------------------------------------------------------------ */
    const canUseMiddleman = !!middlemanService;
    const handleAddToCartInternal = () => {
        handleAddToCart();
    };
    const handleLikeClick = () => {
        if (!userId) {
            showToast('warn', 'Please log in to like products.');
            return;
        }
        handleLike(id, userId).catch(() => {/* handle error */
        });
    };
    const handleDislikeClick = () => {
        if (!userId) {
            showToast('warn', 'Please log in to dislike products.');
            return;
        }
        handleDislike(id, userId).catch(() => {/* handle error */
        });
    };
    const isInWishlist = useMemo(
        () => wishlistItems.some((wlItem) => wlItem.productId === id),
        [wishlistItems, id]
    );
    const handleWishlistClick = () => {
        if (!userId) {
            showToast('warn', 'Please log in to manage your wishlist.');
            return;
        }
        toggleWishlist(id).catch(() => {/* handle error */
        });
    };
    // Swipe handlers for gallery
    const handleNextMedia = useCallback(() => {
        setCurrentMediaIndex((prev) => (prev + 1) % finalGallery.length);
        setActiveVideoIndex(null);
    }, [finalGallery.length]);
    const handlePrevMedia = useCallback(() => {
        setCurrentMediaIndex((prev) => (prev - 1 + finalGallery.length) % finalGallery.length);
        setActiveVideoIndex(null);
    }, [finalGallery.length]);
    // Comments & Messaging
    const handleOpenComments = () => {
        dispatch(
            openCommentsFullModal({
                itemId: id,
                itemType: 'product',
                categoryId: product?.categoryId,
            })
        );
    };
    const handleOpenMessage = () => {
        dispatch(
            openMessageModal({
                productId: id,
                senderId: userId,
                recipientId: product?.userSellerId,
            })
        );
    };
    // Star ratings
    const renderStars = useCallback(() => {
        const fullStars = Math.floor(rating);
        const halfStar = rating - fullStars >= 0.5;
        const starsArray = [];
        for (let i = 0; i < fullStars; i += 1) {
            starsArray.push(<Star key={`star-full-${i}`} color="#FFC107"/>);
        }
        if (halfStar) {
            starsArray.push(<Star key="star-half" color="#FFC107" style={{opacity: 0.6}}/>);
        }
        const emptyStarsCount = 5 - fullStars - (halfStar ? 1 : 0);
        for (let j = 0; j < emptyStarsCount; j += 1) {
            starsArray.push(<Star key={`star-empty-${j}`} color="#ccc" style={{opacity: 0.4}}/>);
        }
        return starsArray;
    }, [rating]);
    // Offers form
    const [isOffersFormOpen, setIsOffersFormOpen] = useState(false);
    const handleOpenOffersForm = () => setIsOffersFormOpen(true);
    const handleCloseOffersForm = () => setIsOffersFormOpen(false);
    // Transform variants
    const transformVariants = useCallback((variantsFromAPI) => {
        return variantsFromAPI.map((variant) => {
            const {Storage, Color, VRAM} = variant.Attributes || {};
            const details = [];
            if (Storage) details.push(`Storage: ${Storage}`);
            if (Color) details.push(`Color: ${Color}`);
            if (VRAM) details.push(`VRAM: ${VRAM}`);
            return {
                id: variant.ID,
                name: variant.SKU,
                details,
            };
        });
    }, []);
    // Inline editing: Title, Description, Price, Condition
    const startEditTitle = useCallback(() => setEditingTitle(true), []);
    const cancelEditTitle = useCallback(() => {
        setTitle(product?.name || 'Untitled Product');
        setEditingTitle(false);
    }, [product?.name]);
    const saveEditTitle = useCallback(() => {
        showToast('success', 'Title updated!');
        setEditingTitle(false);
    }, [showToast]);
    const startEditDescription = useCallback(() => setEditingDescription(true), []);
    const cancelEditDescription = useCallback(() => {
        setDescription(product?.description || 'No description available.');
        setEditingDescription(false);
    }, [product?.description]);
    const saveEditDescription = useCallback(() => {
        showToast('success', 'Description updated!');
        setEditingDescription(false);
    }, [showToast]);
    const startEditPrice = useCallback(() => setEditingPrice(true), []);
    const cancelEditPrice = useCallback(() => {
        setBasePrice(product?.basePrice || '999');
        setEditingPrice(false);
    }, [product?.basePrice]);
    const saveEditPrice = useCallback(() => {
        showToast('success', 'Price updated!');
        setEditingPrice(false);
    }, [showToast]);
    const startEditCondition = useCallback(() => setEditingCondition(true), []);
    const cancelEditCondition = useCallback(() => {
        setCondition(product?.condition || 'new');
        setEditingCondition(false);
    }, [product?.condition]);
    const saveEditCondition = useCallback(() => {
        showToast('success', 'Condition updated!');
        setEditingCondition(false);
    }, [showToast]);
    /* ------------------------------------------------------------------
       Render
    ------------------------------------------------------------------ */
    return (
        <CardContainer>
            {/* Time & Location badges (row at top) */}
            <TimeAgoBadgesRowEditable
                approximateLat={productWithDefaults.lat}
                approximateLon={productWithDefaults.lng}
                timeAgo={dayjs(productWithDefaults.added).fromNow()}
            />
            {/* Seller info & Price row */}
            {userType === 'private' ? (
                <PrivateUserInfoEditable
                    productPrice={
                        editingPrice ? (
                            <EditInlineRow>
                                <EditInput
                                    type="number"
                                    min="0"
                                    step="0.01"
                                    value={basePrice}
                                    onChange={(e) => setBasePrice(e.target.value)}
                                />
                                <ButtonRow>
                                    <SmallButton onClick={saveEditPrice}>Save</SmallButton>
                                    <SmallButtonGray onClick={cancelEditPrice}>Cancel</SmallButtonGray>
                                </ButtonRow>
                            </EditInlineRow>
                        ) : (
                            <InlineDisplay onClick={startEditPrice}>
                                {`${basePrice} €`} <Edit className="icon"/>
                            </InlineDisplay>
                        )
                    }
                    handleAddToCart={handleAddToCartInternal}
                    handleOpenOffersForm={handleOpenOffersForm}
                    canUseMiddleman={canUseMiddleman}
                    sellerId={product?.userSellerId}
                    user={user}
                    rating={rating}
                    productLocation={productLocation}
                    renderStars={renderStars}
                />
            ) : (
                <CommercialUserInfoEditable
                    productPrice={
                        editingPrice ? (
                            <EditInlineRow>
                                <EditInput
                                    type="number"
                                    min="0"
                                    step="0.01"
                                    value={basePrice}
                                    onChange={(e) => setBasePrice(e.target.value)}
                                />
                                <ButtonRow>
                                    <SmallButton onClick={saveEditPrice}>Save</SmallButton>
                                    <SmallButtonGray onClick={cancelEditPrice}>Cancel</SmallButtonGray>
                                </ButtonRow>
                            </EditInlineRow>
                        ) : (
                            <InlineDisplay onClick={startEditPrice}>
                                {`${basePrice} €`} <Edit className="icon"/>
                            </InlineDisplay>
                        )
                    }
                    handleAddToCart={handleAddToCartInternal}
                    handleOpenOffersForm={handleOpenOffersForm}
                    canUseMiddleman={canUseMiddleman}
                    sellerId={product?.userSellerId}
                    user={user}
                    rating={rating}
                    productLocation={productLocation}
                    renderStars={renderStars}
                />
            )}
            {/* Title + Condition inline editing */}
            <TitleRowEditable
                id={id}
                condition={
                    editingCondition ? (
                        <EditInlineRow>
                            <EditSelect
                                value={condition}
                                onChange={(e) => setCondition(e.target.value)}
                            >
                                <option value="new">New</option>
                                <option value="used">Used</option>
                                <option value="broken">Broken</option>
                                <option value="refurbished">Refurbished</option>
                            </EditSelect>
                            <ButtonRow>
                                <SmallButton onClick={saveEditCondition}>Save</SmallButton>
                                <SmallButtonGray onClick={cancelEditCondition}>Cancel</SmallButtonGray>
                            </ButtonRow>
                        </EditInlineRow>
                    ) : (
                        <InlineDisplay onClick={startEditCondition}>
                            {condition} <Edit className="icon"/>
                        </InlineDisplay>
                    )
                }
                productName={
                    editingTitle ? (
                        <EditInlineRow>
                            <EditInput
                                value={finalName}
                                onChange={(e) => setTitle(e.target.value)}
                            />
                            <ButtonRow>
                                <SmallButton onClick={saveEditTitle}>Save</SmallButton>
                                <SmallButtonGray onClick={cancelEditTitle}>Cancel</SmallButtonGray>
                            </ButtonRow>
                        </EditInlineRow>
                    ) : (
                        <InlineDisplay onClick={startEditTitle}>
                            {finalName} <Edit className="icon"/>
                        </InlineDisplay>
                    )
                }
                productPrice={basePrice}
            />
            {/* Shipping / Variants / Product Config */}
            <ProductInfoEditable
                category={productWithDefaults.category}
                shippingOptions={shipping}
                selectedShipping={shipping?.[0]?.code || ''}
                onShippingChange={() => {/* no-op */
                }}
                hasVariants={hasVariants}
                variants={
                    fetchedVariants.length > 0
                        ? transformVariants(fetchedVariants)
                        : productWithDefaults.variants
                }
                selectedVariantId=""
                onVariantChange={() => {/* no-op */
                }}
                stockOptions={stock}
                selectedstock={stock?.[0] || 1}
                onstockChange={() => {/* no-op */
                }}
                productOptions={productOptions}
                selectedProductOptions={{}}
                onProductOptionChange={() => {/* no-op */
                }}
                onQuantityChange={() => {/* no-op */
                }}
            />
            {/* Offers (BuyNow / Lease / Pawn) */}
            {isOffersFormOpen && (
                <OffersModalFormEditable productId={id} onClose={handleCloseOffersForm}/>
            )}
            {/* Description (inline editing) */}
            <DescriptionMobileEditable
                description={
                    editingDescription ? (
                        <EditInlineColumn>
                            <EditTextArea
                                rows={4}
                                value={description}
                                onChange={(e) => setDescription(e.target.value)}
                            />
                            <ButtonRow>
                                <SmallButton onClick={saveEditDescription}>Save</SmallButton>
                                <SmallButtonGray onClick={cancelEditDescription}>Cancel</SmallButtonGray>
                            </ButtonRow>
                        </EditInlineColumn>
                    ) : (
                        <DescDisplay onClick={startEditDescription}>
                            {description}
                            <Edit className="icon"/>
                        </DescDisplay>
                    )
                }
                productTitle={finalName}
            />
            {/* Media gallery (images + videos), read-only with arrow nav */}
            <MediaGalleryMobileEditable
                mediaId={fetchedMediaId}
                gallery={finalGallery}
                currentMediaIndex={currentMediaIndex}
                setCurrentMediaIndex={setCurrentMediaIndex}
                handleNext={handleNextMedia}
                handlePrev={handlePrevMedia}
                productTitle={finalName}
                activeVideoIndex={activeVideoIndex}
                setActiveVideoIndex={setActiveVideoIndex}
            />
            {/* Interactions: like, wishlist, comments, messaging, etc. */}
            <InteractionsMobile
                itemId={id || 'fallback-id'}
                handleLike={handleLikeClick}
                handleDislike={handleDislikeClick}
                handleAddToCart={canUseMiddleman ? handleAddToCartInternal : undefined}
                handleWishlistClick={handleWishlistClick}
                toggleCommentsList={handleOpenComments}
                toggleMessageInput={handleOpenMessage}
                counts={{
                    wishlist: wishlistCount,
                    like: likeCount,
                    dislike: dislikeCount,
                    comments: product?.comments?.length || 0,
                    message: messages || 0,
                }}
                isLoggedIn={!!userId}
                isInWishlist={isInWishlist}
            />
        </CardContainer>
    );
});
MobileCardItemEditable.propTypes = {
    product: PropTypes.shape({
        id: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        name: PropTypes.any,
        description: PropTypes.any,
        basePrice: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        condition: PropTypes.any,
        thumbnail: PropTypes.string,
        category: PropTypes.oneOfType([PropTypes.string, PropTypes.object]),
        userType: PropTypes.string,
        location: PropTypes.string,
        likeCount: PropTypes.number,
        dislikeCount: PropTypes.number,
        rating: PropTypes.number,
        wishlistCount: PropTypes.number,
        messages: PropTypes.number,
        hasVariants: PropTypes.bool,
        shipping: PropTypes.array,
        stock: PropTypes.number,
        options: PropTypes.array,
        userSellerId: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
        added: PropTypes.string,
        comments: PropTypes.array,
    }).isRequired,
};
export default MobileCardItemEditable;
/* ------------------------------------------------------------------
   Styled components & CSS
------------------------------------------------------------------ */
const CardContainer = styled.div`
    position: relative;
    border-radius: 6px;
    background: rgba(255, 255, 255, 0.6);
    padding: 0 2px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    max-height: 560px;
    @media (max-width: 768px) {
        border: none;
        border-top: 1px solid rgba(0, 0, 0, 0.2);
    }
`;
const InlineDisplay = styled.div`
    display: inline-flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
    .icon {
        font-size: 0.8rem;
        color: #999;
    }
    &:hover {
        color: #0077cc;
        .icon {
            color: #0077cc;
        }
    }
`;
const EditInlineRow = styled.div`
    display: flex;
    align-items: center;
    gap: 6px;
`;
const EditInlineColumn = styled.div`
    display: flex;
    flex-direction: column;
    gap: 6px;
`;
const EditInput = styled.input`
    border: 1px solid #aaa;
    border-radius: 4px;
    padding: 4px 8px;
    font-size: 0.9rem;
    &:focus {
        outline: none;
        border-color: #0077cc;
        background-color: #f0f8ff;
    }
`;
const EditSelect = styled.select`
    border: 1px solid #aaa;
    border-radius: 4px;
    padding: 4px 8px;
    font-size: 0.9rem;
    &:focus {
        outline: none;
        border-color: #0077cc;
        background-color: #f0f8ff;
    }
`;
const EditTextArea = styled.textarea`
    border: 1px solid #aaa;
    border-radius: 4px;
    padding: 6px 8px;
    font-size: 0.9rem;
    resize: vertical;
    &:focus {
        outline: none;
        border-color: #0077cc;
        background-color: #f0f8ff;
    }
`;
const ButtonRow = styled.div`
    display: flex;
    gap: 6px;
`;
const SmallButton = styled.button`
    background-color: #0061bb;
    color: #fff;
    border: none;
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
    &:hover {
        background-color: #004c92;
    }
`;
const SmallButtonGray = styled.button`
    background-color: #999;
    color: #fff;
    border: none;
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
    &:hover {
        background-color: #777;
    }
`;
const DescDisplay = styled.div`
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    position: relative;
    .icon {
        font-size: 0.8rem;
        color: #999;
    }
    &:hover {
        color: #0077cc;
        .icon {
            color: #0077cc;
        }
    }
`;
