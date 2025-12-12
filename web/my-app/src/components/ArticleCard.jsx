import Image from "next/image";
import Link from "next/link";

/**
 * @param {{ slug:string; title:string; thumbnail?:string; publishedAt:string; excerpt?:string }} props
 */
export default function ArticleCard({ slug, title, thumbnail = "/images/placeholders/placeholder.jpg", publishedAt, excerpt }) {
  return (
    <article className="border border-gray-200 rounded overflow-hidden hover:shadow-sm transition">
      <Link href={`/${slug}`}>
        <div className="relative w-full h-44 bg-gray-100">
          <Image
            src={thumbnail}
            alt={title}
            fill
            sizes="(max-width:768px) 100vw, (max-width:1200px) 33vw, 25vw"
            className="object-cover"
          />
        </div>
        <div className="p-4">
          <h2 className="font-semibold text-sm mb-2 line-clamp-2">{title}</h2>
          {excerpt && <p className="text-xs text-gray-600 line-clamp-2">{excerpt}</p>}
          <p className="text-[11px] text-gray-400 mt-2">
            {new Date(publishedAt).toLocaleDateString("id-ID", { day: "2-digit", month: "short", year: "numeric" })}
          </p>
        </div>
      </Link>
    </article>
  );
}