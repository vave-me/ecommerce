"use client"
import React, {useState} from 'react';
import {
    Search,
    BellRing,
    ShoppingBag,
    Heart,
    MessageCircle,
    Tag,
    TrendingUp,
    Menu,
    X,
    ChevronDown,
    User,
    Filter
} from '@/icons';
import styles from './page.module.css';
const EcommerceSocialPlatform = () => {
    const [isMenuOpen, setIsMenuOpen] = useState(false);
    const toggleMenu = () => setIsMenuOpen(!isMenuOpen);
    return (
        <div className={styles.appContainer}>
            {/* Header */}
            <header className={styles.header}>
                <div className={styles.container}>
                    <div className={styles.headerContent}>
                        {/* Logo & Mobile Menu Toggle */}
                        <div className={styles.logoContainer}>
                            <button onClick={toggleMenu} className={styles.menuToggle}>
                                {isMenuOpen ? <X size={24}/> : <Menu size={24}/>}
                            </button>
                            <div className={styles.logo}>DealSocial</div>
                        </div>
                        {/* Search */}
                        <div className={styles.searchContainer}>
                            <div className={styles.searchInputWrapper}>
                                <input
                                    type="text"
                                    placeholder="Search for deals, products, or users..."
                                    className={styles.searchInput}
                                />
                                <Search className={styles.searchIcon} size={18}/>
                            </div>
                        </div>
                        {/* Nav Icons */}
                        <div className={styles.navIcons}>
                            <button className={styles.iconButton}>
                                <BellRing size={20}/>
                            </button>
                            <button className={styles.iconButton}>
                                <ShoppingBag size={20}/>
                            </button>
                            <button className={styles.postDealButton}>
                                Post Deal
                            </button>
                            <button className={styles.userButton}>
                                <User size={20}/>
                                <ChevronDown size={16}/>
                            </button>
                        </div>
                    </div>
                </div>
                {/* Mobile Search - Shows below header on mobile */}
                <div className={styles.mobileSearch}>
                    <div className={styles.mobileSearchInputWrapper}>
                        <input
                            type="text"
                            placeholder="Search deals & products..."
                            className={styles.mobileSearchInput}
                        />
                        <Search className={styles.mobileSearchIcon} size={16}/>
                    </div>
                </div>
            </header>
            {/* Mobile Menu */}
            {isMenuOpen && (
                <div className={styles.mobileMenu}>
                    <div className={styles.container}>
                        <div className={styles.mobileMenuContent}>
                            <button className={styles.mobileMenuItem}>
                                <User size={20}/>
                                <span>Profile</span>
                            </button>
                            <button className={styles.mobileMenuItem}>
                                <Heart size={20}/>
                                <span>Saved Deals</span>
                            </button>
                            <button className={styles.mobileMenuItem}>
                                <MessageCircle size={20}/>
                                <span>Messages</span>
                            </button>
                            <button className={styles.mobileMenuItem}>
                                <ShoppingBag size={20}/>
                                <span>Purchases</span>
                            </button>
                            <button className={styles.mobilePostButton}>
                                Post a Deal
                            </button>
                        </div>
                    </div>
                </div>
            )}
            {/* Main Content */}
            <main className={styles.main}>
                <div className={styles.container}>
                    {/* Categories Bar */}
                    <div className={styles.categoriesWrapper}>
                        <div className={styles.categoriesContainer}>
                            <button className={styles.categoryButtonActive}>
                                All Deals
                            </button>
                            <button className={styles.categoryButton}>
                                Electronics
                            </button>
                            <button className={styles.categoryButton}>
                                Fashion
                            </button>
                            <button className={styles.categoryButton}>
                                Home & Garden
                            </button>
                            <button className={styles.categoryButton}>
                                Travel
                            </button>
                            <button className={styles.categoryButton}>
                                Services
                            </button>
                            <button className={styles.categoryButton}>
                                Food
                            </button>
                        </div>
                    </div>
                    {/* Featured & Filters row */}
                    <div className={styles.filtersRow}>
                        <div className={styles.featuredHeader}>
                            <div className={styles.featuredTitle}>🔥 Hot Deals</div>
                            <div className={styles.dealCount}>(237 active deals)</div>
                        </div>
                        <div className={styles.filterButtons}>
                            <button className={styles.filterButton}>
                                <Filter size={14}/>
                                <span>Filters</span>
                            </button>
                            <button className={styles.filterButton}>
                                <TrendingUp size={14}/>
                                <span>Trending</span>
                            </button>
                            <button className={styles.filterButton}>
                                <span>Newest</span>
                                <ChevronDown size={14}/>
                            </button>
                        </div>
                    </div>
                    {/* Deals Grid */}
                    <div className={styles.dealsGrid}>
                        {/* Deal Card 1 */}
                        <div className={styles.dealCard}>
                            <div className={styles.dealImageContainer}>
                                <img src="/api/placeholder/500/300" alt="Product" className={styles.dealImage}/>
                                <div className={styles.discountBadge}>
                                    -70%
                                </div>
                                <button className={styles.favoriteButton}>
                                    <Heart size={18} className={styles.favoriteIcon}/>
                                </button>
                            </div>
                            <div className={styles.dealContent}>
                                <div className={styles.dealHeader}>
                                    <div>
                                        <div className={styles.dealPoster}>Posted by @techdeals · 2h ago</div>
                                        <h3 className={styles.dealTitle}>Wireless Noise Cancelling Headphones - Premium
                                            Sound Quality</h3>
                                    </div>
                                </div>
                                <div className={styles.dealPricing}>
                                    <div className={styles.currentPrice}>$89.99</div>
                                    <div className={styles.originalPrice}>$299.99</div>
                                    <div className={styles.vendorName}>
                                        <span>Amazon</span>
                                    </div>
                                </div>
                                <div className={styles.dealFooter}>
                                    <div className={styles.dealStats}>
                                        <div className={styles.dealStat}>
                                            <MessageCircle size={16} className={styles.commentIcon}/>
                                            <span>23</span>
                                        </div>
                                        <div className={styles.dealStat}>
                                            <TrendingUp size={16} className={styles.trendingIcon}/>
                                            <span>142</span>
                                        </div>
                                    </div>
                                    <div className={styles.shippingInfo}>Free Shipping</div>
                                </div>
                            </div>
                        </div>
                        {/* Deal Card 2 */}
                        <div className={styles.dealCard}>
                            <div className={styles.dealImageContainer}>
                                <img src="/api/placeholder/500/300" alt="Product" className={styles.dealImage}/>
                                <div className={styles.localBadge}>
                                    LOCAL
                                </div>
                                <button className={styles.favoriteButton}>
                                    <Heart size={18} className={styles.favoriteIcon}/>
                                </button>
                            </div>
                            <div className={styles.dealContent}>
                                <div className={styles.dealHeader}>
                                    <div>
                                        <div className={styles.dealPoster}>Posted by @localpicks · 5h ago</div>
                                        <h3 className={styles.dealTitle}>Cozy Corner Café - Buy One Get One Free on All
                                            Coffee Drinks</h3>
                                    </div>
                                </div>
                                <div className={styles.dealPricing}>
                                    <div className={styles.bogoBadge}>
                                        BOGO Deal
                                    </div>
                                    <div className={styles.expiryInfo}>
                                        <span>Expires in 3 days</span>
                                    </div>
                                </div>
                                <div className={styles.dealFooter}>
                                    <div className={styles.dealStats}>
                                        <div className={styles.dealStat}>
                                            <MessageCircle size={16} className={styles.commentIcon}/>
                                            <span>7</span>
                                        </div>
                                        <div className={styles.dealStat}>
                                            <TrendingUp size={16} className={styles.trendingIcon}/>
                                            <span>32</span>
                                        </div>
                                    </div>
                                    <div className={styles.locationInfo}>San Francisco, CA</div>
                                </div>
                            </div>
                        </div>
                        {/* Deal Card 3 */}
                        <div className={styles.dealCard}>
                            <div className={styles.dealImageContainer}>
                                <img src="/api/placeholder/500/300" alt="Product" className={styles.dealImage}/>
                                <div className={styles.bestPriceBadge}>
                                    BEST PRICE
                                </div>
                                <button className={styles.favoriteButton}>
                                    <Heart size={18} className={styles.favoriteActive}/>
                                </button>
                            </div>
                            <div className={styles.dealContent}>
                                <div className={styles.dealHeader}>
                                    <div>
                                        <div className={styles.dealPoster}>Posted by @gadgetpro · 1d ago</div>
                                        <h3 className={styles.dealTitle}>4K Smart TV 55" - Ultra HD with Voice
                                            Control</h3>
                                    </div>
                                </div>
                                <div className={styles.dealPricing}>
                                    <div className={styles.currentPrice}>$399.99</div>
                                    <div className={styles.originalPrice}>$699.99</div>
                                    <div className={styles.vendorNameBlue}>
                                        <span>BestBuy</span>
                                    </div>
                                </div>
                                <div className={styles.dealFooter}>
                                    <div className={styles.dealStats}>
                                        <div className={styles.dealStat}>
                                            <MessageCircle size={16} className={styles.commentIcon}/>
                                            <span>56</span>
                                        </div>
                                        <div className={styles.dealStat}>
                                            <TrendingUp size={16} className={styles.trendingIcon}/>
                                            <span>283</span>
                                        </div>
                                    </div>
                                    <div className={styles.cashbackInfo}>+$20 Cashback</div>
                                </div>
                            </div>
                        </div>
                    </div>
                    {/* Trending Communities Section */}
                    <div className={styles.communitiesSection}>
                        <h2 className={styles.sectionTitle}>👥 Trending Communities</h2>
                        <div className={styles.communitiesGrid}>
                            <div className={styles.communityCard}>
                                <div className={styles.communityContent}>
                                    <div className={styles.communityIconOrange}>
                                        T
                                    </div>
                                    <div className={styles.communityInfo}>
                                        <h3 className={styles.communityName}>TechDeals</h3>
                                        <p className={styles.communityMembers}>124k members</p>
                                    </div>
                                    <button className={styles.joinButton}>
                                        Join
                                    </button>
                                </div>
                            </div>
                            <div className={styles.communityCard}>
                                <div className={styles.communityContent}>
                                    <div className={styles.communityIconBlue}>
                                        T
                                    </div>
                                    <div className={styles.communityInfo}>
                                        <h3 className={styles.communityName}>TravelFinds</h3>
                                        <p className={styles.communityMembers}>87k members</p>
                                    </div>
                                    <button className={styles.joinButton}>
                                        Join
                                    </button>
                                </div>
                            </div>
                            <div className={styles.communityCard}>
                                <div className={styles.communityContent}>
                                    <div className={styles.communityIconGreen}>
                                        H
                                    </div>
                                    <div className={styles.communityInfo}>
                                        <h3 className={styles.communityName}>HomeDecor</h3>
                                        <p className={styles.communityMembers}>62k members</p>
                                    </div>
                                    <button className={styles.joinedButton}>
                                        Joined
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                    {/* Classified Ads Section */}
                    <div className={styles.classifiedsSection}>
                        <div className={styles.classifiedsHeader}>
                            <h2 className={styles.sectionTitle}>📦 Local Classifieds</h2>
                            <button className={styles.viewAllButton}>View all</button>
                        </div>
                        <div className={styles.classifiedsGrid}>
                            <div className={styles.classifiedCard}>
                                <div className={styles.classifiedImageContainer}>
                                    <img src="/api/placeholder/400/300" alt="Classified"
                                         className={styles.classifiedImage}/>
                                    <button className={styles.favoriteButton}>
                                        <Heart size={18} className={styles.favoriteIcon}/>
                                    </button>
                                </div>
                                <div className={styles.classifiedContent}>
                                    <h3 className={styles.classifiedTitle}>Mountain Bike - Great Condition</h3>
                                    <div className={styles.classifiedPrice}>$350</div>
                                    <div className={styles.classifiedMeta}>
                                        <span>Seattle, WA</span>
                                        <span>Posted 2 days ago</span>
                                    </div>
                                </div>
                            </div>
                            <div className={styles.classifiedCard}>
                                <div className={styles.classifiedImageContainer}>
                                    <img src="/api/placeholder/400/300" alt="Classified"
                                         className={styles.classifiedImage}/>
                                    <button className={styles.favoriteButton}>
                                        <Heart size={18} className={styles.favoriteIcon}/>
                                    </button>
                                </div>
                                <div className={styles.classifiedContent}>
                                    <h3 className={styles.classifiedTitle}>Vintage Record Player</h3>
                                    <div className={styles.classifiedPrice}>$120</div>
                                    <div className={styles.classifiedMeta}>
                                        <span>Portland, OR</span>
                                        <span>Posted 1 day ago</span>
                                    </div>
                                </div>
                            </div>
                            <div className={styles.classifiedCard}>
                                <div className={styles.classifiedImageContainer}>
                                    <img src="/api/placeholder/400/300" alt="Classified"
                                         className={styles.classifiedImage}/>
                                    <button className={styles.favoriteButton}>
                                        <Heart size={18} className={styles.favoriteIcon}/>
                                    </button>
                                </div>
                                <div className={styles.classifiedContent}>
                                    <h3 className={styles.classifiedTitle}>IKEA Desk - Like New</h3>
                                    <div className={styles.classifiedPrice}>$75</div>
                                    <div className={styles.classifiedMeta}>
                                        <span>San Francisco, CA</span>
                                        <span>Posted 5 hours ago</span>
                                    </div>
                                </div>
                            </div>
                            <div className={styles.classifiedCard}>
                                <div className={styles.classifiedImageContainer}>
                                    <img src="/api/placeholder/400/300" alt="Classified"
                                         className={styles.classifiedImage}/>
                                    <button className={styles.favoriteButton}>
                                        <Heart size={18} className={styles.favoriteIcon}/>
                                    </button>
                                </div>
                                <div className={styles.classifiedContent}>
                                    <h3 className={styles.classifiedTitle}>iPhone 12 Pro - 256GB</h3>
                                    <div className={styles.classifiedPrice}>$550</div>
                                    <div className={styles.classifiedMeta}>
                                        <span>Chicago, IL</span>
                                        <span>Posted 3 days ago</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </main>
            {/* Footer */}
            <footer className={styles.footer}>
                <div className={styles.container}>
                    <div className={styles.footerContent}>
                        © 2025 DealSocial - The social platform for deals and classifieds
                    </div>
                </div>
            </footer>
            {/* Mobile Action Button */}
            <div className={styles.mobileActionButton}>
                <button className={styles.floatingActionButton}>
                    <Tag size={24}/>
                </button>
            </div>
        </div>
    );
};
export default EcommerceSocialPlatform;