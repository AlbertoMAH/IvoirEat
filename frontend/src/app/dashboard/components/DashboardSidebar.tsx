'use client';

import React from 'react';
import { Icon } from './ui/Icon';

// Composant Sidebar
export const DashboardSidebar = ({ activeView, setActiveView }: { activeView: string, setActiveView: (view: string) => void }) => {
  const navItems = [
    { name: 'Tableau de Bord', icon: 'BarChart', view: 'dashboard' },
    { name: 'Parkings', icon: 'Map', view: 'parkings' },
    { name: 'Utilisateurs', icon: 'Users', view: 'users' },
    { name: 'Statistiques', icon: 'BarChart', view: 'stats' },
    { name: 'Paramètres', icon: 'Settings', view: 'settings' },
  ];

  return (
    <div className="hidden lg:flex flex-col w-64 border-r border-gray-200 bg-white h-full fixed top-0 left-0 pt-20">
      <div className="p-4 flex flex-col space-y-2">
        {navItems.map((item) => (
          <button
            key={item.view}
            onClick={() => setActiveView(item.view)}
            className={`flex items-center space-x-3 p-3 rounded-xl transition-colors ${
              activeView === item.view
                ? 'bg-blue-600 text-white font-semibold shadow-md'
                : 'text-gray-600 hover:bg-gray-100'
            }`}
          >
            <Icon name={item.icon} className="w-5 h-5" />
            <span>{item.name}</span>
          </button>
        ))}
      </div>
      <div className="mt-auto p-4 border-t border-gray-100">
        <button className="flex items-center space-x-3 p-3 rounded-xl text-gray-500 hover:bg-red-50 hover:text-red-600 w-full transition-colors">
            <Icon name="LogOut" className="w-5 h-5" />
            <span>Déconnexion</span>
        </button>
      </div>
    </div>
  );
};
