// components/ReceiptCard.tsx
"use client";

import { FC } from "react";
import {
  Briefcase,
  Car,
  HeartPulse,
  Home,
  ReceiptText,
  ShoppingCart,
  Utensils,
} from "lucide-react";

// Define the shape of a receipt object
interface Receipt {
  ID: number;
  amount: number;
  date: string;
  merchant: string;
  receipt_type: string;
  is_anomaly: boolean;
}

interface ReceiptCardProps {
  receipt: Receipt;
}

const ReceiptCard: FC<ReceiptCardProps> = ({ receipt }) => {
  return (
    <div className="bg-white border rounded-lg shadow-sm overflow-hidden transform hover:-translate-y-1 transition-all duration-200">
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
          {receipt.amount.toFixed(2)}{" "}
          <span className="text-xl font-normal text-gray-500">€</span>
        </p>
        <p className="text-sm text-gray-600">
          {new Date(receipt.date).toLocaleDateString()}
        </p>
      </div>
    </div>
  );
};

export default ReceiptCard;
