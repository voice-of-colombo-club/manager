import { ApolloError, DocumentNode, useQuery } from "@apollo/client";
import { useCallback, useEffect, useRef, useState } from "react";

type PaginationQueryReturn<T> = [T[], boolean, (ApolloError | undefined), (node: any) => void, (string) => void];

export function usePaginationQuery<T>(query: DocumentNode, key: string, pageSize: number = 100, acceptDataParam = true) : PaginationQueryReturn<T> {
  const [list, setList] = useState([] as T[])
  const [hasMore, setHasMore] = useState(true)
  const [pageNumber, setPageNumber] = useState(0)
  const [searchCriteria, setSearchCriteria] = useState({});
  const [acceptData, setAcceptData] = useState(acceptDataParam)

  const processedSearchCriteria = { ...searchCriteria, offset: (pageNumber * pageSize), limit: pageSize };
  // TODO : Performance optimize to useLazyQuery
  const { loading, error, data } = useQuery(query, {
    fetchPolicy: "network-only", 
    variables: processedSearchCriteria
  });

  useEffect(() => {
    if (acceptData) {
      setPaginationData();
    }
  }, [data, acceptData])

  function setPaginationData() {
    if(!data) return;
    const dataList = data[key];
    setList((prevData: T[]) => [...prevData, ...dataList])
    setHasMore(dataList.length > 0)
  }

  function updateWithSearchCriteria(searchCriteria) {
    setSearchCriteria(searchCriteria)
    setPageNumber(0)
    setList([])
    setAcceptData(true);
  }

  const [lastElementRef] = usePaginationDetection(loading, hasMore, setPageNumber)
  
  return ([list, loading, error, lastElementRef, updateWithSearchCriteria]);
}

function usePaginationDetection(loading: boolean, hasMore: boolean, setPageNumber: any) {
  const observer: any = useRef()
  const lastElementRef = useCallback(node => {
    if (loading) return;
    if (observer.current) {
      observer.current.disconnect();
    }
    observer.current = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && hasMore) {
        setPageNumber((prevPageNumber: number) => prevPageNumber + 1)
      }
    });
    if (node) {
      observer.current.observe(node)
    }
  }, [loading, hasMore])

  return [lastElementRef]
}