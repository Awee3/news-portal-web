import Header from "./Header";
import ContentBlock from "./sections/ContentBlock";

/**
 * Reusable template untuk halaman kategori
 * @param {string} categoryName - Nama kategori (e.g., "INTERNATIONAL", "POLITICS")
 * @param {Array} contentBlocks - Array of content blocks (bisa 1-3 blok)
 */
export default function CategoryPageTemplate({ categoryName, contentBlocks = [] }) {
  return (
    <div className="min-h-screen bg-white">
      <Header />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Category Header */}
        <div className="mb-8 pb-4 border-b-2 border-black">
          <h1 className="text-3xl font-bold uppercase tracking-wider">
            {categoryName}
          </h1>
        </div>

        {/* Content Blocks */}
        {contentBlocks.map((block, index) => (
          <ContentBlock
            key={block.id || index}
            topMainSection={block.topMainSection}
            carouselItems={block.carouselItems}
            bottomMainSection={block.bottomMainSection}
            sidebar={block.sidebar}
          />
        ))}

        {/* Empty State */}
        {contentBlocks.length === 0 && (
          <div className="text-center py-16 text-gray-500">
            <p className="text-lg">Belum ada artikel dalam kategori ini.</p>
          </div>
        )}
      </main>
    </div>
  );
}