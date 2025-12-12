import CategoryPageTemplate from "@/components/CategoryPageTemplate";

export default function LifestylePage() {
  const contentBlocks = [
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "Wellness Trends Transform Daily Routines",
            excerpt: "Health-conscious practices gain popularity among urban populations.",
            slug: "wellness-trends-daily-routines",
          },
          {
            title: "Sustainable Fashion Movement Gains Momentum",
            excerpt: "Eco-friendly clothing choices reflect changing consumer values.",
            slug: "sustainable-fashion-movement",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "lifestyle-featured-article",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Home Decor Trends for Modern Living", img: "/test.png", slug: "home-decor-trends" },
        { title: "Travel Destinations Off the Beaten Path", img: "/test.png", slug: "travel-destinations-offbeat" },
        { title: "Culinary Innovations in Local Restaurants", img: "/test.png", slug: "culinary-innovations-local" },
        { title: "Mindfulness Practices for Stress Relief", img: "/test.png", slug: "mindfulness-stress-relief" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Fitness Technology Revolutionizes Workouts",
            excerpt: "Wearable devices and apps enhance personal training experiences.",
            slug: "fitness-technology-workouts",
          },
          {
            title: "Cultural Events Celebrate Local Heritage",
            excerpt: "Festivals and exhibitions showcase traditional arts and crafts.",
            slug: "cultural-events-heritage",
          },
        ],
        featured: {
          img: "/test.png",
          slug: "fitness-tech-featured",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Plant-Based Diet Gains Followers",
          img: "/test.png",
          slug: "plant-based-diet-followers",
        },
        thumbArticles: [
          { title: "Beauty Product Innovations Unveiled", img: "/test.png", slug: "beauty-product-innovations", date: "Senin, 9 Des 2025" },
          { title: "Pet Care Tips for New Owners", img: "/test.png", slug: "pet-care-tips", date: "Senin, 9 Des 2025" },
          { title: "Work-Life Balance Strategies", img: "/test.png", slug: "work-life-balance", date: "Minggu, 8 Des 2025" },
          { title: "Urban Gardening Projects Flourish", img: "/test.png", slug: "urban-gardening-projects", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return <CategoryPageTemplate categoryName="LIFESTYLE" contentBlocks={contentBlocks} />;
}