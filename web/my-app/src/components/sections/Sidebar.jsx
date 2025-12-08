import Image from "next/image";
import Link from "next/link";

/**
 * Right sidebar with main article + 2 thumbnails + pagination
 * @param {Object} mainArticle - Main sidebar article with img, title, excerpt, href
 * @param {Array} thumbArticles - Array of 2 thumbnail articles
 */
export default function Sidebar({ mainArticle, thumbArticles }) {
  return (
    <aside className="col-span-12 lg:col-span-4">
      {/* Main sidebar article */}
      <article className="mb-6 pb-6 border-b border-gray-200">
        <Link href={mainArticle.href}>
          <div className="relative w-full overflow-hidden mb-3" style={{ aspectRatio: '312/203' }}>
            <Image
              src={mainArticle.img}
              alt={mainArticle.title}
              width={312}
              height={203}
              className="object-cover w-full h-full"
              sizes="(max-width:768px) 100vw, 33vw"
            />
          </div>
        </Link>
        <Link href={mainArticle.href}>
          <h3 className="font-bold text-base leading-tight hover:underline mb-2">
            {mainArticle.title}
          </h3>
        </Link>
        <p className="text-gray-600 text-sm leading-relaxed">
          {mainArticle.excerpt}
        </p>
      </article>

      {/* Sidebar thumbs grid */}
      <div className="grid grid-cols-2 gap-6 mb-6">
        {thumbArticles.map((article, i) => (
          <article key={article.href} className="relative">
            <Link href={article.href}>
              <div className="relative w-full overflow-hidden mb-3" style={{ aspectRatio: '1/1' }}>
                <Image
                  src={article.img}
                  alt={article.title}
                  width={140}
                  height={140}
                  className="object-cover w-full h-full"
                  sizes="140px"
                />
              </div>
            </Link>
            <Link href={article.href}>
              <h4 className="font-bold text-xs leading-tight hover:underline line-clamp-3">
                {article.title}
              </h4>
            </Link>
            {/* Vertical divider */}
            {i === 0 && (
              <div className="absolute top-0 bottom-0 -right-3 w-px bg-gray-200" />
            )}
          </article>
        ))}
      </div>

      {/* Pagination dots */}
      <div className="flex items-center justify-center gap-2 mt-4 pb-4 border-b border-gray-200">
        <button
          className="w-6 h-6 text-xs rounded-full border border-gray-300 hover:bg-gray-100 flex items-center justify-center"
          aria-label="Previous page"
        >
          ‹
        </button>
        <div className="flex gap-1">
          <span className="w-2 h-2 rounded-full bg-gray-800" />
          <span className="w-2 h-2 rounded-full bg-gray-300" />
          <span className="w-2 h-2 rounded-full bg-gray-300" />
        </div>
        <button
          className="w-6 h-6 text-xs rounded-full border border-gray-300 hover:bg-gray-100 flex items-center justify-center"
          aria-label="Next page"
        >
          ›
        </button>
      </div>
    </aside>
  );
}