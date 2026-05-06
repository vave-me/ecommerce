#!/usr/bin/env python3
import os
import re

# Define which tools to keep for each category
essential_tools = {
    "basket": ["basket_add_item", "basket_remove_item", "basket_get", "basket_clear", "basket_checkout"],
    "order": ["order_create", "order_get_by_id", "order_get_user_orders", "order_update_status", "order_cancel"],
    "shipping": ["shipping_get_options", "shipping_calculate_cost", "shipping_create_label", "shipping_track"],
    "newsletter": ["newsletter_subscribe", "newsletter_unsubscribe", "newsletter_get_subscriptions"],
    "post": ["post_create", "post_get_by_id", "post_search", "post_update", "post_delete"],
    "service": ["service_create", "service_get_by_id", "service_search", "service_update"],
    "comment": ["comment_create", "comment_get_by_post", "comment_delete"],
    "wishlist": ["wishlist_add_item", "wishlist_remove_item", "wishlist_get"],
    "notification": ["notification_send", "notification_get_user_notifications", "notification_mark_read"],
    "mailer": ["mailer_send_email", "mailer_send_template"],
    "media": ["media_upload", "media_get_by_id", "media_delete"],
    "activity": ["activity_log", "activity_get_user_activities"],
}

def process_file(filepath):
    """Process a tool file to keep only essential methods"""
    category = filepath.split('/')[-1].replace('_tools.go', '')
    
    if category not in essential_tools:
        print(f"Skipping {category} - not in essential list")
        return
    
    keep_methods = essential_tools[category]
    
    with open(filepath, 'r') as f:
        content = f.read()
    
    # Find all tool definitions
    tool_pattern = r'{\s*Type:\s*"function",\s*Function:\s*ai2\.FunctionDef{[^}]*Name:\s*"([^"]+)"'
    
    lines = content.split('\n')
    new_lines = []
    in_tool = False
    tool_name = None
    brace_count = 0
    skip_tool = False
    
    i = 0
    while i < len(lines):
        line = lines[i]
        
        # Check if we're starting a new tool
        if 'Type: "function"' in line:
            in_tool = True
            brace_count = 0
            skip_tool = False
            tool_name = None
            
            # Look ahead to find the tool name
            for j in range(i, min(i+10, len(lines))):
                if 'Name:' in lines[j]:
                    match = re.search(r'Name:\s*"([^"]+)"', lines[j])
                    if match:
                        tool_name = match.group(1)
                        if tool_name not in keep_methods:
                            skip_tool = True
                            new_lines.append(f"\t\t// COMMENTED OUT - {tool_name} (non-essential)")
                            new_lines.append("\t\t/*")
                        break
        
        # Count braces to track when tool ends
        if in_tool:
            brace_count += line.count('{') - line.count('}')
            
        # Add the line
        if skip_tool:
            new_lines.append(line)
        else:
            new_lines.append(line)
            
        # Check if tool ended
        if in_tool and brace_count <= 0 and '},' in line:
            if skip_tool:
                new_lines.append("\t\t*/")
            in_tool = False
            skip_tool = False
            tool_name = None
            
        i += 1
    
    # Write back
    new_content = '\n'.join(new_lines)
    
    # Backup original
    os.rename(filepath, filepath + '.bak')
    
    # Write new content
    with open(filepath, 'w') as f:
        f.write(new_content)
    
    print(f"Processed {category}: kept {len(keep_methods)} essential tools")

# Process all tool files
for filename in os.listdir('.'):
    if filename.endswith('_tools.go') and not filename.endswith('_essential.go') and filename not in ['product_tools.go', 'category_tools.go', 'user_tools.go']:
        process_file(filename)