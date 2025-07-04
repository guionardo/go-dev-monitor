<template>
  <q-expansion-item
    class="bx-1"
    expand-separator
    :icon="icon"
    :label="props.origin.origin"
    :header-class="iconColor"
  >
    <q-card bordered :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-grey-2'">
      <q-card-section>
        <div class="row items-center no-wrap">
          <div class="col">
            <div class="text-h6">{{ props.origin.description }}</div>
          </div>

          <div class="col-auto">
            <q-badge v-if="props.origin.dirtyHosts > 0" color="red" title="Dirty hosts">
              {{ props.origin.dirtyHosts }}
              <q-icon name="warning" color="white" /> </q-badge
            >&nbsp;
            <q-badge color="blue" :title="props.origin.language">
              {{ props.origin.language }}<q-icon :name="props.origin.language_icon" color="white" />
            </q-badge>
          </div>
        </div>
      </q-card-section>
      <q-separator />

      <q-card-section>
        <local-repo-card
          v-for="host in props.origin.hosts"
          :key="host.host_name"
          :host="host"
        ></local-repo-card>
      </q-card-section>
    </q-card>
  </q-expansion-item>
</template>
<script setup>
import LocalRepoCard from './LocalRepoCard.vue'
import { computed } from 'vue'
const props = defineProps({
  origin: {
    type: Object,
    required: true,
  },
})

const icon = computed(() => {
  return props.origin.clean ? 'fa-brands fa-git-alt' : 'fa-solid fa-triangle-exclamation'
})
const iconColor = computed(() => {
  return props.origin.clean ? 'text-primary' : 'text-red'
})
</script>
