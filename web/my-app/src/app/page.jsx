import Header from "../components/Header";
import ContentBlock from "@/components/sections/ContentBlock";

export default function HomePage() {
  // Data untuk 3 blok (nanti akan diambil dari API backend)
  const contentBlocks = [
    // BLOK 1
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
            excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Semper Lacinia...",
            slug: "lorem-ipsum-dolor-sit-amet-1", // Tambahkan slug     
          },
          {
            title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
            excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Cras elementum libero ac...",
            slug: "lorem-ipsum-dolor-sit-amet-2", // Tambahkan slug
          },
        ],
        featured: {
          img: "/test.png",
          slug: "featured-article-1", // Tambahkan slug
          priority: true,
        }
      },
      carouselItems: [
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-1" },
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-2" },
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-3" },
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-4" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
            excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Semper Lacinia...",
            slug: "lorem-ipsum-dolor-sit-amet-3", // Tambahkan slug
          },
          {
            title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
            excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Cras elementum libero ac...",
            slug: "lorem-ipsum-dolor-sit-amet-4", // Tambahkan slug
          },
        ],
        featured: {
          img: "/test.png",
          slug: "featured-article-2", // Tambahkan slug
        }
      },
      sidebar: {
        mainArticle: {
          title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
          img: "/test.png",
          slug: "sidebar-main-1",
        },
        thumbArticles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-1", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-2", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-3", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-4", date: "Senin, 9 Des 2025" },
        ]
      }
    },

    // BLOK 2 (struktur sama, data berbeda)
    {
      id: 2,
      topMainSection: {
        sideTitles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 2.1...", slug: "lorem-ipsum-dolor-sit-amet-5" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 2.2...", slug: "lorem-ipsum-dolor-sit-amet-6" },
        ],
        featured: { img: "/test.png", slug: "featured-article-3" }
      },
      carouselItems: [
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-5" },
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-6" },
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-7" },
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-8" },
      ],
      bottomMainSection: {
        sideTitles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 2.3...", slug: "lorem-ipsum-dolor-sit-amet-7" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 2.4...", slug: "lorem-ipsum-dolor-sit-amet-8" },
        ],
        featured: { img: "/test.png", slug: "featured-article-4" }
      },
      sidebar: {
        mainArticle: { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "sidebar-main-2" },
        thumbArticles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-5", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-6", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-7", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-8", date: "Senin, 9 Des 2025" },
        ]
      }
    },

    // BLOK 3 (struktur sama, data berbeda)
    {
      id: 3,
      topMainSection: {
        sideTitles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 3.1...", slug: "lorem-ipsum-dolor-sit-amet-9" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 3.2...", slug: "lorem-ipsum-dolor-sit-amet-10" },
        ],
        featured: { img: "/test.png", slug: "featured-article-5" }
      },
      carouselItems: [
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-9" },
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-10" },
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-11" },
        { title: "Lorem Ipsum Dolor Sit Amet", img: "/test.png", slug: "carousel-article-12" },
      ],
      bottomMainSection: {
        sideTitles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 3.3...", slug: "lorem-ipsum-dolor-sit-amet-11" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 3.4...", slug: "lorem-ipsum-dolor-sit-amet-12" },
        ],
        featured: { img: "/test.png", slug: "featured-article-6" }
      },
      sidebar: {
        mainArticle: { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "sidebar-main-3" },
        thumbArticles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-9", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-10", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-11", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", slug: "thumb-article-12", date: "Senin, 9 Des 2025" },
        ]
      }
    }
  ];

  return (
    <div className="min-h-screen bg-white">
      <Header />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {contentBlocks.map((block) => (
          <ContentBlock
            key={block.id}
            topMainSection={block.topMainSection}
            carouselItems={block.carouselItems}
            bottomMainSection={block.bottomMainSection}
            sidebar={block.sidebar}
          />
        ))}
      </main>
    </div>
  );
}
