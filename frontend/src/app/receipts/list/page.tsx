"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { Plus, BarChart, Clock, AlertTriangle } from "lucide-react";
import { useAuth } from "../../../contexts/AuthContext";
import KpiCard from "../../../components/KpiCard";
import ReceiptCard from "../../../components/ReceiptCard";

// Define the shape of a receipt object
interface Receipt {
  ID: number;
  amount: number;
  date: string;
  merchant: string;
  receipt_type: string;
  is_anomaly: boolean;
}

export default function ReceiptListPage() {
  const [receipts, setReceipts] = useState<Receipt[]>([]);
  const [message, setMessage] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const { token } = useAuth();

  useEffect(() => {
    const fetchReceipts = async () => {
      if (!token) {
        setIsLoading(false);
        setMessage("Please log in to view receipts.");
        return;
      }

      try {
        const response = await fetch("/api/receipts", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || "Failed to fetch receipts");
        }

        const data = await response.json();
        setReceipts(data.receipts || []);
      } catch (error) {
        if (error instanceof Error) {
          setMessage(`Error: ${error.message}`);
        } else {
          setMessage("An unknown error occurred");
        }
      } finally {
        setIsLoading(false);
      }
    };

    fetchReceipts();
  }, [token]);

  if (isLoading) {
    return <div className="container mx-auto p-4">Loading receipts...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto p-4 sm:p-6 lg:p-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">
            Mes Notes de Frais
          </h1>
          <Link
            href="/receipts/upload"
            className="inline-flex items-center justify-center px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            <Plus className="-ml-1 mr-2 h-5 w-5" />
            Ajouter
          </Link>
        </div>

        {/* KPI Cards Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
          <KpiCard
            title="Dépenses (30 jours)"
            value="1,250.50 €"
            subtitle="+5.2% vs. mois dernier"
            icon={<BarChart className="h-5 w-5 text-green-600" />}
            colorClass="bg-green-100"
          />
          <KpiCard
            title="En Attente"
            value="3"
            subtitle="Pour un total de 180.00 €"
            icon={<Clock className="h-5 w-5 text-yellow-600" />}
            colorClass="bg-yellow-100"
          />
          <KpiCard
            title="Anomalies"
            value="1"
            subtitle="Nécessite une vérification"
            icon={<AlertTriangle className="h-5 w-5 text-red-600" />}
            colorClass="bg-red-100"
          />
        </div>

        {message && <p className="text-red-500 text-center mb-4">{message}</p>}

        {isLoading && <p className="text-center text-gray-500">Chargement...</p>}

      {!isLoading && receipts.length === 0 && !message && (
        <div className="text-center py-12 px-6 bg-gray-50 rounded-lg">
          <svg xmlns="http://www.w3.org/2000/svg" className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900">Aucune note de frais</h3>
          <p className="mt-1 text-sm text-gray-500">Commencez par uploader votre premier reçu.</p>
        </div>
      )}

      {receipts.length > 0 && (
        <>
          <h2 className="text-xl font-bold text-gray-800 mb-4">
            Récents
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {receipts.map((receipt) => (
              <ReceiptCard key={receipt.ID} receipt={receipt} />
            ))}
          </div>
        </>
      )}
      </div>
    </div>
  );
}
