"use client";

import { useAlerts } from "@/lib/api-client";
import { Alert } from "@/lib/types";

const NotificationsPanel = () => {
    const { data: alerts } = useAlerts();

    const getAlertStyle = (type: Alert['type']) => {
        switch (type) {
            case 'Hors Ligne': return 'border-red-500 bg-red-50 text-red-700';
            case 'Saturation': return 'border-orange-500 bg-orange-50 text-orange-700';
            case 'Maintenance': return 'border-blue-500 bg-blue-50 text-blue-700';
            default: return 'border-gray-500 bg-gray-50 text-gray-700';
        }
    };

    return (
        <div className="absolute right-0 mt-3 w-80 bg-white border border-gray-200 rounded-xl shadow-xl overflow-hidden z-20">
            <div className="p-4 border-b border-gray-100">
                <h3 className="text-lg font-semibold text-gray-800">Alertes Système ({alerts?.length || 0})</h3>
            </div>
            <div className="max-h-80 overflow-y-auto divide-y divide-gray-100">
                {alerts && alerts.length > 0 ? (
                    alerts.map((alert) => (
                        <div key={alert.id} className={`p-4 hover:bg-gray-50 border-l-4 ${getAlertStyle(alert.type)}`}>
                            <p className="font-medium text-sm">{alert.type}: {alert.message}</p>
                            <p className="text-xs text-gray-500 mt-1">{new Date(alert.timestamp).toLocaleString('fr-FR')}</p>
                        </div>
                    ))
                ) : (
                    <p className="p-4 text-sm text-gray-500 text-center">Aucune alerte récente.</p>
                )}
            </div>
        </div>
    );
};

export default NotificationsPanel;
