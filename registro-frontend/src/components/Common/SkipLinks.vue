<template>
  <nav :aria-label="t('skipLinks.navAria') || 'Scorciatoie di navigazione rapida'" class="skip-links-nav">
    <a
      href="#main-content"
      class="skip-link"
      @click.prevent="jumpTo('#main-content')"
    >
      {{ t('skipLinks.mainContent') || 'Salta al contenuto principale' }}
    </a>
    <a
      href="#main-nav"
      class="skip-link"
      @click.prevent="jumpTo('#main-nav')"
    >
      {{ t('skipLinks.mainNav') || 'Salta al menu di navigazione' }}
    </a>
    <a
      href="#a11y-panel"
      class="skip-link"
      @click.prevent="jumpTo('#a11y-panel')"
    >
      {{ t('skipLinks.a11yPanel') || 'Apri pannello accessibilità' }}
    </a>
  </nav>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

function jumpTo(targetId) {
  const el = document.querySelector(targetId)
  if (el) {
    el.setAttribute('tabindex', '-1')
    el.focus()
    el.scrollIntoView({ behavior: 'smooth' })
  }
}
</script>

<style scoped>
.skip-links-nav {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 9999;
}

.skip-link {
  position: absolute;
  top: -100px;
  left: 12px;
  background: #1976D2;
  color: #ffffff;
  padding: 10px 16px;
  font-weight: 700;
  border-radius: 0 0 8px 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
  transition: top 0.2s ease-in-out;
  text-decoration: none;
  font-size: 0.95rem;
}

.skip-link:focus {
  top: 0;
  outline: 3px solid #f59e0b;
}
</style>
