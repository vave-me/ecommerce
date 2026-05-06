#!/usr/bin/env node

/**
 * Fix broken function structures from console removal
 */

const fs = require('fs');
const path = require('path');

// Specific fixes for known broken files
const fixes = [
    {
        file: 'src/api/client/mediaApi.jsx',
        fixes: [
            {
                // Fix the first broken catch block around line 115
                find: `        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error details logged for debugging, ... ] }
    } catch (error) {`,
                replace: `        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error details logged for debugging
        }
        throw error;
    }
};

/**
 * GET ALL ITEM IMAGES
 * GET /api/media/item/{itemId}/image
 * Response -> mediapbGetAllItemImagesResponse { images: [...] }
 */
export const getAllItemImages = async (itemId) => {
    try {
        const response = await axiosInstance.get(\`/media/item/\${itemId}/image\`);
        // Returns { images: [...] }
        return response.data;
    } catch (error) {`
            },
            {
                // Fix the second broken catch block
                find: `        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error details logged for debugging, ... ] }
    } catch (error) {`,
                replace: `        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error details logged for debugging
        }
        throw error;
    }
};

/**
 * GET ALL ITEM VIDEOS
 * GET /api/media/item/{itemId}/video
 * Response -> mediapbGetAllItemVideosResponse { videos: [...] }
 */
export const getAllItemVideos = async (itemId) => {
    try {
        const response = await axiosInstance.get(\`/media/item/\${itemId}/video\`);
        // Returns { videos: [...] }
        return response.data;
    } catch (error) {`
            },
            {
                // Fix the comment before getMedia
                find: `        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error details logged for debugging }
 */`,
                replace: `        // Log other errors for debugging
        if (process.env.NODE_ENV === 'development') {
            // Error details logged for debugging
        }
        throw error;
    }
};

/**
 * GET MEDIA BY ID
 * GET /api/media/{mediaId}
 * Response -> mediapbGetMediaResponse { media: {...} }
 */`
            },
            {
                // Fix broken console statements
                find: `, returning empty array\`);`,
                replace: `// Server error - returning empty array`
            },
            {
                // Fix another broken console
                find: `, returning empty results\`);`,
                replace: `// Server error - returning empty results`
            }
        ]
    },
    {
        file: 'src/api/client/messagingApi.jsx',
        fixes: [
            {
                find: `        // Error details logged for debugging, or the entire conversation object.
 */`,
                replace: `        // Error details logged for debugging
        throw err;
    }
};

/**
 * Start a new conversation or get existing one by recipientId and itemId
 */`
            }
        ]
    },
    {
        file: 'src/api/client/productsApi.jsx',
        fixes: [
            {
                find: `        // Error details logged for debugging } or similar
    } catch (error) {`,
                replace: `        // Error details logged for debugging
        throw error;
    }
};

/**
 * GET PRODUCT BY ID
 */
export const getProductById = async (productId) => {
    const endpoint = \`/products/\${productId}\`;
    try {
        const response = await axiosInstance.get(endpoint);
        return response.data;
    } catch (error) {`
            }
        ]
    },
    {
        file: 'src/api/client/userApi.jsx',
        fixes: [
            {
                find: `            + '...');`,
                replace: `            // Token stored successfully`
            }
        ]
    }
];

function fixFile(filePath, fixes) {
    try {
        let content = fs.readFileSync(filePath, 'utf8');
        const originalContent = content;
        
        fixes.forEach(fix => {
            if (content.includes(fix.find)) {
                content = content.replace(fix.find, fix.replace);
                console.log(`✅ Applied fix in ${filePath}`);
            }
        });
        
        if (content !== originalContent) {
            fs.writeFileSync(filePath, content, 'utf8');
            console.log(`✅ Saved ${filePath}`);
            return true;
        }
        
        return false;
    } catch (error) {
        console.error(`❌ Error fixing ${filePath}:`, error.message);
        return false;
    }
}

console.log('🔧 Fixing broken function structures...\n');

let totalFixed = 0;

fixes.forEach(({ file, fixes: fileFixes }) => {
    const fullPath = path.join(process.cwd(), file);
    if (fixFile(fullPath, fileFixes)) {
        totalFixed++;
    }
});

console.log(`\n✨ Fixed ${totalFixed} files`);