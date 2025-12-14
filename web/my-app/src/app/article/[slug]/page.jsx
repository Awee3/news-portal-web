import Image from "next/image";
import Link from "next/link";
import Header from "@/components/Header";
import CommentSection from "@/components/sections/CommentSection";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
const API_BASE = API_URL.replace("/api/v1", "");

function getImageUrl(imagePath) {
  if (!imagePath) return "/test.png";
  if (imagePath.startsWith("http")) return imagePath;
  return `${API_BASE}${imagePath}`;
}

// NEW: Helper untuk format konten artikel menjadi HTML paragraf
// Ganti fungsi yang lama dengan yang ini
function formatArticleContent(content) {
  if (!content) return "";
  
  // Jika sudah ada tag HTML, return as is
  if (content.includes("<p>") || content.includes("<h2>") || content.includes("<div>")) {
    return content;
  }
  
  // REVISI: Split berdasarkan 1 enter saja (\n)
  // Setiap baris baru akan langsung jadi <p> tersendiri
  return content
    .split(/\n+/) // Split setiap ketemu enter (satu atau lebih)
    .filter(p => p.trim()) // Hapus yang kosong
    .map(p => {
      // Kita bungkus setiap pecahan teks langsung menjadi <p>
      return `<p>${p.trim()}</p>`;
    })
    .join('\n');
}

async function getArticle(slug) {
  try {
    const res = await fetch(`${API_URL}/articles/slug/${slug}`, {
      next: { revalidate: 60 }
    });
    
    if (!res.ok) {
      return null;
    }
    
    return res.json();
  } catch (error) {
    console.error('Error fetching article:', error);
    return null;
  }
}

async function getComments(articleId) {
  try {
    const res = await fetch(`${API_URL}/articles/${articleId}/comments`, {
      next: { revalidate: 10 }
    });
    
    if (!res.ok) return [];
    return res.json();
  } catch (error) {
    console.error('Error fetching comments:', error);
    return [];
  }
}

function formatDate(dateString) {
  if (!dateString) return "";
  const date = new Date(dateString);
  const options = { 
    weekday: 'long', 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  };
  return date.toLocaleDateString('id-ID', options);
}

export default async function ArticlePage({ params }) {
  const { slug } = await params;
  const article = await getArticle(slug);
  
  if (!article) {
    return (
      <div className="min-h-screen bg-white">
        <Header />
        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 text-center">
          <h1 className="text-4xl font-bold mb-4">Article Not Found</h1>
          <p className="text-gray-600 mb-8">The article you're looking for doesn't exist.</p>
          <Link href="/" className="text-blue-600 hover:underline">
            Return to Home
          </Link>
        </main>
      </div>
    );
  }
  
  const comments = await getComments(article.artikel_id);

  return (
    <div className="min-h-screen bg-white">
      <Header />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Breadcrumb */}
        <nav className="flex items-center gap-2 text-sm text-gray-600 mb-6">
          <Link href="/" className="hover:text-black transition-colors">
            Home
          </Link>
          <span>/</span>
          {article.kategori && article.kategori.length > 0 && (
            <>
              <Link 
                href={`/${article.kategori[0].nama_kategori.toLowerCase()}`}
                className="hover:text-black transition-colors capitalize"
              >
                {article.kategori[0].nama_kategori}
              </Link>
              <span>/</span>
            </>
          )}
          <span className="text-gray-400">{article.judul.substring(0, 30)}...</span>
        </nav>

        {/* Main Article Content */}
        <article className="max-w-4xl mx-auto">
          {/* Article Header */}
          <header className="mb-8">
            <h1 className="text-4xl md:text-5xl font-bold font-serif leading-tight mb-4">
              {article.judul}
            </h1>
            
            {/* Meta Information */}
            <div className="flex flex-wrap items-center gap-4 text-sm text-gray-600 mb-6">
              <span className="font-semibold">By {article.penulis || "Admin"}</span>
              <span>•</span>
              <time dateTime={article.tanggal_publikasi}>
                {formatDate(article.tanggal_publikasi)}
              </time>
              {article.kategori && article.kategori.length > 0 && (
                <>
                  <span>•</span>
                  <span>{article.kategori[0]?.nama_kategori}</span>
                </>
              )}
            </div>

            {/* Featured Image */}
            {article.gambar_utama && (
              <div className="relative w-full aspect-[16/9] overflow-hidden rounded-lg mb-6">
                <Image
                  src={getImageUrl(article.gambar_utama)}
                  alt={article.judul}
                  fill
                  className="object-cover"
                  sizes="(max-width: 768px) 100vw, 896px"
                  priority
                />
              </div>
            )}

            {/* Excerpt */}
            {article.ringkasan && (
              <p className="text-xl text-gray-700 leading-relaxed font-medium">
                {article.ringkasan}
              </p>
            )}
          </header>

          {/* Article Body - UPDATED */}
          <div 
            className="prose prose-lg max-w-none
              prose-headings:font-serif prose-headings:font-bold
              prose-h2:text-3xl prose-h2:mt-12 prose-h2:mb-6
              prose-p:text-gray-800 prose-p:leading-relaxed prose-p:mb-6
              prose-blockquote:border-l-4 prose-blockquote:border-gray-800 
              prose-blockquote:pl-6 prose-blockquote:italic prose-blockquote:text-gray-700
              prose-a:text-blue-600 prose-a:no-underline hover:prose-a:underline
              prose-img:rounded-lg prose-img:shadow-md"
            dangerouslySetInnerHTML={{ __html: formatArticleContent(article.konten) }}
          />

          {/* Tags */}
          {article.tags && article.tags.length > 0 && (
            <div className="mt-12 pt-8 border-t border-gray-200">
              <div className="flex flex-wrap gap-2">
                <span className="text-sm font-semibold text-gray-600">Tags:</span>
                {article.tags.map((tag) => (
                  <Link
                    key={tag.tag_id}
                    href={`/tag/${tag.nama.toLowerCase()}`}
                    className="px-3 py-1 bg-gray-100 hover:bg-gray-200 text-gray-700 text-sm rounded-full transition-colors"
                  >
                    {tag.nama}
                  </Link>
                ))}
              </div>
            </div>
          )}


          {/* Comment Section */}
          <CommentSection 
            articleId={article.artikel_id} 
            initialComments={comments}
          />
        </article>
      </main>
    </div>
  );
}

// Generate metadata untuk SEO
export async function generateMetadata({ params }) {
  const { slug } = await params;
  const article = await getArticle(slug);
  
  if (!article) {
    return {
      title: 'Article Not Found',
    };
  }
  
  return {
    title: article.judul,
    description: article.ringkasan || article.konten?.substring(0, 160),
    openGraph: {
      title: article.judul,
      description: article.ringkasan,
      images: [getImageUrl(article.gambar_utama)],
      type: 'article',
      publishedTime: article.tanggal_publikasi,
      authors: [article.penulis || "Admin"],
    },
  };
}