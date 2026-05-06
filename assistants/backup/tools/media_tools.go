package tools

import ai2 "middleman/internal/ai"

func createMediaTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "media_upload",
				Description: "Upload media file",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"filename": map[string]interface{}{
							"type":        "string",
							"description": "File name",
						},
						"content_type": map[string]interface{}{
							"type":        "string",
							"description": "Content type",
						},
					},
					"required": []string{"filename", "content_type"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "media_get_by_id",
				Description: "Get media by ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"media_id": map[string]interface{}{
							"type":        "string",
							"description": "Media ID",
						},
					},
					"required": []string{"media_id"},
				},
			},
		},
	}
}