"use client";

import { useState, useEffect } from "react";
import { useAuth } from "../../../contexts/AuthContext";

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
    <div className="container mx-auto p-4 max-w-4xl">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-800">Mes Notes de Frais</h1>
      </div>

      {message && <p className="text-red-500 text-center">{message}</p>}

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

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {receipts.map((receipt) => (
          <div key={receipt.ID} className="bg-white border rounded-lg shadow-sm overflow-hidden transform hover:-translate-y-1 transition-all duration-200">
            <div className="p-5">
              <div className="flex justify-between items-start">
                <h3 className="font-semibold text-lg text-gray-800 truncate" title={receipt.merchant}>
                  {receipt.merchant}
                </h3>
                {receipt.is_anomaly && (
                  <span className="ml-2 text-xs font-semibold text-red-600 bg-red-100 px-2 py-1 rounded-full">
                    Anomalie
                  </span>
                )}
              </div>
              <p className="text-sm text-gray-500 mb-4">{receipt.receipt_type}</p>

              <p className="text-3xl font-bold text-gray-900 mb-4">
                {receipt.amount.toFixed(2)} <span className="text-xl font-normal text-gray-500">€</span>
              </p>

              <p className="text-sm text-gray-600">
                {new Date(receipt.date).toLocaleDateString()}
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
