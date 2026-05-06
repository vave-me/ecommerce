package tools

import ai2 "middleman/internal/ai"

func createCommentTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "comment_create",
				Description: "Create a comment",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"post_id": map[string]interface{}{
							"type":        "string",
							"description": "Post ID",
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "Comment content",
						},
						"author_id": map[string]interface{}{
							"type":        "string",
							"description": "Author ID",
						},
					},
					"required": []string{"post_id", "content", "author_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "comment_get_by_post",
				Description: "Get comments for a post",
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
	}
}