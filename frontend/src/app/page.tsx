import { redirect } from 'next/navigation';

export default function RootPage() {
  // Par défaut, on redirige vers la page de connexion
  redirect('/login');
}
