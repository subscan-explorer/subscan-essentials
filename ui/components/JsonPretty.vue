<template>
  <div class="json-pretty-component">
    <vue-json-pretty
      :data="data"
      :deep="deep"
      :showLine="showLine"
      :showDoubleQuotes="showDoubleQuotes"
      @click="handleClick">
    </vue-json-pretty>
  </div>
</template>
<script>
import VueJsonPretty from 'vue-json-pretty'
export default {
  props: {
    data: {
      type: [Object, Array, String],
    },
    deep: {
      type: Number,
      default: 4
    },
    showLine: {
      type: Boolean,
      default: false
    },
    showDoubleQuotes: {
      type: Boolean,
      default: false
    }
  },
  components: {
    VueJsonPretty
  },
  methods: {
    handleClick() {

    }
  }
};
</script>

<style lang="scss" scoped>
.json-pretty-component {
  /deep/ .vjs-tree {
    font-family: initial;
    color: #333;
    .vjs-value__string {
      color: #333;
      word-break: break-all;
    }
    .vjs-key {
      color: var(--main-color);
    }
    .vjs-value__number {
      color: #333;
    }
  }
}
</style>
