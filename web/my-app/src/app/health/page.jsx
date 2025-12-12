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
            href: "/health/treatment",
          },
          {
            title: "Mental Health Awareness Campaign Launched",
            excerpt: "Initiative aims to reduce stigma and improve access to services.",
            href: "/health/mental-health",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/health/featured-1",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Vaccination Drive Reaches Rural Areas", img: "/test.png", href: "/health/vaccination" },
        { title: "Nutrition Guidelines Updated", img: "/test.png", href: "/health/nutrition" },
        { title: "Hospital Capacity Expands to Meet Demand", img: "/test.png", href: "/health/hospital" },
        { title: "Telemedicine Services Improve Accessibility", img: "/test.png", href: "/health/telemedicine" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Preventive Care Programs Save Lives",
            excerpt: "Early screening initiatives detect diseases before symptoms appear.",
            href: "/health/preventive",
          },
          {
            title: "Medical Research Funding Increases",
            excerpt: "Investment in healthcare innovation promises future breakthroughs.",
            href: "/health/research",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/health/featured-2",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Health Insurance Coverage Expands",
          img: "/test.png",
          href: "/health/insurance",
        },
        thumbArticles: [
          { title: "Sleep Quality Affects Overall Wellness", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Exercise Benefits for All Ages", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Pediatric Care Standards Improved", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
          { title: "Chronic Disease Management Tips", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="HEALTH" contentBlocks={contentBlocks} />;
}