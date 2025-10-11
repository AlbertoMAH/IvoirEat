'use client';

import React, { useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { DashboardSidebar } from './components/DashboardSidebar';
import { DashboardNavbar } from './components/DashboardNavbar';
import { useAlerts } from '@/lib/api-client'; // Assurez-vous que le chemin d'importation est correct

// Créez une seule instance de QueryClient
const queryClient = new QueryClient();

// Le composant "racine" pour le layout du dashboard
const DashboardLayoutContent = ({ children }: { children: React.ReactNode }) => {
    const [activeView, setActiveView] = useState('dashboard');
    const { data: alerts } = useAlerts(); // Récupère les alertes pour la navbar

    return (
        <div className="min-h-screen flex bg-gray-50 font-sans">
            <DashboardSidebar activeView={activeView} setActiveView={setActiveView} />

            <div className="flex flex-col flex-1 lg:ml-64">
                <DashboardNavbar notifications={alerts} />

                <main className="flex-1 mt-16 pb-8">
                    {/* Ici, nous passons une version modifiée des enfants pour leur injecter les props nécessaires */}
                    {React.Children.map(children, child => {
                        if (React.isValidElement(child)) {
                            // @ts-ignore
                            return React.cloneElement(child, { activeView, setActiveView });
                        }
                        return child;
                    })}
                </main>
            </div>
        </div>
    );
};


// Le layout principal qui enveloppe tout avec le QueryProvider
export default function DashboardLayout({ children }: { children: React.ReactNode }) {
    return (
        <QueryClientProvider client={queryClient}>
            <DashboardLayoutContent>{children}</DashboardLayoutContent>
        </QueryClientProvider>
    );
}
