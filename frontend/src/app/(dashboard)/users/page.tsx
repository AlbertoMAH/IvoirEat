"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { PlusCircle } from "lucide-react";
import { CreateAdminForm } from "@/components/app/users/CreateAdminForm";

export default function UsersPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);

  return (
    <>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Gestion des Utilisateurs</CardTitle>
            <CardDescription>
              Ajoutez, modifiez ou supprimez les administrateurs de parking.
            </CardDescription>
          </div>
          <Button onClick={() => setIsModalOpen(true)}>
            <PlusCircle className="mr-2 h-4 w-4" />
            Créer un admin
          </Button>
        </CardHeader>
        <CardContent>
          <UsersTable />
        </CardContent>
      </Card>
      <CreateAdminForm open={isModalOpen} onOpenChange={setIsModalOpen} />
    </>
  );
}
