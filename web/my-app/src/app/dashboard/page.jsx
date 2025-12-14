"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

export default function DashboardPage() {
  const router = useRouter();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    totalArticles: 0,
    publishedArticles: 0,
    draftArticles: 0,
  });

  useEffect(() => {
    const token = localStorage.getItem("access_token");
    const userData = localStorage.getItem("user");

    if (!token || !userData) {
      router.push("/login");
      return;
    }

    const parsedUser = JSON.parse(userData);
    if (parsedUser.role !== "admin" && parsedUser.role !== "editor") {
      router.push("/");
      return;
    }

    setUser(parsedUser);
    setLoading(false);
    fetchStats();
  }, [router]);

  const fetchStats = async () => {
    try {
      const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
      const res = await fetch(`${API_URL}/articles`);
      const articles = await res.json();
      
      if (Array.isArray(articles)) {
        setStats({
          totalArticles: articles.length,
          publishedArticles: articles.filter(a => a.status === "published").length,
          draftArticles: articles.filter(a => a.status === "draft").length,
        });
      }
    } catch (error) {
      console.error("Failed to fetch stats:", error);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user");
    router.push("/login");
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-100">
        <p className="text-gray-600">Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
          <div className="flex items-center gap-4">
            <span className="text-sm text-gray-600">
              {user?.username} ({user?.role})
            </span>
            <Link href="/" className="text-blue-600 hover:underline text-sm">
              Lihat Website
            </Link>
            <button
              onClick={handleLogout}
              className="text-red-600 hover:underline text-sm"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 py-8">
        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-white p-6 rounded-lg shadow">
            <p className="text-gray-500 text-sm">Total Artikel</p>
            <p className="text-3xl font-bold text-gray-900">{stats.totalArticles}</p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow">
            <p className="text-gray-500 text-sm">Published</p>
            <p className="text-3xl font-bold text-green-600">{stats.publishedArticles}</p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow">
            <p className="text-gray-500 text-sm">Draft</p>
            <p className="text-3xl font-bold text-yellow-600">{stats.draftArticles}</p>
          </div>
        </div>

        {/* Quick Actions */}
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Menu</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Link
            href="/dashboard/articles/create"
            className="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow"
          >
            <div className="text-4xl mb-4">📝</div>
            <h3 className="text-xl font-semibold mb-2">Tulis Artikel</h3>
            <p className="text-gray-600 text-sm">Buat artikel berita baru</p>
          </Link>

          <Link
            href="/dashboard/articles"
            className="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow"
          >
            <div className="text-4xl mb-4">📰</div>
            <h3 className="text-xl font-semibold mb-2">Kelola Artikel</h3>
            <p className="text-gray-600 text-sm">Lihat, edit, hapus artikel</p>
          </Link>

          {user?.role === "admin" && (
            <Link
              href="/dashboard/categories"
              className="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow"
            >
              <div className="text-4xl mb-4">📂</div>
              <h3 className="text-xl font-semibold mb-2">Kelola Kategori</h3>
              <p className="text-gray-600 text-sm">Tambah dan edit kategori</p>
            </Link>
          )}
        </div>
      </main>
    </div>
  );
}