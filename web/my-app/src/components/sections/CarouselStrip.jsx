'use client';

import { useState } from "react";
import Image from "next/image";
import Link from "next/link";

export default function CarouselStrip({ items = [] }) {
  // Handle empty items
  if (!items || items.length === 0) {
    return null;
  }

  // Duplicate items if count is low
  const displayItems = items.length <= 2 ? [...items, ...items] : items;

  const [currentSlide, setCurrentSlide] = useState(0);
  const [direction, setDirection] = useState("next");
  const itemsPerPage = 2;
  
  const totalSlides = Math.ceil(displayItems.length / itemsPerPage);

  const nextSlide = () => {
    if (totalSlides <= 1) return;
    setDirection("next");
    setCurrentSlide((prev) => (prev + 1) % totalSlides);
  };
  
  const prevSlide = () => {
    if (totalSlides <= 1) return;
    setDirection("prev");
    setCurrentSlide((prev) => (prev - 1 + totalSlides) % totalSlides);
  };

  const currentItems = displayItems.slice(
    currentSlide * itemsPerPage,
    (currentSlide + 1) * itemsPerPage
  );

  return (
    <section className="mb-8 pb-8 border-b border-gray-300">
      <div className="relative overflow-hidden">
        <div 
          key={currentSlide}
          className={`
            grid grid-cols-1 sm:grid-cols-2 gap-6 items-stretch
            animate-in fade-in slide-in-from-right-4 duration-500
            ${direction === 'prev' ? 'slide-in-from-left-4' : 'slide-in-from-right-4'}
          `}
        >
          {currentItems.map((article, i) => {
            // Ensure href is always a string
            const href = article.slug 
              ? `/article/${article.slug}` 
              : (article.href || '#');

            return (
              <div key={`${currentSlide}-${i}`} className="relative">
                <Link href={href} className="group">
                  <div className="flex items-center gap-4">
                    <div 
                      className="relative flex-shrink-0 overflow-hidden rounded" 
                      style={{ width: '103px', height: '103px' }}
                    >
                      <Image
                        src={article.img || "/test.png"}
                        alt={article.title || "Article"}
                        width={103}
                        height={103}
                        className="object-cover w-full h-full transition-transform duration-300 group-hover:scale-110"
                        sizes="103px"
                      />
                    </div>
                    <h3 className="font-semibold text-sm leading-snug group-hover:underline transition-colors">
                      {article.title || "Untitled"}
                    </h3>
                  </div>
                </Link>
                {i === 0 && currentItems.length > 1 && (
                  <div 
                    className="absolute right-0 top-0 bottom-0 w-0.5 bg-gray-300 hidden sm:block"
                    style={{ right: '-12px' }}
                  ></div>
                )}
              </div>
            );
          })}
        </div>

        {/* Navigation */}
        {totalSlides > 1 && (
          <div className="flex items-center justify-center gap-3 mt-6">
            <button
              onClick={prevSlide}
              className="w-8 h-8 rounded-full border border-gray-300 hover:bg-gray-100 hover:border-gray-400 flex items-center justify-center transition-all duration-200 hover:scale-110 active:scale-95"
            >
              ‹
            </button>

            <div className="flex gap-2">
              {Array.from({ length: totalSlides }).map((_, index) => (
                <button
                  key={index}
                  onClick={() => {
                    setDirection(index > currentSlide ? "next" : "prev");
                    setCurrentSlide(index);
                  }}
                  className={`h-1.5 rounded-full transition-all duration-300 ${
                    currentSlide === index ? "w-6 bg-black" : "w-1.5 bg-gray-300 hover:bg-gray-400"
                  }`}
                />
              ))}
            </div>

            <button
              onClick={nextSlide}
              className="w-8 h-8 rounded-full border border-gray-300 hover:bg-gray-100 hover:border-gray-400 flex items-center justify-center transition-all duration-200 hover:scale-110 active:scale-95"
            >
              ›
            </button>
          </div>
        )}
      </div>

      <style jsx>{`
        .slide-in-from-right-4 { animation: slideInRight 0.4s ease-out forwards; }
        .slide-in-from-left-4 { animation: slideInLeft 0.4s ease-out forwards; }
        
        @keyframes slideInRight {
          from { opacity: 0; transform: translateX(20px); }
          to { opacity: 1; transform: translateX(0); }
        }
        @keyframes slideInLeft {
          from { opacity: 0; transform: translateX(-20px); }
          to { opacity: 1; transform: translateX(0); }
        }
      `}</style>
    </section>
  );
}