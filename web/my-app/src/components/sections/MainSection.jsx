import Image from "next/image";
import Link from "next/link";

export default function MainSection({ sideTitles = [], featured = {} }) {
  // Defensive check
  if (!featured.img || !sideTitles || sideTitles.length === 0) {
    return null;
  }

  return (
    <section className="grid grid-cols-1 md:grid-cols-12 gap-6 mb-8 pb-8 border-b border-gray-300">
      <div className="md:col-span-5 flex flex-col h-full">
        {sideTitles.slice(0, 2).map((article, i) => {
          // Generate href, ensuring it's always a string
          const href = article.slug
            ? `/article/${article.slug}`
            : article.href || "#";

          return (
            <article
              key={article.slug || article.href || i}
              className={`flex-1 flex flex-col justify-center ${
                i < 1 ? "border-b border-gray-300 pb-4" : "pt-4"
              }`}
            >
              <Link href={href}>
                <h2 className="text-xl font-bold leading-snug hover:underline">
                  {article.title || "Untitled"}
                </h2>
              </Link>
              <p className="text-gray-600 text-sm mt-2 line-clamp-3">
                {article.excerpt || ""}
              </p>
            </article>
          );
        })}
      </div>

      <div className="md:col-span-7">
        {/* Generate featured href */}
        {(() => {
          const featuredHref = featured.slug
            ? `/article/${featured.slug}`
            : featured.href || "#";

          return (
            <Link href={featuredHref}>
              <div
                className="relative w-full overflow-hidden rounded bg-gray-100"
                style={{ aspectRatio: "411/268" }}
              >
                <Image
                  src={featured.img || "/test.png"}
                  alt={featured.title || "Featured article"}
                  fill
                  className="object-cover hover:scale-105 transition-transform duration-500"
                  sizes="(max-width:768px) 100vw, 60vw"
                  priority={featured.priority || false}
                />
              </div>
            </Link>
          );
        })()}
      </div>
    </section>
  );
}