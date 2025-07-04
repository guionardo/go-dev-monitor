<template>
  <q-expansion-item expand-separator :icon="props.icon" :label="title" switch-toggle-side>
    <q-list bordered padding class="rounded-borders">
      <q-item v-for="file in props.files" :key="file.name">
        <q-item-section>
          <q-item-label>{{ file.name }}</q-item-label>
          <q-item-label caption>{{ file.time }}</q-item-label>
        </q-item-section>
        <q-item-section side>
          {{ file.status }}
        </q-item-section>
      </q-item>
    </q-list>
  </q-expansion-item>
</template>
<script setup>
import { computed } from 'vue'
const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  files: {
    type: Array,
    required: true,
  },
  icon: {
    type: String,
    required: true,
  },
})
const title = computed(() => `${props.title} (${props.files.length || 'zero'})`)
</script>
