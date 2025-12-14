"use client";

import { useState, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Image from "next/image";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
const API_BASE = API_URL.replace("/api/v1", "");

export default function CreateArticlePage() {
  const router = useRouter();
  const fileInputRef = useRef(null);
  
  const [loading, setLoading] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [categoriesLoading, setCategoriesLoading] = useState(true);
  const [categories, setCategories] = useState([]);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  
  // Image upload state
  const [imagePreview, setImagePreview] = useState(null);

  const [formData, setFormData] = useState({
    judul: "",
    ringkasan: "",
    konten: "",
    status: "draft",
    kategori_ids: [],
    gambar_utama: "",
  });

  useEffect(() => {
    const token = localStorage.getItem("access_token");
    if (!token) {
      router.push("/login");
      return;
    }
    fetchCategories();
  }, [router]);

  const fetchCategories = async () => {
    setCategoriesLoading(true);
    try {
      const res = await fetch(`${API_URL}/categories`);
      const data = await res.json();
      
      // Handle different response structures
      let categoriesData = [];
      if (Array.isArray(data)) {
        categoriesData = data;
      } else if (data.data && Array.isArray(data.data)) {
        categoriesData = data.data;
      } else if (data.categories && Array.isArray(data.categories)) {
        categoriesData = data.categories;
      }
      
      setCategories(categoriesData);
      console.log("Categories loaded:", categoriesData);
    } catch (err) {
      console.error("Failed to fetch categories:", err);
    } finally {
      setCategoriesLoading(false);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleCategoryChange = (categoryId) => {
    setFormData((prev) => {
      const ids = prev.kategori_ids.includes(categoryId)
        ? prev.kategori_ids.filter((id) => id !== categoryId)
        : [...prev.kategori_ids, categoryId];
      return { ...prev, kategori_ids: ids };
    });
  };

  // Handle file selection and upload
  const handleFileSelect = async (e) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Validate file type
    const allowedTypes = ["image/jpeg", "image/png", "image/gif", "image/webp"];
    if (!allowedTypes.includes(file.type)) {
      setError("Format file tidak didukung. Gunakan JPG, PNG, GIF, atau WebP.");
      return;
    }

    // Validate file size (max 10MB)
    if (file.size > 10 * 1024 * 1024) {
      setError("Ukuran file terlalu besar. Maksimal 10MB.");
      return;
    }

    setError("");

    // Create preview
    const reader = new FileReader();
    reader.onload = (e) => setImagePreview(e.target?.result);
    reader.readAsDataURL(file);

    // Upload file
    setUploading(true);
    try {
      const token = localStorage.getItem("access_token");
      const uploadData = new FormData();
      uploadData.append("file", file);

      const res = await fetch(`${API_URL}/editor/upload`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: uploadData,
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.message || data.error || "Gagal upload gambar");
      }

      const imagePath = data.path || data.url || data.file_path;
      setFormData((prev) => ({ ...prev, gambar_utama: imagePath }));
    } catch (err) {
      setError(err.message);
      setImagePreview(null);
    } finally {
      setUploading(false);
    }
  };

  // Remove image
  const handleRemoveImage = () => {
    setImagePreview(null);
    setFormData((prev) => ({ ...prev, gambar_utama: "" }));
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setSuccess("");

    try {
      const token = localStorage.getItem("access_token");

      const response = await fetch(`${API_URL}/editor/articles`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(formData),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || data.error || "Gagal membuat artikel");
      }

      setSuccess("Artikel berhasil dibuat!");
      setFormData({
        judul: "",
        ringkasan: "",
        konten: "",
        status: "draft",
        kategori_ids: [],
        gambar_utama: "",
      });
      setImagePreview(null);

      setTimeout(() => router.push("/dashboard/articles"), 1500);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-4xl mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-gray-900">Tulis Artikel Baru</h1>
          <Link href="/dashboard" className="text-blue-600 hover:underline text-sm">
            ← Kembali
          </Link>
        </div>
      </header>

      {/* Form */}
      <main className="max-w-4xl mx-auto px-4 py-8">
        <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-6">
          {error && (
            <div className="mb-4 p-3 bg-red-100 text-red-700 rounded text-sm">
              {error}
            </div>
          )}
          {success && (
            <div className="mb-4 p-3 bg-green-100 text-green-700 rounded text-sm">
              {success}
            </div>
          )}

          {/* Gambar Thumbnail */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Gambar Thumbnail
            </label>
            
            {imagePreview ? (
              <div className="relative">
                <div className="relative w-full h-48 rounded-md overflow-hidden border border-gray-300">
                  <Image
                    src={imagePreview}
                    alt="Preview"
                    fill
                    className="object-cover"
                  />
                  {uploading && (
                    <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                      <span className="text-white">Mengupload...</span>
                    </div>
                  )}
                </div>
                {!uploading && formData.gambar_utama && (
                  <button
                    type="button"
                    onClick={handleRemoveImage}
                    className="mt-2 text-red-600 hover:text-red-800 text-sm"
                  >
                    Hapus gambar
                  </button>
                )}
              </div>
            ) : (
              <div
                onClick={() => fileInputRef.current?.click()}
                className="border-2 border-dashed border-gray-300 rounded-md p-6 text-center cursor-pointer hover:border-blue-400"
              >
                <p className="text-gray-500">Klik untuk upload gambar</p>
                <p className="text-xs text-gray-400 mt-1">JPG, PNG, GIF, WebP (Maks. 5MB)</p>
              </div>
            )}
            
            <input
              ref={fileInputRef}
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp"
              onChange={handleFileSelect}
              className="hidden"
            />
          </div>

          {/* Judul */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Judul Artikel <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              name="judul"
              value={formData.judul}
              onChange={handleChange}
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Masukkan judul artikel"
            />
          </div>

          {/* Ringkasan */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Ringkasan
            </label>
            <textarea
              name="ringkasan"
              value={formData.ringkasan}
              onChange={handleChange}
              rows={2}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Ringkasan singkat artikel (opsional)"
            />
          </div>

          {/* Konten */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Konten <span className="text-red-500">*</span>
            </label>
            <textarea
              name="konten"
              value={formData.konten}
              onChange={handleChange}
              required
              rows={12}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
              placeholder="Tulis konten artikel di sini..."
            />
            <p className="text-xs text-gray-500 mt-1">
              Gunakan HTML untuk formatting: &lt;p&gt;, &lt;h2&gt;, &lt;strong&gt;, &lt;ul&gt;, dll.
            </p>
          </div>

          {/* Kategori */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Kategori
            </label>
            
            {categoriesLoading ? (
              <div className="text-gray-500 text-sm">Memuat kategori...</div>
            ) : categories.length > 0 ? (
              <div className="flex flex-wrap gap-2">
                {categories.map((cat) => {
                  const catId = cat.kategori_id || cat.id;
                  const catName = cat.nama_kategori || cat.name;
                  const isSelected = formData.kategori_ids.includes(catId);
                  
                  return (
                    <button
                      key={catId}
                      type="button"
                      onClick={() => handleCategoryChange(catId)}
                      className={`px-4 py-2 rounded-full border text-sm font-medium transition-all ${
                        isSelected
                          ? "bg-blue-600 text-white border-blue-600 shadow-md"
                          : "bg-white text-gray-700 border-gray-300 hover:border-blue-400 hover:bg-blue-50"
                      }`}
                    >
                      {isSelected && <span className="mr-1">✓</span>}
                      {catName}
                    </button>
                  );
                })}
              </div>
            ) : (
              <div className="p-4 bg-yellow-50 border border-yellow-200 rounded-md">
                <p className="text-yellow-800 text-sm">
                  Belum ada kategori. 
                  <Link href="/dashboard/categories" className="text-blue-600 hover:underline ml-1">
                    Tambah kategori
                  </Link>
                </p>
              </div>
            )}
            
            {formData.kategori_ids.length > 0 && (
              <p className="text-xs text-gray-500 mt-2">
                {formData.kategori_ids.length} kategori dipilih
              </p>
            )}
          </div>

          {/* Status */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Status
            </label>
            <div className="flex gap-4">
              <label className={`flex items-center gap-2 px-4 py-2 rounded-md border cursor-pointer transition-all ${
                formData.status === "draft" 
                  ? "bg-yellow-50 border-yellow-400 text-yellow-800" 
                  : "bg-white border-gray-300 text-gray-700 hover:border-gray-400"
              }`}>
                <input
                  type="radio"
                  name="status"
                  value="draft"
                  checked={formData.status === "draft"}
                  onChange={handleChange}
                  className="hidden"
                />
                <span className="text-lg">📝</span>
                <span className="font-medium">Draft</span>
              </label>
              
              <label className={`flex items-center gap-2 px-4 py-2 rounded-md border cursor-pointer transition-all ${
                formData.status === "published" 
                  ? "bg-green-50 border-green-400 text-green-800" 
                  : "bg-white border-gray-300 text-gray-700 hover:border-gray-400"
              }`}>
                <input
                  type="radio"
                  name="status"
                  value="published"
                  checked={formData.status === "published"}
                  onChange={handleChange}
                  className="hidden"
                />
                <span className="text-lg">🌐</span>
                <span className="font-medium">Published</span>
              </label>
            </div>
          </div>

          {/* Buttons */}
          <div className="flex gap-3 pt-4 border-t">
            <button
              type="submit"
              disabled={loading || uploading}
              className="px-6 py-2.5 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed font-medium"
            >
              {loading ? "Menyimpan..." : "Simpan Artikel"}
            </button>
            <Link
              href="/dashboard"
              className="px-6 py-2.5 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 font-medium"
            >
              Batal
            </Link>
          </div>
        </form>
      </main>
    </div>
  );
}