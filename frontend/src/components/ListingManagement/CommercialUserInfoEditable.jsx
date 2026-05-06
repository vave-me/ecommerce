import { Briefcase, ShoppingCart } from '@/icons';
import { FaHandshake } from '../../utils/iconImports';
import React from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';
import {useRouter} from "next/navigation";
const colorAccent = '#007bc7';   // Primary accent color
const colorTextStrong = '#1C1C1C';
const colorHoverBg = '#f1f1f1';
const colorFocusOutline = colorAccent;
// Unified spacing scale & breakpoints
const spacing = {
    xs: '4px',
    sm: '6px',
    md: '8px',
    lg: '12px',
    xl: '16px',
};
const breakpoints = {
    mobile: '480px',
    tablet: '768px',
};
// A small set of typical font sizes
const fontSize = {
    xs: '0.75rem',    // ~12px
    sm: '0.875rem',   // ~14px
    md: '1rem',       // ~16px
    lg: '1.125rem',   // ~18px
};
function CommercialUserInfoEditable({
                                sellerId,
                                user = {name: 'Unknown Seller', avatarUrl: '/images/user-user.webp'},
                                rating = 0,
                                renderStars,
                                id,
                                productPrice,
                                canUseMiddleman,
                                handleAddToCart,
                                handleOpenOffersForm,
                            }) {
    const avatarUrl = user?.avatarUrl || '/images/user-user.webp';
    const sellerName = user?.id || 'Unknown Seller';
    const navigate = useRouter();
    const handleClick = () => {
        // navigate to the profile page with the seller ID
        navigate.push(`/profile/${sellerId}`);
    };
    return (
        <Row>
            <LeftGroup onClick={handleClick}>
                <Avatar src={avatarUrl} alt={`${sellerName} avatar`}/>
                <UserMeta>
                    <UserName>
                        {sellerName}
                        <BadgeUserType $type="commercial">
                            <Briefcase className="badge-icon"/>
                            Commercial
                        </BadgeUserType>
                    </UserName>
                    <UserRating>
                        {renderStars?.()}
                        <RatingValue>{rating.toFixed(1)}</RatingValue>
                    </UserRating>
                </UserMeta>
            </LeftGroup>
            <RightGroup>
                <PriceColumn>
                    <PriceTag>
                        {productPrice}
                    </PriceTag>
                    {canUseMiddleman && (
                        <>
                            <IconButton
                                onClick={() =>
                                    handleAddToCart?.(id)}
                                title="Add to Cart"
                                aria-label="Add item to cart"
                            >
                                <ShoppingCart/>
                            </IconButton>
                            <IconButton
                                onClick={
                                    handleOpenOffersForm}
                                title="Send Offer"
                                aria-label="Send an offer"
                            >
                                <FaHandshake/>
                            </IconButton>
                        </>
                    )}
                </PriceColumn>
            </RightGroup>
        </Row>
    );
}
CommercialUserInfoEditable.propTypes = {
    sellerId: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
    user: PropTypes.shape({
        id: PropTypes.string,
        avatarUrl: PropTypes.string,
        name: PropTypes.string,
    }),
    rating: PropTypes.number,
    productLocation: PropTypes.string,
    renderStars: PropTypes.func,
};
CommercialUserInfoEditable.defaultProps = {
    user: {
        name: 'Unknown Seller',
        avatarUrl: '/images/user-user.webp',
    },
    rating: 0,
    productLocation: 'Unknown Location',
    renderStars: null,
};
export default CommercialUserInfoEditable;
/* --------------------------------------------
   STYLES (Compact)
--------------------------------------------- */
const colorBadgeSellerPrivate = '#6C757D';
const colorBadgeSellerCommercial = '#2EA02E';
const colorAccentHover = '#0067a8';
const hoverBgColor = '#4B9EA0';
const Row = styled.div`
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 6px;
    cursor: pointer; /* entire row is clickable */
    @media (max-width: 768px) {
        gap: 5px;
    }
    @media (max-width: 480px) {
        gap: 4px;
    }
`;
const LeftGroup = styled.div`
    display: inline-flex;
    align-items: center;
    gap: 6px;
    @media (max-width: 768px) {
        gap: 5px;
    }
    @media (max-width: 480px) {
        gap: 4px;
    }
`;
const RightGroup = styled.div`
    margin-left: auto;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    @media (max-width: 768px) {
        gap: 5px;
    }
    @media (max-width: 480px) {
        gap: 4px;
    }
`;
const Avatar = styled.img`
    width: 38px;
    height: 38px;
    object-fit: cover;
    border-radius: 50%;
    transition: transform 0.2s;
    &:hover {
        transform: scale(1.05);
    }
    @media (max-width: 768px) {
        width: 36px;
        height: 36px;
    }
    @media (max-width: 480px) {
        width: 34px;
        height: 34px;
    }
`;
const UserMeta = styled.div`
    display: flex;
    flex-direction: column;
    gap: 2px;
`;
const UserName = styled.span`
    font-size: 0.85rem;
    font-weight: 600;
    color: ${colorTextStrong};
    display: inline-flex;
    align-items: center;
    gap: 5px;
    @media (max-width: 768px) {
        font-size: 0.8rem;
        gap: 4px;
    }
    @media (max-width: 480px) {
        font-size: 0.75rem;
    }
`;
const UserRating = styled.button`
    display: inline-flex;
    align-items: center;
    gap: 5px;
    background: none;
    border: none;
    padding: 4px;
    cursor: pointer;
    border-radius: 4px;
    &:hover {
        background-color: ${hoverBgColor};
    }
    .star-icon {
        font-size: 14px;
    }
    @media (max-width: 768px) {
        gap: 4px;
        padding: 3px;
        .star-icon {
            font-size: 16px;
        }
    }
    @media (max-width: 480px) {
        .star-icon {
            font-size: 16px;
        }
    }
`;
const RatingValue = styled.span`
    font-size: 0.75rem;
    color: #777;
    @media (max-width: 768px) {
        font-size: 0.7rem;
    }
    @media (max-width: 480px) {
        font-size: 0.65rem;
    }
`;
const LocationWrapper = styled.button`
    display: inline-flex;
    align-items: center;
    gap: 5px;
    font-size: 0.8rem;
    font-weight: 500;
    color: ${colorTextStrong};
    background: none;
    border: none;
    padding: 4px;
    border-radius: 4px;
    cursor: pointer;
    &:hover {
        background-color: ${hoverBgColor};
    }
    svg {
        width: 14px;
        height: 14px;
        color: ${colorAccent};
        transition: color 0.2s ease;
    }
    &:hover svg {
        color: ${colorAccentHover};
    }
    @media (max-width: 768px) {
        font-size: 0.75rem;
        gap: 4px;
    }
    @media (max-width: 480px) {
        font-size: 0.7rem;
    }
`;
/** Reused BaseBadge (both Private/Commercial) */
const BaseBadge = styled.div`
    display: inline-flex;
    align-items: center;
    gap: 2px;
    font-size: 0.6rem;
    padding: 2px 4px;
    border-radius: 12px;
    font-weight: 600;
    text-transform: uppercase;
    color: ${colorTextStrong};
    .badge-icon {
        font-size: 9px;
        @media (max-width: 768px) {
            font-size: 8px;
        }
        @media (max-width: 480px) {
            font-size: 7px;
        }
    }
    @media (max-width: 768px) {
        font-size: 0.55rem;
        padding: 2px 3px;
    }
    @media (max-width: 480px) {
        font-size: 0.5rem;
    }
`;
const PriceColumn = styled.div`
    display: inline-flex;
    align-items: center;
    gap: ${spacing.sm};
    @media (max-width: ${breakpoints.tablet}) {
        gap: ${spacing.xs};
    }
`;
const PriceTag = styled.div`
    margin: 0;
    color: #3b82f6;
    font-size: 1rem;
    font-weight: bold;
    white-space: nowrap;
    @media (max-width: ${breakpoints.tablet}) {
        font-size: ${fontSize.sm};
    }
    @media (max-width: ${breakpoints.mobile}) {
        font-size: 0.8rem;
    }
`;
const IconButton = styled.button`
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: 1px solid transparent;
    border-radius: 6px;
    color: ${colorTextStrong};
    font-size: ${fontSize.sm};
    padding: ${spacing.sm};
    cursor: pointer;
    transition: background-color 0.2s ease, color 0.2s ease, border-color 0.2s ease, transform 0.2s ease;
    svg {
        width: 20px;
        height: 20px;
    }
    &:hover {
        background-color: ${colorHoverBg};
        border-color: ${colorAccent};
        color: ${colorAccent};
    }
    &:focus {
        outline: 2px solid ${colorFocusOutline};
        outline-offset: 2px;
    }
    &:active {
        transform: scale(0.96);
        background-color: ${colorAccent};
        color: #fff;
        border-color: ${colorAccent};
    }
    @media (max-width: ${breakpoints.tablet}) {
        padding: ${spacing.md};
        svg {
            width: 22px;
            height: 22px;
        }
    }
    @media (max-width: ${breakpoints.mobile}) {
        padding: ${spacing.lg};
        svg {
            width: 24px;
            height: 24px;
        }
    }
`;
export const BadgeUserType = styled(BaseBadge)`
    background-color: ${({$type}) =>
            $type === 'private'
                    ? `${colorBadgeSellerPrivate}22`
                    : `${colorBadgeSellerCommercial}22`};
`;
