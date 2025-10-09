"use client"

import { useState } from "react"
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table"
import { MoreHorizontal } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogClose,
} from "@/components/ui/dialog"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

// --- Type Definition ---
type Parking = {
  ID: number
  name: string
  location: string
  capacity: number
  tenant_id: number
}

// --- API Functions ---
const fetchParkings = async (): Promise<Parking[]> => {
  const token = localStorage.getItem("token")
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/parkings`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) throw new Error("Failed to fetch parkings")
  return response.json()
}

const createParking = async (newParking: Omit<Parking, "ID">) => {
  const token = localStorage.getItem("token")
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/parkings`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(newParking),
  })
  if (!response.ok) throw new Error("Failed to create parking")
  return response.json()
}

const updateParking = async (updatedParking: Parking) => {
  const token = localStorage.getItem("token")
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/parkings/${updatedParking.ID}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(updatedParking),
  })
  if (!response.ok) throw new Error("Failed to update parking")
  return response.json()
}

const deleteParking = async (id: number) => {
  const token = localStorage.getItem("token")
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/v1/parkings/${id}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) throw new Error("Failed to delete parking")
  return response.json()
}


// --- DataTable Component ---
interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
}

function DataTable<TData, TValue>({ columns, data }: DataTableProps<TData, TValue>) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  })

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(header.column.columnDef.header, header.getContext())}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id} data-state={row.getIsSelected() && "selected"}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}

// --- Main Page Component ---
export default function ParkingsPage() {
  const queryClient = useQueryClient()
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [selectedParking, setSelectedParking] = useState<Parking | null>(null)

  // Form state
  const [name, setName] = useState("")
  const [location, setLocation] = useState("")
  const [capacity, setCapacity] = useState(0)
  const [tenantId, setTenantId] = useState(1) // Defaulting to 1 for superadmin

  const { data: parkings, isLoading, isError } = useQuery<Parking[]>({
    queryKey: ["parkings"],
    queryFn: fetchParkings,
  })

  const mutationOptions = {
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["parkings"] })
      setIsDialogOpen(false)
      setSelectedParking(null)
    },
  }

  const createMutation = useMutation({ mutationFn: createParking, ...mutationOptions })
  const updateMutation = useMutation({ mutationFn: updateParking, ...mutationOptions })
  const deleteMutation = useMutation({ mutationFn: deleteParking, ...mutationOptions })

  const handleOpenDialog = (parking: Parking | null = null) => {
    setSelectedParking(parking)
    if (parking) {
      setName(parking.name)
      setLocation(parking.location)
      setCapacity(parking.capacity)
      setTenantId(parking.tenant_id)
    } else {
      setName("")
      setLocation("")
      setCapacity(0)
      setTenantId(1)
    }
    setIsDialogOpen(true)
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const parkingData = { name, location, capacity, tenant_id: tenantId }
    if (selectedParking) {
      updateMutation.mutate({ ...parkingData, ID: selectedParking.ID })
    } else {
      createMutation.mutate(parkingData)
    }
  }

  const columns: ColumnDef<Parking>[] = [
    { accessorKey: "name", header: "Nom" },
    { accessorKey: "location", header: "Adresse" },
    { accessorKey: "capacity", header: "Capacité" },
    { accessorKey: "tenant_id", header: "Client ID" },
    {
      id: "actions",
      cell: ({ row }) => {
        const parking = row.original
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuItem onClick={() => handleOpenDialog(parking)}>
                Modifier
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                className="text-red-600"
                onClick={() => deleteMutation.mutate(parking.ID)}
              >
                Supprimer
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]

  if (isLoading) return <div>Chargement des parkings...</div>
  if (isError) return <div>Erreur lors du chargement des parkings.</div>

  return (
    <div className="space-y-8">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Gestion des Parkings</CardTitle>
            <CardDescription>
              Ajoutez, modifiez ou supprimez des parkings.
            </CardDescription>
          </div>
          <Button onClick={() => handleOpenDialog()}>Ajouter un Parking</Button>
        </CardHeader>
        <CardContent>
          <DataTable columns={columns} data={parkings || []} />
        </CardContent>
      </Card>

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <form onSubmit={handleSubmit}>
            <DialogHeader>
              <DialogTitle>{selectedParking ? "Modifier le Parking" : "Ajouter un Parking"}</DialogTitle>
              <DialogDescription>
                Remplissez les détails du parking ici.
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="name" className="text-right">Nom</Label>
                <Input id="name" value={name} onChange={(e) => setName(e.target.value)} className="col-span-3" required />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="location" className="text-right">Adresse</Label>
                <Input id="location" value={location} onChange={(e) => setLocation(e.target.value)} className="col-span-3" required />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="capacity" className="text-right">Capacité</Label>
                <Input id="capacity" type="number" value={capacity} onChange={(e) => setCapacity(parseInt(e.target.value, 10))} className="col-span-3" required />
              </div>
               <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="tenantId" className="text-right">Client ID</Label>
                <Input id="tenantId" type="number" value={tenantId} onChange={(e) => setTenantId(parseInt(e.target.value, 10))} className="col-span-3" required />
              </div>
            </div>
            <DialogFooter>
              <DialogClose asChild><Button type="button" variant="secondary">Annuler</Button></DialogClose>
              <Button type="submit" disabled={createMutation.isPending || updateMutation.isPending}>
                Enregistrer
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
