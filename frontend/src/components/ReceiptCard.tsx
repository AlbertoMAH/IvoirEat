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

const getCategoryStyle = (category: string) => {
  switch (category) {
    case "Transport":
      return {
        icon: <Car className="h-5 w-5 text-blue-600" />,
        colorClass: "border-blue-200",
        bgColorClass: "bg-blue-50",
      };
    case "Restaurant":
      return {
        icon: <Utensils className="h-5 w-5 text-orange-600" />,
        colorClass: "border-orange-200",
        bgColorClass: "bg-orange-50",
      };
    case "Hébergement":
      return {
        icon: <Home className="h-5 w-5 text-purple-600" />,
        colorClass: "border-purple-200",
        bgColorClass: "bg-purple-50",
      };
    case "Achats":
      return {
        icon: <ShoppingCart className="h-5 w-5 text-green-600" />,
        colorClass: "border-green-200",
        bgColorClass: "bg-green-50",
      };
    case "Services":
      return {
        icon: <Briefcase className="h-5 w-5 text-indigo-600" />,
        colorClass: "border-indigo-200",
        bgColorClass: "bg-indigo-50",
      };
    case "Santé":
      return {
        icon: <HeartPulse className="h-5 w-5 text-pink-600" />,
        colorClass: "border-pink-200",
        bgColorClass: "bg-pink-50",
      };
    default:
      return {
        icon: <ReceiptText className="h-5 w-5 text-gray-600" />,
        colorClass: "border-gray-200",
        bgColorClass: "bg-gray-50",
      };
  }
};


const ReceiptCard: FC<ReceiptCardProps> = ({ receipt }) => {
  const { icon, colorClass, bgColorClass } = getCategoryStyle(receipt.receipt_type);

  return (
    <div className={`bg-white rounded-xl border ${colorClass} shadow-sm overflow-hidden transform hover:-translate-y-1 transition-all duration-300`}>
      <div className="p-5">
        <div className="flex justify-between items-start">
          <div className={`p-3 rounded-lg ${bgColorClass}`}>
            {icon}
          </div>
          {receipt.is_anomaly && (
            <span className="text-xs font-semibold text-red-600 bg-red-100 px-2.5 py-1 rounded-full">
              Anomalie
            </span>
          )}
        </div>

        <div className="mt-4">
          <p className="text-2xl font-bold text-gray-900">
            {receipt.amount.toFixed(2)} <span className="text-lg font-normal text-gray-500">€</span>
          </p>
          <h3 className="mt-1 font-semibold text-md text-gray-800 truncate" title={receipt.merchant}>
            {receipt.merchant}
          </h3>
        </div>
      </div>
      <div className={`px-5 py-3 ${bgColorClass} bg-opacity-50 border-t ${colorClass}`}>
        <p className="text-xs text-gray-600">
          {new Date(receipt.date).toLocaleDateString('fr-FR', { day: '2-digit', month: 'long', year: 'numeric' })}
        </p>
      </div>
    </div>
  );
};

export default ReceiptCard;
