"use client";

import { useState } from 'react';
import { useParkings } from '@/lib/api-client';
import { Parking, ParkingStatus } from '@/lib/types';
import { Card } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { X } from 'lucide-react';

const MapViewPlaceholder = () => {
    const { data: parkings, isLoading } = useParkings();
    const [selectedParking, setSelectedParking] = useState<Parking | null>(null);

    if (isLoading) {
        return <Skeleton className="h-96 w-full rounded-2xl" />;
    }

    const getMarkerColor = (status: ParkingStatus) => {
        switch (status) {
            case 'Disponible': return 'bg-green-500';
            case 'Presque Plein': return 'bg-orange-500';
            case 'Complet': return 'bg-red-500';
            case 'Hors Ligne': return 'bg-gray-800';
            default: return 'bg-blue-500';
        }
    };

    return (
        <Card className="h-96 relative overflow-hidden p-0 rounded-2xl shadow-md">
            <div className="absolute inset-0 bg-gray-200 flex items-center justify-center">
                <p className="text-gray-500 font-medium">Simulation de Carte Interactive</p>
            </div>

            {parkings?.map((p, index) => (
                <div
                    key={p.id}
                    className="absolute flex flex-col items-center cursor-pointer transition-transform hover:scale-110"
                    style={{ top: `${15 + index * 12}%`, left: `${10 + (index % 5) * 18}%` }}
                    onClick={() => setSelectedParking(p)}
                >
                    <div className={`w-4 h-4 rounded-full ${getMarkerColor(p.statut)} ring-4 ring-white shadow-lg`} />
                    <span className="mt-1 text-xs font-medium text-gray-800 bg-white/80 backdrop-blur-sm px-2 py-0.5 rounded-full">{p.nom}</span>
                </div>
            ))}

            <div className={`absolute top-0 right-0 h-full bg-white transition-transform duration-300 shadow-2xl p-6 w-full sm:w-80 border-l border-gray-200 ${selectedParking ? 'translate-x-0' : 'translate-x-full'}`}>
                <div className="flex justify-between items-start">
                    <h3 className="text-xl font-bold truncate text-gray-800">{selectedParking?.nom}</h3>
                    <button onClick={() => setSelectedParking(null)} className="p-1 rounded-full hover:bg-gray-100 text-gray-500">
                        <X className="w-5 h-5" />
                    </button>
                </div>

                {selectedParking && (
                    <div className="mt-4 space-y-3 text-sm">
                        <p><strong>Ville:</strong> {selectedParking.ville}</p>
                        <p><strong>Statut:</strong> <span className={`font-semibold ${getMarkerColor(selectedParking.statut).replace('bg-', 'text-')}`}>{selectedParking.statut}</span></p>
                        <p><strong>Taux d'Occupation:</strong> <span className="font-bold text-lg">{selectedParking.tauxOccupation}%</span></p>
                        <p><strong>Admin Local:</strong> {selectedParking.adminLocal}</p>
                        <div className="h-2 w-full bg-gray-200 rounded-full mt-1">
                            <div className="h-2 bg-green-500 rounded-full" style={{ width: `${selectedParking.tauxOccupation}%` }}></div>
                        </div>
                    </div>
                )}
            </div>
        </Card>
    );
};

export default MapViewPlaceholder;
