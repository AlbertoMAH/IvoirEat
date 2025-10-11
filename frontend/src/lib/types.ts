/** Type pour le statut d'un parking */
export type ParkingStatus = 'Disponible' | 'Presque Plein' | 'Complet' | 'Hors Ligne';

/** Type pour un objet Parking */
export interface Parking {
  id: string;
  nom: string;
  ville: string;
  statut: ParkingStatus;
  tauxOccupation: number; // 0 à 100
  adminLocal: string;
  coordonnees: { lat: number; lng: number; };
}

/** Type pour les indicateurs clés de performance (KPI) */
export interface KPIStats {
  totalParkingsActifs: number;
  tauxOccupationGlobal: number; // 0 à 100
  reservationsDuJour: number;
  revenuTotal: number;
  revenuDuMois: number;
}

/** Type pour une alerte système */
export interface Alert {
  id: string;
  type: 'Hors Ligne' | 'Saturation' | 'Maintenance';
  message: string;
  timestamp: number;
  parkingId: string;
}

/** Type pour les données de statistiques (graphiques) */
export interface StatsData {
  month: string;
  occupancy: number; // Taux d'occupation moyen
  revenue: number; // Revenu en milliers
}

/** Type for a parking administrator */
export interface AdminUser {
  id: string;
  name: string;
  email: string;
  role: 'Admin' | 'Super Admin';
  parkingIds: string[];
}
