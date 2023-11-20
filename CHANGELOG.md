# Changelog

## [1.0.1] - 11/20/23

This release focuses on implementing case-insensitive search and pagination for the search feature. Additionally, a basic error UI and error UI test was added to the frontend.

### Backend

#### 1. TestSearchCaseSensitive
This test asserts that the search function is case-insensitive, such that searching "hAmLeT" will return search results for "hamlet".

In the `Load` method of the `Searcher` struct, both the original text (`CompleteWorks`) and a lowercase version of the text (`LowerCompleteWorks`) are stored. This step preprocesses the test to be lowercase, which is used later for case-insensitive searching.
```
s.CompleteWorks  =  string(dat)
s.LowerCompleteWorks  = strings.ToLower(s.CompleteWorks)
```

Next, a suffix array is created from the lowercase version of the complete works (`LowerCompleteWorks`). This is used for efficient substring searching.
```
s.SuffixArray  = suffixarray.New([]byte(s.LowerCompleteWorks))
```

In the `Search` method, the query is converted to lowercase before performing the search, ensuring that the search is case-insensitive because both the text in the suffix array and query are lowercase. The `Lookup` method of the suffix array is then used to find occurrences of the lowercase query within the lowercase complete works.
```
lowerQuery  := strings.ToLower(query)
idxs  := s.SuffixArray.Lookup([]byte(lowerQuery), -1)
```

Finally, after finding the indexes where the query matches, the method extracts the surrounding test from the original `CompleteWorks` (vs. the original version). This ensures that the original casing of the text is preserved in the search results.
```
start  :=  max(0, idx-250)
end  :=  min(len(s.CompleteWorks), idx+250)
results  =  append(results, s.CompleteWorks[start:end])
```

By processing both the complete works and the query in lowercase, the search becomes insensitive to the case of the letters, allowing, for example, a search for "hamlet" to match occurrences of "Hamlet", "HAMLET", etc.

#### 2. TestSearchDrunk
This test asserts that the search function limits the search results to 20.

The solution here is to default the search `limit` to 20 (in `handleSearch`). Basic pagination has been added, so now we have an `offset` and `limit` argument passed to the method. If the client doesn't specify a `limit`, or provides an invalid one, the server will default to returning 20 results.

### Frontend
#### 1. Load More feature
**JavaScript changes**
Pagination state management: `currentPaginationOffset` is used to keep track of the current offset for search results. It starts at 0 and is incremented by `PAGE_SIZE` each time the "Load More" button is clicked. `PAGE_SIZE` is set to 20, indicating the number of results to load per page. In a future release, we should update this to use a more robust state management system.

`loadMore` function: when the "Load More" button is clicked, `loadMore` is called. It increments `currentPaginationOffset` and calls `fetchAndDisplayResults` with the updated offset. This requests the next set of results from the server starting from the new offset.

`fetchAndDisplayResults` function: makes a request to the server's `/search` endpoint with the query, current offset, and limit (20). It then updates the table with new results. If it's a fresh search (i.e., offset is 0), it replaces the table's contents, else it appends to the existing content.

**Go changes**
`handleSearch` function: extracts the `offset` and `limit` query params from incoming requests, adding reasonable defaults if needed.
`Search` method: now accepts the `offset` and `limit` parameters, and slices the results based on the new parameters (defaulting the limit to 20). These new parameters allow for paginated search results, only returning a slice of results in each response.