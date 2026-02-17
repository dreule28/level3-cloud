<script setup>
async function loginAndCall() {
  // 1) login
  const loginRes = await fetch("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username: "admin", password: "password" }),
  });

  const loginData = await loginRes.json();
  console.log("login:", loginData);

  const token = loginData.access_token; // adjust if your field name differs
  localStorage.setItem("token", token);

  // 2) call protected endpoint with bearer token
  const res = await fetch("/api/instances", {
    headers: { Authorization: `Bearer ${token}` },
  });

  const data = await res.json();
  console.log("instances:", data);
}
</script>

<template>
  <router-view />
</template>
