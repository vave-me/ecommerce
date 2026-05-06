package tools

import ai2 "middleman/internal/ai"

func createPostTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "post_create",
				Description: "Create a new post",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"title": map[string]interface{}{
							"type":        "string",
							"description": "Post title",
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "Post content",
						},
						"author_id": map[string]interface{}{
							"type":        "string",
							"description": "Author ID",
						},
					},
					"required": []string{"title", "content", "author_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "post_get_by_id",
				Description: "Get post by ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"post_id": map[string]interface{}{
							"type":        "string",
							"description": "Post ID",
						},
					},
					"required": []string{"post_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "post_search",
				Description: "Search posts",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query",
						},
					},
					"required": []string{"query"},
				},
			},
		},
	}
}