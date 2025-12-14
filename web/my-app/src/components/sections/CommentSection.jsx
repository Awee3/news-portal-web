'use client';

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

function formatCommentDate(dateString) {
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now - date;
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMs / 3600000);
  const diffDays = Math.floor(diffMs / 86400000);

  if (diffMins < 1) return 'Baru saja';
  if (diffMins < 60) return `${diffMins} menit yang lalu`;
  if (diffHours < 24) return `${diffHours} jam yang lalu`;
  if (diffDays < 7) return `${diffDays} hari yang lalu`;
  
  return date.toLocaleDateString('id-ID', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  });
}

export default function CommentSection({ articleId, initialComments = [] }) {
  const [comments, setComments] = useState(initialComments);
  const [newComment, setNewComment] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [user, setUser] = useState(null);
  const router = useRouter();

  useEffect(() => {
    // Ambil user dari localStorage
    const storedUser = localStorage.getItem("user");
    if (storedUser) {
      try {
        setUser(JSON.parse(storedUser));
      } catch (error) {
        console.error("Error parsing user data:", error);
      }
    }

    // Load komentar dari API
    fetchComments();
  }, [articleId]);

  const fetchComments = async () => {
    try {
      const response = await fetch(`${API_URL}/articles/${articleId}/comments`);
      if (!response.ok) throw new Error('Failed to load comments');
      const data = await response.json();
      setComments(data);
    } catch (error) {
      console.error('Error loading comments:', error);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!newComment.trim()) {
      alert('Komentar harus diisi');
      return;
    }

    if (!user) {
      router.push('/login');
      return;
    }

    setIsSubmitting(true);

    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${API_URL}/articles/${articleId}/comments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          konten: newComment
        })
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Failed to post comment');
      }

      // Normalisasi response dan tambahkan ke state lokal agar user melihat komentarnya segera
      const normalized = {
        komentar_id: data.komentar_id ?? data.id ?? Date.now(),
        username: data.username ?? data.nama_pengguna ?? user.username,
        konten: data.konten ?? newComment,
        tanggal_dibuat: data.tanggal_dibuat ?? new Date().toISOString(),
        status: data.status ?? 'pending',
      };

      setComments((prev) => [normalized, ...prev]);
      setNewComment('');

      if (normalized.status !== 'approved') {
        // beri tahu user kalau comment menunggu moderasi
        alert('Komentar terkirim dan menunggu moderasi.');
      }

    } catch (error) {
      console.error('Error posting comment:', error);
      alert('Gagal mengirim komentar: ' + error.message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <section className="mt-12 pt-12 border-t border-gray-200">
      <h2 className="text-3xl font-bold font-serif mb-8">
        Komentar ({comments.length})
      </h2>

      {/* Comment Form */}
      {user ? (
        <form onSubmit={handleSubmit} className="mb-8 p-6 bg-gray-50 rounded-lg">
          <h3 className="text-lg font-semibold mb-4">Tinggalkan Komentar</h3>
          
          <div className="mb-4">
            <label htmlFor="konten" className="block text-sm font-medium text-gray-700 mb-2">
              Komentar sebagai {user.username}
            </label>
            <textarea
              id="konten"
              value={newComment}
              onChange={(e) => setNewComment(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
              rows="4"
              placeholder="Tulis komentar Anda..."
              required
            />
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="px-6 py-2 bg-black text-white font-semibold rounded-lg hover:bg-gray-800 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            {isSubmitting ? 'Mengirim...' : 'Kirim Komentar'}
          </button>
        </form>
      ) : (
        <div className="mb-8 p-6 bg-gray-50 rounded-lg text-center">
          <p className="text-gray-700 mb-4">Login untuk meninggalkan komentar</p>
          <button
            onClick={() => router.push('/login')}
            className="px-6 py-2 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition-colors"
          >
            Login
          </button>
        </div>
      )}

      {/* Comments List */}
      <div className="space-y-6">
        {comments.length === 0 ? (
          <p className="text-gray-500 text-center py-8">
            Belum ada komentar. Jadilah yang pertama berkomentar!
          </p>
        ) : (
          comments.map((comment) => (
            <div key={comment.komentar_id} className="pb-6 border-b border-gray-200 last:border-0">
              <div className="flex items-start gap-4">
                {/* Avatar */}
                <div className="flex-shrink-0 w-10 h-10 bg-gray-300 rounded-full flex items-center justify-center text-gray-600 font-bold">
                  {comment.username ? comment.username.charAt(0).toUpperCase() : '?'}
                </div>
                
                {/* Comment Content */}
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <span className="font-semibold text-gray-900">
                      {comment.username || 'Anonymous'}
                    </span>
                    <span className="text-sm text-gray-500">
                      {formatCommentDate(comment.tanggal_dibuat)}
                    </span>
                  </div>
                  <p className="text-gray-800 leading-relaxed">
                    {comment.konten}
                  </p>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </section>
  );
}