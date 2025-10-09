"use client"

import { useState } from "react"
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

// Type corresponding to the backend model
interface Parking {
  ID: number;
  name: string;
  location: string;
  capacity: number;
}

// Type for creating/updating a parking
type ParkingInput = Omit<Parking, 'ID'>;

// --- API Functions ---
const fetchParkings = async (): Promise<Parking[]> => {
  const token = localStorage.getItem("token")
  const response = await fetch(`/api/v1/parkings`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) throw new Error("Failed to fetch parkings")
  const data = await response.json();
  return data || []; // Ensure it returns an array even if API returns null
}

const createParking = async (newParking: ParkingInput) => {
  const token = localStorage.getItem("token")
  const response = await fetch(`/api/v1/parkings`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    body: JSON.stringify(newParking),
  })
  if (!response.ok) throw new Error('Failed to create parking')
  return response.json()
}

const updateParking = async (parkingToUpdate: Parking) => {
    const token = localStorage.getItem("token")
    const { ID, ...data } = parkingToUpdate;
    const response = await fetch(`/api/v1/parkings/${ID}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify(data),
    })
    if (!response.ok) throw new Error('Failed to update parking')
    return response.json()
}

const deleteParking = async (id: number) => {
    const token = localStorage.getItem("token")
    const response = await fetch(`/api/v1/parkings/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
    })
    if (!response.ok) throw new Error('Failed to delete parking')
    return { success: true };
}


export default function ParkingsPage() {
  const queryClient = useQueryClient()
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [selectedParking, setSelectedParking] = useState<Parking | null>(null)

  const { data: parkings = [], isLoading, isError, error } = useQuery<Parking[], Error>({
    queryKey: ["parkings"],
    queryFn: fetchParkings,
  })

  const invalidateAndRefetch = () => {
    queryClient.invalidateQueries({ queryKey: ["parkings"] })
  }

  const createMutation = useMutation({ mutationFn: createParking, onSuccess: invalidateAndRefetch })
  const updateMutation = useMutation({ mutationFn: updateParking, onSuccess: invalidateAndRefetch })
  const deleteMutation = useMutation({ mutationFn: deleteParking, onSuccess: invalidateAndRefetch })

  const openForm = (parking: Parking | null) => {
    setSelectedParking(parking)
    setIsFormOpen(true)
  }

  const closeForm = () => {
    setIsFormOpen(false)
    setSelectedParking(null)
  }

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const data: ParkingInput = {
      name: formData.get('name') as string,
      location: formData.get('location') as string,
      capacity: parseInt(formData.get('capacity') as string, 10),
    };

    if (selectedParking) {
      updateMutation.mutate({ ...data, ID: selectedParking.ID });
    } else {
      createMutation.mutate(data);
    }
    closeForm();
  }

  if (isLoading) return <div>Chargement des parkings...</div>
  if (isError) return <div>Erreur: {error.message}</div>

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Gestion des Parkings</h1>
        <Button onClick={() => openForm(null)}>Ajouter un parking</Button>
      </div>

      <Dialog open={isFormOpen} onOpenChange={setIsFormOpen}>
        <DialogContent onInteractOutside={(e) => e.preventDefault()} onEscapeKeyDown={(e) => e.preventDefault()}>
          <form onSubmit={handleSubmit}>
            <DialogHeader>
              <DialogTitle>{selectedParking ? "Modifier le parking" : "Ajouter un nouveau parking"}</DialogTitle>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="name" className="text-right">Nom</Label>
                <Input id="name" name="name" defaultValue={selectedParking?.name} required className="col-span-3" />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="location" className="text-right">Adresse</Label>
                <Input id="location" name="location" defaultValue={selectedParking?.location} required className="col-span-3" />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="capacity" className="text-right">Places</Label>
                <Input id="capacity" name="capacity" type="number" defaultValue={selectedParking?.capacity} required className="col-span-3" />
              </div>
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={closeForm}>Annuler</Button>
              <Button type="submit">Enregistrer</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Nom</TableHead>
              <TableHead>Adresse</TableHead>
              <TableHead>Capacité</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {parkings.length > 0 ? (
              parkings.map((parking) => (
                <TableRow key={parking.ID}>
                  <TableCell>{parking.name}</TableCell>
                  <TableCell>{parking.location}</TableCell>
                  <TableCell>{parking.capacity}</TableCell>
                  <TableCell className="text-right space-x-2">
                    <Button variant="outline" size="sm" onClick={() => openForm(parking)}>Modifier</Button>
                    <AlertDialog>
                      <AlertDialogTrigger asChild>
                        <Button variant="destructive" size="sm">Supprimer</Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Êtes-vous sûr ?</AlertDialogTitle>
                          <AlertDialogDescription>
                            Cette action est irréversible et supprimera définitivement le parking &quot;{parking.name}&quot;.
                          </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Annuler</AlertDialogCancel>
                          <AlertDialogAction onClick={() => deleteMutation.mutate(parking.ID)}>Supprimer</AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={4} className="text-center">Aucun parking trouvé.</TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
