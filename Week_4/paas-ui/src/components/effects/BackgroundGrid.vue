<!--
  BackgroundGrid — Subtle animated grid + radial gradient overlay
  Provides the deep-space cyberpunk atmosphere
-->
<script setup>
import { onMounted, onUnmounted, ref } from "vue";

const canvasRef = ref(null);
let animationId = null;
let particles = [];

function initParticles(w, h) {
  particles = Array.from({ length: 40 }, () => ({
    x: Math.random() * w,
    y: Math.random() * h,
    r: Math.random() * 1.5 + 0.5,
    dx: (Math.random() - 0.5) * 0.3,
    dy: (Math.random() - 0.5) * 0.3,
    alpha: Math.random() * 0.3 + 0.1,
  }));
}

function draw(ctx, w, h) {
  ctx.clearRect(0, 0, w, h);

  // Draw grid
  ctx.strokeStyle = "rgba(59, 130, 246, 0.04)";
  ctx.lineWidth = 1;
  const gridSize = 50;
  for (let x = 0; x < w; x += gridSize) {
    ctx.beginPath();
    ctx.moveTo(x, 0);
    ctx.lineTo(x, h);
    ctx.stroke();
  }
  for (let y = 0; y < h; y += gridSize) {
    ctx.beginPath();
    ctx.moveTo(0, y);
    ctx.lineTo(w, y);
    ctx.stroke();
  }

  // Draw particles
  for (const p of particles) {
    p.x += p.dx;
    p.y += p.dy;
    if (p.x < 0 || p.x > w) p.dx *= -1;
    if (p.y < 0 || p.y > h) p.dy *= -1;

    ctx.beginPath();
    ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2);
    ctx.fillStyle = `rgba(139, 92, 246, ${p.alpha})`;
    ctx.fill();
  }

  animationId = requestAnimationFrame(() => draw(ctx, w, h));
}

onMounted(() => {
  const canvas = canvasRef.value;
  if (!canvas) return;
  const ctx = canvas.getContext("2d");
  const resize = () => {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    initParticles(canvas.width, canvas.height);
  };
  resize();
  window.addEventListener("resize", resize);
  draw(ctx, canvas.width, canvas.height);
});

onUnmounted(() => {
  if (animationId) cancelAnimationFrame(animationId);
});
</script>

<template>
  <div class="pointer-events-none fixed inset-0 z-0">
    <canvas ref="canvasRef" class="absolute inset-0" />
    <!-- Radial gradient overlays for depth -->
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_rgba(59,130,246,0.08)_0%,_transparent_50%)]" />
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_bottom_right,_rgba(139,92,246,0.06)_0%,_transparent_50%)]" />
  </div>
</template>
