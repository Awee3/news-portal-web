import CategoryPageTemplate from "@/components/CategoryPageTemplate";

export default function InternationalPage() {
  const contentBlocks = [
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "Global Summit Discusses Climate Change Initiatives",
            excerpt: "World leaders convene to address pressing environmental concerns and establish new international protocols.",
            slug: "global-summit-climate-change-initiatives",
          },
          {
            title: "Economic Partnership Agreement Signed Between Nations",
            excerpt: "Historic trade agreement promises to boost economic cooperation and regional development.",
            slug: "economic-partnership-agreement-signed",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "international-featured-article",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Diplomatic Relations Strengthened", img: "/test.png", slug: "diplomatic-relations-strengthened" },
        { title: "International Aid Reaches Remote Areas", img: "/test.png", slug: "international-aid-remote-areas" },
        { title: "Global Security Conference Concludes", img: "/test.png", slug: "global-security-conference" },
        { title: "Cross-Border Infrastructure Project Launched", img: "/test.png", slug: "cross-border-infrastructure-project" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "United Nations Passes New Resolution",
            excerpt: "Member states vote unanimously on critical international policy affecting multiple regions.",
            slug: "united-nations-new-resolution",
          },
          {
            title: "Cultural Exchange Program Expands",
            excerpt: "Educational initiatives bring students together from diverse backgrounds worldwide.",
            slug: "cultural-exchange-program-expands",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "un-resolution-featured",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Breaking: International Peace Talks Resume",
          img: "/test.png",
          slug: "international-peace-talks-resume",
        },
        thumbArticles: [
          { title: "Regional Trade Statistics Released", img: "/test.png", slug: "regional-trade-statistics", date: "Senin, 9 Des 2025" },
          { title: "Migration Patterns Shift Globally", img: "/test.png", slug: "migration-patterns-shift", date: "Senin, 9 Des 2025" },
          { title: "International Sports Event Announced", img: "/test.png", slug: "international-sports-event", date: "Minggu, 8 Des 2025" },
          { title: "Global Health Initiative Launched", img: "/test.png", slug: "global-health-initiative", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="INTERNATIONAL" contentBlocks={contentBlocks} />;
}