"use client";

import { useState, FormEvent } from "react";
import { useAuth } from "../../../contexts/AuthContext";

export default function UploadReceiptPage() {
  const [file, setFile] = useState<File | null>(null);
  const [message, setMessage] = useState<string>("");
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const { token } = useAuth();

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

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Upload a Receipt</h1>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="receiptImage" className="block text-sm font-medium text-gray-700">
            Receipt Image
          </label>
          <input
            id="receiptImage"
            name="receiptImage"
            type="file"
            onChange={handleFileChange}
            className="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-violet-50 file:text-violet-700 hover:file:bg-violet-100"
          />
        </div>
        <button
          type="submit"
          disabled={isLoading || !file}
          className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 disabled:bg-gray-400"
        >
          {isLoading ? "Uploading..." : "Upload"}
        </button>
      </form>
      {message && <p className="mt-4 text-sm text-gray-600">{message}</p>}
    </div>
  );
}
