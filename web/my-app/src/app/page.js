import Image from "next/image";
import Link from "next/link";

export default function Home() {
  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
      <header className="border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Top bar with date and user */}
          <div className="flex justify-between items-center py-2 text-sm">
            <div className="text-gray-600">
              SENIN, 29 JULI 2025
            </div>
            <div className="flex items-center gap-4">
              <div className="relative">
                <input
                  type="text"
                  placeholder="Search"
                  className="border border-gray-300 rounded px-3 py-1 text-sm w-48"
                />
                <button className="absolute right-2 top-1/2 transform -translate-y-1/2">
                  🔍
                </button>
              </div>
              <div className="text-gray-600">USERNAME</div>
            </div>
          </div>

          {/* Logo and tagline */}
          <div className="text-center py-8">
            <h1 className="text-4xl font-bold tracking-widest mb-2">
              BINTARO TIMES
            </h1>
            <p className="text-gray-600 text-sm uppercase tracking-wide">
              LOREM IPSUM DOLOR SIT AMET, CONSECTETUR ADIPISCING ELIT, DUB SOLLICITUDIN.
            </p>
          </div>

          {/* Navigation */}
          <nav className="border-t border-b border-gray-200">
            <div className="flex justify-center space-x-8 py-3">
              <Link href="/international" className="text-sm font-medium hover:text-gray-600">
                INTERNATIONAL
              </Link>
              <Link href="/national" className="text-sm font-medium hover:text-gray-600">
                NATIONAL
              </Link>
              <Link href="/politics" className="text-sm font-medium hover:text-gray-600">
                POLITICS
              </Link>
              <Link href="/business" className="text-sm font-medium hover:text-gray-600">
                BUSINESS
              </Link>
              <Link href="/sci-tech" className="text-sm font-medium hover:text-gray-600">
                SCI-TECH
              </Link>
              <Link href="/lifestyle" className="text-sm font-medium hover:text-gray-600">
                LIFESTYLE
              </Link>
              <Link href="/health" className="text-sm font-medium hover:text-gray-600">
                HEALTH
              </Link>
              <Link href="/sports" className="text-sm font-medium hover:text-gray-600">
                SPORTS
              </Link>
              <Link href="/opinion" className="text-sm font-medium hover:text-gray-600">
                OPINION
              </Link>
            </div>
          </nav>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-12 gap-8">
          {/* Left Column - Main Articles */}
          <div className="col-span-8">
            {/* Featured Article */}
            <article className="mb-8">
              <h2 className="text-2xl font-bold mb-4">
                Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.
              </h2>
              <div className="relative mb-4">
                <Image
                  src="/placeholder-news.jpg"
                  alt="Featured article"
                  width={600}
                  height={400}
                  className="w-full h-64 object-cover"
                />
              </div>
              <p className="text-gray-600 text-sm mb-4">
                Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Tempor Lacinia. Ut Elementum...
              </p>
            </article>

            {/* Secondary Articles */}
            <div className="grid grid-cols-2 gap-6 mb-8">
              <article>
                <h3 className="text-lg font-bold mb-2">
                  Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.
                </h3>
                <div className="relative mb-2">
                  <Image
                    src="/placeholder-news-2.jpg"
                    alt="Article"
                    width={280}
                    height={180}
                    className="w-full h-32 object-cover"
                  />
                </div>
                <p className="text-gray-600 text-xs">
                  Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Tempor Lacinia...
                </p>
              </article>

              <article>
                <h3 className="text-lg font-bold mb-2">
                  Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.
                </h3>
                <div className="relative mb-2">
                  <Image
                    src="/placeholder-news-3.jpg"
                    alt="Article"
                    width={280}
                    height={180}
                    className="w-full h-32 object-cover"
                  />
                </div>
                <p className="text-gray-600 text-xs">
                  Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Tempor Lacinia...
                </p>
              </article>
            </div>

            {/* Additional Articles Grid */}
            <div className="space-y-6">
              <article className="flex gap-4">
                <div className="w-24 h-16 flex-shrink-0">
                  <Image
                    src="/placeholder-news-4.jpg"
                    alt="Article thumbnail"
                    width={96}
                    height={64}
                    className="w-full h-full object-cover"
                  />
                </div>
                <div>
                  <h4 className="font-bold text-sm mb-1">
                    Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.
                  </h4>
                  <p className="text-gray-600 text-xs">
                    Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit...
                  </p>
                </div>
              </article>

              <article>
                <h3 className="text-lg font-bold mb-2">
                  Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.
                </h3>
                <div className="relative mb-2">
                  <Image
                    src="/placeholder-news-5.jpg"
                    alt="Article"
                    width={600}
                    height={300}
                    className="w-full h-48 object-cover"
                  />
                </div>
                <p className="text-gray-600 text-sm">
                  Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Tempor Lacinia. Ut Elementum...
                </p>
              </article>
            </div>
          </div>

          {/* Right Sidebar */}
          <div className="col-span-4">
            {/* Main sidebar article */}
            <article className="mb-6">
              <div className="relative mb-3">
                <Image
                  src="/placeholder-news-sidebar.jpg"
                  alt="Sidebar article"
                  width={300}
                  height={200}
                  className="w-full h-48 object-cover"
                />
              </div>
              <h3 className="font-bold mb-2">
                Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.
              </h3>
              <p className="text-gray-600 text-sm">
                Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna...
              </p>
            </article>

            {/* Sidebar thumbnails */}
            <div className="grid grid-cols-2 gap-3">
              <article>
                <div className="relative mb-2">
                  <Image
                    src="/placeholder-thumb-1.jpg"
                    alt="Thumbnail"
                    width={140}
                    height={100}
                    className="w-full h-20 object-cover"
                  />
                </div>
                <h4 className="font-bold text-xs mb-1">
                  Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.
                </h4>
              </article>

              <article>
                <div className="relative mb-2">
                  <Image
                    src="/placeholder-thumb-2.jpg"
                    alt="Thumbnail"
                    width={140}
                    height={100}
                    className="w-full h-20 object-cover"
                  />
                </div>
                <h4 className="font-bold text-xs mb-1">
                  Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.
                </h4>
              </article>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
