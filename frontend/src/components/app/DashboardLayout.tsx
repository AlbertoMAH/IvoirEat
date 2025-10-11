"use client";

import React, { useState } from 'react';
import DashboardSidebar from './DashboardSidebar';
import DashboardNavbar from './DashboardNavbar';

const DashboardLayout = ({ children }: { children: React.ReactNode }) => {
    // Note: l'état activeView est déplacé ici pour être partagé si nécessaire,
    // mais pour l'instant, seul le DashboardSidebar l'utilise.
    const [activeView, setActiveView] = useState('dashboard');

    return (
        <div className="min-h-screen flex bg-gray-50">
            <DashboardSidebar activeView={activeView} setActiveView={setActiveView} />
            <div className="flex flex-col flex-1 lg:ml-64">
                <DashboardNavbar />
                <main className="flex-1 mt-16 p-8">
                    {children}
                </main>
            </div>
        </div>
    );
};

export default DashboardLayout;
