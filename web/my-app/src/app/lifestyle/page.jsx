import CategoryPageTemplate from "@/components/CategoryPageTemplate";

export default function National() {
  const contentBlocks = [
    {
      id: 1,
      topMainSection: {
        sideTitles: [
          {
            title: "Wellness Trends Transform Daily Routines",
            excerpt: "Health-conscious practices gain popularity among urban populations.",
            href: "/lifestyle/wellness",
          },
          {
            title: "Sustainable Fashion Movement Gains Momentum",
            excerpt: "Eco-friendly clothing choices reflect changing consumer values.",
            href: "/lifestyle/fashion",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/lifestyle/featured-1",
          priority: true,
        }
      },
      carouselItems: [
        { title: "Home Decor Trends for Modern Living", img: "/test.png", href: "/lifestyle/home-decor" },
        { title: "Travel Destinations Off the Beaten Path", img: "/test.png", href: "/lifestyle/travel" },
        { title: "Culinary Innovations in Local Restaurants", img: "/test.png", href: "/lifestyle/culinary" },
        { title: "Mindfulness Practices for Stress Relief", img: "/test.png", href: "/lifestyle/mindfulness" },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: "Fitness Technology Revolutionizes Workouts",
            excerpt: "Wearable devices and apps enhance personal training experiences.",
            href: "/lifestyle/fitness-tech",
          },
          {
            title: "Cultural Events Celebrate Local Heritage",
            excerpt: "Festivals and exhibitions showcase traditional arts and crafts.",
            href: "/lifestyle/culture",
          },
        ],
        featured: {
          img: "/test.png",
          href: "/lifestyle/featured-2",
        }
      },
      sidebar: {
        mainArticle: {
          title: "Plant-Based Diet Gains Followers",
          img: "/test.png",
          href: "/lifestyle/plant-based",
        },
        thumbArticles: [
          { title: "Beauty Product Innovations Unveiled", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Pet Care Tips for New Owners", img: "/test.png", href: "#", date: "Senin, 9 Des 2025" },
          { title: "Work-Life Balance Strategies", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
          { title: "Urban Gardening Projects Flourish", img: "/test.png", href: "#", date: "Minggu, 8 Des 2025" },
        ]
      }
    }
  ];

  return (
    <div className="min-h-screen bg-white">
      <CategoryPageTemplate categoryName="LIFESTYLE" contentBlocks={contentBlocks} />
    </div>
  );
}