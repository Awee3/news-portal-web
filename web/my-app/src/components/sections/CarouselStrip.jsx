'use client';

import { useState } from "react";
import Image from "next/image";
import Link from "next/link";

/**
 * Carousel strip showing 2 items per slide with thumbnail + title
 * @param {Array} items - Array of article objects with img, title, href
 */
export default function CarouselStrip({ items }) {
  const [currentSlide, setCurrentSlide] = useState(0);
  const [direction, setDirection] = useState('next'); // Track animation direction
  const itemsPerPage = 2;
  const totalSlides = Math.ceil(items.length / itemsPerPage);

  const nextSlide = () => {
    setDirection('next');
    setCurrentSlide((prev) => (prev + 1) % totalSlides);
  };

  const prevSlide = () => {
    setDirection('prev');
    setCurrentSlide((prev) => (prev - 1 + totalSlides) % totalSlides);
  };

  const goToSlide = (index) => {
    setDirection(index > currentSlide ? 'next' : 'prev');
    setCurrentSlide(index);
  };

  const currentItems = items.slice(
    currentSlide * itemsPerPage,
    (currentSlide + 1) * itemsPerPage
  );

  return (
    <section className="mb-8 pb-8 border-b border-gray-200">
      <div className="relative overflow-hidden">
        {/* Carousel items with animation */}
        <div 
          key={currentSlide}
          className={`
            grid grid-cols-1 sm:grid-cols-2 gap-6 items-stretch
            animate-in fade-in slide-in-from-right-4 duration-500
            ${direction === 'prev' ? 'slide-in-from-left-4' : 'slide-in-from-right-4'}
          `}
        >
          {currentItems.map((article, i) => (
            <div key={`${currentSlide}-${i}`} className="relative">
              <Link href={article.href} className="group">
                <div className="flex items-center gap-4">
                  <div 
                    className="relative flex-shrink-0 overflow-hidden rounded" 
                    style={{ width: '103px', height: '103px' }}
                  >
                    <Image
                      src={article.img}
                      alt={article.title}
                      width={103}
                      height={103}
                      className="object-cover w-full h-full transition-transform duration-300 group-hover:scale-110"
                      sizes="103px"
                    />
                  </div>
                  <h3 className="font-semibold text-sm leading-snug group-hover:underline transition-colors">
                    {article.title}
                  </h3>
                </div>
              </Link>
              {/* Vertical divider between items */}
              {i === 0 && (
                <div className="hidden sm:block absolute top-0 bottom-0 -right-3 w-px bg-gray-200" />
              )}
            </div>
          ))}
        </div>

        {/* Navigation controls */}
        <div className="flex items-center justify-center gap-3 mt-6">
          <button
            onClick={prevSlide}
            disabled={items.length <= itemsPerPage}
            className="w-6 h-6 text-sm rounded-full border border-gray-300 hover:bg-gray-100 hover:border-gray-400 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center transition-all duration-200 hover:scale-110 active:scale-95"
            aria-label="Previous"
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
            disabled={items.length <= itemsPerPage}
            className="w-6 h-6 text-sm rounded-full border border-gray-300 hover:bg-gray-100 hover:border-gray-400 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center transition-all duration-200 hover:scale-110 active:scale-95"
            aria-label="Next"
          >
            ›
          </button>
        </div>
      </div>

      <style jsx>{`
        @keyframes slide-in-from-right {
          from {
            opacity: 0;
            transform: translateX(10px);
          }
          to {
            opacity: 1;
            transform: translateX(0);
          }
        }

        @keyframes slide-in-from-left {
          from {
            opacity: 0;
            transform: translateX(-10px);
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
    </section>
  );
}