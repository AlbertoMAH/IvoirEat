/* eslint-disable */
"use client";

import { useState, FormEvent } from "react";
import { useAuth } from "../../../contexts/AuthContext";

export default function UploadReceiptPage() {
  const [file, setFile] = useState<File | null>(null);
  const [message, setMessage] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const { token } = useAuth();
  const [dragActive, setDragActive] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!file) {
      setMessage("Please select a file to upload.");
      return;
    }
    if (!token) {
      setMessage("You must be logged in to upload a receipt.");
      return;
    }

    setIsLoading(true);
    setMessage("");

    const formData = new FormData();
    formData.append("receiptImage", file);

    try {
      const response = await fetch("/api/receipts/upload", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to upload receipt");
      }

      const result = await response.json();
      setMessage(`Receipt uploaded successfully! Merchant: ${result.receipt.merchant}, Amount: ${result.receipt.amount}`);
      setFile(null); // Clear the file input
    } catch (error) {
      if (error instanceof Error) {
        setMessage(`Error: ${error.message}`);
      } else {
        setMessage("An unknown error occurred");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      setFile(e.dataTransfer.files[0]);
    }
  };

  return (
    <div className="container mx-auto p-4 max-w-lg">
      <h1 className="text-3xl font-bold mb-6 text-center text-gray-800">Uploader un Reçu</h1>
      <form onSubmit={handleSubmit} className="space-y-6">
        <div
          onDragEnter={handleDrag}
          onDragLeave={handleDrag}
          onDragOver={handleDrag}
          onDrop={handleDrop}
          className={`relative border-2 border-dashed rounded-lg p-10 text-center transition-colors duration-200 ${
            dragActive ? "border-indigo-600 bg-indigo-50" : "border-gray-300"
          }`}
        >
          <input
            id="receiptImage"
            name="receiptImage"
            type="file"
            onChange={handleFileChange}
            className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
          />
          <div className="flex flex-col items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" className="w-12 h-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
            <p className="mt-2 text-sm text-gray-600">
              <span className="font-semibold">Cliquez pour choisir un fichier</span> ou glissez-déposez.
            </p>
            <p className="text-xs text-gray-500">PNG, JPG, GIF jusqu'à 10MB</p>
            {file && <p className="mt-4 text-sm font-medium text-green-600">{file.name}</p>}
          </div>
        </div>

        <button
          type="submit"
          disabled={isLoading || !file}
          className="w-full px-4 py-3 bg-indigo-600 text-white font-semibold rounded-md hover:bg-indigo-700 disabled:bg-gray-400 transition-all duration-200 ease-in-out"
        >
          {isLoading ? "Analyse en cours..." : "Analyser et Sauvegarder"}
        </button>
      </form>
      {message && (
        <p className={`mt-4 text-sm text-center ${message.startsWith("Error") ? "text-red-600" : "text-gray-800"}`}>
          {message}
        </p>
      )}
    </div>
  );
}
