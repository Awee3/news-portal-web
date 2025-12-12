import CategoryPageTemplate from "@/components/CategoryPageTemplate";

export default function OpinionPage() {
  const contentBlocks = [
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "Editorial: Future of Democracy in Digital Age",
            excerpt: "Technology transforms civic engagement but raises important questions.",
            href: "/opinion/democracy",
          },
          {
            title: "Columnist: Economic Policy Needs Fresh Perspective",
            excerpt: "Traditional approaches may not address contemporary challenges.",
            href: "/opinion/economy",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/opinion/featured-1",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Letter to Editor: Education Reform Concerns", img: "/test.png", href: "/opinion/education" },
        { title: "Op-Ed: Climate Action Requires Urgency", img: "/test.png", href: "/opinion/climate" },
        { title: "Commentary: Social Media Influence", img: "/test.png", href: "/opinion/social-media" },
        { title: "Analysis: Healthcare System Sustainability", img: "/test.png", href: "/opinion/healthcare" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Guest Column: Innovation and Ethics",
            excerpt: "Technological advancement must consider moral implications.",
            href: "/opinion/ethics",
          },
          {
            title: "Perspective: Cultural Preservation Matters",
            excerpt: "Balancing modernization with traditional values remains crucial.",
            href: "/opinion/culture",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/opinion/featured-2",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Debate: Infrastructure Investment Priorities",
          img: "/test.png",
          href: "/opinion/infrastructure",
        },
        thumbArticles: [
          { title: "Reader Response: Local Governance", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Editorial Board: Press Freedom", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Viewpoint: Urban Planning Challenges", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
          { title: "Commentary: Youth Engagement", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="OPINION" contentBlocks={contentBlocks} />;
}