"use client";

import { useChartData } from "@/lib/api-client";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const StatsDashboard = () => {
    const { data: chartData, isLoading, isError } = useChartData();

    if (isError) return <Card className="p-4 text-red-500">Erreur de chargement des graphiques.</Card>;

    return (
        <Card className="p-6 rounded-2xl shadow-md">
            <h3 className="text-lg font-semibold mb-4 text-gray-800">Évolution du Taux d'Occupation (%)</h3>
            {isLoading ? <Skeleton className="h-64 w-full" /> : (
                <ResponsiveContainer width="100%" height={250}>
                    <LineChart data={chartData} margin={{ top: 5, right: 20, left: -10, bottom: 5 }}>
                        <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
                        <XAxis dataKey="month" stroke="#6b7280" />
                        <YAxis stroke="#6b7280" domain={[0, 100]} />
                        <Tooltip contentStyle={{ borderRadius: '10px', fontSize: '14px' }} />
                        <Legend />
                        <Line type="monotone" dataKey="occupancy" name="Taux d'Occupation" stroke="#3b82f6" strokeWidth={2} />
                        <Line type="monotone" dataKey="revenue" name="Revenus (k€)" stroke="#8b5cf6" strokeWidth={2} />
                    </LineChart>
                </ResponsiveContainer>
            )}
        </Card>
    );
};

export default StatsDashboard;
