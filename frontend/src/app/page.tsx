import KPIStats from "@/components/app/dashboard/KPIStats";
import ParkingTable from "@/components/app/dashboard/ParkingTable";
import MapView from "@/components/app/dashboard/MapView";
import StatsDashboard from "@/components/app/dashboard/StatsDashboard";

export default function DashboardPage() {
  return (
    <div className="space-y-8">
      <KPIStats />
      <MapView />
      <StatsDashboard />
      <ParkingTable />
    </div>
  );
}
