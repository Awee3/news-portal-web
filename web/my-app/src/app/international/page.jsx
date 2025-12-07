import ArticleCard from "../../components/ArticleCard";
export const metadata = { title: "International - Bintaro Times" };


export default function InternationalPage() {
  const articles = [
    {
      id: "1",
      slug: "international-summit",
      title: "International Summit Reaches Agreement",
      excerpt: "World leaders agree on new climate framework...",
      thumbnail: "/images/placeholders/international-1.jpg",
      publishedAt: "2025-09-09T07:30:00Z"
    },
    {
      id: "2",
      slug: "regional-security-talks",
      title: "Regional Security Talks Continue Amid Tension",
      excerpt: "Diplomats push for de‑escalation measures...",
      thumbnail: "/images/placeholders/international-2.jpg",
      publishedAt: "2025-09-09T06:10:00Z"
    }
  ];

  return (
    <main className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">International</h1>
      <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 md:grid-cols-3">
        {articles.map(a => <ArticleCard key={a.id} {...a} />)}
      </div>
    </main>
  );
}