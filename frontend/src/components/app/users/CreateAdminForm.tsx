"use client";

import { useState } from "react";
import { useParkings, useAddUser } from "@/lib/api-client";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export function CreateAdminForm({ open, onOpenChange }: { open: boolean, onOpenChange: (open: boolean) => void }) {
  const { data: parkings } = useParkings();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [parkingId, setParkingId] = useState<string | undefined>(undefined);
  const addUserMutation = useAddUser();

  const handleSubmit = async () => {
    if (!name || !email || !password || !parkingId) {
      alert("Veuillez remplir tous les champs.");
      return;
    }

    await addUserMutation.mutateAsync({
      name,
      email,
      role: 'Admin',
      parkingIds: [parkingId],
    });

    onOpenChange(false); // Fermer la modale après soumission
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Créer un nouvel administrateur</DialogTitle>
          <DialogDescription>
            Remplissez les informations pour le nouvel admin de parking.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="name" className="text-right">Nom</Label>
            <Input id="name" value={name} onChange={(e) => setName(e.target.value)} className="col-span-3" />
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="email" className="text-right">Email</Label>
            <Input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} className="col-span-3" />
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="password"  className="text-right">Mot de passe</Label>
            <Input id="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} className="col-span-3" />
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="parking" className="text-right">Parking</Label>
            <Select onValueChange={setParkingId} value={parkingId}>
              <SelectTrigger className="col-span-3">
                <SelectValue placeholder="Sélectionner un parking" />
              </SelectTrigger>
              <SelectContent>
                {parkings?.map((p) => (
                  <SelectItem key={p.id} value={p.id}>{p.nom}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
        <DialogFooter>
          <Button type="submit" onClick={handleSubmit}>Créer</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
