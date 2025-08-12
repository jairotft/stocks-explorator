<template>
  <main class="min-h-screen bg-slate-50 text-slate-800 p-6">
    <section class="mx-auto w-[90%] space-y-6">
      <router-view />
    </section>
  </main>
</template>

<script setup lang="ts">
// No necesitamos importar componentes específicos ya que router-view los maneja
</script>
