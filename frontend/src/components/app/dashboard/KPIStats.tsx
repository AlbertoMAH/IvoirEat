"use client";

import { useKPIStats } from "@/lib/api-client";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { ParkingCircle, Gauge, CalendarCheck, DollarSign } from 'lucide-react';
import { useMemo } from "react";

const KPIStats = () => {
  const { data, isLoading, isError } = useKPIStats();

  const kpis = useMemo(() => ([
    { title: 'Parkings Actifs', value: data?.totalParkingsActifs, icon: ParkingCircle, format: (v: number) => v.toLocaleString('fr-FR'), color: 'text-blue-600', loading: isLoading },
    { title: 'Occupation Globale', value: data?.tauxOccupationGlobal, icon: Gauge, format: (v: number) => `${v.toFixed(1)}%`, color: 'text-orange-600', loading: isLoading },
    { title: 'Réservations J-1', value: data?.reservationsDuJour, icon: CalendarCheck, format: (v: number) => v.toLocaleString('fr-FR'), color: 'text-green-600', loading: isLoading },
    { title: 'Revenus du Mois', value: data?.revenuDuMois, icon: DollarSign, format: (v: number) => `${(v / 1000).toFixed(1)}k €`, color: 'text-purple-600', loading: isLoading },
  ]), [data, isLoading]);

  if (isError) return <div className="col-span-4 p-4 text-red-500">Erreur de chargement des KPIs.</div>;

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      {kpis.map((kpi, index) => (
        <Card key={index} className="p-6 rounded-2xl shadow-md">
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium text-gray-500">{kpi.title}</p>
            <div className={`p-2 rounded-full bg-opacity-10 ${kpi.color.replace('text-', 'bg-')}`}>
              <kpi.icon className={`w-5 h-5 ${kpi.color}`} />
            </div>
          </div>
          <div className="mt-2">
            {kpi.loading ? (
              <Skeleton className="h-8 w-3/4" />
            ) : (
              <p className="text-3xl font-bold text-gray-900">
                {kpi.value !== undefined ? kpi.format(kpi.value) : 'N/A'}
              </p>
            )}
          </div>
        </Card>
      ))}
    </div>
  );
};

export default KPIStats;
