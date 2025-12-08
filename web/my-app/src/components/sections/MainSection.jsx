import Image from "next/image";
import Link from "next/link";

/**
 * Main section with 2 text articles (left) + 1 big image (right)
 * @param {Array} sideTitles - Array of 2 article objects with title, excerpt, href
 * @param {Object} featured - Featured image object with img, href
 */

export default function MainSection({ sideTitles, featured }) {
  return (
    <section className="grid grid-cols-1 md:grid-cols-12 gap-6 mb-8 pb-8 border-b border-gray-200">
      {/* LEFT: 2 text-only articles - flexbox untuk distribute space */}
      <div className="md:col-span-5 flex flex-col h-full">
        {sideTitles.map((article, i) => (
          <article 
            key={article.href} 
            className={`flex-1 flex flex-col justify-center ${
              i < sideTitles.length - 1 ? "border-b border-gray-200" : ""
            }`}
          >
            <Link href={article.href}>
              <h2 className="text-xl font-bold leading-snug hover:underline">
                {article.title}
              </h2>
            </Link>
            <p className="text-gray-600 text-sm mt-2">{article.excerpt}</p>
          </article>
        ))}
      </div>

      {/* RIGHT: Big image only */}
      <div className="md:col-span-7">
        <Link href={featured.href}>
          <div className="relative w-full overflow-hidden" style={{ aspectRatio: '411/268' }}>
            <Image
              src={featured.img}
              alt="Featured article"
              width={411}
              height={268}
              className="object-cover w-full h-full"
              sizes="(max-width:768px) 100vw, 60vw"
              priority={featured.priority || false}
            />
          </div>
        </Link>
      </div>
    </section>
  );
}