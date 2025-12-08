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
  const itemsPerPage = 2;
  const totalSlides = Math.ceil(items.length / itemsPerPage);

  const nextSlide = () => {
    setCurrentSlide((prev) => (prev + 1) % totalSlides);
  };

  const prevSlide = () => {
    setCurrentSlide((prev) => (prev - 1 + totalSlides) % totalSlides);
  };

  const goToSlide = (index) => {
    setCurrentSlide(index);
  };

  const currentItems = items.slice(
    currentSlide * itemsPerPage,
    (currentSlide + 1) * itemsPerPage
  );

  return (
    <section className="mb-8 pb-8 border-b border-gray-200">
      <div className="relative">
        {/* Carousel items */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 items-stretch">
          {currentItems.map((article, i) => (
            <div key={`${currentSlide}-${i}`} className="relative">
              <Link href={article.href} className="group">
                <div className="flex items-center gap-4">
                  <div 
                    className="relative flex-shrink-0 overflow-hidden" 
                    style={{ width: '103px', height: '103px' }}
                  >
                    <Image
                      src={article.img}
                      alt={article.title}
                      width={103}
                      height={103}
                      className="object-cover w-full h-full"
                      sizes="103px"
                    />
                  </div>
                  <h3 className="font-semibold text-sm leading-snug group-hover:underline">
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
        <div className="flex items-center justify-center gap-3 mt-4">
          <button
            onClick={prevSlide}
            className="w-6 h-6 text-xs rounded-full border border-gray-300 hover:bg-gray-100 flex items-center justify-center transition"
            aria-label="Previous"
          >
            ‹
          </button>

          <div className="flex gap-2">
            {Array.from({ length: totalSlides }).map((_, index) => (
              <button
                key={index}
                onClick={() => goToSlide(index)}
                className={`w-2 h-2 rounded-full transition ${
                  index === currentSlide ? 'bg-gray-800' : 'bg-gray-300'
                }`}
                aria-label={`Go to slide ${index + 1}`}
              />
            ))}
          </div>

          <button
            onClick={nextSlide}
            className="w-6 h-6 text-xs rounded-full border border-gray-300 hover:bg-gray-100 flex items-center justify-center transition"
            aria-label="Next"
          >
            ›
          </button>
        </div>
      </div>
    </section>
  );
}