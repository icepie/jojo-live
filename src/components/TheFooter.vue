<template>
  <nav text-xl mt-6 inline-flex gap-2>
    <button icon-btn @click="toggleDark()">
      <div dark:i-carbon-moon i-carbon-sun />
    </button>

    <a
      icon-btn i-carbon-logo-github
      rel="noreferrer"
      href="https://github.com/icepie/jojo-live"
      target="_blank"
      title="GitHub"
    />


  </nav>

  <nav>
  ©️ powered by <a text-green-5 href="https://github.com/icepie">icepie</a>
  </nav>
</template>
