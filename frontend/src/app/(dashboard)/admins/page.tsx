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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

// --- TypeScript Interfaces ---
interface Parking {
  ID: number;
  name: string;
}

interface Admin {
  ID: number;
  name: string;
  email: string;
  tenant_id: number;
}

type AdminInput = Omit<Admin, 'ID'> & { password?: string };

// --- API Functions ---
const fetchAdmins = async (): Promise<Admin[]> => {
  const token = localStorage.getItem("token")
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/users?role=tenant_admin`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) throw new Error("Failed to fetch admins")
  const data = await response.json()
  return data || []
}

const fetchParkings = async (): Promise<Parking[]> => {
    const token = localStorage.getItem("token")
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/parkings`, {
        headers: { Authorization: `Bearer ${token}` },
    })
    if (!response.ok) throw new Error("Failed to fetch parkings for dropdown")
    const data = await response.json()
    return data || []
}

const createAdmin = async (newAdmin: AdminInput) => {
    const token = localStorage.getItem("token")
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify({ ...newAdmin, role: 'tenant_admin' }),
    })
    if (!response.ok) throw new Error('Failed to create admin')
    return response.json()
}

const updateAdmin = async (adminToUpdate: Omit<Admin, 'tenant_id'> & { tenant_id?: number, password?: string }) => {
    const token = localStorage.getItem("token")
    const { ID, ...data } = adminToUpdate
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/users/${ID}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify(data),
    })
    if (!response.ok) throw new Error('Failed to update admin')
    return response.json()
}

const deleteAdmin = async (id: number) => {
    const token = localStorage.getItem("token")
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/users/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
    })
    if (!response.ok) throw new Error('Failed to delete admin')
    return { success: true }
}

export default function AdminsPage() {
    const queryClient = useQueryClient()
    const [isFormOpen, setIsFormOpen] = useState(false)
    const [selectedAdmin, setSelectedAdmin] = useState<Admin | null>(null)

    const { data: admins = [], isLoading: isLoadingAdmins, isError: isErrorAdmins, error: errorAdmins } = useQuery<Admin[], Error>({
        queryKey: ["admins"],
        queryFn: fetchAdmins,
    });

    const { data: parkings = [], isLoading: isLoadingParkings } = useQuery<Parking[], Error>({
        queryKey: ["parkings"],
        queryFn: fetchParkings,
    });

    const invalidateAndRefetch = () => {
        queryClient.invalidateQueries({ queryKey: ["admins"] })
    }

    const createMutation = useMutation({ mutationFn: createAdmin, onSuccess: invalidateAndRefetch })
    const updateMutation = useMutation({ mutationFn: updateAdmin, onSuccess: invalidateAndRefetch })
    const deleteMutation = useMutation({ mutationFn: deleteAdmin, onSuccess: invalidateAndRefetch })

    const openForm = (admin: Admin | null) => {
        setSelectedAdmin(admin)
        setIsFormOpen(true)
    }

    const closeForm = () => {
        setIsFormOpen(false)
        setSelectedAdmin(null)
    }

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const formData = new FormData(event.currentTarget);
        const password = formData.get('password') as string;

        const data: AdminInput = {
          name: formData.get('name') as string,
          email: formData.get('email') as string,
          tenant_id: parseInt(formData.get('tenant_id') as string, 10),
        };

        if (password) {
            data.password = password;
        }

        if (selectedAdmin) {
            updateMutation.mutate({ ...data, ID: selectedAdmin.ID });
        } else {
            createMutation.mutate(data);
        }
        closeForm();
    }

    const getParkingName = (tenantId: number) => {
        return parkings.find(p => p.ID === tenantId)?.name || 'N/A';
    }

    if (isLoadingAdmins || isLoadingParkings) return <div>Chargement...</div>
    if (isErrorAdmins) return <div>Erreur (admins): {errorAdmins.message}</div>

    return (
        <div className="space-y-4">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold">Gestion des Administrateurs</h1>
                <Button onClick={() => openForm(null)}>Ajouter un administrateur</Button>
            </div>

            <Dialog open={isFormOpen} onOpenChange={setIsFormOpen}>
                <DialogContent>
                    <form onSubmit={handleSubmit}>
                        <DialogHeader>
                            <DialogTitle>{selectedAdmin ? "Modifier l'administrateur" : "Ajouter un administrateur"}</DialogTitle>
                        </DialogHeader>
                        <div className="grid gap-4 py-4">
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="name" className="text-right">Nom</Label>
                                <Input id="name" name="name" defaultValue={selectedAdmin?.name} required className="col-span-3" />
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="email" className="text-right">Email</Label>
                                <Input id="email" name="email" type="email" defaultValue={selectedAdmin?.email} required className="col-span-3" />
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="password" className="text-right">Mot de passe</Label>
                                <Input id="password" name="password" type="password" placeholder={selectedAdmin ? "Laisser vide pour ne pas changer" : "Mot de passe"} className="col-span-3" />
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="tenant_id" className="text-right">Parking</Label>
                                <Select name="tenant_id" defaultValue={selectedAdmin?.tenant_id.toString()} required>
                                    <SelectTrigger className="col-span-3">
                                        <SelectValue placeholder="Sélectionner un parking" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {parkings.map((parking) => (
                                            <SelectItem key={parking.ID} value={parking.ID.toString()}>
                                                {parking.name}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
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
                            <TableHead>Email</TableHead>
                            <TableHead>Parking</TableHead>
                            <TableHead className="text-right">Actions</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {admins.map((admin) => (
                            <TableRow key={admin.ID}>
                                <TableCell>{admin.name}</TableCell>
                                <TableCell>{admin.email}</TableCell>
                                <TableCell>{getParkingName(admin.tenant_id)}</TableCell>
                                <TableCell className="text-right space-x-2">
                                    <Button variant="outline" size="sm" onClick={() => openForm(admin)}>Modifier</Button>
                                    <AlertDialog>
                                        <AlertDialogTrigger asChild><Button variant="destructive" size="sm">Supprimer</Button></AlertDialogTrigger>
                                        <AlertDialogContent>
                                            <AlertDialogHeader>
                                                <AlertDialogTitle>Êtes-vous sûr ?</AlertDialogTitle>
                                                <AlertDialogDescription>
                                                    Cette action supprimera l&apos;administrateur &quot;{admin.name}&quot;.
                                                </AlertDialogDescription>
                                            </AlertDialogHeader>
                                            <AlertDialogFooter>
                                                <AlertDialogCancel>Annuler</AlertDialogCancel>
                                                <AlertDialogAction onClick={() => deleteMutation.mutate(admin.ID)}>Supprimer</AlertDialogAction>
                                            </AlertDialogFooter>
                                        </AlertDialogContent>
                                    </AlertDialog>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </div>
        </div>
    )
}
