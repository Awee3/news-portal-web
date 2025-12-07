import Image from "next/image";
import Link from "next/link";
import CurrentDate from "../components/CurrentDate";

export default function Home() {
  // Data dummy sementara
  const featured = {
    title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
    excerpt:
      "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Tempor Lacinia. Ut Elementum...",
    img: "/test.png",
    href: "/international/international-summit",
  };

  const sideTitles = [
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      excerpt:
        "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Semper Lacinia...",
      href: "/international/regional-security-talks",
    },
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      excerpt:
        "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Cras elementum libero ac...",
      href: "/business/market-update",
    },
  ];

  const stripItems = [
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      img: "/test.png",
      href: "/politics/policy-brief",
    },
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      img: "/test.png",
      href: "/health/health-update",
    },
  ];

  const listItems = [
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit...",
      img: "/test.png",
      href: "/lifestyle/travel",
    },
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      excerpt: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit...",
      img: "/test.png",
      href: "/sci-tech/ai",
    },
  ];

  const bigBelow = {
    title:
      "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenam Quis Magna In Urna Tempor Lacinia.",
    img: "/test.png",
    href: "/sports/match",
  };

  const sidebarMain = {
    title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
    excerpt:
      "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit. Aenean quis magna in urna...",
    img: "/test.png",
    href: "/opinion/editorial",
  };

  const sidebarThumbs = [
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      img: "/test.png",
      href: "/national/update-1",
    },
    {
      title: "Lorem Ipsum Dolor Sit Amet, Consectetur Adipiscing Elit.",
      img: "/test.png",
      href: "/national/update-2",
    },
  ];

  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
      <header className="border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Top utility bar */}
          <div className="flex justify-between items-center py-2 text-sm">
            <div className="flex items-center gap-3">
              <div className="relative">
                <input
                  type="text"
                  placeholder="Search"
                  className="border border-gray-300 rounded px-3 py-1 text-sm w-48"
                />
              </div>
            </div>
            <div className="text-gray-600">USERNAME</div>
          </div>

          {/* Date */}
          <div className="text-left py-1 text-sm">
            <CurrentDate />
          </div>

          {/* Logo + tagline */}
          <div className="text-center py-8">
            <h1 className="text-4xl font-bold tracking-widest mb-2">
              BINTARO TIMES
            </h1>
            <p className="text-gray-600 text-sm uppercase tracking-wide">
              LOREM IPSUM DOLOR SIT AMET, CONSECTETUR ADIPISCING ELIT, DUIS SOLLICITUDIN.
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

      {/* Main */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-12 gap-8">
          {/* Left Column */}
          <div className="col-span-12 lg:col-span-8">
            {/* Featured split: text list + big image */}
            <section className="grid grid-cols-1 md:grid-cols-12 gap-6 mb-8">
              <div className="md:col-span-5 space-y-6">
                {sideTitles.map((a, i) => (
                  <article key={i}>
                    <Link href={a.href}>
                      <h2 className="text-xl font-bold leading-snug hover:underline">
                        {a.title}
                      </h2>
                    </Link>
                    <p className="text-gray-600 text-sm mt-2">{a.excerpt}</p>
                  </article>
                ))}
              </div>
              <div className="md:col-span-7">
                <Link href={featured.href}>
                  <div className="relative w-full h-64 md:h-[22rem]">
                    <Image
                      src={featured.img}
                      alt={featured.title}
                      fill
                      className="object-cover"
                      sizes="(max-width:768px) 100vw, 60vw"
                    />
                  </div>
                </Link>
                <Link href={featured.href}>
                  <h2 className="text-2xl font-bold mt-4 hover:underline">
                    {featured.title}
                  </h2>
                </Link>
                <p className="text-gray-600 text-sm mt-2">{featured.excerpt}</p>
              </div>
            </section>

            {/* Two wide thumbs strip */}
            <section className="mb-8">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 items-stretch">
                {stripItems.map((a, i) => (
                  <Link key={i} href={a.href} className="group">
                    <div className="flex items-center gap-4 border border-gray-200 p-3">
                      <div className="relative w-24 h-16 flex-shrink-0 bg-gray-100">
                        <Image
                          src={a.img}
                          alt={a.title}
                          fill
                          className="object-cover"
                          sizes="96px"
                        />
                      </div>
                      <h3 className="font-semibold text-sm leading-snug group-hover:underline">
                        {a.title}
                      </h3>
                    </div>
                  </Link>
                ))}
              </div>
            </section>

            {/* List rows with small thumbnail */}
            <section className="space-y-6 mb-8">
              {listItems.map((a, i) => (
                <article key={i} className="flex gap-4">
                  <Link href={a.href} className="relative w-24 h-16 flex-shrink-0 bg-gray-100">
                    <Image
                      src={a.img}
                      alt={a.title}
                      fill
                      className="object-cover"
                      sizes="96px"
                    />
                  </Link>
                  <div>
                    <Link href={a.href}>
                      <h4 className="font-bold text-sm mb-1 hover:underline">{a.title}</h4>
                    </Link>
                    <p className="text-gray-600 text-xs">{a.excerpt}</p>
                  </div>
                </article>
              ))}
            </section>

            {/* Big article below */}
            <section>
              <Link href={bigBelow.href}>
                <h3 className="text-lg md:text-xl font-bold mb-2 hover:underline">
                  {bigBelow.title}
                </h3>
              </Link>
              <div className="relative w-full h-56 md:h-72">
                <Image
                  src={bigBelow.img}
                  alt={bigBelow.title}
                  fill
                  className="object-cover"
                  sizes="(max-width:768px) 100vw, 60vw"
                />
              </div>
            </section>
          </div>

          {/* Right Sidebar */}
          <aside className="col-span-12 lg:col-span-4">
            {/* Main sidebar article */}
            <article className="mb-6">
              <Link href={sidebarMain.href}>
                <div className="relative w-full h-56 mb-3">
                  <Image
                    src={sidebarMain.img}
                    alt={sidebarMain.title}
                    fill
                    className="object-cover"
                    sizes="(max-width:768px) 100vw, 33vw"
                  />
                </div>
              </Link>
              <Link href={sidebarMain.href}>
                <h3 className="font-bold mb-2 hover:underline">{sidebarMain.title}</h3>
              </Link>
              <p className="text-gray-600 text-sm">{sidebarMain.excerpt}</p>
            </article>

            {/* Sidebar thumbs grid */}
            <div className="grid grid-cols-2 gap-3">
              {sidebarThumbs.map((a, i) => (
                <article key={i}>
                  <Link href={a.href}>
                    <div className="relative w-full h-20 mb-2">
                      <Image
                        src={a.img}
                        alt={a.title}
                        fill
                        className="object-cover"
                        sizes="150px"
                      />
                    </div>
                  </Link>
                  <Link href={a.href}>
                    <h4 className="font-bold text-xs mb-1 hover:underline">{a.title}</h4>
                  </Link>
                </article>
              ))}
            </div>

            {/* Dots / pager mock */}
            <div className="flex items-center gap-2 mt-4">
              <button
                className="w-6 h-6 text-xs rounded-full border border-gray-300 hover:bg-gray-100"
                aria-label="Prev"
              >
                ‹
              </button>
              <div className="flex gap-1">
                <span className="w-2 h-2 rounded-full bg-gray-400 inline-block" />
                <span className="w-2 h-2 rounded-full bg-gray-300 inline-block" />
                <span className="w-2 h-2 rounded-full bg-gray-300 inline-block" />
              </div>
              <button
                className="w-6 h-6 text-xs rounded-full border border-gray-300 hover:bg-gray-100"
                aria-label="Next"
              >
                ›
              </button>
            </div>
          </aside>
        </div>
      </main>
    </div>
  );
}
