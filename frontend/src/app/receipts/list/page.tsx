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
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Your Receipts</h1>
      {message && <p className="text-red-500">{message}</p>}
      {receipts.length === 0 && !message && (
        <p>You have no receipts yet. Try uploading one!</p>
      )}
      <div className="space-y-4">
        {receipts.map((receipt) => (
          <div key={receipt.ID} className="p-4 border rounded-lg shadow-sm">
            <div className="flex justify-between">
              <span className="font-semibold">{receipt.merchant}</span>
              <span className="font-bold text-lg">{receipt.amount.toFixed(2)} €</span>
            </div>
            <div className="text-sm text-gray-600">
              {new Date(receipt.date).toLocaleDateString()} - {receipt.receipt_type}
            </div>
            {receipt.is_anomaly && (
              <div className="mt-2 text-xs font-semibold text-red-600 bg-red-100 p-1 rounded-md inline-block">
                Anomaly Detected
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
