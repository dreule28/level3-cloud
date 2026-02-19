<!--
  BackgroundGrid — Animated cyberpunk grid + particle network + scanlines
  Full-canvas atmosphere with connected particle network and pulsing grid nodes
-->
<script setup>
import { onMounted, onUnmounted, ref } from "vue";

const canvasRef = ref(null);
let animationId = null;
let particles = [];
let mouse = { x: -1000, y: -1000 };
let time = 0;

function initParticles(w, h) {
  particles = Array.from({ length: 60 }, () => ({
    x: Math.random() * w,
    y: Math.random() * h,
    r: Math.random() * 2 + 0.5,
    dx: (Math.random() - 0.5) * 0.4,
    dy: (Math.random() - 0.5) * 0.4,
    alpha: Math.random() * 0.4 + 0.1,
    hue: Math.random() > 0.5 ? 220 : 270, // blue or purple
  }));
}

function draw(ctx, w, h) {
  ctx.clearRect(0, 0, w, h);
  time += 0.005;

  // ── Animated Grid ──
  const gridSize = 60;
  ctx.lineWidth = 0.5;
  for (let x = 0; x < w; x += gridSize) {
    const pulse = Math.sin(time + x * 0.01) * 0.02 + 0.03;
    ctx.strokeStyle = `rgba(59, 130, 246, ${pulse})`;
    ctx.beginPath();
    ctx.moveTo(x, 0);
    ctx.lineTo(x, h);
    ctx.stroke();
  }
  for (let y = 0; y < h; y += gridSize) {
    const pulse = Math.sin(time + y * 0.01) * 0.02 + 0.03;
    ctx.strokeStyle = `rgba(59, 130, 246, ${pulse})`;
    ctx.beginPath();
    ctx.moveTo(0, y);
    ctx.lineTo(w, y);
    ctx.stroke();
  }

  // ── Pulsing Grid Intersections ──
  for (let x = 0; x < w; x += gridSize) {
    for (let y = 0; y < h; y += gridSize) {
      const dist = Math.hypot(x - mouse.x, y - mouse.y);
      const prox = Math.max(0, 1 - dist / 250);
      const pulse = (Math.sin(time * 2 + x * 0.02 + y * 0.02) + 1) * 0.5;
      const alpha = 0.05 + pulse * 0.1 + prox * 0.4;
      const radius = 1 + prox * 3;
      ctx.beginPath();
      ctx.arc(x, y, radius, 0, Math.PI * 2);
      ctx.fillStyle = `rgba(59, 130, 246, ${alpha})`;
      ctx.fill();
    }
  }

  // ── Particle Network ──
  for (const p of particles) {
    p.x += p.dx;
    p.y += p.dy;
    if (p.x < 0 || p.x > w) p.dx *= -1;
    if (p.y < 0 || p.y > h) p.dy *= -1;

    // Mouse repulsion
    const md = Math.hypot(p.x - mouse.x, p.y - mouse.y);
    if (md < 150) {
      p.x += (p.x - mouse.x) * 0.01;
      p.y += (p.y - mouse.y) * 0.01;
    }

    ctx.beginPath();
    ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2);
    ctx.fillStyle = `hsla(${p.hue}, 80%, 65%, ${p.alpha})`;
    ctx.fill();
  }

  // ── Connection Lines ──
  const maxDist = 120;
  for (let i = 0; i < particles.length; i++) {
    for (let j = i + 1; j < particles.length; j++) {
      const d = Math.hypot(particles[i].x - particles[j].x, particles[i].y - particles[j].y);
      if (d < maxDist) {
        const alpha = (1 - d / maxDist) * 0.15;
        ctx.beginPath();
        ctx.moveTo(particles[i].x, particles[i].y);
        ctx.lineTo(particles[j].x, particles[j].y);
        ctx.strokeStyle = `rgba(139, 92, 246, ${alpha})`;
        ctx.lineWidth = 0.5;
        ctx.stroke();
      }
    }
  }

  // ── Scanline Effect ──
  const scanY = (time * 200) % h;
  const grad = ctx.createLinearGradient(0, scanY - 30, 0, scanY + 30);
  grad.addColorStop(0, "rgba(59, 130, 246, 0)");
  grad.addColorStop(0.5, "rgba(59, 130, 246, 0.03)");
  grad.addColorStop(1, "rgba(59, 130, 246, 0)");
  ctx.fillStyle = grad;
  ctx.fillRect(0, scanY - 30, w, 60);

  animationId = requestAnimationFrame(() => draw(ctx, w, h));
}

function onMouseMove(e) {
  mouse.x = e.clientX;
  mouse.y = e.clientY;
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
  window.addEventListener("mousemove", onMouseMove);
  draw(ctx, canvas.width, canvas.height);
});

onUnmounted(() => {
  if (animationId) cancelAnimationFrame(animationId);
  window.removeEventListener("mousemove", onMouseMove);
});
</script>

<template>
  <div class="pointer-events-none fixed inset-0 z-0">
    <canvas ref="canvasRef" class="absolute inset-0" />
    <!-- Radial gradient overlays for depth -->
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_rgba(59,130,246,0.1)_0%,_transparent_50%)]" />
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_bottom_right,_rgba(139,92,246,0.08)_0%,_transparent_50%)]" />
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_left,_rgba(6,182,212,0.05)_0%,_transparent_40%)]" />
    <!-- Subtle CRT scanline overlay -->
    <div class="absolute inset-0 opacity-[0.015]" style="background: repeating-linear-gradient(0deg, transparent, transparent 2px, rgba(255,255,255,0.03) 2px, rgba(255,255,255,0.03) 4px)" />
  </div>
</template>
