package processor

import (
	"fmt"
	"middleman/assistants/internal/domain"
	ai2 "middleman/internal/ai"
	"strings"
	"sync"
	"time"
)

// LLMProcessor provides a full implementation of LLMProcessor using tool registry

// ToolPerformanceMetrics tracks tool success rates and patterns

// CachedResponse stores semantically similar responses
type CachedResponse struct {
	Response   string                   `json:"response"`
	Actions    []domain.AssistantAction `json:"actions"`
	Confidence float64                  `json:"confidence"`
	Timestamp  time.Time                `json:"timestamp"`
	QueryHash  string                   `json:"query_hash"`
	Similarity float64                  `json:"similarity"`
}

// ConversationPatternAnalyzer identifies conversation patterns
type ConversationPattern struct {
	Pattern     string  `json:"pattern"`
	SuccessRate float64 `json:"success_rate"`
	Frequency   int     `json:"frequency"`
}

// ConversationContext provides enhanced context for conversation analysis
type ConversationContext struct {
	Complexity      float64                      `json:"complexity"`
	Intent          string                       `json:"intent"`
	Entities        []string                     `json:"entities"`
	Patterns        []ConversationPattern        `json:"patterns"`
	RelevantHistory []domain.ConversationMessage `json:"relevant_history"`
	LastToolResult  string                       `json:"last_tool_result"`
	ErrorCount      int                          `json:"error_count"`
	IterationCount  int                          `json:"iteration_count"`
}

// ExecutionPlan defines optimized tool execution strategy
type ExecutionPlan struct {
	Batches  []ExecutionBatch `json:"batches"`
	Strategy string           `json:"strategy"`
}

type ExecutionBatch struct {
	ToolCalls       []ai2.ToolCall `json:"tool_calls"`
	OriginalIndices []int          `json:"original_indices"`
	Priority        int            `json:"priority"`
}

// ConversationMemoryManager handles intelligent conversation history management
type ConversationMemoryManager struct {
	maxTokens        int
	summaryThreshold int
	contextStore     map[string]*ConversationMemory
	mutex            sync.RWMutex
}

type ConversationMemory struct {
	ConversationID    string                       `json:"conversation_id"`
	FullHistory       []domain.ConversationMessage `json:"full_history"`
	KeyFacts          []string                     `json:"key_facts"`
	ConversationState map[string]interface{}       `json:"conversation_state"`
	LastSummary       string                       `json:"last_summary"`
	SummaryTimestamp  time.Time                    `json:"summary_timestamp"`
	TokenUsage        int                          `json:"token_usage"`
	RecentMessages    []domain.ConversationMessage `json:"recent_messages"`
}

type TokenCounter struct {
	// Simple token estimation - 1 token ≈ 4 characters for GPT models
	CharToTokenRatio float64
}

func NewConversationMemoryManager() *ConversationMemoryManager {
	return &ConversationMemoryManager{
		maxTokens:        15000, // Conservative limit for most models
		summaryThreshold: 30,    // Summarize when conversation exceeds 30 messages
		contextStore:     make(map[string]*ConversationMemory),
	}
}

func (t *TokenCounter) EstimateTokens(text string) int {
	if t.CharToTokenRatio == 0 {
		t.CharToTokenRatio = 0.25 // 1 token ≈ 4 characters
	}
	return int(float64(len(text)) * t.CharToTokenRatio)
}

func (cmm *ConversationMemoryManager) GetOptimalHistory(conversationID string, history []domain.ConversationMessage, currentMessage string) (*ConversationMemory, []domain.ConversationMessage, string) {
	cmm.mutex.Lock()
	defer cmm.mutex.Unlock()

	// Get or create conversation memory
	memory, exists := cmm.contextStore[conversationID]
	if !exists {
		memory = &ConversationMemory{
			ConversationID:    conversationID,
			FullHistory:       make([]domain.ConversationMessage, 0),
			KeyFacts:          make([]string, 0),
			ConversationState: make(map[string]interface{}),
			RecentMessages:    make([]domain.ConversationMessage, 0),
		}
		cmm.contextStore[conversationID] = memory
	}

	// Update full history
	memory.FullHistory = history
	memory.RecentMessages = history

	tokenCounter := &TokenCounter{}

	// Calculate optimal history within token limits
	optimalHistory, contextSummary := cmm.calculateOptimalHistory(memory, currentMessage, tokenCounter)

	// Update memory with key facts from current interaction
	cmm.extractAndStoreKeyFacts(memory, currentMessage)

	return memory, optimalHistory, contextSummary
}

func (cmm *ConversationMemoryManager) calculateOptimalHistory(memory *ConversationMemory, currentMessage string, tokenCounter *TokenCounter) ([]domain.ConversationMessage, string) {
	availableTokens := cmm.maxTokens
	currentMessageTokens := tokenCounter.EstimateTokens(currentMessage)
	availableTokens -= currentMessageTokens

	// Reserve tokens for system prompt (estimated)
	systemPromptTokens := 2000
	availableTokens -= systemPromptTokens

	// Build context summary
	contextSummary := cmm.buildContextSummary(memory)
	contextTokens := tokenCounter.EstimateTokens(contextSummary)
	availableTokens -= contextTokens

	// Include as many recent messages as possible within token limit
	selectedHistory := make([]domain.ConversationMessage, 0)
	totalTokens := 0

	// Start from most recent messages and work backwards
	for i := len(memory.FullHistory) - 1; i >= 0; i-- {
		message := memory.FullHistory[i]
		messageTokens := tokenCounter.EstimateTokens(message.Content)

		if totalTokens+messageTokens > availableTokens {
			break
		}

		// Prepend to maintain chronological order
		selectedHistory = append([]domain.ConversationMessage{message}, selectedHistory...)
		totalTokens += messageTokens
	}

	return selectedHistory, contextSummary
}

func (cmm *ConversationMemoryManager) buildContextSummary(memory *ConversationMemory) string {
	summary := fmt.Sprintf("📊 CONVERSATION CONTEXT SUMMARY for conversation %s:\n", memory.ConversationID)

	// Add conversation statistics
	summary += fmt.Sprintf("• Total messages in conversation: %d\n", len(memory.FullHistory))
	summary += fmt.Sprintf("• Conversation started: %s\n",
		memory.FullHistory[0].Timestamp.Format("2006-01-02 15:04"))

	// Add key facts
	if len(memory.KeyFacts) > 0 {
		summary += "\n🔑 KEY FACTS & CONTEXT:\n"
		for i, fact := range memory.KeyFacts {
			if i >= 10 { // Limit to most important facts
				break
			}
			summary += fmt.Sprintf("• %s\n", fact)
		}
	}

	// Add conversation state
	if len(memory.ConversationState) > 0 {
		summary += "\n📋 CONVERSATION STATE:\n"
		for key, value := range memory.ConversationState {
			summary += fmt.Sprintf("• %s: %v\n", key, value)
		}
	}

	// Add previous summary if available
	if memory.LastSummary != "" && !memory.SummaryTimestamp.IsZero() {
		summary += fmt.Sprintf("\n📝 PREVIOUS CONTEXT (from %s):\n%s\n",
			memory.SummaryTimestamp.Format("2006-01-02 15:04"), memory.LastSummary)
	}

	summary += "\n" + strings.Repeat("=", 50) + "\n"
	return summary
}

func (cmm *ConversationMemoryManager) extractAndStoreKeyFacts(memory *ConversationMemory, currentMessage string) {
	// Simple keyword-based fact extraction (can be enhanced with AI later)
	keywords := []string{
		"my name is", "i am", "i work", "i live", "i need", "i want",
		"looking for", "price", "budget", "prefer", "don't like", "favorite",
		"remember", "important", "note that", "keep in mind",
	}

	lowerMessage := strings.ToLower(currentMessage)
	for _, keyword := range keywords {
		if strings.Contains(lowerMessage, keyword) {
			// Extract sentence containing the keyword
			sentences := strings.Split(currentMessage, ".")
			for _, sentence := range sentences {
				if strings.Contains(strings.ToLower(sentence), keyword) {
					fact := strings.TrimSpace(sentence)
					if len(fact) > 10 && !cmm.containsFact(memory.KeyFacts, fact) {
						memory.KeyFacts = append(memory.KeyFacts, fact)
						// Keep only most recent 20 facts
						if len(memory.KeyFacts) > 20 {
							memory.KeyFacts = memory.KeyFacts[1:]
						}
					}
					break
				}
			}
		}
	}
}

// ConversationPatternAnalyzer analyzes conversation patterns for optimization
type ConversationPatternAnalyzer struct {
	patterns map[string]*ConversationPattern
	mutex    sync.RWMutex
}

func NewConversationPatternAnalyzer() *ConversationPatternAnalyzer {
	return &ConversationPatternAnalyzer{
		patterns: make(map[string]*ConversationPattern),
	}
}

// ... existing code ...

// containsFact checks if a fact is already stored
func (cmm *ConversationMemoryManager) containsFact(facts []string, newFact string) bool {
	newFactLower := strings.ToLower(newFact)
	for _, fact := range facts {
		if strings.ToLower(fact) == newFactLower {
			return true
		}
	}
	return false
}

// UNIFIED SYSTEM PROMPT GENERATION - Single source of truth
