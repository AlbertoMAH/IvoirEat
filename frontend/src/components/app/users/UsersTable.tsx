"use client";

import { useUsers } from "@/lib/api-client";
import { AdminUser } from "@/lib/types";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";

export function UsersTable() {
  const { data: users, isLoading, isError } = useUsers();

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Nom</TableHead>
          <TableHead>Email</TableHead>
          <TableHead>Parkings Assignés</TableHead>
          <TableHead>Rôle</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {isLoading &&
          Array(3)
            .fill(0)
            .map((_, i) => (
              <TableRow key={i}>
                <TableCell colSpan={4}>
                  <Skeleton className="h-6 w-full" />
                </TableCell>
              </TableRow>
            ))}
        {isError && (
          <TableRow>
            <TableCell colSpan={4} className="text-center text-red-500">
              Erreur lors du chargement des utilisateurs.
            </TableCell>
          </TableRow>
        )}
        {!isLoading && !isError && users?.map((user: AdminUser) => (
          <TableRow key={user.id}>
            <TableCell className="font-medium">{user.name}</TableCell>
            <TableCell>{user.email}</TableCell>
            <TableCell>{user.parkingIds.join(', ')}</TableCell>
            <TableCell>{user.role}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
