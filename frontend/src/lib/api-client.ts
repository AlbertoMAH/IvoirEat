import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Parking, KPIStats, Alert, StatsData, AdminUser } from './types';

const API_DELAY = 600; // Délai de chargement simulé

// Fonction utilitaire pour simuler l'attente
const wait = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

// Données fictives
let currentMockParkings: Parking[] = [
  { id: 'p1', nom: 'Gare Centrale', ville: 'Paris', statut: 'Disponible', tauxOccupation: 25, adminLocal: 'Marc Dupont', coordonnees: { lat: 48.8566, lng: 2.3522 } },
  { id: 'p2', nom: 'Aeroport Ouest', ville: 'Lyon', statut: 'Complet', tauxOccupation: 100, adminLocal: 'Sophie Martin', coordonnees: { lat: 45.75, lng: 4.85 } },
  { id: 'p3', nom: 'Centre Commercial', ville: 'Marseille', statut: 'Presque Plein', tauxOccupation: 85, adminLocal: 'Luc Bernard', coordonnees: { lat: 43.296, lng: 5.369 } },
  { id: 'p4', nom: 'Hôtel de Ville', ville: 'Bordeaux', statut: 'Hors Ligne', tauxOccupation: 0, adminLocal: 'Elsa Petit', coordonnees: { lat: 44.837, lng: -0.579 } },
  { id: 'p5', nom: 'Pôle Événementiel', ville: 'Paris', statut: 'Disponible', tauxOccupation: 10, adminLocal: 'Marc Dupont', coordonnees: { lat: 48.86, lng: 2.33 } },
  { id: 'p6', nom: 'Zone Industrielle', ville: 'Lyon', statut: 'Presque Plein', tauxOccupation: 70, adminLocal: 'Sophie Martin', coordonnees: { lat: 45.77, lng: 4.83 } },
];

const mockStats: KPIStats = {
  totalParkingsActifs: 5,
  tauxOccupationGlobal: 50.8,
  reservationsDuJour: 45,
  revenuTotal: 345210,
  revenuDuMois: 15600,
};

const mockAlerts: Alert[] = [
  { id: 'a1', type: 'Hors Ligne', message: 'Parking "Hôtel de Ville" hors ligne depuis 15min.', timestamp: Date.now() - 900000, parkingId: 'p4' },
  { id: 'a2', type: 'Saturation', message: 'Parking "Centre Commercial" atteint 95% de sa capacité.', timestamp: Date.now() - 3600000, parkingId: 'p3' },
];

const mockChartData: StatsData[] = [
  { month: 'Jan', occupancy: 45, revenue: 15 },
  { month: 'Fev', occupancy: 55, revenue: 18 },
  { month: 'Mar', occupancy: 62, revenue: 22 },
  { month: 'Avr', occupancy: 58, revenue: 20 },
  { month: 'Mai', occupancy: 70, revenue: 25 },
  { month: 'Juin', occupancy: 65, revenue: 23 },
];

let mockUsers: AdminUser[] = [
    { id: 'u1', name: 'Marc Dupont', email: 'marc.dupont@example.com', role: 'Admin', parkingIds: ['p1', 'p5'] },
    { id: 'u2', name: 'Sophie Martin', email: 'sophie.martin@example.com', role: 'Admin', parkingIds: ['p2', 'p6'] },
];

/** Simule l'appel à une API REST */
async function mockFetch<T>(endpoint: string, data: T): Promise<T> {
  await wait(API_DELAY);

  // Dans une vraie app, on ferait un fetch ici et on vérifierait le token JWT
  // ex: const response = await fetch(`/api${endpoint}`);

  return data;
}

// --- Fonctions de mutation simulées ---
export const mockAddUser = async (newUserData: Omit<AdminUser, 'id'>): Promise<AdminUser> => {
    await wait(API_DELAY);
    const newUser: AdminUser = {
        ...newUserData,
        id: `u${mockUsers.length + 1}`,
    };
    mockUsers = [...mockUsers, newUser];
    return newUser;
};

// Hooks React Query pour le fetching de données
export const useParkings = () => useQuery<Parking[], Error>({
  queryKey: ['parkings'],
  queryFn: () => mockFetch('/api/parkings', currentMockParkings),
});

export const useKPIStats = () => useQuery<KPIStats, Error>({
  queryKey: ['kpiStats'],
  queryFn: () => mockFetch('/api/stats', mockStats),
});

export const useAlerts = () => useQuery<Alert[], Error>({
  queryKey: ['alerts'],
  queryFn: () => mockFetch('/api/alerts', mockAlerts),
});

export const useChartData = () => useQuery<StatsData[], Error>({
  queryKey: ['chartData'],
  queryFn: () => mockFetch('/api/stats', mockChartData),
});

export const useUsers = () => useQuery<AdminUser[], Error>({
    queryKey: ['users'],
    queryFn: () => mockFetch('/api/users', mockUsers),
});

export const useAddUser = () => {
    const queryClient = useQueryClient();
    return useMutation<AdminUser, Error, Omit<AdminUser, 'id'>>({
        mutationFn: mockAddUser,
        onSuccess: () => {
            // Invalider et re-fetcher la query 'users' après un succès
            queryClient.invalidateQueries({ queryKey: ['users'] });
        },
    });
};
