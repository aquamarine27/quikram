import { useEffect, useRef, useState } from "react";

export function useParallax(speed = 0.3, maxOffset = 30) {
  const ref = useRef<HTMLDivElement>(null);
  const [offset, setOffset] = useState(0);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;

    let ticking = false;

    const handleScroll = () => {
      if (!ticking) {
        window.requestAnimationFrame(() => {
          const rect = el.getBoundingClientRect();
          const center = rect.top + rect.height / 2;
          const viewCenter = window.innerHeight / 2;
          const raw = (center - viewCenter) * speed;
          setOffset(Math.max(-maxOffset, Math.min(maxOffset, raw)));
          ticking = false;
        });
        ticking = true;
      }
    };

    window.addEventListener("scroll", handleScroll, { passive: true });
    handleScroll();
    return () => window.removeEventListener("scroll", handleScroll);
  }, [speed, maxOffset]);

  return { ref, style: { transform: `translateY(${offset}px)` } };
}
