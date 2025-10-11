'use client';

import React, { useState } from 'react';
import { Icon } from './ui/Icon';
import { NotificationsPanel } from './NotificationsPanel';
import { Alert } from '@/lib/types';

// Composant Navbar (TopBar)
export const DashboardNavbar = ({ notifications }: { notifications: Alert[] | undefined }) => {
  const [isPanelOpen, setIsPanelOpen] = useState(false);
  const unreadCount = notifications?.length || 0;

  return (
    <nav className="fixed top-0 left-0 right-0 lg:left-64 z-10 h-16 bg-white border-b border-gray-200 shadow-sm flex items-center justify-between px-4 sm:px-8">
      <h1 className="text-xl font-bold text-gray-800 hidden lg:block">Tableau de Bord Super Admin</h1>
      <h1 className="text-xl font-bold text-gray-800 lg:hidden">Admin Dashboard</h1>

      <div className="flex items-center space-x-4">
        <div className="relative">
          <button
            onClick={() => setIsPanelOpen(!isPanelOpen)}
            className="p-2 rounded-full text-gray-600 hover:bg-gray-100 transition-colors"
          >
            <Icon name="Bell" />
            {unreadCount > 0 && (
              <span className="absolute top-0 right-0 block h-3 w-3 rounded-full ring-2 ring-white bg-red-500">
                <span className="text-xs text-white absolute -top-1 -right-1.5">{unreadCount}</span>
              </span>
            )}
          </button>
          {isPanelOpen && <NotificationsPanel notifications={notifications} />}
        </div>

        <div className="flex items-center space-x-2 cursor-pointer p-1 rounded-full hover:bg-gray-100 transition-colors">
            <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold text-sm">SA</div>
            <span className="text-sm font-medium text-gray-700 hidden sm:block">Super Admin</span>
        </div>
      </div>
    </nav>
  );
};
