import CategoryPageTemplate from "@/components/CategoryPageTemplate";

export default function HealthPage() {
  const contentBlocks = [
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "New Treatment Protocol Shows Promising Results",
            excerpt: "Clinical trials demonstrate effectiveness in managing chronic conditions.",
            slug: "treatment-protocol-promising",
          },
          {
            title: "Mental Health Awareness Campaign Launched",
            excerpt: "Initiative aims to reduce stigma and improve access to services.",
            slug: "mental-health-awareness",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "health-featured-article",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Vaccination Drive Reaches Rural Areas", img: "/test.png", slug: "vaccination-drive-rural" },
        { title: "Nutrition Guidelines Updated", img: "/test.png", slug: "nutrition-guidelines-updated" },
        { title: "Hospital Capacity Expands to Meet Demand", img: "/test.png", slug: "hospital-capacity-expands" },
        { title: "Telemedicine Services Improve Accessibility", img: "/test.png", slug: "telemedicine-accessibility" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Preventive Care Programs Save Lives",
            excerpt: "Early screening initiatives detect diseases before symptoms appear.",
            slug: "preventive-care-programs",
          },
          {
            title: "Medical Research Funding Increases",
            excerpt: "Investment in healthcare innovation promises future breakthroughs.",
            slug: "medical-research-funding",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "preventive-care-featured",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Health Insurance Coverage Expands",
          img: "/test.png",
          slug: "health-insurance-expands",
        },
        thumbArticles: [
          { title: "Sleep Quality Affects Overall Wellness", img: "/test.png", slug: "sleep-quality-wellness", date: "Senin, 9 Des 2025" },
          { title: "Exercise Benefits for All Ages", img: "/test.png", slug: "exercise-benefits-ages", date: "Senin, 9 Des 2025" },
          { title: "Pediatric Care Standards Improved", img: "/test.png", slug: "pediatric-care-standards", date: "Minggu, 8 Des 2025" },
          { title: "Chronic Disease Management Tips", img: "/test.png", slug: "chronic-disease-management", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="HEALTH" contentBlocks={contentBlocks} />;
}