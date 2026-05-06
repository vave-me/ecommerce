// File: CategorySection.server.jsx
import React from "react";
import CategoryGrid from "./CategoryGrid.client";
import {getCategories} from "../../api/categories";
export default async function CategorySectionServer() {
    // 1) Fetch categories on the server side
    const categories = await getCategories(); // your own API or data fetch
    // 2) Return the client component
    return <CategoryGrid categories={categories} showGoogleCategoryId={false}/>;
}
