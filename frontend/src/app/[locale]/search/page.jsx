// File: src/app/search/page.jsx
import { Suspense } from "react";
import SearchResults from "./SearchResults";
export default function SearchPage() {
    return (
        <Suspense fallback={<div>Loading search...</div>}>
            <SearchResults />
        </Suspense>
    );
}