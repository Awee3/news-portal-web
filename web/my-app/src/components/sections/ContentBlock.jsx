import MainSection from "./MainSection";
import CarouselStrip from "./CarouselStrip";
import Sidebar from "./Sidebar";

/**
 * Reusable content block: 2 MainSections + 1 CarouselStrip (left) + Sidebar (right)
 * This entire block can be repeated multiple times on the page
 */
export default function ContentBlock({ 
  topMainSection,
  carouselItems,
  bottomMainSection,
  sidebar 
}) {
  // Defensive check: ensure all required data exists
  if (!topMainSection || !bottomMainSection || !sidebar) {
    return null;
  }

  return (
    <section className="">
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        
        {/* LEFT COLUMN: MainSection + Carousel + MainSection */}
        <div className="lg:col-span-8 space-y-8 lg:border-r lg:border-gray-300 lg:pr-8">
          {/* 1. MainSection Atas */}
          <MainSection 
            sideTitles={topMainSection.sideTitles || []} 
            featured={topMainSection.featured || {}} 
          />

          {/* 2. CarouselStrip */}
          {carouselItems && carouselItems.length > 0 && (
            <CarouselStrip items={carouselItems} />
          )}

          {/* 3. MainSection Bawah */}
          <MainSection 
            sideTitles={bottomMainSection.sideTitles || []} 
            featured={bottomMainSection.featured || {}} 
          />
        </div>

        {/* RIGHT COLUMN: Sidebar */}
        <div className="lg:col-span-4">
          <Sidebar 
            mainArticle={sidebar.mainArticle || {}} 
            thumbArticles={sidebar.thumbArticles || []} 
          />
        </div>
      </div>
    </section>
  );
}