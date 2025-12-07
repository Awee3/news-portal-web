'use client';
import { useState, useEffect } from 'react';

function formatIDDate(d) {
  const days = ['MINGGU','SENIN','SELASA','RABU','KAMIS','JUMAT','SABTU'];
  const months = [
    'JANUARI','FEBRUARI','MARET','APRIL','MEI','JUNI',
    'JULI','AGUSTUS','SEPTEMBER','OKTOBER','NOVEMBER','DESEMBER'
  ];
  return `${days[d.getDay()]}, ${String(d.getDate()).padStart(2,'0')} ${months[d.getMonth()]} ${d.getFullYear()}`;
}

export default function CurrentDate() {
  const [text, setText] = useState('');

  useEffect(() => {
    function update() {
      setText(formatIDDate(new Date()));
    }
    update();

    // Hitung waktu menuju pergantian hari untuk refresh otomatis
    const now = new Date();
    const msToMidnight =
      new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1).getTime() - now.getTime();
    const timeout = setTimeout(() => {
      update();
      // Setelah pergantian hari, update tiap 24 jam
      setInterval(update, 24 * 60 * 60 * 1000);
    }, msToMidnight);

    return () => {
      clearTimeout(timeout);
    };
  }, []);

  return <div className="text-gray-600">{text || '...'}</div>;
}