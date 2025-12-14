'use client';

import Image from "next/image";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import CurrentDate from "./CurrentDate";

export default function Header() {
  const pathname = usePathname();
  const router = useRouter();
  const [user, setUser] = useState(null);

  useEffect(() => {
    // Ambil user dari localStorage saat komponen mount
    const storedUser = localStorage.getItem("user");
    if (storedUser) {
      try {
        setUser(JSON.parse(storedUser));
      } catch (error) {
        console.error("Error parsing user data:", error);
      }
    }
  }, []);

  const handleLogout = () => {
    // Hapus token dan user dari localStorage
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user");
    setUser(null);
    router.push("/");
  };

  const navItems = [
    { href: "/international", label: "INTERNATIONAL" },
    { href: "/national", label: "NATIONAL" },
    { href: "/politics", label: "POLITICS" },
    { href: "/business", label: "BUSINESS" },
    { href: "/scitech", label: "SCI-TECH" },
    { href: "/lifestyle", label: "LIFESTYLE" },
    { href: "/health", label: "HEALTH" },
    { href: "/sports", label: "SPORTS" },
    { href: "/opinion", label: "OPINION" },
  ];

  return (
    <header className="">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Top utility bar */}
        <div className="flex justify-between items-center py-2 text-sm">
          <div className="flex items-center gap-3">
            {/* Search Bar dengan Icon */}
            <div className="relative">
              <div className="absolute left-2 top-1/2 -translate-y-1/2 pointer-events-none">
                <Image 
                  src="/searchlogo.svg" // Ganti .svg ke .png jika file Anda png
                  alt="Search" 
                  width={12} 
                  height={12}
                  style={{ filter: 'invert(0.5)' }}
                />
              </div>
              <input
                type="text"
                placeholder="SEARCH"
                className="border border-gray-300 rounded pl-8 pr-3 py-1 text-xs w-40 uppercase tracking-wide"
              />
            </div>
          </div>
          
          {/* Username atau Login/Register */}
          <div className="flex items-center gap-2 text-gray-600 text-xs uppercase tracking-wide">
            {user ? (
              <>
                <span>{user.username}</span>
                <Image 
                  src="/userLogo.svg" // Ganti .svg ke .png jika file Anda png
                  alt="User Profile" 
                  width={16} 
                  height={16}
                  style={{ filter: 'invert(0.5)' }}
                />
                <button 
                  onClick={handleLogout}
                  className="ml-2 text-blue-600 hover:underline"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link href="/login" className="text-blue-600 hover:underline">Login</Link>
                <span>/</span>
                <Link href="/register" className="text-blue-600 hover:underline">Register</Link>
              </>
            )}
          </div>
        </div>

        {/* Date */}
        <div className="text-left py-2 text-xs">
          <CurrentDate />
        </div>

        {/* Logo + tagline */}
        <div className="text-center py-6">
          <Link href="/">
            <h1 className="text-4xl md:text-5xl font-bold tracking-widest mb-2">
              BINTARO TIMES
            </h1>
          </Link>
          <p className="text-gray-600 text-xs uppercase tracking-widest">
            Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit, Duis
            Sollicitudin.
          </p>
        </div>

        {/* Navigation */}
        <nav className="border-t border-b border-gray-200">
          <div className="flex justify-center flex-wrap gap-x-19 gap-y-2 py-3">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`text-xs font-semibold tracking-wide hover:text-gray-500 transition-colors ${
                  pathname === item.href
                    ? "text-black border-b-2 border-black"
                    : "text-gray-700"
                }`}
              >
                {item.label}
              </Link>
            ))}
          </div>
        </nav>
      </div>
    </header>
  );
}