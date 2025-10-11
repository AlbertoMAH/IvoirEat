'use client';

import React, { useState } from 'react';
import { Parking, ParkingStatus } from '@/lib/types';
import { Card } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Icon } from './ui/Icon';

// MapView - Simulation
export const MapViewPlaceholder = ({ parkings }: { parkings: Parking[] | undefined }) => {
    const [selectedParking, setSelectedParking] = useState<Parking | null>(null);

    if (!parkings) {
        return <Skeleton className="h-96 w-full" />;
    }

    // Détermine la couleur du marqueur
    const getMarkerColor = (status: ParkingStatus) => {
        switch (status) {
            case 'Disponible': return 'bg-green-500'; // 🟢
            case 'Presque Plein': return 'bg-orange-500'; // 🟠
            case 'Complet': return 'bg-red-500'; // 🔴
            case 'Hors Ligne': return 'bg-gray-800'; // ⚫
            default: return 'bg-blue-500';
        }
    };

    return (
        <Card className="h-96 relative overflow-hidden p-0">
            {/* Carte Visuelle Simplifiée (fond de carte stylisé) */}
            <div className="absolute inset-0 bg-gray-100 bg-cover bg-center" style={{ backgroundImage: 'url("https://placehold.co/1200x500/E5E7EB/4B5563?text=Vue+Cartographique+Simplifiée")' }}>
                <div className="absolute inset-0 bg-blue-500 opacity-10"></div>

                {/* Marqueurs de Parkings (basé sur une position relative simple) */}
                {parkings.map((p, index) => (
                    <div
                        key={p.id}
                        className={`absolute flex flex-col items-center cursor-pointer transition-transform hover:scale-110`}
                        // Positionnement fictif mais visuellement distinct
                        style={{ top: `${15 + index * 10}%`, left: `${10 + (index % 3) * 30}%` }}
                        onClick={() => setSelectedParking(p)}
                    >
                        <div className={`w-4 h-4 rounded-full ${getMarkerColor(p.statut)} ring-4 ring-white shadow-lg`} />
                        <span className="mt-1 text-xs font-medium text-gray-800 bg-white/80 backdrop-blur-sm px-2 py-0.5 rounded-full">{p.nom}</span>
                    </div>
                ))}
            </div>

            {/* Panneau Latéral (Side Panel) */}
            <div className={`absolute top-0 right-0 h-full bg-white transition-all duration-300 shadow-2xl p-6 w-full sm:w-80 border-l border-gray-200 ${selectedParking ? 'translate-x-0' : 'translate-x-full'}`}>
                <div className="flex justify-between items-start">
                    <h3 className="text-xl font-bold truncate text-gray-800">{selectedParking?.nom}</h3>
                    <button onClick={() => setSelectedParking(null)} className="p-1 rounded-full hover:bg-gray-100 text-gray-500">
                        <Icon name="X" className="w-5 h-5" />
                    </button>
                </div>

                {selectedParking && (
                    <div className="mt-4 space-y-3 text-sm">
                        <p><strong>Ville:</strong> {selectedParking.ville}</p>
                        <p><strong>Statut:</strong> <span className={`font-semibold ${getMarkerColor(selectedParking.statut).replace('bg-', 'text-')}`}>{selectedParking.statut}</span></p>
                        <p><strong>Taux d'Occupation:</strong> <span className="font-bold text-lg">{selectedParking.tauxOccupation}%</span></p>
                        <p><strong>Admin Local:</strong> {selectedParking.adminLocal}</p>
                        <div className="h-2 w-full bg-gray-200 rounded-full mt-1">
                            <div className="h-2 rounded-full" style={{ width: `${selectedParking.tauxOccupation}%`, backgroundColor: selectedParking.tauxOccupation >= 90 ? '#ef4444' : selectedParking.tauxOccupation >= 70 ? '#f97316' : '#22c55e' }}></div>
                        </div>
                        <p className="pt-2 text-xs text-gray-500">Coordonnées: {selectedParking.coordonnees.lat.toFixed(3)}, {selectedParking.coordonnees.lng.toFixed(3)}</p>
                    </div>
                )}
            </div>
        </Card>
    );
};
