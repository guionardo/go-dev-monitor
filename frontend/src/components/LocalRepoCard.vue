<template>
  <q-card bordered class="my-card" :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-grey-2'">
    <q-card-section>
      <div class="row items-center no-wrap">
        <div class="col">
          <div class="text-h7">
            <q-icon name="fas fa-server"></q-icon> {{ props.host.host_name || 'unknown' }}
          </div>
          <div class="text-subtitle2">{{ props.host.foldername }}</div>
        </div>
        <q-badge color="green" v-if="clean">
          clean <q-icon name="check" color="white" class="q-ml-xs" />
        </q-badge>
        <q-badge color="red" v-else>
          dirty <q-icon name="warning" color="white" class="q-ml-xs" /> </q-badge
        >&nbsp;
        <q-badge color="primary" :title="fetchTime">
          {{ fetchTimeDisplay }}
        </q-badge>
      </div>
    </q-card-section>
    <q-separator />

    <q-card-section>
      <q-icon name="fas fa-code-branch"></q-icon>
      {{ props.host.current_branch }}
    </q-card-section>
    <q-separator />
    <q-card-section>
      <q-icon name="folder"></q-icon>
      {{ props.host.folder_name }}
    </q-card-section>

    <q-separator />
    <!-- commit -->
    <q-card-section>
      <div class="row items-center no-wrap">
        <div class="col-auto">
          <div class="text-h7">
            <q-icon name="commit"></q-icon>
          </div>
        </div>
        <div class="col">
          {{ props.host.last_commit.message }}
          <div class="text-subtitle2">{{ props.host.last_commit.author }}</div>
        </div>
        <q-badge color="teal"
          ><a :href="props.host.last_commit.url" target="_blank">#{{ commitHash }}</a></q-badge
        >&nbsp;
        <q-badge color="secondary" :title="commitTime">{{ commitTimeDisplay }}</q-badge>
      </div>
    </q-card-section>
    <q-separator />

    <files-card
      title="Untracked files"
      :files="props.host.untracked_files"
      icon="fa fa-question-circle"
    />
    <files-card title="Changed files" :files="props.host.changed_files" icon="fa fa-edit" />
    <files-card
      title="Last changed files"
      :files="props.host.last_changed_files"
      icon="fa fa-clock"
    />
  </q-card>
</template>
<script setup>
import moment from 'moment'

import FilesCard from './FilesCard.vue'
import { computed } from 'vue'
const props = defineProps({
  host: {
    type: Object,
    required: true,
  },
})

const clean = computed(() => props.host.clean)
const fetchTimeDisplay = computed(() => moment(new Date(props.host.fetch_time)).fromNow())
const fetchTime = computed(() => {
  return new Date(props.host.fetch_time).toLocaleString()
})
const commitHash = computed(() => props.host.last_commit.hash.slice(0, 7))
const commitTime = computed(() => new Date(props.host.last_commit.when).toLocaleString())
const commitTimeDisplay = computed(() => moment(new Date(props.host.last_commit.when)).fromNow())
</script>
<style lang="css" scoped>
.commit > span {
  left: 3em;
}
</style>
