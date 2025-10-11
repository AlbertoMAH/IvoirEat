'use client';

import React, { useState, useMemo } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useParkings, mockAddParking } from '@/lib/api-client';
import { Parking, ParkingStatus } from '@/lib/types';
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Icon } from './ui/Icon';
import { Skeleton } from '@/components/ui/skeleton';
import { AddParkingModal } from './AddParkingModal';

// ParkingTable
export const ParkingTable = () => {
    const queryClient = useQueryClient();
    const { data: parkings, isLoading, isError } = useParkings();
    const [searchTerm, setSearchTerm] = useState('');
    const [currentPage, setCurrentPage] = useState(1);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const itemsPerPage = 5;

    // Logique pour l'ajout d'un nouveau parking
    const handleAddParking = async (newParkingData: any) => {
        try {
            await mockAddParking(newParkingData);
            // Invalider la requête pour forcer un re-fetch et mettre à jour la liste
            await queryClient.invalidateQueries({ queryKey: ['parkings'] });
        } catch (error) {
            // Gérer l'erreur
            console.error("Échec de l'ajout du parking via le formulaire.", error);
        }
    };

    // Logique de recherche et de pagination
    const filteredParkings = useMemo(() => {
        if (!parkings) return [];
        return parkings.filter(p =>
            p.nom.toLowerCase().includes(searchTerm.toLowerCase()) ||
            p.ville.toLowerCase().includes(searchTerm.toLowerCase()) ||
            p.adminLocal.toLowerCase().includes(searchTerm.toLowerCase())
        );
    }, [parkings, searchTerm]);

    const paginatedParkings = useMemo(() => {
        const start = (currentPage - 1) * itemsPerPage;
        const end = start + itemsPerPage;
        return filteredParkings.slice(start, end);
    }, [filteredParkings, currentPage, itemsPerPage]);

    const totalPages = Math.ceil(filteredParkings.length / itemsPerPage);

    const getStatusStyle = (status: ParkingStatus) => {
        switch (status) {
            case 'Disponible': return 'text-green-600 bg-green-50';
            case 'Presque Plein': return 'text-orange-600 bg-orange-50';
            case 'Complet': return 'text-red-600 bg-red-50';
            case 'Hors Ligne': return 'text-gray-600 bg-gray-100';
            default: return 'text-blue-600 bg-blue-50';
        }
    };

    const handleAction = (parking: Parking, action: string) => {
        console.log(`Action '${action}' demandée pour le parking: ${parking.nom}`);
        // Logique d'appel API pour l'action (e.g., /api/parkings/p1/disable)
    };

    return (
        <>
            <AddParkingModal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                onSave={handleAddParking}
            />
            <Card>
                <CardHeader>
                    <div className="flex justify-between items-center flex-wrap gap-4">
                        <CardTitle>Liste des Parkings Partenaires ({parkings?.length || 0})</CardTitle>
                        <div className="flex items-center space-x-4">
                            <input
                                type="text"
                                placeholder="Rechercher par Nom, Ville..."
                                className="p-2 border border-gray-300 rounded-lg shadow-sm text-sm focus:ring-blue-500 focus:border-blue-500 w-full sm:w-48 transition-all"
                                value={searchTerm}
                                onChange={(e) => {
                                    setSearchTerm(e.target.value);
                                    setCurrentPage(1); // Reset page on search
                                }}
                            />
                            <button
                                onClick={() => setIsModalOpen(true)}
                                className="flex items-center space-x-2 px-4 py-2 text-sm font-medium rounded-xl text-white bg-green-600 hover:bg-green-700 transition-colors shadow-md"
                            >
                                <Icon name="Plus" className="w-5 h-5" />
                                <span>Ajouter Parking</span>
                            </button>
                        </div>
                    </div>
                </CardHeader>
                <CardContent>
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead className="bg-gray-50">
                                <tr>
                                    {['Nom', 'Ville', 'Statut', 'Occupation', 'Admin Local', 'Actions'].map((header) => (
                                        <th key={header} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{header}</th>
                                    ))}
                                </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                                {isLoading ? (
                                    Array(itemsPerPage).fill(0).map((_, i) => (
                                        <tr key={i}><td colSpan={6} className="px-6 py-4 whitespace-nowrap"><Skeleton className="h-6 w-full" /></td></tr>
                                    ))
                                ) : isError ? (
                                    <tr><td colSpan={6} className="px-6 py-4 text-center text-red-500">Erreur lors du chargement des parkings.</td></tr>
                                ) : paginatedParkings.length === 0 ? (
                                    <tr><td colSpan={6} className="px-6 py-4 text-center text-gray-500">Aucun résultat trouvé.</td></tr>
                                ) : (
                                    paginatedParkings.map((parking) => (
                                        <tr key={parking.id} className="hover:bg-gray-50">
                                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{parking.nom}</td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{parking.ville}</td>
                                            <td className="px-6 py-4 whitespace-nowrap">
                                                <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusStyle(parking.statut)}`}>
                                                    {parking.statut}
                                                </span>
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 font-semibold">{parking.tauxOccupation}%</td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{parking.adminLocal}</td>
                                            <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                                <button onClick={() => handleAction(parking, 'Détails')} className="text-blue-600 hover:text-blue-900 transition-colors">Détails</button>
                                                <button onClick={() => handleAction(parking, 'Désactiver')} className="text-red-600 hover:text-red-900 ml-4 transition-colors">Désactiver</button>
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                </CardContent>
                <CardFooter>
                    <div className="w-full flex justify-between items-center">
                        <span className="text-sm text-gray-700">
                            Affichage de {paginatedParkings.length} parkings sur {filteredParkings.length}
                        </span>
                        <div className="flex space-x-2">
                            <button
                                onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                                disabled={currentPage === 1 || isLoading}
                                className="px-3 py-1 text-sm font-medium border rounded-lg text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 transition-colors"
                            >
                                Précédent
                            </button>
                            <span className="px-3 py-1 text-sm font-medium text-gray-700">
                                Page {currentPage} de {totalPages}
                            </span>
                            <button
                                onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
                                disabled={currentPage === totalPages || isLoading}
                                className="px-3 py-1 text-sm font-medium border rounded-lg text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 transition-colors"
                            >
                                Suivant
                            </button>
                        </div>
                    </div>
                </CardFooter>
            </Card>
        </>
    );
};
