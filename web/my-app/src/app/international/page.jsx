import CategoryPageTemplate from "@/components/CategoryPageTemplate";

export const metadata = { title: "International - Bintaro Times" };

export default function InternationalPage() {
  // Data dummy - nantinya fetch dari API: /api/v1/categories/international/articles
  const contentBlocks = [
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "Global Summit Discusses Climate Change Initiatives",
            excerpt: "World leaders convene to address pressing environmental concerns and establish new international protocols.",
            href: "/international/climate-summit",
          },
          {
            title: "Economic Partnership Agreement Signed Between Nations",
            excerpt: "Historic trade agreement promises to boost economic cooperation and regional development.",
            href: "/international/trade-agreement",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/international/featured-1",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Diplomatic Relations Strengthened", img: "/test.png", href: "/international/diplomacy" },
        { title: "International Aid Reaches Remote Areas", img: "/test.png", href: "/international/aid" },
        { title: "Global Security Conference Concludes", img: "/test.png", href: "/international/security" },
        { title: "Cross-Border Infrastructure Project Launched", img: "/test.png", href: "/international/infrastructure" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "United Nations Passes New Resolution",
            excerpt: "Member states vote unanimously on critical international policy affecting multiple regions.",
            href: "/international/un-resolution",
          },
          {
            title: "Cultural Exchange Program Expands",
            excerpt: "Educational initiatives bring students together from diverse backgrounds worldwide.",
            href: "/international/cultural-exchange",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/international/featured-2",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Breaking: International Peace Talks Resume",
          img: "/test.png",
          href: "/international/peace-talks",
        },
        thumbArticles: [
          { title: "Regional Trade Statistics Released", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Migration Patterns Shift Globally", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "International Sports Event Announced", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
          { title: "Global Health Initiative Launched", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="INTERNATIONAL" contentBlocks={contentBlocks} />;
}