"use client";

import { useState } from 'react';
import { useAlerts } from '@/lib/api-client';
import { Bell } from 'lucide-react';
import NotificationsPanel from './dashboard/NotificationsPanel';

const DashboardNavbar = () => {
    const { data: alerts } = useAlerts();
    const [isPanelOpen, setIsPanelOpen] = useState(false);
    const unreadCount = alerts?.length || 0;

  return (
    <header className="fixed top-0 left-0 right-0 lg:left-64 z-10 h-16 bg-white border-b border-gray-200 shadow-sm flex items-center justify-between px-8">
      <h1 className="text-xl font-bold text-gray-800 hidden lg:block">Tableau de Bord Super Admin</h1>
      <div className="flex items-center space-x-4">
        <div className="relative">
          <button
            onClick={() => setIsPanelOpen(!isPanelOpen)}
            className="p-2 rounded-full text-gray-600 hover:bg-gray-100 transition-colors"
          >
            <Bell />
            {unreadCount > 0 && (
              <span className="absolute top-0 right-0 block h-2 w-2 rounded-full ring-2 ring-white bg-red-500" />
            )}
          </button>
          {isPanelOpen && <NotificationsPanel />}
        </div>
        <div className="flex items-center space-x-2">
            <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold text-sm">SA</div>
            <span className="text-sm font-medium text-gray-700">Super Admin</span>
        </div>
      </div>
    </header>
  );
};

export default DashboardNavbar;
