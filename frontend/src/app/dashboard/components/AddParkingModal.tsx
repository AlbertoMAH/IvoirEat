'use client';

import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Icon } from './ui/Icon';
import { ParkingStatus } from '@/lib/types';

// Modal pour l'ajout de parking
export const AddParkingModal = ({ isOpen, onClose, onSave }: { isOpen: boolean, onClose: () => void, onSave: (data: any) => Promise<void> }) => {
    const [formData, setFormData] = useState({
        nom: '',
        ville: '',
        capaciteTotale: 100,
        adminLocal: '',
        lat: 48.8566,
        lng: 2.3522,
        statut: 'Disponible' as ParkingStatus,
    });
    const [isSaving, setIsSaving] = useState(false);

    useEffect(() => {
        // Réinitialiser le formulaire à l'ouverture/fermeture
        if (!isOpen) {
            setFormData({
                nom: '',
                ville: '',
                capaciteTotale: 100,
                adminLocal: '',
                lat: 48.8566,
                lng: 2.3522,
                statut: 'Disponible',
            });
        }
    }, [isOpen]);


    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const { name, value, type } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: (type === 'number' || name === 'lat' || name === 'lng') ? (parseFloat(value) || 0) : value,
        }));
    };

    const handleSubmit = async () => {
        setIsSaving(true);
        try {
            await onSave({
                nom: formData.nom,
                ville: formData.ville,
                adminLocal: formData.adminLocal,
                statut: formData.statut,
                coordonnees: { lat: formData.lat, lng: formData.lng },
                capaciteTotale: formData.capaciteTotale
            });
            onClose();
        } catch (error) {
            console.error("Erreur lors de l'ajout du parking:", error);
            // Ici, on afficherait un message d'erreur à l'utilisateur
        } finally {
            setIsSaving(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center z-50 p-4 transition-opacity duration-300">
            <Card className="w-full max-w-lg p-6 relative">
                <div className="flex justify-between items-center mb-4 border-b pb-3">
                    <h3 className="text-2xl font-bold text-gray-800">Ajouter un nouveau Parking</h3>
                    <button onClick={onClose} className="p-2 rounded-full text-gray-500 hover:bg-gray-100 transition-colors">
                        <Icon name="X" className="w-6 h-6" />
                    </button>
                </div>

                <div className="space-y-4 max-h-[70vh] overflow-y-auto pr-2">
                    {/* Nom du Parking */}
                    <div>
                        <label htmlFor="nom" className="block text-sm font-medium text-gray-700">Nom du Parking</label>
                        <input type="text" id="nom" name="nom" value={formData.nom} onChange={handleChange} required className="mt-1 block w-full border border-gray-300 rounded-lg shadow-sm p-2.5 focus:ring-blue-500 focus:border-blue-500" />
                    </div>

                    {/* Ville */}
                    <div>
                        <label htmlFor="ville" className="block text-sm font-medium text-gray-700">Ville</label>
                        <input type="text" id="ville" name="ville" value={formData.ville} onChange={handleChange} required className="mt-1 block w-full border border-gray-300 rounded-lg shadow-sm p-2.5 focus:ring-blue-500 focus:border-blue-500" />
                    </div>

                    {/* Capacité Totale */}
                    <div>
                        <label htmlFor="capaciteTotale" className="block text-sm font-medium text-gray-700">Capacité Totale (Nombre de places)</label>
                        <input type="number" id="capaciteTotale" name="capaciteTotale" value={formData.capaciteTotale} onChange={handleChange} required min="10" className="mt-1 block w-full border border-gray-300 rounded-lg shadow-sm p-2.5 focus:ring-blue-500 focus:border-blue-500" />
                    </div>

                    {/* Admin Local */}
                    <div>
                        <label htmlFor="adminLocal" className="block text-sm font-medium text-gray-700">Administrateur Local</label>
                        <input type="text" id="adminLocal" name="adminLocal" value={formData.adminLocal} onChange={handleChange} required className="mt-1 block w-full border border-gray-300 rounded-lg shadow-sm p-2.5 focus:ring-blue-500 focus:border-blue-500" />
                    </div>

                    {/* Coordonnées */}
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label htmlFor="lat" className="block text-sm font-medium text-gray-700">Latitude</label>
                            <input type="number" step="0.0001" id="lat" name="lat" value={formData.lat} onChange={handleChange} required className="mt-1 block w-full border border-gray-300 rounded-lg shadow-sm p-2.5 focus:ring-blue-500 focus:border-blue-500" />
                        </div>
                        <div>
                            <label htmlFor="lng" className="block text-sm font-medium text-gray-700">Longitude</label>
                            <input type="number" step="0.0001" id="lng" name="lng" value={formData.lng} onChange={handleChange} required className="mt-1 block w-full border border-gray-300 rounded-lg shadow-sm p-2.5 focus:ring-blue-500 focus:border-blue-500" />
                        </div>
                    </div>

                    {/* Statut Initial */}
                    <div>
                        <label htmlFor="statut" className="block text-sm font-medium text-gray-700">Statut Initial</label>
                        <select id="statut" name="statut" value={formData.statut} onChange={handleChange} required className="mt-1 block w-full border border-gray-300 rounded-lg shadow-sm p-2.5 focus:ring-blue-500 focus:border-blue-500 bg-white">
                            <option value="Disponible">Disponible</option>
                            <option value="Hors Ligne">Hors Ligne</option>
                        </select>
                    </div>
                </div>

                <div className="mt-6 pt-4 border-t flex justify-end space-x-3">
                    <button
                        onClick={onClose}
                        disabled={isSaving}
                        className="px-4 py-2 text-sm font-medium rounded-xl border border-gray-300 text-gray-700 bg-white hover:bg-gray-50 transition-colors disabled:opacity-50"
                    >
                        Annuler
                    </button>
                    <button
                        onClick={handleSubmit}
                        disabled={isSaving || !formData.nom || !formData.ville || !formData.adminLocal || formData.capaciteTotale < 1}
                        className="px-4 py-2 text-sm font-medium rounded-xl text-white bg-blue-600 hover:bg-blue-700 transition-colors disabled:bg-blue-400"
                    >
                        {isSaving ? 'Ajout en cours...' : 'Ajouter le Parking'}
                    </button>
                </div>
            </Card>
        </div>
    );
};
