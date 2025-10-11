'use client';

import React from 'react';
import { useParkings } from '@/lib/api-client';
import { Card } from './components/ui/Card';
import { KPIStats } from './components/KPIStats';
import { MapViewPlaceholder } from './components/MapView';
import { ParkingTable } from './components/ParkingTable';
import { StatsDashboard } from './components/StatsDashboard';

// --- DASHBOARD PAGE COMPONENT ---

// La prop `activeView` est maintenant passée depuis le `DashboardLayout`
const DashboardPage = ({ activeView }: { activeView: string }) => {
    const { data: parkingsData } = useParkings();

    const renderContent = () => {
        switch (activeView) {
            case 'dashboard':
            case 'parkings': // Pour la démo, nous affichons la carte et la table sur le 'dashboard'
                return (
                    <div className="space-y-8">
                        {/* KPI Stats (Maintenant 3 cartes) */}
                        <KPIStats />

                        {/* Map View */}
                        <MapViewPlaceholder parkings={parkingsData} />

                        {/* Parking Table */}
                        <ParkingTable />
                    </div>
                );
            case 'stats':
                return <StatsDashboard />;
            case 'users':
                return <Card className="h-64 flex items-center justify-center text-gray-500">Vue de la gestion des Utilisateurs (API: /api/users)</Card>;
            case 'settings':
                return <Card className="h-64 flex items-center justify-center text-gray-500">Vue des Paramètres Globaux</Card>;
            default:
                return <div className="p-8 text-center text-gray-500">Sélectionnez une vue.</div>;
        }
    };

    return (
        <div className="flex-1 p-4 sm:p-8 pt-8 lg:pt-0">
            {renderContent()}
        </div>
    );
};

export default DashboardPage;
