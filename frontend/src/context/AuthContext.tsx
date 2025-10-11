"use client";

import React, { createContext, useContext, useState, ReactNode } from 'react';

// Définir la structure du contexte
interface AuthContextType {
  isAuthenticated: boolean;
  login: (email: string, pass: string) => boolean;
  logout: () => void;
}

// Créer le contexte avec une valeur par défaut
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Créer le fournisseur (Provider) du contexte
export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // Logique de connexion (simulée)
  const login = (email: string, pass: string): boolean => {
    if (email === "falbertmah@gmail.com" && pass === "test") {
      setIsAuthenticated(true);
      return true;
    }
    return false;
  };

  // Logique de déconnexion
  const logout = () => {
    setIsAuthenticated(false);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

// Créer un hook personnalisé pour utiliser facilement le contexte
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
