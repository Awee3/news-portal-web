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
            slug: "future-democracy-digital-age",
          },
          {
            title: "Columnist: Economic Policy Needs Fresh Perspective",
            excerpt: "Traditional approaches may not address contemporary challenges.",
            slug: "economic-policy-fresh-perspective",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "opinion-featured-article",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Letter to Editor: Education Reform Concerns", img: "/test.png", slug: "education-reform-concerns" },
        { title: "Op-Ed: Climate Action Requires Urgency", img: "/test.png", slug: "climate-action-urgency" },
        { title: "Commentary: Social Media Influence", img: "/test.png", slug: "social-media-influence" },
        { title: "Analysis: Healthcare System Sustainability", img: "/test.png", slug: "healthcare-sustainability" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Guest Column: Innovation and Ethics",
            excerpt: "Technological advancement must consider moral implications.",
            slug: "innovation-and-ethics",
          },
          {
            title: "Perspective: Cultural Preservation Matters",
            excerpt: "Balancing modernization with traditional values remains crucial.",
            slug: "cultural-preservation-matters",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "innovation-ethics-featured",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Debate: Infrastructure Investment Priorities",
          img: "/test.png",
          slug: "infrastructure-investment-debate",
        },
        thumbArticles: [
          { title: "Reader Response: Local Governance", img: "/test.png", slug: "local-governance-response", date: "Senin, 9 Des 2025" },
          { title: "Editorial Board: Press Freedom", img: "/test.png", slug: "press-freedom-editorial", date: "Senin, 9 Des 2025" },
          { title: "Viewpoint: Urban Planning Challenges", img: "/test.png", slug: "urban-planning-challenges", date: "Minggu, 8 Des 2025" },
          { title: "Commentary: Youth Engagement", img: "/test.png", slug: "youth-engagement-commentary", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="OPINION" contentBlocks={contentBlocks} />;
}