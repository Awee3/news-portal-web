'use client';

import { useState } from "react";
import Image from "next/image";
import Link from "next/link";

export default function Sidebar({ mainArticle = {}, thumbArticles = [] }) {
  // Defensive check: Pastikan mainArticle ada
  if (!mainArticle || !mainArticle.img) {
    return null;
  }

  // Ensure thumbArticles is array
  const safeThumbArticles = Array.isArray(thumbArticles) ? thumbArticles : [];

  // Ensure we have enough articles to show at least 2 pages (4 items) for the carousel
  // If we have 2 or fewer items, duplicate them to enable sliding behavior
  const displayArticles = safeThumbArticles.length > 0 && safeThumbArticles.length <= 2 
    ? [...safeThumbArticles, ...safeThumbArticles] 
    : safeThumbArticles;

  const [currentSlide, setCurrentSlide] = useState(0);
  const [direction, setDirection] = useState("next");
  const itemsPerPage = 2;
  const totalSlides = displayArticles.length > 0 ? Math.ceil(displayArticles.length / itemsPerPage) : 0;

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

  const goToSlide = (index) => {
    setDirection(index > currentSlide ? "next" : "prev");
    setCurrentSlide(index);
  };

  const currentItems = displayArticles.slice(
    currentSlide * itemsPerPage,
    (currentSlide + 1) * itemsPerPage
  );

  // Generate href for main article (support both slug and href)
  const mainHref = mainArticle.slug 
    ? `/article/${mainArticle.slug}` 
    : (mainArticle.href || '#');

  return (
    <aside className="w-full">
      {/* Main sidebar article */}
      <article className="mb-6 pb-6 border-b border-gray-300">
        <Link href={mainHref} className="group block">
          <div
            className="relative w-full overflow-hidden mb-3 rounded"
            style={{ aspectRatio: "312/203" }}
          >
            <Image
              src={mainArticle.img}
              alt={mainArticle.title || "Article"}
              fill
              className="object-cover transition-transform duration-300 group-hover:scale-105"
              sizes="(max-width:768px) 100vw, 33vw"
            />
          </div>
        </Link>
        <div className="w-full">
          <Link href={mainHref}>
            <h3 className="font-bold text-base leading-tight hover:underline mb-2 transition-colors">
              {mainArticle.title || "Untitled"}
            </h3>
          </Link>
          <p className="text-gray-600 text-sm leading-relaxed">
            {mainArticle.excerpt || ""}
          </p>
        </div>
      </article>

      {/* Sidebar thumbs carousel */}
      {displayArticles.length > 0 && (
        <>
          <div className="mb-6 relative w-full">
            <div
              key={currentSlide}
              className={`
                relative flex justify-between items-start
                animate-in fade-in duration-500
                ${direction === "prev" ? "slide-in-from-left-4" : "slide-in-from-right-4"}
              `}
              style={{ width: "372px", minHeight: "140px", gap: "16px" }}
            >
              {currentItems.map((article, i) => {
                // Generate href (support both slug and href)
                const articleHref = article.slug 
                  ? `/article/${article.slug}` 
                  : (article.href || '#');

                return (
                  <article key={`${currentSlide}-${i}`} className="flex-shrink-0">
                    <Link href={articleHref} className="group block">
                      <div
                        className="relative overflow-hidden mb-3 rounded"
                        style={{ width: "140px", height: "140px" }}
                      >
                        <Image
                          src={article.img || "/test.png"}
                          alt={article.title || "Article"}
                          width={140}
                          height={140}
                          className="object-cover w-full h-full transition-transform duration-300 group-hover:scale-110"
                          sizes="140px"
                        />
                      </div>
                    </Link>
                    <Link href={articleHref}>
                      <h4
                        className="font-bold text-xs leading-tight hover:underline line-clamp-3 transition-colors"
                        style={{ width: "140px" }}
                      >
                        {article.title || "Untitled"}
                      </h4>
                    </Link>
                  </article>
                );
              })}

              {/* Vertical divider centered between thumbnails */}
              {currentItems.length === 2 && (
                <div
                  className="pointer-events-none absolute left-1/2 -translate-x-1/2 w-0.5 bg-gray-300 z-10"
                  style={{ top: 0, height: "140px" }}
                />
              )}
            </div>
          </div>

          {/* Navigation controls */}
          <div className="flex items-center justify-center gap-3 mt-4 pb-4 border-b border-gray-300">
            <button
              onClick={prevSlide}
              disabled={totalSlides <= 1}
              className="w-6 h-6 text-sm rounded-full border border-gray-300 hover:bg-gray-100 hover:border-gray-400 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center transition-all duration-200 hover:scale-110 active:scale-95"
              aria-label="Previous page"
            >
              ‹
            </button>

            <div className="flex gap-2">
                {Array.from({ length: totalSlides }).map((_, index) => (
                  <button
                    key={index}
                    onClick={() => goToSlide(index)}
                    className={`
                      h-2 rounded-full transition-all duration-300
                      ${index === currentSlide 
                        ? 'bg-gray-800 w-6' 
                        : 'bg-gray-300 w-2 hover:bg-gray-400'
                      }
                    `}
                    aria-label={`Go to slide ${index + 1}`}
                    aria-current={index === currentSlide ? 'true' : 'false'}
                  />
                ))}
              </div>

            <button
              onClick={nextSlide}
              disabled={totalSlides <= 1}
              className="w-6 h-6 text-sm rounded-full border border-gray-300 hover:bg-gray-100 hover:border-gray-400 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center transition-all duration-200 hover:scale-110 active:scale-95"
              aria-label="Next page"
            >
              ›
            </button>
          </div>
        </>
      )}

      <style jsx>{`
        @keyframes slide-in-from-right {
          from {
            opacity: 0;
            transform: translateX(20px);
          }
          to {
            opacity: 1;
            transform: translateX(0);
          }
        }

        @keyframes slide-in-from-left {
          from {
            opacity: 0;
            transform: translateX(-20px);
          }
          to {
            opacity: 1;
            transform: translateX(0);
          }
        }

        .slide-in-from-right-4 {
          animation: slide-in-from-right 0.5s cubic-bezier(0.16, 1, 0.3, 1);
        }

        .slide-in-from-left-4 {
          animation: slide-in-from-left 0.5s cubic-bezier(0.16, 1, 0.3, 1);
        }
      `}</style>
    </aside>
  );
}