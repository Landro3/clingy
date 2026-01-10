import type { AxiosResponse } from 'axios';
import { useState, useEffect, useCallback } from 'react';

export function useQuery<T>(apiFunction: () => Promise<AxiosResponse>) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const query = useCallback(() => {
    apiFunction()
      .then((res) => setData(res.data))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [apiFunction]);

  useEffect(query, [query]);

  return { data, loading, error, refetch: query };
}

// TODO: revisit typeing, infer it from function argument and return type
export function useMutation<T, TVariables = any>(
  apiFunction: (variables: TVariables) => Promise<AxiosResponse>,
) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const mutate = useCallback(async (variables: TVariables) => {
    setLoading(true);
    setError(null);

    try {
      const res = await apiFunction(variables);
      setData(res.data);
      return res.data;
    } catch (err: any) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [apiFunction]);

  return { mutate, data, loading, error };
}

