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
            slug: "legislative-session-policy-reforms",
          },
          {
            title: "Coalition Government Announces Strategic Partnership",
            excerpt: "Political parties unite to advance shared policy objectives.",
            slug: "coalition-government-partnership",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "politics-featured-article",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Electoral Commission Updates Voter Registry", img: "/test.png", slug: "electoral-commission-voter-registry" },
        { title: "Political Campaign Finance Reforms Proposed", img: "/test.png", slug: "campaign-finance-reforms" },
        { title: "Parliamentary Debate on Constitutional Amendment", img: "/test.png", slug: "constitutional-amendment-debate" },
        { title: "Regional Autonomy Bill Gains Support", img: "/test.png", slug: "regional-autonomy-bill" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Cabinet Reshuffle Expected Next Month",
            excerpt: "Sources indicate potential changes in key ministerial positions.",
            slug: "cabinet-reshuffle-expected",
          },
          {
            title: "Anti-Corruption Agency Reports Progress",
            excerpt: "Annual review highlights enforcement efforts and ongoing investigations.",
            slug: "anti-corruption-agency-progress",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "cabinet-reshuffle-featured",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Opinion Poll Shows Shifting Public Sentiment",
          img: "/test.png",
          slug: "opinion-poll-public-sentiment",
        },
        thumbArticles: [
          { title: "Local Elections Scheduled for Next Quarter", img: "/test.png", slug: "local-elections-scheduled", date: "Senin, 9 Des 2025" },
          { title: "Political Party Conventions Begin", img: "/test.png", slug: "party-conventions-begin", date: "Senin, 9 Des 2025" },
          { title: "Government Response to Opposition Demands", img: "/test.png", slug: "government-opposition-response", date: "Minggu, 8 Des 2025" },
          { title: "Diplomatic Mission Returns from Abroad", img: "/test.png", slug: "diplomatic-mission-returns", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="POLITICS" contentBlocks={contentBlocks} />;
}