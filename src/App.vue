<script lang="ts" setup>

onMounted(()=>{
  console.log('onMounted')
  // auto dark
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    toggleDark()
  }
})

</script>


<template>
  <main font-sans p="x-4 y-10" bg="gray-light-500 dark:bg-dark-500" text="center gray-700 dark:gray-200">
    <RouterView />
    <TheFooter />
  </main>
</template>
