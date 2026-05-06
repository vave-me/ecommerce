// app/videos/[itemId]/layout.jsx
// No "use client" here — so it’s a Server Component by default.
export const revalidate = 60;
export const metadata = {
    title: 'Short Video Feed',
    description: 'A vertical feed of short videos for a specific itemId.',
    openGraph: {
        // Any OG tags if you wish:
        title: 'Short Video Feed',
        description: 'A vertical feed of short videos for a specific itemId.'
    }
};
export default function Layout({ children }) {
    return (
        <section>
            {children}
        </section>
    );
}
//