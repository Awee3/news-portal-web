import CategoryPageTemplate from "@/components/CategoryPageTemplate";

export default function SciTechPage() {
  const contentBlocks = [
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "Breakthrough in Quantum Computing Research",
            excerpt: "Scientists achieve major milestone in development of next-generation processors.",
            slug: "quantum-computing-breakthrough",
          },
          {
            title: "AI System Demonstrates Advanced Problem-Solving",
            excerpt: "Machine learning model surpasses expectations in complex cognitive tasks.",
            slug: "ai-system-problem-solving",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "scitech-featured-article",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Space Mission Discovers New Exoplanet", img: "/test.png", slug: "space-mission-exoplanet" },
        { title: "5G Network Coverage Expands Nationwide", img: "/test.png", slug: "5g-network-coverage" },
        { title: "Medical Breakthrough in Cancer Treatment", img: "/test.png", slug: "cancer-treatment-breakthrough" },
        { title: "Renewable Energy Technology Advances", img: "/test.png", slug: "renewable-energy-advances" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Cybersecurity Measures Enhanced Across Sectors",
            excerpt: "New protocols implemented to protect critical infrastructure.",
            slug: "cybersecurity-measures-enhanced",
          },
          {
            title: "Biotechnology Research Yields Promising Results",
            excerpt: "Innovative techniques offer potential for disease prevention.",
            slug: "biotechnology-research-results",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "cybersecurity-featured",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Tech Giants Announce Collaboration",
          img: "/test.png",
          slug: "tech-giants-collaboration",
        },
        thumbArticles: [
          { title: "Smartphone Technology Evolves", img: "/test.png", slug: "smartphone-technology-evolves", date: "Senin, 9 Des 2025" },
          { title: "Climate Modeling Improves with AI", img: "/test.png", slug: "climate-modeling-ai", date: "Senin, 9 Des 2025" },
          { title: "Electric Vehicle Sales Surge", img: "/test.png", slug: "electric-vehicle-sales", date: "Minggu, 8 Des 2025" },
          { title: "Data Privacy Regulations Updated", img: "/test.png", slug: "data-privacy-regulations", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="SCI-TECH" contentBlocks={contentBlocks} />;
}