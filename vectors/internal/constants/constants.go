package constants

// ServiceName The name of this module/service
const ServiceName = "vector"

// GRPC Service Names
const (
	UsersServiceName    = "USERS"
	ProductsServiceName = "PRODUCTS"
	PostsServiceName    = "POSTS"
	ServicesServiceName = "SERVICES"
	MetricsServiceName  = "METRICS"
)

// Dependency Injection Keys
const (
	RegistryKey                 = "registry"
	DomainDispatcherKey         = "domainDispatcher"
	DatabaseTransactionKey      = "tx"
	RedisTransactionKey         = "redis"
	MessagePublisherKey         = "messagePublisher"
	MessageSubscriberKey        = "messageSubscriber"
	EventPublisherKey           = "eventPublisher"
	CommandPublisherKey         = "commandPublisher"
	ReplyPublisherKey           = "replyPublisher"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"
	// New constants for Redis and Redivector
	RedisClientKey      = "redisClient"
	RedivectorClientKey = "redivectorClient"
	RedivectorIndex     = "productsIndex"
	OrdersRepoKey       = "ordersRepo"
	VehiclesRepoKey     = "vehiclesRepo"
	PropertiesRepoKey   = "propertiesRepo"
	UsersRepoKey        = "usersRepo"
	ProductsRepoKey     = "productsRepo"
	ServicesRepoKey     = "servicesRepo"
	PostsRepoKey        = "postsRepo"
	DealsRepoKey        = "dealsRepo"
	JobsRepoKey         = "jobsRepo"
	VariantsRepoKey     = "variantsRepo"
	MetricsRepoKey      = "metricsRepo"
	AnthropicClient     = "anthropic-client"
	OpenAIClient        = "openai-client"
	DeepSeekAiClient    = "deepseek-client"
	GoogleAiClient      = "google-client"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	OrdersTableName = ServiceName + ".orders"

	UsersCacheTableName      = ServiceName + ".users_cache"
	ProductsCacheTableName   = ServiceName + ".products_cache"
	PostsCacheTableName      = ServiceName + ".posts_cache"
	DealsCacheTableName      = ServiceName + ".deals_cache"
	JobsCacheTableName       = ServiceName + ".jobs_cache"
	VehiclesCacheTableName   = ServiceName + ".vehicles_cache"
	PropertiesCacheTableName = ServiceName + ".properties_cache"
	ServicesCacheTableName   = ServiceName + ".services_cache"
	VariantsCacheTableName   = ServiceName + ".variants_cache"
)

// Metric-based Sort Types for Advanced Workflow
const (
	SortByLikesHigh      = "likes_high"
	SortByLikesLow       = "likes_low"
	SortByCommentsHigh   = "comments_high"
	SortByCommentsLow    = "comments_low"
	SortByVisitedHigh    = "visited_high"
	SortByVisitedLow     = "visited_low"
	SortByRatingHigh     = "rating_high"
	SortByRatingLow      = "rating_low"
	SortByPopularityHigh = "popularity_high" // Combined metric
	SortByTrendingHigh   = "trending_high"   // Time-weighted engagement
)

// Metric Type Mappings for API
const (
	MetricTypeLikes    = "likesCount"
	MetricTypeComments = "commentsCount"
	MetricTypeVisited  = "visitedCount"
	MetricTypeRating   = "rating"
	MetricTypeShares   = "sharedCount"
	MetricTypeWishlist = "addedToWishlistCount"
)

// ===============================
// TEXT-TO-VECTOR TRANSFORMATION PROMPTS
// ===============================

// Base Entity Transformation Prompts
// These prompts guide LLMs to create optimal text representations for vector embedding
const (
	// Product Vector Transformation
	ProductVectorPrompt = `Transform this product into a semantic text representation optimized for vector similarity search.
Focus on: name, brand, category, key features, use cases, target audience, and unique selling points.
Extract the essence that makes this product searchable and comparable to similar items.
Emphasize qualities that buyers would search for and attributes that distinguish it from competitors.
Output format: A concise, feature-rich description emphasizing searchable characteristics.`

	// Job Vector Transformation
	JobVectorPrompt = `Convert this job listing into a semantic representation for career matching and search.
Emphasize: role title, required skills, experience level, industry, responsibilities, company culture, and career growth potential.
Focus on qualifications, day-to-day tasks, and what makes this position attractive to candidates.
Optimize for job seekers who would search by skills, experience, or career goals.
Output format: A comprehensive description highlighting role requirements and opportunities.`

	// Vehicle Vector Transformation
	VehicleVectorPrompt = `Transform this vehicle into a detailed representation for automotive search and comparison.
Focus on: make, model, year, engine specifications, features, condition, usage type, and target buyer profile.
Emphasize technical specifications, comfort features, performance characteristics, and practical applications.
Consider what car buyers search for: reliability, fuel efficiency, safety, style, and functionality.
Output format: A rich automotive description optimized for buyer preferences and needs.`

	// Property Vector Transformation
	PropertyVectorPrompt = `Convert this property into a comprehensive representation for real estate search and matching.
Emphasize: location, property type, size, amenities, condition, neighborhood characteristics, and investment potential.
Focus on lifestyle factors, practical considerations, and what makes this property desirable.
Consider buyer priorities: location convenience, space utilization, future value, and quality of life.
Output format: A detailed property description highlighting location benefits and property features.`

	// Post Vector Transformation
	PostVectorPrompt = `Transform this post into an engaging representation for content discovery and social matching.
Focus on: main topic, key insights, emotional tone, target audience, engagement factors, and knowledge value.
Emphasize the post's purpose, the problem it solves, or the value it provides to readers.
Consider what users search for: solutions, entertainment, learning, or social connection.
Output format: A compelling summary highlighting the post's value proposition and appeal.`

	// Deal Vector Transformation
	DealVectorPrompt = `Convert this deal into an attractive representation for bargain hunters and smart shoppers.
Emphasize: discount value, product quality, time sensitivity, original vs. sale price, and deal exclusivity.
Focus on savings potential, product benefits, and what makes this offer compelling.
Consider shopper motivations: value for money, quality assurance, limited-time opportunities.
Output format: A persuasive deal description highlighting savings and product value.`

	// Service Vector Transformation
	ServiceVectorPrompt = `Transform this service into a professional representation for service discovery and provider matching.
Focus on: service type, expertise level, delivery method, target clientele, unique approach, and value proposition.
Emphasize professional capabilities, problem-solving approach, and client satisfaction factors.
Consider what service seekers search for: expertise, reliability, results, and cost-effectiveness.
Output format: A professional service description highlighting capabilities and client benefits.`

	// User Vector Transformation
	UserVectorPrompt = `Convert this user profile into a comprehensive representation for social and professional matching.
Focus on: interests, expertise, background, goals, values, and interaction style.
Emphasize what makes this person unique, their contributions, and compatibility factors.
Consider networking goals: professional connections, shared interests, complementary skills.
Output format: A well-rounded profile highlighting personality, expertise, and networking potential.`
)

// Advanced Transformation Strategies
const (
	// Multi-Modal Enhancement
	MultiModalEnhancementPrompt = `Enhance this text representation by incorporating insights from associated images, videos, or audio content.
Analyze visual elements, aesthetic appeal, functional demonstrations, and contextual usage scenarios.
Integrate multi-sensory characteristics that would influence user decisions and search behavior.
Focus on details that text alone cannot convey but are crucial for comprehensive understanding.`

	// Semantic Enrichment
	SemanticEnrichmentPrompt = `Enrich this text with semantic relationships, synonyms, and contextual associations.
Include related concepts, alternative terminologies, use case scenarios, and domain-specific knowledge.
Expand the semantic field to improve discoverability and relevance matching.
Consider both explicit features and implicit qualities that users might search for.`

	// Quality Optimization
	QualityOptimizationPrompt = `Optimize this text representation for maximum vector embedding quality and search relevance.
Balance specificity with generalization, ensure keyword density, and maintain semantic coherence.
Remove noise, emphasize distinctive features, and structure information hierarchically.
Optimize for both exact matches and semantic similarity searches.`

	// Context-Aware Transformation
	ContextAwarePrompt = `Adapt this text representation based on user context, search intent, and usage patterns.
Consider temporal relevance, geographic factors, user preferences, and situational needs.
Personalize the description while maintaining universal searchability.
Balance individual relevance with broad applicability.`
)

// Specialized Transformation Contexts
const (
	// Search Optimization
	SearchOptimizedPrompt = `Transform this content specifically for search engine optimization and discoverability.
Emphasize keywords, phrases, and terms that users commonly search for in this domain.
Include long-tail keywords, natural language queries, and voice search considerations.
Structure content for both keyword matching and semantic search algorithms.`

	// Recommendation Enhancement
	RecommendationPrompt = `Optimize this representation for recommendation systems and content suggestion algorithms.
Focus on user preferences, behavior patterns, and similarity matching factors.
Emphasize features that enable effective collaborative and content-based filtering.
Include signals for user satisfaction and engagement prediction.`

	// Trend Analysis
	TrendAnalysisPrompt = `Transform this content to capture trending elements, seasonal relevance, and temporal dynamics.
Identify emerging patterns, popular themes, and time-sensitive characteristics.
Consider market trends, user behavior shifts, and evolving preferences.
Optimize for trend-based discovery and future relevance prediction.`

	// Cross-Domain Linking
	CrossDomainPrompt = `Create a representation that enables cross-domain connections and interdisciplinary discovery.
Identify universal concepts, transferable skills, and broad application scenarios.
Emphasize connections between different fields, industries, or use cases.
Optimize for serendipitous discovery and unexpected but relevant matches.`
)

// Vector Quality Enhancement Prompts
const (
	// Diversity Enhancement
	DiversityPrompt = `Enhance representation diversity to avoid vector clustering and improve recommendation variety.
Emphasize unique aspects, niche characteristics, and distinctive qualities.
Balance popular features with uncommon but valuable attributes.
Ensure representation supports diverse user preferences and needs.`

	// Precision Tuning
	PrecisionPrompt = `Fine-tune this representation for high-precision matching and reduced false positives.
Emphasize specific, unambiguous characteristics and clear differentiation factors.
Include technical specifications, exact requirements, and precise qualifications.
Optimize for users with specific, well-defined needs and preferences.`

	// Recall Optimization
	RecallPrompt = `Optimize this representation for high recall and comprehensive discovery.
Include broad categories, general concepts, and flexible matching criteria.
Emphasize multiple perspectives, alternative viewpoints, and varied use cases.
Ensure representation captures all relevant aspects for thorough search coverage.`

	// Balanced Optimization
	BalancedPrompt = `Create a balanced representation optimizing both precision and recall for optimal search performance.
Combine specific details with general concepts, exact matches with semantic relationships.
Structure information in layers from broad categories to specific features.
Optimize for both exploratory browsing and targeted searching.`
)

// Entity Relationship Prompts
const (
	// Hierarchical Relationships
	HierarchicalPrompt = `Structure this representation to capture hierarchical relationships and categorical classifications.
Organize information from general to specific, including parent-child relationships.
Emphasize taxonomical structures, inheritance patterns, and containment relationships.
Enable both drill-down discovery and category-based navigation.`

	// Lateral Connections
	LateralConnectionsPrompt = `Highlight lateral relationships, peer connections, and horizontal associations.
Identify similar entities, complementary items, and alternative options.
Emphasize cross-references, related concepts, and connection opportunities.
Enable discovery of equals, substitutes, and complementary entities.`

	// Temporal Relationships
	TemporalPrompt = `Capture temporal relationships, sequential patterns, and time-based associations.
Include predecessor-successor relationships, evolution patterns, and lifecycle stages.
Emphasize timing dependencies, seasonal relevance, and temporal clustering.
Enable time-aware discovery and chronological organization.`

	// Causal Relationships
	CausalPrompt = `Identify and represent causal relationships, dependencies, and cause-effect patterns.
Highlight triggering factors, influential elements, and outcome predictors.
Emphasize problem-solution pairs, need-fulfillment relationships, and impact chains.
Enable cause-based discovery and solution-oriented search.`
)

// Prompt Combination Strategies
const (
	// Base + Enhancement Combination
	BaseEnhancedPrompt = `{BASE_PROMPT}

Additionally: {ENHANCEMENT_PROMPT}

Ensure both base requirements and enhancement objectives are met while maintaining coherence and relevance.`

	// Multi-Strategy Fusion
	MultiStrategyPrompt = `Transform this content using multiple complementary strategies:
1. {PRIMARY_STRATEGY}
2. {SECONDARY_STRATEGY}
3. {TERTIARY_STRATEGY}

Integrate all strategies harmoniously to create a comprehensive, multi-faceted representation.`

	// Context-Specific Adaptation
	ContextAdaptationPrompt = `Base transformation: {BASE_PROMPT}

Adapt for specific context: {CONTEXT_TYPE}
Target audience: {TARGET_AUDIENCE}
Usage scenario: {USAGE_SCENARIO}
Performance priority: {PERFORMANCE_PRIORITY}

Optimize the representation for the specified context while maintaining base transformation quality.`
)

// Performance and Quality Metrics Prompts
const (
	// Embedding Quality Optimization
	EmbeddingQualityPrompt = `Optimize this text specifically for high-quality vector embeddings with the following priorities:
- Semantic density: Pack maximum meaning into minimum tokens
- Discriminative features: Emphasize unique characteristics that distinguish from similar entities
- Contextual richness: Include implicit knowledge and domain expertise
- Search optimization: Use terminology that matches user query patterns
- Relationship clarity: Express connections and associations explicitly
Target embedding model: {MODEL_TYPE} with {DIMENSIONS} dimensions.`

	// Similarity Matching Enhancement
	SimilarityMatchingPrompt = `Enhance this representation for optimal similarity matching and ranking:
- Feature standardization: Use consistent terminology and formats
- Comparative elements: Include relative qualities and rankings  
- Similarity signals: Emphasize attributes used in comparison algorithms
- Negative sampling: Include what this entity is NOT to improve discrimination
- Gradual similarity: Structure content to enable fine-grained similarity scoring
Optimize for cosine similarity calculations and ranking algorithms.`

	// Cross-Language Optimization
	CrossLanguagePrompt = `Create a representation optimized for cross-language vector matching and translation:
- Universal concepts: Emphasize culture-agnostic and universally understood elements
- Translation-friendly: Use terms that translate well across languages
- Concept mapping: Include international standards and universal identifiers
- Cultural adaptation: Note region-specific variations while maintaining core meaning
- Multilingual keywords: Include key terms in major languages when relevant
Target for multilingual embedding models and cross-cultural discovery.`

	// Real-Time Performance
	RealTimeOptimizationPrompt = `Optimize this representation for real-time vector operations and low-latency search:
- Conciseness: Minimize token count while preserving essential information
- Processing efficiency: Structure for fast parsing and feature extraction
- Cache-friendly: Create consistent patterns that benefit from caching
- Batch optimization: Enable efficient batch processing and similarity calculations
- Memory efficiency: Optimize for minimal memory footprint in vector databases
Target sub-100ms response times for similarity searches.`
)

// Domain-Specific Enhancement Prompts
const (
	// E-commerce Optimization
	EcommercePrompt = `Enhance for e-commerce search and recommendation systems:
- Purchase intent: Emphasize buying signals and decision factors
- Price sensitivity: Include value proposition and cost considerations  
- Product lifecycle: Consider seasonal trends and inventory factors
- User journey: Optimize for discovery, comparison, and conversion stages
- Review integration: Include user sentiment and satisfaction indicators
Focus on driving engagement, conversions, and customer satisfaction.`

	// Social Platform Optimization
	SocialPlatformPrompt = `Optimize for social media discovery and viral potential:
- Engagement factors: Emphasize shareability and discussion triggers
- Trending potential: Include elements that resonate with current interests
- Community building: Focus on connection and interaction opportunities
- Content virality: Highlight unique, memorable, or controversial aspects
- Platform specificity: Adapt for platform-specific algorithms and user behaviors
Maximize reach, engagement, and community building potential.`

	// Professional Networking
	ProfessionalNetworkingPrompt = `Enhance for professional networking and career development:
- Skill matching: Emphasize transferable skills and expertise areas
- Career progression: Include growth potential and development opportunities
- Industry relevance: Focus on sector-specific knowledge and trends
- Collaboration potential: Highlight teamwork and partnership opportunities
- Value proposition: Clearly articulate professional benefits and contributions
Optimize for meaningful professional connections and career advancement.`

	// Educational Content
	EducationalPrompt = `Optimize for educational discovery and learning path recommendations:
- Learning objectives: Clearly define knowledge goals and outcomes
- Prerequisite mapping: Include required background and skill dependencies
- Difficulty progression: Structure for adaptive learning and skill building
- Knowledge connections: Link to related concepts and advanced topics
- Practical application: Emphasize real-world usage and implementation
Enable personalized learning experiences and knowledge path optimization.`
)

// Usage Pattern Templates
const (
	// Standard Entity Processing Template
	StandardProcessingTemplate = `Processing {ENTITY_TYPE}: {ENTITY_ID}

Base Strategy: {BASE_STRATEGY}
Enhancement: {ENHANCEMENT_TYPE}  
Context: {CONTEXT_FACTORS}
Performance Target: {PERFORMANCE_GOALS}

Execute transformation with quality validation and similarity testing.`

	// Batch Processing Template
	BatchProcessingTemplate = `Batch processing {ENTITY_COUNT} {ENTITY_TYPE} entities

Consistency Strategy: {CONSISTENCY_APPROACH}
Quality Threshold: {QUALITY_MINIMUM}
Performance Budget: {TIME_BUDGET}
Output Format: {OUTPUT_SPECIFICATION}

Ensure uniform quality across all entities while meeting performance targets.`

	// Real-Time Processing Template
	RealTimeTemplate = `Real-time transformation for {ENTITY_TYPE}

Latency Target: {MAX_LATENCY}ms
Quality Level: {QUALITY_TIER}
Cache Strategy: {CACHE_APPROACH}
Fallback Plan: {FALLBACK_STRATEGY}

Optimize for speed while maintaining acceptable quality thresholds.`
)

/*
USAGE EXAMPLES:

1. Product Transformation:
   prompt := constants.ProductVectorPrompt + "\n\n" + constants.EcommercePrompt

2. Multi-Strategy Job Processing:
   prompt := strings.Replace(constants.MultiStrategyPrompt, "{PRIMARY_STRATEGY}", constants.JobVectorPrompt, 1)
   prompt = strings.Replace(prompt, "{SECONDARY_STRATEGY}", constants.SemanticEnrichmentPrompt, 1)

3. Context-Aware Service Description:
   prompt := strings.Replace(constants.ContextAdaptationPrompt, "{BASE_PROMPT}", constants.ServiceVectorPrompt, 1)
   prompt = strings.Replace(prompt, "{CONTEXT_TYPE}", "B2B_ENTERPRISE", 1)

4. High-Performance Real-Time Processing:
   prompt := constants.PropertyVectorPrompt + "\n\n" + constants.RealTimeOptimizationPrompt

5. Cross-Language Vehicle Search:
   prompt := constants.VehicleVectorPrompt + "\n\n" + constants.CrossLanguagePrompt

INTEGRATION WITH EMBEDDING SERVICE:

func (e *EmbeddingService) GenerateEntityEmbeddingWithPrompt(
    ctx context.Context,
    entityType string,
    entityData map[string]interface{},
    promptStrategy string,
) ([]float32, error) {

    // Select appropriate base prompt
    basePrompt := e.getBasePromptForEntity(entityType)

    // Combine with strategy
    fullPrompt := basePrompt + "\n\n" + promptStrategy

    // Transform entity data using LLM with prompt
    transformedText := e.llmTransform(entityData, fullPrompt)

    // Generate embedding from transformed text
    return e.GenerateEmbedding(ctx, transformedText)
}
*/
