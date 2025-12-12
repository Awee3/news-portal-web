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
            slug: "stock-market-all-time-high",
          },
          {
            title: "Tech Startup Secures Record Funding Round",
            excerpt: "Innovative company attracts major investment for expansion plans.",
            slug: "tech-startup-record-funding",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "business-featured-article",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Corporate Earnings Exceed Projections", img: "/test.png", slug: "corporate-earnings-exceed" },
        { title: "Merger Deal Reshapes Industry Landscape", img: "/test.png", slug: "merger-deal-reshapes-industry" },
        { title: "Export Volume Increases Significantly", img: "/test.png", slug: "export-volume-increases" },
        { title: "E-Commerce Platform Expands Operations", img: "/test.png", slug: "ecommerce-platform-expands" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Central Bank Announces Interest Rate Decision",
            excerpt: "Monetary policy adjustment aims to balance growth and inflation.",
            slug: "central-bank-interest-rate",
          },
          {
            title: "Small Business Support Program Launched",
            excerpt: "Government initiative provides resources for entrepreneurship development.",
            slug: "small-business-support-program",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "interest-rate-featured",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Global Supply Chain Disruptions Ease",
          img: "/test.png",
          slug: "supply-chain-disruptions-ease",
        },
        thumbArticles: [
          { title: "Cryptocurrency Market Volatility Continues", img: "/test.png", slug: "cryptocurrency-volatility", date: "Senin, 9 Des 2025" },
          { title: "Real Estate Sector Shows Growth", img: "/test.png", slug: "real-estate-growth", date: "Senin, 9 Des 2025" },
          { title: "Manufacturing Output Increases", img: "/test.png", slug: "manufacturing-output-increases", date: "Minggu, 8 Des 2025" },
          { title: "Consumer Spending Trends Analyzed", img: "/test.png", slug: "consumer-spending-trends", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="BUSINESS" contentBlocks={contentBlocks} />;
}