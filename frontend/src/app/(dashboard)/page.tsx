"use client"

import { useQuery } from "@tanstack/react-query"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts"
import { useEffect } from "react";
import { useRouter } from "next/navigation";


// Define the structure of the stats data
interface DashboardStats {
  total_parkings: number
  total_spots: number
  available_spots: number
  total_users: number
  today_reservations: number
  weekly_reservations: { date: string; count: number }[]
}

// Function to fetch dashboard statistics
const fetchDashboardStats = async (): Promise<DashboardStats> => {
  const token = localStorage.getItem("token")
  if (!token) {
    throw new Error("No authentication token found")
  }

  const response = await fetch(
    `/api/v1/dashboard/stats`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  )

  if (!response.ok) {
    if (response.status === 401) {
      // Handle unauthorized access, e.g., redirect to login
      localStorage.removeItem("token");
    }
    throw new Error("Failed to fetch dashboard stats")
  }

  return response.json()
}

// Format data for the chart
const formatChartData = (
  data: { date: string; count: number }[] | undefined
) => {
  if (!data) return []
  const days = ["Dim", "Lun", "Mar", "Mer", "Jeu", "Ven", "Sam"]
  return data.map((item) => ({
    name: days[new Date(item.date).getDay()],
    reservations: item.count,
  }))
}

export default function DashboardPage() {
    const router = useRouter()
  const { data, isLoading, isError, error } = useQuery<DashboardStats, Error>({
    queryKey: ["dashboardStats"],
    queryFn: fetchDashboardStats,
      retry: (failureCount, error) => {
      if (error.message.includes("401")) {
        return false
      }
      return failureCount < 3
    },
  })

   useEffect(() => {
    if (error?.message.includes("401") || error?.message.includes("No authentication token found")) {
      router.push("/login")
    }
  }, [error, router])


  if (isLoading) {
    return <div>Chargement des statistiques...</div>
  }

  if (isError) {
    return (
      <div>
        Erreur lors de la récupération des données: {error?.message}
      </div>
    )
  }

  const chartData = formatChartData(data?.weekly_reservations)

  return (
    <div className="space-y-8">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader>
            <CardTitle>Parkings</CardTitle>
            <CardDescription>Nombre total de parkings</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-4xl font-bold">{data?.total_parkings}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Places</CardTitle>
            <CardDescription>Total / Disponibles</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-4xl font-bold">
              {data?.total_spots} / {data?.available_spots}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Utilisateurs</CardTitle>
            <CardDescription>Nombre total d&apos;admins</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-4xl font-bold">{data?.total_users}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Réservations</CardTitle>
            <CardDescription>Aujourd&apos;hui</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-4xl font-bold">{data?.today_reservations}</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Réservations de la semaine</CardTitle>
          <CardDescription>
            Graphique des tendances des réservations sur les 7 derniers jours.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Bar dataKey="reservations" fill="#8884d8" />
            </BarChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Localisation des parkings</CardTitle>
          <CardDescription>
            Carte affichant les parkings actifs.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-64 w-full bg-gray-200 flex items-center justify-center">
            <p className="text-gray-500">
              [Placeholder pour la carte Leaflet / MapLibre]
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
