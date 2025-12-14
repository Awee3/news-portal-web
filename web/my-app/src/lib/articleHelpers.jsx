const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
const API_BASE = API_URL.replace("/api/v1", "");

const CATEGORY_MAP = {
  'scitech': 'Sci-Tech',
  'sci-tech': 'Sci-Tech',
  'health': 'Health',
  'international': 'International',
  'national': 'National',
  'politics': 'Politics',
  'business': 'Business',
  'sports': 'Sports',
  'lifestyle': 'Lifestyle',
  'opinion': 'Opinion',
};

export async function getArticlesByCategory(categoryName) {
  try {
    const res = await fetch(`${API_URL}/articles?status=published&kategori=${categoryName}&limit=50`, {
      next: { revalidate: 60 },
    });

    if (!res.ok) {
      return [];
    }

    const data = await res.json();
    return Array.isArray(data) ? data : [];
  } catch (error) {
    console.error("Error fetching articles:", error);
    return [];
  }
}

export async function getAllArticles() {
  try {
    const res = await fetch(`${API_URL}/articles?status=published&limit=50`, {
      cache: 'no-store',
    });

    if (!res.ok) {
      console.error("Failed to fetch articles:", res.status);
      return [];
    }

    const data = await res.json();
    return Array.isArray(data) ? data : [];
  } catch (error) {
    console.error("Error fetching articles:", error);
    return [];
  }
}

export function getImageUrl(article) {
  if (article?.gambar_utama) {
    if (article.gambar_utama.startsWith("http")) {
      return article.gambar_utama;
    }
    return `${API_BASE}${article.gambar_utama}`;
  }
  return "/test.png";
}

export function formatDate(dateString) {
  if (!dateString) return "";
  const date = new Date(dateString);
  return date.toLocaleDateString("id-ID", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
}

export function transformToContentBlocks(articles) {
  if (articles.length === 0) {
    return [];
  }

  let workingArticles = [...articles];
  while (workingArticles.length < 13) {
    workingArticles = [...workingArticles, ...articles];
  }

  const shuffled = workingArticles.sort(() => Math.random() - 0.5);
  const blocksData = [];
  const articlesPerBlock = 13;
  const maxBlocks = Math.min(3, Math.floor(shuffled.length / articlesPerBlock));

  for (let blockIndex = 0; blockIndex < maxBlocks; blockIndex++) {
    const startIdx = blockIndex * articlesPerBlock;
    const blockArticles = shuffled.slice(startIdx, startIdx + articlesPerBlock);

    const block = {
      id: blockIndex + 1,
      topMainSection: {
        sideTitles: [
          {
            title: blockArticles[0]?.judul || "Lorem Ipsum",
            excerpt: blockArticles[0]?.ringkasan || blockArticles[0]?.konten?.substring(0, 150) || "Lorem ipsum...",
            slug: blockArticles[0]?.slug || "#",
          },
          {
            title: blockArticles[1]?.judul || "Lorem Ipsum",
            excerpt: blockArticles[1]?.ringkasan || blockArticles[1]?.konten?.substring(0, 150) || "Lorem ipsum...",
            slug: blockArticles[1]?.slug || "#",
          },
        ],
        featured: {
          img: getImageUrl(blockArticles[2]),
          slug: blockArticles[2]?.slug || "#",
          priority: blockIndex === 0,
        },
      },
      carouselItems: [
        {
          title: blockArticles[3]?.judul || "Lorem Ipsum",
          img: getImageUrl(blockArticles[3]),
          slug: blockArticles[3]?.slug || "#",
        },
        {
          title: blockArticles[4]?.judul || "Lorem Ipsum",
          img: getImageUrl(blockArticles[4]),
          slug: blockArticles[4]?.slug || "#",
        },
        {
          title: blockArticles[5]?.judul || "Lorem Ipsum",
          img: getImageUrl(blockArticles[5]),
          slug: blockArticles[5]?.slug || "#",
        },
        {
          title: blockArticles[6]?.judul || "Lorem Ipsum",
          img: getImageUrl(blockArticles[6]),
          slug: blockArticles[6]?.slug || "#",
        },
      ],
      bottomMainSection: {
        sideTitles: [
          {
            title: blockArticles[7]?.judul || "Lorem Ipsum",
            excerpt: blockArticles[7]?.ringkasan || blockArticles[7]?.konten?.substring(0, 150) || "Lorem ipsum...",
            slug: blockArticles[7]?.slug || "#",
          },
          {
            title: blockArticles[8]?.judul || "Lorem Ipsum",
            excerpt: blockArticles[8]?.ringkasan || blockArticles[8]?.konten?.substring(0, 150) || "Lorem ipsum...",
            slug: blockArticles[8]?.slug || "#",
          },
        ],
        featured: {
          img: getImageUrl(blockArticles[9]),
          slug: blockArticles[9]?.slug || "#",
        },
      },
      sidebar: {
        mainArticle: {
          title: blockArticles[10]?.judul || "Lorem Ipsum",
          img: getImageUrl(blockArticles[10]),
          slug: blockArticles[10]?.slug || "#",
        },
        thumbArticles: [
          {
            title: blockArticles[11]?.judul || "Lorem Ipsum",
            img: getImageUrl(blockArticles[11]),
            slug: blockArticles[11]?.slug || "#",
            date: formatDate(blockArticles[11]?.tanggal_publikasi),
          },
          {
            title: blockArticles[12]?.judul || "Lorem Ipsum",
            img: getImageUrl(blockArticles[12]),
            slug: blockArticles[12]?.slug || "#",
            date: formatDate(blockArticles[12]?.tanggal_publikasi),
          },
          {
            title: blockArticles[0]?.judul || "Lorem Ipsum",
            img: getImageUrl(blockArticles[0]),
            slug: blockArticles[0]?.slug || "#",
            date: formatDate(blockArticles[0]?.tanggal_publikasi),
          },
          {
            title: blockArticles[1]?.judul || "Lorem Ipsum",
            img: getImageUrl(blockArticles[1]),
            slug: blockArticles[1]?.slug || "#",
            date: formatDate(blockArticles[1]?.tanggal_publikasi),
          },
        ],
      },
    };

    blocksData.push(block);
  }

  return blocksData;
}