package tools

import ai2 "middleman/internal/ai"

func createReviewTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_create_new",
				Description: "Create a new review for an item",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"sender_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user leaving the review",
						},
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the item being reviewed",
						},
						"item_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of item being reviewed",
							"enum":        []string{"product", "service", "order", "seller"},
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "Review content",
							"minLength":   10,
							"maxLength":   2000,
						},
						"category_id": map[string]interface{}{
							"type":        "string",
							"description": "Category ID for the review",
						},
						"parent_id": map[string]interface{}{
							"type":        "string",
							"description": "Parent review ID if this is a reply",
						},
						"rating": map[string]interface{}{
							"type":        "integer",
							"description": "Rating score",
							"minimum":     1,
							"maximum":     5,
						},
					},
					"required": []string{"sender_id", "item_id", "item_type", "content", "rating"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_get_by_id",
				Description: "Get a review by its ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"review_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the review to retrieve",
						},
					},
					"required": []string{"review_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_get_all_for_item",
				Description: "Get all reviews for a specific item",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the item to get reviews for",
						},
						"item_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of the item",
							"enum":        []string{"product", "service", "order", "seller"},
						},
					},
					"required": []string{"item_id", "item_type"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_edit_content",
				Description: "Edit the content of an existing review",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"review_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the review to edit",
						},
						"new_content": map[string]interface{}{
							"type":        "string",
							"description": "New review content",
							"minLength":   10,
							"maxLength":   2000,
						},
					},
					"required": []string{"review_id", "new_content"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_delete_by_id",
				Description: "Delete a review by its ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"review_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the review to delete",
						},
					},
					"required": []string{"review_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_approve_by_id",
				Description: "Approve a review by moderator",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"review_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the review to approve",
						},
					},
					"required": []string{"review_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_reject_by_id",
				Description: "Reject a review by moderator",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"review_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the review to reject",
						},
						"reason": map[string]interface{}{
							"type":        "string",
							"description": "Reason for rejection",
						},
					},
					"required": []string{"review_id", "reason"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_flag_as_inappropriate",
				Description: "Flag a review as inappropriate",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"review_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the review to flag",
						},
						"reason": map[string]interface{}{
							"type":        "string",
							"description": "Reason for flagging",
							"enum":        []string{"spam", "offensive", "irrelevant", "fake", "other"},
						},
					},
					"required": []string{"review_id", "reason"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_get_by_sender_id",
				Description: "Get all reviews written by a specific user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"sender_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user who wrote the reviews",
						},
					},
					"required": []string{"sender_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_get_all_approved",
				Description: "Get all approved reviews",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_get_most_reviewed_items",
				Description: "Get the most reviewed items across all categories",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_get_most_reviewed_by_category",
				Description: "Get the most reviewed items in a specific category",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the category",
						},
					},
					"required": []string{"category_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_update_content_by_id",
				Description: "Update the content and rating of a review",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"review_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the review to update",
						},
						"new_content": map[string]interface{}{
							"type":        "string",
							"description": "Updated review content",
							"minLength":   10,
							"maxLength":   2000,
						},
						"new_rating": map[string]interface{}{
							"type":        "integer",
							"description": "Updated rating score",
							"minimum":     1,
							"maximum":     5,
						},
					},
					"required": []string{"review_id", "new_content", "new_rating"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_remove_permanently",
				Description: "Permanently remove a review from the system",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"review_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the review to permanently remove",
						},
					},
					"required": []string{"review_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_unflag_by_id",
				Description: "Remove inappropriate flag from a review",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"review_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the review to unflag",
						},
					},
					"required": []string{"review_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "review_search_by_keyword",
				Description: "Search reviews by keyword",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type":        "string",
							"description": "Keyword to search for in reviews",
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     50,
						},
					},
					"required": []string{"keyword"},
				},
			},
		},
	}
}