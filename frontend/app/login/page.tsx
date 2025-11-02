'use client';

import { useState, useEffect } from 'react'; // Dodaj useEffect
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const router = useRouter();
  const { login, token, isLoading } = useAuth(); // Dodaj token i isLoading

  // Efekt do przekierowania, jeśli użytkownik jest już zalogowany
  useEffect(() => {
    if (!isLoading && token) {
      router.push('/dashboard'); // Przekieruj do panelu
    }
  }, [token, isLoading, router]);


  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    // ... (logika handleSubmit bez zmian) ...
     event.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const response = await fetch('http://localhost:8080/users/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      console.log('Login successful, received token:', data.token);

      login(data.token);

      console.log('Token saved, redirecting...');
      router.push('/dashboard'); // Zmieniamy przekierowanie na /dashboard

    } catch (err: unknown) {
      console.error('Login failed:', err);
      let errorMessage = 'Wystąpił nieoczekiwany błąd logowania.';
      if (err instanceof Error) {
        errorMessage = err.message;
      } else if (typeof err === 'string') {
        errorMessage = err;
      }
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // Nie renderuj formularza, dopóki trwa sprawdzanie sesji i potencjalne przekierowanie
  if (isLoading || token) {
     return <div>Sprawdzanie sesji...</div>;
  }


  // ... (reszta komponentu - formularz JSX - bez zmian) ...
   return (
    <div className="flex min-h-screen items-center justify-center bg-gray-100">
      <div className="w-full max-w-md rounded-lg bg-white p-8 shadow-md">
        <h2 className="mb-6 text-center text-2xl font-bold text-gray-900">
          Zaloguj się do raitometer
        </h2>
        <form onSubmit={handleSubmit} className="space-y-6">
           <div>
            <label
              htmlFor="email"
              className="block text-sm font-medium text-gray-700"
            >
              Adres email
            </label>
            <input
              id="email"
              name="email"
              type="email"
              autoComplete="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
              placeholder="ty@example.com"
            />
          </div>

          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-gray-700"
            >
              Hasło
            </label>
            <input
              id="password"
              name="password"
              type="password"
              autoComplete="current-password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
              placeholder="••••••••"
            />
          </div>


          {error && (
            <div className="rounded-md bg-red-50 p-4">
              <p className="text-sm text-red-700">{error}</p>
            </div>
          )}

          <div>
            <button
              type="submit"
              disabled={loading}
              className={`flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 ${
                loading ? 'cursor-not-allowed opacity-50' : ''
              }`}
            >
              {loading ? 'Logowanie...' : 'Zaloguj się'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}