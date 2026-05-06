"use client"
import React, { useMemo, useCallback, useState, useEffect, useRef, memo } from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';
import ReactDOM from 'react-dom';
import { X } from '@/icons';
/**
 * DescriptionMobileEditable:
 * - Shows a truncated string ("2 lines") with a "More info" button if it overflows.
 * - If `description` is a **React node** instead of a string, no truncation is done,
 *   and we simply render it directly.
 * - Clicking "More info" opens a modal with the full text (if it's a string) or a fallback notice.
 */
const DescriptionMobileEditable = memo(({description, productTitle}) => {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const descriptionRef = useRef(null);
    const [isOverflowed, setIsOverflowed] = useState(false);
    // Memoize modal handlers to prevent unnecessary re-renders
    const openModal = useCallback((e) => {
        e.stopPropagation();
        setIsModalOpen(true);
    }, []);
    const closeModal = useCallback((e) => {
        e.stopPropagation();
        setIsModalOpen(false);
    }, []);
    // Memoize hashtag/mention handlers
    const handleHashtagClick = useCallback((tag) => {
        // implement or remove
        // Debug log removed for production
    }, []);
    const handleMentionClick = useCallback((mention) => {
        // implement or remove
        // Debug log removed for production
    }, []);
    // Memoize string check to avoid recalculation
    const isStringDescription = useMemo(() => typeof description === 'string', [description]);
    // Measure if the text is overflowed - optimized with useCallback
    const checkOverflow = useCallback(() => {
        if (!isStringDescription || !descriptionRef.current) {
            setIsOverflowed(false);
            return;
        }
        const element = descriptionRef.current;
        const isOverflowing = element.scrollHeight > element.clientHeight;
        setIsOverflowed(isOverflowing);
    }, [isStringDescription]);
    useEffect(() => {
        checkOverflow();
        // Add resize observer for responsive overflow detection
        const resizeObserver = new ResizeObserver(checkOverflow);
        if (descriptionRef.current) {
            resizeObserver.observe(descriptionRef.current);
        }
        return () => {
            resizeObserver.disconnect();
        };
    }, [checkOverflow]);
    /**
     * Splits and highlights #hashtags / @mentions if description is a string.
     * If it's a node, just return it directly.
     * Memoized to prevent unnecessary re-processing
     */
    const renderDescription = useMemo(() => {
        if (!description) return null;
        // If already a node, just render it
        if (!isStringDescription) {
            return <>{description}</>;
        }
        // Otherwise, do the highlight logic
        const regex = /(#\w[\w-]*)|(@\w[\w-]*)|([^#@\s]+)/g;
        const tokens = description.match(regex) || [];
        return tokens.map((token, index) => {
            const key = `${token}-${index}`; // More stable key
            if (token.startsWith('#')) {
                return (
                    <InteractiveText
                        key={key}
                        onClick={() => handleHashtagClick(token)}
                        onKeyPress={(e) => {
                            if (e.key === 'Enter' || e.key === ' ') {
                                handleHashtagClick(token);
                            }
                        }}
                        role="button"
                        aria-label={`Hashtag ${token.substring(1)}`}
                        tabIndex={0}
                    >
                        {token}{' '}
                    </InteractiveText>
                );
            } else if (token.startsWith('@')) {
                return (
                    <InteractiveText
                        key={key}
                        onClick={() => handleMentionClick(token)}
                        onKeyPress={(e) => {
                            if (e.key === 'Enter' || e.key === ' ') {
                                handleMentionClick(token);
                            }
                        }}
                        role="button"
                        aria-label={`Mention ${token.substring(1)}`}
                        tabIndex={0}
                    >
                        {token}{' '}
                    </InteractiveText>
                );
            } else {
                return <span key={key}>{token} </span>;
            }
        });
    }, [description, isStringDescription, handleHashtagClick, handleMentionClick]);
    // Memoize modal title
    const modalTitle = useMemo(() => productTitle || 'Product Info', [productTitle]);
    return (
        <DescriptionWrapper>
            <DescriptionText ref={descriptionRef}>{renderDescription}</DescriptionText>
            {/** If it's a string AND we see overflow => "More info" button. */}
            {isStringDescription && isOverflowed && (
                <ToggleExpandButton onClick={openModal} aria-controls="description-modal">
                    More info
                </ToggleExpandButton>
            )}
            {/* The modal displays full text if it's a string, or fallback if it's a node */}
            <DescriptionModal
                isOpen={isModalOpen}
                onClose={closeModal}
                title={modalTitle}
            >
                {isStringDescription ? (
                    <FullDescription>{description}</FullDescription>
                ) : (
                    <FullDescription>(No expanded view for editing mode)</FullDescription>
                )}
            </DescriptionModal>
        </DescriptionWrapper>
    );
});
DescriptionMobileEditable.displayName = 'DescriptionMobileEditable';
DescriptionMobileEditable.propTypes = {
    // can be a string or a React node
    description: PropTypes.oneOfType([
        PropTypes.string,
        PropTypes.node,
    ]),
    productTitle: PropTypes.string,
};
DescriptionMobileEditable.defaultProps = {
    description: '',
    productTitle: 'Description',
};
export default DescriptionMobileEditable;
/* ------------------------------------------------------------------
   MODAL (with some sample tabbed content you can remove or customize)
-------------------------------------------------------------------- */
const DescriptionModal = memo(({isOpen, onClose, title, children}) => {
    const modalRef = useRef(null);
    const [activeTab, setActiveTab] = useState('overview');
    // Memoize keyboard handlers
    const handleKeyDown = useCallback(
        (e) => {
            if (e.key === 'Escape') onClose();
        },
        [onClose]
    );
    // Memoize tab change handler
    const handleTabChange = useCallback((tabName) => {
        setActiveTab(tabName);
    }, []);
    // Add/remove keydown listener; lock/unlock body scroll
    useEffect(() => {
        if (isOpen) {
            document.addEventListener('keydown', handleKeyDown);
            document.body.style.overflow = 'hidden';
        } else {
            document.removeEventListener('keydown', handleKeyDown);
            document.body.style.overflow = 'unset';
        }
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
            document.body.style.overflow = 'unset';
        };
    }, [isOpen, handleKeyDown]);
    // Auto-focus the modal
    useEffect(() => {
        if (isOpen && modalRef.current) {
            modalRef.current.focus();
        }
    }, [isOpen]);
    // Focus trapping - memoized for performance
    const handleTabKey = useCallback((e) => {
        if (e.key !== 'Tab') return;
        const focusableElements = modalRef.current?.querySelectorAll(
            'a[href], area[href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), button:not([disabled]), [tabindex="0"], [contenteditable]'
        );
        if (!focusableElements?.length) return;
        const firstElement = focusableElements[0];
        const lastElement = focusableElements[focusableElements.length - 1];
        if (!e.shiftKey && document.activeElement === lastElement) {
            e.preventDefault();
            firstElement.focus();
        } else if (e.shiftKey && document.activeElement === firstElement) {
            e.preventDefault();
            lastElement.focus();
        }
    }, []);
    useEffect(() => {
        if (isOpen) {
            document.addEventListener('keydown', handleTabKey);
        } else {
            document.removeEventListener('keydown', handleTabKey);
        }
        return () => {
            document.removeEventListener('keydown', handleTabKey);
        };
    }, [isOpen, handleTabKey]);
    // Memoize tab content to prevent unnecessary re-renders
    const renderTabContent = useCallback(() => {
        switch (activeTab) {
            case 'overview':
                return (
                    <TabPanel>
                        <h3>Overview</h3>
                        <p>
                            <strong>Price:</strong> $3999 <br/>
                            <strong>Condition:</strong> Used <br/>
                            <strong>Category:</strong> Computers / Laptops <br/>
                            <strong>Location:</strong> Berlin, 10365
                        </p>
                        <p>
                            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin vel
                            ligula eu libero euismod placerat in in nibh.
                        </p>
                    </TabPanel>
                );
            case 'shipping':
                return (
                    <TabPanel>
                        <h3>Shipping Options</h3>
                        <ul>
                            <li>Standard Shipping (3-5 days)</li>
                            <li>Express Shipping (1-2 days)</li>
                            <li>Local Pickup Available</li>
                        </ul>
                        <p>
                            Additional shipping details: tracking info, packaging, and potential
                            insurance can be described here.
                        </p>
                    </TabPanel>
                );
            case 'specs':
                return (
                    <TabPanel>
                        <h3>Technical Specs</h3>
                        <ul>
                            <li>Processor: Intel i7 10th Gen</li>
                            <li>RAM: 16GB DDR4</li>
                            <li>Storage: 512GB SSD</li>
                            <li>Graphics: NVIDIA RTX 2060</li>
                            <li>Screen: 15.6" FHD (1920x1080)</li>
                        </ul>
                        <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</p>
                    </TabPanel>
                );
            case 'seller':
                return (
                    <TabPanel>
                        <h3>Seller Info</h3>
                        <p>
                            Seller: <strong>John Doe</strong> <br/>
                            Rating: 4.5/5 <br/>
                            Location: Berlin, Germany <br/>
                            Return Policy: 14-day returns accepted
                        </p>
                        <p>Seller's additional information or disclaimers can be listed here.</p>
                    </TabPanel>
                );
            case 'reviews':
                return (
                    <TabPanel>
                        <h3>Reviews</h3>
                        <div>
                            <p>
                                <strong>JaneSmith:</strong> "Great product, fast shipping!"
                            </p>
                            <p>
                                <strong>Mike89:</strong> "Item as described, recommended."
                            </p>
                            <p>
                                <strong>SarahD:</strong> "Packaging was a bit damaged, but item works well."
                            </p>
                        </div>
                    </TabPanel>
                );
            case 'policies':
                return (
                    <TabPanel>
                        <h3>Policies &amp; Returns</h3>
                        <p>
                            Here you can find all important legal information, including your right
                            of return (AGB) and any special policies:
                        </p>
                        <ul>
                            <li>
                                <strong>Return Policy:</strong> 14-day returns (buyer pays return shipping)
                            </li>
                            <li>
                                <strong>Warranty:</strong> Covered under manufacturer warranty for 1 year
                            </li>
                            <li>
                                <strong>Refunds:</strong> Issued within 7 days after item is returned
                            </li>
                            <li>
                                <strong>AGB (Terms &amp; Conditions):</strong> Please review our full AGB at{' '}
                                <a href="#">www.example.com/agb</a>
                            </li>
                        </ul>
                        <p>
                            By purchasing, you agree to comply with relevant laws and the store's terms
                            of service. If you have any questions, please contact our support.
                        </p>
                    </TabPanel>
                );
            default:
                return null;
        }
    }, [activeTab]);
    // Memoize tab buttons to prevent unnecessary re-renders
    const tabButtons = useMemo(() => [
        { id: 'overview', label: 'Overview' },
        { id: 'shipping', label: 'Shipping' },
        { id: 'specs', label: 'Specs' },
        { id: 'seller', label: 'Seller' },
        { id: 'reviews', label: 'Reviews' },
        { id: 'policies', label: 'Policies' },
    ], []);
    if (!isOpen) return null;
    return ReactDOM.createPortal(
        <ModalOverlay onClick={onClose} isOpen={isOpen} aria-hidden={!isOpen}>
            <ModalContent
                onClick={(e) => e.stopPropagation()}
                role="dialog"
                aria-modal="true"
                aria-labelledby="modal-title"
                tabIndex="-1"
                ref={modalRef}
                isOpen={isOpen}
            >
                <ModalHeader>
                    <ModalTitle id="modal-title">{title}</ModalTitle>
                    <CloseButton onClick={onClose} aria-label="Close modal">
                        <X/>
                    </CloseButton>
                </ModalHeader>
                <TabsHeader role="tablist" aria-label="Product Info Tabs">
                    {tabButtons.map(tab => (
                        <TabButton
                            key={tab.id}
                            isActive={activeTab === tab.id}
                            onClick={() => handleTabChange(tab.id)}
                            role="tab"
                            aria-selected={activeTab === tab.id}
                        >
                            {tab.label}
                        </TabButton>
                    ))}
                </TabsHeader>
                <ModalBody>{renderTabContent()}</ModalBody>
            </ModalContent>
        </ModalOverlay>,
        document.body
    );
});
DescriptionModal.displayName = 'DescriptionModal';
DescriptionModal.propTypes = {
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    title: PropTypes.string,
    children: PropTypes.node.isRequired,
};
DescriptionModal.defaultProps = {
    title: 'Description',
};
/* ------------------------------------------------------------------
   STYLED COMPONENTS (modal + truncated text)
-------------------------------------------------------------------- */
// Design tokens for consistency and performance
const designTokens = {
    fonts: {
        primary: `'Inter', system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Fira Sans', 'Droid Sans', 'Helvetica Neue', Arial, sans-serif`,
    },
    colors: {
        text: '#171616',
        textSecondary: '#555555',
        accent: '#00A699',
        overlay: 'rgba(0, 0, 0, 0.4)',
        modalBg: '#fff',
        border: '#eee',
        hover: '#f0f0f0',
    },
    spacing: {
        xs: '4px',
        sm: '6px',
        md: '8px',
        lg: '10px',
        xl: '12px',
        xxl: '14px',
    },
    fontSize: {
        xs: '0.75rem',
        sm: '0.8rem',
        md: '0.85rem',
        lg: '0.9rem',
        xl: '0.95rem',
    },
    transitions: {
        fast: '0.2s ease-out',
        medium: '0.3s ease-out',
    },
    breakpoints: {
        mobile: '480px',
        tablet: '768px',
    },
};
// Base components for reusability
const BaseButton = styled.button`
    background: none;
    border: none;
    cursor: pointer;
    font-family: inherit;
    transition: color ${designTokens.transitions.fast};
    &:focus {
        outline: 2px solid ${designTokens.colors.accent};
        outline-offset: 2px;
    }
`;
const BaseText = styled.p`
    margin: 0;
    font-family: ${designTokens.fonts.primary};
    color: ${designTokens.colors.text};
`;
// Optimized main components
const DescriptionWrapper = styled.div`
    width: 100%;
    padding: ${designTokens.spacing.xs} ${designTokens.spacing.sm};
    font-family: ${designTokens.fonts.primary};
    color: ${designTokens.colors.text};
    position: relative;
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        padding: 0 ${designTokens.spacing.xs};
    }
`;
const DescriptionText = styled(BaseText)`
    font-size: ${designTokens.fontSize.md};
    line-height: 1.4;
    position: relative;
    overflow: hidden;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
    max-height: 2.8em;
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        font-size: ${designTokens.fontSize.sm};
        line-height: 1.3;
        max-height: 2.6em;
    }
`;
const ToggleExpandButton = styled(BaseButton)`
    color: ${designTokens.colors.accent};
    padding: 0;
    margin-top: ${designTokens.spacing.xs};
    font-size: ${designTokens.fontSize.sm};
    font-weight: 500;
    &:hover {
        text-decoration: underline;
    }
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        font-size: ${designTokens.fontSize.xs};
        margin-top: ${designTokens.spacing.xs};
    }
`;
const InteractiveText = styled.span`
    color: ${designTokens.colors.accent};
    cursor: pointer;
    font-weight: 600;
    &:hover,
    &:focus {
        text-decoration: underline;
    }
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        font-size: ${designTokens.fontSize.md};
    }
`;
export const FullDescription = styled(BaseText)`
    margin: ${designTokens.spacing.sm} 0;
    word-wrap: break-word;
    white-space: pre-wrap;
    font-size: ${designTokens.fontSize.md};
    line-height: 1.4;
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        font-size: ${designTokens.fontSize.sm};
        margin: ${designTokens.spacing.xs} 0;
    }
`;
/* MODAL COMPONENTS - Optimized for performance */
const ModalOverlay = styled.div`
    position: fixed;
    inset: 0;
    background-color: ${designTokens.colors.overlay};
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    opacity: ${(props) => (props.isOpen ? '1' : '0')};
    transition: opacity ${designTokens.transitions.fast};
    padding: 0 ${designTokens.spacing.lg};
    box-sizing: border-box;
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        padding: 0 ${designTokens.spacing.sm};
    }
`;
const ModalContent = styled.div`
    background-color: ${designTokens.colors.modalBg};
    width: 100%;
    max-width: 600px;
    border-radius: ${designTokens.spacing.sm};
    overflow: hidden;
    display: flex;
    flex-direction: column;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    transform: ${(props) => (props.isOpen ? 'translateY(0)' : 'translateY(-10px)')};
    transition: transform ${designTokens.transitions.medium}, box-shadow ${designTokens.transitions.medium};
    outline: none;
    font-family: ${designTokens.fonts.primary};
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        max-width: 500px;
    }
    @media (max-width: ${designTokens.breakpoints.mobile}) {
        max-width: 400px;
    }
`;
const ModalHeader = styled.div`
    padding: ${designTokens.spacing.lg} ${designTokens.spacing.xxl};
    background-color: #f5f5f5;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid ${designTokens.colors.border};
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        padding: ${designTokens.spacing.md} ${designTokens.spacing.xl};
    }
`;
const ModalTitle = styled.h2`
    margin: 0;
    font-size: ${designTokens.fontSize.xl};
    font-weight: 600;
    color: ${designTokens.colors.text};
    line-height: 1.3;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    max-height: calc(2 * 1.3em);
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        font-size: ${designTokens.fontSize.lg};
        max-height: calc(2 * 1.2em);
    }
`;
const CloseButton = styled(BaseButton)`
    color: #666;
    font-size: 1rem;
    &:hover {
        color: ${designTokens.colors.text};
    }
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        font-size: ${designTokens.fontSize.lg};
    }
`;
const TabsHeader = styled.div`
    display: flex;
    align-items: center;
    background-color: #fafafa;
    border-bottom: 1px solid ${designTokens.colors.border};
`;
const TabButton = styled(BaseButton)`
    flex: 1;
    padding: ${designTokens.spacing.lg} 0;
    background-color: ${(props) => (props.isActive ? '#fff' : 'transparent')};
    font-size: ${designTokens.fontSize.md};
    font-weight: 600;
    color: ${(props) => (props.isActive ? designTokens.colors.accent : '#444')};
    border-bottom: ${(props) => (props.isActive ? `3px solid ${designTokens.colors.accent}` : 'none')};
    transition: background-color ${designTokens.transitions.medium}, color ${designTokens.transitions.medium};
    &:hover {
        background-color: ${designTokens.colors.hover};
    }
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        font-size: ${designTokens.fontSize.sm};
        padding: ${designTokens.spacing.md} 0;
    }
`;
const ModalBody = styled.div`
    padding: ${designTokens.spacing.xl};
    overflow-y: auto;
    max-height: 70vh;
    background-color: #fff;
    @media (max-width: ${designTokens.breakpoints.tablet}) {
        padding: ${designTokens.spacing.lg};
    }
`;
const TabPanel = styled.div`
    h3 {
        margin-top: 0;
        font-size: ${designTokens.fontSize.xl};
        font-weight: 600;
        color: ${designTokens.colors.accent};
        line-height: 1.3;
        margin-bottom: ${designTokens.spacing.sm};
        @media (max-width: ${designTokens.breakpoints.tablet}) {
            font-size: ${designTokens.fontSize.lg};
        }
    }
    ul {
        margin-left: 18px;
        list-style: disc;
        padding-left: 0;
    }
    p,
    li {
        font-size: ${designTokens.fontSize.md};
        line-height: 1.4;
        color: #333;
        margin-bottom: 0.6em;
        @media (max-width: ${designTokens.breakpoints.tablet}) {
            font-size: ${designTokens.fontSize.sm};
            line-height: 1.3;
            margin-bottom: 0.5em;
        }
    }
    a {
        color: ${designTokens.colors.accent};
        text-decoration: none;
        &:hover {
            text-decoration: underline;
        }
    }
`;
