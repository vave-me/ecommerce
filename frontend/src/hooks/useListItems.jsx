// File: src/hooks/useListItems.js
import { useSelector } from 'react-redux';
import { useInfiniteQuery } from '@tanstack/react-query';
import {fetchProductsByFilters} from "../api/productsApi";
import {fetchPostsByFilters} from "../api/postsApi";
export default function useListItems() {
    const listingFilters = useSelector((state) => state.listingFilters);
    const { listingType, ...filters } = listingFilters;
    const pageSize = 10; // or whatever default you want
    const {
        data,
        isLoading,
        isError,
        error,
        fetchNextPage,
        hasNextPage,
        isFetching,
    } = useInfiniteQuery({
        queryKey: ['listItems', listingType, filters],
        queryFn: async ({ pageParam = 1 }) => {
            if (listingType === 'products') {
                return fetchProductsByFilters({
                    ...filters,
                    page: pageParam,
                    pageSize,
                });
            } else {
                // 'posts', 'cars', 'jobs', etc.
                return fetchPostsByFilters({
                    ...filters,
                    page: pageParam,
                    pageSize,
                });
            }
        },
        getNextPageParam: (lastPage, allPages) => {
            // e.g. if your API returns totalCount
            const totalPages = Math.ceil(lastPage.totalCount / pageSize);
            const nextPage = allPages.length + 1;
            return nextPage <= totalPages ? nextPage : undefined;
        },
        keepPreviousData: true,
        staleTime: 60000,
    });
    // Flatten the data into a single items array:
    let items = [];
    if (data?.pages?.length) {
        if (listingType === 'products') {
            items = data.pages.flatMap((p) => p.products || []);
        } else {
            items = data.pages.flatMap((p) => p.posts || []);
        }
    }
    return {
        items,
        isLoading,
        isError,
        error,
        isFetching,
        fetchNextPage,
        hasNextPage,
    };
}
