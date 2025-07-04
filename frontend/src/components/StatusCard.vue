<template>
  <q-card class="my-card" :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-grey-2'">
    <q-card-section>
      <div class="row items-center no-wrap">
        <div class="col-auto">
          <div class="text-h2"><q-icon :class="props.icon" :color="props.icon_color" /></div>
        </div>
        <div class="col" style="min-height: 84.6px">
          <div class="text-h6">{{ props.quantity }}</div>
          <div class="text-subtitle2">{{ props.title }}</div>
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  icon: {
    type: String,
    required: true,
  },
  quantity: {
    type: Number,
    required: true,
  },
  icon_color: {
    type: String,
    default: '',
  },
})
</script>
<style>
.my-card {
  min-width: 380px;
}
</style>
