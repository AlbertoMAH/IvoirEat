'use client';

import React, { useMemo } from 'react';
import { useKPIStats } from '@/lib/api-client';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Icon } from './ui/Icon';
import { Skeleton } from '@/components/ui/skeleton';

// KPIStats
export const KPIStats = () => {
  const { data, isLoading, isError } = useKPIStats();

  const kpis = useMemo(() => ([
    { title: 'Parkings Actifs', value: data?.totalParkingsActifs, icon: 'ParkingCircle', format: (v: number) => v.toLocaleString('fr-FR'), color: 'text-blue-600', loading: isLoading },
    { title: 'Occupation Globale', value: data?.tauxOccupationGlobal, icon: 'Gauge', format: (v: number) => `${v.toFixed(1)}%`, color: 'text-orange-600', loading: isLoading },
    { title: 'Réservations J-1', value: data?.reservationsDuJour, icon: 'CalendarCheck', format: (v: number) => v.toLocaleString('fr-FR'), color: 'text-green-600', loading: isLoading },
  ]), [data, isLoading]);

  if (isError) return <div className="col-span-3 p-4 text-red-500">Erreur de chargement des KPIs.</div>;

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
      {kpis.map((kpi, index) => (
        <Card key={index}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-gray-500">{kpi.title}</CardTitle>
            <div className={`p-2 rounded-full bg-opacity-10 ${kpi.color.replace('text-', 'bg-')}`}>
              <Icon name={kpi.icon} className={`w-5 h-5 ${kpi.color}`} />
            </div>
          </CardHeader>
          <CardContent>
            {kpi.loading ? (
              <Skeleton className="h-8 w-3/4" />
            ) : (
              <div className="text-3xl font-bold text-gray-900">
                {kpi.value !== undefined ? kpi.format(kpi.value) : 'N/A'}
              </div>
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  );
};
