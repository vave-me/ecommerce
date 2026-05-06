import { readFileSync } from 'fs';
import { join } from 'path';
import { NextResponse } from 'next/server';
// Define ads.txt content directly in case file read fails
const fallbackAdsTxt = 'google.com, pub-7872277873986607, DIRECT, f08c47fec0942fa0';
/**
 * Route handler for /ads.txt
 * This ensures the ads.txt file is served properly even if the static file approach fails
 */
export async function GET() {
  try {
    // Try to read from the public directory
    const filePath = join(process.cwd(), 'public', 'ads.txt');
    const content = readFileSync(filePath, 'utf8');
    return new NextResponse(content, {
      status: 200,
      headers: {
        'Content-Type': 'text/plain',
        'Cache-Control': 'public, max-age=86400', // Cache for 24 hours
      },
    });
  } catch (error) {
    // Use fallback content if file read fails
    return new NextResponse(fallbackAdsTxt, {
      status: 200,
      headers: {
        'Content-Type': 'text/plain',
        'Cache-Control': 'public, max-age=86400', // Cache for 24 hours
      },
    });
  }
} 