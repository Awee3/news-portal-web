import Image from "next/image";
import Link from "next/link";
import Header from "@/components/Header";
import CommentSection from "@/components/sections/CommentSection";

// Fetch artikel berdasarkan slug
async function getArticle(slug) {
  try {
    const res = await fetch(`http://localhost:8080/api/v1/articles/slug/${slug}`, {
      next: { revalidate: 60 }
    });
    
    if (!res.ok) {
      if (res.status === 404) throw new Error('Article not found');
      throw new Error('Failed to fetch article');
    }
    
    return res.json();
  } catch (error) {
    console.error('Error fetching article:', error);
    // Return dummy data for development
    return {
      artikel_id: 1,
      judul: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Duis Sollicitudin.",
      slug: slug,
      konten: `<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aenean quis magna in urna semper lacinia. Cras elementum, libero ac consectetur adipiscing elit. Fusce vel dolor nec nunc ultricies tincidunt.</p>
      
      <p>Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
      
      <h2>Subtitle Section Here</h2>
      
      <p>Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>`,
      excerpt: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aenean quis magna in urna semper lacinia.",
      gambar_utama: "/test.png",
      penulis: "John Doe",
      tanggal_publikasi: "2025-12-09T10:30:00Z",
      kategori: [
        { kategori_id: 1, nama_kategori: "International" }
      ],
      tags: [
        { tag_id: 1, nama: "Politics" },
        { tag_id: 2, nama: "Economy" }
      ]
    };
  }
}

// Fetch comments untuk artikel
async function getComments(articleId) {
  try {
    const res = await fetch(`http://localhost:8080/api/v1/articles/${articleId}/comments`, {
      next: { revalidate: 10 }
    });
    
    if (!res.ok) throw new Error('Failed to fetch comments');
    return res.json();
  } catch (error) {
    console.error('Error fetching comments:', error);
    // Return dummy data
    return [
      {
        comment_id: 1,
        artikel_id: articleId,
        user_id: null,
        nama_pengguna: "Anonymous User",
        konten: "Great article! Very informative and well-written.",
        tanggal_dibuat: "2025-12-10T08:30:00Z"
      },
      {
        comment_id: 2,
        artikel_id: articleId,
        user_id: 123,
        nama_pengguna: "Jane Smith",
        konten: "I found this really helpful. Thanks for sharing!",
        tanggal_dibuat: "2025-12-10T10:15:00Z"
      },
      {
        comment_id: 3,
        artikel_id: articleId,
        user_id: null,
        nama_pengguna: "Reader",
        konten: "Looking forward to more content like this.",
        tanggal_dibuat: "2025-12-10T14:45:00Z"
      }
    ];
  }
}

// Format tanggal ke Indonesian locale
function formatDate(dateString) {
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
  const { slug } = params;
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
              <span className="font-semibold">By {article.penulis}</span>
              <span>•</span>
              <time dateTime={article.tanggal_publikasi}>
                {formatDate(article.tanggal_publikasi)}
              </time>
              <span>•</span>
              <span>{article.kategori[0]?.nama_kategori}</span>
            </div>

            {/* Featured Image */}
            <div className="relative w-full aspect-[16/9] overflow-hidden rounded-lg mb-6">
              <Image
                src={article.gambar_utama}
                alt={article.judul}
                fill
                className="object-cover"
                sizes="(max-width: 768px) 100vw, 896px"
                priority
              />
            </div>

            {/* Excerpt */}
            <p className="text-xl text-gray-700 leading-relaxed font-medium">
              {article.excerpt}
            </p>
          </header>

          {/* Article Body */}
          <div 
            className="prose prose-lg max-w-none
              prose-headings:font-serif prose-headings:font-bold
              prose-h2:text-3xl prose-h2:mt-12 prose-h2:mb-6
              prose-p:text-gray-800 prose-p:leading-relaxed prose-p:mb-6
              prose-blockquote:border-l-4 prose-blockquote:border-gray-800 
              prose-blockquote:pl-6 prose-blockquote:italic prose-blockquote:text-gray-700
              prose-a:text-blue-600 prose-a:no-underline hover:prose-a:underline
              prose-img:rounded-lg prose-img:shadow-md"
            dangerouslySetInnerHTML={{ __html: article.konten }}
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

          {/* Share Buttons */}
          <div className="mt-8 pt-8 border-t border-gray-200">
            <p className="text-sm font-semibold text-gray-600 mb-3">Share this article:</p>
            <div className="flex gap-3">
              <button className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded transition-colors">
                Facebook
              </button>
              <button className="px-4 py-2 bg-sky-500 hover:bg-sky-600 text-white text-sm rounded transition-colors">
                Twitter
              </button>
              <button className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white text-sm rounded transition-colors">
                WhatsApp
              </button>
            </div>
          </div>

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
  const { slug } = params;
  const article = await getArticle(slug);
  
  if (!article) {
    return {
      title: 'Article Not Found | Bintaro Times',
    };
  }
  
  return {
    title: `${article.judul} | Bintaro Times`,
    description: article.excerpt,
    openGraph: {
      title: article.judul,
      description: article.excerpt,
      images: [article.gambar_utama],
      type: 'article',
      publishedTime: article.tanggal_publikasi,
      authors: [article.penulis],
    },
  };
}