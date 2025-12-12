import CategoryPageTemplate from "@/components/CategoryPageTemplate";

export default function SportsPage() {
  const contentBlocks = [
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "Championship Finals Deliver Thrilling Conclusion",
            excerpt: "Intense competition culminates in dramatic victory for underdog team.",
            slug: "championship-finals-thrilling",
          },
          {
            title: "Olympic Athletes Begin Training Camp",
            excerpt: "National team prepares for upcoming international competition.",
            slug: "olympic-athletes-training",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "sports-featured-article",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Transfer Window Sees Major Signings", img: "/test.png", slug: "transfer-window-signings" },
        { title: "Youth Sports Development Program Launched", img: "/test.png", slug: "youth-sports-program" },
        { title: "Stadium Renovation Project Completed", img: "/test.png", slug: "stadium-renovation-completed" },
        { title: "Sports Science Advances Training Methods", img: "/test.png", slug: "sports-science-training" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Record-Breaking Performance Stuns Audience",
            excerpt: "Athlete surpasses previous benchmarks in spectacular fashion.",
            slug: "record-breaking-performance",
          },
          {
            title: "Regional Tournament Schedule Announced",
            excerpt: "Teams prepare for highly anticipated competitive season.",
            slug: "regional-tournament-schedule",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "record-breaking-featured",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Coach Shares Strategy Insights",
          img: "/test.png",
          slug: "coach-strategy-insights",
        },
        thumbArticles: [
          { title: "Player Injury Update Released", img: "/test.png", slug: "player-injury-update", date: "Senin, 9 Des 2025" },
          { title: "Sports Equipment Technology Evolves", img: "/test.png", slug: "sports-equipment-tech", date: "Senin, 9 Des 2025" },
          { title: "Fan Engagement Initiatives Launched", img: "/test.png", slug: "fan-engagement-initiatives", date: "Minggu, 8 Des 2025" },
          { title: "Sports Broadcasting Rights Negotiated", img: "/test.png", slug: "broadcasting-rights", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="SPORTS" contentBlocks={contentBlocks} />;
}