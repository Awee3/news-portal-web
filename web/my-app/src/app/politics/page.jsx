import Header from "@/components/Header";
import ContentBlock from "@/components/sections/ContentBlock";
import { getArticlesByCategory, transformToContentBlocks } from "@/lib/articleHelpers";

export default async function PoliticsPage() {
  const articles = await getArticlesByCategory("Politics");
  const contentBlocks = transformToContentBlocks(articles);

  if (contentBlocks.length === 0) {
    return (
      <div className="min-h-screen bg-white">
        <Header />
        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="text-center py-16">
            <p className="text-gray-500">Belum ada artikel Politics yang dipublikasikan.</p>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white">
      <Header />
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {contentBlocks.map((block) => (
          <ContentBlock
            key={block.id}
            topMainSection={block.topMainSection}
            carouselItems={block.carouselItems}
            bottomMainSection={block.bottomMainSection}
            sidebar={block.sidebar}
          />
        ))}
      </main>
    </div>
  );
}