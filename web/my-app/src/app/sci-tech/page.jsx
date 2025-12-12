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
            href: "/sci-tech/quantum",
          },
          {
            title: "AI System Demonstrates Advanced Problem-Solving",
            excerpt: "Machine learning model surpasses expectations in complex cognitive tasks.",
            href: "/sci-tech/ai-advancement",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/sci-tech/featured-1",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Space Mission Discovers New Exoplanet", img: "/test.png", href: "/sci-tech/exoplanet" },
        { title: "5G Network Coverage Expands Nationwide", img: "/test.png", href: "/sci-tech/5g" },
        { title: "Medical Breakthrough in Cancer Treatment", img: "/test.png", href: "/sci-tech/cancer" },
        { title: "Renewable Energy Technology Advances", img: "/test.png", href: "/sci-tech/renewable" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Cybersecurity Measures Enhanced Across Sectors",
            excerpt: "New protocols implemented to protect critical infrastructure.",
            href: "/sci-tech/cybersecurity",
          },
          {
            title: "Biotechnology Research Yields Promising Results",
            excerpt: "Innovative techniques offer potential for disease prevention.",
            href: "/sci-tech/biotech",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/sci-tech/featured-2",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Tech Giants Announce Collaboration",
          img: "/test.png",
          href: "/sci-tech/collaboration",
        },
        thumbArticles: [
          { title: "Smartphone Technology Evolves", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Climate Modeling Improves with AI", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Electric Vehicle Sales Surge", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
          { title: "Data Privacy Regulations Updated", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="SCI-TECH" contentBlocks={contentBlocks} />;
}