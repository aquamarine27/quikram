import { useEffect, useRef } from "react";

function isMobile() {
  return window.innerWidth < 768;
}

export function useSmoothSnap() {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const container = containerRef.current;
    if (!container || isMobile()) return;

    let isScrolling = false;
    let timeout: number | null = null;
    let touchStartY = 0;
    let touchDir = 0;

    const scrollToSection = (dir: number) => {
      if (isScrolling) return;
      const sections = container.querySelectorAll<HTMLElement>(".home-section, .footer");
      const currentScroll = container.scrollTop;
      let target: HTMLElement | null = null;

      if (dir > 0) {
        for (const s of sections) {
          if (s.offsetTop > currentScroll + 40) { target = s; break; }
        }
        if (!target) target = sections[sections.length - 1];
      } else {
        let prev: HTMLElement | null = null;
        for (const s of sections) {
          if (s.offsetTop >= currentScroll - 40) break;
          prev = s;
        }
        target = prev || sections[0];
      }

      if (target) {
        isScrolling = true;
        container.scrollTo({ top: target.offsetTop, behavior: "smooth" });
        if (timeout) clearTimeout(timeout);
        timeout = window.setTimeout(() => { isScrolling = false; }, 900);
      }
    };

    const onWheel = (e: WheelEvent) => {
      e.preventDefault();
      scrollToSection(e.deltaY > 0 ? 1 : -1);
    };

    const onTouchStart = (e: TouchEvent) => {
      touchStartY = e.touches[0].clientY;
      touchDir = 0;
    };

    const onTouchMove = (e: TouchEvent) => {
      if (touchDir !== 0) return;
      const dy = e.touches[0].clientY - touchStartY;
      if (Math.abs(dy) > 50) {
        touchDir = dy < 0 ? 1 : -1;
        e.preventDefault();
        scrollToSection(touchDir);
      }
    };

    container.addEventListener("wheel", onWheel, { passive: false });
    container.addEventListener("touchstart", onTouchStart, { passive: true });
    container.addEventListener("touchmove", onTouchMove, { passive: false });

    return () => {
      container.removeEventListener("wheel", onWheel);
      container.removeEventListener("touchstart", onTouchStart);
      container.removeEventListener("touchmove", onTouchMove);
      if (timeout) clearTimeout(timeout);
    };
  }, []);

  return containerRef;
}
