export async function fetchApi<T>(
  endpoint: string,
  token: string | null,
  options: RequestInit = {}
): Promise<T> {
  const defaultHeaders: HeadersInit = {
    'Content-Type': 'application/json',
  };
  if (token) {
    defaultHeaders['Authorization'] = `Bearer ${token}`;
  }

  const config: RequestInit = {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
  };

  const response = await fetch(`http://localhost:8080${endpoint}`, config);

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ message: 'Unknown API error' }));
    throw new Error(errorData.message || `API error: ${response.status}`);
  }

  if (response.status === 204) {
    return null as T;
  }

  return response.json() as Promise<T>;
}