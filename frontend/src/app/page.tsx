"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";

export default function Home() {
  const [count, setCount] = useState(0);

  const increment = () => {
    setCount(count + 1);
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-4">Hello World</h1>
        <p className="text-2xl mb-8">
          Compteur : <span className="font-mono font-bold">{count}</span>
        </p>
        <Button onClick={increment}>
          Incrémenter
        </Button>
      </div>
    </main>
  );
}
