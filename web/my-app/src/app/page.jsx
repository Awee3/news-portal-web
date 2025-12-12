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
            href: "/international/article-1",
          },
          {
            title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
            excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Cras elementum libero ac...",
            href: "/business/article-2",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/featured/1",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
            excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Semper Lacinia...",
            href: "/sports/article-1",
          },
          {
            title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
            excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Cras elementum libero ac...",
            href: "/sports/article-2",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/featured/2",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
          img: "/test.png",
          href: "/sidebar/main-1",
        },
        thumbArticles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
        ]
      }
    },

    // BLOK 2 (struktur sama, data berbeda)
    {
      id: 2,
      topMainSection: {
        sideTitles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 2.1...", href: "/article-2-1" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 2.2...", href: "/article-2-2" },
        ],
        featured: { img: "/test.png", href: "/featured/3" }
      },
      carouselItems: [
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
      ],
      bottomMainSection: {
        sideTitles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 2.3...", href: "/article-2-3" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 2.4...", href: "/article-2-4" },
        ],
        featured: { img: "/test.png", href: "/featured/4" }
      },
      sidebar: {
        mainArticle: { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "/sidebar/main-2" },
        thumbArticles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
        ]
      }
    },

    // BLOK 3 (struktur sama, data berbeda)
    {
      id: 3,
      topMainSection: {
        sideTitles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 3.1...", href: "/article-3-1" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 3.2...", href: "/article-3-2" },
        ],
        featured: { img: "/test.png", href: "/featured/5" }
      },
      carouselItems: [
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
        { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#" },
      ],
      bottomMainSection: {
        sideTitles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 3.3...", href: "/article-3-3" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", excerpt: "Excerpt 3.4...", href: "/article-3-4" },
        ],
        featured: { img: "/test.png", href: "/featured/6" }
      },
      sidebar: {
        mainArticle: { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "/sidebar/main-3" },
        thumbArticles: [
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
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
