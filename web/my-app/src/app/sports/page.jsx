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
            href: "/sports/championship",
          },
          {
            title: "Olympic Athletes Begin Training Camp",
            excerpt: "National team prepares for upcoming international competition.",
            href: "/sports/olympics",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/sports/featured-1",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Transfer Window Sees Major Signings", img: "/test.png", href: "/sports/transfer" },
        { title: "Youth Sports Development Program Launched", img: "/test.png", href: "/sports/youth" },
        { title: "Stadium Renovation Project Completed", img: "/test.png", href: "/sports/stadium" },
        { title: "Sports Science Advances Training Methods", img: "/test.png", href: "/sports/science" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Record-Breaking Performance Stuns Audience",
            excerpt: "Athlete surpasses previous benchmarks in spectacular fashion.",
            href: "/sports/record",
          },
          {
            title: "Regional Tournament Schedule Announced",
            excerpt: "Teams prepare for highly anticipated competitive season.",
            href: "/sports/tournament",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/sports/featured-2",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Coach Shares Strategy Insights",
          img: "/test.png",
          href: "/sports/coach",
        },
        thumbArticles: [
          { title: "Player Injury Update Released", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Sports Equipment Technology Evolves", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Fan Engagement Initiatives Launched", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
          { title: "Sports Broadcasting Rights Negotiated", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="SPORTS" contentBlocks={contentBlocks} />;
}