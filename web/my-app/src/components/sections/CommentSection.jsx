'use client';

import { useState } from "react";

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
  const [newComment, setNewComment] = useState({
    nama_pengguna: '',
    konten: ''
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!newComment.konten.trim() || !newComment.nama_pengguna.trim()) {
      alert('Nama dan komentar harus diisi');
      return;
    }

    setIsSubmitting(true);

    try {
      // TODO: Replace with actual API call
      const res = await fetch(`http://localhost:8080/api/v1/articles/${articleId}/comments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          nama_pengguna: newComment.nama_pengguna,
          konten: newComment.konten
        })
      });

      if (!res.ok) throw new Error('Failed to post comment');

      const savedComment = await res.json();
      
      // Add new comment to list
      setComments([savedComment, ...comments]);
      
      // Reset form
      setNewComment({ nama_pengguna: '', konten: '' });
      
    } catch (error) {
      console.error('Error posting comment:', error);
      
      // Fallback: Add dummy comment locally for development
      const dummyComment = {
        comment_id: Date.now(),
        artikel_id: articleId,
        user_id: null,
        nama_pengguna: newComment.nama_pengguna,
        konten: newComment.konten,
        tanggal_dibuat: new Date().toISOString()
      };
      
      setComments([dummyComment, ...comments]);
      setNewComment({ nama_pengguna: '', konten: '' });
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
      <form onSubmit={handleSubmit} className="mb-8 p-6 bg-gray-50 rounded-lg">
        <h3 className="text-lg font-semibold mb-4">Tinggalkan Komentar</h3>
        
        <div className="mb-4">
          <label htmlFor="nama_pengguna" className="block text-sm font-medium text-gray-700 mb-2">
            Nama
          </label>
          <input
            type="text"
            id="nama_pengguna"
            value={newComment.nama_pengguna}
            onChange={(e) => setNewComment({ ...newComment, nama_pengguna: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Masukkan nama Anda"
            required
          />
        </div>

        <div className="mb-4">
          <label htmlFor="konten" className="block text-sm font-medium text-gray-700 mb-2">
            Komentar
          </label>
          <textarea
            id="konten"
            value={newComment.konten}
            onChange={(e) => setNewComment({ ...newComment, konten: e.target.value })}
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

      {/* Comments List */}
      <div className="space-y-6">
        {comments.length === 0 ? (
          <p className="text-gray-500 text-center py-8">
            Belum ada komentar. Jadilah yang pertama berkomentar!
          </p>
        ) : (
          comments.map((comment) => (
            <div key={comment.comment_id} className="pb-6 border-b border-gray-200 last:border-0">
              <div className="flex items-start gap-4">
                {/* Avatar */}
                <div className="flex-shrink-0 w-10 h-10 bg-gray-300 rounded-full flex items-center justify-center text-gray-600 font-bold">
                  {comment.nama_pengguna.charAt(0).toUpperCase()}
                </div>
                
                {/* Comment Content */}
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <span className="font-semibold text-gray-900">
                      {comment.nama_pengguna}
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