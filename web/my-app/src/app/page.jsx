import Header from "../components/Header";
import MainSection from "../components/sections/MainSection";
import CarouselStrip from "../components/sections/CarouselStrip";
import Sidebar from "../components/sections/Sidebar";

export default function Home() {
  // Data constants (nanti fetch dari API)
  const topFeatured = {
    img: "/test.png",
    href: "/international/international-summit",
    priority: true, // For LCP optimization
  };

  const topSideTitles = [
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Semper Lacinia...",
      href: "/international/regional-security-talks",
    },
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Cras elementum libero ac...",
      href: "/business/market-update",
    },
  ];

  const stripItems = [
    { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "/politics/policy-brief" },
    { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "/health/health-update" },
    { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "/tech/innovation" },
    { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "/business/economy" },
  ];

  const bottomFeatured = {
    img: "/test.png",
    href: "/sports/featured-match",
  };

  const bottomSideTitles = [
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Semper Lacinia...",
      href: "/sports/match-analysis",
    },
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Cras elementum libero ac...",
      href: "/sports/player-interview",
    },
  ];

  const sidebarMain = {
    title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
    excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenean quis magna in urna...",
    img: "/test.png",
    href: "/opinion/editorial",
  };

  const sidebarThumbs = [
    { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "/national/update-1" },
    { title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.", img: "/test.png", href: "/national/update-2" },
  ];

  return (
    <div className="min-h-screen bg-white">
      <Header />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-12 gap-8 relative">
          {/* Left Column */}
          <div className="col-span-12 lg:col-span-8">
            <MainSection sideTitles={topSideTitles} featured={topFeatured} />
            <CarouselStrip items={stripItems} />
            <MainSection sideTitles={bottomSideTitles} featured={bottomFeatured} />
          </div>

          {/* Vertical divider */}
          <div className="hidden lg:block absolute left-[calc(66.666%+0.3rem)] top-0 bottom-0 w-px bg-gray-200" />

          {/* Right Sidebar */}
          <Sidebar mainArticle={sidebarMain} thumbArticles={sidebarThumbs} />
        </div>
      </main>
    </div>
  );
}
