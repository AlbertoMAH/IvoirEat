'use client';

import React from 'react';
import { useChartData } from '@/lib/api-client';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

// StatsDashboard (Recharts)
export const StatsDashboard = () => {
    const { data: chartData, isLoading, isError } = useChartData();

    if (isError) return <div className="p-4 text-red-500">Erreur de chargement des graphiques.</div>;

    // Graphique 1: Évolution du Taux d'Occupation
    const OccupancyChart = () => (
        <Card>
            <CardHeader>
                <CardTitle>Évolution du Taux d'Occupation (%)</CardTitle>
            </CardHeader>
            <CardContent>
                {isLoading ? <Skeleton className="h-[250px] w-full" /> : (
                    <ResponsiveContainer width="100%" height={250}>
                        <LineChart data={chartData} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
                            <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
                            <XAxis dataKey="month" stroke="#6b7280" />
                            <YAxis stroke="#6b7280" domain={[0, 100]} label={{ value: '%', angle: -90, position: 'insideLeft' }} />
                            <Tooltip contentStyle={{ borderRadius: '10px', fontSize: '14px' }} formatter={(value: number) => [`${value.toFixed(1)}%`, 'Taux moyen']} />
                            <Legend wrapperStyle={{ paddingTop: '10px' }} />
                            <Line type="monotone" dataKey="occupancy" name="Taux d'Occupation" stroke="#3b82f6" strokeWidth={2} dot={{ fill: '#3b82f6', r: 4 }} activeDot={{ r: 8 }} />
                        </LineChart>
                    </ResponsiveContainer>
                )}
            </CardContent>
        </Card>
    );

    return (
        <div className="space-y-6">
            <OccupancyChart />
        </div>
    );
};
