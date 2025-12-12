import CategoryPageTemplate from "@/components/CategoryPageTemplate";

export default function BusinessPage() {
  const contentBlocks = [
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "Stock Market Reaches New All-Time High",
            excerpt: "Investor confidence drives unprecedented growth across major indices.",
            href: "/business/stock-market",
          },
          {
            title: "Tech Startup Secures Record Funding Round",
            excerpt: "Innovative company attracts major investment for expansion plans.",
            href: "/business/startup-funding",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/business/featured-1",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Corporate Earnings Exceed Projections", img: "/test.png", href: "/business/earnings" },
        { title: "Merger Deal Reshapes Industry Landscape", img: "/test.png", href: "/business/merger" },
        { title: "Export Volume Increases Significantly", img: "/test.png", href: "/business/export" },
        { title: "E-Commerce Platform Expands Operations", img: "/test.png", href: "/business/ecommerce" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Central Bank Announces Interest Rate Decision",
            excerpt: "Monetary policy adjustment aims to balance growth and inflation.",
            href: "/business/interest-rate",
          },
          {
            title: "Small Business Support Program Launched",
            excerpt: "Government initiative provides resources for entrepreneurship development.",
            href: "/business/sme-program",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/business/featured-2",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Global Supply Chain Disruptions Ease",
          img: "/test.png",
          href: "/business/supply-chain",
        },
        thumbArticles: [
          { title: "Cryptocurrency Market Volatility Continues", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Real Estate Sector Shows Growth", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Manufacturing Output Increases", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
          { title: "Consumer Spending Trends Analyzed", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="BUSINESS" contentBlocks={contentBlocks} />;
}