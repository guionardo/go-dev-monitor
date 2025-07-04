<template>
  <q-page class="flex flex-center">
    <div class="q-pa-md row-equal-width">
      <div class="row">
        <div class="col">
          <status-card
            title="Total origins"
            :quantity="originCount"
            icon="fa-brands fa-git-alt"
            icon_color="primary"
          />
        </div>
        <div class="col">
          <status-card
            title="Dirty origins"
            :quantity="dirtyCount"
            icon="fa-solid fa-triangle-exclamation"
            icon_color="warning"
          />
        </div>
      </div>
      <div class="row">
        <div class="col-6" v-for="origin in data" :key="origin">
          <origin-card :origin="origin"></origin-card>
        </div>
      </div>
    </div>
    <!-- <img alt="Quasar logo" src="~assets/quasar-logo-vertical.svg" style="width: 200px; height: 200px"> -->
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import StatusCard from '../components/StatusCard.vue'
import OriginCard from '../components/OriginCard.vue'
import useAPI from '../src/data.js'

const { getData } = useAPI()

const data = ref()
const dirtyCount = ref()
const originCount = ref()
onMounted(async () => {
  const d = await getData()
  console.info('Mounted data', d)
  data.value = d.origins
  dirtyCount.value = d.dirtyOrigins
  originCount.value = d.countOrigins
})
</script>
<style lang="sass">
.row-equal-width
  .row > div
    padding: 10px 15px
    background: rgba(#999,.15)
    border: 0px solid rgba(#999,.2)
  .row + .row
    margin-top: 1rem
</style>
