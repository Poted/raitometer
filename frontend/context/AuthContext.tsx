'use client';

import React, { createContext, useState, useContext, useEffect, ReactNode } from 'react';

interface AuthContextType {
  token: string | null;
  login: (newToken: string) => void;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true); 

  useEffect(() => {
    try {
      const storedToken = localStorage.getItem('authToken');
      if (storedToken) {
        setToken(storedToken);
      }
    } catch (error) {
       console.error("Could not access localStorage:", error);
    } finally {
        setIsLoading(false);
    }
  }, []);

  const login = (newToken: string) => {
    setToken(newToken);
    try {
        localStorage.setItem('authToken', newToken);
    } catch (error) {
        console.error("Could not write to localStorage:", error);
    }

  };

  const logout = () => {
    setToken(null);
     try {
        localStorage.removeItem('authToken');
    } catch (error) {
        console.error("Could not remove item from localStorage:", error);
    }
  };

  if (isLoading) {
    return <div>Loading session...</div>;
  }


  return (
    <AuthContext.Provider value={{ token, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};