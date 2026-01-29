<script setup lang="ts">
import { computed } from 'vue'

interface GradientTextProps {
  text: string
  className?: string
  colors?: string[]
  animationSpeed?: number
  showBorder?: boolean
}

const props = withDefaults(defineProps<GradientTextProps>(), {
  text: '',
  className: '',
  colors: () => ['#ffaa40', '#9c40ff', '#ffaa40'],
  animationSpeed: 8,
  showBorder: false,
})

const gradientStyle = computed(() => ({
  backgroundImage: `linear-gradient(to right, ${props.colors.join(', ')})`,
  animationDuration: `${props.animationSpeed}s`,
  backgroundSize: '300% 100%',
  '--animation-duration': `${props.animationSpeed}s`,
}))

const borderStyle = computed(() => ({
  ...gradientStyle.value,
}))

const textStyle = computed(() => ({
  ...gradientStyle.value,
  backgroundClip: 'text',
  WebkitBackgroundClip: 'text',
}))
</script>

<template>
  <div
    :class="`relative mx-auto flex max-w-fit flex-row items-center justify-center font-medium backdrop-blur transition-shadow duration-500 overflow-hidden cursor-pointer ${className}`"
  >
    <div
      v-if="showBorder"
      class="absolute inset-0 bg-cover z-0 pointer-events-none animate-gradient"
      :style="borderStyle"
    >
      <div
        class="absolute inset-0 bg-black z-[-1]"
        style="
          width: calc(100% - 2px);
          height: calc(100% - 2px);
          left: 50%;
          top: 50%;
          transform: translate(-50%, -50%);
        "
      />
    </div>

    <div
      class="inline-block relative z-2 text-transparent bg-cover animate-gradient"
      :style="textStyle"
    >
      <slot></slot>
    </div>
  </div>
</template>

<style scoped>
@keyframes gradient {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.animate-gradient {
  animation: gradient var(--animation-duration, 8s) linear infinite;
}
</style>
