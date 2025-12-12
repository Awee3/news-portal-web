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
            slug: "government-infrastructure-development-plan",
          },
          {
            title: "National Education Reform Takes Effect",
            excerpt: "New curriculum standards designed to enhance learning outcomes nationwide.",
            slug: "national-education-reform",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "national-featured-article",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Healthcare System Modernization Underway", img: "/test.png", slug: "healthcare-system-modernization" },
        { title: "Tourism Sector Shows Strong Recovery", img: "/test.png", slug: "tourism-sector-recovery" },
        { title: "Agricultural Export Reaches Record High", img: "/test.png", slug: "agricultural-export-record" },
        { title: "Digital Infrastructure Expansion Announced", img: "/test.png", slug: "digital-infrastructure-expansion" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "National Census Results Released",
            excerpt: "Population data reveals demographic shifts and urbanization trends.",
            slug: "national-census-results",
          },
          {
            title: "Environmental Protection Law Strengthened",
            excerpt: "New regulations aim to preserve natural resources for future generations.",
            slug: "environmental-protection-law",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "census-results-featured",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Parliament Debates Budget Allocation",
          img: "/test.png",
          slug: "parliament-budget-debate",
        },
        thumbArticles: [
          { title: "Transportation Hub Opens in Capital", img: "/test.png", slug: "transportation-hub-opens", date: "Senin, 9 Des 2025" },
          { title: "National Athletes Prepare for Games", img: "/test.png", slug: "national-athletes-prepare", date: "Senin, 9 Des 2025" },
          { title: "Cultural Heritage Sites Restored", img: "/test.png", slug: "cultural-heritage-restored", date: "Minggu, 8 Des 2025" },
          { title: "Renewable Energy Target Announced", img: "/test.png", slug: "renewable-energy-target", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="NATIONAL" contentBlocks={contentBlocks} />;
}