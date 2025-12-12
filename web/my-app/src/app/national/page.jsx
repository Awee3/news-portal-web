import CategoryPageTemplate from "@/components/CategoryPageTemplate";

export default function NationalPage() {
  const contentBlocks = [
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "Government Announces New Infrastructure Development Plan",
            excerpt: "Ambitious project aims to improve connectivity across major cities and rural areas.",
            href: "/national/infrastructure",
          },
          {
            title: "National Education Reform Takes Effect",
            excerpt: "New curriculum standards designed to enhance learning outcomes nationwide.",
            href: "/national/education-reform",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/national/featured-1",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Healthcare System Modernization Underway", img: "/test.png", href: "/national/healthcare" },
        { title: "Tourism Sector Shows Strong Recovery", img: "/test.png", href: "/national/tourism" },
        { title: "Agricultural Export Reaches Record High", img: "/test.png", href: "/national/agriculture" },
        { title: "Digital Infrastructure Expansion Announced", img: "/test.png", href: "/national/digital" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "National Census Results Released",
            excerpt: "Population data reveals demographic shifts and urbanization trends.",
            href: "/national/census",
          },
          {
            title: "Environmental Protection Law Strengthened",
            excerpt: "New regulations aim to preserve natural resources for future generations.",
            href: "/national/environment",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/national/featured-2",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Parliament Debates Budget Allocation",
          img: "/test.png",
          href: "/national/budget",
        },
        thumbArticles: [
          { title: "Transportation Hub Opens in Capital", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "National Athletes Prepare for Games", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Cultural Heritage Sites Restored", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
          { title: "Renewable Energy Target Announced", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="NATIONAL" contentBlocks={contentBlocks} />;
}