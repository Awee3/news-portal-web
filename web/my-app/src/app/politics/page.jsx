import CategoryPageTemplate from "@/components/CategoryPageTemplate";

export default function PoliticsPage() {
  const contentBlocks = [
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "Legislative Session Addresses Key Policy Reforms",
            excerpt: "Lawmakers debate critical legislation affecting various sectors of governance.",
            href: "/politics/legislative-session",
          },
          {
            title: "Coalition Government Announces Strategic Partnership",
            excerpt: "Political parties unite to advance shared policy objectives.",
            href: "/politics/coalition",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/politics/featured-1",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Electoral Commission Updates Voter Registry", img: "/test.png", href: "/politics/electoral" },
        { title: "Political Campaign Finance Reforms Proposed", img: "/test.png", href: "/politics/finance" },
        { title: "Parliamentary Debate on Constitutional Amendment", img: "/test.png", href: "/politics/constitution" },
        { title: "Regional Autonomy Bill Gains Support", img: "/test.png", href: "/politics/autonomy" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Cabinet Reshuffle Expected Next Month",
            excerpt: "Sources indicate potential changes in key ministerial positions.",
            href: "/politics/cabinet",
          },
          {
            title: "Anti-Corruption Agency Reports Progress",
            excerpt: "Annual review highlights enforcement efforts and ongoing investigations.",
            href: "/politics/corruption",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/politics/featured-2",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Opinion Poll Shows Shifting Public Sentiment",
          img: "/test.png",
          href: "/politics/poll",
        },
        thumbArticles: [
          { title: "Local Elections Scheduled for Next Quarter", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Political Party Conventions Begin", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Government Response to Opposition Demands", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
          { title: "Diplomatic Mission Returns from Abroad", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="POLITICS" contentBlocks={contentBlocks} />;
}