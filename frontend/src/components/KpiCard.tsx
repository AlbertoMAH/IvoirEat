// components/KpiCard.tsx
"use client";

import { FC, ReactNode } from 'react';

interface KpiCardProps {
  title: string;
  value: string;
  subtitle?: string;
  icon: ReactNode;
  colorClass: string;
}

const KpiCard: FC<KpiCardProps> = ({ title, value, subtitle, icon, colorClass }) => {
  return (
    <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm flex flex-col justify-between">
      <div>
        <div className="flex justify-between items-start">
          <p className="text-sm font-medium text-gray-500">{title}</p>
          <div className={`p-1.5 rounded-lg ${colorClass}`}>
            {icon}
          </div>
        </div>
        <p className="mt-2 text-3xl font-bold text-gray-900">{value}</p>
      </div>
      {subtitle && <p className="mt-4 text-xs text-gray-500">{subtitle}</p>}
    </div>
  );
};

export default KpiCard;
