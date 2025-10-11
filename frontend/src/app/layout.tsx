import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Providers from "@/components/Providers";
import { AuthProvider } from "@/context/AuthContext";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Babipark - Super Admin Dashboard",
  description: "Dashboard de gestion pour le réseau de parkings Babipark.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="fr">
      <body className={inter.className}>
        <AuthProvider>
          <Providers>
            {children}
          </Providers>
        </AuthProvider>
      </body>
    </html>
  );
}
