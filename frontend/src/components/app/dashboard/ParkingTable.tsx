"use client";

import { useMemo, useState } from "react";
import { useParkings } from "@/lib/api-client";
import { Parking, ParkingStatus } from "@/lib/types";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button"; // Assurez-vous d'avoir ce composant

const ParkingTable = () => {
    const { data: parkings, isLoading, isError } = useParkings();
    const [searchTerm, setSearchTerm] = useState('');
    const [currentPage, setCurrentPage] = useState(1);
    const itemsPerPage = 5;

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

    return (
        <Card className="p-0">
            <div className="p-6 border-b border-gray-100 flex justify-between items-center">
                <h2 className="text-xl font-semibold text-gray-800">Liste des Parkings ({parkings?.length || 0})</h2>
                <Input
                    type="text"
                    placeholder="Rechercher..."
                    className="w-full sm:w-48"
                    value={searchTerm}
                    onChange={(e) => {
                        setSearchTerm(e.target.value);
                        setCurrentPage(1);
                    }}
                />
            </div>

            <div className="overflow-x-auto">
                <Table>
                    <TableHeader>
                        <TableRow>
                            {['Nom', 'Ville', 'Statut', 'Occupation', 'Admin Local', 'Actions'].map((header) => (
                                <TableHead key={header}>{header}</TableHead>
                            ))}
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {isLoading ? (
                            Array(itemsPerPage).fill(0).map((_, i) => (
                                <TableRow key={i}><TableCell colSpan={6}><Skeleton className="h-6 w-full" /></TableCell></TableRow>
                            ))
                        ) : isError ? (
                            <TableRow><TableCell colSpan={6} className="text-center text-red-500">Erreur de chargement.</TableCell></TableRow>
                        ) : paginatedParkings.length === 0 ? (
                            <TableRow><TableCell colSpan={6} className="text-center text-gray-500">Aucun résultat.</TableCell></TableRow>
                        ) : (
                            paginatedParkings.map((parking) => (
                                <TableRow key={parking.id}>
                                    <TableCell className="font-medium">{parking.nom}</TableCell>
                                    <TableCell>{parking.ville}</TableCell>
                                    <TableCell>
                                        <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusStyle(parking.statut)}`}>
                                            {parking.statut}
                                        </span>
                                    </TableCell>
                                    <TableCell className="font-semibold">{parking.tauxOccupation}%</TableCell>
                                    <TableCell>{parking.adminLocal}</TableCell>
                                    <TableCell>
                                        <Button variant="link" size="sm">Détails</Button>
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>

            <div className="p-4 flex justify-between items-center border-t border-gray-100">
                <span className="text-sm text-gray-700">
                    Page {currentPage} de {totalPages}
                </span>
                <div className="flex space-x-2">
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                        disabled={currentPage === 1 || isLoading}
                    >
                        Précédent
                    </Button>
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
                        disabled={currentPage === totalPages || isLoading}
                    >
                        Suivant
                    </Button>
                </div>
            </div>
        </Card>
    );
};

export default ParkingTable;
